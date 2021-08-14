package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/VarityPlatform/scraping/common"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// ProcessorHandler is an interface used by a Processor to process individual kafka messages
// and return a list of messages to send to subsequent kafka topic
type ProcessorHandler interface {
	Process(*kafka.Message) ([][]byte, error)
}

// ProcessorOpts is a struct used to pass configuration parameters to the NewProcessor constructor
type ProcessorOpts struct {
	BootstrapServers string
	Username         string
	Password         string
	InputTopic       string
	OutputTopic      string
	CheckpointKey    string
}

// Processor reads a Kafka topic and parses the messages for tickers.
type Processor struct {
	offsetManager    *OffsetManager
	writer           *kafka.Writer
	dialer           *kafka.Dialer
	inputTopic       string
	bootstrapServers string
	partitionsCount  int // Number of partitions for each topic
	checkpointKey    string
}

// NewProcessor initializes a new Processor
func NewProcessor(ctx context.Context, offsetManager *OffsetManager, opts ProcessorOpts) *Processor {

	// SASL mechanism
	mechanism := plain.Mechanism{
		Username: opts.Username,
		Password: opts.Password,
	}

	// kafka Dialer, used by kafka.NewReader()
	dialer := &kafka.Dialer{
		Timeout:       2 * time.Second,
		DualStack:     true,
		TLS:           &tls.Config{},
		SASLMechanism: mechanism,
	}

	// Initialize a kafka writer
	writer := &kafka.Writer{
		Addr:      kafka.TCP(opts.BootstrapServers),
		Topic:     opts.OutputTopic,
		BatchSize: 1,
		Transport: &kafka.Transport{
			TLS:  &tls.Config{},
			SASL: mechanism,
		},
	}

	// Return pointer
	return &Processor{
		offsetManager:    offsetManager,
		writer:           writer,
		dialer:           dialer,
		inputTopic:       opts.InputTopic,
		bootstrapServers: opts.BootstrapServers,
		checkpointKey:    opts.CheckpointKey,
		partitionsCount:  common.KafkaPartitionsCount,
	}
}

// Close closes the ticker processor's connection
func (p *Processor) Close() error {
	if err := p.writer.Close(); err != nil {
		return err
	}

	return p.offsetManager.Close()
}

// ProcessTopic processes a kafka topic for ticker mentions
func (p *Processor) ProcessTopic(ctx context.Context, handler ProcessorHandler) (int, error) {
	count := 0

	// Fetch offsets from firstore
	offsets, err := p.offsetManager.Fetch(ctx, p.checkpointKey)
	if err != nil {
		return count, err
	}

	// Create list off error channels to watch for finish
	var errcList []<-chan error

	// Read messages from kafka
	msgs, readErrc := p.readMessages(ctx, p.partitionsCount, offsets)
	errcList = append(errcList, readErrc)

	// apply handler
	newMsgs, handlerErrc := p.handleMessages(ctx, handler, msgs, &count)
	errcList = append(errcList, handlerErrc)

	// sink messages
	writeErrc := p.sinkMessages(ctx, newMsgs)
	errcList = append(errcList, writeErrc)

	// Wait for pipelines to finish
	if err := common.WaitForPipeline(errcList...); err != nil {
		return count, err
	}

	// Commit offsets to firestore
	if err := p.offsetManager.Save(ctx, p.checkpointKey, offsets); err != nil {
		return count, err
	}

	return count, nil
}

// readMessages reads the latest messages from a kafka topic
func (p *Processor) readMessages(ctx context.Context, partitionsCount int, offsets map[string]int) (<-chan *kafka.Message, <-chan error) {
	mu := sync.Mutex{}
	out := make(chan *kafka.Message)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		// Iterate over each partition and create a unique reader for each
		for i := 0; i < partitionsCount; i++ {
			reader, err := p.createReader(ctx, i, offsets)
			if err != nil {
				errc <- err
				return
			}
			defer reader.Close()

			// Read batches of messages
			for {
				timeoutCtx, cancel := context.WithTimeout(ctx, ConsumerTimeout)
				defer cancel()

				count := p.readBatch(timeoutCtx, reader, out, errc)
				if err != nil {
					return
				}
				if count == 0 {
					mu.Lock()
					offsets[fmt.Sprint(i)] = int(reader.Offset())
					mu.Unlock()
					break
				}
			}
		}
	}()

	return out, errc
}

// readBatch reads a batch of messages from a Kafka topic
func (p *Processor) readBatch(ctx context.Context, reader *kafka.Reader, out chan *kafka.Message, errc chan error) int {
	count := 0
	for {
		// Read messages
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			switch err {
			case context.DeadlineExceeded:
				return count
			default:
				errc <- fmt.Errorf("kafka.ReadMessage: %v", err)
				return count
			}
		}

		// Send on output channel
		select {
		case out <- &msg:
			count++
		case <-ctx.Done():
			return count
		}
	}
}

// handleMessages applies the handler function to each message in the input channel
func (p *Processor) handleMessages(ctx context.Context, handler ProcessorHandler, in <-chan *kafka.Message, count *int) (<-chan []byte, <-chan error) {

	out := make(chan []byte)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for msg := range in {
			// Send message through handler
			newMsgs, err := handler.Process(msg)
			if err != nil {
				errc <- fmt.Errorf("handler.Process: %v", err)
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
				for _, newMsg := range newMsgs {
					out <- newMsg
					*count++
				}
			}

		}
	}()

	return out, errc

}

// Write messages back to kafka
func (p *Processor) sinkMessages(ctx context.Context, in <-chan []byte) <-chan error {

	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		for msg := range in {
			select {
			case <-ctx.Done():
				return
			default:
				if err := p.writer.WriteMessages(ctx, kafka.Message{Value: msg}); err != nil {
					errc <- fmt.Errorf("kafka.WriteMessages: %v", err)
					return
				}
			}
		}
	}()

	return errc
}

func (p *Processor) createReader(ctx context.Context, partition int, offsets map[string]int) (*kafka.Reader, error) {
	// Initialize a kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{p.bootstrapServers},
		Topic:     p.inputTopic,
		Partition: partition,
		Dialer:    p.dialer,
		MaxWait:   ConsumerTimeout,
	})

	// Set offset
	offset := offsets[fmt.Sprint(partition)]
	if err := reader.SetOffset(int64(offset)); err != nil {
		return nil, fmt.Errorf("kafka.SetOffset: %v", err)
	}

	return reader, nil
}
