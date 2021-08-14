package main

import (
	"context"
	"log"
	"os"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/labstack/echo/v4"

	"github.com/spf13/viper"
)

// Entrypoint method
func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Initialize the reddit scrapers
	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}
	submissionsScraper, err := initSubmissionsScraper(
		ctx,
		credentials,
		scrapers.MemoryOpts{CollectionName: common.RedditSubmissions + "-v2-" + viper.GetString("deploymentMode")},
	)
	if err != nil {
		log.Fatalf("initSubmissionsScraper: %v", err)
	}
	defer submissionsScraper.Close()

	commentsScraper, err := initCommentsScraper(
		ctx,
		credentials,
		scrapers.MemoryOpts{CollectionName: common.RedditComments + "-v2-" + viper.GetString("deploymentMode")},
	)
	if err != nil {
		log.Fatalf("initCommentsScraper: %v", err)
	}
	defer commentsScraper.Close()

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	err = setupRoutes(web, submissionsScraper, commentsScraper)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
