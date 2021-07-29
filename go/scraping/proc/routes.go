package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VarityPlatform/scraping/common"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, psClient *pubsub.Client, tracer *trace.Tracer, allTickers []common.IEXTicker) error {

	web.GET("/reddit/submissions", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.submissions")
		defer span.End()

		// Process reddit submissions
		count, err := readRedditSubmission(readCtx, psClient, allTickers)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	web.GET("/reddit/comments", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := ctx.Request().Context()

		// Create tracer span
		readCtx, span := (*tracer).Start(readCtx, "proc.reddit.comments")
		defer span.End()

		// Process reddit comment
		count, err := readRedditComment(readCtx, psClient, allTickers)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "critical error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := Response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	return nil
}
