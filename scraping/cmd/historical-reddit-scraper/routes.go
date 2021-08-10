package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/data/kafka"
	"github.com/VarityPlatform/scraping/scrapers/reddit/historical"
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
func initRoutes(web *echo.Echo, submissionsScraper *historical.SubmissionsScraper, commentsScraper *historical.CommentsScraper, publisher *kafka.Publisher) {

	web.POST("/scraping/reddit/historical/submissions", func(c echo.Context) error {

		body := new(ScrapingRequestBody)
		if err := c.Bind(body); err != nil {
			log.Println(err)
			return err
		}

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

		ctx := c.Request().Context()

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
		publisher.Publish(msgs, common.RedditSubmissions)

		// Save seen msgs to memory
		err = submissionsScraper.CommitSeen(ctx, submissions)
		if err != nil {
			log.Printf("submissions.CommitSeen: %v", err)
			return err
		}

		response := response{Message: fmt.Sprintf("Scraped %d reddit submissions from r/%s.", len(submissions), body.Subreddit)}
		return c.JSON(http.StatusOK, response)

	})

	web.POST("/scraping/reddit/historical/comments", func(c echo.Context) error {

		body := new(ScrapingRequestBody)
		if err := c.Bind(body); err != nil {
			log.Println(err)
			return err
		}

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

		ctx := c.Request().Context()

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
		publisher.Publish(msgs, common.RedditComments)

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
