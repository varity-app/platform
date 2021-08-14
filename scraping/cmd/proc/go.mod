module proc

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

replace github.com/VarityPlatform/scraping/transforms => ../../transforms

replace github.com/VarityPlatform/scraping/data/kafka => ../../data/kafka

replace github.com/VarityPlatform/scraping/protobuf/common => ../../protobuf/common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../../protobuf/reddit

require (
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/data/kafka v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/protobuf/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/protobuf/reddit v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/transforms v0.0.0-00010101000000-000000000000
	github.com/google/wire v0.5.0
	github.com/labstack/echo/v4 v4.4.0
	github.com/segmentio/kafka-go v0.4.17
	github.com/spf13/viper v1.8.1
	google.golang.org/api v0.51.0 // indirect
	google.golang.org/protobuf v1.27.1
)
