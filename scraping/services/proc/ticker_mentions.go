package main

import (
	"fmt"

	pb "github.com/VarityPlatform/scraping/protobuf/common"

	"google.golang.org/protobuf/proto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Convert a kafka message into a ticker mention
func kafkaToTickerMention(msg *kafka.Message) (interface{}, error) {
	mention := &pb.TickerMention{}
	if err := proto.Unmarshal(msg.Value, mention); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	return mention, nil
}
