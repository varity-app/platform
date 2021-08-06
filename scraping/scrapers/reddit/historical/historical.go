package historical

import (
	"github.com/google/wire"

	"github.com/VarityPlatform/scraping/scrapers"
)

// ScraperOpts defines options for initialize a new historical scraper
type ScraperOpts struct {
	ProxyURL string
}

//SuperSet is the wire superset for the scrapers module
var SuperSet = wire.NewSet(scrapers.NewMemory, NewSubmissionsScraper, NewCommentsScraper)
