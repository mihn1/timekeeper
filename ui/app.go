package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/mihn1/timekeeper/core"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/internal/tzutil"
	"github.com/mihn1/timekeeper/platforms"
	"github.com/mihn1/timekeeper/ui/dtos"
)

// App struct
type App struct {
	ctx        context.Context
	timekeeper *core.TimeKeeper
	logger     *slog.Logger
	config     AppConfig

	prefsMu sync.RWMutex
	prefs   *models.UserPreferences
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

	// Load user preferences from storage.
	if prefs, err := a.timekeeper.Storage.Preferences().GetPreferences(); err == nil {
		a.prefs = prefs
	} else {
		a.logger.Warn("Failed to load preferences, using defaults", "error", err)
		a.prefs = models.DefaultPreferences()
	}
	a.timekeeper.SetMinEventDurationMs(a.prefs.MinEventDurationMs)

	// Set up the platform observer
	observer := platforms.NewPlatformObserver(a.timekeeper.PushEvent, false, a.logger)
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

// getTimezone returns the user's configured IANA timezone string.
func (a *App) getTimezone() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	if a.prefs != nil && a.prefs.Timezone != "" {
		return a.prefs.Timezone
	}
	return "UTC"
}

// categoryNameMap returns a map of CategoryId → display name.
func (a *App) categoryNameMap() map[models.CategoryId]string {
	m := make(map[models.CategoryId]string)
	if cats, err := a.timekeeper.Storage.Categories().GetCategories(); err == nil {
		for _, cat := range cats {
			m[cat.Id] = cat.Name
		}
	}
	return m
}

// buildAppUsageItemsFromEvents aggregates per-app time from raw events.
func (a *App) buildAppUsageItemsFromEvents(events []*models.AppSwitchEvent) ([]*dtos.AppUsageItem, error) {
	appTotals := tzutil.AggregateEventsByApp(events)
	appCats := tzutil.AppCategoryMap(events)
	catNames := a.categoryNameMap()

	result := make([]*dtos.AppUsageItem, 0, len(appTotals))
	for appName, ms := range appTotals {
		catId := appCats[appName]
		catName := catNames[catId]
		if catName == "" {
			catName = "Undefined"
		}
		result = append(result, &dtos.AppUsageItem{
			AppName:      appName,
			TimeElapsed:  ms,
			CategoryId:   int(catId),
			CategoryName: catName,
		})
	}
	return result, nil
}

// buildAppUsageItemsFromAggr falls back to aggregation tables (UTC-date-keyed).
// Used in inmem mode where raw events are not persisted.
func (a *App) buildAppUsageItemsFromAggr(date datatypes.DateOnly) ([]*dtos.AppUsageItem, error) {
	aggregations, err := a.timekeeper.Storage.AppAggregations().GetAppAggregationsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to load app usage data: %w", err)
	}

	appCatMap := make(map[string]models.CategoryId)
	events, _ := a.timekeeper.Storage.Events().GetEventsByDate(date)
	for _, ev := range events {
		appCatMap[ev.AppName] = ev.CategoryId
	}

	catNames := a.categoryNameMap()
	result := make([]*dtos.AppUsageItem, 0, len(aggregations))
	for _, aggr := range aggregations {
		catId, ok := appCatMap[aggr.AppName]
		if !ok {
			catId = models.UNDEFINED
		}
		catName := catNames[catId]
		if catName == "" {
			catName = "Undefined"
		}
		result = append(result, &dtos.AppUsageItem{
			AppName:      aggr.AppName,
			TimeElapsed:  aggr.TimeElapsed,
			CategoryId:   int(catId),
			CategoryName: catName,
		})
	}
	return result, nil
}

