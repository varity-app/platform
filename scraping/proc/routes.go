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

	// Extract tickers from reddit submissions
	web.GET("/proc/reddit/submissions/extract", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.submissions.extract")
		defer span.End()

		// Process reddit submissions
		count, err := extractTickersKafka(readCtx, fsClient, producer, consumer, allTickers, tracer, common.REDDIT_SUBMISSIONS, handleKafkaSubmissions)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Sink reddit submissions to bigquery
	web.GET("/proc/reddit/submissions/sink", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.submissions.sink")
		defer span.End()

		// Process reddit submissions
		count, err := sinkKafkaToBigquery(readCtx, fsClient, bqClient, consumer, tracer, common.REDDIT_SUBMISSIONS, common.BIGQUERY_TABLE_REDDIT_SUBMISSIONS, kafkaToRedditSubmission)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Extract ticker mentions from reddit comments
	web.GET("/proc/reddit/comments/extract", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.comments.extract")
		defer span.End()

		// Process reddit comment
		count, err := extractTickersKafka(readCtx, fsClient, producer, consumer, allTickers, tracer, common.REDDIT_COMMENTS, handleKafkaComments)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Sink reddit comments to bigquery
	web.GET("/proc/reddit/comments/sink", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.comments.sink")
		defer span.End()

		// Process reddit comments
		count, err := sinkKafkaToBigquery(readCtx, fsClient, bqClient, consumer, tracer, common.REDDIT_COMMENTS, common.BIGQUERY_TABLE_REDDIT_COMMENTS, kafkaToRedditComment)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Sink ticker mentions to bigquery
	web.GET("/proc/tickerMentions/sink", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.tickerMentions.sink")
		defer span.End()

		// Process reddit comments
		count, err := sinkKafkaToBigquery(readCtx, fsClient, bqClient, consumer, tracer, common.TICKER_MENTIONS, common.BIGQUERY_TABLE_TICKER_MENTIONS, kafkaToTickerMention)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d ticker mentions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	return nil
}
