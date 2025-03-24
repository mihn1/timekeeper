package interfaces

import (
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type EventStore interface {
	GetEventsByDate(date datatypes.DateOnly) ([]*models.AppSwitchEvent, error)
	AddEvent(event *models.AppSwitchEvent) error
}
