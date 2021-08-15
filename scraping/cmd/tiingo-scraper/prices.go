package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VarityPlatform/scraping/data/prices"
	"github.com/VarityPlatform/scraping/data/tickers"
	"github.com/VarityPlatform/scraping/scrapers/tiingo"
)

func pricesPipeline(ctx context.Context, scraper *tiingo.EODPricesScraper, sink *prices.EODPriceSink, allTickers []tickers.IEXTicker, startDate, endDate time.Time) (<-chan error, <-chan error) {
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
	go func() {
		defer close(tickersChan)

		// Send a chunk of tickers to the ticker channel
		for _, ticker := range allTickers {
			select {
			case <-ctx.Done():
				return
			case tickersChan <- ticker:
			}
		}
	}()

	return scrapeErrc, sinkErrc
}
