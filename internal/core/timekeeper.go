package core

import (
	"log"
	"time"

	"github.com/mihn1/timekeeper/internal/data"
	"golang.org/x/exp/rand"
)

const (
	AppName = "TimeKeeper"
)

type TimeKeeper struct {
	curAppEvent         *data.AppSwitchEvent
	appAggregration     map[string]*data.AppAggregation
	categoryStore       data.CategoryStore
	categoryAggregation map[data.CategoryId]*data.CategoryAggregation
	isEnabled           bool
	eventChannel        chan data.AppSwitchEvent
}

func NewTimeKeeperInMem() *TimeKeeper {
	return &TimeKeeper{
		categoryStore:       data.NewCategoryStore_Memory_Impl(),
		appAggregration:     make(map[string]*data.AppAggregation),
		categoryAggregation: make(map[data.CategoryId]*data.CategoryAggregation),
		eventChannel:        make(chan data.AppSwitchEvent),
	}
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
			t.handleEvent(event)
		}
	}()
}

func (t *TimeKeeper) PushEvent(event data.AppSwitchEvent) {
	if t.IsEnabled() {
		t.eventChannel <- event
	}
}

func (t *TimeKeeper) handleEvent(event data.AppSwitchEvent) {
	// TODO: gracefully handle the case when the timekeeper is disabled
	if !t.isEnabled {
		return
	}

	log.Printf("Received event: %v\n", event)
	elapsedTime := t.aggregateEvent(event)
	log.Printf("App Aggregated: %v\n", t.appAggregration)
	t.aggregateCategory(elapsedTime) // Call after aggregateEvent
	log.Printf("Category Aggregated: %v\n", t.categoryAggregation)
}

func (t *TimeKeeper) aggregateEvent(event data.AppSwitchEvent) int {
	if t.curAppEvent == nil {
		t.curAppEvent = &event
		return 0
	}

	key := t.curAppEvent.GetEventKey()
	aggr, ok := t.appAggregration[key]

	if !ok {
		aggr = &data.AppAggregation{
			AppName:    t.curAppEvent.AppName,
			SubAppName: t.curAppEvent.SubAppName,
		}
		t.appAggregration[key] = aggr
	}

	elapsedTime := int(event.Time.Sub(t.curAppEvent.Time).Seconds())
	aggr.TimeElapsed += elapsedTime
	t.curAppEvent = &event
	return elapsedTime
}

func (t *TimeKeeper) aggregateCategory(elapsedTime int) {
	cat, err := t.getCategoryFromApp(t.curAppEvent)
	if err != nil {
		log.Printf("Error aggregating category: %v\n", err)
		return
	}

	aggr, ok := t.categoryAggregation[cat.Id]

	if !ok {
		aggr = &data.CategoryAggregation{
			CategoryId: cat.Id,
		}
		t.categoryAggregation[cat.Id] = aggr
	}

	aggr.TimeElapsed += elapsedTime
}

func (t *TimeKeeper) getCategoryFromApp(event *data.AppSwitchEvent) (data.Category, error) {
	// This is a dummy implementation
	rand.Seed(uint64(time.Now().UnixNano()))
	id := data.CategoryId(rand.Intn(3))
	cat, err := t.categoryStore.GetCategory(id)
	return cat, err
}
