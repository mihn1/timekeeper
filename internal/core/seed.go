package core

import "github.com/mihn1/timekeeper/internal/models"

func SeedDataInMem(t *TimeKeeper) {
	t.categoryStore.AddCategory(models.Category{Id: models.WORK, Name: "Work"})
	t.categoryStore.AddCategory(models.Category{Id: models.ENTERTAINMENT, Name: "Entertainment"})
	t.categoryStore.AddCategory(models.Category{Id: models.UNDEFINED, Name: "Undefined"})

	t.ruleStore.AddRule(models.CategoryRule{RuleId: 1, CategoryId: models.ENTERTAINMENT, AppName: "Google Chrome"})
	t.ruleStore.AddRule(models.CategoryRule{RuleId: 2, CategoryId: models.WORK, AppName: "Code"})
	t.ruleStore.AddRule(models.CategoryRule{RuleId: 3, CategoryId: models.WORK, AppName: "Ghostty"})
	t.ruleStore.AddRule(models.CategoryRule{RuleId: 4, CategoryId: models.WORK, AppName: "ChatGPT"})
}
