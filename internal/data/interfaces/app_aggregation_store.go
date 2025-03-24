package interfaces

import (
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type AppAggregationStore interface {
	AggregateAppEvent(event *models.AppSwitchEvent, elapsedTime int64) (*models.AppAggregation, error)
	GetAppAggregations() ([]*models.AppAggregation, error)
	GetAppAggregationsByDate(date datatypes.DateOnly) ([]*models.AppAggregation, error)
}
