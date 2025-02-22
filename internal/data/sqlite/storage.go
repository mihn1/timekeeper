package sqlite

import (
	"database/sql"

	"github.com/mihn1/timekeeper/internal/data"
)

type SqliteStorage struct {
	categoryStore            data.CategoryStore
	ruleStore                data.RuleStore
	appAggregationStore      data.AppAggregationStore
	categoryAggregationStore data.CategoryAggregationStore
}

func NewSqliteStorage(db *sql.DB) *SqliteStorage {
	return &SqliteStorage{
		categoryStore:            NewCategoryStore(db, "categories"),
		ruleStore:                NewRuleStore(db, "rules"),
		appAggregationStore:      NewAppAggregationStore(db, "app_aggregations"),
		categoryAggregationStore: NewCategoryAggregationStore(db, "category_aggregations"),
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
