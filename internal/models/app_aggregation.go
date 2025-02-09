package models

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/utils"
)

type AppAggregation struct {
	AppName        string
	AdditionalData interface{} // E.g. title of a browser tab or url // TODO: redesign this to support multiple fields
	Date           datatypes.Date
	TimeElapsed    int64 // in milliseconds
}

func GetAppAggregationKey(event *AppSwitchEvent) string {
	return event.AppName + "-" + event.GetEventDate().String()
}

func (a *AppAggregation) String() string {
	fullAppName := a.AppName
	return fmt.Sprintf("%s: %v", fullAppName, utils.FormatTimeElapsed(a.TimeElapsed))
}