// GetAppUsageData returns per-app time totals for the given local date, enriched with category info.
// Uses the user's timezone preference to compute the correct UTC day window.
func (a *App) GetAppUsageData(dateStr string) ([]*dtos.AppUsageItem, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	tz := a.getTimezone()
	start, end, err := tzutil.LocalDayToUTCRange(dateStr, tz)
	if err != nil {
		return nil, fmt.Errorf("invalid date %q: %w", dateStr, err)
	}

	events, qErr := a.timekeeper.Storage.Events().GetEventsByTimeRange(start, end)
	if qErr == nil && len(events) > 0 {
		return a.buildAppUsageItemsFromEvents(events)
	}

	// Fall back to aggregation tables (inmem mode or first-run before events are stored).
	date, parseErr := datatypes.NewDateOnlyFromStr(dateStr)
	if parseErr != nil {
		return nil, fmt.Errorf("invalid date %q: %w", dateStr, parseErr)
	}
	return a.buildAppUsageItemsFromAggr(date)
}

// GetUncategorizedApps returns app names that resolved to the UNDEFINED category on the given date.
func (a *App) GetUncategorizedApps(dateStr string) ([]string, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	items, err := a.GetAppUsageData(dateStr)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, item := range items {
		if models.CategoryId(item.CategoryId) == models.UNDEFINED {
			result = append(result, item.AppName)
		}
	}
	return result, nil
}

func (a *App) GetCategoryUsageData(dateStr string) ([]*dtos.CategoryUsageItem, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	tz := a.getTimezone()
	catNames := a.categoryNameMap()

	// Try timezone-accurate aggregation from raw events first.
	start, end, err := tzutil.LocalDayToUTCRange(dateStr, tz)
	if err == nil {
		events, qErr := a.timekeeper.Storage.Events().GetEventsByTimeRange(start, end)
		if qErr == nil && len(events) > 0 {
			catTotals := tzutil.AggregateEventsByCategory(events)
			result := make([]*dtos.CategoryUsageItem, 0, len(catTotals))
			for catId, ms := range catTotals {
				name := catNames[catId]
				if name == "" {
					name = "Undefined"
				}
				result = append(result, &dtos.CategoryUsageItem{
					Id:          int(catId),
					Name:        name,
					TimeElapsed: ms,
				})
			}
			return result, nil
		}
	}

	// Fall back to pre-computed aggregation table.
	date, parseErr := datatypes.NewDateOnlyFromStr(dateStr)
	if parseErr != nil {
		return nil, fmt.Errorf("invalid date %q: %w", dateStr, parseErr)
	}
	data, err := a.timekeeper.Storage.CategoryAggregations().GetCategoryAggregationsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to load category usage data: %w", err)
	}

	result := make([]*dtos.CategoryUsageItem, 0, len(data))
	for _, catAggr := range data {
		name := catNames[catAggr.CategoryId]
		if name == "" {
			name = "Undefined"
		}
		result = append(result, &dtos.CategoryUsageItem{
			Id:          int(catAggr.CategoryId),
			Name:        name,
			TimeElapsed: catAggr.TimeElapsed,
		})
	}
	return result, nil
}

// GetCategoryUsageTotals returns category time totals summed across a local date range.
// Uses the user's timezone to compute accurate UTC boundaries, then queries raw events
// (same approach as single-day queries) so UTC+ users get correct data.
func (a *App) GetCategoryUsageTotals(startDate, endDate string) ([]*dtos.CategoryUsageItem, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	tz := a.getTimezone()

	// Convert local start date → UTC start of that local day.
	utcStart, _, err := tzutil.LocalDayToUTCRange(startDate, tz)
	if err != nil {
		return nil, fmt.Errorf("invalid start date %q: %w", startDate, err)
	}
	// Convert local end date → UTC end of that local day.
	_, utcEnd, err := tzutil.LocalDayToUTCRange(endDate, tz)
	if err != nil {
		return nil, fmt.Errorf("invalid end date %q: %w", endDate, err)
	}

	catNames := a.categoryNameMap()

	// Try raw events first — timezone-accurate.
	events, qErr := a.timekeeper.Storage.Events().GetEventsByTimeRange(utcStart, utcEnd)
	if qErr == nil && len(events) > 0 {
		catTotals := tzutil.AggregateEventsByCategory(events)
		result := make([]*dtos.CategoryUsageItem, 0, len(catTotals))
		for catId, ms := range catTotals {
			name := catNames[catId]
			if name == "" {
				name = "Undefined"
			}
			result = append(result, &dtos.CategoryUsageItem{
				Id:          int(catId),
				Name:        name,
				TimeElapsed: ms,
			})
		}
		return result, nil
	}

	// Fall back to aggregation tables. Convert UTC boundaries back to DateOnly for the query.
	startDateOnly := datatypes.NewDateOnly(utcStart)
	endDateOnly := datatypes.NewDateOnly(utcEnd.Add(-time.Second))
	aggrs, err := a.timekeeper.Storage.CategoryAggregations().GetCategoryAggregationsByDateRange(startDateOnly, endDateOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to load category totals: %w", err)
	}

	totals := make(map[models.CategoryId]int64)
	for _, aggr := range aggrs {
		totals[aggr.CategoryId] += aggr.TimeElapsed
	}

	result := make([]*dtos.CategoryUsageItem, 0, len(totals))
	for catId, ms := range totals {
		name := catNames[catId]
		if name == "" {
			name = "Undefined"
		}
		result = append(result, &dtos.CategoryUsageItem{
			Id:          int(catId),
			Name:        name,
			TimeElapsed: ms,
		})
	}
	return result, nil
}

