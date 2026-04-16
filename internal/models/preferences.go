package models

// UserPreferences holds user-level application preferences persisted in storage.
type UserPreferences struct {
	Timezone             string `json:"timezone"`             // IANA timezone string, e.g. "Asia/Ho_Chi_Minh"
	MinEventDurationMs   int64  `json:"minEventDurationMs"`   // Minimum event duration to count; shorter events are discarded as noise
}

// DefaultPreferences returns safe defaults used when no preferences have been saved yet.
func DefaultPreferences() *UserPreferences {
	return &UserPreferences{
		Timezone:           "UTC",
		MinEventDurationMs: 1000,
	}
}
