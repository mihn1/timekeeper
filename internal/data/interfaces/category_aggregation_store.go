package interfaces

import (
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryAggregationStore interface {
	AggregateCategory(cat models.Category, date datatypes.DateOnly, elapsedTime int64) (*models.CategoryAggregation, error)
	GetCategoryAggregation(categoryId models.CategoryId, date datatypes.DateOnly) (*models.CategoryAggregation, bool)
	GetCategoryAggregations() ([]*models.CategoryAggregation, error)
	GetCategoryAggregationsByDate(date datatypes.DateOnly) ([]*models.CategoryAggregation, error)
}
