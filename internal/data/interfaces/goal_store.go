package interfaces

import "github.com/mihn1/timekeeper/internal/models"

type GoalStore interface {
	GetGoals() ([]*models.CategoryGoal, error)
	SetGoal(categoryId models.CategoryId, targetMs int64) error
	DeleteGoal(categoryId models.CategoryId) error
}
