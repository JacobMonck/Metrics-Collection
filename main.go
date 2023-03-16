package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/bot"
	"github.com/jacobmonck/metrics-collection/src/utils"
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

	err = utils.ParseConfig("./config/config.yaml")
	if err != nil {
		logrus.Fatalf("Failed to load config: %s", err)
	}

	err = db.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database.")
	}

	logrus.Info("Initialization complete.")
}

func main() {
	client, err := bot.Start()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start bot.")
	}
	logrus.Info("Connected to the Discord gateway.")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	<-exit
	client.Close(context.Background())
	logrus.Info("Bot connections are closed.")
}
