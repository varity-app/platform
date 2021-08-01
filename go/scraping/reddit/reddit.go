package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"

	"go.opentelemetry.io/otel/trace"

	"github.com/VarityPlatform/scraping/common"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-querystring/query"
	"google.golang.org/protobuf/proto"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const PRODUCER_TIMEOUT_MS = 5 * 1000 // 5 seconds

// Initialize the reddit client.  Load credentials from environment variables.
func initReddit() (*reddit.Client, error) {
	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}

	userAgent := os.Getenv("REDDIT_USER_AGENT")

	client, err := reddit.NewClient(credentials, reddit.WithUserAgent(userAgent))
	if err != nil {
		return nil, fmt.Errorf("redditClient.Create: %v", err)
	}

	return client, nil
}

// Scrape submissions from reddit
func scrapeSubmissions(ctx context.Context, redditClient *reddit.Client, fsClient *firestore.Client, psClient *pubsub.Client, producer *kafka.Producer, subreddit string, tracer *trace.Tracer) (int, error) {

	// Fetch posts
	ctx, span := (*tracer).Start(ctx, "query.Reddit")
	posts, _, err := redditClient.Subreddit.NewPosts(ctx, subreddit, &reddit.ListOptions{
		Limit: LIMIT,
	})
	if err != nil {
		return 0, fmt.Errorf("fetch.Post: %v", err)
	}
	span.End()

	// Check with datastore for unseen posts
	ctx, span = (*tracer).Start(ctx, "query.Datastore")
	newPosts, err := getNewDSSubmissions(ctx, fsClient, posts)
	if err != nil {
		return 0, err
	}
	span.End()

	// Send submissions to kafka
	_, span = (*tracer).Start(ctx, "kafka.PublishSubmissions")
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

	for _, submission := range newPosts {
		// Convert to proto
		submissionProto := submissionToProto(submission)

		// Serialize
		serializedSubmission, err := proto.Marshal(submissionProto)
		if err != nil {
			return 0, fmt.Errorf("serialize.Submission: %v", err)
		}

		wg.Add(1) // Add wait counter
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: stringToPtr(common.REDDIT_SUBMISSIONS), Partition: kafka.PartitionAny},
			Value:          serializedSubmission,
		}, deliveryChan)
	}

	// Flush queue
	producer.Flush(PRODUCER_TIMEOUT_MS)

	// Wait for all messages to process in the delivery report handler
	wg.Wait()
	if writeErr != nil {
		return 0, writeErr
	}
	span.End()

	// Save newly processed submissions to datastore
	ctx, span = (*tracer).Start(ctx, "datastore.WriteSubmissions")
	if err = saveNewDSSubmissions(ctx, fsClient, newPosts); err != nil {
		return 0, fmt.Errorf("datastore.SaveSubmissions: %v", err)
	}
	span.End()

	return len(newPosts), nil
}

// Scrape comments from reddit
func scrapeComments(ctx context.Context, redditClient *reddit.Client, fsClient *firestore.Client, psClient *pubsub.Client, producer *kafka.Producer, subreddit string, tracer *trace.Tracer) (int, error) {

	// Fetch comments
	ctx, span := (*tracer).Start(ctx, "query.Reddit")
	comments, _, err := getNewRedditComments(ctx, redditClient, subreddit, &reddit.ListOptions{
		Limit: LIMIT,
	})
	if err != nil {
		return 0, fmt.Errorf("fetch.Comment: %v", err)
	}
	span.End()

	// Check with datastore for unseen comments
	ctx, span = (*tracer).Start(ctx, "query.Datastore")
	newComments, err := getNewDSComments(ctx, fsClient, comments)
	if err != nil {
		return 0, err
	}
	span.End()

	// Send comments to kafka
	_, span = (*tracer).Start(ctx, "kafka.PublishComments")
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

	for _, comment := range newComments {
		// Convert to proto
		commentProto := commentToProto(comment)

		// Serialize
		serializedComment, err := proto.Marshal(commentProto)
		if err != nil {
			return 0, fmt.Errorf("serialize.Comment: %v", err)
		}

		wg.Add(1) // Add wait counter
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: stringToPtr(common.REDDIT_COMMENTS), Partition: kafka.PartitionAny},
			Value:          serializedComment,
		}, deliveryChan)
	}

	// Flush queue
	producer.Flush(PRODUCER_TIMEOUT_MS)

	// Wait for all messages to process in the delivery report handler
	wg.Wait()
	if writeErr != nil {
		return 0, writeErr
	}
	span.End()

	ctx, span = (*tracer).Start(ctx, "datastore.WriteComments")
	// Save newly processed comments to datastore
	if err = saveNewDSComments(ctx, fsClient, newComments); err != nil {
		return 0, fmt.Errorf("datastore.SaveComments: %v", err)
	}
	span.End()

	return len(newComments), nil
}

// Get comments.  This method and everything in things.go had to be included
// Because go-reddit did not implement a comments endpoint
func getNewRedditComments(ctx context.Context, redditClient *reddit.Client, subreddit string, opts *reddit.ListOptions) ([]*reddit.Comment, *reddit.Response, error) {
	// Add options to HTTP path
	path, err := addOptions("r/"+subreddit+"/comments", opts)
	if err != nil {
		return nil, nil, err
	}

	// Create the request
	req, err := redditClient.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	// Run the request
	t := new(thing)
	resp, err := redditClient.Do(ctx, req, t)
	if err != nil {
		return nil, nil, err
	}

	// Parse data to comments
	l, _ := t.Listing()

	return l.Comments(), resp, nil
}

// Add options to an HTTP path (reddit API)
// This is a copy of the same method in go-reddit only used for getNewRedditComments().
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// Utility function
func stringToPtr(s string) *string {
	return &s
}
