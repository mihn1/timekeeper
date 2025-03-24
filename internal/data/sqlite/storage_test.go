package sqlite

import (
	"database/sql"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSqliteStorage(t *testing.T) {
	// Create temp DB file
	tmpFile := path.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", tmpFile)
	assert.NoError(t, err)

	storage := NewSqliteStorage(db)

	// Test adding and retrieving categories
	err = storage.Categories().UpsertCategory(&models.Category{
		Id:   models.WORK,
		Name: "Work",
	})
	assert.NoError(t, err)

	cat, err := storage.Categories().GetCategory(models.WORK)
	assert.NoError(t, err)
	assert.Equal(t, "Work", cat.Name)

	// Test rule storage
	rule := &models.CategoryRule{
		RuleId:     1,
		CategoryId: models.WORK,
		AppName:    "Code",
	}

	err = storage.Rules().UpsertRule(rule)
	assert.NoError(t, err)

	rules, err := storage.Rules().GetRules()
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
}
