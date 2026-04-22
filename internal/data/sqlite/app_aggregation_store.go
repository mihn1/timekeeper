package sqlite

import (
	"database/sql"
	"sync"

	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type AppAggregationStore struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewAppAggregationStore(db *sql.DB, mu *sync.RWMutex, tableName string) *AppAggregationStore {
	store := &AppAggregationStore{
		db:        db,
		mu:        mu,
		tableName: tableName,
	}

	_, err := store.db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			key TEXT PRIMARY KEY,
			app_name TEXT NOT NULL,
			date DATETIME NOT NULL,
			time_elapsed INTEGER NOT NULL,
			additional_data TEXT
		)`)

	if err != nil {
		panic(err)
	}

	return store
}

func (s *AppAggregationStore) AggregateAppEvent(event *models.AppSwitchEvent, elapsedTime int64) (*models.AppAggregation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := models.GetAppAggregationKey(event)
	var appAggr *models.AppAggregation = &models.AppAggregation{}
	row := s.db.QueryRow("SELECT app_name, date, time_elapsed, additional_data FROM "+s.tableName+" WHERE key = ?", key)

	err := row.Scan(&appAggr.AppName, &appAggr.Date, &appAggr.TimeElapsed, &appAggr.AdditionalData)
	if err != nil {
		if err == sql.ErrNoRows {
			appAggr = &models.AppAggregation{
				AppName: event.AppName,
				Date:    event.GetEventDate(),
			}
			_, err = s.db.Exec("INSERT INTO "+s.tableName+" (key, app_name, date, time_elapsed) VALUES (?, ?, ?, ?)", key, appAggr.AppName, appAggr.Date, appAggr.TimeElapsed)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	appAggr.TimeElapsed += elapsedTime
	_, err = s.db.Exec("UPDATE "+s.tableName+" SET time_elapsed = ? WHERE key = ?", appAggr.TimeElapsed, key)
	if err != nil {
		return nil, err
	}

	return appAggr, nil
}

func (s *AppAggregationStore) DeductAppEvent(event *models.AppSwitchEvent, elapsedTime int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := models.GetAppAggregationKey(event)
	_, err := s.db.Exec(
		"UPDATE "+s.tableName+" SET time_elapsed = MAX(0, time_elapsed - ?) WHERE key = ?",
		elapsedTime, key,
	)
	return err
}

func (s *AppAggregationStore) GetAppAggregations() ([]*models.AppAggregation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT app_name, date, time_elapsed, additional_data FROM " + s.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aggregations []*models.AppAggregation
	for rows.Next() {
		var appAggr models.AppAggregation
		err = rows.Scan(&appAggr.AppName, &appAggr.Date, &appAggr.TimeElapsed, &appAggr.AdditionalData)
		if err != nil {
			return nil, err
		}
		aggregations = append(aggregations, &appAggr)
	}

	return aggregations, nil

}

func (s *AppAggregationStore) GetAppAggregationsByDate(date datatypes.DateOnly) ([]*models.AppAggregation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT app_name, date, time_elapsed, additional_data FROM "+s.tableName+" WHERE date = ?", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aggregations []*models.AppAggregation = make([]*models.AppAggregation, 0)
	for rows.Next() {
		var appAggr models.AppAggregation
		err = rows.Scan(&appAggr.AppName, &appAggr.Date, &appAggr.TimeElapsed, &appAggr.AdditionalData)
		if err != nil {
			return nil, err
		}
		aggregations = append(aggregations, &appAggr)
	}

	return aggregations, nil
}

func (s *AppAggregationStore) ReplaceAppAggregationsForDates(dates []datatypes.DateOnly, aggrs []*models.AppAggregation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, d := range dates {
		if _, err := tx.Exec("DELETE FROM "+s.tableName+" WHERE date = ?", d); err != nil {
			return err
		}
	}

	stmt, err := tx.Prepare("INSERT INTO " + s.tableName + " (key, app_name, date, time_elapsed) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, aggr := range aggrs {
		key := aggr.AppName + "-" + aggr.Date.String()
		if _, err := stmt.Exec(key, aggr.AppName, aggr.Date, aggr.TimeElapsed); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *AppAggregationStore) GetAppAggregationsByDateRange(start, end datatypes.DateOnly) ([]*models.AppAggregation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT app_name, date, time_elapsed, additional_data FROM "+s.tableName+" WHERE date >= ? AND date <= ?", start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aggregations []*models.AppAggregation = make([]*models.AppAggregation, 0)
	for rows.Next() {
		var appAggr models.AppAggregation
		err = rows.Scan(&appAggr.AppName, &appAggr.Date, &appAggr.TimeElapsed, &appAggr.AdditionalData)
		if err != nil {
			return nil, err
		}
		aggregations = append(aggregations, &appAggr)
	}

	return aggregations, nil
}
