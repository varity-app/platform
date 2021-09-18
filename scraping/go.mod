module github.com/varity-app/platform/scraping

go 1.16

require (
	cloud.google.com/go/bigquery v1.22.0
	cloud.google.com/go/cloudtasks v0.1.0
	cloud.google.com/go/firestore v1.5.0
	cloud.google.com/go/kms v0.1.0 // indirect
	cloud.google.com/go/pubsub v1.3.1
	github.com/go-pg/pg/v10 v10.10.5
	github.com/google/go-querystring v1.1.0
	github.com/google/wire v0.5.0
	github.com/influxdata/influxdb-client-go/v2 v2.5.0
	github.com/labstack/echo/v4 v4.5.0
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/segmentio/kafka-go v0.4.18
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/vartanbeno/go-reddit/v2 v2.0.1
	google.golang.org/api v0.56.0
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.1
	gorm.io/gorm v1.21.15
	gotest.tools v2.2.0+incompatible
)
