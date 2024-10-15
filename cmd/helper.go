package cmd

import (
	"log"
	"os"
	"path/filepath"
)

func GetConfigPath(filename string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(homeDir, ".config", "flawa", filename)
}
