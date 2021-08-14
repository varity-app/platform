package kafka

import "github.com/google/wire"

// SuperSet is the wire superset for the kafka module
var SuperSet = wire.NewSet(NewPublisher, NewProcessor, NewOffsetManager, NewBigquerySink)
