package core

import (
	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/internal/models"
)

func SeedData(t *TimeKeeper) {
	t.logger.Info("Start seeding interfaces...")
	cat, err := t.Storage.Categories().GetCategories()
	if err != nil {
		panic(err)
	}

	if len(cat) > 0 {
		t.logger.Info("Data already seeded.")
		return
	}

	t.Storage.Categories().AddCategory(models.Category{Id: models.WORK, Name: "Work"})
	t.Storage.Categories().AddCategory(models.Category{Id: models.ENTERTAINMENT, Name: "Entertainment"})
	t.Storage.Categories().AddCategory(models.Category{Id: models.PERSONAL, Name: "Personal"})
	t.Storage.Categories().AddCategory(models.Category{Id: models.UNDEFINED, Name: "Undefined"})

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
		t.Storage.Rules().AddRule(rule)
	}

	t.logger.Info("Data seeding completed.")
}
