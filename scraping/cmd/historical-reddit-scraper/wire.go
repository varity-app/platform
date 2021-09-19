// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/varity-app/platform/scraping/internal/data/kafka"
	"github.com/varity-app/platform/scraping/internal/scrapers"
	"github.com/varity-app/platform/scraping/internal/scrapers/reddit/historical"

	"github.com/google/wire"
)

func initSubmissionsScraper(ctx context.Context, memoryOpts scrapers.MemoryOpts) (*historical.SubmissionsScraper, error) {
	wire.Build(historical.SuperSet)
	return &historical.SubmissionsScraper{}, nil
}

func initCommentsScraper(ctx context.Context, memoryOpts scrapers.MemoryOpts) (*historical.CommentsScraper, error) {
	wire.Build(historical.SuperSet)
	return &historical.CommentsScraper{}, nil
}

func initPublisher(ctx context.Context, opts kafka.PublisherOpts) (*kafka.Publisher, error) {
	wire.Build(kafka.SuperSet)
	return &kafka.Publisher{}, nil
}