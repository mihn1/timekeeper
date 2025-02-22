package core

import (
	"github.com/mihn1/timekeeper/internal/constants"
	"github.com/mihn1/timekeeper/internal/models"
)

func SeedDataInMem(t *TimeKeeper) {

	t.storage.Categories().AddCategory(models.Category{Id: models.WORK, Name: "Work"})
	t.storage.Categories().AddCategory(models.Category{Id: models.ENTERTAINMENT, Name: "Entertainment"})
	t.storage.Categories().AddCategory(models.Category{Id: models.PERSONAL, Name: "Personal"})
	t.storage.Categories().AddCategory(models.Category{Id: models.UNDEFINED, Name: "Undefined"})

	rules := make([]models.CategoryRule, 0)
	rules = append(rules, models.CategoryRule{CategoryId: models.PERSONAL, AppName: constants.GOOGLE_CHROME})
	rules = append(rules, models.CategoryRule{CategoryId: models.WORK, AppName: "Code"})
	rules = append(rules, models.CategoryRule{CategoryId: models.WORK, AppName: "Ghostty"})
	rules = append(rules, models.CategoryRule{CategoryId: models.PERSONAL, AppName: "ChatGPT"})
	rules = append(rules, models.CategoryRule{CategoryId: models.PERSONAL, AppName: "Notion"})

	// Rules based on browsers' tabs
	rules = append(rules, models.CategoryRule{
		CategoryId:        models.WORK,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          3,
		Expression:        "github.com"})
	rules = append(rules, models.CategoryRule{
		CategoryId:        models.PERSONAL,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          1,
		Expression:        "chatgpt.com"})
	rules = append(rules, models.CategoryRule{
		CategoryId:        models.WORK,
		AppName:           constants.ALL_APPS,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          2,
		Expression:        "developer"})
	rules = append(rules, models.CategoryRule{
		CategoryId:        models.WORK,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_TITLE,
		Priority:          2,
		Expression:        "^work",
		IsRegex:           true})
	rules = append(rules, models.CategoryRule{
		CategoryId:        models.ENTERTAINMENT,
		AppName:           constants.ALL_APPS,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          1,
		Expression:        "youtube.com"})
	rules = append(rules, models.CategoryRule{
		CategoryId:        models.ENTERTAINMENT,
		AppName:           constants.ALL_APPS,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          1,
		Expression:        "twitch.tv"})

	for idx, rule := range rules {
		rule.RuleId = idx + 1
		t.storage.Rules().AddRule(rule)
	}
}
