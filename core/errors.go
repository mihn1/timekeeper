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
	e := &EventExcludedError{
		AppName:     appName,
		Description: description,
		LogInfo:     logInfo,
	}
	if rule != nil {
		e.RuleId = rule.RuleId
		e.Expression = rule.Expression
	}
	return e
}
