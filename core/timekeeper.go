package core

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/core/resolvers"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/data/inmem"
	"github.com/mihn1/timekeeper/internal/data/interfaces"
	"github.com/mihn1/timekeeper/internal/data/sqlite"
	"github.com/mihn1/timekeeper/internal/models"
)

const (
	AppName = "TimeKeeper"
)

var (
	defaultResolver resolvers.CategoryResolver
)

type TimeKeeperOptions struct {
	StoreEvents bool
	StoragePath string
	Logger      *slog.Logger
}

type TimeKeeper struct {
	curAppEvent  *models.AppSwitchEvent
	opts         TimeKeeperOptions
	Storage      interfaces.Storage
	isEnabled    bool
	eventChannel chan models.AppSwitchEvent
	observers    []Observer
	logger       *slog.Logger
}

func NewTimeKeeperInMem(opts TimeKeeperOptions) *TimeKeeper {
	t := &TimeKeeper{
		Storage:      inmem.NewInmemStorage(),
		eventChannel: make(chan models.AppSwitchEvent),
		opts:         opts,
		logger:       opts.Logger,
	}
	return t
}

func NewTimeKeeperSqlite(opts TimeKeeperOptions) *TimeKeeper {
	if opts.StoragePath == "" {
		opts.StoragePath = "./timekeeper.db"
	}

	db, err := sql.Open("sqlite3", opts.StoragePath)
	if err != nil {
		opts.Logger.Error("Error opening database", "error", err, "path", opts.StoragePath)
		os.Exit(1) // Maintain fatal behavior
	}

	db.SetMaxOpenConns(1)
	_, err = db.Exec("PRAGMA busy_timeout = 5000;")
	if err != nil {
		opts.Logger.Error("Error setting busy_timeout", "error", err)
		os.Exit(1) // Maintain fatal behavior
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		opts.Logger.Error("Error setting journal_mode", "error", err)
		os.Exit(1) // Maintain fatal behavior
	}

	t := &TimeKeeper{
		Storage:      sqlite.NewSqliteStorage(db),
		eventChannel: make(chan models.AppSwitchEvent),
		opts:         opts,
		logger:       opts.Logger,
	}
	return t
}

func (t *TimeKeeper) SetLogger(logger *slog.Logger) {
	t.logger = logger
}

func (t *TimeKeeper) Logger() *slog.Logger {
	return t.logger
}

func (t *TimeKeeper) AddObserver(o Observer) {
	if t.observers == nil {
		t.observers = make([]Observer, 0)
	}
	t.observers = append(t.observers, o)
}

func (t *TimeKeeper) Close() {
	t.logger.Info("Closing TimeKeeper...")
	t.Disable()
	t.Storage.Close()

	if t.observers != nil {
		for _, obs := range t.observers {
			obs.Stop()
		}
	}
}

func (t *TimeKeeper) Disable() {
	t.isEnabled = false
}

func (t *TimeKeeper) IsEnabled() bool {
	return t.isEnabled
}

func (t *TimeKeeper) StartTracking() {
	t.logger.Info("Starting TimeKeeper...")

	t.isEnabled = true
	defaultResolver = resolvers.NewDefaultCategoryResolver(t.Storage.Rules(), t.Storage.Categories())

	// Start all observers
	if t.observers != nil {
		for _, obs := range t.observers {
			go obs.Start()
		}
	}

	// Start listening for events
	go func() {
		for event := range t.eventChannel {
			t.handleEvent(&event)
		}
	}()

	t.logger.Info("TimeKeeper started")
}

func (t *TimeKeeper) Report(date datatypes.DateOnly) {
	t.logger.Info("TimeKeeper Report", "date", date)

	appAggrs, _ := t.Storage.AppAggregations().GetAppAggregationsByDate(date)
	catAggrs, _ := t.Storage.CategoryAggregations().GetCategoryAggregationsByDate(date)

	t.logger.Info("App Aggregation", "data", appAggrs)
	t.logger.Info("Category Aggregation", "data", catAggrs)
	t.logger.Info("-------------------------------------------------------------------------------------------------")
}

func (t *TimeKeeper) PushEvent(event models.AppSwitchEvent) {
	if t.IsEnabled() {
		t.eventChannel <- event
	}
}

