package main

import "github.com/mihn1/timekeeper/internal/models"

func (a *App) GetRules() []models.CategoryRule {
	rule, err := a.timekeeper.Storage.Rules().GetRules()

	if err != nil {
		a.logger.Error("Error getting rule", "Error", err)
		return []models.CategoryRule{}
	}

	return rule
}

func (a *App) GetRule(ruleId int) (models.CategoryRule, error) {
	rule, err := a.timekeeper.Storage.Rules().GetRule(ruleId)

	if err != nil {
		a.logger.Error("Error getting rule", "Error", err)
		return models.CategoryRule{}, err
	}

	return rule, nil
}

func (a *App) AddRule(rule models.CategoryRule) error {
	err := a.timekeeper.Storage.Rules().AddRule(rule)

	if err != nil {
		a.logger.Error("Error adding rule", "Error", err)
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
