package core

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mihn1/timekeeper/core/resolvers"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

// RerunJobState represents the lifecycle state of a rerun job.
type RerunJobState string

const (
	RerunStateIdle      RerunJobState = "idle"
	RerunStateRunning   RerunJobState = "running"
	RerunStateCompleted RerunJobState = "completed"
	RerunStateFailed    RerunJobState = "failed"
)

const defaultMaxRerunRangeDays = 7

// RerunJobStatus is a snapshot of the current (or last) rerun job.
type RerunJobStatus struct {
	State           RerunJobState `json:"state"`
	StartDate       string        `json:"startDate"`
	EndDate         string        `json:"endDate"`
	TotalEvents     int           `json:"totalEvents"`
	ProcessedEvents int           `json:"processedEvents"`
	ErrorMessage    string        `json:"errorMessage,omitempty"`
	StartedAt       time.Time     `json:"startedAt"`
	CompletedAt     time.Time     `json:"completedAt,omitempty"`
}

// rerunController holds runtime job state and coordination primitives.
type rerunController struct {
	mu                sync.RWMutex
	status            RerunJobStatus
	running           bool
	maxRangeDays      int
	statusCallback    func(RerunJobStatus)
}

func newRerunController() *rerunController {
	return &rerunController{
		status:       RerunJobStatus{State: RerunStateIdle},
		maxRangeDays: defaultMaxRerunRangeDays,
	}
}

// SetMaxRerunRangeDays configures the maximum inclusive date-range size for rerun jobs.
// Values < 1 are coerced to 1.
func (t *TimeKeeper) SetMaxRerunRangeDays(days int) {
	t.ensureRerunController()
	if days < 1 {
		days = 1
	}
	t.rerun.mu.Lock()
	t.rerun.maxRangeDays = days
	t.rerun.mu.Unlock()
}

// SetRerunStatusCallback registers a callback invoked on every rerun status change.
// The callback runs synchronously on the worker goroutine; keep it fast (e.g. emit a Wails event).
func (t *TimeKeeper) SetRerunStatusCallback(cb func(RerunJobStatus)) {
	t.ensureRerunController()
	t.rerun.mu.Lock()
	t.rerun.statusCallback = cb
	t.rerun.mu.Unlock()
}

// GetRerunJobStatus returns the current rerun job status snapshot.
func (t *TimeKeeper) GetRerunJobStatus() RerunJobStatus {
	t.ensureRerunController()
	t.rerun.mu.RLock()
	defer t.rerun.mu.RUnlock()
	return t.rerun.status
}

// GetMaxRerunRangeDays returns the configured max rerun window.
func (t *TimeKeeper) GetMaxRerunRangeDays() int {
	t.ensureRerunController()
	t.rerun.mu.RLock()
	defer t.rerun.mu.RUnlock()
	return t.rerun.maxRangeDays
}

func (t *TimeKeeper) ensureRerunController() {
	t.stateMu.Lock()
	defer t.stateMu.Unlock()
	if t.rerun == nil {
		t.rerun = newRerunController()
	}
}

// StartRerunRules launches a background rerun across the local date range [startDate, endDate].
// Returns an error synchronously if a job is already running, the range is invalid, or the range
// exceeds the configured maximum.
func (t *TimeKeeper) StartRerunRules(startDate, endDate datatypes.DateOnly, timezone string) error {
	t.ensureRerunController()

	if endDate.Time.Before(startDate.Time) {
		return fmt.Errorf("end date must be on or after start date")
	}

	days := int(endDate.Time.Sub(startDate.Time).Hours()/24) + 1

	t.rerun.mu.Lock()
	if t.rerun.running {
		t.rerun.mu.Unlock()
		return errors.New("a rerun job is already in progress")
	}
	if days > t.rerun.maxRangeDays {
		maxDays := t.rerun.maxRangeDays
		t.rerun.mu.Unlock()
		return fmt.Errorf("date range %d days exceeds maximum of %d days", days, maxDays)
	}

	now := time.Now()
	t.rerun.running = true
	t.rerun.status = RerunJobStatus{
		State:     RerunStateRunning,
		StartDate: startDate.String(),
		EndDate:   endDate.String(),
		StartedAt: now,
	}
	cb := t.rerun.statusCallback
	snapshot := t.rerun.status
	t.rerun.mu.Unlock()

	if cb != nil {
		cb(snapshot)
	}

	go t.runRerunJob(startDate, endDate, timezone)
	return nil
}

