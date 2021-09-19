package config

import "github.com/spf13/viper"

// Initialize postgres configuration
func initPostgres() error {
	err := viper.BindEnv("postgres.username", "POSTGRES_USERNAME")
	if err != nil {
		return err
	}

	err = viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	if err != nil {
		return err
	}

	viper.SetDefault("postgres.address", "localhost:1234")
	err = viper.BindEnv("postgres.address", "POSTGRES_ADDRESS")
	if err != nil {
		return err
	}

	viper.SetDefault("postgres.database", "finance")
	err = viper.BindEnv("postgres.database", "POSTGRES_DB")
	if err != nil {
		return err
	}

	err = viper.BindEnv("postgres.network", "POSTGRES_NETWORK")
	if err != nil {
		return err
	}

	return nil
}
