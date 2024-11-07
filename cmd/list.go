package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var listFilesCmd = &cobra.Command{
	Use:   "list",
	Short: "Get info of a specific repo by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listFiles(args[0])
	},
}

type Config struct {
	Ignore struct {
		Directories []string `toml:"directories"`
		Files       []string `toml:"files"`
	} `toml:"ignore"`
}

func dirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) || !stat.IsDir() {
		return false, nil
	}
	return false, err
}

func listFiles(arg string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	dir := arg

	checkDir, err := dirExists(dir)
	if err != nil {
		log.Error().Msgf("Error while looking for directory: %v", err)
		os.Exit(1)
	}

	if !checkDir {
		log.Warn().Msg("Directory doesn't exist")
		os.Exit(1)
	}

	config, err := loadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load config file")
		os.Exit(1)
	}

	fmt.Printf("Directory tree for %s (ignoring %v directories and %v files)\n", dir, len(config.Ignore.Directories), len(config.Ignore.Files))
	var totalFiles, totalDirs int
	printTree(dir, "", &totalFiles, &totalDirs, config.Ignore.Directories, config.Ignore.Files)
	log.Info().Msgf("Done! Found %d files and %d directories\n", totalFiles, totalDirs)
}

func loadConfig() (Config, error) {
	var config Config
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "flawa", "config.toml")
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return config, err
	}
	return config, nil
}

func printTree(path, indent string, fileCount, dirCount *int, ignoreDirs, ignoreFiles []string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Error().Err(err).Msgf("Could not read directory: %s", path)
		return
	}

	for i, entry := range entries {
		if entry.IsDir() && contains(ignoreDirs, entry.Name()) {
			continue
		}
		if !entry.IsDir() && contains(ignoreFiles, entry.Name()) {
			continue
		}

		isLast := i == len(entries)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		fmt.Printf("%s%s%s\n", indent, prefix, entry.Name())

		if entry.IsDir() {
			*dirCount++
			newIndent := indent + "│   "
			if isLast {
				newIndent = indent + "    "
			}
			printTree(filepath.Join(path, entry.Name()), newIndent, fileCount, dirCount, ignoreDirs, ignoreFiles)
		} else {
			*fileCount++
		}
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
