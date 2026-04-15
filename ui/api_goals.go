package main

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/ui/dtos"
)

// GetGoals returns all goals, enriched with category names.
func (a *App) GetGoals() ([]*dtos.GoalItem, error) {
	if a.timekeeper == nil {
		return nil, fmt.Errorf("timekeeper is not initialized")
	}

	goals, err := a.timekeeper.Storage.Goals().GetGoals()
	if err != nil {
		return nil, fmt.Errorf("failed to load goals: %w", err)
	}

	catNames := a.categoryNameMap()

	result := make([]*dtos.GoalItem, 0, len(goals))
	for _, g := range goals {
		catIds := make([]int, len(g.CategoryIds))
		catNamesSlice := make([]string, len(g.CategoryIds))
		for i, cid := range g.CategoryIds {
			catIds[i] = int(cid)
			name := catNames[cid]
			if name == "" {
				name = "Unknown"
			}
			catNamesSlice[i] = name
		}
		result = append(result, &dtos.GoalItem{
			Id:            g.Id,
			Name:          g.Name,
			IsActive:      g.IsActive,
			CategoryIds:   catIds,
			CategoryNames: catNamesSlice,
			Frequency:     int(g.Frequency),
			TargetMs:      g.TargetMs,
		})
	}
	return result, nil
}

// AddGoal creates a new goal and returns its ID.
func (a *App) AddGoal(name string, categoryIds []int, frequency int, targetMs int64) (int64, error) {
	if a.timekeeper == nil {
		return 0, fmt.Errorf("timekeeper is not initialized")
	}

	catIds := make([]models.CategoryId, len(categoryIds))
	for i, id := range categoryIds {
		catIds[i] = models.CategoryId(id)
	}

	goal := &models.CategoryGoal{
		Name:        name,
		IsActive:    true,
		CategoryIds: catIds,
		Frequency:   models.GoalFrequency(frequency),
		TargetMs:    targetMs,
	}

	id, err := a.timekeeper.Storage.Goals().AddGoal(goal)
	if err != nil {
		return 0, fmt.Errorf("failed to add goal: %w", err)
	}
	return id, nil
}

// UpdateGoal updates an existing goal by ID.
func (a *App) UpdateGoal(id int64, name string, categoryIds []int, frequency int, targetMs int64, isActive bool) error {
	if a.timekeeper == nil {
		return fmt.Errorf("timekeeper is not initialized")
	}

	catIds := make([]models.CategoryId, len(categoryIds))
	for i, cid := range categoryIds {
		catIds[i] = models.CategoryId(cid)
	}

	goal := &models.CategoryGoal{
		Id:          id,
		Name:        name,
		IsActive:    isActive,
		CategoryIds: catIds,
		Frequency:   models.GoalFrequency(frequency),
		TargetMs:    targetMs,
	}

	if err := a.timekeeper.Storage.Goals().UpdateGoal(goal); err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}
	return nil
}

// DeleteGoal removes a goal by ID.
func (a *App) DeleteGoal(id int64) error {
	if a.timekeeper == nil {
		return fmt.Errorf("timekeeper is not initialized")
	}
	return a.timekeeper.Storage.Goals().DeleteGoal(id)
}
