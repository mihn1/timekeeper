package inmem

import (
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type GoalStore struct {
	mu    sync.RWMutex
	goals map[models.CategoryId]*models.CategoryGoal
}

func NewGoalStore() *GoalStore {
	return &GoalStore{
		goals: make(map[models.CategoryId]*models.CategoryGoal),
	}
}

func (s *GoalStore) GetGoals() ([]*models.CategoryGoal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.CategoryGoal, 0, len(s.goals))
	for _, g := range s.goals {
		result = append(result, g)
	}
	return result, nil
}

func (s *GoalStore) SetGoal(categoryId models.CategoryId, targetMs int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.goals[categoryId] = &models.CategoryGoal{
		CategoryId:    categoryId,
		DailyTargetMs: targetMs,
		Enabled:       true,
	}
	return nil
}

func (s *GoalStore) DeleteGoal(categoryId models.CategoryId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.goals, categoryId)
	return nil
}
