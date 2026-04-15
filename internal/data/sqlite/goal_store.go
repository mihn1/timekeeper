package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		name         TEXT    NOT NULL DEFAULT '',
		is_active    INTEGER NOT NULL DEFAULT 1,
		category_ids TEXT    NOT NULL DEFAULT '',
		frequency    INTEGER NOT NULL DEFAULT 1,
		target_ms    INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		panic(fmt.Sprintf("failed to create %s table: %v", s.tableName, err))
	}

	// Detect old schema by checking for presence of "name" column.
	rows, err := s.db.Query(`PRAGMA table_info(` + s.tableName + `)`)
	if err != nil {
		return
	}
	defer rows.Close()

	hasName := false
	hasOldCategoryId := false
	for rows.Next() {
		var cid int
		var colName, colType string
		var notNull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &colName, &colType, &notNull, &dflt, &pk); err != nil {
			continue
		}
		if colName == "name" {
			hasName = true
		}
		if colName == "category_id" {
			hasOldCategoryId = true
		}
	}
	rows.Close()

	if hasName {
		return // Already on new schema.
	}

	if !hasOldCategoryId {
		return // Empty or unknown table, no migration needed.
	}

	// Migrate old schema: had (category_id, goal_type TEXT, target_ms, enabled).
	newTable := s.tableName + "_new"
	_, err = s.db.Exec(`CREATE TABLE ` + newTable + ` (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		name         TEXT    NOT NULL DEFAULT '',
		is_active    INTEGER NOT NULL DEFAULT 1,
		category_ids TEXT    NOT NULL DEFAULT '',
		frequency    INTEGER NOT NULL DEFAULT 1,
		target_ms    INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		return
	}

	// Map goal_type string → frequency int.
	s.db.Exec(`INSERT INTO ` + newTable + ` (name, is_active, category_ids, frequency, target_ms)
		SELECT
			'',
			COALESCE(enabled, 1),
			CAST(category_id AS TEXT),
			CASE goal_type
				WHEN 'weekly'  THEN 2
				WHEN 'monthly' THEN 3
				ELSE 1
			END,
			target_ms
		FROM ` + s.tableName)

	s.db.Exec(`DROP TABLE ` + s.tableName)
	s.db.Exec(`ALTER TABLE ` + newTable + ` RENAME TO ` + s.tableName)
}

func parseCategoryIds(s string) []models.CategoryId {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]models.CategoryId, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err == nil {
			result = append(result, models.CategoryId(n))
		}
	}
	return result
}

func serializeCategoryIds(ids []models.CategoryId) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.FormatInt(int64(id), 10)
	}
	return strings.Join(parts, ",")
}

func (s *GoalStore) GetGoals() ([]*models.CategoryGoal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(`SELECT id, name, is_active, category_ids, frequency, target_ms FROM ` + s.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*models.CategoryGoal
	for rows.Next() {
		goal := &models.CategoryGoal{}
		var isActive int
		var categoryIdsStr string
		if err := rows.Scan(&goal.Id, &goal.Name, &isActive, &categoryIdsStr, &goal.Frequency, &goal.TargetMs); err != nil {
			return nil, err
		}
		goal.IsActive = isActive != 0
		goal.CategoryIds = parseCategoryIds(categoryIdsStr)
		goals = append(goals, goal)
	}
	return goals, nil
}

func (s *GoalStore) AddGoal(goal *models.CategoryGoal) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(
		`INSERT INTO `+s.tableName+` (name, is_active, category_ids, frequency, target_ms) VALUES (?, ?, ?, ?, ?)`,
		goal.Name,
		boolToInt(goal.IsActive),
		serializeCategoryIds(goal.CategoryIds),
		int(goal.Frequency),
		goal.TargetMs,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	goal.Id = id
	return id, nil
}

func (s *GoalStore) UpdateGoal(goal *models.CategoryGoal) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`UPDATE `+s.tableName+` SET name=?, is_active=?, category_ids=?, frequency=?, target_ms=? WHERE id=?`,
		goal.Name,
		boolToInt(goal.IsActive),
		serializeCategoryIds(goal.CategoryIds),
		int(goal.Frequency),
		goal.TargetMs,
		goal.Id,
	)
	return err
}

func (s *GoalStore) DeleteGoal(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM `+s.tableName+` WHERE id=?`, id)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
