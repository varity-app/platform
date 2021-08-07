package kafka

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/VarityPlatform/scraping/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OffsetManager manages kafka offsets and saves them to firestore
type OffsetManager struct {
	fsClient       *firestore.Client
	collectionName string
}

// OffsetManagerOpts is a config struct for initializing a new OffsetManager
type OffsetManagerOpts struct {
	collectionName string
}

// NewOffsetManager initializes a new OffsetManager
func NewOffsetManager(ctx context.Context, opts OffsetManagerOpts) (*OffsetManager, error) {
	fsClient, err := firestore.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.GetClient: %v", err)
	}

	return &OffsetManager{
		fsClient:       fsClient,
		collectionName: opts.collectionName,
	}, nil
}

// Close closes the offset manager's connection
func (manager *OffsetManager) Close() error {
	return manager.fsClient.Close()
}

// Fetch takes a kafka topic name and an offset firestore document key, and returns an array of kafka TopicPartitions
func (manager *OffsetManager) Fetch(ctx context.Context, topic string, key string) ([]kafka.TopicPartition, map[string]int, error) {
	collection := manager.fsClient.Collection(manager.collectionName)
	ref := collection.Doc(key)

	// Get snapshot
	snap, err := ref.Get(ctx)
	if status.Code(err) == codes.NotFound {
		offsets := make(map[string]int)
		return manager.mapOffsetsToPartitions(topic, offsets), offsets, nil
	} else if err != nil {
		return nil, nil, fmt.Errorf("offsetManager.GetOffsets: %v", err)
	}

	data := snap.Data()
	offsets := make(map[string]int)

	// Format to correct map type
	for key, value := range data {
		newValue, err := strconv.Atoi(fmt.Sprint(value))
		if err != nil {
			return nil, nil, fmt.Errorf("offsetManager.DecodeOffsets: %v", err)
		}
		offsets[key] = newValue
	}

	topicPartitions := manager.mapOffsetsToPartitions(topic, offsets)

	return topicPartitions, offsets, err
}

// mapOffsetsToPartitions takes a kafka topic and a map of offsets and returns an array of kafka TopicPartitions
// ready to be used by the kafka client library
func (manager *OffsetManager) mapOffsetsToPartitions(topic string, offsets map[string]int) []kafka.TopicPartition {

	partitions := []kafka.TopicPartition{}

	for partition := 0; partition < common.KafkaPartitionsCount; partition++ {
		partitions = append(partitions, kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(partition),
			Offset:    kafka.Offset(offsets[fmt.Sprint(partition)]),
		})
	}

	return partitions
}

// Save kafka offsets to firestore
func (manager *OffsetManager) Save(ctx context.Context, key string, offsets map[string]int) error {
	collection := manager.fsClient.Collection(manager.collectionName)
	ref := collection.Doc(key)

	_, err := ref.Set(ctx, offsets)
	if err != nil {
		return fmt.Errorf("firestore.SaveOffsets: %v", err)
	}

	return nil
}
