package config

import "github.com/spf13/viper"

// Initialize influxdb configuration
func initInflux() error {
	err := viper.BindEnv("influxdb.url", "INFLUX_URL")
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
