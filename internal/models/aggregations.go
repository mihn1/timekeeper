package models

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/datatypes"
)

type AppAggregation struct {
	AppName        string
	AdditionalData interface{} // E.g. title of a browser tab or url // TODO: redesign this to support multiple fields
	Date           datatypes.Date
	TimeElapsed    int
}

func (a AppAggregation) String() string {
	fullAppName := a.AppName
	return fmt.Sprintf("%s - %d", fullAppName, a.TimeElapsed)
}

type CategoryAggregation struct {
	CategoryId  CategoryId
	Date        datatypes.Date
	TimeElapsed int
}

func (c CategoryAggregation) String() string {
	return fmt.Sprintf("%d - %d", c.CategoryId, c.TimeElapsed)
}
