package sqlite

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryStore struct {
	db        *sql.DB
	tableName string
	mu        sync.Mutex // Add a mutex to protect critical sections
}

func NewCategoryStore(db *sql.DB, tableName string) *CategoryStore {
	store := &CategoryStore{
		db:        db,
		tableName: tableName,
	}

	_, err := store.db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + tableName + ` (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT
        )`)

	if err != nil {
		panic(err)
	}

	return store
}

func (s *CategoryStore) AddCategory(category models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("INSERT INTO "+s.tableName+" (id, name, description) VALUES (?, ?, ?)",
		category.Id, category.Name, category.Description)
	return err
}

func (s *CategoryStore) GetCategory(id models.CategoryId) (models.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow("SELECT id, name, description FROM "+s.tableName+" WHERE id = ?", id)
	var category models.Category
	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, fmt.Errorf("category with id %s not found", id)
		}
		return models.Category{}, err
	}

	return category, nil
}

func (s *CategoryStore) GetCategories() ([]models.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT id, name, description FROM " + s.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
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
