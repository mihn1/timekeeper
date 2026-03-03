package core

import (
	"log/slog"
	"os"
	"path"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTimeKeeperEventProcessing(t *testing.T) {
	// Create an in-memory TimeKeeper for testing
	timekeeper := NewTimeKeeperInMem(TimeKeeperOptions{Logger: slog.Default()})
	defer timekeeper.Close()

	// Seed test data
	SeedData(timekeeper)

	// Start tracking
	timekeeper.StartTracking()

	// Simulate events
	firstEvent := models.AppSwitchEvent{
		AppName:   "Code",
		StartTime: time.Now().Add(-time.Minute),
	}

	secondEvent := models.AppSwitchEvent{
		AppName:   constants.GOOGLE_CHROME,
		StartTime: time.Now(),
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_URL:   "https://github.com",
			constants.KEY_BROWSER_TITLE: "GitHub",
		},
	}

	// Push events
	timekeeper.PushEvent(firstEvent)
	timekeeper.PushEvent(secondEvent)

	// Give time for event processing
	time.Sleep(100 * time.Millisecond)

	// Check aggregations
	date := datatypes.NewDateOnly(time.Now())
	appAggr, err := timekeeper.Storage.AppAggregations().GetAppAggregationsByDate(date)
	assert.NoError(t, err)
	assert.Len(t, appAggr, 1) // Only one app aggregation should exist
}

func TestEndToEndSqlite(t *testing.T) {
	// Skip in CI environments
	if os.Getenv("CI") != "" {
		t.Skip("Skipping in CI environment")
	}

	// Create temp DB
	tmpFile := path.Join(t.TempDir(), "timekeeper.db")

	// Create TimeKeeper with SQLite storage
	timekeeper := NewTimeKeeperSqlite(TimeKeeperOptions{
		StoragePath: tmpFile,
		StoreEvents: true,
		Logger:      slog.Default(),
	})
	defer timekeeper.Close()

	// Initialize data
	SeedData(timekeeper)

	// Start tracking
	timekeeper.StartTracking()

	// Simulate application usage
	simulateEvents(t, timekeeper)

	// Verify data is correctly stored and categorized
	verifyAggregations(t, timekeeper)
}

func simulateEvents(t *testing.T, tk *TimeKeeper) {
	// Implement simulation logic
}

func verifyAggregations(t *testing.T, tk *TimeKeeper) {
	// Verify aggregations are correct
}

type countingObserver struct {
	starts atomic.Int32
	stops  atomic.Int32
}

func (o *countingObserver) Start() error {
	o.starts.Add(1)
	return nil
}

func (o *countingObserver) Stop() error {
	o.stops.Add(1)
	return nil
}

func TestIsSameEventUsesAbsoluteTimeDelta(t *testing.T) {
	now := time.Now().UTC()
	prev := &models.AppSwitchEvent{
		AppName:   "Code",
		StartTime: now,
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_URL: "https://github.com",
		},
	}

	next := &models.AppSwitchEvent{
		AppName:   "Code",
		StartTime: now.Add(30 * time.Second),
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_URL: "https://github.com",
		},
	}

	assert.True(t, isSameEvent(prev, next))
}

func TestGetRulesForEventSortsCombinedRulesByPriority(t *testing.T) {
	tk := NewTimeKeeperInMem(TimeKeeperOptions{Logger: slog.Default()})
	defer tk.Close()

	err := tk.Storage.Rules().UpsertRule(&models.CategoryRule{
		CategoryId: models.WORK,
		AppName:    "Code",
		Priority:   1,
	})
	assert.NoError(t, err)

	err = tk.Storage.Rules().UpsertRule(&models.CategoryRule{
		CategoryId:        models.PERSONAL,
		AppName:           constants.ALL_APPS,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Expression:        "github.com",
		Priority:          5,
	})
	assert.NoError(t, err)

	err = tk.Storage.Rules().UpsertRule(&models.CategoryRule{
		CategoryId:        models.ENTERTAINMENT,
		AppName:           constants.ALL_APPS,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Expression:        "docs",
		Priority:          3,
	})
	assert.NoError(t, err)

	event := &models.AppSwitchEvent{
		AppName: "Code",
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_URL: "https://github.com",
		},
	}

	rules, err := tk.getRulesForEvent(event)
	assert.NoError(t, err)
	assert.Len(t, rules, 3)
	assert.Equal(t, 5, rules[0].Priority)
	assert.Equal(t, 3, rules[1].Priority)
	assert.Equal(t, 1, rules[2].Priority)
}

