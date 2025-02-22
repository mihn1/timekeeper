package inmem

import (
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryStore struct {
	data map[models.CategoryId]models.Category
	mu   sync.Mutex
}

func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		data: make(map[models.CategoryId]models.Category),
		mu:   sync.Mutex{},
	}
}

func (c *CategoryStore) AddCategory(category models.Category) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[category.Id] = category
	return nil
}

func (c *CategoryStore) GetCategory(id models.CategoryId) (models.Category, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	category, ok := c.data[id]
	if !ok {
		return models.Category{}, fmt.Errorf("category with id %s not found", id)
	}

	return category, nil
}

func (c *CategoryStore) GetCategories() ([]models.Category, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	categories := make([]models.Category, 0, len(c.data))
	for _, category := range c.data {
		categories = append(categories, category)
	}

	return categories, nil
}

func (c *CategoryStore) DeleteCategory(id models.CategoryId) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, id)
	return nil
}