// GetCategoryUsageRange returns per-category daily summaries for a date range (trend chart data).
func (a *App) GetCategoryUsageRange(startDate, endDate string) ([]*dtos.DailyCategorySummary, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	start, err := datatypes.NewDateOnlyFromStr(startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date %q: %w", startDate, err)
	}
	end, err := datatypes.NewDateOnlyFromStr(endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date %q: %w", endDate, err)
	}

	aggrs, err := a.timekeeper.Storage.CategoryAggregations().GetCategoryAggregationsByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to load category range data: %w", err)
	}

	categoryNameMap := make(map[models.CategoryId]string)
	if cats, err := a.timekeeper.Storage.Categories().GetCategories(); err == nil {
		for _, cat := range cats {
			categoryNameMap[cat.Id] = cat.Name
		}
	}

	result := make([]*dtos.DailyCategorySummary, 0, len(aggrs))
	for _, aggr := range aggrs {
		catName := categoryNameMap[aggr.CategoryId]
		if catName == "" {
			catName = "Undefined"
		}
		result = append(result, &dtos.DailyCategorySummary{
			Date:         aggr.Date.String(),
			CategoryId:   int(aggr.CategoryId),
			CategoryName: catName,
			TimeElapsed:  aggr.TimeElapsed,
		})
	}

	// Sort by date ascending.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result, nil
}

// GetActivityCalendar returns daily activity totals for the heatmap calendar.
func (a *App) GetActivityCalendar(year int) ([]*dtos.DayActivity, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	start := datatypes.NewDateOnly(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC))
	end := datatypes.NewDateOnly(time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC))

	aggrs, err := a.timekeeper.Storage.CategoryAggregations().GetCategoryAggregationsByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to load calendar data: %w", err)
	}

	// Group by date: sum total ms, track max-time category per day (excluding EXCLUDED).
	type dayStats struct {
		totalMs       int64
		catTimes      map[models.CategoryId]int64
	}
	byDate := make(map[string]*dayStats)

	for _, aggr := range aggrs {
		dateStr := aggr.Date.String()
		ds, ok := byDate[dateStr]
		if !ok {
			ds = &dayStats{catTimes: make(map[models.CategoryId]int64)}
			byDate[dateStr] = ds
		}
		ds.totalMs += aggr.TimeElapsed
		if aggr.CategoryId != models.EXCLUDED {
			ds.catTimes[aggr.CategoryId] += aggr.TimeElapsed
		}
	}

	result := make([]*dtos.DayActivity, 0, len(byDate))
	for dateStr, ds := range byDate {
		topCatId := int(models.UNDEFINED)
		var topMs int64
		for catId, ms := range ds.catTimes {
			if ms > topMs {
				topMs = ms
				topCatId = int(catId)
			}
		}
		result = append(result, &dtos.DayActivity{
			Date:          dateStr,
			TotalMs:       ds.totalMs,
			TopCategoryId: topCatId,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result, nil
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
