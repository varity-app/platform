package live

import (
	"github.com/google/wire"

	"github.com/varity-app/platform/scraping/internal/scrapers"
)

//SuperSet is the wire superset for the scrapers module
var SuperSet = wire.NewSet(scrapers.NewMemory, NewRedditCommentsScraper, NewRedditSubmissionsScraper)
