package db

import (
	"github.com/jacobmonck/metrics-collection/src/api/db/models"
	"github.com/jacobmonck/metrics-collection/src/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Session *gorm.DB

func Init() error {
	db, err := gorm.Open(
		postgres.Open(utils.EnvOr(
			"DB_DSN",
			"postgresql://postgres:postgres@localhost:5432/metrics-collection",
		)),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return err
	}

	Session = db

	logrus.Info("Successfully connected to database.")

	err = db.AutoMigrate(
		&models.Category{},
		&models.Channel{},
		&models.Thread{},
		&models.User{},
		&models.Message{},
	)
	if err != nil {
		return err
	}

	logrus.Info("Successfully migrated database.")

	return nil
}