func (t *TimeKeeper) runRerunJob(startDate, endDate datatypes.DateOnly, timezone string) {
	err := t.doRerun(startDate, endDate, timezone)

	t.rerun.mu.Lock()
	t.rerun.running = false
	t.rerun.status.CompletedAt = time.Now()
	if err != nil {
		t.rerun.status.State = RerunStateFailed
		t.rerun.status.ErrorMessage = err.Error()
	} else {
		t.rerun.status.State = RerunStateCompleted
	}
	cb := t.rerun.statusCallback
	snapshot := t.rerun.status
	t.rerun.mu.Unlock()

	if cb != nil {
		cb(snapshot)
	}
}

func (t *TimeKeeper) updateRerunProgress(processed, total int) {
	t.rerun.mu.Lock()
	t.rerun.status.ProcessedEvents = processed
	t.rerun.status.TotalEvents = total
	cb := t.rerun.statusCallback
	snapshot := t.rerun.status
	t.rerun.mu.Unlock()

	if cb != nil {
		cb(snapshot)
	}
}

// doRerun is the actual work. Reads events in range, re-resolves categories against
// current rules, rebuilds aggregations in memory, then applies the result to the
// storage in one batch per table.
func (t *TimeKeeper) doRerun(startDate, endDate datatypes.DateOnly, timezone string) error {
	startUTC, endUTC, err := localDayRangeToUTC(startDate, endDate, timezone)
	if err != nil {
		return fmt.Errorf("invalid date range: %w", err)
	}

	events, err := t.Storage.Events().GetEventsByTimeRange(startUTC, endUTC)
	if err != nil {
		return fmt.Errorf("failed to load events: %w", err)
	}

	total := len(events)
	t.updateRerunProgress(0, total)

	if total == 0 {
		// Still clear stale aggregations for the selected dates so the UI reflects an empty range.
		dates := enumerateDates(startDate, endDate)
		if err := t.Storage.AppAggregations().ReplaceAppAggregationsForDates(dates, nil); err != nil {
			return fmt.Errorf("failed to clear app aggregations: %w", err)
		}
		if err := t.Storage.CategoryAggregations().ReplaceCategoryAggregationsForDates(dates, nil); err != nil {
			return fmt.Errorf("failed to clear category aggregations: %w", err)
		}
		return nil
	}

	appAggr := make(map[string]*models.AppAggregation)
	catAggr := make(map[string]*models.CategoryAggregation)
	dateSet := make(map[string]datatypes.DateOnly)

	type eventUpdate struct {
		id    models.EventId
		catId models.CategoryId
	}
	updates := make([]eventUpdate, 0, total)

	for idx, ev := range events {
		evDate := ev.GetEventDate()
		dateSet[evDate.String()] = evDate

		elapsed := ev.EndTime.Sub(ev.StartTime).Milliseconds()
		if elapsed <= 0 {
			if (idx+1)%50 == 0 || idx+1 == total {
				t.updateRerunProgress(idx+1, total)
			}
			continue
		}

		catId, excluded, err := t.resolveCategoryForRerun(ev)
		if err != nil {
			t.logger.Error("Rerun: failed to resolve category", "event", ev.Id, "error", err)
			if (idx+1)%50 == 0 || idx+1 == total {
				t.updateRerunProgress(idx+1, total)
			}
			continue
		}

		// Always update the stored event's category so downstream views stay consistent.
		updates = append(updates, eventUpdate{id: ev.Id, catId: catId})

		if excluded {
			if (idx+1)%50 == 0 || idx+1 == total {
				t.updateRerunProgress(idx+1, total)
			}
			continue
		}

		// App aggregation (always counted, even for UNDEFINED category)
		appKey := ev.AppName + "-" + evDate.String()
		aa, ok := appAggr[appKey]
		if !ok {
			aa = &models.AppAggregation{AppName: ev.AppName, Date: evDate}
			appAggr[appKey] = aa
		}
		aa.TimeElapsed += elapsed

		// Category aggregation (EXCLUDED already short-circuited above)
		catKey := models.GetCategoryAggregationKey(catId, evDate)
		ca, ok := catAggr[catKey]
		if !ok {
			ca = &models.CategoryAggregation{CategoryId: catId, Date: evDate}
			catAggr[catKey] = ca
		}
		ca.TimeElapsed += elapsed

		if (idx+1)%50 == 0 || idx+1 == total {
			t.updateRerunProgress(idx+1, total)
		}
	}

	dates := make([]datatypes.DateOnly, 0, len(dateSet))
	for _, d := range dateSet {
		dates = append(dates, d)
	}

	newAppAggrs := make([]*models.AppAggregation, 0, len(appAggr))
	for _, a := range appAggr {
		newAppAggrs = append(newAppAggrs, a)
	}
	newCatAggrs := make([]*models.CategoryAggregation, 0, len(catAggr))
	for _, c := range catAggr {
		newCatAggrs = append(newCatAggrs, c)
	}

	if err := t.Storage.AppAggregations().ReplaceAppAggregationsForDates(dates, newAppAggrs); err != nil {
		return fmt.Errorf("failed to replace app aggregations: %w", err)
	}
	if err := t.Storage.CategoryAggregations().ReplaceCategoryAggregationsForDates(dates, newCatAggrs); err != nil {
		return fmt.Errorf("failed to replace category aggregations: %w", err)
	}

	for _, u := range updates {
		if err := t.Storage.Events().UpdateEventCategory(u.id, u.catId); err != nil {
			t.logger.Error("Rerun: failed to update event category", "event", u.id, "error", err)
		}
	}

	t.updateRerunProgress(total, total)
	return nil
}

