package kafka

import "github.com/google/wire"

// Opts is a struct containing parameters for a new kafka producer
type Opts struct {
	BootstrapServers string
	Username         string
	Password         string
}

// SuperSet is the wire superset for the kafka module
var SuperSet = wire.NewSet(NewPublisher, NewProcessor, NewOffsetManager, NewBigquerySink)
