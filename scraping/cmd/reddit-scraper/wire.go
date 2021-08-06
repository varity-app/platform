// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/VarityPlatform/scraping/scrapers"
	redditLive "github.com/VarityPlatform/scraping/scrapers/reddit/live"

	"github.com/google/wire"
)

func initSubmissionsScraper(ctx context.Context, redditCredentials reddit.Credentials, memoryOpts scrapers.MemoryOpts) (*redditLive.RedditSubmissionsScraper, error) {
	wire.Build(redditLive.SuperSet)
	return &redditLive.RedditSubmissionsScraper{}, nil
}

func initCommentsScraper(ctx context.Context, redditCredentials reddit.Credentials, memoryOpts scrapers.MemoryOpts) (*redditLive.RedditCommentsScraper, error) {
	wire.Build(redditLive.SuperSet)
	return &redditLive.RedditCommentsScraper{}, nil
}
