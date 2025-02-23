package sqlite

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type RuleStore struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewRuleStore(db *sql.DB, mu *sync.RWMutex, tableName string) *RuleStore {
	store := &RuleStore{
		db:        db,
		mu:        mu,
		tableName: tableName,
	}

	_, err := store.db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + tableName + ` (
            rule_id INTEGER PRIMARY KEY,
            category_id TEXT NOT NULL,
            app_name TEXT NOT NULL,
            additional_data_key TEXT,
            expression TEXT NOT NULL,
            is_regex BOOLEAN NOT NULL,
            priority INTEGER NOT NULL
        )`)

	if err != nil {
		panic(err)
	}

	return store
}

func (s *RuleStore) GetRules() ([]models.CategoryRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM " + s.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.CategoryRule
	for rows.Next() {
		var rule models.CategoryRule
		err = rows.Scan(&rule.RuleId, &rule.CategoryId, &rule.AppName, &rule.AdditionalDataKey, &rule.Expression, &rule.IsRegex, &rule.Priority)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (s *RuleStore) GetRule(ruleId int) (models.CategoryRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM "+s.tableName+" WHERE rule_id = ?", ruleId)
	var rule models.CategoryRule
	err := row.Scan(&rule.RuleId, &rule.CategoryId, &rule.AppName, &rule.AdditionalDataKey, &rule.Expression, &rule.IsRegex, &rule.Priority)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.CategoryRule{}, fmt.Errorf("rule with id %d not found", ruleId)
		}
		return models.CategoryRule{}, err
	}

	return rule, nil
}

func (s *RuleStore) GetRulesByCategory(categoryId models.CategoryId) ([]models.CategoryRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM "+s.tableName+" WHERE category_id = ?", categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.CategoryRule
	for rows.Next() {
		var rule models.CategoryRule
		err = rows.Scan(&rule.RuleId, &rule.CategoryId, &rule.AppName, &rule.AdditionalDataKey, &rule.Expression, &rule.IsRegex, &rule.Priority)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (s *RuleStore) GetRulesByApp(appName string) ([]models.CategoryRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM "+s.tableName+" WHERE app_name = ?", appName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.CategoryRule
	for rows.Next() {
		var rule models.CategoryRule
		err = rows.Scan(&rule.RuleId, &rule.CategoryId, &rule.AppName, &rule.AdditionalDataKey, &rule.Expression, &rule.IsRegex, &rule.Priority)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (s *RuleStore) AddRule(rule models.CategoryRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("INSERT INTO "+s.tableName+" (rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority) VALUES (?, ?, ?, ?, ?, ?, ?)",
		rule.RuleId, rule.CategoryId, rule.AppName, rule.AdditionalDataKey, rule.Expression, rule.IsRegex, rule.Priority)
	return err
}

func (s *RuleStore) DeleteRule(ruleId int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM "+s.tableName+" WHERE rule_id = ?", ruleId)
	return err
}
