package service

import (
	"fmt"
	"repo-mon/internal/database"
	"repo-mon/internal/git"
	"repo-mon/internal/models"
)

func AddRepository(name, path string, pollInterval int) (*models.Repository, error) {
	if !git.IsGitRepo(path) {
		return nil, fmt.Errorf("not a git repository: %s", path)
	}
	if pollInterval <= 0 {
		pollInterval = 30
	}
	repo := &models.Repository{
		Name:         name,
		Path:         path,
		PollInterval: pollInterval,
	}
	result := database.DB.Create(repo)
	if result.Error != nil {
		return nil, result.Error
	}
	return repo, nil
}

func RemoveRepository(id uint) error {
	var repo models.Repository
	if err := database.DB.First(&repo, id).Error; err != nil {
		return err
	}
	if err := database.DB.Model(&repo).Association("Tags").Clear(); err != nil {
		return err
	}
	return database.DB.Delete(&models.Repository{}, id).Error
}

func GetRepositories() ([]models.Repository, error) {
	var repos []models.Repository
	result := database.DB.Preload("Tags").Order("sort_order asc, id asc").Find(&repos)
	return repos, result.Error
}

func UpdateSortOrder(ids []uint) error {
	tx := database.DB.Begin()
	for i, id := range ids {
		if err := tx.Model(&models.Repository{}).Where("id = ?", id).Update("sort_order", i).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func GetRepository(id uint) (*models.Repository, error) {
	var repo models.Repository
	result := database.DB.Preload("Tags").First(&repo, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &repo, nil
}

func UpdatePollInterval(id uint, seconds int) error {
	return database.DB.Model(&models.Repository{}).Where("id = ?", id).Update("poll_interval", seconds).Error
}

func SaveLastStatus(id uint, statusJSON string) error {
	return database.DB.Model(&models.Repository{}).Where("id = ?", id).Update("last_status", statusJSON).Error
}
