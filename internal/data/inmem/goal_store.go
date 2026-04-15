package inmem

import (
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type GoalStore struct {
	mu     sync.RWMutex
	goals  map[int64]*models.CategoryGoal
	nextId int64
}

func NewGoalStore() *GoalStore {
	return &GoalStore{
		goals:  make(map[int64]*models.CategoryGoal),
		nextId: 1,
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

func (s *GoalStore) AddGoal(goal *models.CategoryGoal) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextId
	s.nextId++
	goal.Id = id
	// Store a copy.
	cp := *goal
	ids := make([]models.CategoryId, len(goal.CategoryIds))
	copy(ids, goal.CategoryIds)
	cp.CategoryIds = ids
	s.goals[id] = &cp
	return id, nil
}

func (s *GoalStore) UpdateGoal(goal *models.CategoryGoal) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.goals[goal.Id]; !ok {
		return fmt.Errorf("goal %d not found", goal.Id)
	}
	cp := *goal
	ids := make([]models.CategoryId, len(goal.CategoryIds))
	copy(ids, goal.CategoryIds)
	cp.CategoryIds = ids
	s.goals[goal.Id] = &cp
	return nil
}

func (s *GoalStore) DeleteGoal(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.goals, id)
	return nil
}
