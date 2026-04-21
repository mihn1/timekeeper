package interfaces

import (
	"time"

	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type EventStore interface {
	GetEventsByDate(date datatypes.DateOnly) ([]*models.AppSwitchEvent, error)
	// GetEventsByTimeRange returns events whose start_time falls in [start, end).
	// Used for timezone-aware day queries.
	GetEventsByTimeRange(start, end time.Time) ([]*models.AppSwitchEvent, error)
	GetEvent(id models.EventId) (*models.AppSwitchEvent, error)
	AddEvent(event *models.AppSwitchEvent) error
	DeleteEvent(id models.EventId) error
}
