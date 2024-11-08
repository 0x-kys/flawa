package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flawa",
	Short: "A CLI tool to generate docs from code",
	Long:  `A command-line tool to generate docs from your code files`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listFilesCmd)
	rootCmd.AddCommand(generateDocumentCmd)
}
