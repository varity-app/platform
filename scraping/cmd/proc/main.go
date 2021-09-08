package main

import (
	"log"

	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/data/kafka"

	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"
)

// Entrypoint method
func main() {
	initConfig()

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

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	err = setupRoutes(web, offsetOpts, allTickers)
	if err != nil {
		log.Fatalln(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))

}
