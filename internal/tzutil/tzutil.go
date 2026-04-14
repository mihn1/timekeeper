// Package tzutil provides helpers for timezone-aware date handling.
package tzutil

import (
	"time"

	"github.com/mihn1/timekeeper/internal/models"
)

// LocalDayToUTCRange converts a local calendar date ("YYYY-MM-DD") to the
// corresponding UTC [start, end) interval for that timezone.
// Falls back to UTC if the timezone string is invalid.
func LocalDayToUTCRange(dateStr, timezone string) (start, end time.Time, err error) {
	loc := loadLoc(timezone)
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return
	}
	start = t.UTC()
	end = t.AddDate(0, 0, 1).UTC()
	return
}

// LocalDateForTime returns the "YYYY-MM-DD" string for a UTC time in the given timezone.
func LocalDateForTime(t time.Time, timezone string) string {
	return t.In(loadLoc(timezone)).Format("2006-01-02")
}

// LoadLoc returns the time.Location for an IANA timezone, falling back to UTC.
func LoadLoc(timezone string) *time.Location {
	return loadLoc(timezone)
}

func loadLoc(timezone string) *time.Location {
	if timezone == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.UTC
	}
	return loc
}

// AggregateEventsByApp sums elapsed time (ms) per app name from a slice of stored events.
// Events without a closed EndTime are skipped (still in progress).
func AggregateEventsByApp(events []*models.AppSwitchEvent) map[string]int64 {
	totals := make(map[string]int64)
	for _, ev := range events {
		if ev.EndTime.IsZero() {
			continue
		}
		if ms := ev.EndTime.Sub(ev.StartTime).Milliseconds(); ms > 0 {
			totals[ev.AppName] += ms
		}
	}
	return totals
}

// AggregateEventsByCategory sums elapsed time (ms) per CategoryId from a slice of stored events.
func AggregateEventsByCategory(events []*models.AppSwitchEvent) map[models.CategoryId]int64 {
	totals := make(map[models.CategoryId]int64)
	for _, ev := range events {
		if ev.EndTime.IsZero() {
			continue
		}
		if ms := ev.EndTime.Sub(ev.StartTime).Milliseconds(); ms > 0 {
			totals[ev.CategoryId] += ms
		}
	}
	return totals
}

// AppCategoryMap builds appName → latest CategoryId from a slice of events.
// Later events overwrite earlier ones for the same app.
func AppCategoryMap(events []*models.AppSwitchEvent) map[string]models.CategoryId {
	result := make(map[string]models.CategoryId)
	for _, ev := range events {
		result[ev.AppName] = ev.CategoryId
	}
	return result
}
