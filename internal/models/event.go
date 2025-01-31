package models

import (
	"fmt"
	"time"

	"github.com/mihn1/timekeeper/internal/datatypes"
)

type AppSwitchEvent struct {
	AppName        string
	SubAppName     string
	Time           time.Time
	AdditionalData interface{}
}

func (e AppSwitchEvent) String() string {
	return fmt.Sprintf("App Changed: %s, Time: %s", e.AppName, e.Time.Format(time.DateTime))
}

func (e AppSwitchEvent) GetEventKey() string {
	return e.AppName + "-" + e.SubAppName
}

func (e AppSwitchEvent) GetEventDate() datatypes.Date {
	return datatypes.NewDate(e.Time)
}
