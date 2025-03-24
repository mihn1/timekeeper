package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func createMenu(app *App) *menu.Menu {
	appMenu := menu.NewMenu()

	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Export Data...", keys.CmdOrCtrl("E"), func(_ *menu.CallbackData) {
		// Implement export functionality
		app.logger.Info("Export functionality not yet implemented")
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Quit", keys.CmdOrCtrl("Q"), func(_ *menu.CallbackData) {
		runtime.Quit(app.ctx)
	})

	// TimeKeeper menu
	tkMenu := appMenu.AddSubmenu("TimeKeeper")
	tkMenu.AddText("Start Tracking", nil, func(_ *menu.CallbackData) {
		app.EnableTracking()
	})
	tkMenu.AddText("Stop Tracking", nil, func(_ *menu.CallbackData) {
		app.DisableTracking()
	})
	tkMenu.AddSeparator()
	tkMenu.AddText("Preferences...", keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
		// Implement preferences UI
		app.logger.Info("Preferences not yet implemented")
	})

	// View menu
	viewMenu := appMenu.AddSubmenu("View")
	viewMenu.AddText("Refresh Data", keys.CmdOrCtrl("R"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(app.ctx, "timekeeper:data-updated")
	})

	// Help menu
	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("About TimeKeeper", nil, func(_ *menu.CallbackData) {
		// Show about dialog
		app.logger.Info("About dialog not yet implemented")
	})

	return appMenu
}
