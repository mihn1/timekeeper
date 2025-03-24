package resolvers

import (
	"github.com/mihn1/timekeeper/internal/data/interfaces"
	"github.com/mihn1/timekeeper/internal/models"
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

func (r *DefaultCategoryResolver) ResolveCategory(event *models.AppSwitchEvent, rules []*models.CategoryRule) (models.CategoryId, error) {
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
