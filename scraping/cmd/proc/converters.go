package main

import (
	"fmt"

	pb "github.com/VarityPlatform/scraping/protobuf/common"
	rpb "github.com/VarityPlatform/scraping/protobuf/reddit"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"google.golang.org/protobuf/proto"
)

// convertKafkaToComment converts a kafka message to a RedditComment
func convertKafkaToComment(msg *kafka.Message) (interface{}, error) {
	comment := &rpb.RedditComment{}
	if err := proto.Unmarshal(msg.Value, comment); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	return comment, nil
}

// convertKafkaToSubmission converts a kafka message to a RedditSubmission
func convertKafkaToSubmission(msg *kafka.Message) (interface{}, error) {
	submission := &rpb.RedditSubmission{}
	if err := proto.Unmarshal(msg.Value, submission); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	return submission, nil
}

// convertKafkaToTickerMention converts a kafka message to a TickerMention
func convertKafkaToTickerMention(msg *kafka.Message) (interface{}, error) {
	mention := &pb.TickerMention{}
	if err := proto.Unmarshal(msg.Value, mention); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	return mention, nil
}
