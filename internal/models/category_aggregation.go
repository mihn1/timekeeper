package models

import (
	"fmt"

	"github.com/mihn1/timekeeper/datatypes"
	"github.com/mihn1/timekeeper/utils"
)

type CategoryAggregation struct {
	CategoryId  CategoryId
	Date        datatypes.DateOnly
	TimeElapsed int64 // in miliseconds
}

func GetCategoryAggregationKey(categoryId CategoryId, date datatypes.DateOnly) string {
	return fmt.Sprintf("%v-%s", categoryId, date)
}

func (c *CategoryAggregation) String() string {
	return fmt.Sprintf("%v: %s", c.CategoryId, utils.FormatTimeElapsed(c.TimeElapsed))
}
