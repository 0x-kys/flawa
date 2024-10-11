package main

import (
	"flawa/cmd"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
}

func main() {
	cmd.Execute()
}
