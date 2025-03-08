package resolvers

import (
	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/internal/data/interfaces"
	"github.com/mihn1/timekeeper/internal/models"
	"golang.org/x/exp/slices"
)

type DefaultCategoryResolver struct {
	RuleStore     interfaces.RuleStore
	CategoryStore interfaces.CategoryStore
}

func NewDefaultCategoryResolver(ruleStore interfaces.RuleStore, categoryStore interfaces.CategoryStore) *DefaultCategoryResolver {
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

	// Get rules that are applied to all apps
	globalRules, err := r.RuleStore.GetRulesByApp(constants.ALL_APPS)
	if err == nil {
		rules = append(rules, globalRules...)
	}

	slices.SortStableFunc(rules, models.CmpRules)

	// iterate through the rules to find the first match
	for _, rule := range rules {
		match, err := rule.IsMatch(event)
		if err != nil {
			// we want to keep checking other rules here if there is an error
			continue
		}

		if match {
			return rule.CategoryId, nil
		}
	}

	// if no rules match, return the default category
	return models.UNDEFINED, nil
}
