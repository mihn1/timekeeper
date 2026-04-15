package interfaces

import "github.com/mihn1/timekeeper/internal/models"

type GoalStore interface {
	GetGoals() ([]*models.CategoryGoal, error)
	AddGoal(goal *models.CategoryGoal) (int64, error)
	UpdateGoal(goal *models.CategoryGoal) error
	DeleteGoal(id int64) error
}
