package main

import (
	"fmt"
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
		startDate, err := time.Parse(tiingo.DateFormat, c.Param("start_date"))
		if err != nil {
			log.Printf("time.Parse: %v", err)
			return echo.ErrBadRequest
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

			tickersChan := make(chan tickers.IEXTicker)
			priceBatches := make(chan []prices.EODPrice)

			scrapeErrc := make(chan error, 1)
			sinkErrc := make(chan error, 1)

			// Create scraping thread
			go func() {
				defer close(scrapeErrc)
				defer close(priceBatches)

				for ticker := range tickersChan {

					// Fetch prices from Tiingo
					scrapedPrices, err := scraper.Scrape(ctx, ticker.Symbol, startDate, endDate)
					if err != nil {
						scrapeErrc <- fmt.Errorf("scraper.Scrape: %v", err)
						return
					}
					if len(scrapedPrices) == 0 {
						continue
					}

					// Send batch to output channel
					select {
					case <-ctx.Done():
						return
					case priceBatches <- scrapedPrices:
					}
				}
			}()

			// Create sink thread
			go func() {
				defer close(sinkErrc)

				for batch := range priceBatches {
					select {
					case <-ctx.Done():
						return
					default:
						// Save prices to BigQuery
						err := sink.Sink(ctx, batch)
						if err != nil {
							sinkErrc <- fmt.Errorf("sink.Sink: %v", err)
							return
						}
					}
				}
			}()

			// Send tickers to the channel
			go func(i int) {
				defer close(tickersChan)
				chunkSize := len(allTickers)/numThreads + 1

				// Send a chunk of tickers to the ticker channel
				for _, ticker := range allTickers[i*chunkSize : min((i+1)*chunkSize, len(allTickers))] {
					select {
					case <-ctx.Done():
						return
					case tickersChan <- ticker:
					}
				}
			}(i)

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
