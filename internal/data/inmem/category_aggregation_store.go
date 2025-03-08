package inmem

import (
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/utils"
)

type CategoryAggregationStore struct {
	Aggregations map[string]*models.CategoryAggregation
}

func NewCategoryAggregationStore() *CategoryAggregationStore {
	return &CategoryAggregationStore{
		Aggregations: map[string]*models.CategoryAggregation{},
	}
}

func (store *CategoryAggregationStore) AggregateCategory(cat models.Category, date datatypes.DateOnly, elapsedTime int64) (*models.CategoryAggregation, error) {
	key := models.GetCategoryAggregationKey(cat.Id, date)
	aggr, ok := store.Aggregations[key]

	if !ok {
		aggr = &models.CategoryAggregation{
			CategoryId: cat.Id,
			Date:       date,
		}
		store.Aggregations[key] = aggr
	}

	aggr.TimeElapsed += elapsedTime
	return aggr, nil
}

func (store *CategoryAggregationStore) GetCategoryAggregation(categoryId models.CategoryId, date datatypes.DateOnly) (*models.CategoryAggregation, bool) {
	key := models.GetCategoryAggregationKey(categoryId, date)
	aggr, ok := store.Aggregations[key]
	return aggr, ok
}

func (store *CategoryAggregationStore) GetCategoryAggregations() ([]*models.CategoryAggregation, error) {
	return utils.GetMapValues(store.Aggregations), nil
}

func (store *CategoryAggregationStore) GetCategoryAggregationsByDate(date datatypes.DateOnly) ([]*models.CategoryAggregation, error) {
	var aggregations []*models.CategoryAggregation
	for _, aggregation := range store.Aggregations {
		if aggregation.Date == date {
			aggregations = append(aggregations, aggregation)
		}
	}
	return aggregations, nil
}
