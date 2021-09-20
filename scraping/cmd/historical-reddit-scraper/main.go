package main

import (
	"context"
	"fmt"
	"log"

	"github.com/varity-app/platform/scraping/internal/config"
	"github.com/varity-app/platform/scraping/internal/logging"

	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"

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

	// Initialize scrapers
	submissionsScraper, err := initSubmissionsScraper(ctx, rdb)
	if err != nil {
		logger.Fatal(fmt.Errorf("scraper.Init: %v", err))
	}

	commentsScraper, err := initCommentsScraper(ctx, rdb)
	if err != nil {
		logger.Fatal(fmt.Errorf("scraper.Init: %v", err))
	}

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	initRoutes(web, submissionsScraper, commentsScraper)

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
