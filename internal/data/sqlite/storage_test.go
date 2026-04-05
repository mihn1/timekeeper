package sqlite

import (
	"database/sql"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSqliteStorage(t *testing.T) {
	// Create temp DB file
	tmpFile := path.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", tmpFile)
	assert.NoError(t, err)

	storage := NewSqliteStorage(db)
	t.Cleanup(func() { storage.Close() })

	// Test adding and retrieving categories
	workCategory := &models.Category{Name: "Work"}
	err = storage.Categories().UpsertCategory(workCategory)
	assert.NoError(t, err)

	cat, err := storage.Categories().GetCategory(workCategory.Id)
	assert.NoError(t, err)
	assert.Equal(t, "Work", cat.Name)

	// Test rule storage
	rule := &models.CategoryRule{
		CategoryId: models.WORK,
		AppName:    "Code",
	}

	err = storage.Rules().UpsertRule(rule)
	assert.NoError(t, err)

	rules, err := storage.Rules().GetRules()
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
}

func TestSqliteEventStoreReturnsPersistedEventsByDate(t *testing.T) {
	tmpFile := path.Join(t.TempDir(), "events.db")
	db, err := sql.Open("sqlite3", tmpFile)
	assert.NoError(t, err)

	storage := NewSqliteStorage(db)
	t.Cleanup(func() { storage.Close() })

	start := time.Now().UTC().Add(-2 * time.Minute)
	end := start.Add(1 * time.Minute)
	event := &models.AppSwitchEvent{
		CategoryId: models.WORK,
		AppName:    "Code",
		StartTime:  start,
		EndTime:    end,
		AdditionalData: map[string]string{
			"source": "test",
		},
	}

	err = storage.Events().AddEvent(event)
	assert.NoError(t, err)

	events, err := storage.Events().GetEventsByDate(datatypes.NewDateOnly(start))
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "Code", events[0].AppName)
	assert.Equal(t, models.WORK, events[0].CategoryId)
}
