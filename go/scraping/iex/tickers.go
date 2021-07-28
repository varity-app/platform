package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/VarityPlatform/scraping/common"
	"github.com/go-pg/pg/v10"
)

const SHORT_NAME_FREQUENCY uint32 = 10
const NUM_PG_WORKERS int = 12

// Fetch tickers form IEX Cloud
func extractTickers(client *http.Client, token string) []common.IEXTicker {
	// Send request to IEX cloud
	resp, err := iexRequest(client, IEX_BASE_URL+"/ref-data/symbols", token)
	if err != nil {
		log.Fatalln("Error fetching tickers from IEXCloud:", err.Error())
	}
	defer resp.Body.Close()

	// Read respose body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err.Error())
	}

	// log.Println(string(body))

	tickers := []common.IEXTicker{}
	err = json.Unmarshal(body, &tickers)
	if err != nil {
		log.Fatalln("Error unmarshaling ticker response:", err.Error())
	}

	return tickers

}

// Extract short_name column for each ticker
func addShortName(tickers []common.IEXTicker) []common.IEXTicker {
	// Generate list of all ticker names
	names := []string{}
	for _, ticker := range tickers {
		names = append(names, ticker.Name)
	}

	// Concatenate all words in names into one giant list
	nameString := strings.Join(names, " ")
	words := strings.Split(nameString, " ")

	// Create list of alphanumeric words
	alphaRegex := regexp.MustCompile(`[^a-zA-Z\d\s:]`)
	filteredWords := []string{}
	for _, word := range words {
		if !alphaRegex.MatchString(word) {
			filteredWords = append(filteredWords, word)
		}
	}

	// Count frequencies of each word
	frequencies := make(map[string]uint32)
	for _, word := range filteredWords {
		frequencies[word] += 1
	}

	// Assign short name to a word if it's infrequent
	for idx, ticker := range tickers {
		nameWords := strings.Split(ticker.Name, " ")

		for _, word := range nameWords {
			if !alphaRegex.MatchString(word) && frequencies[word] <= SHORT_NAME_FREQUENCY {
				tickers[idx].ShortName = word
				break
			}
		}
	}

	return tickers
}

// Save IEX Tickers to Postgres
func loadTickers(db *pg.DB, tickers []common.IEXTicker) error {
	// Create channel
	tickersChan := make(chan common.IEXTicker)

	go func() {
		for _, ticker := range tickers {
			tickersChan <- ticker
		}
		close(tickersChan)
	}()

	// Create workers
	wg := new(sync.WaitGroup)
	for i := 0; i < NUM_PG_WORKERS; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()

			// Read tickers from channel
			for ticker := range tickersChan {

				// Populate scraped_on column
				ticker.ScrapedOn = time.Now()

				// Save to postgres
				_, err := db.Model(&ticker).
					OnConflict("(symbol) DO UPDATE").
					Set("name = EXCLUDED.name").
					Set("short_name = EXCLUDED.short_name").
					Set("scraped_on = EXCLUDED.scraped_on").
					Insert()
				if err != nil {
					log.Fatalln("Error writing to postgres:", err.Error())
				}
			}
		}()
	}

	// Wait for all messages to publish
	wg.Wait()

	return nil
}

// Run entire tickers ETL
func tickersETL(db *pg.DB, client *http.Client, token string) {
	// Fetch and save tickers
	log.Println("Fetching tickers...")
	tickers := extractTickers(client, token)

	// Transform tickers
	log.Println("Transforming tickers...")
	tickers = addShortName(tickers)

	// Load tickers into postgres
	log.Println("Loading", len(tickers), "tickers...")
	err := loadTickers(db, tickers)
	if err != nil {
		log.Fatalln("Error loading tickers:", err.Error())
	}
}
