package data

import (
	"fmt"
	"time"
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
