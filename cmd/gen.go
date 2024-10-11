package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen <repo-name> <file>",
	Short: "Generate docs for codebase of a repo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]
		fileName := ""
		if len(args) > 1 {
			fileName = args[1]
		}
		genDocs(repoName, fileName)
	},
}

func genDocs(repo string, file string) {
	fmt.Println(repo)
	fmt.Println(file)
	// TODO: check if current directory if a git repo
	// TODO: check if current directory name matches repo name
	// TODO: if repo not available git clone https://github.com/user.login/reponame
	// TODO: check if file exists in local repo (relative path)
	// TODO: if (file) { get functions store in chunks }
	// TODO: if (repo) { fetch project info; fetch main source }
	// TODO: generateDocs(content) -> storeTo("README.md")
}
