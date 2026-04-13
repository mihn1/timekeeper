package sqlite

import (
	"database/sql"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type GoalStore struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewGoalStore(db *sql.DB, mu *sync.RWMutex, tableName string) *GoalStore {
	store := &GoalStore{
		db:        db,
		mu:        mu,
		tableName: tableName,
	}

	_, err := store.db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			category_id INTEGER PRIMARY KEY,
			daily_target_ms INTEGER NOT NULL,
			enabled INTEGER NOT NULL DEFAULT 1
		)`)
	if err != nil {
		panic(err)
	}

	return store
}

func (s *GoalStore) GetGoals() ([]*models.CategoryGoal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT category_id, daily_target_ms, enabled FROM " + s.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*models.CategoryGoal
	for rows.Next() {
		goal := &models.CategoryGoal{}
		var enabled int
		err = rows.Scan(&goal.CategoryId, &goal.DailyTargetMs, &enabled)
		if err != nil {
			return nil, err
		}
		goal.Enabled = enabled != 0
		goals = append(goals, goal)
	}
	return goals, nil
}

func (s *GoalStore) SetGoal(categoryId models.CategoryId, targetMs int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		INSERT INTO `+s.tableName+` (category_id, daily_target_ms, enabled)
		VALUES (?, ?, 1)
		ON CONFLICT(category_id) DO UPDATE SET daily_target_ms = excluded.daily_target_ms, enabled = 1`,
		int(categoryId), targetMs)
	return err
}

func (s *GoalStore) DeleteGoal(categoryId models.CategoryId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM "+s.tableName+" WHERE category_id = ?", int(categoryId))
	return err
}
