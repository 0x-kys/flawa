package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/go-github/v66/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var genCmd = &cobra.Command{
	Use:   "gen <repo-name> [file] <ssh>",
	Short: "Generate docs for codebase of a repo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]
		fileName := ""
		useSSH := true
		if len(args) > 1 {
			fileName = args[1]
		}
		genDocs(repoName, fileName, useSSH)
	},
}

func genDocs(repo string, file string, useSSH bool) {
	fmt.Printf("Repository Name: %s\n", repo)
	fmt.Printf("File: %s\n", file)

	if !isGitRepo() && !dirExists(repo) {
		fmt.Println("Current directory is not a Git repository.")

		fmt.Printf("Cloning repository '%s'...\n", repo)
		err := cloneRepo(repo, useSSH)
		if err != nil {
			fmt.Printf("Error cloning repository: %v\n", err)
			return
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	dirName := filepath.Base(dir)
	if dirName != repo {
		fmt.Printf("Directory name '%s' does not match repository name '%s'.\n", dirName, repo)
		return
	}

	if len(file) > 0 {
		if !fileExists(file) {
			fmt.Printf("File '%s' does not exist in the repository.\n", file)
			return
		}

		fmt.Println("File found, proceeding with documentation generation.")
	}

	// TODO: if (file) { get functions store in chunks }
	// TODO: if (repo) { fetch project info; fetch main source }
	// TODO: generateDocs(content) -> storeTo("README.md")
}

func fetchUsername() string {
	token, err := os.ReadFile("token.txt")
	if err != nil {
		logrus.Fatalf("Error reading token: %v", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(token)})
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		logrus.Fatalf("Error fetching authenticated user info: %v", err)
	}

	return *user.Login
}

func isGitRepo() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

func cloneRepo(repo string, useSSH bool) error {
	user := fetchUsername()

	var cloneURL string

	if !useSSH {
		cloneURL = fmt.Sprintf("https://github.com/%s/%s.git", user, repo)
	} else {
		cloneURL = fmt.Sprintf("git@kys:%s/%s.git .", user, repo)
	}

	cmd := exec.Command("git", "clone", cloneURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to clone the repository: %v", err)
	}
	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
