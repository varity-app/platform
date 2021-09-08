package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/varity-app/platform/scraping/internal/common"
	"github.com/varity-app/platform/scraping/internal/data/bigquery2influx"

	"github.com/go-redis/redis/v8"

	"cloud.google.com/go/bigquery"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Wrap a DBT Bigquery table name with the GCP project and dataset name.
func wrapDBTTable(tableName string) string {
	return fmt.Sprintf(
		"%s.%s_%s.%s",
		common.GCPProjectID,
		common.BigqueryDatasetDBTScraping,
		viper.GetString("deployment.mode"),
		tableName,
	)
}

// Entrypoint method
func main() {

	// Initialize viper configuration
	err := initConfig()
	if err != nil {
		log.Fatalf("viper.BindEnv: %v", err)
	}

	// Init bigquery client
	ctx := context.Background()
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer bqClient.Close()

	// Init InfluxDB client
	influxClient := influxdb2.NewClient(viper.GetString("influxdb.url"), viper.GetString("influxdb.token"))
	defer influxClient.Close()

	// Init repos
	redisOptions := &redis.Options{
		Addr:     viper.GetString("redis.bigquery.endpoint"),
		Password: viper.GetString("redis.bigquery.password"),
	}
	membershipRepo := bigquery2influx.NewCohortMembershipRepo(bqClient, redisOptions)
	mentionsRepo := bigquery2influx.NewMentionsRepo(bqClient, redisOptions)

	// Specify hour to scrape
	year := 2021
	month := 9
	day := 1
	hour := 18

	// Fetch cohort memberships from the past month
	tableName := wrapDBTTable(bigquery2influx.CohortTableRedditorsAge)
	memberships, err := membershipRepo.GetMemberships(ctx, tableName, year, month-1)
	if err != nil {
		log.Fatalf("membershipRepo.GetMemberships: %v", err)
	}

	// Fetch mentions for the given hour
	tableName = wrapDBTTable(bigquery2influx.MentionsTableReddit)
	mentions, err := mentionsRepo.GetMentions(ctx, tableName, year, month, day, hour)
	if err != nil {
		log.Fatalf("mentionsRepo.GetMentions: %v", err)
	}

	// Aggregate mentions
	time := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)
	points := aggregate(memberships, mentions, time)

	log.Println(len(points), points[0])

	// Write points to InfluxDB
	writeAPI := influxClient.WriteAPIBlocking(viper.GetString("influxdb.org"), viper.GetString("influxdb.bucket"))
	for _, point := range points {
		if err := writeAPI.WritePoint(ctx, point); err != nil {
			log.Fatalf("influxdb.WritePoint: %v", err)
		}
	}

	log.Println(len(points), points[0])

}
