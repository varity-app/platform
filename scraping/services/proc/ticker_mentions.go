package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/protobuf/proto"

	"github.com/VarityPlatform/scraping/common"
	commonPB "github.com/VarityPlatform/scraping/protobuf/common"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"go.opentelemetry.io/otel/trace"
)

type kafkaExtractorHandler func(producer *kafka.Producer, readChan chan *kafka.Message, deliveryChan chan kafka.Event, mu *sync.Mutex, readErr *error, readWG *sync.WaitGroup, deliveryWG *sync.WaitGroup, count *int, allTickers []common.IEXTicker)

// Convert a kafka message into a ticker mention
func kafkaToTickerMention(msg *kafka.Message) (interface{}, error) {
	mention := &commonPB.TickerMention{}
	if err := proto.Unmarshal(msg.Value, mention); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	return mention, nil
}

// General purpose method for extracting tickers from a kafka stream
func extractTickersKafka(ctx context.Context, fsClient *firestore.Client, producer *kafka.Producer, consumer *kafka.Consumer, allTickers []common.IEXTicker, tracer *trace.Tracer, topic string, handler kafkaExtractorHandler) (int, error) {
	var mu sync.Mutex
	count := 0

	// Subscribe to topic
	err := consumer.Unassign()
	if err != nil {
		return count, fmt.Errorf("kafka.Unassign: %v", err)
	}
	_, span := (*tracer).Start(ctx, "firestore.GetOffsets")
	checkpointKey := topic + "-proc"
	offsets, err := getOffsetsFromFS(ctx, fsClient, checkpointKey)
	if err != nil {
		return count, err
	}
	partitions := offsetMapToPartitions(topic, offsets)
	err = consumer.Assign(partitions)
	if err != nil {
		return count, fmt.Errorf("kafka.SubscribeTopic: %v", err)
	}
	span.End()

	// Delivery report handler for produced messages
	deliveryChan := make(chan kafka.Event)
	deliveryWG := new(sync.WaitGroup)
	var writeErr error
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					mu.Lock()
					writeErr = fmt.Errorf("kafka.ProduceMessage: %v", ev.TopicPartition)
					mu.Unlock()
				}
				deliveryWG.Done()
			}
		}
	}()

	// Handle read messages from kafka
	readChan := make(chan *kafka.Message)
	readWG := new(sync.WaitGroup)
	var readErr error
	go handler(producer, readChan, deliveryChan, &mu, &readErr, readWG, deliveryWG, &count, allTickers)
	span.End()

	// Consume messages from Kafka
	_, span = (*tracer).Start(ctx, "kafka.ReadMessage")
	for {
		msg, err := consumer.ReadMessage(CONSUMER_TIMEOUT)
		if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
			break
		} else if err != nil {
			return count, fmt.Errorf("kafka.ReadMessage: %v", err)
		}

		offset, err := strconv.Atoi(msg.TopicPartition.Offset.String())
		if err != nil {
			return count, fmt.Errorf("strconv.Atoi: %v", err)
		}
		offsets[fmt.Sprint(msg.TopicPartition.Partition)] = offset + 1

		// Check for errors occurring in goroutines
		if readErr != nil {
			return count, readErr
		} else if writeErr != nil {
			return count, writeErr
		}

		readWG.Add(1)
		readChan <- msg

	}

	span.End()
	_, span = (*tracer).Start(ctx, "kafka.Produce")

	producer.Flush(PRODUCER_TIMEOUT_MS)
	readWG.Wait()
	deliveryWG.Wait()
	span.End()

	// Check for errors occurring in goroutines
	if readErr != nil {
		return count, readErr
	} else if writeErr != nil {
		return count, writeErr
	}

	// Save offsets to firestore
	_, span = (*tracer).Start(ctx, "firestore.SaveOffsets")
	if err := saveOffsetsToFS(ctx, fsClient, checkpointKey, offsets); err != nil {
		return count, err
	}
	span.End()

	return count, nil
}
