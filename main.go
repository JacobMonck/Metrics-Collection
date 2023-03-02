package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/bot"
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

	err = db.Init()
	if err != nil {
		panic(err)
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

	select {
	case <-exit:
		client.Close(context.Background())
		logrus.Info("Bot connections are closed.")
	}
}
