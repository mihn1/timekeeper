package main

import (
	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/macos"
)

func main() {
	timekeeper := core.NewTimeKeeperInMem()
	core.SeedDataInMem(timekeeper)

	timekeeper.StartTracking()
	var observer core.Observer = macos.NewObserver()
	go observer.StartObserving(timekeeper)

	select {}
}
