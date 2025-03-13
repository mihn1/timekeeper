package main

import (
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/ui/dtos"
)

// GetCategories returns a list of all categories
func (a *App) GetCategories() []*dtos.CategoryListItem {
	categories, err := a.timekeeper.Storage.Categories().GetCategories()

	if err != nil {
		a.logger.Error("Error getting categories", "Error", err)
		return nil
	}

	return dtos.CategoryListFromModels(categories)
}

// GetCategory returns a single category by ID
func (a *App) GetCategory(id int) (*dtos.CategoryDetail, error) {
	category, err := a.timekeeper.Storage.Categories().GetCategory(models.CategoryId(id))

	if err != nil {
		a.logger.Error("Error getting category", "Error", err)
		return nil, err
	}

	return dtos.CategoryDetailFromModel(category), nil
}

// AddCategory creates a new category
func (a *App) AddCategory(categorydtos *dtos.CategoryCreate) error {
	category := categorydtos.ToModel()
	err := a.timekeeper.Storage.Categories().UpsertCategory(category)

	if err != nil {
		a.logger.Error("Error adding category", "Error", err)
	}

	return err
}

// UpdateCategory updates an existing category
func (a *App) UpdateCategory(id int, categorydtos *dtos.CategoryUpdate) error {
	// First get the existing category
	existingCategory, err := a.timekeeper.Storage.Categories().GetCategory(models.CategoryId(id))
	if err != nil {
		a.logger.Error("Error getting category for update", "Error", err)
		return err
	}

	// Update fields
	existingCategory.Name = categorydtos.Name
	existingCategory.Description = categorydtos.Description

	// Save the updated category
	err = a.timekeeper.Storage.Categories().UpsertCategory(existingCategory)
	if err != nil {
		a.logger.Error("Error updating category", "Error", err)
	}

	return err
}

// DeleteCategory deletes a category by ID
func (a *App) DeleteCategory(id int) error {
	err := a.timekeeper.Storage.Categories().DeleteCategory(models.CategoryId(id))

	if err != nil {
		a.logger.Error("Error deleting category", "Error", err)
	}

	return err
}
