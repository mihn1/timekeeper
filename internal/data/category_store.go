package data

import (
	"fmt"
	"sync"
)

type CategoryStore interface {
	AddCategory(c Category) error
	GetCategory(id CategoryId) (Category, error)
	GetCategories() []Category
	DeleteCategory(id CategoryId) error
}

type CategoryStore_Memory_Impl struct {
	data map[CategoryId]Category
	mu   sync.RWMutex
}

func NewCategoryStore_Memory_Impl() *CategoryStore_Memory_Impl {
	return &CategoryStore_Memory_Impl{
		data: make(map[CategoryId]Category),
		mu:   sync.RWMutex{},
	}
}

func (c *CategoryStore_Memory_Impl) AddCategory(category Category) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[category.Id] = category
	return nil
}

func (c *CategoryStore_Memory_Impl) GetCategory(id CategoryId) (Category, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	category, ok := c.data[id]
	if !ok {
		return Category{}, fmt.Errorf("category with id %d not found", id)
	}

	return category, nil
}

func (c *CategoryStore_Memory_Impl) GetCategories() []Category {
	c.mu.RLock()
	defer c.mu.RUnlock()

	categories := make([]Category, 0, len(c.data))
	for _, category := range c.data {
		categories = append(categories, category)
	}

	return categories
}

func (c *CategoryStore_Memory_Impl) DeleteCategory(id CategoryId) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, id)
	return nil
}
