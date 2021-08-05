package main

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/VarityPlatform/scraping/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

// Get checkpointed kafka offsets from firestore
func getOffsetsFromFS(ctx context.Context, fsClient *firestore.Client, key string) (map[string]int, error) {
	collection := fsClient.Collection(FIRESTORE_KAFKA_OFFSETS + "-" + viper.GetString("deploymentMode"))
	ref := collection.Doc(key)

	// Get snapshot
	snap, err := ref.Get(ctx)
	if status.Code(err) == codes.NotFound {
		offsets := make(map[string]int)
		return offsets, nil
	} else if err != nil {
		return nil, fmt.Errorf("firestore.GetOffsets: %v", err)
	}

	data := snap.Data()
	offsets := make(map[string]int)

	// Format to correct map type
	for key, value := range data {
		newValue, err := strconv.Atoi(fmt.Sprint(value))
		if err != nil {
			return nil, fmt.Errorf("firestore.Decode: %v", err)
		}
		offsets[key] = newValue
	}

	return offsets, err
}

func offsetMapToPartitions(topic string, offsets map[string]int) []kafka.TopicPartition {

	partitions := []kafka.TopicPartition{}

	for partition := 0; partition < common.KAFKA_PARTITION_COUNT; partition++ {
		partitions = append(partitions, kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(partition),
			Offset:    kafka.Offset(offsets[fmt.Sprint(partition)]),
		})
	}

	return partitions
}

// Save kafka offsets to firestore
func saveOffsetsToFS(ctx context.Context, fsClient *firestore.Client, key string, offsets map[string]int) error {
	collection := fsClient.Collection(FIRESTORE_KAFKA_OFFSETS + "-" + viper.GetString("deploymentMode"))
	ref := collection.Doc(key)

	_, err := ref.Set(ctx, offsets)
	if err != nil {
		return fmt.Errorf("firestore.SaveOffsets: %v", err)
	}

	return nil
}
