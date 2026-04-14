package interfaces

import "github.com/mihn1/timekeeper/internal/models"

type PreferencesStore interface {
	GetPreferences() (*models.UserPreferences, error)
	SavePreferences(*models.UserPreferences) error
}
