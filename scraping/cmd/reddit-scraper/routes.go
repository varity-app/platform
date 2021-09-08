package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	username := os.Getenv("KAFKA_AUTH_KEY")
	password := os.Getenv("KAFKA_AUTH_SECRET")

	// Scrape most recent submissions
	web.GET("/scraping/reddit/submissions/:subreddit", func(c echo.Context) error {

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
			log.Printf("kafka.NewPublisher: %v", err)
			return err
		}
		defer publisher.Close()

		// Scrape submissions from reddit
		submissions, err := submissionsScraper.Scrape(ctx, subreddit, limit)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Serialize messages
		serializedMsgs, err := serializeSubmissions(submissions)
		if err != nil {
			log.Println(err)
			return err
		}

		// Publish submissions to kafka
		err = publisher.Publish(ctx, serializedMsgs)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Commit submissions to memory
		err = submissionsScraper.CommitSeen(ctx, submissions)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Scraped %d reddit submissions from r/%s.", len(submissions), subreddit)}
		return c.JSON(http.StatusOK, response)
	})

	// Scrape most recent comments
	web.GET("/scraping/reddit/comments/:subreddit", func(c echo.Context) error {

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
			log.Printf("kafka.NewPublisher: %v", err)
			return err
		}
		defer publisher.Close()

		// Scrape comments from reddit
		comments, err := commentsScraper.Scrape(ctx, subreddit, limit)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Serialize messages
		serializedMsgs, err := serializeComments(comments)
		if err != nil {
			log.Println(err)
			return err
		}

		// Publish submissions to kafka
		err = publisher.Publish(ctx, serializedMsgs)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Commit comments to memory
		err = commentsScraper.CommitSeen(ctx, comments)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Scraped %d reddit comments from r/%s.", len(comments), subreddit)}
		return c.JSON(http.StatusOK, response)
	})

	return nil
}
