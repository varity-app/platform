package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
)

const WALLSTREETBETS string = "wallstreetbets"

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

	// Init pubsub client
	psClient, err := pubsub.NewClient(ctx, GCP_PROJECT_ID)
	if err != nil {
		log.Fatal("Error initializing pubsub client:", err.Error())
	}
	defer psClient.Close()

	scrapeSubmissions(ctx, redditClient, fsClient, psClient, WALLSTREETBETS)
}
