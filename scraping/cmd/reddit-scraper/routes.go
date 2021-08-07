package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/scrapers/reddit/live"

	"github.com/VarityPlatform/scraping/data/kafka"

	"github.com/labstack/echo/v4"
)

type response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, submissionsScraper *live.RedditSubmissionsScraper, commentsScraper *live.RedditCommentsScraper, publisher *kafka.Publisher) error {

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
		err = publisher.Publish(serializedMsgs, common.RedditSubmissions)
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
		err = publisher.Publish(serializedMsgs, common.RedditComments)
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
