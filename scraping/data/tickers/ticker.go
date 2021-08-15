package tickers

import "time"

// IEXTicker is an object used by the ORM for storing tickers fetched from the
// IEX Cloud API
type IEXTicker struct {
	tableName    struct{}  `json:"-" pg:"tickers"` //lint:ignore U1000 PG table name does not need a special name
	Symbol       string    `pg:"symbol,pk"`
	Exchange     string    `pg:"exchange"`
	ExchangeName string    `pg:"exchange_name"`
	Name         string    `pg:"name"`
	ScrapedOn    time.Time `pg:"scraped_on"`
	Type         string    `pg:"type"`
	Region       string    `pg:"region"`
	Currency     string    `pg:"currency"`
	IsEnabled    bool      `pg:"is_enabled"`
	Cik          string    `pg:"cik"`

	// Custom fields
	ShortName string `pg:"short_name"`
}
