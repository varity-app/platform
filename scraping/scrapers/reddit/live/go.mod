module github.com/VarityPlatform/varity-scraping/scraping/scrapers/reddit/live

go 1.16

replace github.com/VarityPlatform/scraping/scrapers => ../../

replace github.com/VarityPlatform/scraping/common => ../../../common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../../../protobuf/reddit

require (
	github.com/VarityPlatform/scraping/protobuf/reddit v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/scrapers v0.0.0-00010101000000-000000000000
	github.com/google/go-querystring v1.0.0
	github.com/google/wire v0.5.0
	github.com/vartanbeno/go-reddit/v2 v2.0.1
	google.golang.org/protobuf v1.27.1
)
