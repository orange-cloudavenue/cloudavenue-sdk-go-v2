package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "Generator is a command to generate commands/endpoints code from the definitions.",
	Long:  `Generator is a command to generate commands/endpoints code from the definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", "", "The path to the file to generate commands from")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")

	if err := rootCmd.MarkPersistentFlagRequired("path"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
