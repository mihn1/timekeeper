package inmem

import (
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type RuleStore struct {
	data map[int]*models.CategoryRule // Changed to store pointers
	mu   sync.Mutex                   // Added mutex for thread safety
}

func NewRuleStore() *RuleStore {
	return &RuleStore{
		data: make(map[int]*models.CategoryRule), // Changed to store pointers
		mu:   sync.Mutex{},
	}
}

func (r *RuleStore) GetRules() ([]*models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rules := make([]*models.CategoryRule, 0, len(r.data))
	for _, rule := range r.data {
		// Create a deep copy to prevent external mutations
		copiedRule := &models.CategoryRule{
			RuleId:            rule.RuleId,
			CategoryId:        rule.CategoryId,
			AppName:           rule.AppName,
			AdditionalDataKey: rule.AdditionalDataKey,
			Expression:        rule.Expression,
			IsRegex:           rule.IsRegex,
			Priority:          rule.Priority,
		}
		rules = append(rules, copiedRule)
	}
	return rules, nil
}

func (r *RuleStore) GetRule(ruleId int) (*models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, rule := range r.data {
		if rule.RuleId == ruleId {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("rule with id %d not found", ruleId)
}

func (r *RuleStore) GetRulesByCategory(categoryId models.CategoryId) ([]*models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]*models.CategoryRule, 0)
	for _, rule := range r.data {
		if rule.CategoryId == categoryId {
			res = append(res, rule)
		}
	}
	return res, nil
}

func (r *RuleStore) GetRulesByApp(appName string) ([]*models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]*models.CategoryRule, 0)
	for _, rule := range r.data {
		if rule.AppName == appName {
			res = append(res, rule)
		}
	}
	return res, nil
}

func (r *RuleStore) UpsertRule(rule *models.CategoryRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if rule.RuleId == 0 {
		// Generate a new rule id
		maxId := 0
		for _, r := range r.data {
			if r.RuleId > maxId {
				maxId = r.RuleId
			}
		}
		rule.RuleId = maxId + 1
	}

	r.data[rule.RuleId] = rule
	return nil
}

func (r *RuleStore) DeleteRule(ruleId int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[ruleId]; !ok {
		return fmt.Errorf("rule with id %d not found", ruleId)
	}

	delete(r.data, ruleId)
	return nil
}
