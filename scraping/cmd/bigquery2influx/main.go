package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/varity-app/platform/scraping/internal/config"
	"github.com/varity-app/platform/scraping/internal/data"
	b2i "github.com/varity-app/platform/scraping/internal/services/bigquery2influx"
)

// Entrypoint method
func main() {

	// Initialize viper configuration
	err := config.InitConfig()
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
		Postgres: data.PostgresOpts{
			Username: viper.GetString("postgres.username"),
			Password: viper.GetString("postgres.password"),
			Address:  viper.GetString("postgres.address"),
			Database: viper.GetString("postgres.database"),
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
