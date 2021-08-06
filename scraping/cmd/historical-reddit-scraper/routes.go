package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/VarityPlatform/scraping/scrapers/reddit/historical"
)

// ScrapingRequestBody is the definition of a scraping HTTP request body
type ScrapingRequestBody struct {
	Subreddit string `json:"subreddit"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Limit     int    `json:"limit"`
}

// ISOTimeFormat is an example time format (ISO 8601) for time.Parse
const ISOTimeFormat string = "2000-01-01T00:00:00Z"

// Initialize routes
func initRoutes(web *echo.Echo, submissionsScraper *historical.SubmissionsScraper, commentsScraper *historical.CommentsScraper) {

	web.POST("/scraping/reddit/historical/submissions", func(c echo.Context) error {

		body := new(ScrapingRequestBody)
		if err := c.Bind(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
		}

		// Parse time fields
		before, err := time.Parse(time.RFC3339, body.Before)
		if err != nil {
			log.Println(err, before)
			return echo.NewHTTPError(http.StatusBadRequest, "Field `before` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		after, err := time.Parse(time.RFC3339, body.After)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Field `after` is not of valid date format (`YYYY-MM-DDThh:mm:ssZ`)")
		}

		ctx := c.Request().Context()

		submissions, err := submissionsScraper.Scrape(ctx, body.Subreddit, before, after, body.Limit)
		if err != nil {
			log.Printf("submissionsScraper.Scrape: %v", err)
			return err
		}

		for _, submission := range submissions {
			log.Println(submission.SubmissionId, submission.Timestamp.AsTime(), submission.Title)
		}

		return nil

	})

}
