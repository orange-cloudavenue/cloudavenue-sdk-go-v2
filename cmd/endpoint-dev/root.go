package main

import (
	"github.com/spf13/cobra"
)

var (
	loggerLevel string
	mockFlag    bool
	configFile  string

	rootCmd = &cobra.Command{
		Use:   "endpoint-dev",
		Short: "CloudAvenue Endpoint development CLI",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&loggerLevel, "logger", "info", "Set the logger level (e.g., debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&mockFlag, "mock", false, "Use the mock client (default: false)")
	rootCmd.PersistentFlags().StringVar(&configFile, "file", "", "Path to config file (default: $HOME/.sdkv2/devtools/config.yaml)")
}
