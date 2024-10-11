package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v66/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var infoCmd = &cobra.Command{
	Use:   "info <repo-name>",
	Short: "Get info of a specific repo by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]
		fetchRepoInfo(repoName)
	},
}

func fetchRepoInfo(repoName string) {
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
	owner := *user.Login

	repo, _, err := client.Repositories.Get(context.Background(), owner, repoName)
	if err != nil {
		logrus.Fatalf("Error fetching repository info: %v", err)
	}

	commits, _, err := client.Repositories.ListCommits(context.Background(), owner, repoName, nil)
	if err != nil {
		logrus.Fatalf("Error fetching commits: %v", err)
	}
	totalCommits := len(commits)

	var stats []*github.ContributorStats
	retries := 5
	for retries > 0 {
		stats, _, err = client.Repositories.ListContributorsStats(context.Background(), owner, repoName)
		if err != nil && strings.Contains(err.Error(), "job scheduled on GitHub side") {
			logrus.Warnf("GitHub is processing stats, retrying in 10 seconds... (%d retries left)", retries)
			time.Sleep(10 * time.Second)
			retries--
		} else {
			break
		}
	}
	if err != nil {
		logrus.Fatalf("Error fetching repository stats: %v", err)
	}

	additions := 0
	deletions := 0
	for _, stat := range stats {
		for _, week := range stat.Weeks {
			additions += *week.Additions
			deletions += *week.Deletions
		}
	}

	fmt.Printf("Repo Name: %s\n", repoName)
	fmt.Printf("Total Commits: %d\n", totalCommits)
	fmt.Printf("Created At: %s\n", repo.GetCreatedAt().Format("2006-01-02"))
	fmt.Printf("Updated At: %s\n", repo.GetUpdatedAt().Format("2006-01-02"))
	fmt.Printf("Additions: %d\n", additions)
	fmt.Printf("Deletions: %d\n", deletions)
	fmt.Printf("Visibility: %s\n", repo.GetVisibility())
}
