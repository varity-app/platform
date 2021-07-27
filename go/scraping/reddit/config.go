package main

import (
	"log"

	"github.com/spf13/viper"
)

// Initialize viper
func initConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("deploymentMode", "dev")
	err := viper.BindEnv("deploymentMode", "DEPLOYMENT_MODE")
	if err != nil {
		log.Fatal("Error binding env variable:", err.Error())
	}
}
