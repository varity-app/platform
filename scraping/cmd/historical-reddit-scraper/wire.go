// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/VarityPlatform/scraping/scrapers/reddit/historical"

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
