package sqlite

import (
	"database/sql"
	"sync"

	"github.com/mihn1/timekeeper/internal/data/interfaces"
)

type SqliteStorage struct {
	db                       *sql.DB
	categoryStore            interfaces.CategoryStore
	ruleStore                interfaces.RuleStore
	appAggregationStore      interfaces.AppAggregationStore
	categoryAggregationStore interfaces.CategoryAggregationStore
	eventStore               interfaces.EventStore
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

func (s *SqliteStorage) Categories() interfaces.CategoryStore {
	return s.categoryStore
}

func (s *SqliteStorage) Rules() interfaces.RuleStore {
	return s.ruleStore
}

func (s *SqliteStorage) AppAggregations() interfaces.AppAggregationStore {
	return s.appAggregationStore
}

func (s *SqliteStorage) CategoryAggregations() interfaces.CategoryAggregationStore {
	return s.categoryAggregationStore
}

func (s *SqliteStorage) Events() interfaces.EventStore {
	return s.eventStore
}

func (s *SqliteStorage) Close() error {
	return s.db.Close()
}
