package main

import (
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
	rule := ruledtos.ToModel()
	err := a.timekeeper.Storage.Rules().UpsertRule(rule)

	if err != nil {
		a.logger.Error("Error adding rule", "Error", err)
	}

	return err
}

func (a *App) UpdateRule(ruleId int, ruledtos *dtos.RuleUpdate) error {
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
