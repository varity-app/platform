package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/data/kafka"
	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"
)

type response struct {
	Message string `json:"message"`
}

// Set up echo routes
func setupRoutes(web *echo.Echo, offsetOpts kafka.OffsetManagerOpts, allTickers []common.IEXTicker) error {
	datasetName := common.BigqueryDatasetScraping + "_" + viper.GetString("deploymentMode")

	// Kafka auth
	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	username := os.Getenv("KAFKA_AUTH_KEY")
	password := os.Getenv("KAFKA_AUTH_SECRET")

	// Extract tickers from reddit submissions
	web.GET("/scraping/proc/reddit/submissions/extract", func(c echo.Context) error {

		// Create context for new process
		ctx := c.Request().Context()

		// Initialize processor
		processor, err := initProcessor(ctx, kafka.ProcessorOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			InputTopic:       common.RedditSubmissions,
			OutputTopic:      common.TickerMentions,
			CheckpointKey:    common.RedditSubmissions + "-proc",
		}, offsetOpts)
		if err != nil {
			log.Fatalf("kafkaTickerProcessor.New: %v", err)
		}
		defer processor.Close()

		// Initialize handler
		handler := NewRedditSubmissionHandler(allTickers)

		// Process reddit submissions
		count, err := processor.ProcessTopic(ctx, handler)
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

		// Initialize bigquery sink
		sink, err := initSink(ctx, kafka.BigquerySinkOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			InputTopic:       common.RedditSubmissions,
			CheckpointKey:    common.RedditSubmissions + "-sink",
			DatasetName:      datasetName,
			TableName:        common.BigqueryTableRedditSubmissions,
		}, offsetOpts)
		if err != nil {
			log.Fatalf("kafkaTickerProcessor.New: %v", err)
		}
		defer sink.Close()

		// Process reddit submissions
		count, err := sink.SinkTopic(ctx, SinkBatchSize, convertKafkaToSubmission)
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

		processor, err := initProcessor(ctx, kafka.ProcessorOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			InputTopic:       common.RedditComments,
			OutputTopic:      common.TickerMentions,
			CheckpointKey:    common.RedditComments + "-proc",
		}, offsetOpts)
		if err != nil {
			log.Fatalf("kafkaTickerProcessor.New: %v", err)
		}
		defer processor.Close()

		// Initialize handler
		handler := NewRedditCommentHandler(allTickers)

		// Process reddit submissions
		count, err := processor.ProcessTopic(ctx, handler)
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

		// Initialize bigquery sink
		sink, err := initSink(ctx, kafka.BigquerySinkOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			InputTopic:       common.RedditComments,
			CheckpointKey:    common.RedditComments + "-sink",
			DatasetName:      datasetName,
			TableName:        common.BigqueryTableRedditComments,
		}, offsetOpts)
		if err != nil {
			log.Fatalf("kafkaTickerProcessor.New: %v", err)
		}
		defer sink.Close()

		// Process reddit comments
		count, err := sink.SinkTopic(ctx, SinkBatchSize, convertKafkaToComment)
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

		// Initialize bigquery sink
		sink, err := initSink(ctx, kafka.BigquerySinkOpts{
			BootstrapServers: bootstrapServers,
			Username:         username,
			Password:         password,
			InputTopic:       common.TickerMentions,
			CheckpointKey:    common.TickerMentions + "-sink",
			DatasetName:      datasetName,
			TableName:        common.BigqueryTableTickerMentions,
		}, offsetOpts)
		if err != nil {
			log.Fatalf("kafkaTickerProcessor.New: %v", err)
		}
		defer sink.Close()

		// Process reddit submissions
		count, err := sink.SinkTopic(ctx, SinkBatchSize, convertKafkaToTickerMention)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		response := response{Message: fmt.Sprintf("Processed %d ticker mentions.", count)}
		return c.JSON(http.StatusOK, response)
	})

	return nil
}
