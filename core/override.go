package core

import (
	"errors"
	"fmt"

	"github.com/mihn1/timekeeper/internal/models"
)

// OverrideEventCategory manually sets a stored event's category and rebalances
// the category aggregation for its date.
//
// Concurrency: holds rerunController.mu for the duration so a rerun job cannot
// start while an override is in flight (a rerun would otherwise clobber the
// override via ReplaceCategoryAggregationsForDates / UpdateEventCategory).
// A running rerun causes this call to return an error; the UI should also
// disable the override action while a rerun is running.
//
// App aggregation is intentionally NOT modified: it is keyed by app name +
// date, which does not change when a category override is applied.
func (t *TimeKeeper) OverrideEventCategory(id models.EventId, newCatId models.CategoryId) error {
	t.ensureRerunController()

	t.rerun.mu.Lock()
	defer t.rerun.mu.Unlock()

	if t.rerun.running {
		return errors.New("cannot override event category while a rerun job is in progress")
	}

	ev, err := t.Storage.Events().GetEvent(id)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	if ev.CategoryId == newCatId {
		return nil
	}

	var newCat *models.Category
	if newCatId != models.EXCLUDED {
		newCat, err = t.Storage.Categories().GetCategory(newCatId)
		if err != nil {
			return fmt.Errorf("invalid category %d: %w", newCatId, err)
		}
	}

	elapsedMs := ev.EndTime.Sub(ev.StartTime).Milliseconds()
	date := ev.GetEventDate()

	if elapsedMs > 0 {
		if ev.CategoryId != models.EXCLUDED {
			if err := t.Storage.CategoryAggregations().DeductCategory(ev.CategoryId, date, elapsedMs); err != nil {
				t.logger.Error("Override: failed to deduct old category aggregation", "event", id, "oldCat", ev.CategoryId, "error", err)
			}
		}
		if newCat != nil {
			if _, err := t.Storage.CategoryAggregations().AggregateCategory(newCat, date, elapsedMs); err != nil {
				return fmt.Errorf("failed to update new category aggregation: %w", err)
			}
		}
	}

	if err := t.Storage.Events().UpdateEventCategory(id, newCatId); err != nil {
		return fmt.Errorf("failed to update event category: %w", err)
	}

	t.logger.Info("Event category overridden", "event", id, "oldCat", ev.CategoryId, "newCat", newCatId, "elapsedMs", elapsedMs)
	return nil
}
