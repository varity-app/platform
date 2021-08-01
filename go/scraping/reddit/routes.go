package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"

	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, redditClient *reddit.Client, fsClient *firestore.Client, psClient *pubsub.Client, producer *kafka.Producer, tracer *trace.Tracer) error {

	web.GET("/scraping/reddit/submissions/:subreddit", func(ctx echo.Context) error {

		subreddit := ctx.Param("subreddit")

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "scrape.reddit.submissions")
		defer span.End()

		// Process reddit submissions
		count, err := scrapeSubmissions(readCtx, redditClient, fsClient, psClient, producer, subreddit, tracer)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Scraped %d reddit submissions from r/%s.", count, subreddit)}
		return ctx.JSON(http.StatusOK, response)
	})

	web.GET("/scraping/reddit/comments/:subreddit", func(ctx echo.Context) error {

		subreddit := ctx.Param("subreddit")

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "scrape.reddit.comments")
		defer span.End()

		// Process reddit comments
		count, err := scrapeComments(readCtx, redditClient, fsClient, psClient, producer, subreddit, tracer)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Scraped %d reddit comments from r/%s.", count, subreddit)}
		return ctx.JSON(http.StatusOK, response)
	})

	return nil
}
