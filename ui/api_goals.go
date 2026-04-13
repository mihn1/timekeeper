package main

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/ui/dtos"
)

// GetGoals returns all category goals, enriched with category names.
func (a *App) GetGoals() ([]*dtos.GoalItem, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	goals, err := a.timekeeper.Storage.Goals().GetGoals()
	if err != nil {
		return nil, fmt.Errorf("failed to load goals: %w", err)
	}

	categoryNameMap := make(map[models.CategoryId]string)
	if cats, err := a.timekeeper.Storage.Categories().GetCategories(); err == nil {
		for _, cat := range cats {
			categoryNameMap[cat.Id] = cat.Name
		}
	}

	result := make([]*dtos.GoalItem, 0, len(goals))
	for _, g := range goals {
		catName := categoryNameMap[g.CategoryId]
		if catName == "" {
			catName = "Unknown"
		}
		result = append(result, &dtos.GoalItem{
			CategoryId:    int(g.CategoryId),
			CategoryName:  catName,
			DailyTargetMs: g.DailyTargetMs,
			Enabled:       g.Enabled,
		})
	}
	return result, nil
}

// SetGoal creates or updates the daily time target for a category.
func (a *App) SetGoal(categoryId int, targetMs int64) error {
	if a.timekeeper == nil {
		return fmt.Errorf("timekeeper is not initialized")
	}
	return a.timekeeper.Storage.Goals().SetGoal(models.CategoryId(categoryId), targetMs)
}

// DeleteGoal removes the daily time target for a category.
func (a *App) DeleteGoal(categoryId int) error {
	if a.timekeeper == nil {
		return fmt.Errorf("timekeeper is not initialized")
	}
	return a.timekeeper.Storage.Goals().DeleteGoal(models.CategoryId(categoryId))
}
