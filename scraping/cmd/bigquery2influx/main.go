package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
	b2i "github.com/varity-app/platform/scraping/internal/services/bigquery2influx"

	"github.com/go-redis/redis/v8"
)

// Entrypoint method
func main() {

	// Initialize viper configuration
	err := initConfig()
	if err != nil {
		log.Fatalf("viper.BindEnv: %v", err)
	}

	// Init bigquery client
	ctx := context.Background()

	// Init service
	opts := b2i.ServiceOpts{
		Influx: b2i.InfluxOpts{
			Addr:   viper.GetString("influxdb.url"),
			Token:  viper.GetString("influxdb.token"),
			Org:    viper.GetString("influxdb.org"),
			Bucket: viper.GetString("influxdb.bucket"),
		},
		Redis: &redis.Options{
			Addr:     viper.GetString("redis.bigquery.endpoint"),
			Password: viper.GetString("redis.bigquery.password"),
		},
		DeploymentMode: viper.GetString("deployment.mode"),
	}

	web, err := b2i.NewService(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))

}
