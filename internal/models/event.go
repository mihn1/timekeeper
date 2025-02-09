package models

import (
	"fmt"
	"time"

	"github.com/mihn1/timekeeper/internal/datatypes"
)

type AppSwitchEvent struct {
	AppName        string
	StartTime      time.Time
	EndTime        time.Time
	AdditionalData map[string]string
}

func (e *AppSwitchEvent) String() string {
	return fmt.Sprintf("App Changed: %s, Time: %s, AdditionalData: %s", e.AppName, e.StartTime.Format(time.DateTime), e.AdditionalData)
}

func (e *AppSwitchEvent) GetEventDate() datatypes.Date {
	return datatypes.NewDate(e.StartTime)
}
