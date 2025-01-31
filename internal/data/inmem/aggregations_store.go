package inmem

import (
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/utils"
)

type AppAggregationStore struct {
	Aggregations map[string]models.AppAggregation
}

func NewAppAggregationStore() *AppAggregationStore {
	return &AppAggregationStore{
		Aggregations: make(map[string]models.AppAggregation),
	}
}

func (store *AppAggregationStore) AggregateAppEvent(event *models.AppSwitchEvent, elapsedTime int) (models.AppAggregation, error) {
	key := event.GetEventKey()
	aggr, ok := store.GetAppAggregation(key, datatypes.NewDate(event.Time))

	if !ok {
		aggr = models.AppAggregation{
			AppName: event.AppName,
			// SubAppName: t.curAppEvent.SubAppName,
			Date: event.GetEventDate(),
		}
		// store.Aggregations[key] = aggr
	}

	aggr.TimeElapsed += elapsedTime
	store.Aggregations[key] = aggr
	return aggr, nil
}

func (store *AppAggregationStore) GetAppAggregation(key string, date datatypes.Date) (models.AppAggregation, bool) {
	aggr, ok := store.Aggregations[key]
	// TODO: check if the date is the same as well
	return aggr, ok
}

func (store *AppAggregationStore) GetAppAggregations() ([]models.AppAggregation, error) {
	return utils.GetMapValues(store.Aggregations), nil
}

func (store *AppAggregationStore) GetAppAggregationsByDate(date datatypes.Date) ([]models.AppAggregation, error) {
	var aggregations []models.AppAggregation
	for _, aggregation := range store.Aggregations {
		if aggregation.Date == date {
			aggregations = append(aggregations, aggregation)
		}
	}
	return aggregations, nil
}

type CategoryAggregationStore struct {
	Aggregations map[models.CategoryId]models.CategoryAggregation
}

func NewCategoryAggregationStore() *CategoryAggregationStore {
	return &CategoryAggregationStore{
		Aggregations: map[models.CategoryId]models.CategoryAggregation{},
	}
}

func (store *CategoryAggregationStore) AggregateCategory(cat models.Category, date datatypes.Date, elapsedTime int) (models.CategoryAggregation, error) {
	aggr, ok := store.Aggregations[cat.Id]

	if !ok {
		aggr = models.CategoryAggregation{
			CategoryId: cat.Id,
			Date:       date,
		}
		// store.Aggregations[cat.Id] = aggr
	}

	aggr.TimeElapsed += elapsedTime
	store.Aggregations[cat.Id] = aggr
	return aggr, nil
}

func (store *CategoryAggregationStore) GetCategoryAggregation(categoryId models.CategoryId, date datatypes.Date) (models.CategoryAggregation, bool) {
	aggr, ok := store.Aggregations[categoryId]
	return aggr, ok
}

func (store *CategoryAggregationStore) GetCategoryAggregations() ([]models.CategoryAggregation, error) {
	return utils.GetMapValues(store.Aggregations), nil
}

func (store *CategoryAggregationStore) GetCategoryAggregationsByDate(date datatypes.Date) ([]models.CategoryAggregation, error) {
	var aggregations []models.CategoryAggregation
	for _, aggregation := range store.Aggregations {
		if aggregation.Date == date {
			aggregations = append(aggregations, aggregation)
		}
	}
	return aggregations, nil
}
