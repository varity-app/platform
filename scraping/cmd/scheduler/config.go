package main

import (
	"github.com/varity-app/platform/scraping/internal/common"

	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AutomaticEnv()

	// Deployment mode should be `dev` or `prod`
	viper.SetDefault("deployment.mode", common.DeploymentModeDev)
	err := viper.BindEnv("deployment.mode", "DEPLOYMENT_MODE")
	if err != nil {
		return err
	}

	// Port to run webserver on
	viper.SetDefault("port", 8000)
	err = viper.BindEnv("port", "PORT")
	if err != nil {
		return err
	}

	return nil
}
