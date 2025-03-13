package interfaces

import "github.com/mihn1/timekeeper/internal/models"

type CategoryStore interface {
	UpsertCategory(c *models.Category) error
	GetCategory(id models.CategoryId) (*models.Category, error)
	GetCategories() ([]*models.Category, error)
	DeleteCategory(id models.CategoryId) error
}
