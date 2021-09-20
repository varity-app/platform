package main

import (
	"fmt"
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
	web.POST("/api/scraping/reddit/historical/v1/submissions", func(c echo.Context) error {

		// Create context
		ctx := c.Request().Context()

		// Parse request body
		body := new(ScrapingRequestBody)
		if err := c.Bind(body); err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Initialize publisher
		publisher, err := initPublisher(ctx, kafka.PublisherOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			Topic:            common.RedditSubmissions,
		})
		if err != nil {
			logger.Error(fmt.Errorf("kafka.NewPublisher: %v", err))
			return echo.ErrInternalServerError
		}
		defer publisher.Close()

		// Parse time fields
		before, err := time.Parse(time.RFC3339, body.Before)
		if err != nil {
			logger.Error(fmt.Errorf("time.Parse: %v - %s", err, before))
			return echo.NewHTTPError(http.StatusBadRequest, "Field `before` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		after, err := time.Parse(time.RFC3339, body.After)
		if err != nil {
			logger.Error(fmt.Errorf("time.Parse: %v - %s", err, after))
			return echo.NewHTTPError(http.StatusBadRequest, "Field `after` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		submissions, err := submissionsScraper.Scrape(ctx, body.Subreddit, before, after, body.Limit)
		if err != nil {
			logger.Error(fmt.Errorf("submissionsScraper.Scrape: %v", err))
			return echo.ErrInternalServerError
		}

		// Serialize submissions
		msgs, err := serializeSubmissions(submissions)
		if err != nil {
			logger.Error(fmt.Errorf("submissions.Serialize: %v", err))
			return echo.ErrInternalServerError
		}

		// Publish to kafka
		if err := publisher.Publish(ctx, msgs); err != nil {
			logger.Error(fmt.Errorf("publisher.Publish: %v", err))
		}

		// Save seen msgs to memory
		err = submissionsScraper.CommitSeen(ctx, submissions)
		if err != nil {
			logger.Error(fmt.Errorf("submissions.CommitSeen: %v", err))
			return echo.ErrInternalServerError
		}

		msg := fmt.Sprintf("Scraped %d reddit submissions from r/%s.", len(submissions), body.Subreddit)
		logger.Debug(msg)
		return c.JSON(http.StatusOK, response{Message: msg})
	})

	// Scrape comments
	web.POST("/api/scraping/reddit/historical/v1/comments", func(c echo.Context) error {

		// Create context
		ctx := c.Request().Context()

		// Parse request body
		body := new(ScrapingRequestBody)
		if err := c.Bind(body); err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Initialize publisher
		publisher, err := initPublisher(ctx, kafka.PublisherOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			Topic:            common.RedditComments,
		})
		if err != nil {
			logger.Error(fmt.Errorf("kafka.NewPublisher: %v", err))
			return echo.ErrInternalServerError
		}
		defer publisher.Close()

		// Parse time fields
		before, err := time.Parse(time.RFC3339, body.Before)
		if err != nil {
			logger.Error(fmt.Errorf("time.Parse: %v - %s", err, before))
			return echo.NewHTTPError(http.StatusBadRequest, "Field `before` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		after, err := time.Parse(time.RFC3339, body.After)
		if err != nil {
			logger.Error(fmt.Errorf("time.Parse: %v - %s", err, after))
			return echo.NewHTTPError(http.StatusBadRequest, "Field `after` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		comments, err := commentsScraper.Scrape(ctx, body.Subreddit, before, after, body.Limit)
		if err != nil {
			logger.Error(fmt.Errorf("commentsScraper.Scrape: %v", err))
			return echo.ErrInternalServerError
		}

		// Serialize comments
		msgs, err := serializeComments(comments)
		if err != nil {
			logger.Error(fmt.Errorf("comments.Serialize: %v", err))
			return echo.ErrInternalServerError
		}

		// Publish to kafka
		if err := publisher.Publish(ctx, msgs); err != nil {
			logger.Error(fmt.Errorf("publisher.Publish: %v", err))
		}

		// Save seen msgs to memory
		err = commentsScraper.CommitSeen(ctx, comments)
		if err != nil {
			logger.Error(fmt.Errorf("comments.CommitSeen: %v", err))
			return echo.ErrInternalServerError
		}

		msg := fmt.Sprintf("Scraped %d reddit comments from r/%s.", len(comments), body.Subreddit)
		logger.Debug(msg)
		return c.JSON(http.StatusOK, response{Message: msg})
	})

}
