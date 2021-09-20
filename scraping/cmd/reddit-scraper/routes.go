package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/viper"

	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/scrapers/reddit/live"

	"github.com/varity-app/platform/scraping/internal/data/kafka"

	"github.com/labstack/echo/v4"
)

type response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, submissionsScraper *live.RedditSubmissionsScraper, commentsScraper *live.RedditCommentsScraper) error {

	// Kafka auth
	bootstrapServers := viper.GetString("kafka.boostrap_servers")
	username := viper.GetString("kafka.auth_key")
	password := viper.GetString("kafka.auth_secret")

	// Scrape most recent submissions
	web.POST("/api/scraping/reddit/live/v1/submissions/:subreddit", func(c echo.Context) error {

		subreddit := c.Param("subreddit")
		qLimit := c.QueryParam("limit")
		limit, err := strconv.Atoi(qLimit)
		if err != nil {
			limit = 100
		} else if limit > 500 {
			limit = 500
		} else if limit < 1 {
			limit = 100
		}

		// Create context for new process
		ctx := c.Request().Context()

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

		// Scrape submissions from reddit
		submissions, err := submissionsScraper.Scrape(ctx, subreddit, limit)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Serialize messages
		serializedMsgs, err := serializeSubmissions(submissions)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Publish submissions to kafka
		err = publisher.Publish(ctx, serializedMsgs)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Commit submissions to memory
		err = submissionsScraper.CommitSeen(ctx, submissions)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		msg := fmt.Sprintf("Scraped %d reddit submissions from r/%s.", len(submissions), subreddit)
		logger.Debug(msg)
		return c.JSON(http.StatusOK, response{Message: msg})
	})

	// Scrape most recent comments
	web.POST("/api/scraping/reddit/live/v1/comments/:subreddit", func(c echo.Context) error {

		subreddit := c.Param("subreddit")
		qLimit := c.QueryParam("limit")
		limit, err := strconv.Atoi(qLimit)
		if err != nil {
			limit = 100
		} else if limit > 500 {
			limit = 500
		} else if limit < 1 {
			limit = 100
		}

		// Create context for new process
		ctx := c.Request().Context()

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

		// Scrape comments from reddit
		comments, err := commentsScraper.Scrape(ctx, subreddit, limit)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Serialize messages
		serializedMsgs, err := serializeComments(comments)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Publish submissions to kafka
		err = publisher.Publish(ctx, serializedMsgs)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		// Commit comments to memory
		err = commentsScraper.CommitSeen(ctx, comments)
		if err != nil {
			logger.Error(err)
			return echo.ErrInternalServerError
		}

		msg := fmt.Sprintf("Scraped %d reddit comments from r/%s.", len(comments), subreddit)
		logger.Debug(msg)
		return c.JSON(http.StatusOK, response{Message: msg})
	})

	return nil
}
