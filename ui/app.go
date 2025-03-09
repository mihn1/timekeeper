package main

import (
	"context"
	"log/slog"

	"github.com/mihn1/timekeeper/core"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/macos"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	timekeeper *core.TimeKeeper
	logger     *slog.Logger
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Create Wails slog handler
	wailsHandler := NewWailsHandler(ctx)
	a.logger = slog.New(wailsHandler)
	slog.SetDefault(a.logger)

	// Initialize TimeKeeper directly (not in goroutine)
	a.initTimekeeper()

	// Debug message to verify logger is working
	runtime.LogInfo(ctx, "TimeKeeper initialized")
}

func (a *App) Shutdown(ctx context.Context) {
	if a.timekeeper != nil {
		a.logger.Info("Shutting down TimeKeeper...")
		a.timekeeper.Close()
		a.timekeeper = nil // Prevent double-close
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	a.logger.Info("Hello {Name}, It's show time!", "Name", name)
	return name
}

// Expose TimeKeeper functionality to JavaScript
func (a *App) GetAppUsageData(dateStr string) any {
	date, _ := datatypes.NewDateOnlyFromStr(dateStr)
	data, _ := a.timekeeper.Storage.AppAggregations().GetAppAggregationsByDate(date)
	return data
}

func (a *App) EnableTracking() {
	if a.timekeeper != nil && !a.timekeeper.IsEnabled() {
		a.logger.Info("Enabling TimeKeeper tracking")
		a.timekeeper.StartTracking()
	}
}

func (a *App) DisableTracking() {
	if a.timekeeper != nil && a.timekeeper.IsEnabled() {
		a.logger.Info("Disabling TimeKeeper tracking")
		a.timekeeper.Disable()
	}
}

// Add this method to be called from JS
func (a *App) IsTrackingEnabled() bool {
	if a.timekeeper != nil {
		return a.timekeeper.IsEnabled()
	}
	return false
}

func (a *App) ForceCleanup() {
	a.logger.Info("Force cleaning up resources...")
	a.Shutdown(a.ctx)
}

func (a *App) initTimekeeper() {
	opts := core.TimeKeeperOptions{
		StoreEvents: true,
		StoragePath: "../db/timekeeper_wails.db",
		Logger:      a.logger,
	}

	// Create a new TimeKeeper instance
	a.timekeeper = core.NewTimeKeeperSqlite(opts)
	core.SeedData(a.timekeeper)

	// Set up the macOS observer
	observer := macos.NewObserver(a.timekeeper.PushEvent, false, a.logger)
	a.timekeeper.AddObserver(observer)

	// Start tracking
	a.timekeeper.StartTracking()
}

// Add more methods to expose TimeKeeper functionality...
