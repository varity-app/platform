package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/VarityPlatform/scraping/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	rpb "github.com/VarityPlatform/scraping/protobuf/reddit"
	"google.golang.org/protobuf/proto"
)

// Scrape submissions from reddit
func publishSubmissions(ctx context.Context, producer *kafka.Producer, submissions []*rpb.RedditSubmission) error {

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

	for _, submission := range submissions {

		// Serialize
		serializedSubmission, err := proto.Marshal(submission)
		if err != nil {
			return fmt.Errorf("serialize.Submission: %v", err)
		}

		wg.Add(1) // Add wait counter
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: common.StringToPtr(common.RedditSubmissions), Partition: kafka.PartitionAny},
			Value:          serializedSubmission,
		}, deliveryChan)
		if err != nil {
			return fmt.Errorf("kafka.Produce: %v", err)
		}
	}

	// Flush queue
	producer.Flush(ProducerTimeoutMS)

	// Wait for all messages to process in the delivery report handler
	wg.Wait()
	if writeErr != nil {
		return writeErr
	}

	return nil
}
