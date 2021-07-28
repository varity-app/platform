package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/VarityPlatform/scraping/common"
)

// Entrypoint method
func main() {
	initConfig()

	// Initialize postgres
	db := common.InitPostgres()

	// Fetch list of tickers from postgres
	allTickers, err := fetchTickers(db)
	if err != nil {
		log.Fatalln("Error fetching tickers:", err.Error())
	}

	// Initialize pubsub client
	ctx := context.Background()
	psClient, err := pubsub.NewClient(ctx, common.GCP_PROJECT_ID)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
	}
	defer psClient.Close()

	// Process reddit submissions
	err = readRedditSubmission(ctx, psClient, allTickers)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	// Process reddit comments
	err = readRedditComment(ctx, psClient, allTickers)
	if err != nil {
		log.Fatalln("Error:", err)
	}

}
