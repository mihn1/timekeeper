package models

import (
	"regexp"
	"strings"
	"sync"
)

var regexCache sync.Map

func getCompiledRegex(pattern string) (*regexp.Regexp, error) {
	if cached, ok := regexCache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	actual, _ := regexCache.LoadOrStore(pattern, compiled)
	return actual.(*regexp.Regexp), nil
}

type CategoryRule struct {
	RuleId            int
	CategoryId        CategoryId
	AppName           string
	AdditionalDataKey string
	Expression        string
	IsRegex           bool
	Priority          int
	IsExclusion       bool
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
		pattern, err := getCompiledRegex("(?i)" + r.Expression)
		if err != nil {
			return false, nil
		}

		return pattern.MatchString(val), nil
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
