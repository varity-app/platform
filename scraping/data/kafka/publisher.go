package kafka

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// PublisherOpts is a struct used to pass configuration parameters to the NewPublisher constructor
type PublisherOpts struct {
	BootstrapServers string
	Username         string
	Password         string
	Topic            string
}

// Publisher publishes messages to a kafka topic
type Publisher struct {
	writer *kafka.Writer
}

// NewPublisher initializes a new kafka Publisher
func NewPublisher(opts PublisherOpts) (*Publisher, error) {

	// SASL mechanism
	mechanism := plain.Mechanism{
		Username: opts.Username,
		Password: opts.Password,
	}

	// Initialize a kafka writer
	writer := &kafka.Writer{
		Addr:      kafka.TCP(opts.BootstrapServers),
		Topic:     opts.Topic,
		BatchSize: 1,
		Transport: &kafka.Transport{
			TLS:  &tls.Config{},
			SASL: mechanism,
		},
	}

	return &Publisher{
		writer: writer,
	}, nil
}

// Close closes a Publisher's connection to the server
func (p *Publisher) Close() error {
	return p.writer.Close()
}

// Publish sends messages to a specific kafka topic
func (p *Publisher) Publish(ctx context.Context, msgs [][]byte) error {
	// Send submissions to kafka

	var kafkaMsgs []kafka.Message
	for _, msg := range msgs {
		kafkaMsgs = append(kafkaMsgs, kafka.Message{Value: msg})
	}

	if err := p.writer.WriteMessages(ctx, kafkaMsgs...); err != nil {
		return fmt.Errorf("kafka.WriteMessages: %v", err)
	}

	return nil
}
