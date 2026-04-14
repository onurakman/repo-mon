package models

import "gorm.io/gorm"

type UserSettings struct {
	gorm.Model
	Theme              string `json:"theme" gorm:"default:neutral-carbon"`
	DarkMode           bool   `json:"darkMode" gorm:"default:true"`
	ViewMode           string `json:"viewMode" gorm:"default:grid"`
	GlobalPollInterval int    `json:"globalPollInterval" gorm:"default:30"`
}
