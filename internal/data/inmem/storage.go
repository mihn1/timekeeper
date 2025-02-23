package inmem

import (
	"github.com/mihn1/timekeeper/internal/data"
)

type InmemStorage struct {
	categoryStore            data.CategoryStore
	ruleStore                data.RuleStore
	appAggregationStore      data.AppAggregationStore
	categoryAggregationStore data.CategoryAggregationStore
}

func NewInmemStorage() *InmemStorage {
	return &InmemStorage{
		categoryStore:            NewCategoryStore(),
		ruleStore:                NewRuleStore(),
		appAggregationStore:      NewAppAggregationStore(),
		categoryAggregationStore: NewCategoryAggregationStore(),
	}
}

func (s *InmemStorage) Categories() data.CategoryStore {
	return s.categoryStore
}

func (s *InmemStorage) Rules() data.RuleStore {
	return s.ruleStore
}

func (s *InmemStorage) AppAggregations() data.AppAggregationStore {
	return s.appAggregationStore
}

func (s *InmemStorage) CategoryAggregations() data.CategoryAggregationStore {
	return s.categoryAggregationStore
}

func (s *InmemStorage) Close() error {
	return nil
}
