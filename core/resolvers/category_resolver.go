package resolvers

import (
	"github.com/mihn1/timekeeper/internal/models"
)

// resolve category from an app switch event
type CategoryResolver interface {
	ResolveCategory(event *models.AppSwitchEvent, rules []*models.CategoryRule) (models.CategoryId, error)
}
