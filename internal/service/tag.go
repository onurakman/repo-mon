package service

import (
	"repo-mon/internal/database"
	"repo-mon/internal/models"
)

func AddTag(name, color string) (*models.Tag, error) {
	tag := &models.Tag{Name: name, Color: color}
	result := database.DB.Create(tag)
	if result.Error != nil {
		return nil, result.Error
	}
	return tag, nil
}

func RemoveTag(id uint) error {
	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		return err
	}
	if err := database.DB.Model(&tag).Association("Repositories").Clear(); err != nil {
		return err
	}
	return database.DB.Delete(&models.Tag{}, id).Error
}

func GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	result := database.DB.Find(&tags)
	return tags, result.Error
}

func AssignTag(repoID, tagID uint) error {
	var repo models.Repository
	if err := database.DB.First(&repo, repoID).Error; err != nil {
		return err
	}
	var tag models.Tag
	if err := database.DB.First(&tag, tagID).Error; err != nil {
		return err
	}
	return database.DB.Model(&repo).Association("Tags").Append(&tag)
}

func UnassignTag(repoID, tagID uint) error {
	var repo models.Repository
	if err := database.DB.First(&repo, repoID).Error; err != nil {
		return err
	}
	var tag models.Tag
	if err := database.DB.First(&tag, tagID).Error; err != nil {
		return err
	}
	return database.DB.Model(&repo).Association("Tags").Delete(&tag)
}
