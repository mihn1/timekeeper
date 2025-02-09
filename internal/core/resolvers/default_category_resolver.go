package resolvers

import (
	"github.com/mihn1/timekeeper/internal/data"
	"github.com/mihn1/timekeeper/internal/models"
	"golang.org/x/exp/slices"
)

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

func (r *DefaultCategoryResolver) ResolveCategory(event *models.AppSwitchEvent) (models.CategoryId, error) {
	rules, err := r.RuleStore.GetRulesByApp(event.AppName)
	if err != nil {
		return models.UNDEFINED, err
	}

	slices.SortStableFunc(rules, models.CmpRules)

	// iterate through the rules to find the first match
	for _, rule := range rules {
		match, err := rule.IsMatch(event)
		if err != nil {
			return models.UNDEFINED, err
		}

		if match {
			return rule.CategoryId, nil
		}
	}

	// if no rules match, return the default category
	return models.UNDEFINED, nil
}
