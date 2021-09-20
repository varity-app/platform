package config

import "github.com/spf13/viper"

// Initialize kafka configuration
func initKafka() error {
	err := viper.BindEnv("kafka.boostrap_servers", "KAFKA_BOOTSTRAP_SERVERS")
	if err != nil {
		return err
	}

	err = viper.BindEnv("kafka.auth_key", "KAFKA_AUTH_KEY")
	if err != nil {
		return err
	}

	err = viper.BindEnv("kafka.auth_secret", "KAFKA_AUTH_SECRET")
	if err != nil {
		return err
	}

	return nil
}
