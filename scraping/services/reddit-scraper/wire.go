// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/google/wire"
)

func initSubmissionsScraper(ctx context.Context, redditCredentials reddit.Credentials, memoryOpts scrapers.MemoryOpts) (*scrapers.RedditSubmissionsScraper, error) {
	wire.Build(scrapers.SuperSet)
	return &scrapers.RedditSubmissionsScraper{}, nil
}

func initCommentsScraper(ctx context.Context, redditCredentials reddit.Credentials, memoryOpts scrapers.MemoryOpts) (*scrapers.RedditCommentsScraper, error) {
	wire.Build(scrapers.SuperSet)
	return &scrapers.RedditCommentsScraper{}, nil
}
