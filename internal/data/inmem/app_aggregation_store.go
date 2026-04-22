package inmem

import (
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/utils"
)

type AppAggregationStore struct {
	Aggregations map[string]*models.AppAggregation
}

func NewAppAggregationStore() *AppAggregationStore {
	return &AppAggregationStore{
		Aggregations: make(map[string]*models.AppAggregation),
	}
}

func (store *AppAggregationStore) AggregateAppEvent(event *models.AppSwitchEvent, elapsedTime int64) (*models.AppAggregation, error) {
	key := models.GetAppAggregationKey(event)
	aggr, ok := store.Aggregations[key]

	if !ok {
		aggr = &models.AppAggregation{
			AppName: event.AppName,
			Date:    event.GetEventDate(),
		}
		store.Aggregations[key] = aggr
	}

	aggr.TimeElapsed += elapsedTime
	return aggr, nil
}

func (store *AppAggregationStore) DeductAppEvent(event *models.AppSwitchEvent, elapsedTime int64) error {
	key := models.GetAppAggregationKey(event)
	if aggr, ok := store.Aggregations[key]; ok {
		aggr.TimeElapsed -= elapsedTime
		if aggr.TimeElapsed < 0 {
			aggr.TimeElapsed = 0
		}
	}
	return nil
}

func (store *AppAggregationStore) GetAppAggregations() ([]*models.AppAggregation, error) {
	return utils.GetMapValues(store.Aggregations), nil
}

func (store *AppAggregationStore) GetAppAggregationsByDate(date datatypes.DateOnly) ([]*models.AppAggregation, error) {
	var aggregations []*models.AppAggregation
	for _, aggregation := range store.Aggregations {
		if aggregation.Date == date {
			aggregations = append(aggregations, aggregation)
		}
	}
	return aggregations, nil
}

func (store *AppAggregationStore) ReplaceAppAggregationsForDates(dates []datatypes.DateOnly, aggrs []*models.AppAggregation) error {
	dateSet := make(map[string]struct{}, len(dates))
	for _, d := range dates {
		dateSet[d.String()] = struct{}{}
	}

	for key, aggr := range store.Aggregations {
		if _, ok := dateSet[aggr.Date.String()]; ok {
			delete(store.Aggregations, key)
		}
	}

	for _, aggr := range aggrs {
		key := aggr.AppName + "-" + aggr.Date.String()
		copy := *aggr
		store.Aggregations[key] = &copy
	}

	return nil
}

func (store *AppAggregationStore) GetAppAggregationsByDateRange(start, end datatypes.DateOnly) ([]*models.AppAggregation, error) {
	var aggregations []*models.AppAggregation
	for _, aggregation := range store.Aggregations {
		if !aggregation.Date.Time.Before(start.Time) && !aggregation.Date.Time.After(end.Time) {
			aggregations = append(aggregations, aggregation)
		}
	}
	return aggregations, nil
}
