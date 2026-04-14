package service

import (
	"repo-mon/internal/database"
	"repo-mon/internal/models"
)

func GetSettings() (*models.UserSettings, error) {
	var settings models.UserSettings
	result := database.DB.First(&settings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &settings, nil
}

func UpdateSettings(settings models.UserSettings) error {
	return database.DB.Model(&models.UserSettings{}).Where("id = ?", 1).Updates(map[string]interface{}{
		"theme":                settings.Theme,
		"dark_mode":            settings.DarkMode,
		"view_mode":            settings.ViewMode,
		"global_poll_interval": settings.GlobalPollInterval,
	}).Error
}
