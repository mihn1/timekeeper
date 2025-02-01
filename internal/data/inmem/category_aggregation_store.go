package inmem

import (
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/utils"
)

type CategoryAggregationStore struct {
	Aggregations map[models.CategoryId]*models.CategoryAggregation
}

func NewCategoryAggregationStore() *CategoryAggregationStore {
	return &CategoryAggregationStore{
		Aggregations: map[models.CategoryId]*models.CategoryAggregation{},
	}
}

func (store *CategoryAggregationStore) AggregateCategory(cat models.Category, date datatypes.Date, elapsedTime int) (*models.CategoryAggregation, error) {
	aggr, ok := store.Aggregations[cat.Id]

	if !ok {
		aggr = &models.CategoryAggregation{
			CategoryId: cat.Id,
			Date:       date,
		}
		store.Aggregations[cat.Id] = aggr
	}

	aggr.TimeElapsed += elapsedTime
	return aggr, nil
}

func (store *CategoryAggregationStore) GetCategoryAggregation(categoryId models.CategoryId, date datatypes.Date) (*models.CategoryAggregation, bool) {
	aggr, ok := store.Aggregations[categoryId]
	return aggr, ok
}

func (store *CategoryAggregationStore) GetCategoryAggregations() ([]*models.CategoryAggregation, error) {
	return utils.GetMapValues(store.Aggregations), nil
}

func (store *CategoryAggregationStore) GetCategoryAggregationsByDate(date datatypes.Date) ([]*models.CategoryAggregation, error) {
	var aggregations []*models.CategoryAggregation
	for _, aggregation := range store.Aggregations {
		if aggregation.Date == date {
			aggregations = append(aggregations, aggregation)
		}
	}
	return aggregations, nil
}
