---
name: TimeKeeper architecture and file map
description: Key file locations, storage pattern, data flow, and Storage interface layout
type: project
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
## Storage interface pattern

All storage accessed via `interfaces.Storage`:
```
Categories() CategoryStore
Rules() RuleStore
AppAggregations() AppAggregationStore
CategoryAggregations() CategoryAggregationStore
Events() EventStore
Goals() GoalStore
Preferences() PreferencesStore
Close() error
```

Two implementations: `internal/data/sqlite/` and `internal/data/inmem/`.
**Rule:** Any new store must be added to BOTH implementations and the interface, or the build breaks.

## Key file locations

| Concern | File |
|---------|------|
| Core orchestration | `core/timekeeper.go` |
| Windows observer | `platforms/windows/observer.go` |
| macOS observer | `platforms/darwin/observer.go` |
| Windows browser URL | `platforms/windows/browsers/url_extractor.go` |
| SQLite storage init | `internal/data/sqlite/storage.go` |
| Inmem storage | `internal/data/inmem/storage.go` |
| Wails API surface | `ui/app.go`, `ui/api_rules.go`, `ui/api_categories.go`, `ui/api_events.go`, `ui/api_goals.go`, `ui/api_preferences.go` |
| DTOs (JSON-facing) | `ui/dtos/` |
| Timezone utility | `internal/tzutil/tzutil.go` |
| Frontend entry | `ui/frontend/src/App.svelte` |
| Frontend stores | `ui/frontend/src/stores/` |
| Frontend utils | `ui/frontend/src/utils/` |

## Aggregation date keys

- App: `appName-YYYY-MM-DD` (UTC date)
- Category: `categoryId-YYYY-MM-DD` (UTC date)
- Aggregation tables are UTC-date-keyed. For timezone-aware single-day queries, use `GetEventsByTimeRange` + `tzutil.LocalDayToUTCRange` to get accurate data from raw events. Multi-day range queries (trend, calendar) use aggregation tables directly (acceptable approximation).

## Wails binding regeneration

After adding new Go API methods: run `wails generate module` from `ui/` directory. Generated files in `ui/frontend/wailsjs/` — never hand-edit.

**How to apply:** Check here first when adding a new feature to locate the right file.
