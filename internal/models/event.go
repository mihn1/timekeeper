package models

import (
	"fmt"
	"time"

	"github.com/mihn1/timekeeper/datatypes"
)

type EventId int

type AppSwitchEvent struct {
	Id             EventId
	AppName        string
	StartTime      time.Time
	EndTime        time.Time
	CategoryId     CategoryId
	AdditionalData map[string]string
}

func (e *AppSwitchEvent) String() string {
	return fmt.Sprintf("App Changed: %s, Time: %s, AdditionalData: %s", e.AppName, e.StartTime.Format(time.DateTime), e.AdditionalData)
}

func (e *AppSwitchEvent) GetEventDate() datatypes.DateOnly {
	return datatypes.NewDateOnly(e.StartTime)
}
