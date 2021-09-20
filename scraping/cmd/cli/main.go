package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root cobra command
var rootCmd = &cobra.Command{}

// Initialize cobra commands
func init() {
	initRedditHistorical(rootCmd)
	initB2I(rootCmd)

	rootCmd.PersistentFlags().StringP("deployment", "d", "dev", "varity deployment mode")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
}

// Entrypoint method
func main() {
	cobra.CheckErr(rootCmd.Execute())
}
