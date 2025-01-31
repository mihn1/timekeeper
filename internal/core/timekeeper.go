package core

import (
	"log"

	"github.com/mihn1/timekeeper/internal/data"
	"github.com/mihn1/timekeeper/internal/data/inmem"
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

const (
	AppName = "TimeKeeper"
)

var (
	defaultResolver CategoryResolver
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
	defaultResolver = NewDefaultCategoryResolver(t.storage.RuleStore, t.storage.CategoryStore)
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

	// Start listening for events
	go func() {
		for event := range t.eventChannel {
			t.handleEvent(&event)
		}
	}()
}

func (t *TimeKeeper) Report(date datatypes.Date) {
	log.Println("-------------TimeKeeper Report-------------")
	appAggr, _ := t.storage.AppAggregationStore.GetAppAggregationsByDate(date)
	catAggr, _ := t.storage.CategoryAggregationStore.GetCategoryAggregationsByDate(date)
	log.Printf("App Aggregation: %v\n", appAggr)
	log.Printf("Category Aggregation: %v\n", catAggr)
}

func (t *TimeKeeper) PushEvent(event models.AppSwitchEvent) {
	if t.IsEnabled() {
		t.eventChannel <- event
	}
}

func (t *TimeKeeper) handleEvent(event *models.AppSwitchEvent) {
	// TODO: gracefully handle the case when the timekeeper is disabled
	if !t.isEnabled {
		return
	}

	log.Printf("Received event: %v\n", event)
	t.aggregateEvent(event)
	t.Report(datatypes.NewDate(event.Time))
}

func (t *TimeKeeper) aggregateEvent(event *models.AppSwitchEvent) {
	if t.curAppEvent == nil {
		t.curAppEvent = event
		return
	}

	elapsedTime := int(event.Time.Sub(t.curAppEvent.Time).Seconds())

	// key := t.curAppEvent.GetEventKey()
	// aggr, ok := t.appAggregration[key]

	// if !ok {
	// 	aggr = &models.AppAggregation{
	// 		AppName: t.curAppEvent.AppName,
	// 		// SubAppName: t.curAppEvent.SubAppName,
	// 		Date: datatypes.NewDate(t.curAppEvent.Time),
	// 	}
	// 	t.appAggregration[key] = aggr
	// }

	// aggr.TimeElapsed += elapsedTime

	_, err := t.storage.AppAggregationStore.AggregateAppEvent(t.curAppEvent, elapsedTime)
	if err != nil {
		log.Printf("Error aggregating app event: %v\n", err)
	}

	t.aggregateCategory(t.curAppEvent, elapsedTime) // Call after aggregateEvent
	t.curAppEvent = event
}

func (t *TimeKeeper) aggregateCategory(event *models.AppSwitchEvent, elapsedTime int) {
	cat, err := t.getCategoryFromApp(event, defaultResolver)
	if err != nil {
		log.Printf("Error aggregating category: %v\n", err)
		return
	}

	// aggr, ok := t.categoryAggregation[cat.Id]

	// if !ok {
	// 	aggr = &models.CategoryAggregation{
	// 		CategoryId: cat.Id,
	// 		Date:       datatypes.NewDate(event.Time),
	// 	}
	// 	t.categoryAggregation[cat.Id] = aggr
	// }

	// aggr.TimeElapsed += elapsedTime

	date := datatypes.NewDate(event.Time)
	_, err = t.storage.CategoryAggregationStore.AggregateCategory(cat, date, elapsedTime)
	if err != nil {
		log.Printf("Error aggregating category: %v\n", err)
	}

	// log.Printf("Category aggregated: %v\n", aggr)
}

func (t *TimeKeeper) getCategoryFromApp(event *models.AppSwitchEvent, resovler CategoryResolver) (models.Category, error) {
	cat, err := resovler.ResolveCategory(event)
	log.Printf("Category resolved for %v: %v", event.GetEventKey(), cat.Name)
	return cat, err
}
