package main

import (
	"regexp"

	"github.com/VarityPlatform/scraping/common"
	"github.com/go-pg/pg/v10"
)

const TICKER_SPAM_THRESHOLD int = 5

var urls *regexp.Regexp = regexp.MustCompile(`https?:\/\/.*[\r\n]*`)
var alphaNumeric *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 \n\.]`)
var capsSpam *regexp.Regexp = regexp.MustCompile(`([A-Z]{1,3}.?[A-Z]{1,3})(\W[A-Z].?[A-Z]+)+`)
var tickerRegex *regexp.Regexp = regexp.MustCompile(`[A-Z][A-Z0-9.]*[A-Z0-9]`)
var wordRegex *regexp.Regexp = regexp.MustCompile(`\w+`)

// Extract tickers from string
func extractTickersString(s string, tickerList []common.IEXTicker) []common.IEXTicker {

	// Remove urls
	s = urls.ReplaceAllString(s, "")

	// Remove alphanumeric characters
	s = alphaNumeric.ReplaceAllString(s, "")

	// Remove caps spam
	s = capsSpam.ReplaceAllString(s, "")

	// Find ticker look-a-likes
	tickers := tickerRegex.FindAllString(s, -1)

	// Cross reference with provided list of valid tickers
	validTickers := []common.IEXTicker{}
	for _, ticker := range tickers {
		for _, realTicker := range tickerList {
			if ticker == realTicker.Symbol {
				validTickers = append(validTickers, realTicker)
			}
		}
	}

	// If there are more than N tickers, count as spam and discard
	if len(validTickers) > TICKER_SPAM_THRESHOLD {
		return []common.IEXTicker{}
	}

	return validTickers
}

// Extract short name mentions from strings
func extractShortNamesString(s string, tickerList []common.IEXTicker) []common.IEXTicker {
	// Extract all alphanumeric words
	words := wordRegex.FindAllString(s, -1)

	// Find mentioned tickers
	mentionedTickers := []common.IEXTicker{}
	for _, ticker := range tickerList {

		if ticker.ShortName == "" {
			continue
		}

		for _, word := range words {
			if word == ticker.ShortName {
				mentionedTickers = append(mentionedTickers, ticker)
			}
		}
	}

	return mentionedTickers
}

// Get frequency counts of tickers
func calcTickerFrequency(tickerList []common.IEXTicker) ([]common.IEXTicker, map[string]int) {
	frequencies := make(map[string]int)
	uniqueTickers := []common.IEXTicker{}

	for _, ticker := range tickerList {
		if frequencies[ticker.Symbol] == 0 {
			uniqueTickers = append(uniqueTickers, ticker)
		}
		frequencies[ticker.Symbol] += 1
	}

	return uniqueTickers, frequencies
}

// Fetch tickers from postgres
func fetchTickers(db *pg.DB) ([]common.IEXTicker, error) {
	var tickers []common.IEXTicker
	_, err := db.Query(&tickers, `SELECT symbol, short_name FROM tickers WHERE exchange IN ('NYS', 'NAS')`)
	return tickers, err
}
