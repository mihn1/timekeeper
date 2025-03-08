package resolvers

import (
	"testing"

	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestRuleMatching(t *testing.T) {
	// Test simple app name matching
	rule := models.CategoryRule{
		RuleId:     1,
		CategoryId: models.WORK,
		AppName:    "Code",
	}

	event := &models.AppSwitchEvent{
		AppName: "Code",
	}

	match, err := rule.IsMatch(event)
	assert.NoError(t, err)
	assert.True(t, match)

	// Test additional data matching
	rule = models.CategoryRule{
		RuleId:            2,
		CategoryId:        models.ENTERTAINMENT,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Expression:        "youtube.com",
	}

	event = &models.AppSwitchEvent{
		AppName: constants.GOOGLE_CHROME,
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_URL: "https://www.youtube.com/watch",
		},
	}

	match, err = rule.IsMatch(event)
	assert.NoError(t, err)
	assert.True(t, match)

	// Test regex matching
	rule = models.CategoryRule{
		RuleId:            3,
		CategoryId:        models.WORK,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_TITLE,
		Expression:        "^work",
		IsRegex:           true,
	}

	event = &models.AppSwitchEvent{
		AppName: constants.GOOGLE_CHROME,
		AdditionalData: map[string]string{
			constants.KEY_BROWSER_TITLE: "work - Some Project",
		},
	}

	match, err = rule.IsMatch(event)
	assert.NoError(t, err)
	assert.True(t, match)
}
