package inmem

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/utils"
)

type RuleStore struct {
	data map[int]models.CategoryRule
}

func NewRuleStore() *RuleStore {
	return &RuleStore{
		data: make(map[int]models.CategoryRule),
	}
}

func (r *RuleStore) GetRules() ([]models.CategoryRule, error) {
	return utils.GetMapValues(r.data), nil
}

func (r *RuleStore) GetRule(ruleId int) (models.CategoryRule, error) {
	for _, rule := range r.data {
		if rule.RuleId == ruleId {
			return rule, nil
		}
	}
	return models.CategoryRule{}, fmt.Errorf("rule with id %d not found", ruleId)
}

func (r *RuleStore) GetRulesByCategory(categoryId models.CategoryId) ([]models.CategoryRule, error) {
	res := make([]models.CategoryRule, 0)
	for _, rule := range r.data {
		if rule.CategoryId == categoryId {
			res = append(res, rule)
		}
	}
	return res, nil
}

func (r *RuleStore) GetRulesByApp(appName string) ([]models.CategoryRule, error) {
	res := make([]models.CategoryRule, 0)
	for _, rule := range r.data {
		if rule.AppName == appName {
			res = append(res, rule)
		}
	}
	return res, nil
}

func (r *RuleStore) AddRule(rule models.CategoryRule) error {
	if _, ok := r.data[rule.RuleId]; ok {
		return fmt.Errorf("rule with id %d already exists", rule.RuleId)
	}

	r.data[rule.RuleId] = rule
	return nil
}

func (r *RuleStore) DeleteRule(ruleId int) error {
	if _, ok := r.data[ruleId]; !ok {
		return fmt.Errorf("rule with id %d not found", ruleId)
	}

	delete(r.data, ruleId)
	return nil
}
