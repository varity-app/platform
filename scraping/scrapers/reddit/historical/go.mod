module github.com/VarityPlatform/varity-scraping/scraping/scrapers/reddit/historical

go 1.16

replace github.com/VarityPlatform/scraping/scrapers => ../../

replace github.com/VarityPlatform/scraping/common => ../../../common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../../../protobuf/reddit

require (
	github.com/VarityPlatform/scraping/protobuf/reddit v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/scrapers v0.0.0-00010101000000-000000000000
	github.com/google/wire v0.5.0
	google.golang.org/protobuf v1.27.1
)
