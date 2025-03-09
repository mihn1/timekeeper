package sqlite

import (
	"database/sql"
	"sync"

	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryAggregations struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewCategoryAggregationStore(db *sql.DB, mu *sync.RWMutex, tableName string) *CategoryAggregations {
	s := &CategoryAggregations{
		db:        db,
		mu:        mu,
		tableName: tableName,
	}

	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			Key TEXT PRIMARY KEY,
			category_id TEXT NOT NULL,
			date DATETIME NOT NULL,
			time_elapsed INTEGER NOT NULL
		)`)

	if err != nil {
		panic(err)
	}

	return s
}

func (s *CategoryAggregations) AggregateCategory(cat models.Category, date datatypes.DateOnly, elapsedTime int64) (*models.CategoryAggregation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := models.GetCategoryAggregationKey(cat.Id, date)
	var aggr *models.CategoryAggregation = &models.CategoryAggregation{}
	row := s.db.QueryRow("SELECT category_id, date, time_elapsed FROM "+s.tableName+" WHERE key = ?", key)
	err := row.Scan(&aggr.CategoryId, &aggr.Date, &aggr.TimeElapsed)
	if err != nil {
		if err == sql.ErrNoRows {
			aggr = &models.CategoryAggregation{
				CategoryId: cat.Id,
				Date:       date,
			}
			_, err = s.db.Exec("INSERT INTO "+s.tableName+" (key, category_id, date, time_elapsed) VALUES (?, ?, ?, ?)", key, aggr.CategoryId, aggr.Date, aggr.TimeElapsed)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	aggr.TimeElapsed += elapsedTime
	_, err = s.db.Exec("UPDATE "+s.tableName+" SET time_elapsed = ? WHERE key = ?", aggr.TimeElapsed, key)
	if err != nil {
		return nil, err
	}

	return aggr, nil
}

func (s *CategoryAggregations) GetCategoryAggregation(categoryId models.CategoryId, date datatypes.DateOnly) (*models.CategoryAggregation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := models.GetCategoryAggregationKey(categoryId, date)
	row := s.db.QueryRow("SELECT category_id, date, time_elapsed FROM "+s.tableName+" WHERE key = ?", key)
	aggregation := &models.CategoryAggregation{}
	err := row.Scan(&aggregation.CategoryId, &aggregation.Date, &aggregation.TimeElapsed)
	if err != nil {
		return nil, false
	}
	return aggregation, true
}

func (s *CategoryAggregations) GetCategoryAggregations() ([]*models.CategoryAggregation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT category_id, date, time_elapsed FROM " + s.tableName)
	if err != nil {
		return nil, err
	}

	var aggregations []*models.CategoryAggregation
	for rows.Next() {
		aggregation := &models.CategoryAggregation{}
		err = rows.Scan(&aggregation.CategoryId, &aggregation.Date, &aggregation.TimeElapsed)
		if err != nil {
			return nil, err
		}
		aggregations = append(aggregations, aggregation)
	}

	return aggregations, nil
}

func (s *CategoryAggregations) GetCategoryAggregationsByDate(date datatypes.DateOnly) ([]*models.CategoryAggregation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT category_id, date, time_elapsed FROM "+s.tableName+" WHERE date = ?", date)
	if err != nil {
		return nil, err
	}

	var aggregations []*models.CategoryAggregation = make([]*models.CategoryAggregation, 0)
	for rows.Next() {
		aggregation := &models.CategoryAggregation{}
		err = rows.Scan(&aggregation.CategoryId, &aggregation.Date, &aggregation.TimeElapsed)
		if err != nil {
			return nil, err
		}
		aggregations = append(aggregations, aggregation)
	}

	return aggregations, nil
}
