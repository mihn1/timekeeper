package dtos

import (
	"github.com/mihn1/timekeeper/internal/models"
)

type CategoryListItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryDetail struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (cc *CategoryCreate) ToModel() *models.Category {
	return &models.Category{
		Name:        cc.Name,
		Description: cc.Description,
	}
}

func (cc *CategoryUpdate) ToModel() *models.Category {
	return &models.Category{
		Id:          models.CategoryId(cc.ID),
		Name:        cc.Name,
		Description: cc.Description,
	}
}

func CategoryDetailFromModel(category *models.Category) *CategoryDetail {
	return &CategoryDetail{
		ID:          int(category.Id),
		Name:        category.Name,
		Description: category.Description,
	}
}

func CategoryListFromModels(categories []*models.Category) []*CategoryListItem {
	result := make([]*CategoryListItem, len(categories))
	for i, category := range categories {
		result[i] = &CategoryListItem{
			ID:          int(category.Id),
			Name:        category.Name,
			Description: category.Description,
		}
	}
	return result
}
