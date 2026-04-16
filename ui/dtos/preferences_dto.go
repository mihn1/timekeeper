package dtos

import "github.com/mihn1/timekeeper/internal/models"

// PreferencesDto is the JSON-facing shape for user preferences exposed to the Wails frontend.
type PreferencesDto struct {
	Timezone           string `json:"timezone"`
	MinEventDurationMs int64  `json:"minEventDurationMs"`
}

func PreferencesDtoFromModel(m *models.UserPreferences) *PreferencesDto {
	return &PreferencesDto{
		Timezone:           m.Timezone,
		MinEventDurationMs: m.MinEventDurationMs,
	}
}

func (d *PreferencesDto) ToModel() *models.UserPreferences {
	return &models.UserPreferences{
		Timezone:           d.Timezone,
		MinEventDurationMs: d.MinEventDurationMs,
	}
}
