package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/varity-app/platform/scraping/internal/config"
	"github.com/varity-app/platform/scraping/internal/logging"
	b2i "github.com/varity-app/platform/scraping/internal/services/bigquery2influx"
)

// Entrypoint method
func main() {

	// Initialize viper configuration
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("viper.BindEnv: %v", err)
	}

	// Create logger
	logger := logging.NewLogger(viper.GetString("logging.level"))

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
		DeploymentMode: viper.GetString("deployment.mode"),
	}

	service, err := b2i.NewService(ctx, logger, opts)
	if err != nil {
		logger.Fatal(err)
	}
	defer service.Close()

	// Start webserver
	port := viper.GetString("port")
	service.Start(":" + port)

}
