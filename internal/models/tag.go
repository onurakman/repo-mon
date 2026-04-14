package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name         string       `json:"name" gorm:"uniqueIndex"`
	Color        string       `json:"color"`
	Repositories []Repository `json:"repositories" gorm:"many2many:repository_tags;"`
}
