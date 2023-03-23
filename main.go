package main

import (
	"github.com/jacobmonck/metrics-collection/src/calico"
	"github.com/jacobmonck/metrics-collection/src/calico/listeners"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.TraceLevel)

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file, please ensure you have created one in the root directory.")
	}

	err = db.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database.")
	}

	logrus.Info("Initialization complete.")
}

func main() {
	config, err := utils.ParseConfig("./config/config.yaml")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
	}

	b, err := calico.New(config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create the calico struct.")
	}

	err = b.Setup(
		listeners.GuildReady(b),
		listeners.GuildMessageCreate(b),
		listeners.GuildMessageDelete(b),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Error setting up the calico.")
	}

	err = b.Start()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to the Discord Gateway.")
	}

	logrus.Info("Connected to the Discord gateway.")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}
