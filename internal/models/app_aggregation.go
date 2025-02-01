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

func GetAppAggregationKey(event *AppSwitchEvent) string {
	return event.AppName + "-" + event.GetEventDate().String()
}

func (a AppAggregation) String() string {
	fullAppName := a.AppName
	return fmt.Sprintf("%s: %ds", fullAppName, a.TimeElapsed)
}
