package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"

	"github.com/VarityPlatform/scraping/common"
)

type kafkaProtoConverter func(*kafka.Message) (interface{}, error)

const SINK_BATCH_SIZE int = 250

// Save a batch to bigquery
func saveBigqueryBatch(ctx context.Context, inserter *bigquery.Inserter, items []interface{}, writeErr *error, writeWG *sync.WaitGroup, tracer *trace.Tracer) {
	// Upload to bigquery
	_, bqSpan := (*tracer).Start(ctx, "bigquery.Insert")
	if err := inserter.Put(ctx, items); err != nil {
		*writeErr = fmt.Errorf("bigquery.Insert: %v", err)
	}
	bqSpan.End()
	writeWG.Done()
}

// Sink a kafka stream to bigquery
func sinkKafkaToBigquery(ctx context.Context, fsClient *firestore.Client, bqClient *bigquery.Client, consumer *kafka.Consumer, tracer *trace.Tracer, topic string, tableName string, converter kafkaProtoConverter) (int, error) {
	count := 0

	// Get kafka offsets and assign to the consumer
	err := consumer.Unassign()
	if err != nil {
		return count, fmt.Errorf("kafka.Unassign: %v", err)
	}
	_, span := (*tracer).Start(ctx, "firestore.GetOffsets")
	checkpointKey := topic + "-sink"
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

	// Create bigquery inserter
	inserter := bqClient.Dataset(
		common.BIGQUERY_DATASET_SCRAPING + "_" + viper.GetString("deploymentMode"),
	).Table(tableName).Inserter()

	// Consume messages from Kafka
	consumeCtx, span := (*tracer).Start(ctx, "kafka.ReadMessage")
	var items []interface{}
	var writeErr error
	writeWG := new(sync.WaitGroup)
	for {
		msg, err := consumer.ReadMessage(CONSUMER_TIMEOUT)
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
		if len(items) >= SINK_BATCH_SIZE {
			writeWG.Add(1)
			go saveBigqueryBatch(consumeCtx, inserter, items, &writeErr, writeWG, tracer)
			items = []interface{}{}
		}

	}
	span.End()

	// Insert remaining items into bigquery
	if len(items) > 0 {
		writeWG.Add(1)
		go saveBigqueryBatch(consumeCtx, inserter, items, &writeErr, writeWG, tracer)
	}

	writeWG.Wait()
	if writeErr != nil {
		return 0, writeErr
	}

	// Save offsets to firestore
	_, span = (*tracer).Start(ctx, "firestore.SaveOffsets")
	if err := saveOffsetsToFS(ctx, fsClient, checkpointKey, offsets); err != nil {
		return count, err
	}
	span.End()

	return count, nil
}
