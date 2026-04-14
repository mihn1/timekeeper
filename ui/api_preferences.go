package main

import (
	"fmt"

	"github.com/mihn1/timekeeper/ui/dtos"
)

// GetPreferences returns the current user preferences from storage.
func (a *App) GetPreferences() (*dtos.PreferencesDto, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}
	prefs, err := a.timekeeper.Storage.Preferences().GetPreferences()
	if err != nil {
		return nil, err
	}
	return dtos.PreferencesDtoFromModel(prefs), nil
}

// SavePreferences persists the given preferences to storage and updates the
// in-memory cache so subsequent data queries immediately use the new timezone.
func (a *App) SavePreferences(dto *dtos.PreferencesDto) error {
	if a.timekeeper == nil {
		return fmt.Errorf("timekeeper is not initialized")
	}
	if dto == nil {
		return fmt.Errorf("preferences payload is required")
	}

	prefs := dto.ToModel()
	if err := a.timekeeper.Storage.Preferences().SavePreferences(prefs); err != nil {
		return fmt.Errorf("failed to save preferences: %w", err)
	}

	a.prefsMu.Lock()
	a.prefs = prefs
	a.prefsMu.Unlock()

	a.logger.Info("Preferences saved", "timezone", prefs.Timezone)
	return nil
}
