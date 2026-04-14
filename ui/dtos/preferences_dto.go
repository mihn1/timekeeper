package dtos

import "github.com/mihn1/timekeeper/internal/models"

// PreferencesDto is the JSON-facing shape for user preferences exposed to the Wails frontend.
type PreferencesDto struct {
	Timezone string `json:"timezone"`
}

func PreferencesDtoFromModel(m *models.UserPreferences) *PreferencesDto {
	return &PreferencesDto{Timezone: m.Timezone}
}

func (d *PreferencesDto) ToModel() *models.UserPreferences {
	return &models.UserPreferences{Timezone: d.Timezone}
}
