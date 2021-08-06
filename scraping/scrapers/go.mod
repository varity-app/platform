//Package scrapers defines all webscrapers used in varity
module scrapers

go 1.16

replace github.com/VarityPlatform/scraping/common => ../common

replace github.com/VarityPlatform/scraping/protobuf/reddit => ../protobuf/reddit

require (
	cloud.google.com/go/firestore v1.5.0
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.49.0 // indirect
	google.golang.org/genproto v0.0.0-20210624174822-c5cf32407d0a // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
