package data

import (
	"fmt"
	"sync"

	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryStore interface {
	AddCategory(c models.Category) error
	GetCategory(id models.CategoryId) (models.Category, error)
	GetCategories() []models.Category
	DeleteCategory(id models.CategoryId) error
}

type CategoryStore_Memory_Impl struct {
	data map[models.CategoryId]models.Category
	mu   sync.RWMutex
}

func NewCategoryStore_Memory_Impl() *CategoryStore_Memory_Impl {
	return &CategoryStore_Memory_Impl{
		data: make(map[models.CategoryId]models.Category),
		mu:   sync.RWMutex{},
	}
}

func (c *CategoryStore_Memory_Impl) AddCategory(category models.Category) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[category.Id] = category
	return nil
}

func (c *CategoryStore_Memory_Impl) GetCategory(id models.CategoryId) (models.Category, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	category, ok := c.data[id]
	if !ok {
		return models.Category{}, fmt.Errorf("category with id %d not found", id)
	}

	return category, nil
}

func (c *CategoryStore_Memory_Impl) GetCategories() []models.Category {
	c.mu.RLock()
	defer c.mu.RUnlock()

	categories := make([]models.Category, 0, len(c.data))
	for _, category := range c.data {
		categories = append(categories, category)
	}

	return categories
}

func (c *CategoryStore_Memory_Impl) DeleteCategory(id models.CategoryId) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, id)
	return nil
}
