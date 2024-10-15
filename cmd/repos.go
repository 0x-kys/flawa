package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v66/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List all your GitHub repositories (public and private)",
	Run: func(cmd *cobra.Command, args []string) {
		fetchRepositories()
	},
}

func fetchRepositories() {
	token, err := os.ReadFile(GetConfigPath(".token"))
	if err != nil {
		log.Fatalf("Error reading token: %v", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(token)})
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	repos, _, err := client.Repositories.ListByAuthenticatedUser(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error fetching repositories: %v", err)
	}

	fmt.Printf("%-40s | %-12s | %-12s | %-12s\n", "Repo Name", "Visibility", "Created At", "Updated At")
	fmt.Println(strings.Repeat("-", 100))

	for _, repo := range repos {
		repoName := *repo.Name
		repoCreatedAt := FormatDate(*repo.CreatedAt)
		repoUpdatedAt := FormatDate(*repo.UpdatedAt)
		repoVisibility := repo.GetVisibility()

		fmt.Printf("%-40s | %-12s | %-12s | %-12s\n", repoName, repoVisibility, repoCreatedAt, repoUpdatedAt)
	}
}

func FormatDate(t github.Timestamp) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
}
