# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Snapshot

TimeKeeper is a macOS-focused time tracking app written in Go.
It tracks active applications and browser tab changes, maps activity to categories using rules, and stores daily aggregations in SQLite (or in-memory for tests/dev).

There are two runnable surfaces:

- `cmd/cli/main.go`: headless tracker process
- `ui/main.go`: Wails desktop app with Svelte frontend

## Commands

From repo root:

```bash
# Run all tests
go test ./...

# Run a single test
go test ./core/... -run TestName

# Build CLI binary
go build -o timekeeper ./cmd/cli

# Run CLI with SQLite
./timekeeper --db sqlite --dbpath ./db/timekeeper.db

# Run CLI with in-memory store
./timekeeper --db inmem
```

Wails desktop app (from `ui/`):

```bash
cd ui && wails dev
```

Frontend only (from `ui/frontend/`):

```bash
npm install
npm run dev
npm run build
```

## High-Level Architecture

Event flow:

1. macOS observer emits `AppSwitchEvent`
2. `core.TimeKeeper` receives event through `PushEvent`
3. Previous active event is closed when new one arrives
4. Elapsed time is aggregated by app and category
5. Optional raw event persistence (SQLite mode)
6. UI/CLI reads daily aggregates for reporting

Core components:

- `macos/observer.go`: subscribes to app activation/launch/terminate notifications
- `macos/browsers/tab_observer.go` + `tab_observer.m`: accessibility + AppleScript browser tab observer bridge
- `core/timekeeper.go`: orchestration, event handling, aggregation, exclusion handling
- `core/resolvers/default_category_resolver.go`: rule-based category resolution
- `internal/data/*`: storage abstraction and implementations (`sqlite`, `inmem`)

## Repository Map

- `cmd/cli/main.go` — CLI entrypoint, flags, standalone observer startup
- `core/` — domain logic (`TimeKeeper`, seeding, resolver, errors)
- `constants/` — known browser app names + additional data keys
- `datatypes/` — `DateOnly` custom type used in stores/queries
- `internal/models/` — domain models (`AppSwitchEvent`, `CategoryRule`, aggregations)
- `internal/data/interfaces/` — storage interfaces
- `internal/data/sqlite/` — SQLite-backed stores, table creation at startup
- `internal/data/inmem/` — in-memory stores for tests/dev
- `macos/` — macOS observation layer
- `ui/` — Wails Go app + bound API methods
- `ui/dtos/` — JSON-facing DTOs for Wails methods
- `ui/frontend/` — Svelte UI (dashboard, rules, categories)
- `ui/frontend/wailsjs/` — **generated** Wails bindings — never hand-edit, regenerate via `wails dev`

## Runtime Modes

### CLI mode (`cmd/cli/main.go`)

- Supports `--db sqlite|inmem`
- SQLite mode enables raw event persistence (`StoreEvents=true`)
- In-memory mode forces seed-only behavior
- Uses standalone macOS observer (`isStandalone=true`)

### Wails desktop mode (`ui/main.go`)

- Starts Wails app and binds methods from `ui/app.go`
- Uses non-standalone observer (`isStandalone=false`)
- Configuration is env-driven via `ui/config.go` (`LoadAppConfig`):
  - `TIMEKEEPER_DB` — `sqlite` (default) or `inmem`
  - `TIMEKEEPER_DB_PATH` — SQLite file path (default: `../db/timekeeper.db`)
  - `TIMEKEEPER_SEED_MODE` — `if-empty` (default), `always`, or `never`
- Frontend calls Go methods through generated `wailsjs` bindings
- `timekeeper:data-updated` event is emitted (e.g. View > Refresh Data, Cmd+R) to trigger frontend refresh

## Data and Persistence

SQLite store constructors create tables if missing:

- `categories`
- `rules`
- `app_aggregations`
- `category_aggregations`
- `events`

Aggregations are per date (`datatypes.DateOnly`) and keyed as:

- app key: `appName-YYYY-MM-DD`
- category key: `categoryId-YYYY-MM-DD`

`TimeKeeperOptions.StoreEvents` controls raw event storage in `events` table.

## Rule and Category Resolution

Rule model: `internal/models/rule.go`

- Match by app name if no `AdditionalDataKey`/`Expression`
- Else evaluate against `event.AdditionalData[key]`
- Supports plain substring match or regex (`IsRegex`); compiled regex cached in `sync.Map`
- `Priority` determines ordering (higher first)
- `IsExclusion` rules short-circuit processing for matching events

Resolution order in `core/timekeeper.go`:

1. App-specific rules (`GetRulesByApp(event.AppName)`)
2. Exclusion check
3. Global rules (`GetRulesByApp(constants.ALL_APPS)`)
4. Exclusion check
5. Resolver selects first matching rule
6. Fallback category is `models.UNDEFINED`

## Frontend/Backend Contract (Wails)

Primary bound methods used by UI:

- Tracking: `IsTrackingEnabled`, `EnableTracking`, `DisableTracking`
- Dashboard: `GetAppUsageData(dateStr)`, `GetCategoryUsageData(dateStr)`
- Rules CRUD: `GetRules`, `GetRule`, `AddRule`, `UpdateRule`, `DeleteRule`
- Categories CRUD: `GetCategories`, `GetCategory`, `AddCategory`, `UpdateCategory`, `DeleteCategory`

Files to inspect first:

- Go API: `ui/app.go`, `ui/api_rules.go`, `ui/api_categories.go`
- UI entry: `ui/frontend/src/App.svelte`
- Views: `ui/frontend/src/components/Dashboard.svelte`, `ui/frontend/src/components/rules/Rules.svelte`, `ui/frontend/src/components/categories/Categories.svelte`

## Change Cookbook

**New Wails API method:** implement in `ui/app.go` or a new `ui/api_*.go` file → run `wails dev` to regenerate `wailsjs` bindings → call from Svelte.

**Change rule/category matching:** update `internal/models/rule.go` → adjust resolver flow in `core/timekeeper.go` if needed → add tests in `core/resolvers/resolvers_test.go` and `core/timekeeper_test.go`.

**Change storage schema:** update interface in `internal/data/interfaces/` → implement in both `internal/data/sqlite/` and `internal/data/inmem/` → update call sites in `core/` and `ui/`.

**Update dashboard data shape:** adjust Go return DTO/model in `ui/app.go` → update generated binding usage in `ui/frontend` → update chart/table components.

## Known Quirks

- macOS-only: tracking relies on darwinkit, AX observer, and AppleScript.
- Events within 60 seconds of an identical prior event are deduplicated (`isSameEvent` in `core/timekeeper.go`).
- `wailsjs/` files are generated artifacts; do not edit manually.
- `GetCategoryUsageData` returns `[]map[string]any` (not a typed DTO) — enriches aggregations with category names inline.

## Recommended First Reads

1. `core/timekeeper.go`
2. `ui/app.go`
3. `ui/frontend/src/App.svelte`
4. `internal/data/sqlite/storage.go`

## Workflow Rules

- Do not create commits unless the user explicitly asks for one in the current conversation.
- Default behavior is to leave changes in the working tree and let the user commit manually.
- Do not push branches or open pull requests unless explicitly requested.
