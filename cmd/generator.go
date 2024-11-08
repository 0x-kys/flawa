package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"flawa/cfg"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var generateDocumentCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate document of specific file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		generateDocument(args[0])
	},
}

func generateDocument(filePath string) {
	if err := cfg.LoadConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Fatal().Msgf("Couldn't find this file: %s", filePath)
	}

	if fileInfo.IsDir() {
		log.Fatal().Msgf("Selected file is a directory: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read file content")
	}

	data := map[string]interface{}{
		"model":      cfg.Config.Ollama.Model,
		"prompt":     cfg.Config.Ollama.BasePrompt + string(content),
		"stream":     cfg.Config.Ollama.Stream,
		"keep_alive": 0,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal JSON data")
	}

	start := time.Now()

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send request to Ollama API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal().Msgf("Ollama API returned an error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read Ollama API response")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse Ollama API response")
	}

	if responseText, ok := response["response"].(string); ok {
		fmt.Println("Ollama Response:", responseText)
	} else {
		fmt.Println("Ollama Response: unexpected format")
	}

	elapsed := time.Since(start)
	fmt.Printf("Request completed in %v\n", elapsed)
}
