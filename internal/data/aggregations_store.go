package data

import (
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type AppAggregationStore interface {
	AggregateAppEvent(event *models.AppSwitchEvent, elapsedTime int) (models.AppAggregation, error)
	GetAppAggregation(key string, date datatypes.Date) (models.AppAggregation, bool)
	GetAppAggregations() ([]models.AppAggregation, error)
	GetAppAggregationsByDate(date datatypes.Date) ([]models.AppAggregation, error)
}

type CategoryAggregationStore interface {
	AggregateCategory(cat models.Category, date datatypes.Date, elapsedTime int) (models.CategoryAggregation, error)
	GetCategoryAggregation(categoryId models.CategoryId, date datatypes.Date) (models.CategoryAggregation, bool)
	GetCategoryAggregations() ([]models.CategoryAggregation, error)
	GetCategoryAggregationsByDate(date datatypes.Date) ([]models.CategoryAggregation, error)
}
