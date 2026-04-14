package inmem

import (
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type PreferencesStore struct {
	mu    sync.RWMutex
	prefs *models.UserPreferences
}

func NewPreferencesStore() *PreferencesStore {
	return &PreferencesStore{prefs: models.DefaultPreferences()}
}

func (s *PreferencesStore) GetPreferences() (*models.UserPreferences, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := *s.prefs
	return &cp, nil
}

func (s *PreferencesStore) SavePreferences(prefs *models.UserPreferences) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *prefs
	s.prefs = &cp
	return nil
}
