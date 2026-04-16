package sqlite

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type PreferencesStore struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewPreferencesStore(db *sql.DB, mu *sync.RWMutex, tableName string) *PreferencesStore {
	store := &PreferencesStore{db: db, mu: mu, tableName: tableName}
	_, err := store.db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			key   TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`)
	if err != nil {
		panic(err)
	}
	return store
}

func (s *PreferencesStore) GetPreferences() (*models.UserPreferences, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prefs := models.DefaultPreferences()
	rows, err := s.db.Query("SELECT key, value FROM " + s.tableName)
	if err != nil {
		return prefs, nil // return defaults on read error
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}
		switch key {
		case "timezone":
			if value != "" {
				prefs.Timezone = value
			}
		case "min_event_duration_ms":
			var ms int64
			if _, err := fmt.Sscan(value, &ms); err == nil && ms >= 0 {
				prefs.MinEventDurationMs = ms
			}
		}
	}
	return prefs, nil
}

func (s *PreferencesStore) SavePreferences(prefs *models.UserPreferences) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	upsert := "INSERT INTO " + s.tableName + " (key, value) VALUES (?, ?)" +
		" ON CONFLICT(key) DO UPDATE SET value = excluded.value"

	if _, err := s.db.Exec(upsert, "timezone", prefs.Timezone); err != nil {
		return err
	}
	if _, err := s.db.Exec(upsert, "min_event_duration_ms", fmt.Sprintf("%d", prefs.MinEventDurationMs)); err != nil {
		return err
	}
	return nil
}
