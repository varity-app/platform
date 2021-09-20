package main

import (
	"context"
	"fmt"
	"log"

	"github.com/varity-app/platform/scraping/internal/config"
	"github.com/varity-app/platform/scraping/internal/logging"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/labstack/echo/v4"

	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
)

var logger = logging.NewLogger()

// Entrypoint method
func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.scraping.address"),
		Password: viper.GetString("redis.scraping.password"),
	})
	defer rdb.Close()

	// Initialize the reddit scrapers
	credentials := reddit.Credentials{
		ID:       viper.GetString("reddit.client_id"),
		Secret:   viper.GetString("reddit.client_secret"),
		Username: viper.GetString("reddit.username"),
		Password: viper.GetString("reddit.password"),
	}
	submissionsScraper, err := initSubmissionsScraper(
		ctx,
		credentials,
		rdb,
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("initSubmissionsScraper: %v", err))
	}

	commentsScraper, err := initCommentsScraper(
		ctx,
		credentials,
		rdb,
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("initCommentsScraper: %v", err))
	}

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
