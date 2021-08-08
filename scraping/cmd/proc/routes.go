package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/data/kafka"
	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"
)

type response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, processor *kafka.Processor, sink *kafka.BigquerySink, allTickers []common.IEXTicker) error {
	datasetName := common.BigqueryDatasetScraping + "_" + viper.GetString("deploymentMode")

	// Extract tickers from reddit submissions
	web.GET("/scraping/proc/reddit/submissions/extract", func(c echo.Context) error {

		// Create context for new process
		ctx := c.Request().Context()

		// Initialize handler
		handler := NewRedditSubmissionHandler(allTickers)
		checkpointKey := common.RedditSubmissions + "-proc"

		// Process reddit submissions
		count, err := processor.ProcessTopic(ctx, common.RedditSubmissions, checkpointKey, handler)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return c.JSON(http.StatusOK, response)
	})

	// Sink reddit submissions to bigquery
	web.GET("/scraping/proc/reddit/submissions/sink", func(c echo.Context) error {

		// Create context for new process
		ctx := c.Request().Context()

		// Process reddit submissions
		checkpointKey := common.RedditSubmissions + "-proc"
		count, err := sink.SinkTopic(ctx, common.RedditSubmissions, checkpointKey, datasetName, common.BigqueryTableRedditSubmissions, SinkBatchSize, convertKafkaToSubmission)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit submissions.", count)}
		return c.JSON(http.StatusOK, response)
	})

	// Extract tickers from reddit comments
	web.GET("/scraping/proc/reddit/comments/extract", func(c echo.Context) error {

		// Create context for new process
		ctx := c.Request().Context()

		// Initialize handler
		handler := NewRedditCommentHandler(allTickers)
		checkpointKey := common.RedditComments + "-proc"

		// Process reddit submissions
		count, err := processor.ProcessTopic(ctx, common.RedditComments, checkpointKey, handler)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return c.JSON(http.StatusOK, response)
	})

	// Sink reddit comments to bigquery
	web.GET("/scraping/proc/reddit/comments/sink", func(c echo.Context) error {

		// Create context for new process
		ctx := c.Request().Context()

		// Process reddit comments
		checkpointKey := common.RedditComments + "-proc"
		count, err := sink.SinkTopic(ctx, common.RedditComments, checkpointKey, datasetName, common.BigqueryTableRedditComments, SinkBatchSize, convertKafkaToComment)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d reddit comments.", count)}
		return c.JSON(http.StatusOK, response)
	})

	// Sink ticker mentions to bigquery
	web.GET("/scraping/proc/tickerMentions/sink", func(c echo.Context) error {

		// Create context for new process
		ctx := c.Request().Context()
		// Process reddit comments
		checkpointKey := common.TickerMentions + "-proc"
		count, err := sink.SinkTopic(ctx, common.TickerMentions, checkpointKey, datasetName, common.BigqueryTableTickerMentions, SinkBatchSize, convertKafkaToTickerMention)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d ticker mentions.", count)}
		return c.JSON(http.StatusOK, response)
	})

	return nil
}
