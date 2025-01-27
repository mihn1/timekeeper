package main

import (
	"fmt"

	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/macos"
)

func main() {
	eventCh := make(chan core.AppSwitchEvent)
	var observer core.Observer = macos.NewObserver()
	go observer.StartObserving(eventCh)

	for event := range eventCh {
		// TODO: aggregate events
		fmt.Println(event)
	}
}
