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

		fmt.Println("File found, reading chunks.")

		chunks, err := readFileChunks(file)
		if err != nil {
			fmt.Printf("Error reading file chunks: %v\n", err)
			return
		}

		combinedContent := combinedChunks(chunks)

		betterInstructions := "Your task is to generate write report of the code i share in markdown. Code that i'll share could be of any programming language. I want you to divide it into chunks and mention the main section of it (identifier most probably) and write a line mentioning what it does. For example: for `#include <stdio.h>` you need to mention a heading 'imports' and below that you can make a list of all the imports after that for functions you can mention function as heading and then function name, parameter with it's type and then in next line mention what that function does. Do not send anything other than the main content that is requested by the user. Only read the code and make a proper document style doc explaining code in the way i mentioned."

		response, err := sentToStudio(combinedContent, betterInstructions)
		if err != nil {
			fmt.Printf("Error generating documentation: %v\n", err)
			return
		}

		saveReadme(response)
	}

	// WIP: if (file) { get functions store in chunks }
	// TODO: if (repo) { fetch project info; fetch main source }
	// WIP: generateDocs(content) -> storeTo("README.md")
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
		/**
		* TODO: if this part is being executed
		* - ask user for default (git@github.com)
		* - ask user for custom (git@host)
		 */
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
