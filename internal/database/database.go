package database

import (
	"repo-mon/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.Repository{}, &models.Tag{}, &models.UserSettings{})
	if err != nil {
		return err
	}

	var count int64
	DB.Model(&models.UserSettings{}).Count(&count)
	if count == 0 {
		DB.Create(&models.UserSettings{
			GlobalPollInterval: 30,
		})
	}

	return nil
}
