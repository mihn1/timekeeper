package sqlite

import (
	"database/sql"
	"sync"

	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type AppAggregationStore struct {
	db        *sql.DB
	tableName string
	mu        sync.Mutex // Add a mutex to protect critical sections
}

func NewAppAggregationStore(db *sql.DB, tableName string) *AppAggregationStore {
	store := &AppAggregationStore{
		db:        db,
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
	s.mu.Lock()         // Lock the mutex
	defer s.mu.Unlock() // Unlock the mutex when the function returns

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
			_, err = s.db.Exec("INSERT INTO "+s.tableName+" (key, app_name, date, time_elapsed) VALUES (?, ?, ?, ?)", key, appAggr.AppName, appAggr.Date.DateTime(), appAggr.TimeElapsed)
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

func (store *AppAggregationStore) GetAppAggregations() ([]*models.AppAggregation, error) {
	store.mu.Lock()         // Lock the mutex
	defer store.mu.Unlock() // Unlock the mutex when the function returns

	rows, err := store.db.Query("SELECT app_name, date, time_elapsed, additional_data FROM " + store.tableName)
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

func (store *AppAggregationStore) GetAppAggregationsByDate(date datatypes.Date) ([]*models.AppAggregation, error) {
	store.mu.Lock()         // Lock the mutex
	defer store.mu.Unlock() // Unlock the mutex when the function returns

	rows, err := store.db.Query("SELECT app_name, date, time_elapsed, additional_data FROM "+store.tableName+" WHERE date = ?", date.DateTime())
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
