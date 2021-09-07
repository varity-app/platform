package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"

	"github.com/VarityPlatform/scraping/common"
	"github.com/VarityPlatform/scraping/data/bigquery2influx"
	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
)

// Entrypoint method
func main() {

	err := initConfig()
	if err != nil {
		log.Fatalf("viper.BindEnv: %v", err)
	}

	ctx := context.Background()
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}

	redisOptions := &redis.Options{
		Addr:     viper.GetString("redis.bigquery.endpoint"),
		DB:       viper.GetInt("redis.bigquery.database"),
		Password: viper.GetString("redis.bigquery.password"),
	}

	// membershipRepo := bigquery2influx.NewCohortMembershipRepo(bqClient, redisOptions)

	// tableName := fmt.Sprintf("%s.%s_%s.%s", common.GCPProjectID, common.BigqueryDatasetDBTScraping, viper.GetString("deployment.mode"), bigquery2influx.CohortTableRedditorsAge)
	// memberships, err := membershipRepo.GetMemberships(ctx, tableName, 2021, 8)
	// if err != nil {
	// 	log.Fatalf("membershipRepo.GetMemberships: %v", err)
	// }

	// log.Println(len(memberships), memberships[0])

	mentionsRepo := bigquery2influx.NewMentionsRepo(bqClient, redisOptions)

	tableName := fmt.Sprintf("%s.%s_%s.%s", common.GCPProjectID, common.BigqueryDatasetDBTScraping, viper.GetString("deployment.mode"), bigquery2influx.MentionsTableReddit)
	mentions, err := mentionsRepo.GetMentions(ctx, tableName, 2021, 9, 1, 18)
	if err != nil {
		log.Fatalf("mentionsRepo.GetMentions: %v", err)
	}

	log.Println(len(mentions), mentions[0])

}
