module reddit-scraper

go 1.16

replace github.com/VarityPlatform/scraping/scrapers => ../../scrapers

replace github.com/VarityPlatform/scraping/common => ../../common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../../protobuf/reddit

require (
	cloud.google.com/go/firestore v1.5.0
	cloud.google.com/go/pubsub v1.3.1
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.21.0
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/protobuf/reddit v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/scrapers v0.0.0-00010101000000-000000000000
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/google/go-querystring v1.0.0
	github.com/google/subcommands v1.2.0 // indirect
	github.com/google/wire v0.5.0
	github.com/labstack/echo/v4 v4.4.0
	github.com/spf13/viper v1.8.1
	github.com/vartanbeno/go-reddit/v2 v2.0.1
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	google.golang.org/api v0.51.0 // indirect
	google.golang.org/genproto v0.0.0-20210726200206-e7812ac95cc0 // indirect
	google.golang.org/protobuf v1.27.1
)
