package main

import (
	"context"
	"log"

	"github.com/VarityPlatform/scraping/common"

	"github.com/spf13/viper"

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
	tracer := otel.GetTracerProvider().Tracer("varity.app/proc")

	// Initialize postgres
	db := common.InitPostgres()

	// Fetch list of tickers from postgres
	ctx, span := tracer.Start(ctx, "proc.fetchTickers")
	var allTickers []common.IEXTicker
	_, err = db.Query(&allTickers, `SELECT symbol, short_name FROM tickers WHERE exchange IN ('NYS', 'NAS')`)
	if err != nil {
		log.Fatalf("pg.FetchTickers: %v", err)
	}
	span.End()

	offsetsCollectionName := FirestoreKafkaOffsets + "-" + viper.GetString("deploymentMode")
	processor, err := NewKafkaTickerProcessor(ctx, OffsetManagerOpts{collectionName: offsetsCollectionName}, allTickers)
	if err != nil {
		log.Fatalf("kafkaTickerProcessor.New: %v", err)
	}
	defer processor.Close()

	sink, err := NewKafkaBigquerySink(ctx, OffsetManagerOpts{collectionName: offsetsCollectionName})
	if err != nil {
		log.Fatalf("kafkaBigquerySink.New: %v", err)
	}
	defer sink.Close()

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	err = setupRoutes(web, processor, sink, &tracer, allTickers)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))

}
