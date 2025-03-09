package main

import (
	"context"
	"log/slog"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// WailsHandler implements slog.Handler that redirects logs to Wails runtime
type WailsHandler struct {
	ctx context.Context
}

func NewWailsHandler(ctx context.Context) *WailsHandler {
	return &WailsHandler{ctx: ctx}
}

// Enabled implements slog.Handler.
func (h *WailsHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.ctx != nil // Only enabled if context is set
}

// Handle implements slog.Handler.
func (h *WailsHandler) Handle(_ context.Context, record slog.Record) error {
	if h.ctx == nil {
		return nil
	}

	// Convert record to message
	message := record.Message

	// Add attributes if any exist
	attrs := make(map[string]any)
	hasAttrs := false

	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		hasAttrs = true
		return true
	})

	// Send to appropriate Wails log function based on level
	switch {
	case record.Level >= slog.LevelError:
		if hasAttrs {
			runtime.LogErrorf(h.ctx, "%s %v", message, attrs)
		} else {
			runtime.LogError(h.ctx, message)
		}
	case record.Level >= slog.LevelWarn:
		if hasAttrs {
			runtime.LogWarningf(h.ctx, "%s %v", message, attrs)
		} else {
			runtime.LogWarning(h.ctx, message)
		}
	case record.Level >= slog.LevelInfo:
		if hasAttrs {
			runtime.LogInfof(h.ctx, "%s %v", message, attrs)
		} else {
			runtime.LogInfo(h.ctx, message)
		}
	default:
		if hasAttrs {
			runtime.LogDebugf(h.ctx, "%s %v", message, attrs)
		} else {
			runtime.LogDebug(h.ctx, message)
		}
	}

	return nil
}

// WithAttrs implements slog.Handler.
func (h *WailsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we'll just return the same handler
	// In a production-grade implementation, you would clone and store the attrs
	return h
}

// WithGroup implements slog.Handler.
func (h *WailsHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we'll just return the same handler
	return h
}
