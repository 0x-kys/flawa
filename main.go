package main

import (
	"flawa/cmd"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load(cmd.GetConfigPath(".env"))
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	logrus.Println("Loaded .env")
	logrus.Println(os.Getenv("CLIENT_ID"))
}

func main() {
	cmd.Execute()
}
