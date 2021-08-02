package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"cloud.google.com/go/firestore"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/VarityPlatform/scraping/common"
	commonPB "github.com/VarityPlatform/scraping/protobuf/common"
	redditPB "github.com/VarityPlatform/scraping/protobuf/reddit"
)

var questionRegex *regexp.Regexp = regexp.MustCompile(`\?`)

// Read reddit submissions from Pub/Sub and process them
func extractRedditSubmissionsTickers(ctx context.Context, fsClient *firestore.Client, producer *kafka.Producer, consumer *kafka.Consumer, allTickers []common.IEXTicker, tracer *trace.Tracer) (int, error) {

	var mu sync.Mutex
	count := 0

	// Subscribe to topic
	err := consumer.Unassign()
	if err != nil {
		return count, fmt.Errorf("kafka.Unassign: %v", err)
	}
	_, span := (*tracer).Start(ctx, "firestore.GetOffsets")
	checkpointKey := common.REDDIT_SUBMISSIONS + "-proc"
	offsets, err := getOffsetsFromFS(ctx, fsClient, checkpointKey)
	if err != nil {
		return count, err
	}
	partitions := offsetMapToPartitions(common.REDDIT_SUBMISSIONS, offsets)
	err = consumer.Assign(partitions)
	if err != nil {
		return count, fmt.Errorf("kafka.SubscribeTopic: %v", err)
	}
	span.End()

	// Delivery report handler for produced messages
	deliveryChan := make(chan kafka.Event)
	deliveryWG := new(sync.WaitGroup)
	var writeErr error
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					mu.Lock()
					writeErr = fmt.Errorf("kafka.ProduceMessage: %v", ev.TopicPartition)
					mu.Unlock()
				}
				deliveryWG.Done()
			}
		}
	}()

	// Handle read messages from kafka
	readChan := make(chan *kafka.Message)
	readWG := new(sync.WaitGroup)
	var readErr error
	go func() {
		for msg := range readChan {
			submission := &redditPB.RedditSubmission{}
			if err := proto.Unmarshal(msg.Value, submission); err != nil {
				mu.Lock()
				readErr = fmt.Errorf("protobuf.Unmarshal: %v", err)
				readWG.Done()
				mu.Unlock()
				continue
			}

			mu.Lock()
			count++
			mu.Unlock()
			readWG.Done()

			// Handle each mention
			mentions := procRedditSubmission(submission, allTickers)
			for _, mention := range mentions {

				// Serialize
				serializedMention, err := proto.Marshal(&mention)
				if err != nil {
					mu.Lock()
					readErr = fmt.Errorf("protobuf.Marshal: %v", err)
					mu.Unlock()
					continue
				}

				// Write to ticker mentions topic
				deliveryWG.Add(1)
				err = producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: common.StringToPtr(common.TICKER_MENTIONS), Partition: kafka.PartitionAny},
					Value:          serializedMention,
				}, deliveryChan)
				if err != nil {
					readErr = fmt.Errorf("kafka.Produce: %v", err)
				}

			}
		}
	}()
	span.End()

	// Consume messages from Kafka
	_, span = (*tracer).Start(ctx, "kafka.ReadMessage")
	for {
		msg, err := consumer.ReadMessage(CONSUMER_TIMEOUT)
		if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
			break
		} else if err != nil {
			return count, fmt.Errorf("kafka.ReadMessage: %v", err)
		}

		offset, err := strconv.Atoi(msg.TopicPartition.Offset.String())
		if err != nil {
			return count, fmt.Errorf("strconv.Atoi: %v", err)
		}
		offsets[fmt.Sprint(msg.TopicPartition.Partition)] = offset + 1

		// Check for errors occurring in goroutines
		if readErr != nil {
			return count, readErr
		} else if writeErr != nil {
			return count, writeErr
		}

		readWG.Add(1)
		readChan <- msg

	}

	span.End()
	_, span = (*tracer).Start(ctx, "kafka.Produce")

	producer.Flush(PRODUCER_TIMEOUT_MS)
	readWG.Wait()
	deliveryWG.Wait()
	span.End()

	// Check for errors occurring in goroutines
	if readErr != nil {
		return count, readErr
	} else if writeErr != nil {
		return count, writeErr
	}

	// Save offsets to firestore
	_, span = (*tracer).Start(ctx, "firestore.SaveOffsets")
	if err := saveOffsetsToFS(ctx, fsClient, checkpointKey, offsets); err != nil {
		return count, err
	}
	span.End()

	return count, nil

	return count, nil

}

