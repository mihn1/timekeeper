package interfaces

import (
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryAggregationStore interface {
	AggregateCategory(cat *models.Category, date datatypes.DateOnly, elapsedTime int64) (*models.CategoryAggregation, error)
	DeductCategory(categoryId models.CategoryId, date datatypes.DateOnly, elapsedTime int64) error
	GetCategoryAggregation(categoryId models.CategoryId, date datatypes.DateOnly) (*models.CategoryAggregation, bool)
	GetCategoryAggregations() ([]*models.CategoryAggregation, error)
	GetCategoryAggregationsByDate(date datatypes.DateOnly) ([]*models.CategoryAggregation, error)
	GetCategoryAggregationsByDateRange(start, end datatypes.DateOnly) ([]*models.CategoryAggregation, error)
	// ReplaceCategoryAggregationsForDates deletes all existing category aggregations for the
	// given dates and inserts the provided aggregations atomically (best-effort per-store).
	ReplaceCategoryAggregationsForDates(dates []datatypes.DateOnly, aggrs []*models.CategoryAggregation) error
}
