package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mihn1/timekeeper/core"
	"github.com/mihn1/timekeeper/datatypes"
)

// App struct
type App struct {
	ctx        context.Context
	timekeeper *core.TimeKeeper
	logger     *slog.Logger
}

// NewApp creates a new App application struct
func NewApp() *App {
	opts := core.TimeKeeperOptions{
		StoreEvents: true,
		StoragePath: "../db/timekeeper.db", // Use relative path to existing DB
	}

	// Create a new TimeKeeper instance
	tk := core.NewTimeKeeperInMem(opts)

	return &App{
		timekeeper: tk,
	}
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

	// Set the logger in TimeKeeper
	a.timekeeper.SetLogger(a.logger)

	// Now start tracking
	a.timekeeper.StartTracking()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Expose TimeKeeper functionality to JavaScript
func (a *App) GetAppUsageData(dateStr string) interface{} {
	date, _ := datatypes.NewDateOnlyFromStr(dateStr)
	data, _ := a.timekeeper.Storage.AppAggregations().GetAppAggregationsByDate(date)
	return data
}

// Add more methods to expose TimeKeeper functionality...
