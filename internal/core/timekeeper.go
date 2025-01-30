package core

import (
	"log"

	"github.com/mihn1/timekeeper/internal/data"
	"github.com/mihn1/timekeeper/internal/models"
)

const (
	AppName = "TimeKeeper"
)

var (
	defaultResolver CategoryResolver
)

type TimeKeeper struct {
	curAppEvent         *models.AppSwitchEvent
	appAggregration     map[string]*models.AppAggregation
	categoryAggregation map[models.CategoryId]*models.CategoryAggregation
	categoryStore       data.CategoryStore
	ruleStore           data.RuleStore
	isEnabled           bool
	eventChannel        chan models.AppSwitchEvent
}

func NewTimeKeeperInMem() *TimeKeeper {
	t := &TimeKeeper{
		categoryStore:       data.NewCategoryStore_Memory_Impl(),
		ruleStore:           data.NewRuleStore_InMemory_Impl(),
		appAggregration:     make(map[string]*models.AppAggregation),
		categoryAggregation: make(map[models.CategoryId]*models.CategoryAggregation),
		eventChannel:        make(chan models.AppSwitchEvent),
	}
	defaultResolver = NewDefaultCategoryResolver(t.ruleStore, t.categoryStore)
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
	elapsedTime := t.aggregateEvent(event)
	log.Printf("App Aggregated: %v\n", t.appAggregration)
	t.aggregateCategory(event, elapsedTime) // Call after aggregateEvent
	log.Printf("Category Aggregated: %v\n", t.categoryAggregation)
}

func (t *TimeKeeper) aggregateEvent(event *models.AppSwitchEvent) int {
	if t.curAppEvent == nil {
		t.curAppEvent = event
		return 0
	}

	key := t.curAppEvent.GetEventKey()
	aggr, ok := t.appAggregration[key]

	if !ok {
		aggr = &models.AppAggregation{
			AppName:    t.curAppEvent.AppName,
			SubAppName: t.curAppEvent.SubAppName,
		}
		t.appAggregration[key] = aggr
	}

	elapsedTime := int(event.Time.Sub(t.curAppEvent.Time).Seconds())
	aggr.TimeElapsed += elapsedTime
	t.curAppEvent = event
	return elapsedTime
}

func (t *TimeKeeper) aggregateCategory(event *models.AppSwitchEvent, elapsedTime int) {
	cat, err := t.getCategoryFromApp(event, defaultResolver)
	if err != nil {
		log.Printf("Error aggregating category: %v\n", err)
		return
	}

	aggr, ok := t.categoryAggregation[cat.Id]

	if !ok {
		aggr = &models.CategoryAggregation{
			CategoryId: cat.Id,
		}
		t.categoryAggregation[cat.Id] = aggr
	}

	aggr.TimeElapsed += elapsedTime
}

func (t *TimeKeeper) getCategoryFromApp(event *models.AppSwitchEvent, resovler CategoryResolver) (models.Category, error) {
	cat, err := resovler.ResolveCategory(event)
	log.Printf("Category resolved for %v: %v", event.GetEventKey(), cat.Name)
	return cat, err
}
