package inmem

import (
	"log"

	"github.com/mihn1/timekeeper/internal/datatypes"
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
	log.Println("Aggregating app event inmem", event)
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
