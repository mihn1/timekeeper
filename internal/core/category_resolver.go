package core

import (
	"github.com/mihn1/timekeeper/internal/data"
	"github.com/mihn1/timekeeper/internal/models"
)

// resolve category from an app switch event
type CategoryResolver interface {
	ResolveCategory(event *models.AppSwitchEvent) (models.Category, error)
}

type DefaultCategoryResolver struct {
	RuleStore     data.RuleStore
	CategoryStore data.CategoryStore
}

func NewDefaultCategoryResolver(ruleStore data.RuleStore, categoryStore data.CategoryStore) *DefaultCategoryResolver {
	return &DefaultCategoryResolver{
		RuleStore:     ruleStore,
		CategoryStore: categoryStore,
	}
}

func (r *DefaultCategoryResolver) ResolveCategory(event *models.AppSwitchEvent) (models.Category, error) {
	// get the rules for the app
	rules, err := r.RuleStore.GetRulesByApp(event.AppName)
	if err != nil {
		return models.Category{}, err
	}

	// if there are no rules, return the default category
	if len(rules) == 0 {
		return r.CategoryStore.GetCategory(models.UNDEFINED)
	}

	// iterate through the rules to find the first match
	for _, rule := range rules {
		if rule.IsMatch(event) {
			return r.CategoryStore.GetCategory(rule.CategoryId)
		}
	}

	// if no rules match, return the default category
	return r.CategoryStore.GetCategory(models.UNDEFINED)
}
