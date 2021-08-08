package main

import (
	"context"
	"log"
	"os"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/data/kafka"

	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"
)

// Entrypoint method
func main() {
	initConfig()

	ctx := context.Background()

	// Initialize postgres
	db := common.InitPostgres()

	// Fetch list of tickers from postgres
	var allTickers []common.IEXTicker
	_, err := db.Query(&allTickers, `SELECT symbol, short_name FROM tickers WHERE exchange IN ('NYS', 'NAS')`)
	if err != nil {
		log.Fatalf("pg.FetchTickers: %v", err)
	}

	// Initialize kafka structs
	offsetsCollectionName := FirestoreKafkaOffsets + "-" + viper.GetString("deploymentMode")
	offsetOpts := kafka.OffsetManagerOpts{CollectionName: offsetsCollectionName}
	kafkaOpts := kafka.KafkaOpts{
		BootstrapServers: os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		Username:         os.Getenv("KAFKA_AUTH_KEY"),
		Password:         os.Getenv("KAFKA_AUTH_SECRET"),
	}
	processor, err := initProcessor(ctx, kafkaOpts, offsetOpts)
	if err != nil {
		log.Fatalf("kafkaTickerProcessor.New: %v", err)
	}
	defer processor.Close()

	sink, err := initSink(ctx, kafkaOpts, offsetOpts)
	if err != nil {
		log.Fatalf("kafkaBigquerySink.New: %v", err)
	}
	defer sink.Close()

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	err = setupRoutes(web, processor, sink, allTickers)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))

}
