package inmem

import (
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

type EventStore struct {
	events map[string][]*models.AppSwitchEvent
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: make(map[string][]*models.AppSwitchEvent),
	}
}

func (s *EventStore) AddEvent(event *models.AppSwitchEvent) error {
	date := event.GetEventDate().String()
	if _, ok := s.events[date]; !ok {
		s.events[date] = []*models.AppSwitchEvent{}
	}
	s.events[date] = append(s.events[date], event)
	return nil
}

func (s *EventStore) GetEventsByDate(date datatypes.DateOnly) ([]*models.AppSwitchEvent, error) {
	return s.events[date.String()], nil
}
