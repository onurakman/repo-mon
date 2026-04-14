package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"repo-mon/internal/database"
	"repo-mon/internal/service"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	app := NewApp()

	// Pre-init DB to read saved window state
	width, height := 1200, 800
	startMaximised := false

	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	dbDir := filepath.Join(configDir, "repo-mon")
	_ = os.MkdirAll(dbDir, 0755)
	dbPath := filepath.Join(dbDir, "repo-mon.db")

	if err := database.Initialize(dbPath); err == nil {
		if settings, err := service.GetSettings(); err == nil {
			if settings.WindowWidth > 0 && settings.WindowHeight > 0 {
				width = settings.WindowWidth
				height = settings.WindowHeight
			}
			startMaximised = settings.WindowMaximised
		}
	} else {
		fmt.Println("Database init error:", err)
	}

	err = wails.Run(&options.App{
		Title:            "Repo Monitor",
		Width:            width,
		Height:           height,
		MinWidth:         800,
		MinHeight:        600,
		StartHidden:      false,
		Frameless:        true,
		DisableResize:    false,
		BackgroundColour: &options.RGBA{R: 23, G: 23, B: 23, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:     app.startup,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			Icon:                appIcon,
			WindowIsTranslucent: false,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "Repo Monitor",
				Message: "Git Repository Monitor",
				Icon:    appIcon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	// Save maximised state after window closes
	_ = startMaximised // used in startup

	if err != nil {
		println("Error:", err.Error())
	}
}
