package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/macos"
)

func main() {
	// Define flags
	dbType := flag.String("db", "sqlite", "Database type: 'sqlite' or 'inmem'")
	dbPath := flag.String("dbpath", "./db/timekeeper.db", "Path to SQLite database file")
	seed := flag.Bool("seed", true, "Seed initial data")
	seedOnly := flag.Bool("seedonly", false, "Seed initial data")

	flag.Parse()

	var timekeeper *core.TimeKeeper
	opts := core.TimeKeeperOptions{
		StoreEvents: false,
	}

	switch *dbType {
	case "sqlite":
		if dbPath == nil {
			defaultDbPath := "./timekeeper.db"
			dbPath = &defaultDbPath // Default db path
		}

		log.Println("Starting sqlite Timekeeper with database path:", *dbPath)
		opts.StoragePath = *dbPath
		opts.StoreEvents = true // Only store events in SQLite
		timekeeper = core.NewTimeKeeperSqlite(opts)
	case "inmem":
		log.Println("Starting inmem Timekeeper")
		*seedOnly = true // Always seed data for inmem
		timekeeper = core.NewTimeKeeperInMem(opts)
	default:
		panic(fmt.Sprintf("Invalid database type %s", *dbType))
	}

	if *seed {
		seedData(timekeeper)
		if *seedOnly {
			return
		}
	}

	defer timekeeper.Close()

	timekeeper.StartTracking()
	observer := macos.NewObserver(timekeeper)
	go observer.StartObserving()

	// Set up channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)
}

func seedData(timekeeper *core.TimeKeeper) {
	core.SeedData(timekeeper)
}
