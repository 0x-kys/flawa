package main

import (
	"flawa/cfg"
	"flawa/cmd"
	"fmt"
)

func main() {
	if err := cfg.LoadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	cmd.Execute()
}