// Read reddit comments from kafka, parse ticker mentions, and publish to kafka
func extractRedditCommentsTickers(ctx context.Context, fsClient *firestore.Client, producer *kafka.Producer, consumer *kafka.Consumer, allTickers []common.IEXTicker, tracer *trace.Tracer) (int, error) {

	var mu sync.Mutex
	count := 0

	// Subscribe to topic
	err := consumer.Unassign()
	if err != nil {
		return count, fmt.Errorf("kafka.Unassign: %v", err)
	}
	_, span := (*tracer).Start(ctx, "firestore.GetOffsets")
	checkpointKey := common.REDDIT_COMMENTS + "-proc"
	offsets, err := getOffsetsFromFS(ctx, fsClient, checkpointKey)
	if err != nil {
		return count, err
	}
	partitions := offsetMapToPartitions(common.REDDIT_COMMENTS, offsets)
	err = consumer.Assign(partitions)
	if err != nil {
		return count, fmt.Errorf("kafka.SubscribeTopic: %v", err)
	}
	span.End()

	// Delivery report handler for produced messages
	deliveryChan := make(chan kafka.Event)
	deliveryWG := new(sync.WaitGroup)
	var writeErr error
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					mu.Lock()
					writeErr = fmt.Errorf("kafka.ProduceMessage: %v", ev.TopicPartition)
					mu.Unlock()
				}
				deliveryWG.Done()
			}
		}
	}()

	// Handle read messages from kafka
	readChan := make(chan *kafka.Message)
	readWG := new(sync.WaitGroup)
	var readErr error
	go func() {
		for msg := range readChan {
			comment := &redditPB.RedditComment{}
			if err := proto.Unmarshal(msg.Value, comment); err != nil {
				mu.Lock()
				readErr = fmt.Errorf("protobuf.Unmarshal: %v", err)
				readWG.Done()
				mu.Unlock()
				continue
			}

			mu.Lock()
			count++
			mu.Unlock()
			readWG.Done()

			// Handle each mention
			mentions := procRedditComment(comment, allTickers)
			for _, mention := range mentions {

				// Serialize
				serializedMention, err := proto.Marshal(&mention)
				if err != nil {
					mu.Lock()
					readErr = fmt.Errorf("protobuf.Marshal: %v", err)
					mu.Unlock()
					continue
				}

				// Write to ticker mentions topic
				deliveryWG.Add(1)
				err = producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: common.StringToPtr(common.TICKER_MENTIONS), Partition: kafka.PartitionAny},
					Value:          serializedMention,
				}, deliveryChan)
				if err != nil {
					readErr = fmt.Errorf("kafka.Produce: %v", err)
				}

			}
		}
	}()
	span.End()

	// Consume messages from Kafka
	_, span = (*tracer).Start(ctx, "kafka.ReadMessage")
	for {
		msg, err := consumer.ReadMessage(CONSUMER_TIMEOUT)
		if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
			break
		} else if err != nil {
			return count, fmt.Errorf("kafka.ReadMessage: %v", err)
		}

		offset, err := strconv.Atoi(msg.TopicPartition.Offset.String())
		if err != nil {
			return count, fmt.Errorf("strconv.Atoi: %v", err)
		}
		offsets[fmt.Sprint(msg.TopicPartition.Partition)] = offset + 1

		// Check for errors occurring in goroutines
		if readErr != nil {
			return count, readErr
		} else if writeErr != nil {
			return count, writeErr
		}

		readWG.Add(1)
		readChan <- msg

	}

	span.End()
	_, span = (*tracer).Start(ctx, "kafka.Produce")

	producer.Flush(PRODUCER_TIMEOUT_MS)
	readWG.Wait()
	deliveryWG.Wait()
	span.End()

	// Check for errors occurring in goroutines
	if readErr != nil {
		return count, readErr
	} else if writeErr != nil {
		return count, writeErr
	}

	// Save offsets to firestore
	_, span = (*tracer).Start(ctx, "firestore.SaveOffsets")
	if err := saveOffsetsToFS(ctx, fsClient, checkpointKey, offsets); err != nil {
		return count, err
	}
	span.End()

	return count, nil

}

// Process reddit submissions
func procRedditSubmission(submission *redditPB.RedditSubmission, allTickers []common.IEXTicker) []commonPB.TickerMention {
	// Parse tickers
	titleMentions := procPost(submission.Title, submission.SubmissionId, common.PARENT_SOURCE_REDDIT_SUBMISSION_TITLE, submission.Timestamp, allTickers)
	bodyMentions := procPost(submission.Body, submission.SubmissionId, common.PARENT_SOURCE_REDDIT_SUBMISSION_BODY, submission.Timestamp, allTickers)

	// Concat tickers into one array
	allMentions := append(titleMentions, bodyMentions...)

	return allMentions
}

// Process reddit comment
func procRedditComment(comment *redditPB.RedditComment, allTickers []common.IEXTicker) []commonPB.TickerMention {
	// Parse tickers
	mentions := procPost(comment.Body, comment.CommentId, common.PARENT_SOURCE_REDDIT_COMMENT, comment.Timestamp, allTickers)

	return mentions
}

// Process a post text field
func procPost(s string, id string, parentSource string, timestamp *timestamppb.Timestamp, allTickers []common.IEXTicker) []commonPB.TickerMention {
	if s == "" {
		return []commonPB.TickerMention{}
	}

	// Extract tickers and shortname mentions from string
	tickers := extractTickersString(s, allTickers)
	nameTickers := extractShortNamesString(s, tickers)

	// Calculate frequencies
	uniqTickers, tickerFrequencies := calcTickerFrequency(tickers)
	_, nameFrequencies := calcTickerFrequency(nameTickers)

	// Calculate extra metrics
	questionCount := len(questionRegex.FindAllString(s, -1))
	wordCount := len(wordRegex.FindAllString(s, -1))

	// Generate list of ticker mentions
	mentions := []commonPB.TickerMention{}
	for _, ticker := range uniqTickers {
		mentions = append(mentions, commonPB.TickerMention{
			Symbol:            ticker.Symbol,
			ParentId:          id,
			ParentSource:      parentSource,
			Timestamp:         timestamp,
			SymbolCounts:      uint32(tickerFrequencies[ticker.Symbol]),
			ShortNameCounts:   uint32(nameFrequencies[ticker.Symbol]),
			WordCount:         uint32(wordCount),
			QuestionMarkCount: uint32(questionCount),
		})
	}

	return mentions
}
