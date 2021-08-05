package scrapers

import (
	"github.com/google/wire"
)

//SuperSet is the wire superset for the scrapers module
var SuperSet = wire.NewSet(NewMemory, NewRedditCommentsScraper, NewRedditSubmissionsScraper)
