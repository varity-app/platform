module tiingo

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

replace github.com/VarityPlatform/scraping/data/prices => ../../data/prices

replace github.com/VarityPlatform/scraping/data/tickers => ../../data/tickers

replace github.com/VarityPlatform/scraping/scrapers/tiingo => ../../scrapers/tiingo

require (
	cloud.google.com/go/bigquery v1.20.1
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/data/prices v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/data/tickers v0.0.0-00010101000000-000000000000
	github.com/VarityPlatform/scraping/scrapers/tiingo v0.0.0-00010101000000-000000000000
	github.com/go-pg/pg/v10 v10.10.3
	github.com/labstack/echo/v4 v4.5.0
	github.com/spf13/viper v1.8.1
)
