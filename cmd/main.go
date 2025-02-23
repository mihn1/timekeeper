package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/macos"
)

func main() {
	// Define flags
	dbType := flag.String("db", "sqlite", "Database type: 'sqlite' or 'inmem'")
	dbPath := flag.String("dbpath", "./timekeeper.db", "Path to SQLite database file")
	seedData := flag.Bool("seed", false, "Seed initial data")

	// Parse flags
	flag.Parse()

	var timekeeper *core.TimeKeeper

	switch *dbType {
	case "sqlite":
		if dbPath == nil {
			defaultDbPath := "./timekeeper.db"
			dbPath = &defaultDbPath // Default db path
		}
		log.Println("Starting sqlite Timekeeper with database path:", *dbPath)
		timekeeper = core.NewTimeKeeperSqlite(*dbPath)
	case "inmem":
		log.Println("Starting inmem Timekeeper")
		timekeeper = core.NewTimeKeeperInMem()
	default:
		panic(fmt.Sprintf("Invalid database type %s", *dbType))
	}

	if seedData != nil && *seedData {
		core.SeedData(timekeeper)
	}

	defer timekeeper.Close()

	timekeeper.StartTracking()
	observer := macos.NewObserver(timekeeper)
	go observer.StartObserving()

	select {}
}
