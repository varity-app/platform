package config

import "github.com/spf13/viper"

// Initialize redis configuration
func initRedis() error {
	err := viper.BindEnv("redis.scraping.password", "REDIS_SCRAPING_PASSWORD")
	if err != nil {
		return err
	}

	err = viper.BindEnv("redis.scraping.address", "REDIS_SCRAPING_ADDRESS")
	if err != nil {
		return err
	}

	return nil
}
