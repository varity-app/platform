package kafka

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"cloud.google.com/go/bigquery"
	"github.com/VarityPlatform/scraping/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProtoConverter is a function that converts a kafka message into a format
// serializable by bigquery
type ProtoConverter func(*kafka.Message) (interface{}, error)

// BigquerySink reads from a Kafka topic and saves the messages to BigQuery
type BigquerySink struct {
	bqClient      *bigquery.Client
	consumer      *kafka.Consumer
	offsetManager *OffsetManager
}

// NewBigquerySink initializes a new KafkaBigQuery sink
func NewBigquerySink(ctx context.Context, offsetManager *OffsetManager, opts Opts) (*BigquerySink, error) {
	// Initialize bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}

	// Initialize kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        opts.BootstrapServers,
		"security.protocol":        "SASL_SSL",
		"sasl.mechanisms":          "PLAIN",
		"sasl.username":            opts.Username,
		"sasl.password":            opts.Password,
		"group.id":                 "sink",
		"enable.auto.offset.store": "false",
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.NewConsumer: %v", err)
	}

	return &BigquerySink{
		bqClient:      bqClient,
		consumer:      consumer,
		offsetManager: offsetManager,
	}, nil
}

// Close closes the sink connection
func (s *BigquerySink) Close() error {
	if err := s.consumer.Close(); err != nil {
		return err
	} else if err := s.bqClient.Close(); err != nil {
		return err
	}
	return s.offsetManager.Close()
}

// SinkTopic reads a topic and saves messages to bigquery
func (s *BigquerySink) SinkTopic(ctx context.Context, topic, checkpointKey, datasetName, tableName string, batchSize int, converter ProtoConverter) (int, error) {
	count := 0

	// Get kafka offsets and assign to the consumer
	offsets, err := s.subscribe(ctx, topic, checkpointKey)
	if err != nil {
		return count, err
	}

	// Create bigquery inserter
	inserter := s.bqClient.Dataset(datasetName).Table(tableName).Inserter()

	// Consume messages from Kafka
	var items []interface{}
	var writeErr error
	writeWG := new(sync.WaitGroup)
	for {
		msg, err := s.consumer.ReadMessage(ConsumerTimeout)
		if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
			break
		} else if err != nil {
			return count, fmt.Errorf("kafka.ReadMessage: %v", err)
		}

		// Convert kafka messages
		item, err := converter(msg)
		if err != nil {
			return count, err
		}
		items = append(items, item)
		count++

		// Update offsets
		offset, err := strconv.Atoi(msg.TopicPartition.Offset.String())
		if err != nil {
			return count, fmt.Errorf("strconv.Atoi: %v", err)
		}
		offsets[fmt.Sprint(msg.TopicPartition.Partition)] = offset + 1

		// Insert into bigquery
		if len(items) >= batchSize {
			writeWG.Add(1)
			go s.saveBigqueryBatch(ctx, inserter, items, &writeErr, writeWG)
			items = []interface{}{}
		}

	}

	// Insert remaining items into bigquery
	if len(items) > 0 {
		writeWG.Add(1)
		go s.saveBigqueryBatch(ctx, inserter, items, &writeErr, writeWG)
	}

	// Wait for remaining writes to bigquery
	writeWG.Wait()
	if writeErr != nil {
		return 0, writeErr
	}

	// Save offsets to firestore
	if err := s.offsetManager.Save(ctx, checkpointKey, offsets); err != nil {
		return count, err
	}

	return count, nil
}

// subscribe fetches offsets with the offsetManager and subscribes the consumer
func (s *BigquerySink) subscribe(ctx context.Context, topic string, checkpointKey string) (map[string]int, error) {
	err := s.consumer.Unassign()
	if err != nil {
		return nil, fmt.Errorf("kafka.Unassign: %v", err)
	}
	partitions, offsets, err := s.offsetManager.Fetch(ctx, topic, checkpointKey)
	if err != nil {
		return nil, fmt.Errorf("kafka.FetchOffsets: %v", err)
	}
	err = s.consumer.Assign(partitions)
	if err != nil {
		return nil, fmt.Errorf("kafka.SubscribeTopic: %v", err)
	}

	return offsets, nil
}

// Save a batch to bigquery
func (s *BigquerySink) saveBigqueryBatch(ctx context.Context, inserter *bigquery.Inserter, items []interface{}, writeErr *error, writeWG *sync.WaitGroup) {
	// Upload to bigquery
	if err := inserter.Put(ctx, items); err != nil {
		*writeErr = fmt.Errorf("bigquery.Insert: %v", err)
	}
	writeWG.Done()
}
