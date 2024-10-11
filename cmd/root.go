package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flawa",
	Short: "A CLI tool to interact with GitHub",
	Long:  `A command-line tool to log into GitHub and fetch repositories.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(reposCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(genCmd)
}
