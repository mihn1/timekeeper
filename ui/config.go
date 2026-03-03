package main

import "strings"

const (
	defaultDBType   = "sqlite"
	defaultDBPath   = "../db/timekeeper.db"
	defaultSeedMode = "if-empty"
)

type AppConfig struct {
	DBType   string
	DBPath   string
	SeedMode string
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

	return AppConfig{
		DBType:   dbType,
		DBPath:   dbPath,
		SeedMode: seedMode,
	}
}
