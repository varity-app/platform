package main

import (
	"context"
	"log"

	"github.com/VarityPlatform/scraping/common"
	"github.com/spf13/viper"

	"cloud.google.com/go/pubsub"
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

	// Initialize postgres
	db := common.InitPostgres()

	// Fetch list of tickers from postgres
	tracer := otel.GetTracerProvider().Tracer("varity.app/proc")
	ctx, span := tracer.Start(ctx, "proc.fetchTickers")
	allTickers, err := fetchTickers(db)
	if err != nil {
		log.Fatalf("pg.FetchTickers: %v", err)
	}
	log.Println(len(allTickers))
	span.End()

	// Initialize pubsub client
	psClient, err := pubsub.NewClient(ctx, common.GCP_PROJECT_ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer psClient.Close()

	// Initialize webserver
	web := echo.New()
	err = setupRoutes(web, psClient, &tracer, allTickers)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))

}
