package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VarityPlatform/scraping/common"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, processor *KafkaTickerProcessor, sink *KafkaBigquerySink, tracer *trace.Tracer, allTickers []common.IEXTicker) error {

	// Extract tickers from reddit submissions
	web.GET("/scraping/proc/reddit/submissions/extract", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.submissions.extract")
		defer span.End()

		// Process reddit submissions
		count, err := processor.ProcessTopic(readCtx, common.RedditSubmissions, handleSubmissionTickerProcessing)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Sink reddit submissions to bigquery
	web.GET("/scraping/proc/reddit/submissions/sink", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.submissions.sink")
		defer span.End()

		// Process reddit submissions
		count, err := sink.SinkTopic(readCtx, common.RedditSubmissions, common.BigqueryTableRedditSubmissions, kafkaToRedditSubmission)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Extract tickers from reddit comments
	web.GET("/scraping/proc/reddit/comments/extract", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.comments.extract")
		defer span.End()

		// Process reddit comments
		count, err := processor.ProcessTopic(readCtx, common.RedditComments, handleCommentTickerProcessing)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Sink reddit comments to bigquery
	web.GET("/scraping/proc/reddit/comments/sink", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.comments.sink")
		defer span.End()

		// Process reddit comments
		count, err := sink.SinkTopic(readCtx, common.RedditComments, common.BigqueryTableRedditComments, kafkaToRedditComment)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	// Sink ticker mentions to bigquery
	web.GET("/scraping/proc/tickerMentions/sink", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.tickerMentions.sink")
		defer span.End()

		// Process reddit comments
		count, err := sink.SinkTopic(readCtx, common.TickerMentions, common.BigqueryTableTickerMentions, kafkaToTickerMention)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d ticker mentions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	return nil
}
