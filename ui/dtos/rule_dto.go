package dtos

import (
	"github.com/mihn1/timekeeper/internal/models"
)

type RuleListItem struct {
	ID                int    `json:"id"`
	CategoryID        int    `json:"categoryId"`
	AppName           string `json:"appName"`
	AdditionalDataKey string `json:"additionalDataKey"`
	Expression        string `json:"expression"`
	IsRegex           bool   `json:"isRegex"`
	Priority          int    `json:"priority"`
}

type RuleDetail struct {
	ID                int    `json:"id"`
	CategoryID        int    `json:"categoryId"`
	AppName           string `json:"appName"`
	AdditionalDataKey string `json:"additionalDataKey"`
	Expression        string `json:"expression"`
	IsRegex           bool   `json:"isRegex"`
	Priority          int    `json:"priority"`
}

type RuleCreate struct {
	CategoryID        int    `json:"categoryId"`
	AppName           string `json:"appName"`
	AdditionalDataKey string `json:"additionalDataKey"`
	Expression        string `json:"expression"`
	IsRegex           bool   `json:"isRegex"`
	Priority          int    `json:"priority"`
}

type RuleUpdate struct {
	ID                int    `json:"id"`
	CategoryID        int    `json:"categoryId"`
	AppName           string `json:"appName"`
	AdditionalDataKey string `json:"additionalDataKey"`
	Expression        string `json:"expression"`
	IsRegex           bool   `json:"isRegex"`
	Priority          int    `json:"priority"`
}

func (rc *RuleCreate) ToModel() *models.CategoryRule {
	return &models.CategoryRule{
		CategoryId:        models.CategoryId(rc.CategoryID),
		AppName:           rc.AppName,
		AdditionalDataKey: rc.AdditionalDataKey,
		Expression:        rc.Expression,
		IsRegex:           rc.IsRegex,
		Priority:          rc.Priority,
	}
}

func (ru *RuleUpdate) ToModel() *models.CategoryRule {
	return &models.CategoryRule{
		RuleId:            ru.ID,
		CategoryId:        models.CategoryId(ru.CategoryID),
		AppName:           ru.AppName,
		AdditionalDataKey: ru.AdditionalDataKey,
		Expression:        ru.Expression,
		IsRegex:           ru.IsRegex,
		Priority:          ru.Priority,
	}
}

func RuleDetailFromModel(rule *models.CategoryRule) *RuleDetail {
	return &RuleDetail{
		ID:                rule.RuleId,
		CategoryID:        int(rule.CategoryId),
		AppName:           rule.AppName,
		AdditionalDataKey: rule.AdditionalDataKey,
		Expression:        rule.Expression,
		IsRegex:           rule.IsRegex,
		Priority:          rule.Priority,
	}
}

func RuleListFromModels(rules []*models.CategoryRule) []*RuleListItem {
	result := make([]*RuleListItem, len(rules))
	for i, rule := range rules {
		result[i] = &RuleListItem{
			ID:                rule.RuleId,
			CategoryID:        int(rule.CategoryId),
			AppName:           rule.AppName,
			Expression:        rule.Expression,
			IsRegex:           rule.IsRegex,
			AdditionalDataKey: rule.AdditionalDataKey,
			Priority:          rule.Priority,
		}
	}
	return result
}
