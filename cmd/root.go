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
	// cmds
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(reposCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(genCmd)

	// flags
	genCmd.Flags().String("file", "", "Pass file path for document generation")
	genCmd.Flags().String("dir", "", "Pass a directory to read files from it")
	genCmd.Flags().String("domain", "github.com", "Use a different domain to clone your repo instead of default github.com")
	genCmd.Flags().Bool("ssh", false, "Use github ssh for cloning")
}
