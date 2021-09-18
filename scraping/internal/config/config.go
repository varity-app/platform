package config

import (
	"github.com/spf13/viper"
	"github.com/varity-app/platform/scraping/internal/common"
)

// Initialize all viper configuration settings
func InitConfig() error {
	viper.AutomaticEnv()

	err := initCommon()
	if err != nil {
		return err
	}

	err = initInflux()
	if err != nil {
		return err
	}

	err = initPostgres()
	if err != nil {
		return err
	}

	return nil
}

// Initialize common configuration settings
func initCommon() error {
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
