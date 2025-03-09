package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Set up signal handling for clean shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		println("Received shutdown signal, cleaning up resources...")
		app.Shutdown(context.TODO())
	}()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "TimeKeeper",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             createMenu(app), // Add a menu
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []any{
			app,
		},
		// Add this to handle uncaught Go errors
		OnBeforeClose: func(ctx context.Context) bool {
			println("Application closing, cleaning up resources...")
			app.Shutdown(ctx)
			return false
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func createMenu(a *App) *menu.Menu {
	appMenu := menu.NewMenu()

	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Export Data...", nil, func(_ *menu.CallbackData) {
		// Export data functionality
	})

	// TimeKeeper menu
	tkMenu := appMenu.AddSubmenu("TimeKeeper")
	tkMenu.AddCheckbox("Enable Tracking", true, nil, func(cd *menu.CallbackData) {
		if cd.MenuItem.Checked {
			a.EnableTracking()
		} else {
			a.DisableTracking()
		}
	})

	// Debug menu
	debugMenu := appMenu.AddSubmenu("Debug")
	debugMenu.AddText("Force Cleanup", nil, func(_ *menu.CallbackData) {
		a.ForceCleanup()
	})

	return appMenu
}
