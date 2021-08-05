package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/VarityPlatform/scraping/common"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaTickerProcessorHandler is a function type that handles kafka messages and sends ticker mentions on another topic.
type KafkaTickerProcessorHandler func(producer *kafka.Producer, readChan chan *kafka.Message, deliveryChan chan kafka.Event, mu *sync.Mutex, readErr *error, readWG *sync.WaitGroup, deliveryWG *sync.WaitGroup, count *int, allTickers []common.IEXTicker)

// KafkaTickerProcessor reads a Kafka topic and parses the messages for tickers.
type KafkaTickerProcessor struct {
	allTickers    []common.IEXTicker
	consumer      *kafka.Consumer
	producer      *kafka.Producer
	offsetManager *OffsetManager
}

// NewKafkaTickerProcessor initializes a new KafkaTickerProcessor
func NewKafkaTickerProcessor(ctx context.Context, opts OffsetManagerOpts, allTickers []common.IEXTicker) (*KafkaTickerProcessor, error) {
	// Initialize kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     os.Getenv("KAFKA_AUTH_KEY"),
		"sasl.password":     os.Getenv("KAFKA_AUTH_SECRET"),
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.GetProducer: %v", err)
	}

	// Initialize kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol":        "SASL_SSL",
		"sasl.mechanisms":          "PLAIN",
		"sasl.username":            os.Getenv("KAFKA_AUTH_KEY"),
		"sasl.password":            os.Getenv("KAFKA_AUTH_SECRET"),
		"group.id":                 "proc",
		"enable.auto.offset.store": "false",
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.NewConsumer: %v", err)
	}

	offsetManager, err := NewOffsetManager(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return pointer
	return &KafkaTickerProcessor{
		allTickers:    allTickers,
		consumer:      consumer,
		producer:      producer,
		offsetManager: offsetManager,
	}, nil
}

// Close closes the ticker processor's connection
func (processor *KafkaTickerProcessor) Close() error {
	processor.producer.Close()
	if err := processor.consumer.Close(); err != nil {
		return err
	}

	return processor.offsetManager.Close()
}

// ProcessTopic processes a kafka topic for ticker mentions
func (processor *KafkaTickerProcessor) ProcessTopic(ctx context.Context, topic string, handler KafkaTickerProcessorHandler) (int, error) {
	var mu sync.Mutex
	count := 0

	// Subscribe to topic
	err := processor.consumer.Unassign()
	if err != nil {
		return count, fmt.Errorf("kafka.Unassign: %v", err)
	}
	checkpointKey := topic + "-proc"
	partitions, offsets, err := processor.offsetManager.Fetch(ctx, topic, checkpointKey)
	if err != nil {
		return count, fmt.Errorf("kafka.FetchOffsets: %v", err)
	}
	err = processor.consumer.Assign(partitions)
	if err != nil {
		return count, fmt.Errorf("kafka.SubscribeTopic: %v", err)
	}

	// Delivery report handler for produced messages
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
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
	defer close(readChan)
	readWG := new(sync.WaitGroup)
	var readErr error
	go handler(processor.producer, readChan, deliveryChan, &mu, &readErr, readWG, deliveryWG, &count, processor.allTickers)

	// Consume messages from Kafka
	for {
		msg, err := processor.consumer.ReadMessage(ConsumerTimeout)
		if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
			break
		} else if err != nil {
			return count, fmt.Errorf("kafka.ReadMessage: %v", err)
		}

		// Update offsets
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

	// Wait for processes to finish
	processor.producer.Flush(ProducerTimeoutMS)
	readWG.Wait()
	deliveryWG.Wait()

	// Check for errors occurring in goroutines
	if readErr != nil {
		return count, readErr
	} else if writeErr != nil {
		return count, writeErr
	}

	// Save offsets to firestore
	if err := processor.offsetManager.Save(ctx, checkpointKey, offsets); err != nil {
		return count, err
	}

	return count, nil
}
