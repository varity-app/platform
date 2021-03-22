# Create Kafka Topics on Confluent Cloud

```
ccloud kafka topic create reddit-submissions --partitions 6 --config retention.ms=604800000
ccloud kafka topic create reddit-comments --partitions 6 --config retention.ms=604800000
ccloud kafka topic create ticker-mentions --partitions 6 --config retention.ms=604800000
ccloud kafka topic create scraped-posts --partitions 6 --config retention.ms=604800000
ccloud kafka topic create post-sentiment --partitions 6 --config retention.ms=604800000
```