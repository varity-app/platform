package config

import "github.com/spf13/viper"

// Initialize microservice urls configuration
func initURLs() error {
	err := viper.BindEnv("urls.b2i", "URL_BIGQUERY_TO_INFLUX")
	if err != nil {
		return err
	}

	return nil
}
