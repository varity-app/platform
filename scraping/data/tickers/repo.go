package tickers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-pg/pg/v10"
)

// TickerRepo is an abstraction class for interfacing with the Postgres database
// in which tickers are stored.
type IEXTickerRepo struct {
	db *pg.DB
}

// NewIEXTickerRepo contructs a new IEXTickerRepo.
func NewIEXTickerRepo(db *pg.DB) *IEXTickerRepo {
	return &IEXTickerRepo{
		db: db,
	}
}

// Close a repo's connection to the database.
func (r *IEXTickerRepo) Close() error {
	return r.db.Close()
}

// List returns all tickers listed in the NYSE and NASDAQ exchanges.
func (r *IEXTickerRepo) List() ([]IEXTicker, error) {

	// Fetch list of tickers from postgres
	var allTickers []IEXTicker
	_, err := r.db.Query(&allTickers, `SELECT symbol, short_name FROM tickers WHERE exchange IN ('NYS', 'NAS')`)
	if err != nil {
		return nil, fmt.Errorf("pg.FetchTickers: %v", err)
	}

	return allTickers, nil
}

// Save saves a list of tickers to the database.  TODO: this method need integrated with error channels.
func (r *IEXTickerRepo) Save(ctx context.Context, tickers []IEXTicker, threads int) error {
	// Create channel
	tickersChan := make(chan IEXTicker)

	go func() {
		defer close(tickersChan)
		for _, ticker := range tickers {
			select {
			case tickersChan <- ticker:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Create workers
	wg := new(sync.WaitGroup)
	var workerErr error
	for i := 0; i < threads; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()

			// Read tickers from channel
			for ticker := range tickersChan {

				if workerErr != nil {
					return
				}

				// Populate scraped_on column
				ticker.ScrapedOn = time.Now()

				// Save to postgres
				select {
				case <-ctx.Done():
					return
				default:
					_, err := r.db.Model(&ticker).
						OnConflict("(symbol) DO UPDATE").
						Set("name = EXCLUDED.name").
						Set("short_name = EXCLUDED.short_name").
						Set("scraped_on = EXCLUDED.scraped_on").
						Insert()
					if err != nil {
						workerErr = fmt.Errorf("postgres.Insert: %v", err)
						return
					}
				}
			}
		}()
	}

	// Wait for all messages to publish
	wg.Wait()

	return workerErr
}
