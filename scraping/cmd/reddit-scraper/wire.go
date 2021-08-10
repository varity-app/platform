// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/VarityPlatform/scraping/scrapers/reddit/live"

	"github.com/VarityPlatform/scraping/data/kafka"

	"github.com/google/wire"
)

func initSubmissionsScraper(ctx context.Context, redditCredentials reddit.Credentials, memoryOpts scrapers.MemoryOpts) (*live.RedditSubmissionsScraper, error) {
	wire.Build(live.SuperSet)
	return &live.RedditSubmissionsScraper{}, nil
}

func initCommentsScraper(ctx context.Context, redditCredentials reddit.Credentials, memoryOpts scrapers.MemoryOpts) (*live.RedditCommentsScraper, error) {
	wire.Build(live.SuperSet)
	return &live.RedditCommentsScraper{}, nil
}

func initPublisher(ctx context.Context, kafkaOpts kafka.Opts) (*kafka.Publisher, error) {
	wire.Build(kafka.SuperSet)
	return &kafka.Publisher{}, nil
}
