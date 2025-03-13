package models

import (
	"regexp"
	"strings"
)

type CategoryRule struct {
	RuleId            int
	CategoryId        CategoryId
	AppName           string
	AdditionalDataKey string
	Expression        string
	IsRegex           bool
	Priority          int
}

func (r *CategoryRule) IsMatch(event *AppSwitchEvent) (bool, error) {
	if r.AdditionalDataKey == "" || r.Expression == "" || event.AdditionalData == nil {
		return r.AppName == event.AppName, nil
	}

	val, ok := event.AdditionalData[r.AdditionalDataKey]
	if !ok {
		return false, nil
	}

	if r.IsRegex {
		// Implement regex matching
		// TODO: can we cache the compiled regex?
		pattern, err := regexp.Compile("(?i)" + r.Expression)
		if err != nil {
			return false, nil
		}

		return pattern.MatchString(val), err
	}

	return strings.Contains(strings.ToLower(val), strings.ToLower(r.Expression)), nil
}

func CmpRules(x, y *CategoryRule) int {
	if x.Priority > y.Priority {
		return -1
	} else if x.Priority < y.Priority {
		return 1
	}
	return 0
}
