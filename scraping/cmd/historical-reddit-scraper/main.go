package main

import (
	"context"
	"log"

	"github.com/VarityPlatform/scraping/common"
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

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	initRoutes(web, submissionsScraper, commentsScraper)

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
