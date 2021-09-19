package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/varity-app/platform/scraping/internal/config"
	"github.com/varity-app/platform/scraping/internal/logging"
	scheduler "github.com/varity-app/platform/scraping/internal/services/scheduler"
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

	// Create background context
	ctx := context.Background()

	// Initialize service
	service, err := scheduler.NewService(ctx, logger, scheduler.ServiceOpts{
		DeploymentMode: viper.GetString("deployment.mode"),
		B2IUrl:         viper.GetString("urls.b2i"),
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer service.Close()

	// Start webserver
	port := viper.GetString("port")
	service.Start(":" + port)
}
