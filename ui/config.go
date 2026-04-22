package main

import (
	"strconv"
	"strings"
)

const (
	defaultDBType            = "sqlite"
	defaultDBPath            = "../db/timekeeper.db"
	defaultSeedMode          = "if-empty"
	defaultMaxRerunRangeDays = 7
)

type AppConfig struct {
	DBType            string
	DBPath            string
	SeedMode          string
	MaxRerunRangeDays int
}

func LoadAppConfig(getenv func(string) string) AppConfig {
	dbType := strings.ToLower(strings.TrimSpace(getenv("TIMEKEEPER_DB")))
	if dbType == "" {
		dbType = defaultDBType
	}

	dbPath := strings.TrimSpace(getenv("TIMEKEEPER_DB_PATH"))
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	seedMode := strings.ToLower(strings.TrimSpace(getenv("TIMEKEEPER_SEED_MODE")))
	if seedMode != "always" && seedMode != "never" && seedMode != "if-empty" {
		seedMode = defaultSeedMode
	}

	maxRerun := defaultMaxRerunRangeDays
	if v := strings.TrimSpace(getenv("TIMEKEEPER_MAX_RERUN_RANGE_DAYS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxRerun = n
		}
	}

	return AppConfig{
		DBType:            dbType,
		DBPath:            dbPath,
		SeedMode:          seedMode,
		MaxRerunRangeDays: maxRerun,
	}
}
