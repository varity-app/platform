package main

import (
	"fmt"

	"github.com/VarityPlatform/scraping/common"

	"github.com/spf13/viper"
)

// Initialize viper
func initConfig() error {
	viper.AutomaticEnv()

	// Application deployment mode (dev or prod)
	viper.SetDefault("deploymentMode", common.DEPLOYMENT_MODE_DEV)
	err := viper.BindEnv("deploymentMode", "DEPLOYMENT_MODE")
	if err != nil {
		return fmt.Errorf("error binding env variable: %v", err)
	}

	// Port to run webserver on
	viper.SetDefault("port", 8000)
	err = viper.BindEnv("port", "PORT")
	if err != nil {
		return fmt.Errorf("error binding env variable: %v", err)
	}

	viper.SetDefault("tracing.probability", 0.01)
	err = viper.BindEnv("tracing.probability", "TRACING_PROBABILITY")
	if err != nil {
		return fmt.Errorf("error binding env variable: %v", err)
	}

	return nil
}
