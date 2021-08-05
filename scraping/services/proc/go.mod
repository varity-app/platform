module proc

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

replace github.com/VarityPlatform/scraping/transforms => ../../transforms

replace github.com/VarityPlatform/scraping/protobuf/common => ../../protobuf/common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../../protobuf/reddit

require (
	cloud.google.com/go/bigquery v1.19.0
	cloud.google.com/go/firestore v1.1.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.21.0
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/protobuf/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/protobuf/reddit v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/transforms v0.0.0-00010101000000-000000000000
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/labstack/echo/v4 v4.4.0
	github.com/spf13/viper v1.8.1
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	google.golang.org/api v0.51.0 // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)
