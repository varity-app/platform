// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/VarityPlatform/scraping/data/kafka"

	"github.com/google/wire"
)

func initProcessor(ctx context.Context, opts kafka.ProcessorOpts, offsetOpts kafka.OffsetManagerOpts) (*kafka.Processor, error) {
	wire.Build(kafka.SuperSet)
	return &kafka.Processor{}, nil
}

func initSink(ctx context.Context, opts kafka.BigquerySinkOpts, offsetOpts kafka.OffsetManagerOpts) (*kafka.BigquerySink, error) {
	wire.Build(kafka.SuperSet)
	return &kafka.BigquerySink{}, nil
}
