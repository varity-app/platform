module reddit-scraper

go 1.13

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../protobuf/reddit

require (
	cloud.google.com/go/firestore v1.5.0
	github.com/VarityPlatform/scraping/protobuf/reddit v0.0.0-00010101000000-000000000000
	github.com/google/go-querystring v1.0.0
	github.com/spf13/viper v1.8.1
	github.com/vartanbeno/go-reddit v1.0.0
	github.com/vartanbeno/go-reddit/v2 v2.0.1
	google.golang.org/api v0.51.0 // indirect
	google.golang.org/genproto v0.0.0-20210726200206-e7812ac95cc0 // indirect
	google.golang.org/protobuf v1.27.1
)