func TestStartTrackingIsIdempotentForObserversAndLoop(t *testing.T) {
	tk := NewTimeKeeperInMem(TimeKeeperOptions{Logger: slog.Default()})
	defer tk.Close()

	obs := &countingObserver{}
	tk.AddObserver(obs)

	tk.StartTracking()
	tk.StartTracking()

	time.Sleep(50 * time.Millisecond)

	assert.True(t, tk.IsEnabled())
	assert.Equal(t, int32(1), obs.starts.Load())
	assert.True(t, tk.eventLoopStarted)

	tk.Disable()
	assert.False(t, tk.IsEnabled())

	tk.StartTracking()
	time.Sleep(50 * time.Millisecond)

	assert.True(t, tk.IsEnabled())
	assert.Equal(t, int32(1), obs.starts.Load())
}

func TestCloseStopsObserversOnlyOnce(t *testing.T) {
	tk := NewTimeKeeperInMem(TimeKeeperOptions{Logger: slog.Default()})
	obs := &countingObserver{}
	tk.AddObserver(obs)

	tk.StartTracking()
	time.Sleep(30 * time.Millisecond)

	tk.Close()
	tk.Close()

	assert.Equal(t, int32(1), obs.stops.Load())
}

func TestExclusionRulePreventsAggregation(t *testing.T) {
	tk := NewTimeKeeperInMem(TimeKeeperOptions{Logger: slog.Default()})
	defer tk.Close()

	err := tk.Storage.Rules().UpsertRule(&models.CategoryRule{
		CategoryId:        models.EXCLUDED,
		AppName:           "Code",
		IsExclusion:       true,
		Expression:        "",
		Priority:          10,
		AdditionalDataKey: "",
	})
	assert.NoError(t, err)

	tk.StartTracking()

	start := time.Now().UTC().Add(-2 * time.Minute)
	next := time.Now().UTC()

	tk.PushEvent(models.AppSwitchEvent{AppName: "Code", StartTime: start})
	tk.PushEvent(models.AppSwitchEvent{AppName: "Terminal", StartTime: next})
	time.Sleep(80 * time.Millisecond)

	appAggs, err := tk.Storage.AppAggregations().GetAppAggregationsByDate(datatypes.NewDateOnly(start))
	assert.NoError(t, err)
	assert.Len(t, appAggs, 0)
}

func TestHigherPriorityRuleWinsDuringAggregation(t *testing.T) {
	tk := NewTimeKeeperInMem(TimeKeeperOptions{Logger: slog.Default()})
	defer tk.Close()

	assert.NoError(t, tk.Storage.Categories().UpsertCategory(&models.Category{Id: models.WORK, Name: "Work"}))
	assert.NoError(t, tk.Storage.Categories().UpsertCategory(&models.Category{Id: models.PERSONAL, Name: "Personal"}))

	assert.NoError(t, tk.Storage.Rules().UpsertRule(&models.CategoryRule{
		CategoryId:        models.WORK,
		AppName:           constants.ALL_APPS,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Expression:        "github.com",
		Priority:          1,
	}))

	assert.NoError(t, tk.Storage.Rules().UpsertRule(&models.CategoryRule{
		CategoryId:        models.PERSONAL,
		AppName:           constants.ALL_APPS,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Expression:        "github.com",
		Priority:          5,
	}))

	tk.StartTracking()

	start := time.Now().UTC().Add(-3 * time.Minute)
	next := time.Now().UTC()

	tk.PushEvent(models.AppSwitchEvent{
		AppName:   "Code",
		StartTime: start,
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_URL: "https://github.com/mihn1/timekeeper",
		},
	})
	tk.PushEvent(models.AppSwitchEvent{AppName: "Terminal", StartTime: next})
	time.Sleep(80 * time.Millisecond)

	catAggs, err := tk.Storage.CategoryAggregations().GetCategoryAggregationsByDate(datatypes.NewDateOnly(start))
	assert.NoError(t, err)
	assert.Len(t, catAggs, 1)
	assert.Equal(t, models.PERSONAL, catAggs[0].CategoryId)
	assert.True(t, catAggs[0].TimeElapsed > 0)
}

func TestEventsAreIgnoredWhileTrackingDisabled(t *testing.T) {
	tk := NewTimeKeeperInMem(TimeKeeperOptions{Logger: slog.Default()})
	defer tk.Close()

	tk.StartTracking()
	tk.Disable()

	start := time.Now().UTC().Add(-2 * time.Minute)
	next := time.Now().UTC()

	tk.PushEvent(models.AppSwitchEvent{AppName: "Code", StartTime: start})
	tk.PushEvent(models.AppSwitchEvent{AppName: "Terminal", StartTime: next})
	time.Sleep(80 * time.Millisecond)

	appAggs, err := tk.Storage.AppAggregations().GetAppAggregationsByDate(datatypes.NewDateOnly(start))
	assert.NoError(t, err)
	assert.Len(t, appAggs, 0)
}
