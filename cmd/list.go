package cmd

import (
	"flawa/cfg"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var listFilesCmd = &cobra.Command{
	Use:   "list [directory]",
	Short: "List files in the specified directory or current working directory if no directory is provided",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var directory string
		if len(args) == 0 {
			cwd, err := os.Getwd()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get current working directory")
				return
			}
			directory = cwd
		} else {
			directory = args[0]
		}

		listFiles(directory)
	},
}

type ListConfig struct {
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

	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get current working directory")
			return
		}
		dir = cwd
		relativeDir := "./" + filepath.Base(cwd)
		fmt.Printf("Directory tree for %s (ignoring %v directories and %v files)\n", relativeDir, len(cfg.Config.Ignore.Directories), len(cfg.Config.Ignore.Files))
	} else {
		dir = expandHomePath(dir)
		relativeDir := formatDirForDisplay(dir)
		fmt.Printf("Directory tree for %s (ignoring %v directories and %v files)\n", relativeDir, len(cfg.Config.Ignore.Directories), len(cfg.Config.Ignore.Files))
	}

	checkDir, err := dirExists(dir)
	if err != nil {
		log.Error().Msgf("Error while looking for directory: %v", err)
		os.Exit(1)
	}
	if !checkDir {
		log.Warn().Msg("Directory doesn't exist")
		os.Exit(1)
	}

	var totalFiles, totalDirs int
	printTree(dir, "", &totalFiles, &totalDirs, cfg.Config.Ignore.Directories, cfg.Config.Ignore.Files)
	log.Info().Msgf("Done! Found %d files and %d directories\n", totalFiles, totalDirs)
}

func expandHomePath(path string) string {
	if path == "~" {
		usr, err := user.Current()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get current user")
			return path
		}
		return usr.HomeDir
	} else if len(path) > 1 && path[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get current user")
			return path
		}
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}

func formatDirForDisplay(dir string) string {
	usr, err := user.Current()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current user")
		return dir
	}
	homeDir := usr.HomeDir
	if len(dir) >= len(homeDir) && dir[:len(homeDir)] == homeDir {
		return "~" + dir[len(homeDir):]
	}
	return dir
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
