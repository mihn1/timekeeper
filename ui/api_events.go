package main

import (
	"fmt"

	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/internal/tzutil"
	"github.com/mihn1/timekeeper/ui/dtos"
)

// GetEventLog returns raw app-switch events for the given local date (YYYY-MM-DD),
// using the user's timezone preference to determine the correct UTC window.
// Returns an empty slice in inmem mode (no events are persisted there).
func (a *App) GetEventLog(dateStr string) ([]*dtos.EventLogItem, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	tz := a.getTimezone()
	loc := tzutil.LoadLoc(tz)

	// Prefer timezone-accurate range query.
	start, end, err := tzutil.LocalDayToUTCRange(dateStr, tz)
	if err == nil {
		events, qErr := a.timekeeper.Storage.Events().GetEventsByTimeRange(start, end)
		if qErr == nil {
			return dtos.EventLogFromModels(events, loc), nil
		}
	}

	// Fall back to UTC-date query.
	date, err := datatypes.NewDateOnlyFromStr(dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format %q: %w", dateStr, err)
	}
	events, err := a.timekeeper.Storage.Events().GetEventsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %w", err)
	}
	return dtos.EventLogFromModels(events, loc), nil
}

// DeleteEvent removes a single event and deducts its duration from both aggregation tables.
func (a *App) DeleteEvent(id int) error {
	if a.timekeeper == nil {
		return fmt.Errorf("timekeeper is not initialized")
	}

	eventId := models.EventId(id)
	ev, err := a.timekeeper.Storage.Events().GetEvent(eventId)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	elapsedMs := ev.EndTime.Sub(ev.StartTime).Milliseconds()
	if elapsedMs > 0 {
		date := ev.GetEventDate()
		_ = a.timekeeper.Storage.AppAggregations().DeductAppEvent(ev, elapsedMs)
		_ = a.timekeeper.Storage.CategoryAggregations().DeductCategory(ev.CategoryId, date, elapsedMs)
	}

	return a.timekeeper.Storage.Events().DeleteEvent(eventId)
}
