package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/VarityPlatform/scraping/common"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, psClient *pubsub.Client, allTickers []common.IEXTicker) error {

	web.GET("/reddit/submissions", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := context.Background()
		readCtx, cancel := context.WithCancel(readCtx)

		// Process reddit submissions
		count, err := readRedditSubmission(readCtx, psClient, allTickers)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Cancel process
		cancel()

		response := Response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	web.GET("/reddit/comments", func(ctx echo.Context) error {

		// Create context for new process
		readCtx := context.Background()
		readCtx, cancel := context.WithCancel(readCtx)

		// Process reddit comment
		count, err := readRedditComment(readCtx, psClient, allTickers)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Cancel process
		cancel()

		response := Response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return ctx.JSON(http.StatusOK, response)
	})

	return nil
}
