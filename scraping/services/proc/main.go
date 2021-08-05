package main

import (
	"context"
	"log"
	"os"

	"github.com/VarityPlatform/scraping/common"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/spf13/viper"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Entrypoint method
func main() {
	initConfig()

	ctx := context.Background()

	// Create trace exporter
	exporter, err := texporter.NewExporter(texporter.WithProjectID(common.GCP_PROJECT_ID))
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
	tracer := otel.GetTracerProvider().Tracer("varity.app/proc")

	// Initialize postgres
	db := common.InitPostgres()

	// Fetch list of tickers from postgres
	ctx, span := tracer.Start(ctx, "proc.fetchTickers")
	allTickers, err := fetchTickers(db)
	if err != nil {
		log.Fatalf("pg.FetchTickers: %v", err)
	}
	log.Println(len(allTickers))
	span.End()

	// Init firestore client
	fsClient, err := firestore.NewClient(ctx, common.GCP_PROJECT_ID)
	if err != nil {
		log.Fatalf("firestore.GetClient: %v", err)
	}
	defer fsClient.Close()

	// Initialize kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     os.Getenv("KAFKA_AUTH_KEY"),
		"sasl.password":     os.Getenv("KAFKA_AUTH_SECRET"),
	})
	if err != nil {
		log.Fatalf("kafka.GetProducer: %v", err)
	}

	// Initialize kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol":        "SASL_SSL",
		"sasl.mechanisms":          "PLAIN",
		"sasl.username":            os.Getenv("KAFKA_AUTH_KEY"),
		"sasl.password":            os.Getenv("KAFKA_AUTH_SECRET"),
		"group.id":                 "consumers",
		"enable.auto.offset.store": "false",
	})
	if err != nil {
		log.Fatalf("kafka.NewConsumer: %v", err)
	}
	defer consumer.Close()

	// Initialize bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GCP_PROJECT_ID)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize webserver
	web := echo.New()
	err = setupRoutes(web, fsClient, producer, consumer, bqClient, &tracer, allTickers)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))

}
