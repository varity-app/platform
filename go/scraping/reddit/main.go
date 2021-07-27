package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

const WALLSTREETBETS string = "wallstreetbets"

// type thing struct {
// 	Kind string      `json:"kind"`
// 	Data interface{} `json:"data"`
// }

// type things struct {
// 	Comments          []*reddit.Comment
// 	Mores             []*reddit.More
// 	Users             []*reddit.User
// 	Posts             []*reddit.Post
// 	Subreddits        []*reddit.Subreddit
// 	ModActions        []*reddit.ModAction
// 	Multis            []*reddit.Multi
// 	LiveThreads       []*reddit.LiveThread
// 	LiveThreadUpdates []*reddit.LiveThreadUpdate
// }

// type listing struct {
// 	things things
// 	after  string
// }

// func (l *listing) Posts() []*reddit.Post {
// 	if l == nil {
// 		return nil
// 	}
// 	return l.things.Posts
// }

// func (t *thing) Listing() (v *listing, ok bool) {
// 	v, ok = t.Data.(*listing)
// 	return
// }

// Entrypoint method
func main() {
	initConfig()
	redditClient := initReddit()

	ctx := context.Background()

	// Init firestore client
	fsClient, err := firestore.NewClient(ctx, GCP_PROJECT_ID)
	if err != nil {
		log.Fatal("Error initializing firestore client:", err.Error())
	}
	defer fsClient.Close()

	scrapeComments(ctx, redditClient, fsClient, WALLSTREETBETS)
}
