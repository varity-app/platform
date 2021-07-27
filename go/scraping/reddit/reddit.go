package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"cloud.google.com/go/firestore"

	"github.com/google/go-querystring/query"
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
func scrapeSubmissions(ctx context.Context, redditClient *reddit.Client, fsClient *firestore.Client, subreddit string) {

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

	for _, post := range newPosts {
		log.Println(post.ID)
	}
}

// Scrape comments from reddit
func scrapeComments(ctx context.Context, redditClient *reddit.Client, fsClient *firestore.Client, subreddit string) {

	// Fetch comments
	comments, _, err := getNewRedditComments(ctx, redditClient, subreddit, &reddit.ListOptions{
		Limit: LIMIT,
	})
	if err != nil {
		log.Fatal("Error retrieving comments:", err.Error())
	}

	newComments := getNewDSComments(ctx, fsClient, comments)
	log.Printf("Fetched %d unseen comments...\n", len(newComments))

	for _, comment := range newComments {
		log.Println(comment.ID)
	}
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
