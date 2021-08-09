package main

import (
	"context"
	"log"
	"os"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/data/kafka"
	"github.com/VarityPlatform/scraping/scrapers"

	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"
)

// Entrypoint method
func main() {
	initConfig()

	ctx := context.Background()

	// Initialize scrapers
	submissionsScraper, err := initSubmissionsScraper(ctx, scrapers.MemoryOpts{
		CollectionName: common.RedditSubmissions + "-" + viper.GetString("deploymentMode"),
	})
	if err != nil {
		log.Fatalf("scraper.Init: %v", err)
	}

	commentsScraper, err := initCommentsScraper(ctx, scrapers.MemoryOpts{
		CollectionName: common.RedditComments + "-" + viper.GetString("deploymentMode"),
	})
	if err != nil {
		log.Fatalf("scraper.Init: %v", err)
	}

	// Initialize publisher
	publisher, err := initPublisher(ctx, kafka.KafkaOpts{
		BootstrapServers: os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		Username:         os.Getenv("KAFKA_AUTH_KEY"),
		Password:         os.Getenv("KAFKA_AUTH_SECRET"),
	})
	if err != nil {
		log.Fatalf("kafka.NewPublisher: %v", err)
	}
	defer publisher.Close()

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	initRoutes(web, submissionsScraper, commentsScraper, publisher)

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
