package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v66/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var genCmd = &cobra.Command{
	Use:   "gen <repo-name>",
	Short: "Generate docs",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domainChanged := cmd.Flags().Changed("domain")

		repoName := args[0]
		fileName, _ := cmd.Flags().GetString("file")
		dirName, _ := cmd.Flags().GetString("dir")
		useSSH, _ := cmd.Flags().GetBool("ssh")
		domain, _ := cmd.Flags().GetString("domain")

		if len(repoName) < 1 {
			logrus.Fatalln("No repo name provided")
		}

		if len(fileName) > 0 && len(dirName) > 0 {
			logrus.Fatalln("Mention either a file (--file) or a directory (--dir)")
		}

		if !domainChanged && useSSH {
			domain = "git@github.com"
		}

		genDocs(repoName, fileName, dirName, useSSH, domain)
	},
}

func genDocs(repo string, file string, directory string, useSSH bool, domain string) {
	if !isGitRepo() && !dirExists(repo) {
		fmt.Println("Current directory is not a Git repository.")

		fmt.Printf("Cloning repository '%s'...\n", repo)
		err := cloneRepo(repo, useSSH, domain)
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

	if len(directory) == 0 && len(file) > 0 {
		if !fileExists(file) {
			logrus.Fatalf("File '%s' does not exist in the repository.\n", file)
			return
		}

		fmt.Println("File found, reading chunks.")

		chunks, err := readFileChunks(file)
		if err != nil {
			logrus.WithError(err).Fatalln("Error reading file chunks")
		}

		combinedContent := combinedChunks(chunks)

		betterInstructions := "Your task is to generate a well-structured Markdown report for the code I provide, which may be written in any programming language. Analyze the code and identify key sections, such as imports, functions, classes, or constants. For imports, create a heading titled 'Imports' and list each import or library along with a brief note on its purpose if identifiable. For functions, use the function name as a subheading, list the parameters with their types, and provide a brief description of what the function does. For other sections like classes or constants, follow a similar pattern by using appropriate headings or subheadings and adding concise descriptions of their roles. Make sure the report is clean, readable, and properly formatted in Markdown. Only provide the requested content and avoid adding any additional commentary or text beyond the specified structure. Keep the documentation minimal."

		response, err := sentToStudio(combinedContent, betterInstructions)
		if err != nil {
			logrus.WithError(err).Fatalln("Error generating documentation")
			return
		}

		saveReadme(response)
	} else if len(file) == 0 && len(directory) > 0 {
		logrus.Warnln("THIS FEATURE IS WORK IN PROGRESS!!")
	}
}

func sentToStudio(content string, instructions string) (string, error) {
	apiURL := "http://localhost:1234/v1/chat/completions"

	requestBody := map[string]interface{}{
		"model": "llama-3.2-1b-instruct",
		"messages": []map[string]string{
			{"role": "system", "content": instructions},
			{"role": "user", "content": content},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(body))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	output := response["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return output, nil
}

func combinedChunks(chunks []string) string {
	return strings.Join(chunks, "\n\n")
}

func fetchUsername() string {
	token, err := os.ReadFile(GetConfigPath(".token"))
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

func cloneRepo(repo string, useSSH bool, domain string) error {
	user := fetchUsername()

	var cloneURL string

	if useSSH {
		cloneURL = fmt.Sprintf("%s/%s/%s.git", domain, user, repo)
	} else {
		cloneURL = fmt.Sprintf("https://%s/%s/%s.git", domain, user, repo)
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

func readFileChunks(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var chunks []string
	var chunk strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		chunk.WriteString(line + "\n")

		if strings.HasSuffix(line, "}") || len(strings.Split(chunk.String(), "\n")) >= 50 {
			chunks = append(chunks, chunk.String())
			chunk.Reset()
		}
	}

	if chunk.Len() > 0 {
		chunks = append(chunks, chunk.String())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return chunks, nil
}

func saveReadme(content string) {
	err := os.WriteFile("doc.md", []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error saving README.md: %v\n", err)
		return
	}
	fmt.Println("doc generated successfully.")
}
