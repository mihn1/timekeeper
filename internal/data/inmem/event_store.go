package inmem

import (
	"fmt"
	"time"

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

func (s *EventStore) GetEvent(id models.EventId) (*models.AppSwitchEvent, error) {
	for _, evs := range s.events {
		for _, ev := range evs {
			if ev.Id == id {
				return ev, nil
			}
		}
	}
	return nil, fmt.Errorf("event %d not found", id)
}

func (s *EventStore) DeleteEvent(id models.EventId) error {
	for dateKey, evs := range s.events {
		for i, ev := range evs {
			if ev.Id == id {
				s.events[dateKey] = append(evs[:i], evs[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (s *EventStore) GetEventsByTimeRange(start, end time.Time) ([]*models.AppSwitchEvent, error) {
	var result []*models.AppSwitchEvent
	for _, evs := range s.events {
		for _, ev := range evs {
			if !ev.StartTime.Before(start) && ev.StartTime.Before(end) {
				result = append(result, ev)
			}
		}
	}
	return result, nil
}
