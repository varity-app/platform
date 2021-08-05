package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, submissionsScraper *scrapers.RedditSubmissionsScraper, commentsScraper *scrapers.RedditCommentsScraper, producer *kafka.Producer, tracer *trace.Tracer) error {

	web.GET("/scraping/reddit/submissions/:subreddit", func(ctx echo.Context) error {

		subreddit := ctx.Param("subreddit")

		// Create context for new process
		reqCtx := ctx.Request().Context()

		// Create tracer span
		reqCtx, span := (*tracer).Start(reqCtx, "scrape.reddit.submissions")
		defer span.End()

		// Scrape submissions from reddit
		submissions, err := submissionsScraper.Scrape(reqCtx, subreddit, LIMIT)
		// count, err := scrapeSubmissions(reqCtx, redditClient, fsClient, psClient, producer, subreddit, tracer)
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

		response := Response{Message: fmt.Sprintf("Scraped %d reddit submissions from r/%s.", len(submissions), subreddit)}
		return ctx.JSON(http.StatusOK, response)
	})

	// web.GET("/scraping/reddit/comments/:subreddit", func(ctx echo.Context) error {

	// 	subreddit := ctx.Param("subreddit")

	// 	// Create context for new process
	// 	readCtx := ctx.Request().Context()

	// 	// Create tracer span
	// 	readCtx, span := (*tracer).Start(readCtx, "scrape.reddit.comments")
	// 	defer span.End()

	// 	// Process reddit comments
	// 	count, err := scrapeComments(readCtx, redditClient, fsClient, psClient, producer, subreddit, tracer)
	// 	if err != nil {
	// 		log.Println(err)
	// 		span.RecordError(err)
	// 		span.SetStatus(codes.Error, "critical error")
	// 		return echo.NewHTTPError(http.StatusInternalServerError)
	// 	}

	// 	response := Response{Message: fmt.Sprintf("Scraped %d reddit comments from r/%s.", count, subreddit)}
	// 	return ctx.JSON(http.StatusOK, response)
	// })

	return nil
}
