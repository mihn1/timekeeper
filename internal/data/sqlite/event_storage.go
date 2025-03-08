package sqlite

import (
	"database/sql"
	"encoding/json"
	"sync"

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

func (s *EventStore) GetEventsByDate(date datatypes.DateOnly) ([]*models.AppSwitchEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT id, category_id, app_name, additional_data, start_time, end_time FROM "+s.tableName+" WHERE date = ?", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.AppSwitchEvent
	for rows.Next() {
		var event models.AppSwitchEvent
		var rawData string
		err = rows.Scan(&event.Id, &event.CategoryId, &event.AppName, &rawData, &event.StartTime, &event.EndTime)
		if err != nil {
			return nil, err
		}

		if rawData != "{}" {
			err = json.Unmarshal([]byte(rawData), &event.AdditionalData)
			if err != nil {
				return nil, err
			}
		}
	}

	return events, nil
}
