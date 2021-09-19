// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"context"
	"github.com/varity-app/platform/scraping/internal/data/kafka"
)

// Injectors from wire.go:

func initProcessor(ctx context.Context, opts kafka.ProcessorOpts, offsetOpts kafka.OffsetManagerOpts) (*kafka.Processor, error) {
	offsetManager, err := kafka.NewOffsetManager(ctx, offsetOpts)
	if err != nil {
		return nil, err
	}
	processor := kafka.NewProcessor(ctx, offsetManager, opts)
	return processor, nil
}

func initSink(ctx context.Context, opts kafka.BigquerySinkOpts, offsetOpts kafka.OffsetManagerOpts) (*kafka.BigquerySink, error) {
	offsetManager, err := kafka.NewOffsetManager(ctx, offsetOpts)
	if err != nil {
		return nil, err
	}
	bigquerySink, err := kafka.NewBigquerySink(ctx, offsetManager, opts)
	if err != nil {
		return nil, err
	}
	return bigquerySink, nil
}