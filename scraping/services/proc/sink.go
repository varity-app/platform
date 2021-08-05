package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"cloud.google.com/go/bigquery"
	"github.com/VarityPlatform/scraping/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

type kafkaProtoConverter func(*kafka.Message) (interface{}, error)

// KafkaBigquerySink reads from a Kafka topic and saves the messages to BigQuery
type KafkaBigquerySink struct {
	bqClient      *bigquery.Client
	consumer      *kafka.Consumer
	offsetManager *OffsetManager
}

// NewKafkaBigquerySink initializes a new KafkaBigQuery sink
func NewKafkaBigquerySink(ctx context.Context, opts OffsetManagerOpts) (*KafkaBigquerySink, error) {
	// Initialize bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GcpProjectId)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}

	// Initialize kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol":        "SASL_SSL",
		"sasl.mechanisms":          "PLAIN",
		"sasl.username":            os.Getenv("KAFKA_AUTH_KEY"),
		"sasl.password":            os.Getenv("KAFKA_AUTH_SECRET"),
		"group.id":                 "sink",
		"enable.auto.offset.store": "false",
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.NewConsumer: %v", err)
	}

	offsetManager, err := NewOffsetManager(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &KafkaBigquerySink{
		bqClient:      bqClient,
		consumer:      consumer,
		offsetManager: offsetManager,
	}, nil
}

// Close closes the sink connection
func (sink *KafkaBigquerySink) Close() error {
	if err := sink.consumer.Close(); err != nil {
		return err
	} else if err := sink.bqClient.Close(); err != nil {
		return err
	}
	return sink.offsetManager.Close()
}

// SinkTopic reads a topic and saves messages to bigquery
func (sink *KafkaBigquerySink) SinkTopic(ctx context.Context, topic string, tableName string, converter kafkaProtoConverter) (int, error) {
	count := 0

	// Get kafka offsets and assign to the consumer
	err := sink.consumer.Unassign()
	if err != nil {
		return count, fmt.Errorf("kafka.Unassign: %v", err)
	}
	checkpointKey := topic + "-sink"
	partitions, offsets, err := sink.offsetManager.Fetch(ctx, topic, checkpointKey)
	if err != nil {
		return count, fmt.Errorf("kafka.FetchOffsets: %v", err)
	}
	err = sink.consumer.Assign(partitions)
	if err != nil {
		return count, fmt.Errorf("kafka.SubscribeTopic: %v", err)
	}

	// Create bigquery inserter
	inserter := sink.bqClient.Dataset(
		common.BigqueryDatasetScraping + "_" + viper.GetString("deploymentMode"),
	).Table(tableName).Inserter()

	// Consume messages from Kafka
	var items []interface{}
	var writeErr error
	writeWG := new(sync.WaitGroup)
	for {
		msg, err := sink.consumer.ReadMessage(ConsumerTimeout)
		if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
			break
		} else if err != nil {
			return count, fmt.Errorf("kafka.ReadMessage: %v", err)
		}

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
		if len(items) >= SinkBatchSize {
			writeWG.Add(1)
			go sink.saveBigqueryBatch(ctx, inserter, items, &writeErr, writeWG)
			items = []interface{}{}
		}

	}

	// Insert remaining items into bigquery
	if len(items) > 0 {
		writeWG.Add(1)
		go sink.saveBigqueryBatch(ctx, inserter, items, &writeErr, writeWG)
	}

	// Wait for remaining writes to bigquery
	writeWG.Wait()
	if writeErr != nil {
		return 0, writeErr
	}

	// Save offsets to firestore
	if err := sink.offsetManager.Save(ctx, checkpointKey, offsets); err != nil {
		return count, err
	}

	return count, nil
}

// Sink a kafka stream to bigquery
// func sinkKafkaToBigquery(ctx context.Context, fsClient *firestore.Client, bqClient *bigquery.Client, consumer *kafka.Consumer, tracer *trace.Tracer, topic string, tableName string, converter kafkaProtoConverter) (int, error) {

// }

func (sink *KafkaBigquerySink) saveBigqueryBatch(ctx context.Context, inserter *bigquery.Inserter, items []interface{}, writeErr *error, writeWG *sync.WaitGroup) {
	// Upload to bigquery
	if err := inserter.Put(ctx, items); err != nil {
		*writeErr = fmt.Errorf("bigquery.Insert: %v", err)
	}
	writeWG.Done()
}
