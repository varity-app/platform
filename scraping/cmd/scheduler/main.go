package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
	scheduler "github.com/varity-app/platform/scraping/internal/services/scheduler"
)

// Entrypoint method
func main() {

	// Initialize viper configuration
	err := initConfig()
	if err != nil {
		log.Fatalf("viper.BindEnv: %v", err)
	}

	// Create background context
	ctx := context.Background()

	// Initialize service
	web, err := scheduler.NewService(ctx, scheduler.ServiceOpts{
		DeploymentMode: viper.GetString("deployment.mode"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
