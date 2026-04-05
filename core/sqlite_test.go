//go:build !nosqlite

package core

import (
	"log/slog"
	"path"
	"testing"
)

func TestEndToEndSqlite(t *testing.T) {
	tmpFile := path.Join(t.TempDir(), "timekeeper.db")

	tk := NewTimeKeeperSqlite(TimeKeeperOptions{
		StoragePath: tmpFile,
		StoreEvents: true,
		Logger:      slog.Default(),
	})
	defer tk.Close()

	SeedData(tk)
	tk.StartTracking()

	simulateEvents(t, tk)
	assertAggregations(t, tk)
}
