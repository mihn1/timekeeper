package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

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
	config     AppConfig
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
	a.config = LoadAppConfig(os.Getenv)
	a.logger.Info("Loaded app config", "dbType", a.config.DBType, "dbPath", a.config.DBPath, "seedMode", a.config.SeedMode)

	a.timekeeper = a.newTimeKeeperFromConfig()
	a.seedIfNeeded()

	// Set up the macOS observer
	observer := macos.NewObserver(a.timekeeper.PushEvent, false, a.logger)
	a.timekeeper.AddObserver(observer)

	// Start tracking
	a.timekeeper.StartTracking()
}

func (a *App) newTimeKeeperFromConfig() *core.TimeKeeper {
	opts := core.TimeKeeperOptions{Logger: a.logger}

	if a.config.DBType == "inmem" {
		a.logger.Info("Starting TimeKeeper with in-memory storage")
		return core.NewTimeKeeperInMem(opts)
	}

	dbDir := filepath.Dir(a.config.DBPath)
	if dbDir != "" && dbDir != "." {
		if err := os.MkdirAll(dbDir, 0o755); err != nil {
			a.logger.Warn("Failed to create DB directory", "dir", dbDir, "error", err)
		}
	}

	opts.StoragePath = a.config.DBPath
	opts.StoreEvents = true
	a.logger.Info("Starting TimeKeeper with sqlite storage", "path", opts.StoragePath)
	return core.NewTimeKeeperSqlite(opts)
}

func (a *App) seedIfNeeded() {
	if a.timekeeper == nil {
		return
	}

	switch a.config.SeedMode {
	case "never":
		a.logger.Info("Skipping data seeding", "mode", a.config.SeedMode)
		return
	case "always":
		a.logger.Info("Seeding data", "mode", a.config.SeedMode)
		core.SeedData(a.timekeeper)
		return
	}

	categories, err := a.timekeeper.Storage.Categories().GetCategories()
	if err != nil {
		a.logger.Warn("Unable to inspect categories before seeding, skipping", "error", err)
		return
	}

	if len(categories) > 0 {
		a.logger.Info("Skipping data seeding: categories already exist", "count", len(categories))
		return
	}

	a.logger.Info("Seeding data", "mode", a.config.SeedMode)
	core.SeedData(a.timekeeper)
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
