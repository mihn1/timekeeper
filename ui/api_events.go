package main

import (
	"fmt"

	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/ui/dtos"
)

// GetEventLog returns raw app-switch events for the given date (YYYY-MM-DD).
// Returns an empty slice in inmem mode (no events are persisted there).
func (a *App) GetEventLog(dateStr string) ([]*dtos.EventLogItem, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	date, err := datatypes.NewDateOnlyFromStr(dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format %q: %w", dateStr, err)
	}

	events, err := a.timekeeper.Storage.Events().GetEventsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %w", err)
	}

	return dtos.EventLogFromModels(events), nil
}
