package main

import (
	"context"
	"log"

	"github.com/VarityPlatform/scraping/common"
	"github.com/spf13/viper"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo/v4"
)

// Entrypoint method
func main() {
	initConfig()

	// Initialize postgres
	db := common.InitPostgres()

	// Fetch list of tickers from postgres
	allTickers, err := fetchTickers(db)
	if err != nil {
		log.Fatalf("error fetching tickers: %v", err)
	}

	// Initialize pubsub client
	ctx := context.Background()
	psClient, err := pubsub.NewClient(ctx, common.GCP_PROJECT_ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer psClient.Close()

	// Initialize webserver
	web := echo.New()
	err = setupRoutes(web, psClient, allTickers)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))

}
