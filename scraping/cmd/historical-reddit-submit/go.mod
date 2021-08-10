module historical-reddit-submit

go 1.16

replace github.com/VarityPlatform/scraping/common => ../../common

require (
	cloud.google.com/go v0.90.0
	github.com/VarityPlatform/scraping/common v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	google.golang.org/genproto v0.0.0-20210728212813-7823e685a01f
)
