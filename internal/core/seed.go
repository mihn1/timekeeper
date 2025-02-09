package core

import (
	"github.com/mihn1/timekeeper/internal/constants"
	"github.com/mihn1/timekeeper/internal/models"
)

func SeedDataInMem(t *TimeKeeper) {
	t.storage.CategoryStore.AddCategory(models.Category{Id: models.WORK, Name: "Work"})
	t.storage.CategoryStore.AddCategory(models.Category{Id: models.ENTERTAINMENT, Name: "Entertainment"})
	t.storage.CategoryStore.AddCategory(models.Category{Id: models.PERSONAL, Name: "Personal"})
	t.storage.CategoryStore.AddCategory(models.Category{Id: models.UNDEFINED, Name: "Undefined"})

	t.storage.RuleStore.AddRule(models.CategoryRule{RuleId: 1, CategoryId: models.PERSONAL, AppName: constants.GOOGLE_CHROME})
	t.storage.RuleStore.AddRule(models.CategoryRule{RuleId: 2, CategoryId: models.WORK, AppName: "Code"})
	t.storage.RuleStore.AddRule(models.CategoryRule{RuleId: 3, CategoryId: models.WORK, AppName: "Ghostty"})
	t.storage.RuleStore.AddRule(models.CategoryRule{RuleId: 4, CategoryId: models.PERSONAL, AppName: "ChatGPT"})
	t.storage.RuleStore.AddRule(models.CategoryRule{RuleId: 5, CategoryId: models.PERSONAL, AppName: "Notion"})

	// Rules based on chromes' tabs
	t.storage.RuleStore.AddRule(models.CategoryRule{
		RuleId:            6,
		CategoryId:        models.WORK,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          3,
		Expression:        "github.com"})

	t.storage.RuleStore.AddRule(models.CategoryRule{
		RuleId:            7,
		CategoryId:        models.PERSONAL,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          1,
		Expression:        "chatgpt.com"})

	t.storage.RuleStore.AddRule(models.CategoryRule{
		RuleId:            8,
		CategoryId:        models.WORK,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          2,
		Expression:        "developer"})

	t.storage.RuleStore.AddRule(models.CategoryRule{
		RuleId:            9,
		CategoryId:        models.WORK,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_TITLE,
		Priority:          2,
		Expression:        "^work",
		IsRegex:           true})

	t.storage.RuleStore.AddRule(models.CategoryRule{
		RuleId:            10,
		CategoryId:        models.ENTERTAINMENT,
		AppName:           constants.GOOGLE_CHROME,
		AdditionalDataKey: constants.KEY_BROWSER_URL,
		Priority:          1,
		Expression:        "youtube.com"})
}
