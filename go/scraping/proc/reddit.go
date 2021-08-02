package main

import (
	"fmt"
	"regexp"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/VarityPlatform/scraping/common"
	commonPB "github.com/VarityPlatform/scraping/protobuf/common"
	redditPB "github.com/VarityPlatform/scraping/protobuf/reddit"
)

var questionRegex *regexp.Regexp = regexp.MustCompile(`\?`)

// Handler for processing kafka messages containing reddit submissions
func handleKafkaSubmissions(producer *kafka.Producer, readChan chan *kafka.Message, deliveryChan chan kafka.Event, mu *sync.Mutex, readErr *error, readWG *sync.WaitGroup, deliveryWG *sync.WaitGroup, count *int, allTickers []common.IEXTicker) {
	for msg := range readChan {
		submission := &redditPB.RedditSubmission{}
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
		mentions := procRedditSubmission(submission, allTickers)
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
				TopicPartition: kafka.TopicPartition{Topic: common.StringToPtr(common.TICKER_MENTIONS), Partition: kafka.PartitionAny},
				Value:          serializedMention,
			}, deliveryChan)
			if err != nil {
				*readErr = fmt.Errorf("kafka.Produce: %v", err)
			}

		}
	}
}

// Handler for processing kafka messages containing reddit submissions
func handleKafkaComments(producer *kafka.Producer, readChan chan *kafka.Message, deliveryChan chan kafka.Event, mu *sync.Mutex, readErr *error, readWG *sync.WaitGroup, deliveryWG *sync.WaitGroup, count *int, allTickers []common.IEXTicker) {
	for msg := range readChan {
		comment := &redditPB.RedditComment{}
		if err := proto.Unmarshal(msg.Value, comment); err != nil {
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
		mentions := procRedditComment(comment, allTickers)
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
				TopicPartition: kafka.TopicPartition{Topic: common.StringToPtr(common.TICKER_MENTIONS), Partition: kafka.PartitionAny},
				Value:          serializedMention,
			}, deliveryChan)
			if err != nil {
				*readErr = fmt.Errorf("kafka.Produce: %v", err)
			}

		}
	}
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
