package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/data/prices"
	"github.com/VarityPlatform/scraping/data/tickers"
	"github.com/VarityPlatform/scraping/scrapers/tiingo"
)

const numThreads int = 10

func setupRoutes(web *echo.Echo, allTickers []tickers.IEXTicker) {
	datasetName := common.BigqueryDatasetScraping + "_" + viper.GetString("deploymentMode")
	token := os.Getenv("TIINGO_TOKEN")

	// Scrape historical stock prices from Tiingo
	web.GET("/scraping/tiingo/prices/:start_date", func(c echo.Context) error {

		// Get request context
		ctx := c.Request().Context()

		// Parse start and end date from request body
		var startDate time.Time
		if c.Param("start_date") == "1d" {
			startDate = time.Now().Add(-1 * 24 * time.Hour)
		} else if c.Param("start_date") == "3d" {
			startDate = time.Now().Add(-3 * 24 * time.Hour)
		} else if c.Param("start_date") == "7d" {
			startDate = time.Now().Add(-7 * 24 * time.Hour)
		} else {
			d, err := time.Parse(tiingo.DateFormat, c.Param("start_date"))
			if err != nil {
				log.Printf("time.Parse: %v", err)
				return echo.ErrBadRequest
			}
			startDate = d
		}
		endDate := time.Now()

		// Initialize scraper and sink
		scraper := tiingo.NewEODPricesScraper(token)
		sink, err := prices.NewEODPriceSink(ctx, prices.EODPriceSinkOpts{
			DatasetName: datasetName,
			TableName:   common.BigqueryTableEODPrices,
		})
		if err != nil {
			log.Printf("sink.New: %v", err)
			return err
		}
		defer sink.Close()

		// Spawn off a certain number of threads
		var errcs []<-chan error
		for i := 0; i < numThreads; i++ {
			chunkSize := len(allTickers)/numThreads + 1
			tickersChunk := allTickers[i*chunkSize : min((i+1)*chunkSize, len(allTickers))]

			scrapeErrc, sinkErrc := pricesPipeline(ctx, scraper, sink, tickersChunk, startDate, endDate)
			errcs = append(errcs, scrapeErrc, sinkErrc)
		}

		// Wait for pipelines to finish
		if err := common.WaitForPipeline(errcs...); err != nil {
			log.Println(err)
			return err
		}

		return c.JSON(200, map[string]string{})
	})
}
