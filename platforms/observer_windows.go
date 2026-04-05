//go:build windows
// +build windows

package platforms

import (
	"log/slog"

	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/platforms/windows"
)

func NewPlatformObserver(cb func(models.AppSwitchEvent), standalone bool, logger *slog.Logger) interface {
	Start() error
	Stop() error
} {
	return windows.NewObserver(cb, standalone, logger)
}
