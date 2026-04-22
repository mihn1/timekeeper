package main

import (
	"fmt"

	"github.com/mihn1/timekeeper/core"
	"github.com/mihn1/timekeeper/datatypes"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// RerunJobStatusDTO mirrors core.RerunJobStatus for JSON serialization to the frontend.
type RerunJobStatusDTO struct {
	State           string `json:"state"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	TotalEvents     int    `json:"totalEvents"`
	ProcessedEvents int    `json:"processedEvents"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
	StartedAt       string `json:"startedAt,omitempty"`
	CompletedAt     string `json:"completedAt,omitempty"`
	MaxRangeDays    int    `json:"maxRangeDays"`
}

func (a *App) rerunStatusToDTO(s core.RerunJobStatus) RerunJobStatusDTO {
	dto := RerunJobStatusDTO{
		State:           string(s.State),
		StartDate:       s.StartDate,
		EndDate:         s.EndDate,
		TotalEvents:     s.TotalEvents,
		ProcessedEvents: s.ProcessedEvents,
		ErrorMessage:    s.ErrorMessage,
	}
	if !s.StartedAt.IsZero() {
		dto.StartedAt = s.StartedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if !s.CompletedAt.IsZero() {
		dto.CompletedAt = s.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if a.timekeeper != nil {
		dto.MaxRangeDays = a.timekeeper.GetMaxRerunRangeDays()
	}
	return dto
}

// StartRerunRules kicks off a background rerun of category rules across the given local date range.
// Dates are YYYY-MM-DD strings in the user's configured timezone.
func (a *App) StartRerunRules(startDate, endDate string) error {
	if a.timekeeper == nil {
		return fmt.Errorf("timekeeper is not initialized")
	}

	start, err := datatypes.NewDateOnlyFromStr(startDate)
	if err != nil {
		return fmt.Errorf("invalid start date %q: %w", startDate, err)
	}
	end, err := datatypes.NewDateOnlyFromStr(endDate)
	if err != nil {
		return fmt.Errorf("invalid end date %q: %w", endDate, err)
	}

	return a.timekeeper.StartRerunRules(start, end, a.getTimezone())
}

// GetRerunJobStatus returns the current (or last) rerun job status.
func (a *App) GetRerunJobStatus() RerunJobStatusDTO {
	if a.timekeeper == nil {
		return RerunJobStatusDTO{State: "idle"}
	}
	return a.rerunStatusToDTO(a.timekeeper.GetRerunJobStatus())
}

// emitRerunStatus forwards backend status updates to the frontend via a runtime event.
func (a *App) emitRerunStatus(s core.RerunJobStatus) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "timekeeper:rerun-status", a.rerunStatusToDTO(s))

	// When the job completes successfully, also fire a generic data-updated event so
	// open views (dashboard, event log) refetch fresh numbers.
	if s.State == core.RerunStateCompleted {
		runtime.EventsEmit(a.ctx, "timekeeper:data-updated")
	}
}
