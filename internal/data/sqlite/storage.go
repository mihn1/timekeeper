package sqlite

import (
	"database/sql"
	"sync"

	"github.com/mihn1/timekeeper/internal/data"
)

type SqliteStorage struct {
	db                       *sql.DB
	categoryStore            data.CategoryStore
	ruleStore                data.RuleStore
	appAggregationStore      data.AppAggregationStore
	categoryAggregationStore data.CategoryAggregationStore
	eventStore               data.EventStore
}

func NewSqliteStorage(db *sql.DB) *SqliteStorage {
	mu := &sync.RWMutex{}
	return &SqliteStorage{
		db:                       db,
		categoryStore:            NewCategoryStore(db, mu, "categories"),
		ruleStore:                NewRuleStore(db, mu, "rules"),
		appAggregationStore:      NewAppAggregationStore(db, mu, "app_aggregations"),
		categoryAggregationStore: NewCategoryAggregationStore(db, mu, "category_aggregations"),
		eventStore:               NewEventStore(db, mu, "events"),
	}
}

func (s *SqliteStorage) Categories() data.CategoryStore {
	return s.categoryStore
}

func (s *SqliteStorage) Rules() data.RuleStore {
	return s.ruleStore
}

func (s *SqliteStorage) AppAggregations() data.AppAggregationStore {
	return s.appAggregationStore
}

func (s *SqliteStorage) CategoryAggregations() data.CategoryAggregationStore {
	return s.categoryAggregationStore
}

func (s *SqliteStorage) Events() data.EventStore {
	return s.eventStore
}

func (s *SqliteStorage) Close() error {
	return s.db.Close()
}
