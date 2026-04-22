package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type EventStore struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewEventStore(db *sql.DB, mu *sync.RWMutex, tableName string) *EventStore {
	store := &EventStore{
		db:        db,
		mu:        mu,
		tableName: tableName,
	}

	_, err := store.db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id TEXT,
			app_name TEXT NOT NULL,
			additional_data TEXT,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			date DATETIME NOT NULL
		)`)

	if err != nil {
		panic(err)
	}

	return store
}

func (s *EventStore) AddEvent(event *models.AppSwitchEvent) error {
	// dump additional data to json
	var rawData string
	if event.AdditionalData == nil {
		rawData = "{}"
	} else {
		bytes, err := json.Marshal(event.AdditionalData)
		if err != nil {
			rawData = "{}"
		} else {
			rawData = string(bytes)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("INSERT INTO "+s.tableName+" (category_id, app_name, additional_data, start_time, end_time, date) VALUES (?, ?, ?, ?, ?, ?)", event.CategoryId, event.AppName, rawData, event.StartTime, event.EndTime, datatypes.NewDateOnly(event.StartTime))
	return err
}

func (s *EventStore) GetEventsByTimeRange(start, end time.Time) ([]*models.AppSwitchEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		"SELECT id, category_id, app_name, additional_data, start_time, end_time FROM "+s.tableName+
			" WHERE start_time >= ? AND start_time < ? ORDER BY start_time",
		start.UTC(), end.UTC(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEventRows(rows)
}

func (s *EventStore) GetEventsByDate(date datatypes.DateOnly) ([]*models.AppSwitchEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT id, category_id, app_name, additional_data, start_time, end_time FROM "+s.tableName+" WHERE date = ?", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEventRows(rows)
}

func (s *EventStore) GetEvent(id models.EventId) (*models.AppSwitchEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		"SELECT id, category_id, app_name, additional_data, start_time, end_time FROM "+s.tableName+" WHERE id = ?", id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events, err := scanEventRows(rows)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("event %d not found", id)
	}
	return events[0], nil
}

func (s *EventStore) UpdateEventCategory(id models.EventId, categoryId models.CategoryId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("UPDATE "+s.tableName+" SET category_id = ? WHERE id = ?", categoryId, id)
	return err
}

func (s *EventStore) DeleteEvent(id models.EventId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM "+s.tableName+" WHERE id = ?", id)
	return err
}

func scanEventRows(rows *sql.Rows) ([]*models.AppSwitchEvent, error) {
	var events []*models.AppSwitchEvent
	for rows.Next() {
		var event models.AppSwitchEvent
		var rawData string
		err := rows.Scan(&event.Id, &event.CategoryId, &event.AppName, &rawData, &event.StartTime, &event.EndTime)
		if err != nil {
			return nil, err
		}

		if rawData != "{}" {
			err = json.Unmarshal([]byte(rawData), &event.AdditionalData)
			if err != nil {
				return nil, err
			}
		}

		events = append(events, &event)
	}

	return events, nil
}
