package data

import "fmt"

type AppAggregation struct {
	AppName     string
	SubAppName  string // E.g. title of a browser tab
	TimeElapsed int
}

func (a AppAggregation) String() string {
	fullAppName := a.AppName
	if a.SubAppName != "" {
		fullAppName += "-" + a.SubAppName
	}
	return fmt.Sprintf("%s - %d", fullAppName, a.TimeElapsed)
}

type CategoryAggregation struct {
	CategoryId  CategoryId
	TimeElapsed int
}

func (c CategoryAggregation) String() string {
	return fmt.Sprintf("%d - %d", c.CategoryId, c.TimeElapsed)
}
