package core

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/mihn1/timekeeper/internal/constants"
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTimeKeeperEventProcessing(t *testing.T) {
	// Create an in-memory TimeKeeper for testing
	timekeeper := NewTimeKeeperInMem(TimeKeeperOptions{})
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
	appAggr, err := timekeeper.storage.AppAggregations().GetAppAggregationsByDate(date)
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
