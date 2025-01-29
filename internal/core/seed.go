package core

import "github.com/mihn1/timekeeper/internal/data"

func SeedDataInMem(t *TimeKeeper) {
	t.categoryStore.AddCategory(data.Category{Id: data.WORK, Name: "Work"})
	t.categoryStore.AddCategory(data.Category{Id: data.ENTERTAINMENT, Name: "Entertainment"})
	t.categoryStore.AddCategory(data.Category{Id: data.UNDEFINED, Name: "Undefined"})

	t.ruleStore.AddRule(data.CategoryRule{RuleId: 1, CategoryId: data.ENTERTAINMENT, AppName: "Google Chrome"})
	t.ruleStore.AddRule(data.CategoryRule{RuleId: 2, CategoryId: data.WORK, AppName: "Code"})
	t.ruleStore.AddRule(data.CategoryRule{RuleId: 3, CategoryId: data.WORK, AppName: "Ghostty"})
	t.ruleStore.AddRule(data.CategoryRule{RuleId: 4, CategoryId: data.WORK, AppName: "ChatGPT"})
}
