package main

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

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

// formatAttrs formats a map of attributes as key=value pairs in a clean format
func formatAttrs(attrs map[string]any) string {
	if len(attrs) == 0 {
		return ""
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build formatted string
	var pairs []string
	for _, k := range keys {
		v := attrs[k]
		// Handle different value types for better formatting
		switch val := v.(type) {
		case string:
			pairs = append(pairs, fmt.Sprintf("%s=\"%s\"", k, val))
		case error:
			pairs = append(pairs, fmt.Sprintf("%s=\"%s\"", k, val.Error()))
		default:
			pairs = append(pairs, fmt.Sprintf("%s=%v", k, val))
		}
	}

	return strings.Join(pairs, " ")
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

	// Format the log message
	var logMsg string
	if hasAttrs {
		formattedAttrs := formatAttrs(attrs)
		logMsg = fmt.Sprintf("%s %s", message, formattedAttrs)
	} else {
		logMsg = message
	}

	// Send to appropriate Wails log function based on level
	switch {
	case record.Level >= slog.LevelError:
		runtime.LogError(h.ctx, logMsg)
	case record.Level >= slog.LevelWarn:
		runtime.LogWarning(h.ctx, logMsg)
	case record.Level >= slog.LevelInfo:
		runtime.LogInfo(h.ctx, logMsg)
	default:
		runtime.LogDebug(h.ctx, logMsg)
	}

	return nil
}

// WithAttrs implements slog.Handler.
func (h *WailsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we'll just return the same handler
	// TODO: clone and store the attrs
	return h
}

// WithGroup implements slog.Handler.
func (h *WailsHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we'll just return the same handler
	return h
}
