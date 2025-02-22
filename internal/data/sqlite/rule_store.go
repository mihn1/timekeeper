package sqlite

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type RuleStore struct {
	db        *sql.DB
	tableName string
	mu        sync.Mutex // Add a mutex to protect critical sections
}

func NewRuleStore(db *sql.DB, tableName string) *RuleStore {
	store := &RuleStore{
		db:        db,
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

func (r *RuleStore) GetRules() ([]models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM " + r.tableName)
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

func (r *RuleStore) GetRule(ruleId int) (models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	row := r.db.QueryRow("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM "+r.tableName+" WHERE rule_id = ?", ruleId)
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

func (r *RuleStore) GetRulesByCategory(categoryId models.CategoryId) ([]models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM "+r.tableName+" WHERE category_id = ?", categoryId)
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

func (r *RuleStore) GetRulesByApp(appName string) ([]models.CategoryRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query("SELECT rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority FROM "+r.tableName+" WHERE app_name = ?", appName)
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

func (r *RuleStore) AddRule(rule models.CategoryRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec("INSERT INTO "+r.tableName+" (rule_id, category_id, app_name, additional_data_key, expression, is_regex, priority) VALUES (?, ?, ?, ?, ?, ?, ?)",
		rule.RuleId, rule.CategoryId, rule.AppName, rule.AdditionalDataKey, rule.Expression, rule.IsRegex, rule.Priority)
	return err
}

func (r *RuleStore) DeleteRule(ruleId int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec("DELETE FROM "+r.tableName+" WHERE rule_id = ?", ruleId)
	return err
}
