package main

import (
	"context"
	"log"

	"github.com/VarityPlatform/scraping/common"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"

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
	redditClient, err := initReddit()
	if err != nil {
		log.Fatal(err)
	}

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
	tracer := otel.GetTracerProvider().Tracer("varity.app/scraping")

	// Init firestore client
	fsClient, err := firestore.NewClient(ctx, GCP_PROJECT_ID)
	if err != nil {
		log.Fatalf("firestore.GetClient: %v", err)
	}
	defer fsClient.Close()

	// Init pubsub client
	psClient, err := pubsub.NewClient(ctx, GCP_PROJECT_ID)
	if err != nil {
		log.Fatalf("pubsub.GetClient: %v", err)
	}
	defer psClient.Close()

	// Initialize webserver
	web := echo.New()
	err = setupRoutes(web, redditClient, fsClient, psClient, &tracer)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
