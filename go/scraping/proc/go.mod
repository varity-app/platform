module proc

go 1.16

replace github.com/VarityPlatform/scraping/common => ../common

replace github.com/VarityPlatform/scraping/protobuf/common => ../protobuf/common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../protobuf/reddit

require (
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/protobuf/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/protobuf/reddit v0.0.0-00010101000000-000000000000
	github.com/go-pg/pg/v10 v10.10.3
	google.golang.org/protobuf v1.25.0
)
