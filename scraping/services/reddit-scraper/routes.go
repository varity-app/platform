package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/otel/trace"
)

type response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, submissionsScraper *scrapers.RedditSubmissionsScraper, commentsScraper *scrapers.RedditCommentsScraper, producer *kafka.Producer, tracer *trace.Tracer) error {

	web.GET("/scraping/reddit/submissions/:subreddit", func(ctx echo.Context) error {

		subreddit := ctx.Param("subreddit")
		qLimit := ctx.QueryParam("limit")
		limit, err := strconv.Atoi(qLimit)
		if err != nil {
			limit = 100
		} else if limit > 500 {
			limit = 500
		} else if limit < 1 {
			limit = 100
		}

		// Create context for new process
		reqCtx := ctx.Request().Context()

		// Create tracer span
		reqCtx, span := (*tracer).Start(reqCtx, "scrape.reddit.submissions")
		defer span.End()

		// Scrape submissions from reddit
		submissions, err := submissionsScraper.Scrape(reqCtx, subreddit, limit)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Publish submissions to kafka
		err = publishSubmissions(reqCtx, producer, submissions)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Commit submissions to memory
		err = submissionsScraper.CommitSeen(reqCtx, submissions)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Scraped %d reddit submissions from r/%s.", len(submissions), subreddit)}
		return ctx.JSON(http.StatusOK, response)
	})

	web.GET("/scraping/reddit/comments/:subreddit", func(ctx echo.Context) error {

		subreddit := ctx.Param("subreddit")
		qLimit := ctx.QueryParam("limit")
		limit, err := strconv.Atoi(qLimit)
		if err != nil {
			limit = 100
		} else if limit > 500 {
			limit = 500
		} else if limit < 1 {
			limit = 100
		}

		// Create context for new process
		reqCtx := ctx.Request().Context()

		// Create tracer span
		reqCtx, span := (*tracer).Start(reqCtx, "scrape.reddit.comments")
		defer span.End()

		// Scrape comments from reddit
		comments, err := commentsScraper.Scrape(reqCtx, subreddit, limit)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Publish comments to kafka
		err = publishComments(reqCtx, producer, comments)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Commit comments to memory
		err = commentsScraper.CommitSeen(reqCtx, comments)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Scraped %d reddit comments from r/%s.", len(comments), subreddit)}
		return ctx.JSON(http.StatusOK, response)
	})

	return nil
}
