package core

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/models"
)

// EventExcludedError represents an error when an event is excluded by rules
type EventExcludedError struct {
	AppName     string
	RuleId      int
	Expression  string
	Description string
	LogInfo     bool
}

// Error implements the error interface
func (e *EventExcludedError) Error() string {
	if e.LogInfo {
		return fmt.Sprintf("event for app '%s' excluded by rule %d (%s): %s",
			e.AppName, e.RuleId, e.Expression, e.Description)
	}

	return "event excluded"
}

func NewEventExcludedError(appName string, rule *models.CategoryRule, description string, logInfo bool) *EventExcludedError {
	return &EventExcludedError{
		AppName:     appName,
		RuleId:      rule.RuleId,
		Expression:  rule.Expression,
		Description: description,
		LogInfo:     logInfo,
	}
}
