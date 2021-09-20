package main

import (
	"fmt"

	rpb "github.com/varity-app/platform/scraping/api/reddit/v1"
	"github.com/varity-app/platform/scraping/internal/common"
	transforms "github.com/varity-app/platform/scraping/internal/transforms/tickerext"

	"github.com/segmentio/kafka-go"

	"google.golang.org/protobuf/proto"
)

// RedditSubmissionHandler is a handler for the KafkaProcessor reading the reddit-submissions topic
type RedditSubmissionHandler struct {
	Extractor *transforms.TickerExtractor
}

// NewRedditSubmissionHandler initializes a new RedditSubmissionHandler
func NewRedditSubmissionHandler(tickers []common.IEXTicker) *RedditSubmissionHandler {
	return &RedditSubmissionHandler{
		Extractor: transforms.NewTickerExtractor(tickers),
	}
}

// Process processes a kafka message for a reddit submission, extracts tickers,
// and returns the serialized array of tickers.
func (h *RedditSubmissionHandler) Process(msg *kafka.Message) ([][]byte, error) {

	// Parse submission from kafka message
	submission := &rpb.RedditSubmission{}
	if err := proto.Unmarshal(msg.Value, submission); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	// Parse tickers from submission
	mentions := transforms.TransformRedditSubmission(h.Extractor, submission)

	// Serialize tickers into bytes
	results := [][]byte{}
	for _, mention := range mentions {
		serializedMention, err := proto.Marshal(mention)
		if err != nil {
			return nil, fmt.Errorf("protobuf.Marshal: %v", err)
		}
		results = append(results, serializedMention)
	}

	return results, nil
}

// RedditCommentHandler is a handler for the KafkaProcessor reading the reddit-submissions topic
type RedditCommentHandler struct {
	Extractor *transforms.TickerExtractor
}

// NewRedditCommentHandler initializes a new RedditCommentHandler
func NewRedditCommentHandler(tickers []common.IEXTicker) *RedditCommentHandler {
	return &RedditCommentHandler{
		Extractor: transforms.NewTickerExtractor(tickers),
	}
}

// Process processes a kafka message for a reddit comment, extracts tickers,
// and returns the serialized array of tickers.
func (h *RedditCommentHandler) Process(msg *kafka.Message) ([][]byte, error) {

	// Parse comment from kafka message
	comment := &rpb.RedditComment{}
	if err := proto.Unmarshal(msg.Value, comment); err != nil {
		return nil, fmt.Errorf("protobuf.Unmarshal: %v", err)
	}

	// Parse tickers from comment
	mentions := transforms.TransformRedditComment(h.Extractor, comment)

	// Serialize tickers into bytes
	results := [][]byte{}
	for _, mention := range mentions {
		serializedMention, err := proto.Marshal(mention)
		if err != nil {
			return nil, fmt.Errorf("protobuf.Marshal: %v", err)
		}
		results = append(results, serializedMention)
	}

	return results, nil
}
