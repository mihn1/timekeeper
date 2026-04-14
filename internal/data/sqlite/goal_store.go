package sqlite

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type GoalStore struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewGoalStore(db *sql.DB, mu *sync.RWMutex, tableName string) *GoalStore {
	store := &GoalStore{db: db, mu: mu, tableName: tableName}
	store.initSchema()
	return store
}

func (s *GoalStore) initSchema() {
	// Create with new schema if table doesn't exist.
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS ` + s.tableName + ` (
		category_id INTEGER NOT NULL,
		goal_type   TEXT    NOT NULL DEFAULT 'daily',
		target_ms   INTEGER NOT NULL,
		enabled     INTEGER NOT NULL DEFAULT 1,
		PRIMARY KEY (category_id, goal_type)
	)`)
	if err != nil {
		panic(fmt.Sprintf("failed to create %s table: %v", s.tableName, err))
	}

	// Migrate old schema if goal_type column is missing.
	rows, err := s.db.Query(`PRAGMA table_info(` + s.tableName + `)`)
	if err != nil {
		return
	}
	defer rows.Close()

	hasGoalType := false
	hasDailyTargetMs := false
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			continue
		}
		if name == "goal_type" {
			hasGoalType = true
		}
		if name == "daily_target_ms" {
			hasDailyTargetMs = true
		}
	}
	rows.Close()

	if hasGoalType {
		return // Already on new schema.
	}

	// Migrate: old schema used daily_target_ms + category_id as sole PK.
	newTable := s.tableName + "_new"
	_, err = s.db.Exec(`CREATE TABLE ` + newTable + ` (
		category_id INTEGER NOT NULL,
		goal_type   TEXT    NOT NULL DEFAULT 'daily',
		target_ms   INTEGER NOT NULL,
		enabled     INTEGER NOT NULL DEFAULT 1,
		PRIMARY KEY (category_id, goal_type)
	)`)
	if err != nil {
		return
	}

	if hasDailyTargetMs {
		s.db.Exec(`INSERT INTO ` + newTable + ` (category_id, goal_type, target_ms, enabled)
			SELECT category_id, 'daily', daily_target_ms, enabled FROM ` + s.tableName)
	}

	s.db.Exec(`DROP TABLE ` + s.tableName)
	s.db.Exec(`ALTER TABLE ` + newTable + ` RENAME TO ` + s.tableName)
}

func (s *GoalStore) GetGoals() ([]*models.CategoryGoal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(`SELECT category_id, goal_type, target_ms, enabled FROM ` + s.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*models.CategoryGoal
	for rows.Next() {
		goal := &models.CategoryGoal{}
		var enabled int
		if err := rows.Scan(&goal.CategoryId, &goal.GoalType, &goal.TargetMs, &enabled); err != nil {
			return nil, err
		}
		goal.Enabled = enabled != 0
		goals = append(goals, goal)
	}
	return goals, nil
}

func (s *GoalStore) SetGoal(categoryId models.CategoryId, goalType models.GoalType, targetMs int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		INSERT INTO `+s.tableName+` (category_id, goal_type, target_ms, enabled)
		VALUES (?, ?, ?, 1)
		ON CONFLICT(category_id, goal_type) DO UPDATE SET target_ms = excluded.target_ms, enabled = 1`,
		int(categoryId), string(goalType), targetMs)
	return err
}

func (s *GoalStore) DeleteGoal(categoryId models.CategoryId, goalType models.GoalType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM `+s.tableName+` WHERE category_id = ? AND goal_type = ?`,
		int(categoryId), string(goalType))
	return err
}
