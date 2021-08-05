package main

import (
	"net/http"
	"os"
)

// IexBaseURL is the base URL for the IEX Cloud API
const IexBaseURL string = "https://cloud.iexapis.com/stable"

// Send http GET request to IEX Cloud
func iexRequest(client *http.Client, path string, token string) (*http.Response, error) {
	// Create request
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Add token parameter
	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	return client.Do(req)
}

// Entrypoint method
func main() {
	token := os.Getenv("IEX_TOKEN")
	client := &http.Client{}
	db := initPostgres()

	tickersETL(db, client, token)
}
