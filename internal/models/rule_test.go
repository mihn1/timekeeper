package models

import (
	"testing"

	"github.com/mihn1/timekeeper/constants"
	"github.com/stretchr/testify/assert"
)

func TestRegexRuleMatchingIsCaseInsensitive(t *testing.T) {
	rule := &CategoryRule{
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_TITLE,
		Expression:        "^work",
		IsRegex:           true,
	}

	event := &AppSwitchEvent{
		AppName: constants.GOOGLE_CHROME,
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_TITLE: "WORK - sprint planning",
		},
	}

	matched, err := rule.IsMatch(event)
	assert.NoError(t, err)
	assert.True(t, matched)
}

func TestInvalidRegexRuleDoesNotFailMatchingFlow(t *testing.T) {
	rule := &CategoryRule{
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_TITLE,
		Expression:        "[invalid",
		IsRegex:           true,
	}

	event := &AppSwitchEvent{
		AppName: constants.GOOGLE_CHROME,
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_TITLE: "work - sprint planning",
		},
	}

	matched, err := rule.IsMatch(event)
	assert.NoError(t, err)
	assert.False(t, matched)
}
