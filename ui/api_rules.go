package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/core/resolvers"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/ui/dtos"
)

func (a *App) GetRules() []*dtos.RuleListItem {
	rules, err := a.timekeeper.Storage.Rules().GetRules()

	if err != nil {
		a.logger.Error("Error getting rules", "Error", err)
		return nil
	}

	return dtos.RuleListFromModels(rules)
}

func (a *App) GetRule(ruleId int) (*dtos.RuleDetail, error) {
	rule, err := a.timekeeper.Storage.Rules().GetRule(ruleId)

	if err != nil {
		a.logger.Error("Error getting rule", "Error", err)
		return nil, err
	}

	return dtos.RuleDetailFromModel(rule), nil
}

func (a *App) AddRule(ruledtos *dtos.RuleCreate) error {
	if ruledtos == nil {
		return fmt.Errorf("rule payload is required")
	}

	if strings.TrimSpace(ruledtos.AppName) == "" {
		return fmt.Errorf("rule app name is required")
	}

	if ruledtos.CategoryID <= 0 && !ruledtos.IsExclusion {
		return fmt.Errorf("rule categoryId must be a positive integer")
	}

	rule := ruledtos.ToModel()
	err := a.timekeeper.Storage.Rules().UpsertRule(rule)

	if err != nil {
		a.logger.Error("Error adding rule", "Error", err)
	}

	return err
}

func (a *App) UpdateRule(ruleId int, ruledtos *dtos.RuleUpdate) error {
	if ruledtos == nil {
		return fmt.Errorf("rule payload is required")
	}

	if strings.TrimSpace(ruledtos.AppName) == "" {
		return fmt.Errorf("rule app name is required")
	}

	if ruledtos.CategoryID <= 0 && !ruledtos.IsExclusion {
		return fmt.Errorf("rule categoryId must be a positive integer")
	}

	// First get the existing rule
	existingRule, err := a.timekeeper.Storage.Rules().GetRule(ruleId)
	if err != nil {
		a.logger.Error("Error getting rule for update", "Error", err)
		return err
	}

	// Update fields
	existingRule.CategoryId = models.CategoryId(ruledtos.CategoryID)
	existingRule.AppName = ruledtos.AppName
	existingRule.AdditionalDataKey = ruledtos.AdditionalDataKey
	existingRule.Expression = ruledtos.Expression
	existingRule.IsRegex = ruledtos.IsRegex
	existingRule.Priority = ruledtos.Priority
	existingRule.IsExclusion = ruledtos.IsExclusion

	// Save the updated rule
	err = a.timekeeper.Storage.Rules().UpsertRule(existingRule)
	if err != nil {
		a.logger.Error("Error updating rule", "Error", err)
	}

	return err
}

func (a *App) DeleteRule(ruleId int) error {
	err := a.timekeeper.Storage.Rules().DeleteRule(ruleId)

	if err != nil {
		a.logger.Error("Error deleting rule", "Error", err)
	}

	return err
}

// TestRuleMatch simulates rule resolution for a given app name and optional additional data.
// additionalDataKey and value correspond to the key and value in the event's AdditionalData map
// (e.g. key="url", value="https://github.com/...").
func (a *App) TestRuleMatch(appName, additionalDataKey, value string) (*dtos.RuleMatchResult, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	event := &models.AppSwitchEvent{AppName: appName}
	if additionalDataKey != "" && value != "" {
		event.AdditionalData = map[string]string{additionalDataKey: value}
	}

	appRules, _ := a.timekeeper.Storage.Rules().GetRulesByApp(appName)
	globalRules, _ := a.timekeeper.Storage.Rules().GetRulesByApp(constants.ALL_APPS)
	allRules := append(appRules, globalRules...)
	sort.Slice(allRules, func(i, j int) bool {
		return allRules[i].Priority > allRules[j].Priority
	})

	resolver := resolvers.NewDefaultCategoryResolver(
		a.timekeeper.Storage.Rules(),
		a.timekeeper.Storage.Categories(),
	)

	catId, err := resolver.ResolveCategory(event, allRules)
	if err != nil {
		return nil, err
	}

	var matchedRule *dtos.RuleDetail
	for _, rule := range allRules {
		match, err := rule.IsMatch(event)
		if err == nil && match {
			matchedRule = dtos.RuleDetailFromModel(rule)
			break
		}
	}

	catName := "Undefined"
	if cat, err := a.timekeeper.Storage.Categories().GetCategory(catId); err == nil {
		catName = cat.Name
	}

	return &dtos.RuleMatchResult{
		Matched:      matchedRule != nil,
		CategoryId:   int(catId),
		CategoryName: catName,
		MatchedRule:  matchedRule,
	}, nil
}
