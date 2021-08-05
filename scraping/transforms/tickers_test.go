package transforms

import (
	"testing"

	"github.com/VarityPlatform/scraping/common"
)

// TestTickersOne is a unit test
func TestTickersOne(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy GME you cucks.")
	if len(results) != 1 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 1)
	}

	result := results[0]
	if result.Symbol != "GME" {
		t.Errorf("Ticker symbol incorrect, got: %s, want: %s", result.Symbol, "GME")
	}
}

// TestTickersTwo is a unit test
func TestTickersTwo(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
		{Symbol: "GOOGL"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy GME you cucks.  Also buy GOOGL!")
	if len(results) != 2 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 2)
	}

	result := results[1]
	if result.Symbol != "GOOGL" {
		t.Errorf("Ticker symbol incorrect, got: %s, want: %s", result.Symbol, "GOOGL")
	}
}

// TestTickersTwoSameSymbol is a unit test
func TestTickersTwoSameSymbol(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy GME you cucks.  Also buy more GME!")
	if len(results) != 2 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 2)
	}

	result := results[1]
	if result.Symbol != "GME" {
		t.Errorf("Ticker symbol incorrect, got: %s, want: %s", result.Symbol, "GME")
	}
}

// TestTickersTwoSameSymbolFrequency is a unit test
func TestTickersTwoSameSymbolFrequency(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy GME you cucks.  Also buy more GME!")
	if len(results) != 2 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 2)
	}

	uniqTickers, frequencies := extractor.CalcTickerFrequency(results)
	if len(uniqTickers) != 1 {
		t.Errorf("Number of unique tickers was incorrect, got: %d, want: %d", len(results), 1)
	}

	if frequencies["GME"] != 2 {
		t.Errorf("Ticker frequency was incorrect, got: %d, want: %d", frequencies["GME"], 2)
	}
}

// TestTickerFrequencies is a unit test
func TestTickerFrequencies(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
		{Symbol: "GOOGL"},
		{Symbol: "AAPL"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy GME you cucks.  Also buy more GME!  Noone likes GOOGL anymore.  They suck.  AAPL is okay.")
	if len(results) != 4 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 2)
	}

	uniqTickers, frequencies := extractor.CalcTickerFrequency(results)
	if len(uniqTickers) != 3 {
		t.Errorf("Number of unique tickers was incorrect, got: %d, want: %d", len(results), 1)
	}

	if frequencies["GME"] != 2 {
		t.Errorf("Ticker frequency was incorrect, got: %d, want: %d", frequencies["GME"], 2)
	} else if frequencies["GOOGL"] != 1 {
		t.Errorf("Ticker frequency was incorrect, got: %d, want: %d", frequencies["GOOGL"], 1)
	} else if frequencies["AAPL"] != 1 {
		t.Errorf("Ticker frequency was incorrect, got: %d, want: %d", frequencies["AAPL"], 1)
	}
}

// TestTickersNotFound is a unit test
func TestTickersNotFound(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy AAPL you cucks.")
	if len(results) != 0 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 0)
	}
}

// TestTickersOneShortName is a unit test
func TestTickersOneShortName(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME", ShortName: "Gamestop"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractShortNames("Buy AAPL you cucks.  Noway it's better than Gamestop!", tickerList)
	if len(results) != 1 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 1)
	}
}

// TestTickerInURL is a unit test
func TestTickerInURL(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy from https://google.com/search?query=GME!")
	if len(results) != 0 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 0)
	}
}

// TestTickersSpam is a unit test
func TestTickersSpam(t *testing.T) {
	tickerList := []common.IEXTicker{
		{Symbol: "GME"},
		{Symbol: "AAPL"},
		{Symbol: "GOOGL"},
		{Symbol: "EHTH"},
		{Symbol: "MSFT"},
		{Symbol: "TLRY"},
	}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy all of these $GME $AAPL $EHTH $MSFT $TLRY!")
	if len(results) != 0 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 0)
	}
}

// TestNoTickersNotFound is a unit test
func TestNoTickersNotFound(t *testing.T) {
	tickerList := []common.IEXTicker{}

	extractor := NewTickerExtractor(tickerList)

	results := extractor.ExtractTickers("Buy AAPL you cucks.")
	if len(results) != 0 {
		t.Errorf("Number of tickers was incorrect, got: %d, want: %d", len(results), 0)
	}
}
