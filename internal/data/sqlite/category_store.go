package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryStore struct {
	db        *sql.DB
	mu        *sync.RWMutex
	tableName string
}

func NewCategoryStore(db *sql.DB, mu *sync.RWMutex, tableName string) *CategoryStore {
	store := &CategoryStore{
		db:        db,
		mu:        mu,
		tableName: tableName,
	}

	_, err := store.db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + tableName + ` (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            description TEXT
        )`)

	if err != nil {
		panic(err)
	}

	return store
}

func (s *CategoryStore) UpsertCategory(category *models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Trim whitespace from category name
	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}

	var result sql.Result
	var err error

	if category.Id == 0 {
		// Insert new category
		result, err = s.db.Exec("INSERT INTO "+s.tableName+" (name, description) VALUES (?, ?)",
			category.Name, category.Description)

		if err == nil {
			id, _ := result.LastInsertId()
			category.Id = models.CategoryId(id)
		}
	} else {
		// Update existing category
		_, err = s.db.Exec("UPDATE "+s.tableName+" SET name = ?, description = ? WHERE id = ?",
			category.Name, category.Description, category.Id)
	}

	return err
}

func (s *CategoryStore) GetCategory(id models.CategoryId) (*models.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow("SELECT id, name, description FROM "+s.tableName+" WHERE id = ?", id)
	category := &models.Category{}
	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with id %v not found", id)
		}
		return nil, err
	}

	return category, nil
}

func (s *CategoryStore) GetCategories() ([]*models.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT id, name, description FROM " + s.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize with capacity for better performance
	categories := make([]*models.Category, 0)

	for rows.Next() {
		// Create a new Category for each row
		category := &models.Category{}
		err = rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *CategoryStore) DeleteCategory(id models.CategoryId) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM "+s.tableName+" WHERE id = ?", id)
	return err
}
