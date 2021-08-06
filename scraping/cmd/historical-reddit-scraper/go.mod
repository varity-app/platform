module historical-reddit-scraper

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

replace github.com/VarityPlatform/scraping/scrapers => ../../scrapers

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../../protobuf/reddit

replace github.com/VarityPlatform/scraping/scrapers/reddit/historical => ../../scrapers/reddit/historical

require (
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/scrapers v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/scrapers/reddit/historical v0.0.0-00010101000000-000000000000
	github.com/google/wire v0.5.0
	github.com/labstack/echo/v4 v4.5.0
	github.com/spf13/viper v1.8.1
)
