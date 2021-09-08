module bigquery-to-influx

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

replace github.com/VarityPlatform/scraping/transforms => ../../transforms

replace github.com/VarityPlatform/scraping/protobuf/common => ../../protobuf/common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../../protobuf/reddit

replace github.com/VarityPlatform/scraping/data/bigquery2influx => ../../data/bigquery2influx

require (
	cloud.google.com/go/bigquery v1.22.0
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/data/bigquery2influx v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/transforms v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.3
	github.com/influxdata/influxdb-client-go/v2 v2.5.0
	github.com/spf13/viper v1.8.1
)
