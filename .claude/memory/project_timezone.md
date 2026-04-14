---
name: Timezone support implementation
description: How timezone preferences flow from storage through backend queries to frontend display
type: project
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
## Storage

`user_preferences` table (SQLite key-value): `key TEXT PRIMARY KEY, value TEXT NOT NULL`.
Currently stores one row: `key='timezone', value='Asia/Ho_Chi_Minh'` (or whatever user sets).
Extensible: add new preference keys without schema change.

## Backend flow

1. `App.Startup` loads `UserPreferences` from `storage.Preferences().GetPreferences()` into `App.prefs` (RWMutex protected).
2. `App.getTimezone()` reads `a.prefs.Timezone` under read lock — called by all data API methods.
3. `App.SavePreferences(dto)` persists to DB AND updates `a.prefs` in-memory — data queries immediately use the new timezone without restart.

## Timezone-aware single-day queries

`internal/tzutil/tzutil.go`:
```go
LocalDayToUTCRange("2026-04-14", "Asia/Ho_Chi_Minh")
// → start: 2026-04-13 17:00 UTC, end: 2026-04-14 17:00 UTC
```

`GetAppUsageData` and `GetCategoryUsageData`:
1. Compute UTC range via `LocalDayToUTCRange`
2. Query `Events().GetEventsByTimeRange(start, end)` (queries by `start_time` column)
3. Aggregate from raw events using `tzutil.AggregateEventsByApp/Category`
4. Falls back to pre-computed aggregation tables if no events (inmem mode)

`GetEventLog`: same range query + formats `StartTime`/`EndTime` in user's `*time.Location`.

## Multi-day range queries

`GetCategoryUsageRange` and `GetActivityCalendar` still use aggregation tables (UTC-date-keyed). The error at day boundaries is acceptable for multi-day trend/calendar views.

## Frontend

`stores/preferences.js`:
- Initializes from localStorage instantly (fast first render)
- Syncs from backend via `GetPreferences()` on startup
- `preferences.setTimezone(tz)` is async — updates store + calls `SavePreferences` to backend

`stores/timezone` (derived store): IANA string, consumed by Dashboard and EventLog.

Dashboard `today` is reactive to `$timezone`: `$: today = todayInTz($timezone)`.

`utils/dateUtils.js`: `todayInTz(tz)`, `shiftDateStr(dateStr, days)`, `formatDateDisplay(dateStr, tz)` — all timezone-aware, no UTC-offset bugs.

**How to apply:** When debugging "wrong day shown" or "events attributed to wrong date", check `LocalDayToUTCRange` and the `GetEventsByTimeRange` query.
