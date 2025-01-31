package data

import "github.com/mihn1/timekeeper/internal/models"

type CategoryStore interface {
	AddCategory(c models.Category) error
	GetCategory(id models.CategoryId) (models.Category, error)
	GetCategories() []models.Category
	DeleteCategory(id models.CategoryId) error
}
