package config

import "github.com/spf13/viper"

// Initialize reddit configuration
func initReddit() error {
	err := viper.BindEnv("reddit.username", "REDDIT_USERNAME")
	if err != nil {
		return err
	}

	err = viper.BindEnv("reddit.password", "REDDIT_PASSWORD")
	if err != nil {
		return err
	}

	err = viper.BindEnv("reddit.client_id", "REDDIT_CLIENT_ID")
	if err != nil {
		return err
	}

	err = viper.BindEnv("reddit.client_secret", "REDDIT_CLIENT_SECRET")
	if err != nil {
		return err
	}

	return nil
}
