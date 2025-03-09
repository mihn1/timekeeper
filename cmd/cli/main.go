package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mihn1/timekeeper/core"
	"github.com/mihn1/timekeeper/macos"
)

func main() {
	// Define flags
	dbType := flag.String("db", "sqlite", "Database type: 'sqlite' or 'inmem'")
	dbPath := flag.String("dbpath", "./db/timekeeper.db", "Path to SQLite database file")
	seed := flag.Bool("seed", true, "Seed initial data")
	seedOnly := flag.Bool("seedonly", false, "Seed initial data")

	flag.Parse()

	logger := slog.Default()
	var timekeeper *core.TimeKeeper
	opts := core.TimeKeeperOptions{
		StoreEvents: false,
		Logger:      logger,
	}

	switch *dbType {
	case "sqlite":
		if dbPath == nil {
			logger.Info("No database path provided, using default")
			defaultDbPath := "./db/timekeeper.db"
			dbPath = &defaultDbPath // Default db path
		}

		logger.Info("Starting sqlite Timekeeper", "DbPath", *dbPath)
		opts.StoragePath = *dbPath
		opts.StoreEvents = true // Only store events in SQLite
		timekeeper = core.NewTimeKeeperSqlite(opts)
	case "inmem":
		logger.Info("Starting inmem Timekeeper")
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

	observer := macos.NewObserver(timekeeper.PushEvent, logger)
	timekeeper.AddObserver(observer)
	timekeeper.StartTracking()

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
