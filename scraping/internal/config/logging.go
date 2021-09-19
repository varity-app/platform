package config

import "github.com/spf13/viper"

// Initialize logging configuration
func initLogging() error {
	err := viper.BindEnv("logging.level", "LOG_LEVEL")
	if err != nil {
		return err
	}

	return nil
}
