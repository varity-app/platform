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
		return fmt.Errorf("viper.BindEnv: %v", err)
	}

	viper.SetDefault("proxy.url", "http://localhost:2933")
	err = viper.BindEnv("proxy.url", "PROXY_URL")
	if err != nil {
		return fmt.Errorf("viper.BindEnv: %v", err)
	}

	return nil
}
