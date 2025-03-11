package main

import "github.com/mihn1/timekeeper/internal/models"

// GetCategories returns all categories
func (a *App) GetCategories() ([]models.Category, error) {
	categories, err := a.timekeeper.Storage.Categories().GetCategories()
	if err != nil {
		a.logger.Error("Error getting categories", "Error", err)
		return nil, err
	}
	return categories, nil
}

// GetCategory returns a category by ID
func (a *App) GetCategory(id models.CategoryId) (models.Category, error) {
	category, err := a.timekeeper.Storage.Categories().GetCategory(id)
	if err != nil {
		a.logger.Error("Error getting category", "Error", err)
		return models.Category{}, err
	}
	return category, nil
}

// AddCategory adds a new category
func (a *App) AddCategory(category models.Category) error {
	err := a.timekeeper.Storage.Categories().AddCategory(category)
	if err != nil {
		a.logger.Error("Error adding category", "Error", err)
	}
	return err
}

// DeleteCategory deletes a category by ID
func (a *App) DeleteCategory(id models.CategoryId) error {
	err := a.timekeeper.Storage.Categories().DeleteCategory(id)
	if err != nil {
		a.logger.Error("Error deleting category", "Error", err)
	}
	return err
}
