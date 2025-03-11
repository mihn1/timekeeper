package main

import (
	"context"
	"log/slog"

	"github.com/mihn1/timekeeper/core"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/macos"
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
	var wailsHandler slog.Handler = NewWailsHandler(ctx)
	a.logger = slog.New(wailsHandler)
	slog.SetDefault(a.logger)

	// Initialize TimeKeeper directly (not in goroutine)
	opts := core.TimeKeeperOptions{
		StoreEvents: true,
		StoragePath: "../db/timekeeper.db",
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

func (a *App) Shutdown(ctx context.Context) {
	if a.timekeeper != nil {
		a.logger.Info("Shutting down TimeKeeper...")
		a.timekeeper.Close()
		a.timekeeper = nil // Prevent double-close
	}
}

// Expose TimeKeeper functionality to JavaScript
func (a *App) GetAppUsageData(dateStr string) []*models.AppAggregation {
	date, _ := datatypes.NewDateOnlyFromStr(dateStr)
	data, _ := a.timekeeper.Storage.AppAggregations().GetAppAggregationsByDate(date)
	return data
}

func (a *App) GetCategoryUsageData(dateStr string) any {
	date, _ := datatypes.NewDateOnlyFromStr(dateStr)
	data, _ := a.timekeeper.Storage.CategoryAggregations().GetCategoryAggregationsByDate(date)

	// Enrich with category names
	result := make([]map[string]any, 0, len(data))
	for _, catAggr := range data {
		cat, _ := a.timekeeper.Storage.Categories().GetCategory(catAggr.CategoryId)
		result = append(result, map[string]any{
			"Id":          catAggr.CategoryId,
			"Name":        cat.Name,
			"TimeElapsed": catAggr.TimeElapsed,
		})
	}

	return result
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

// Add more methods to expose TimeKeeper functionality...
