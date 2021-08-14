package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/VarityPlatform/scraping/common"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// ProtoConverter is a function that converts a kafka message into a format
// serializable by bigquery
type ProtoConverter func(*kafka.Message) (interface{}, error)

// BigquerySinkOpts is a helper struct that is used for specifying configuration details when
// calling NewBigquerySink
type BigquerySinkOpts struct {
	BootstrapServers string
	Username         string
	Password         string
	InputTopic       string
	CheckpointKey    string
	DatasetName      string
	TableName        string
}

// BigquerySink reads from a Kafka topic and saves the messages to BigQuery
type BigquerySink struct {
	bqClient         *bigquery.Client
	offsetManager    *OffsetManager
	dialer           *kafka.Dialer
	inputTopic       string
	bootstrapServers string
	partitionsCount  int // Number of partitions for each topic
	checkpointKey    string
	datasetName      string
	tableName        string
}

// NewBigquerySink initializes a new KafkaBigQuery sink
func NewBigquerySink(ctx context.Context, offsetManager *OffsetManager, opts BigquerySinkOpts) (*BigquerySink, error) {
	// Initialize bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}

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

	return &BigquerySink{
		bqClient:         bqClient,
		offsetManager:    offsetManager,
		dialer:           dialer,
		inputTopic:       opts.InputTopic,
		bootstrapServers: opts.BootstrapServers,
		checkpointKey:    opts.CheckpointKey,
		datasetName:      opts.DatasetName,
		tableName:        opts.TableName,
		partitionsCount:  common.KafkaPartitionsCount,
	}, nil
}

// Close closes the sink connection
func (s *BigquerySink) Close() error {
	if err := s.bqClient.Close(); err != nil {
		return err
	}
	return s.offsetManager.Close()
}

// SinkTopic reads a topic and saves messages to bigquery
func (s *BigquerySink) SinkTopic(ctx context.Context, batchSize int, converter ProtoConverter) (int, error) {

	// Create bigquery inserter
	inserter := s.bqClient.Dataset(s.datasetName).Table(s.tableName).Inserter()

	// Keep track of the number of messages processed
	count := 0

	// Fetch offsets from firstore
	offsets, err := s.offsetManager.Fetch(ctx, s.checkpointKey)
	if err != nil {
		return count, err
	}

	// Create list off error channels to watch for finish
	var errcList []<-chan error

	// Read messages from kafka
	msgs, readErrc := s.readMessages(ctx, s.partitionsCount, offsets)
	errcList = append(errcList, readErrc)

	// apply handler
	newMsgs, handlerErrc := s.handleMessages(ctx, converter, msgs, &count)
	errcList = append(errcList, handlerErrc)

	// sink messages
	writeErrc := s.sinkMessages(ctx, inserter, newMsgs, batchSize)
	errcList = append(errcList, writeErrc)

	// Wait for pipelines to finish
	if err := common.WaitForPipeline(errcList...); err != nil {
		return count, err
	}

	// Commit offsets to firestore
	if err := s.offsetManager.Save(ctx, s.checkpointKey, offsets); err != nil {
		return count, err
	}

	return count, nil
}

// readMessages reads the latest messages from a kafka topic.
func (s *BigquerySink) readMessages(ctx context.Context, partitionsCount int, offsets map[string]int) (<-chan *kafka.Message, <-chan error) {
	mu := sync.Mutex{}
	out := make(chan *kafka.Message)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		// Iterate over each partition and create a unique reader for each.
		for i := 0; i < partitionsCount; i++ {
			reader, err := s.createReader(ctx, i, offsets)
			if err != nil {
				errc <- err
				return
			}
			defer reader.Close()

			// Read batches of messages.
			for {
				timeoutCtx, cancel := context.WithTimeout(ctx, ConsumerTimeout)
				defer cancel()

				// Get new messages.
				countNew := s.readBatch(timeoutCtx, reader, out, errc)
				if countNew == 0 { // If no new messages, we're done.
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
func (s *BigquerySink) readBatch(ctx context.Context, reader *kafka.Reader, out chan *kafka.Message, errc chan error) int {
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
func (s *BigquerySink) handleMessages(ctx context.Context, converter ProtoConverter, in <-chan *kafka.Message, count *int) (<-chan interface{}, <-chan error) {

	out := make(chan interface{})
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for msg := range in {
			// Send message through protobuf converter
			newMsg, err := converter(msg)
			if err != nil {
				errc <- fmt.Errorf("bigquerySink.ConvertProto: %v", err)
				return
			}

			// Write to output channel
			select {
			case out <- newMsg:
				*count++
			case <-ctx.Done():
				return
			}

		}
	}()

	return out, errc
}

// Save messages to bigquery in batches
func (s *BigquerySink) sinkMessages(ctx context.Context, inserter *bigquery.Inserter, in <-chan interface{}, batchSize int) <-chan error {

	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		var batch []interface{}
		for msg := range in {
			select {
			case <-ctx.Done():
				return
			default:
				batch = append(batch, msg)
				if len(batch) == batchSize {
					if err := inserter.Put(ctx, batch); err != nil {
						errc <- fmt.Errorf("bigquery.Insert: %v", err)
						return
					}
				}
			}
		}

		// Save remaining messages
		select {
		case <-ctx.Done():
			return
		default:
			if err := inserter.Put(ctx, batch); err != nil {
				errc <- fmt.Errorf("bigquery.Insert: %v", err)
				return
			}
		}

	}()

	return errc
}

// Create a new kafka Reader
func (s *BigquerySink) createReader(ctx context.Context, partition int, offsets map[string]int) (*kafka.Reader, error) {
	// Initialize a kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{s.bootstrapServers},
		Topic:     s.inputTopic,
		Partition: partition,
		Dialer:    s.dialer,
		MaxWait:   ConsumerTimeout,
	})

	// Set offset
	offset := offsets[fmt.Sprint(partition)]
	if err := reader.SetOffset(int64(offset)); err != nil {
		return nil, fmt.Errorf("kafka.SetOffset: %v", err)
	}

	return reader, nil
}
