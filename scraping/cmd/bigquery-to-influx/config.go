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

	// Redis variables
	err = viper.BindEnv("redis.bigquery.endpoint", "REDIS_BIGQUERY_ENDPOINT")
	if err != nil {
		return err
	}
	err = viper.BindEnv("redis.bigquery.password", "REDIS_BIGQUERY_PASSWORD")
	if err != nil {
		return err
	}

	// InfluxDB variables
	err = viper.BindEnv("influxdb.url", "INFLUX_URL")
	if err != nil {
		return err
	}
	err = viper.BindEnv("influxdb.token", "INFLUX_TOKEN")
	if err != nil {
		return err
	}
	viper.SetDefault("influxdb.org", "Varity")
	err = viper.BindEnv("influxdb.org", "INFLUX_ORG")
	if err != nil {
		return err
	}
	viper.SetDefault("influxdb.bucket", "dev")
	err = viper.BindEnv("influxdb.bucket", "INFLUX_BUCKET")
	if err != nil {
		return err
	}

	return nil
}
