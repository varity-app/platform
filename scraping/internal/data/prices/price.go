package prices

import (
	"time"
)

// EODPrice is a historical EOD price entry for a specific stock.
// For more information on each field, visit https://api.tiingo.com/documentation/end-of-day
type EODPrice struct {
	Symbol    string    `json:"symbol" bigquery:"symbol"`
	Date      time.Time `json:"date" bigquery:"date"`
	Open      float32   `json:"open" bigquery:"open"`
	High      float32   `json:"high" bigquery:"high"`
	Low       float32   `json:"low" bigquery:"low"`
	Close     float32   `json:"close" bigquery:"close"`
	Volume    int64     `json:"volume" bigquery:"volume"`
	AdjOpen   float32   `json:"adjOpen" bigquery:"adj_open"`
	AdjHigh   float32   `json:"adjHigh" bigquery:"adj_high"`
	AdjLow    float32   `json:"adjLow" bigquery:"adj_low"`
	AdjClose  float32   `json:"adjClose" bigquery:"adj_close"`
	AdjVolume int64     `json:"adjVolume" bigquery:"adj_volume"`
	Dividend  float32   `json:"divCash" bigquery:"dividend"`
	Split     float32   `json:"splitFactor" bigquery:"split"`

	// Custom
	ScrapedOn time.Time `json:"-" bigquery:"scraped_on"`
}
