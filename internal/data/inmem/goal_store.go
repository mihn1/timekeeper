package inmem

import (
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type goalKey struct {
	CategoryId models.CategoryId
	GoalType   models.GoalType
}

type GoalStore struct {
	mu    sync.RWMutex
	goals map[goalKey]*models.CategoryGoal
}

func NewGoalStore() *GoalStore {
	return &GoalStore{
		goals: make(map[goalKey]*models.CategoryGoal),
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

func (s *GoalStore) SetGoal(categoryId models.CategoryId, goalType models.GoalType, targetMs int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := goalKey{CategoryId: categoryId, GoalType: goalType}
	s.goals[key] = &models.CategoryGoal{
		CategoryId: categoryId,
		GoalType:   goalType,
		TargetMs:   targetMs,
		Enabled:    true,
	}
	return nil
}

func (s *GoalStore) DeleteGoal(categoryId models.CategoryId, goalType models.GoalType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.goals, goalKey{CategoryId: categoryId, GoalType: goalType})
	return nil
}
