package tiingo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/varity-app/platform/scraping/internal/data/prices"
)

// PricesEndpoint is the Tiingo API endpoint for scraping daily prices
const PricesEndpoint string = "https://api.tiingo.com/tiingo/daily/%s/prices"

// DateFormat is the format used for Date fields
const DateFormat string = "2006-01-02"

// EODPricesScraperOpts is a struct containing configuration parameters for NewEODPricesScraper
type EODPricesScraperOpts struct {
	Token string
}

// EODPricesScraper scrapes EOD prices from the Tiingo API
type EODPricesScraper struct {
	token  string
	client *http.Client
}

// NewEODPricesScraper constructs a new EODPricesScraper
func NewEODPricesScraper(token string) *EODPricesScraper {
	return &EODPricesScraper{
		token:  token,
		client: &http.Client{},
	}
}

// Scrape daily prices for a given ticker from the Tiingo API
func (s *EODPricesScraper) Scrape(ctx context.Context, symbol string, startDate, endDate time.Time) ([]prices.EODPrice, error) {

	// Make request
	resp, err := s.request(symbol, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Parse response
	eodPrices, err := s.parse(resp)
	if err != nil {
		return nil, err
	}

	var labeledPrices []prices.EODPrice
	for _, price := range eodPrices {
		price.Symbol = symbol
		price.ScrapedOn = time.Now()
		labeledPrices = append(labeledPrices, price)
	}

	return labeledPrices, nil
}

func (s *EODPricesScraper) request(symbol string, startDate, endDate time.Time) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(PricesEndpoint, symbol), nil)
	if err != nil {
		return nil, err
	}

	// Add token parameter
	q := req.URL.Query()
	q.Add("token", s.token)
	q.Add("startDate", startDate.Format(DateFormat))
	q.Add("endDate", endDate.Format(DateFormat))
	q.Add("resampleFreq", "daily")

	req.URL.RawQuery = q.Encode()

	return s.client.Do(req)
}

// parse a response body for
func (s *EODPricesScraper) parse(resp *http.Response) ([]prices.EODPrice, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("scraper.ReadyBody: %v", err)
	}

	eodPrices := []prices.EODPrice{}
	err = json.Unmarshal(body, &eodPrices)
	if err != nil {
		return []prices.EODPrice{}, nil
	}

	return eodPrices, nil
}
