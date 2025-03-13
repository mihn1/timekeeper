package inmem

import (
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryStore struct {
	data map[models.CategoryId]*models.Category // Change to store pointers
	mu   sync.Mutex
}

func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		data: make(map[models.CategoryId]*models.Category), // Update map type
		mu:   sync.Mutex{},
	}
}

func (c *CategoryStore) UpsertCategory(category *models.Category) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if category.Id == 0 {
		// Generate a new ID
		var maxId models.CategoryId = 0
		for id := range c.data {
			if id > maxId {
				maxId = id
			}
		}
		category.Id = maxId + 1
	}

	c.data[category.Id] = category
	return nil
}

func (c *CategoryStore) GetCategory(id models.CategoryId) (*models.Category, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	category, ok := c.data[id]
	if !ok {
		return nil, fmt.Errorf("category with id %v not found", id)
	}

	// Return a copy to prevent external mutations
	return &models.Category{
		Id:          category.Id,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c *CategoryStore) GetCategories() ([]*models.Category, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	categories := make([]*models.Category, 0, len(c.data))
	for _, category := range c.data {
		// Create a copy to avoid external mutations
		copiedCategory := &models.Category{
			Id:          category.Id,
			Name:        category.Name,
			Description: category.Description,
		}
		categories = append(categories, copiedCategory)
	}

	return categories, nil
}

func (c *CategoryStore) DeleteCategory(id models.CategoryId) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, id)
	return nil
}