func (t *TimeKeeper) handleEvent(event *models.AppSwitchEvent) {
	t.logger.Debug("Received event", "event", event)

	if t.curAppEvent == nil {
		t.curAppEvent = event
		return
	}

	// TODO: gracefully handle the case when the timekeeper is disabled
	if !t.isEnabled {
		return
	}

	if isSameEvent(t.curAppEvent, event) {
		t.logger.Debug("Same event detected", "event", event)
		return
	}

	t.aggregateNewEvent(event)
	t.curAppEvent.EndTime = event.StartTime

	// store the current app event
	if t.opts.StoreEvents {
		err := t.Storage.Events().AddEvent(t.curAppEvent)
		if err != nil {
			t.logger.Error("Error storing event", "error", err)
		}
	}

	t.curAppEvent = event
	t.Report(datatypes.NewDateOnly(event.StartTime))
}

func (t *TimeKeeper) aggregateNewEvent(event *models.AppSwitchEvent) {
	elapsedTime := event.StartTime.Sub(t.curAppEvent.StartTime).Milliseconds()

	rules, err := t.getRulesForEvent(t.curAppEvent)
	if err != nil {
		t.logger.Error("Error getting rules for event", "error", err)
		return
	}

	_, err = t.Storage.AppAggregations().AggregateAppEvent(t.curAppEvent, elapsedTime)
	if err != nil {
		t.logger.Error("Error aggregating app event", "app", event.AppName, "error", err)
		return
	}

	catId, err := t.aggregateCategory(t.curAppEvent, elapsedTime, rules) // Call after aggregateEvent
	if err != nil {
		t.logger.Error("Error aggregating category", "app", event.AppName, "error", err)
		return
	}

	t.curAppEvent.CategoryId = catId
}

func (t *TimeKeeper) aggregateCategory(event *models.AppSwitchEvent, elapsedTime int64, rules []*models.CategoryRule) (models.CategoryId, error) {
	catId, err := defaultResolver.ResolveCategory(event, rules)
	if err != nil {
		return catId, err
	}

	if catId == models.EXCLUDED {
		t.logger.Info("Category excluded", "event", event)
		return catId, nil
	}

	cat, err := t.Storage.Categories().GetCategory(catId)
	if err != nil {
		t.logger.Error("Error getting category", "error", err)
		return catId, err
	}

	t.logger.Info("Category resolved", "CatName", cat.Name, "event", event)

	_, err = t.Storage.CategoryAggregations().AggregateCategory(cat, event.GetEventDate(), elapsedTime)
	if err != nil {
		t.logger.Error("Error saving category aggregation", "error", err)
		return catId, err
	}

	return catId, nil
}

func (t *TimeKeeper) getRulesForEvent(event *models.AppSwitchEvent) ([]*models.CategoryRule, error) {
	rules, err := t.Storage.Rules().GetRulesByApp(event.AppName)
	if err != nil {
		return nil, err
	}

	if isEventExcluded(event, rules) {
		return nil, fmt.Errorf("event excluded")
	}

	// Get rules that are applied to all apps
	globalRules, err := t.Storage.Rules().GetRulesByApp(constants.ALL_APPS)
	if err == nil {
		if isEventExcluded(event, globalRules) {
			return nil, fmt.Errorf("event excluded")
		}
	}

	rules = append(rules, globalRules...)
	slices.SortStableFunc(globalRules, models.CmpRules)

	return rules, nil
}

func isSameEvent(e1, e2 *models.AppSwitchEvent) bool {
	if e1.AppName != e2.AppName {
		return false
	}

	for key, val := range e1.AdditionalData {
		if val2, ok := e2.AdditionalData[key]; !ok || val != val2 {
			return false
		}
	}

	// If the events are within 60 seconds of each other, consider them the same event
	return e1.StartTime.Sub(e2.StartTime).Seconds() <= 60
}

func isEventExcluded(event *models.AppSwitchEvent, rules []*models.CategoryRule) bool {
	for _, rule := range rules {
		if !rule.IsExclusion {
			continue
		}

		log.Printf("[DEBUG] Checking exclusion rule: %v", rule)
		match, err := rule.IsMatch(event)
		if err != nil {
			// we want to keep checking other rules here if there is an error
			continue
		}

		if match {
			return true
		}
	}

	return false
}
