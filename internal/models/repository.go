package models

import "gorm.io/gorm"

type Repository struct {
	gorm.Model
	Name         string `json:"name"`
	Path         string `json:"path" gorm:"uniqueIndex"`
	PollInterval int    `json:"pollInterval" gorm:"default:30"`
	Tags         []Tag  `json:"tags" gorm:"many2many:repository_tags;"`
}
