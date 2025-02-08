package models

import (
	"fmt"
	"strconv"

	"github.com/mihn1/timekeeper/internal/datatypes"
)

type CategoryAggregation struct {
	CategoryId  CategoryId
	Date        datatypes.Date
	TimeElapsed int64 // in miliseconds
}

func (c CategoryAggregation) String() string {
	return fmt.Sprintf("%d: %dms", c.CategoryId, c.TimeElapsed)
}

func GetCategoryAggregationKey(categoryId CategoryId, date datatypes.Date) string {
	return strconv.Itoa(int(categoryId)) + "-" + date.String()
}
