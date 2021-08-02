package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VarityPlatform/scraping/common"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, fsClient *firestore.Client, producer *kafka.Producer, consumer *kafka.Consumer, bqClient *bigquery.Client, tracer *trace.Tracer, allTickers []common.IEXTicker) error {

	web.GET("/proc/reddit/submissions/extract", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.submissions.extract")
		defer span.End()

		// Process reddit submissions
		count, err := extractRedditSubmissionsTickers(readCtx, fsClient, producer, consumer, allTickers, tracer)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	web.GET("/proc/reddit/comments/extract", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.comments.extract")
		defer span.End()

		// Process reddit comment
		count, err := extractRedditCommentsTickers(readCtx, fsClient, producer, consumer, allTickers, tracer)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// web.GET("/proc/ticker_mentions", func(ctx echo.Context) error {

	// 	// Create context for new process
	// 	readCtx := ctx.Request().Context()

	// 	// Create tracer span
	// 	readCtx, span := (*tracer).Start(readCtx, "proc.ticker_mentions")
	// 	defer span.End()

	// 	// Process ticker mentions
	// 	count, err := readTickerMentions(readCtx, psClient, bqClient, tracer)
	// 	if err != nil {
	// 		log.Println(err)
	// 		span.RecordError(err)
	// 		span.SetStatus(codes.Error, "critical error")
	// 		return echo.NewHTTPError(http.StatusInternalServerError)
	// 	}

	// 	response := Response{Message: fmt.Sprintf("Processed %d ticker mentions.", count)}
	// 	return ctx.JSON(http.StatusOK, response)
	// })

	return nil
}
