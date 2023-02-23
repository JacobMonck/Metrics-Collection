package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file, please ensure you have created one in the root directory.")
	}

	logrus.Info("Initialization complete.")
}

func main() {}
