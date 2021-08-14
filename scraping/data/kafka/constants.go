package kafka

import "time"

// ConsumerTimeout is the duration for a consumer to wait to receive
// messages from kafka before cancelling
const ConsumerTimeout time.Duration = 1 * time.Second

// ProducerTimeoutMS is how long a producer will wait before timing out
const ProducerTimeoutMS int = 1500
