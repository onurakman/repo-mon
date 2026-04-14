package models

import "gorm.io/gorm"

type UserSettings struct {
	gorm.Model
	Theme              string `json:"theme" gorm:"default:neutral-carbon"`
	DarkMode           bool   `json:"darkMode" gorm:"default:true"`
	ViewMode           string `json:"viewMode" gorm:"default:grid"`
	GlobalPollInterval int    `json:"globalPollInterval" gorm:"default:30"`
	WindowWidth        int    `json:"windowWidth" gorm:"default:1200"`
	WindowHeight       int    `json:"windowHeight" gorm:"default:800"`
	WindowX            int    `json:"windowX" gorm:"default:-1"`
	WindowY            int    `json:"windowY" gorm:"default:-1"`
	WindowMaximised    bool   `json:"windowMaximised" gorm:"default:false"`
}
