package main

import (
	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/macos"
)

func main() {
	timekeeper := core.NewTimeKeeperInMem()
	core.SeedDataInMem(timekeeper)

	timekeeper.StartTracking()
	observer := macos.NewObserver(timekeeper)
	go observer.StartObserving()

	select {}
}
