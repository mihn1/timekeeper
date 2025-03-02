package inmem

import "github.com/mihn1/timekeeper/internal/data/interfaces"

type InmemStorage struct {
	categoryStore            interfaces.CategoryStore
	ruleStore                interfaces.RuleStore
	appAggregationStore      interfaces.AppAggregationStore
	categoryAggregationStore interfaces.CategoryAggregationStore
	eventStore               interfaces.EventStore
}

func NewInmemStorage() *InmemStorage {
	return &InmemStorage{
		categoryStore:            NewCategoryStore(),
		ruleStore:                NewRuleStore(),
		appAggregationStore:      NewAppAggregationStore(),
		categoryAggregationStore: NewCategoryAggregationStore(),
		eventStore:               NewEventStore(),
	}
}

func (s *InmemStorage) Categories() interfaces.CategoryStore {
	return s.categoryStore
}

func (s *InmemStorage) Rules() interfaces.RuleStore {
	return s.ruleStore
}

func (s *InmemStorage) AppAggregations() interfaces.AppAggregationStore {
	return s.appAggregationStore
}

func (s *InmemStorage) CategoryAggregations() interfaces.CategoryAggregationStore {
	return s.categoryAggregationStore
}

func (s *InmemStorage) Events() interfaces.EventStore {
	return s.eventStore
}

func (s *InmemStorage) Close() error {
	return nil
}
