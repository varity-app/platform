module github.com/VarityPlatform/varity-scraping/scraping/data/prices

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

require (
	cloud.google.com/go/bigquery v1.20.1
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
)