// resolveCategoryForRerun mirrors handleEvent's rule + exclusion flow without touching aggregations.
// Returns (catId, excluded, err). An excluded event should still have its CategoryId set to EXCLUDED
// on the stored event but should not contribute to aggregations.
func (t *TimeKeeper) resolveCategoryForRerun(event *models.AppSwitchEvent) (models.CategoryId, bool, error) {
	rules, err := t.getRulesForEvent(event)
	if err != nil {
		if _, ok := err.(*EventExcludedError); ok {
			return models.EXCLUDED, true, nil
		}
		return models.UNDEFINED, false, err
	}

	if defaultResolver == nil {
		defaultResolver = resolvers.NewDefaultCategoryResolver(t.Storage.Rules(), t.Storage.Categories())
	}

	catId, err := defaultResolver.ResolveCategory(event, rules)
	if err != nil {
		return models.UNDEFINED, false, err
	}

	if catId == models.EXCLUDED {
		return catId, true, nil
	}

	return catId, false, nil
}

func localDayRangeToUTC(startDate, endDate datatypes.DateOnly, timezone string) (time.Time, time.Time, error) {
	loc, err := loadTimezone(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	startLocal := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	endLocal := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, 1)
	return startLocal.UTC(), endLocal.UTC(), nil
}

func loadTimezone(tz string) (*time.Location, error) {
	if tz == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(tz)
}

func enumerateDates(start, end datatypes.DateOnly) []datatypes.DateOnly {
	var dates []datatypes.DateOnly
	cur := start.Time
	for !cur.After(end.Time) {
		dates = append(dates, datatypes.NewDateOnly(cur))
		cur = cur.AddDate(0, 0, 1)
	}
	return dates
}
