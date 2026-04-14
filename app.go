package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"repo-mon/internal/database"
	"repo-mon/internal/models"
	"repo-mon/internal/monitor"
	"repo-mon/internal/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	scheduler *monitor.Scheduler
}

func NewApp() *App {
	return &App{
		scheduler: monitor.NewScheduler(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	dbDir := filepath.Join(configDir, "repo-mon")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		fmt.Println("Failed to create config dir:", err)
		return
	}
	dbPath := filepath.Join(dbDir, "repo-mon.db")

	if err := database.Initialize(dbPath); err != nil {
		fmt.Println("Database init error:", err)
		return
	}

	repos, err := service.GetRepositories()
	if err == nil {
		for _, repo := range repos {
			a.scheduler.Start(repo.ID, repo.Path, repo.PollInterval)
		}
	}
}

func (a *App) shutdown(ctx context.Context) {
	a.scheduler.StopAll()
}

// --- Window Controls ---

func (a *App) WindowMinimise() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) WindowToggleMaximise() {
	runtime.WindowToggleMaximise(a.ctx)
}

func (a *App) WindowClose() {
	runtime.Quit(a.ctx)
}

// --- Repository Management ---

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Git Repository",
	})
}

func (a *App) AddRepository(path string) (*models.Repository, error) {
	name := filepath.Base(path)
	settings, _ := service.GetSettings()
	interval := 30
	if settings != nil {
		interval = settings.GlobalPollInterval
	}

	repo, err := service.AddRepository(name, path, interval)
	if err != nil {
		return nil, err
	}

	a.scheduler.Start(repo.ID, repo.Path, repo.PollInterval)
	return repo, nil
}

func (a *App) RemoveRepository(id uint) error {
	a.scheduler.Stop(id)
	return service.RemoveRepository(id)
}

func (a *App) GetRepositories() ([]models.Repository, error) {
	return service.GetRepositories()
}

// --- Status ---

func (a *App) GetRepoStatus(id uint) (*monitor.RepoStatus, error) {
	status := a.scheduler.GetStatus(id)
	if status == nil {
		return nil, fmt.Errorf("no status for repo %d", id)
	}
	return status, nil
}

func (a *App) GetAllStatuses() map[uint]*monitor.RepoStatus {
	return a.scheduler.GetAllStatuses()
}

func (a *App) RefreshRepository(id uint) error {
	repo, err := service.GetRepository(id)
	if err != nil {
		return err
	}
	a.scheduler.Refresh(repo.ID, repo.Path)
	return nil
}

func (a *App) RefreshAll() error {
	repos, err := service.GetRepositories()
	if err != nil {
		return err
	}
	for _, repo := range repos {
		a.scheduler.Refresh(repo.ID, repo.Path)
	}
	return nil
}

func (a *App) UpdateSortOrder(ids []uint) error {
	return service.UpdateSortOrder(ids)
}

// --- Polling ---

func (a *App) UpdatePollInterval(id uint, seconds int) error {
	repo, err := service.GetRepository(id)
	if err != nil {
		return err
	}
	if err := service.UpdatePollInterval(id, seconds); err != nil {
		return err
	}
	a.scheduler.UpdateInterval(repo.ID, repo.Path, seconds)
	return nil
}

func (a *App) SetGlobalPollInterval(seconds int) error {
	settings, err := service.GetSettings()
	if err != nil {
		return err
	}
	settings.GlobalPollInterval = seconds
	return service.UpdateSettings(*settings)
}

// --- Tags ---

func (a *App) AddTag(name, color string) (*models.Tag, error) {
	return service.AddTag(name, color)
}

func (a *App) RemoveTag(id uint) error {
	return service.RemoveTag(id)
}

func (a *App) GetTags() ([]models.Tag, error) {
	return service.GetTags()
}

func (a *App) AssignTag(repoID, tagID uint) error {
	return service.AssignTag(repoID, tagID)
}

func (a *App) UnassignTag(repoID, tagID uint) error {
	return service.UnassignTag(repoID, tagID)
}

func (a *App) AssignTagToRepos(repoIDs []uint, tagID uint) error {
	return service.AssignTagToRepos(repoIDs, tagID)
}

// --- Settings ---

func (a *App) GetSettings() (*models.UserSettings, error) {
	return service.GetSettings()
}

func (a *App) UpdateSettings(settings models.UserSettings) error {
	return service.UpdateSettings(settings)
}
