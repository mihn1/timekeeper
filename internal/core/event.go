package core

import (
	"fmt"
	"time"
)

type AppSwitchEvent struct {
	AppName        string
	Time           time.Time
	AdditionalData interface{}
}

func (e AppSwitchEvent) String() string {
	return fmt.Sprintf("App Changed: %s, Time: %s", e.AppName, e.Time.Format(time.DateTime))
}
