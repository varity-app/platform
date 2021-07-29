package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Initialize viper
func initConfig() error {
	viper.AutomaticEnv()

	// Deployment mode should be `dev` or `prod`
	viper.SetDefault("deploymentMode", "dev")
	err := viper.BindEnv("deploymentMode", "DEPLOYMENT_MODE")
	if err != nil {
		return fmt.Errorf("viper.BindEnv: %v", err)
	}

	// Port to run webserver on
	viper.SetDefault("port", 8000)
	err = viper.BindEnv("port", "PORT")
	if err != nil {
		return fmt.Errorf("error binding env variable: %v", err)
	}

	// Tracing probability
	viper.SetDefault("tracing.probability", 0.01)
	err = viper.BindEnv("tracing.probability", "TRACING_PROBABILITY")
	if err != nil {
		return fmt.Errorf("error binding env variable: %v", err)
	}

	return nil
}
