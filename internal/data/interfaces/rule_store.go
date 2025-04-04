package interfaces

import (
	"github.com/mihn1/timekeeper/internal/models"
)

type RuleStore interface {
	GetRules() ([]*models.CategoryRule, error)
	GetRule(ruleId int) (*models.CategoryRule, error)
	GetRulesByCategory(categoryId models.CategoryId) ([]*models.CategoryRule, error)
	GetRulesByApp(appName string) ([]*models.CategoryRule, error)
	UpsertRule(rule *models.CategoryRule) error
	DeleteRule(ruleId int) error
}
