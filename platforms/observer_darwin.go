//go:build darwin
// +build darwin

package platforms

import (
	"log/slog"

	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/platforms/darwin"
)

// NewPlatformObserver returns the macOS observer on darwin.
func NewPlatformObserver(cb func(models.AppSwitchEvent), standalone bool, logger *slog.Logger) interface {
	Start() error
	Stop() error
} {
	return darwin.NewObserver(cb, standalone, logger)
}
