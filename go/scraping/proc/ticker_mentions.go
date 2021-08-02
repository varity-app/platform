package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"google.golang.org/protobuf/proto"

	"github.com/VarityPlatform/scraping/common"
	commonPB "github.com/VarityPlatform/scraping/protobuf/common"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/spf13/viper"

	"go.opentelemetry.io/otel/trace"
)

type kafkaExtractorHandler func(producer *kafka.Producer, readChan chan *kafka.Message, deliveryChan chan kafka.Event, mu *sync.Mutex, readErr *error, readWG *sync.WaitGroup, deliveryWG *sync.WaitGroup, count *int, allTickers []common.IEXTicker)

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

// Read ticker mentions from Pub/Sub and process them.
func readTickerMentions(ctx context.Context, psClient *pubsub.Client, bqClient *bigquery.Client, tracer *trace.Tracer) (int, error) {

	sub := psClient.Subscription(common.SUBSCRIPTION_TICKER_MENTIONS + "-" + viper.GetString("deploymentMode"))
	topic := psClient.Topic(common.TOPIC_TICKER_MENTIONS + "-" + viper.GetString("deploymentMode"))

	cm := make(chan *pubsub.Message)
	defer close(cm)

	readWG := new(sync.WaitGroup)
	writeWG := new(sync.WaitGroup)
	var mu sync.Mutex

	// Count messages processed
	count := 0
	totalCount := 0

	// Create bigquery streaming insert
	inserter := bqClient.Dataset(
		common.BIGQUERY_DATASET_SCRAPING + "_" + viper.GetString("deploymentMode"),
	).Table(common.BIGQUERY_TABLE_TICKER_MENTIONS).Inserter()

	// Async errors
	var readErr error = nil
	var writeErr error = nil

	// Spawn thread for parsing protobuf messages
	mentions := []*commonPB.TickerMention{}
	go func() {
		for msg := range cm {
			mention := &commonPB.TickerMention{}
			if err := proto.Unmarshal(msg.Data, mention); err != nil {
				log.Println(err)
				readErr = fmt.Errorf("tickerMention.Unmarshal: %v", err)
				readWG.Done()
				return
			}
			mu.Lock()
			mentions = append(mentions, mention)
			mu.Unlock()

			msg.Ack()
			readWG.Done()
		}
	}()

	// Receive messages for N sec
	for {
		pubCtx, cancel := context.WithTimeout(ctx, WAIT_INTERVAL)
		defer cancel()

		// Reset count
		count = 0
		mentions = []*commonPB.TickerMention{}

		ctx, span := (*tracer).Start(ctx, "pubsub.Receive")
		// Receive blocks until the context is cancelled or an error occurs.
		err := sub.Receive(pubCtx, func(ctx context.Context, msg *pubsub.Message) {
			if readErr != nil {
				cancel()
				return
			}
			readWG.Add(1)
			cm <- msg
			count++
		})
		if err != nil {
			log.Println(err)
			return totalCount, err
		}

		readWG.Wait()
		if readErr != nil {
			return totalCount, readErr
		}
		span.End()

		ctx, span = (*tracer).Start(ctx, "bigquery.Upload")
		// Upload to bigquery
		if err := inserter.Put(ctx, mentions); err != nil {
			// Republish mentions to pubsub if there's an error
			for _, mention := range mentions {

				serializedMention, serErr := proto.Marshal(mention)
				if serErr != nil {
					return totalCount, fmt.Errorf("mention.SerializeProto: %v AND bigquery.Insert: %v", serErr, err)
				}

				// Asynchronously publish to pubsub
				result := topic.Publish(ctx, &pubsub.Message{
					Data: serializedMention,
				})

				writeWG.Add(1)
				go func(res *pubsub.PublishResult) {
					defer writeWG.Done()
					_, err := result.Get(ctx)
					if err != nil {
						mu.Lock()
						writeErr = fmt.Errorf("mention.Publish: %v", err)
						mu.Unlock()
					}
					writeWG.Done()
				}(result)
			}

			writeWG.Wait()

			// Handle
			if writeErr != nil {
				return totalCount, fmt.Errorf("pubsub.PublishResult: %v AND bigquery.Insert: %v", writeErr, err)
			}
			return totalCount, fmt.Errorf("bigquery.Insert: %v", err)
		}
		span.End()

		totalCount += count

		// If no messages were received in this window, break
		if count == 0 {
			break
		}
	}

	return totalCount, nil

}
