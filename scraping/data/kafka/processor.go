package kafka

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProcessorHandler is an interface used by a Processor to process individual kafka messages
// and return a list of messages to send to subsequent kafka topic
type ProcessorHandler interface {
	Process(*kafka.Message) ([][]byte, error)
}

// Processor reads a Kafka topic and parses the messages for tickers.
type Processor struct {
	consumer      *kafka.Consumer
	producer      *kafka.Producer
	offsetManager *OffsetManager
}

// NewProcessor initializes a new Processor
func NewProcessor(ctx context.Context, offsetManager *OffsetManager, opts Opts) (*Processor, error) {
	// Initialize kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": opts.BootstrapServers,
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     opts.Username,
		"sasl.password":     opts.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.NewProducer: %v", err)
	}

	// Initialize kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        opts.BootstrapServers,
		"security.protocol":        "SASL_SSL",
		"sasl.mechanisms":          "PLAIN",
		"sasl.username":            opts.Username,
		"sasl.password":            opts.Password,
		"group.id":                 "proc",
		"enable.auto.offset.store": "false",
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.NewConsumer: %v", err)
	}

	// Return pointer
	return &Processor{
		consumer:      consumer,
		producer:      producer,
		offsetManager: offsetManager,
	}, nil
}

// Close closes the ticker processor's connection
func (p *Processor) Close() error {
	p.producer.Close()
	if err := p.consumer.Close(); err != nil {
		return err
	}

	return p.offsetManager.Close()
}

// ProcessTopic processes a kafka topic for ticker mentions
func (p *Processor) ProcessTopic(ctx context.Context, topic string, checkpointKey string, handler ProcessorHandler) (int, error) {
	var mu sync.Mutex
	count := 0

	// Subscribe to topic
	offsets, err := p.subscribe(ctx, topic, checkpointKey)
	if err != nil {
		return count, err
	}

	// Delivery report handler for produced messages
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	deliveryWG := new(sync.WaitGroup)
	var writeErr error
	go p.handleDeliveries(deliveryChan, deliveryWG, &mu, &writeErr)

	// Handle read messages from kafka
	readChan := make(chan *kafka.Message)
	defer close(readChan)
	readWG := new(sync.WaitGroup)
	var readErr error
	go func() {
		for msg := range readChan {
			mu.Lock()
			count++
			mu.Unlock()
			readWG.Done()

			newMsgs, err := handler.Process(msg)
			if err != nil {
				readErr = fmt.Errorf("handler.Process: %v", err)
			}

			for _, newMsg := range newMsgs {

				// Write to ticker mentions topic
				deliveryWG.Add(1)
				err = p.producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          newMsg,
				}, deliveryChan)
				if err != nil {
					readErr = fmt.Errorf("kafka.Produce: %v", err)
				}

			}
		}
	}()

	// Consume messages from Kafka
	for {
		msg, err := p.consumer.ReadMessage(ConsumerTimeout)
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
	p.producer.Flush(ProducerTimeoutMS)
	readWG.Wait()
	deliveryWG.Wait()

	// Check for errors occurring in goroutines
	if readErr != nil {
		return count, readErr
	} else if writeErr != nil {
		return count, writeErr
	}

	// Save offsets to firestore
	if err := p.offsetManager.Save(ctx, checkpointKey, offsets); err != nil {
		return count, err
	}

	return count, nil
}

// subscribe fetches offsets with the offsetManager and subscribes the consumer
func (p *Processor) subscribe(ctx context.Context, topic string, checkpointKey string) (map[string]int, error) {
	err := p.consumer.Unassign()
	if err != nil {
		return nil, fmt.Errorf("kafka.Unassign: %v", err)
	}
	partitions, offsets, err := p.offsetManager.Fetch(ctx, topic, checkpointKey)
	if err != nil {
		return nil, fmt.Errorf("kafka.FetchOffsets: %v", err)
	}
	err = p.consumer.Assign(partitions)
	if err != nil {
		return nil, fmt.Errorf("kafka.SubscribeTopic: %v", err)
	}

	return offsets, nil
}

// Handle kafka.Events and ensure messages are successfully produced
func (p *Processor) handleDeliveries(deliveryChan chan kafka.Event, deliveryWG *sync.WaitGroup, mu *sync.Mutex, writeErr *error) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				mu.Lock()
				*writeErr = fmt.Errorf("kafka.ProduceMessage: %v", ev.TopicPartition)
				mu.Unlock()
			}
			deliveryWG.Done()
		}
	}
}
