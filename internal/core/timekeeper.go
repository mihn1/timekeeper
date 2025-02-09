package core

import (
	"log"

	"github.com/mihn1/timekeeper/internal/core/resolvers"
	"github.com/mihn1/timekeeper/internal/data"
	"github.com/mihn1/timekeeper/internal/data/inmem"
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

const (
	AppName = "TimeKeeper"
)

var (
	defaultResolver resolvers.CategoryResolver
)

type TimeKeeper struct {
	curAppEvent  *models.AppSwitchEvent
	storage      *data.Storage
	isEnabled    bool
	eventChannel chan models.AppSwitchEvent
}

func NewTimeKeeperInMem() *TimeKeeper {
	t := &TimeKeeper{
		storage: data.NewStorage(
			inmem.NewCategoryStore(),
			inmem.NewRuleStore(),
			inmem.NewAppAggregationStore(),
			inmem.NewCategoryAggregationStore(),
		),
		eventChannel: make(chan models.AppSwitchEvent),
	}
	return t
}

func (t *TimeKeeper) Disable() {
	t.isEnabled = false
}

func (t *TimeKeeper) IsEnabled() bool {
	return t.isEnabled
}

func (t *TimeKeeper) StartTracking() {
	t.isEnabled = true

	if !t.IsEnabled() {
		// TODO: stop the event listener
	}

	defaultResolver = resolvers.NewDefaultCategoryResolver(t.storage.RuleStore, t.storage.CategoryStore)

	// Start listening for events
	go func() {
		for event := range t.eventChannel {
			t.handleEvent(&event)
		}
	}()
}

func (t *TimeKeeper) Report(date datatypes.Date) {
	log.Printf("-------------TimeKeeper Report for %s-------------\n", date)
	appAggrs, _ := t.storage.AppAggregationStore.GetAppAggregationsByDate(date)
	catAggrs, _ := t.storage.CategoryAggregationStore.GetCategoryAggregationsByDate(date)
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

	t.aggregateEvent(event)

	// TODO: store events
	t.curAppEvent.EndTime = event.StartTime
	t.curAppEvent = event

	t.Report(datatypes.NewDate(event.StartTime))
}

func (t *TimeKeeper) aggregateEvent(event *models.AppSwitchEvent) {
	elapsedTime := event.StartTime.Sub(t.curAppEvent.StartTime).Milliseconds()

	_, err := t.storage.AppAggregationStore.AggregateAppEvent(t.curAppEvent, elapsedTime)
	if err != nil {
		log.Printf("Error aggregating app event: %v\n", err)
	}

	t.aggregateCategory(t.curAppEvent, elapsedTime) // Call after aggregateEvent
}

func (t *TimeKeeper) aggregateCategory(event *models.AppSwitchEvent, elapsedTime int64) {
	catId, err := defaultResolver.ResolveCategory(event)
	if err != nil {
		log.Printf("Error resolving category: %v\n", err)
	}

	cat, err := t.storage.CategoryStore.GetCategory(catId)
	if err != nil {
		log.Printf("Error getting category: %v\n", err)
		return
	}

	log.Printf("%v resolved for %v - %v", cat.Name, event.AppName, event.AdditionalData)
	_, err = t.storage.CategoryAggregationStore.AggregateCategory(cat, event.GetEventDate(), elapsedTime)
	if err != nil {
		log.Printf("Error aggregating category: %v\n", err)
	}
}
