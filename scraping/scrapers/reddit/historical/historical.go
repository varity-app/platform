package historical

import (
	"github.com/google/wire"

	"github.com/VarityPlatform/scraping/scrapers"
)

//SuperSet is the wire superset for the scrapers module
var SuperSet = wire.NewSet(scrapers.NewMemory, NewSubmissionsScraper, NewCommentsScraper)
