package models

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/utils"
)

type CategoryAggregation struct {
	CategoryId  CategoryId
	Date        datatypes.DateOnly
	TimeElapsed int64 // in miliseconds
}

func GetCategoryAggregationKey(categoryId CategoryId, date datatypes.DateOnly) string {
	return string(categoryId) + "-" + date.String()
}

func (c *CategoryAggregation) String() string {
	return fmt.Sprintf("%s: %s", c.CategoryId, utils.FormatTimeElapsed(c.TimeElapsed))
}
