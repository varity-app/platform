# Create Kafka Topics on Confluent Cloud

```
ccloud kafka topic create reddit-submissions --partitions 2 --config retention.ms=604800000
ccloud kafka topic create reddit-comments --partitions 2 --config retention.ms=604800000
ccloud kafka topic create ticker-mentions --partitions 2 --config retention.ms=604800000
ccloud kafka topic create scraped-posts --partitions 2 --config retention.ms=604800000
ccloud kafka topic create post-sentiment --partitions 2 --config retention.ms=604800000
ccloud kafka topic create varity-faust-app-__assignor-__leader --partitions 2
ccloud kafka topic create logs --partitions 2
```