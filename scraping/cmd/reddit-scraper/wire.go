// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/varity-app/platform/scraping/internal/scrapers/reddit/live"

	"github.com/varity-app/platform/scraping/internal/data/kafka"

	"github.com/go-redis/redis/v8"

	"github.com/google/wire"
)

func initSubmissionsScraper(ctx context.Context, redditCredentials reddit.Credentials, rdb *redis.Client) (*live.RedditSubmissionsScraper, error) {
	wire.Build(live.SuperSet)
	return &live.RedditSubmissionsScraper{}, nil
}

func initCommentsScraper(ctx context.Context, redditCredentials reddit.Credentials, rdb *redis.Client) (*live.RedditCommentsScraper, error) {
	wire.Build(live.SuperSet)
	return &live.RedditCommentsScraper{}, nil
}

func initPublisher(ctx context.Context, opts kafka.PublisherOpts) (*kafka.Publisher, error) {
	wire.Build(kafka.SuperSet)
	return &kafka.Publisher{}, nil
}
