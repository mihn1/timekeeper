package core

import (
	"database/sql"
	"log/slog"
	"os"
	"slices"
	"sync"
	"time"

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
	done         chan struct{}
	observers    []Observer
	logger       *slog.Logger

	stateMu              sync.RWMutex
	eventLoopStarted     bool
	observersStarted     bool
	closed               bool
	closeOnce            sync.Once
	minEventDurationMs   int64

	rerun *rerunController
}

func NewTimeKeeperInMem(opts TimeKeeperOptions) *TimeKeeper {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	t := &TimeKeeper{
		Storage:      inmem.NewInmemStorage(),
		eventChannel: make(chan models.AppSwitchEvent),
		done:         make(chan struct{}),
		opts:         opts,
		logger:       opts.Logger,
	}
	return t
}

func NewTimeKeeperSqlite(opts TimeKeeperOptions) *TimeKeeper {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	if opts.StoragePath == "" {
		opts.StoragePath = "./timekeeper.db"
	}

	db, err := sql.Open("sqlite", opts.StoragePath)
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
		done:         make(chan struct{}),
		opts:         opts,
		logger:       opts.Logger,
	}
	return t
}

func (t *TimeKeeper) SetLogger(logger *slog.Logger) {
	t.logger = logger
}

// SetMinEventDurationMs sets the minimum elapsed time (ms) for an event to be
// counted in aggregations. Events shorter than this threshold are discarded as
// noise (e.g. fast app-switches). Thread-safe; takes effect on the next event.
func (t *TimeKeeper) SetMinEventDurationMs(ms int64) {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()
	if ms < 0 {
		ms = 0
	}
	t.minEventDurationMs = ms
}

func (t *TimeKeeper) Logger() *slog.Logger {
	return t.logger
}

func (t *TimeKeeper) AddObserver(o Observer) {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	if t.observers == nil {
		t.observers = make([]Observer, 0)
	}
	t.observers = append(t.observers, o)
}

func (t *TimeKeeper) Close() {
	t.closeOnce.Do(func() {
		t.logger.Info("Closing TimeKeeper...")

		t.stateMu.Lock()
		t.isEnabled = false
		t.closed = true
		close(t.done)
		observers := append([]Observer(nil), t.observers...)
		t.stateMu.Unlock()

		for _, obs := range observers {
			obs.Stop()
		}

		t.Storage.Close()
	})
}

func (t *TimeKeeper) Disable() {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	t.isEnabled = false
}

func (t *TimeKeeper) IsEnabled() bool {
	t.stateMu.RLock()
	defer t.stateMu.RUnlock()

	return t.isEnabled
}

func (t *TimeKeeper) StartTracking() {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()

	if t.closed {
		t.logger.Warn("StartTracking ignored: TimeKeeper already closed")
		return
	}

	if t.isEnabled {
		t.logger.Debug("StartTracking ignored: already enabled")
		return
	}

	t.logger.Info("Starting TimeKeeper...")

	t.isEnabled = true
	defaultResolver = resolvers.NewDefaultCategoryResolver(t.Storage.Rules(), t.Storage.Categories())

	if !t.observersStarted && t.observers != nil {
		for _, obs := range t.observers {
			go obs.Start()
		}
		t.observersStarted = true
	}

	if !t.eventLoopStarted {
		go func() {
			for {
				select {
				case event := <-t.eventChannel:
					t.handleEvent(&event)
				case <-t.done:
					return
				}
			}
		}()

		t.eventLoopStarted = true
	}

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
	t.stateMu.RLock()
	enabled := t.isEnabled && !t.closed
	t.stateMu.RUnlock()

	if !enabled {
		return
	}

	select {
	case t.eventChannel <- event:
	case <-t.done:
	}
}

func (t *TimeKeeper) handleEvent(event *models.AppSwitchEvent) {
	t.logger.Debug("Received event", "event", event)

	if t.curAppEvent == nil {
		t.curAppEvent = event
		return
	}

	// TODO: gracefully handle the case when the timekeeper is disabled
	if !t.IsEnabled() {
		return
	}

	if isSameEvent(t.curAppEvent, event) {
		t.logger.Debug("Same event detected", "event", event)
		return
	}

	t.aggregateNewEvent(event)
	t.curAppEvent.EndTime = event.StartTime

	// store the current app event (skip synthetic markers)
	if t.opts.StoreEvents && t.curAppEvent.AppName != constants.SYSTEM_PAUSED {
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

	t.stateMu.RLock()
	minMs := t.minEventDurationMs
	t.stateMu.RUnlock()

	if minMs > 0 && elapsedTime < minMs {
		t.logger.Debug("Event discarded as noise (too short)", "appName", t.curAppEvent.AppName, "elapsedMs", elapsedTime, "minMs", minMs)
		return
	}

	rules, err := t.getRulesForEvent(t.curAppEvent)
	if err != nil {
		if _, ok := err.(*EventExcludedError); ok {
			t.logger.Info("Event excluded", "by", err.Error())
			return
		}

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
	// The SYSTEM_PAUSED marker is emitted by the platform observer to close the
	// active tracking period cleanly when the machine is locked or sleeping.
	// It must never accumulate time in any aggregation.
	if event.AppName == constants.SYSTEM_PAUSED {
		return nil, NewEventExcludedError(event.AppName, nil, "system pause marker", false)
	}

	rules, err := t.Storage.Rules().GetRulesByApp(event.AppName)
	if err != nil {
		return nil, err
	}

	logInfoWhenExcluded := true // TODO: make this configurable

	if excluded, rule := isEventExcluded(event, rules); excluded {
		return nil, NewEventExcludedError(event.AppName, rule, "Event matched exlusion rule", logInfoWhenExcluded)
	}

	// Get rules that are applied to all apps
	globalRules, err := t.Storage.Rules().GetRulesByApp(constants.ALL_APPS)
	if err == nil {
		if excluded, rule := isEventExcluded(event, globalRules); excluded {
			return nil, NewEventExcludedError(event.AppName, rule, "Event matched global exlusion rule", logInfoWhenExcluded)
		}
	}

	rules = append(rules, globalRules...)
	slices.SortStableFunc(rules, models.CmpRules)

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

	// If the events are within 60 seconds of each other, consider them the same event.
	delta := e1.StartTime.Sub(e2.StartTime)
	if delta < 0 {
		delta = -delta
	}

	return delta <= 60*time.Second
}

func isEventExcluded(event *models.AppSwitchEvent, rules []*models.CategoryRule) (excluded bool, rule *models.CategoryRule) {
	for _, rule := range rules {
		if !rule.IsExclusion {
			continue
		}

		match, err := rule.IsMatch(event)
		if err != nil {
			// we want to keep checking other rules here if there is an error
			continue
		}

		if match {
			return true, rule
		}
	}

	return false, nil
}
