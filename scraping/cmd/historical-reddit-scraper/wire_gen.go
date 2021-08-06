// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"context"
	"github.com/VarityPlatform/scraping/scrapers"
	"github.com/VarityPlatform/scraping/scrapers/reddit/historical"
)

// Injectors from wire.go:

func initSubmissionsScraper(ctx context.Context, memoryOpts scrapers.MemoryOpts) (*historical.SubmissionsScraper, error) {
	memory, err := scrapers.NewMemory(ctx, memoryOpts)
	if err != nil {
		return nil, err
	}
	submissionsScraper, err := historical.NewSubmissionsScraper(memoryOpts, memory)
	if err != nil {
		return nil, err
	}
	return submissionsScraper, nil
}

func initCommentsScraper(ctx context.Context, memoryOpts scrapers.MemoryOpts) (*historical.CommentsScraper, error) {
	memory, err := scrapers.NewMemory(ctx, memoryOpts)
	if err != nil {
		return nil, err
	}
	commentsScraper, err := historical.NewCommentsScraper(memoryOpts, memory)
	if err != nil {
		return nil, err
	}
	return commentsScraper, nil
}
