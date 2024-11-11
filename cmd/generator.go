package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"flawa/cfg"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	filePathFlag      string
	directoryPathFlag string
)

var generateDocumentCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate document of specific file",
	Run: func(cmd *cobra.Command, args []string) {
		if filePathFlag != "" {
			generateDocument(parseInputPath(filePathFlag))
		} else if directoryPathFlag != "" {
			generateDocsFromDir(parseInputPath(directoryPathFlag))
		} else {
			fmt.Println("Specify either file (-f) or directory (-d)")
			os.Exit(1)
		}
	},
}

func init() {
	generateDocumentCmd.Flags().StringVarP(&filePathFlag, "file", "f", "", "Path to target file")
	generateDocumentCmd.Flags().StringVarP(&directoryPathFlag, "dir", "d", "", "Path to target directory")
}

func parseInputPath(normPath string) string {
	if normPath == "." {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err)
		}

		return pwd
	}

	return normPath
}

func generateDocsFromDir(directory string) {
	files := parseDirectory(directory)

	for _, file := range files {
		generateDocument(file)
	}

	log.Warn().Msg("This feature is under development")
}

func parseDirectory(directory string) []string {
	var files []string

	log.Info().Msg("Parsing" + " " + directory)
	// TODO: parse directory and get all file names

	return files
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

	fmt.Println("flawafying...")

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

	var responseText string
	if resText, ok := response["response"].(string); ok {
		elapsed := time.Since(start)
		fmt.Printf("flawafied in %v\n", elapsed)
		responseText = resText
	} else {
		fmt.Println("Unexpected format")
	}

	defaultOutputDir := filepath.Dir(filePath)
	fmt.Printf("Enter output directory (default: %s): ", defaultOutputDir)
	var outputDir string
	fmt.Scanln(&outputDir)
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	defaultOutputFile := fmt.Sprintf("%s-flawafied.md", strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)))
	fmt.Printf("Enter output filename (default: %s): ", defaultOutputFile)
	var outputFile string
	fmt.Scanln(&outputFile)
	if outputFile == "" {
		outputFile = defaultOutputFile
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := os.WriteFile(outputPath, []byte(responseText), 0644); err != nil {
		log.Fatal().Err(err).Msg("Failed to save the output file")
	} else {
		fmt.Printf("Output saved to %s\n", outputPath)
	}
}
