package cfg

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type FlawaConfig struct {
	Ignore struct {
		Directories []string `toml:"directories"`
		Files       []string `toml:"files"`
	} `toml:"ignore"`

	Ollama struct {
		Model   string `toml:"model"`
		Stream  bool   `toml:"stream"`
		NumGpu  int    `toml:"num-gpu"`
		MainGpu int    `toml:"main-gpu"`
	} `toml:"ollama"`
}

var Config FlawaConfig

func LoadConfig() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	configPath := filepath.Join(usr.HomeDir, ".config", "flawa", "config.toml")

	if _, err := toml.DecodeFile(configPath, &Config); err != nil {
		return fmt.Errorf("failed to decode TOML file: %w", err)
	}

	return nil
}
