package main

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"cloud.google.com/go/pubsub"

	"github.com/VarityPlatform/scraping/common"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonPB "github.com/VarityPlatform/scraping/protobuf/common"
	redditPB "github.com/VarityPlatform/scraping/protobuf/reddit"
)

var questionRegex *regexp.Regexp = regexp.MustCompile(`\?`)

// Read reddit submissions from Pub/Sub and process them
func readRedditSubmission(ctx context.Context, psClient *pubsub.Client, allTickers []common.IEXTicker) (int, error) {

	sub := psClient.Subscription(common.SUBSCRIPTION_REDDIT_SUBSCRIPTIONS + "-" + viper.GetString("deploymentMode"))
	topic := psClient.Topic(common.TOPIC_TICKER_MENTIONS + "-" + viper.GetString("deploymentMode"))

	cm := make(chan *pubsub.Message)
	defer close(cm)

	wg := new(sync.WaitGroup)
	var mu sync.Mutex

	count := 0
	totalCount := 0

	var readErr error = nil

	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			submission := &redditPB.RedditSubmission{}
			if err := proto.Unmarshal(msg.Data, submission); err != nil {
				mu.Lock()
				readErr = fmt.Errorf("error unmarshalling protobuf message: %v", err)
				mu.Unlock()
				return
			}

			// Publish each mention to Pub/Sub
			mentions := procRedditSubmission(submission, allTickers)
			for _, mention := range mentions {
				mu.Lock()
				wg.Add(1)
				mu.Unlock()

				// Serialize
				serializedMention, err := proto.Marshal(&mention)
				if err != nil {
					mu.Lock()
					readErr = fmt.Errorf("error serializing mention protobuf message: %v", err)
					mu.Unlock()
					return
				}

				// Asynchronously publish to pubsub
				result := topic.Publish(ctx, &pubsub.Message{
					Data: serializedMention,
				})

				go func(res *pubsub.PublishResult) {
					defer wg.Done()
					_, err := result.Get(ctx)
					if err != nil {
						mu.Lock()
						readErr = fmt.Errorf("error publishing mention message: %v", err)
						mu.Unlock()
						return
					}
				}(result)

			}

			// Update counter
			mu.Lock()
			count++
			mu.Unlock()

			msg.Ack()
		}
	}()

	// Receive messages for N sec
	for {
		ctx, cancel := context.WithTimeout(ctx, WAIT_INTERVAL)
		defer cancel()

		// Reset count
		count = 0

		// Receive blocks until the context is cancelled or an error occurs.
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			if readErr != nil {
				cancel()
			}
			cm <- msg
		})
		if err != nil {
			return totalCount, err
		} else if readErr != nil {
			return totalCount, readErr
		}

		totalCount += count

		// If no messages were received in this window, break
		if count == 0 {
			break
		}
	}

	wg.Wait()

	return totalCount, nil

}

// Read reddit comment from Pub/Sub and process them
func readRedditComment(ctx context.Context, psClient *pubsub.Client, allTickers []common.IEXTicker) (int, error) {

	sub := psClient.Subscription(common.SUBSCRIPTION_REDDIT_COMMENTS + "-" + viper.GetString("deploymentMode"))
	topic := psClient.Topic(common.TOPIC_TICKER_MENTIONS + "-" + viper.GetString("deploymentMode"))

	cm := make(chan *pubsub.Message)
	defer close(cm)

	wg := new(sync.WaitGroup)
	var mu sync.Mutex

	// Message counts
	count := 0
	totalCount := 0

	var readErr error = nil

	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			comment := &redditPB.RedditComment{}
			if err := proto.Unmarshal(msg.Data, comment); err != nil {
				mu.Lock()
				readErr = fmt.Errorf("error unmarshalling protobuf message: %v", err)
				mu.Unlock()
				return
			}

			// Publish each mention to Pub/Sub
			mentions := procRedditComment(comment, allTickers)
			for _, mention := range mentions {
				mu.Lock()
				wg.Add(1)
				mu.Unlock()

				// Serialize
				serializedMention, err := proto.Marshal(&mention)
				if err != nil {
					mu.Lock()
					readErr = fmt.Errorf("error serializing mention protobuf message: %v", err)
					mu.Unlock()
					return
				}

				// Asynchronously publish to pubsub
				result := topic.Publish(ctx, &pubsub.Message{
					Data: serializedMention,
				})

				go func(res *pubsub.PublishResult) {
					defer wg.Done()
					_, err := result.Get(ctx)
					if err != nil {
						mu.Lock()
						readErr = fmt.Errorf("error publishing mention message: %v", err)
						mu.Unlock()
						return
					}
				}(result)

			}

			// Update counter
			mu.Lock()
			count++
			mu.Unlock()

			msg.Ack()
		}
	}()

	// Receive messages for N sec
	for {
		ctx, cancel := context.WithTimeout(ctx, WAIT_INTERVAL)
		defer cancel()

		// Reset count
		count = 0

		// Receive blocks until the context is cancelled or an error occurs.
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			if readErr != nil {
				cancel()
			}
			cm <- msg
		})
		if err != nil {
			return totalCount, err
		} else if readErr != nil {
			return totalCount, readErr
		}

		totalCount += count

		// If no messages were received in this window, break
		if count == 0 {
			break
		}
	}

	wg.Wait()

	return totalCount, nil

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
