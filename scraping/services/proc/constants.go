package main

import "time"

// FirestoreKafkaOffsets is the name of a firestore collection for checkpointing kafka offsets
const FirestoreKafkaOffsets string = "kafka-offsets"

// ConsumerTimeout is the duration for a consumer to wait to receive
// messages from kafka before cancelling
const ConsumerTimeout time.Duration = 2 * time.Second

// ProducerTimeoutMS is the duration in milliseconds before timing out
// when flushing the kafka producer
const ProducerTimeoutMS int = 1000 // 1 second

// SinkBatchSize is the number of rows of each batch to commit to bigquery
const SinkBatchSize int = 250
