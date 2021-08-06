package main

import (
	"context"
	"log"
	"os"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/labstack/echo/v4"

	"github.com/spf13/viper"
)

// Entrypoint method
func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Initialize the reddit scrapers
	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}
	submissionsScraper, err := initSubmissionsScraper(
		ctx,
		credentials,
		scrapers.MemoryOpts{CollectionName: common.RedditSubmissions + "-v2-" + viper.GetString("deploymentMode")},
	)
	if err != nil {
		log.Fatalf("initSubmissionsScraper: %v", err)
	}
	defer submissionsScraper.Close()

	commentsScraper, err := initCommentsScraper(
		ctx,
		credentials,
		scrapers.MemoryOpts{CollectionName: common.RedditComments + "-v2-" + viper.GetString("deploymentMode")},
	)
	if err != nil {
		log.Fatalf("initCommentsScraper: %v", err)
	}
	defer commentsScraper.Close()

	// Create trace exporter
	exporter, err := texporter.NewExporter(texporter.WithProjectID(common.GCPProjectID))
	if err != nil {
		log.Fatalf("texporter.NewExporter: %v", err)
	}

	// Create trace provider with the exporter
	traceProbability := viper.GetFloat64("tracing.probability")
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(traceProbability)),
	)
	defer tp.ForceFlush(ctx) // flushes any pending spans
	otel.SetTracerProvider(tp)
	tracer := otel.GetTracerProvider().Tracer("varity.app/scraping")

	// Init kafka client
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     os.Getenv("KAFKA_AUTH_KEY"),
		"sasl.password":     os.Getenv("KAFKA_AUTH_SECRET"),
	})
	if err != nil {
		log.Fatalf("kafka.NewProducer: %v", err)
	}
	defer producer.Close()

	// Initialize webserver
	web := echo.New()
	err = setupRoutes(web, submissionsScraper, commentsScraper, producer, &tracer)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
