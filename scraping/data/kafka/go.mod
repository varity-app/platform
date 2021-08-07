module github.com/VarityPlatform/varity-scraping/scraping/data/kafka

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

replace github.com/VarityPlatform/scraping/protobuf/common => ../../protobuf/common

require (
	cloud.google.com/go/bigquery v1.8.0
	cloud.google.com/go/firestore v1.5.0
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/google/wire v0.5.0
	google.golang.org/api v0.44.0 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/grpc v1.39.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
