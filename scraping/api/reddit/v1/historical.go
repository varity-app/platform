package redditv1

// HistoricalRequest is the request payload for a historical reddit scraping request.
type HistoricalRequest struct {
	Subreddit string `json:"subreddit"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Limit     int    `json:"limit"`
}
