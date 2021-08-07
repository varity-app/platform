package kafka

import (
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Publisher publishes messages to a kafka topic
type Publisher struct {
	producer *kafka.Producer
}

// NewPublisher initializes a new kafka Publisher
func NewPublisher(opts KafkaOpts) (*Publisher, error) {

	// Create producer
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

	return &Publisher{
		producer: producer,
	}, nil
}

// Close closes a Publisher's connection to the server
func (p *Publisher) Close() {
	p.producer.Close()
}

// Publish sends messages to a specific kafka topic
func (p *Publisher) Publish(msgs [][]byte, topic string) error {
	// Send submissions to kafka
	wg := new(sync.WaitGroup)
	var writeErr error = nil

	// Delivery report handler for produced messages
	deliveryChan := make(chan kafka.Event)
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					writeErr = fmt.Errorf("kafka.ProduceMessage: %v", ev.TopicPartition)
				}
				wg.Done()
			}
		}
	}()

	for _, msg := range msgs {
		wg.Add(1) // Add wait counter
		err := p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          msg,
		}, deliveryChan)
		if err != nil {
			return fmt.Errorf("kafka.Produce: %v", err)
		}
	}

	// Flush queue
	p.producer.Flush(ProducerTimeoutMS)

	// Wait for all messages to process in the delivery report handler
	wg.Wait()
	if writeErr != nil {
		return writeErr
	}

	return nil
}
