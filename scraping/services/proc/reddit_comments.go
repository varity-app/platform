package main

import (
	"fmt"
	"sync"

	"github.com/VarityPlatform/scraping/common"
	redditPB "github.com/VarityPlatform/scraping/protobuf/reddit"
	"github.com/VarityPlatform/scraping/transforms"

	"google.golang.org/protobuf/proto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Convert a kafka message into a reddit submission proto
func kafkaToRedditComment(msg *kafka.Message) (interface{}, error) {
	submission := &redditPB.RedditComment{}
	if err := proto.Unmarshal(msg.Value, submission); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	return submission, nil
}

// Handler for processing kafka messages containing reddit submissions
func handleCommentTickerProcessing(producer *kafka.Producer, readChan chan *kafka.Message, deliveryChan chan kafka.Event, mu *sync.Mutex, readErr *error, readWG *sync.WaitGroup, deliveryWG *sync.WaitGroup, count *int, allTickers []common.IEXTicker) {
	extractor := transforms.NewTickerExtractor(allTickers)

	for msg := range readChan {
		submission := &redditPB.RedditComment{}
		if err := proto.Unmarshal(msg.Value, submission); err != nil {
			mu.Lock()
			*readErr = fmt.Errorf("protobuf.Unmarshal: %v", err)
			readWG.Done()
			mu.Unlock()
			continue
		}

		mu.Lock()
		*count++
		mu.Unlock()
		readWG.Done()

		// Handle each mention
		mentions := transforms.TransformRedditComment(extractor, submission)
		for _, mention := range mentions {

			// Serialize
			serializedMention, err := proto.Marshal(&mention)
			if err != nil {
				mu.Lock()
				*readErr = fmt.Errorf("protobuf.Marshal: %v", err)
				mu.Unlock()
				continue
			}

			// Write to ticker mentions topic
			deliveryWG.Add(1)
			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: common.StringToPtr(common.TickerMentions), Partition: kafka.PartitionAny},
				Value:          serializedMention,
			}, deliveryChan)
			if err != nil {
				*readErr = fmt.Errorf("kafka.Produce: %v", err)
			}

		}
	}
}
