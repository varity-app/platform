package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/data/kafka"
	"github.com/varity-app/platform/scraping/internal/scrapers/reddit/historical"
)

type response struct {
	Message string `json:"message"`
}

// ScrapingRequestBody is the definition of a scraping HTTP request body
type ScrapingRequestBody struct {
	Subreddit string `json:"subreddit"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Limit     int    `json:"limit"`
}

// Initialize routes
func initRoutes(web *echo.Echo, submissionsScraper *historical.SubmissionsScraper, commentsScraper *historical.CommentsScraper) {

	// Kafka auth
	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	username := os.Getenv("KAFKA_AUTH_KEY")
	password := os.Getenv("KAFKA_AUTH_SECRET")

	// Scrape submissions
	web.POST("/scraping/reddit/historical/submissions", func(c echo.Context) error {

		// Create context
		ctx := c.Request().Context()

		// Parse request body
		body := new(ScrapingRequestBody)
		if err := c.Bind(body); err != nil {
			log.Println(err)
			return err
		}

		// Initialize publisher
		publisher, err := initPublisher(ctx, kafka.PublisherOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			Topic:            common.RedditSubmissions,
		})
		if err != nil {
			log.Printf("kafka.NewPublisher: %v", err)
			return err
		}
		defer publisher.Close()

		// Parse time fields
		before, err := time.Parse(time.RFC3339, body.Before)
		if err != nil {
			log.Println(err, before)
			return echo.NewHTTPError(http.StatusBadRequest, "Field `before` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		after, err := time.Parse(time.RFC3339, body.After)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Field `after` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		submissions, err := submissionsScraper.Scrape(ctx, body.Subreddit, before, after, body.Limit)
		if err != nil {
			log.Printf("submissionsScraper.Scrape: %v", err)
			return err
		}

		// Serialize submissions
		msgs, err := serializeSubmissions(submissions)
		if err != nil {
			log.Printf("submissions.Serialize: %v", err)
			return err
		}

		// Publish to kafka
		if err := publisher.Publish(ctx, msgs); err != nil {
			log.Printf("publisher.Publish: %v", err)
		}

		// Save seen msgs to memory
		err = submissionsScraper.CommitSeen(ctx, submissions)
		if err != nil {
			log.Printf("submissions.CommitSeen: %v", err)
			return err
		}

		response := response{Message: fmt.Sprintf("Scraped %d reddit submissions from r/%s.", len(submissions), body.Subreddit)}
		return c.JSON(http.StatusOK, response)

	})

	// Scrape comments
	web.POST("/scraping/reddit/historical/comments", func(c echo.Context) error {

		// Create context
		ctx := c.Request().Context()

		// Parse request body
		body := new(ScrapingRequestBody)
		if err := c.Bind(body); err != nil {
			log.Println(err)
			return err
		}

		// Initialize publisher
		publisher, err := initPublisher(ctx, kafka.PublisherOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			Topic:            common.RedditComments,
		})
		if err != nil {
			log.Println(fmt.Errorf("kafka.NewPublisher: %v", err))
			return err
		}
		defer publisher.Close()

		// Parse time fields
		before, err := time.Parse(time.RFC3339, body.Before)
		if err != nil {
			log.Println(err, before)
			return echo.NewHTTPError(http.StatusBadRequest, "Field `before` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		after, err := time.Parse(time.RFC3339, body.After)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Field `after` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		comments, err := commentsScraper.Scrape(ctx, body.Subreddit, before, after, body.Limit)
		if err != nil {
			log.Printf("commentsScraper.Scrape: %v", err)
			return err
		}

		// Serialize comments
		msgs, err := serializeComments(comments)
		if err != nil {
			log.Printf("comments.Serialize: %v", err)
			return err
		}

		// Publish to kafka
		if err := publisher.Publish(ctx, msgs); err != nil {
			log.Printf("publisher.Publish: %v", err)
		}

		// Save seen msgs to memory
		err = commentsScraper.CommitSeen(ctx, comments)
		if err != nil {
			log.Printf("comments.CommitSeen: %v", err)
			return err
		}

		response := response{Message: fmt.Sprintf("Scraped %d reddit comments from r/%s.", len(comments), body.Subreddit)}
		return c.JSON(http.StatusOK, response)

	})

}
