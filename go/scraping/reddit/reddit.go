package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"

	"google.golang.org/protobuf/proto"

	"github.com/google/go-querystring/query"
	"github.com/spf13/viper"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Initialize the reddit client.  Load credentials from environment variables.
func initReddit() *reddit.Client {
	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}

	userAgent := os.Getenv("REDDIT_USER_AGENT")

	client, err := reddit.NewClient(credentials, reddit.WithUserAgent(userAgent))
	if err != nil {
		log.Fatal("Error creating reddit client:", err.Error())
	}

	return client
}

// Scrape submissions from reddit
func scrapeSubmissions(ctx context.Context, redditClient *reddit.Client, fsClient *firestore.Client, psClient *pubsub.Client, subreddit string) {

	// Fetch posts
	posts, _, err := redditClient.Subreddit.NewPosts(ctx, subreddit, &reddit.ListOptions{
		Limit: LIMIT,
	})
	if err != nil {
		log.Fatal("Error retrieving posts:", err.Error())
	}

	// log.Println("Response", resp.Rate.Remaining, len(posts))

	newPosts := getNewDSSubmissions(ctx, fsClient, posts)
	log.Printf("Fetched %d unseen posts...\n", len(newPosts))

	// Send submissions to pubsub
	topic := psClient.Topic(REDDIT_SUBMISSIONS + "-" + viper.GetString("deploymentMode"))
	wg := new(sync.WaitGroup)

	for _, post := range newPosts {
		// Convert to proto
		postProto := submissionToProto(post)

		// Serialize
		serializedPost, err := proto.Marshal(postProto)
		if err != nil {
			log.Fatal("Error serializing submission proto:", err.Error())
		}

		// Asynchronously publish to pubsub
		result := topic.Publish(ctx, &pubsub.Message{
			Data: serializedPost,
		})

		wg.Add(1) // Add wait counter
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			_, err := result.Get(ctx)
			if err != nil {
				log.Fatalf("Error publishing msg: %v", err)
			}
		}(result)

	}

	// Wait for all messages to publish
	wg.Wait()
}

// Scrape comments from reddit
func scrapeComments(ctx context.Context, redditClient *reddit.Client, fsClient *firestore.Client, psClient *pubsub.Client, subreddit string) {

	// Fetch comments
	comments, _, err := getNewRedditComments(ctx, redditClient, subreddit, &reddit.ListOptions{
		Limit: LIMIT,
	})
	if err != nil {
		log.Fatal("Error retrieving comments:", err.Error())
	}

	newComments := getNewDSComments(ctx, fsClient, comments)
	log.Printf("Fetched %d unseen comments...\n", len(newComments))

	// Send submissions to pubsub
	topic := psClient.Topic(REDDIT_COMMENTS + "-" + viper.GetString("deploymentMode"))
	wg := new(sync.WaitGroup)

	for _, comment := range newComments {
		// Convert to proto
		commentProto := commentToProto(comment)

		// Serialize
		serializedComment, err := proto.Marshal(commentProto)
		if err != nil {
			log.Fatal("Error serializing comment proto:", err.Error())
		}

		// Asynchronously publish to pubsub
		result := topic.Publish(ctx, &pubsub.Message{
			Data: serializedComment,
		})

		wg.Add(1) // Add wait counter
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			_, err := result.Get(ctx)
			if err != nil {
				log.Fatalf("Error publishing msg: %v", err)
			}
		}(result)

	}

	// Wait for all messages to publish
	wg.Wait()
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
