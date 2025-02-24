package core

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mihn1/timekeeper/internal/core/resolvers"
	"github.com/mihn1/timekeeper/internal/data"
	"github.com/mihn1/timekeeper/internal/data/inmem"
	"github.com/mihn1/timekeeper/internal/data/sqlite"
	"github.com/mihn1/timekeeper/internal/datatypes"
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
}

type TimeKeeper struct {
	curAppEvent  *models.AppSwitchEvent
	opts         TimeKeeperOptions
	storage      data.Storage
	isEnabled    bool
	eventChannel chan models.AppSwitchEvent
}

func NewTimeKeeperInMem(opts TimeKeeperOptions) *TimeKeeper {
	t := &TimeKeeper{
		storage:      inmem.NewInmemStorage(),
		eventChannel: make(chan models.AppSwitchEvent),
		opts:         opts,
	}
	return t
}

func NewTimeKeeperSqlite(opts TimeKeeperOptions) *TimeKeeper {
	if opts.StoragePath == "" {
		log.Println("No storage path provided. Using default path.")
		opts.StoragePath = "./timekeeper.db"
	}

	db, err := sql.Open("sqlite3", opts.StoragePath)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	_, err = db.Exec("PRAGMA busy_timeout = 5000;")
	if err != nil {
		log.Fatalf("Error setting busy_timeout: %v\n", err)
	}

	t := &TimeKeeper{
		storage:      sqlite.NewSqliteStorage(db),
		eventChannel: make(chan models.AppSwitchEvent),
		opts:         opts,
	}
	return t
}

func (t *TimeKeeper) Close() {
	log.Println("Closing TimeKeeper...")
	t.Disable()
	t.storage.Close()
}

func (t *TimeKeeper) Disable() {
	t.isEnabled = false
}

func (t *TimeKeeper) IsEnabled() bool {
	return t.isEnabled
}

func (t *TimeKeeper) StartTracking() {
	t.isEnabled = true
	defaultResolver = resolvers.NewDefaultCategoryResolver(t.storage.Rules(), t.storage.Categories())

	// Start listening for events
	go func() {
		for event := range t.eventChannel {
			t.handleEvent(&event)
		}
	}()
}

func (t *TimeKeeper) Report(date datatypes.DateOnly) {
	log.Printf("-------------TimeKeeper Report for %s-------------\n", date)
	appAggrs, _ := t.storage.AppAggregations().GetAppAggregationsByDate(date)
	catAggrs, _ := t.storage.CategoryAggregations().GetCategoryAggregationsByDate(date)
	log.Printf("App Aggregation: %v\n", appAggrs)
	log.Printf("Category Aggregation: %v\n", catAggrs)
	log.Println("-------------------------------------------------------------------------------------------------")
}

func (t *TimeKeeper) PushEvent(event models.AppSwitchEvent) {
	if t.IsEnabled() {
		t.eventChannel <- event
	}
}

func (t *TimeKeeper) handleEvent(event *models.AppSwitchEvent) {
	log.Printf("Received event: %v\n", event)

	if t.curAppEvent == nil {
		t.curAppEvent = event
		return
	}

	// TODO: gracefully handle the case when the timekeeper is disabled
	if !t.isEnabled {
		return
	}

	if isSameEvent(t.curAppEvent, event) {
		log.Printf("Same event detected: %v\n", event)
		return
	}

	t.aggregateNewEvent(event)
	t.curAppEvent.EndTime = event.StartTime

	// store the current app event
	if t.opts.StoreEvents {
		err := t.storage.Events().AddEvent(t.curAppEvent)
		if err != nil {
			log.Printf("Error storing event: %v\n", err)
		}
	}

	t.curAppEvent = event
	t.Report(datatypes.NewDateOnly(event.StartTime))
}

func (t *TimeKeeper) aggregateNewEvent(event *models.AppSwitchEvent) {
	elapsedTime := event.StartTime.Sub(t.curAppEvent.StartTime).Milliseconds()

	_, err := t.storage.AppAggregations().AggregateAppEvent(t.curAppEvent, elapsedTime)
	if err != nil {
		log.Printf("Error aggregating app event for %s: %v\n", event.AppName, err)
		return
	}

	catId, err := t.aggregateCategory(t.curAppEvent, elapsedTime) // Call after aggregateEvent
	if err != nil {
		log.Printf("Error aggregating category for %s: %v\n", event.AppName, err)
		return
	}

	t.curAppEvent.CategoryId = catId
}

func (t *TimeKeeper) aggregateCategory(event *models.AppSwitchEvent, elapsedTime int64) (models.CategoryId, error) {
	catId, err := defaultResolver.ResolveCategory(event)
	if err != nil {
		return catId, err
	}

	cat, err := t.storage.Categories().GetCategory(catId)
	if err != nil {
		log.Printf("Error getting category: %v\n", err)
		return catId, err
	}

	log.Printf("%v resolved for %v - %v", cat.Name, event.AppName, event.AdditionalData)
	_, err = t.storage.CategoryAggregations().AggregateCategory(cat, event.GetEventDate(), elapsedTime)
	if err != nil {
		log.Printf("Error saving category aggregation: %v\n", err)
		return catId, err
	}

	return catId, nil
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
