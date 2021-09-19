package main

import (
	"github.com/spf13/cobra"
)

// Root cobra command
var rootCmd = &cobra.Command{}

// Initialize cobra commands
func init() {
	initRedditHistorical(rootCmd)
}

// Entrypoint method
func main() {
	cobra.CheckErr(rootCmd.Execute())
}
