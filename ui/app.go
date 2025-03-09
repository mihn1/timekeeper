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

	// Create a logger with the Wails handler
	a.logger = slog.New(wailsHandler)

	// Set as default logger
	slog.SetDefault(a.logger)

	// Debug message to verify logger is working
	runtime.LogInfo(ctx, "Wails logger initialized")
	a.logger.Info("Testing wailsLogger directly")

	go a.initTimekeeper()
}

func (a *App) Shutdown(ctx context.Context) {
	if a.timekeeper != nil {
		a.timekeeper.Close()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	a.logger.Info("Hello {Name}, It's show time!", "Name", name)
	return name
}

// Expose TimeKeeper functionality to JavaScript
func (a *App) GetAppUsageData(dateStr string) interface{} {
	date, _ := datatypes.NewDateOnlyFromStr(dateStr)
	data, _ := a.timekeeper.Storage.AppAggregations().GetAppAggregationsByDate(date)
	return data
}

func (a *App) initTimekeeper() {
	opts := core.TimeKeeperOptions{
		StoreEvents: true,
		StoragePath: "../db/timekeeper_wails.db", // Use relative path to existing DB
	}

	// Create a new TimeKeeper instance
	a.timekeeper = core.NewTimeKeeperSqlite(opts)

	a.timekeeper.SetLogger(a.logger)
	observer := macos.NewObserver(a.timekeeper.PushEvent, a.logger)
	a.timekeeper.AddObserver(observer)
	core.SeedData(a.timekeeper)
	a.timekeeper.StartTracking()
}

// Add more methods to expose TimeKeeper functionality...
