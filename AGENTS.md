# AGENTS.md

## Project Snapshot

TimeKeeper is a macOS-focused time tracking app written in Go.
It tracks active applications and browser tab changes, maps activity to categories using rules, and stores daily aggregations in SQLite (or in-memory for tests/dev).

There are two runnable surfaces:

- `cmd/cli/main.go`: headless tracker process
- `ui/main.go`: Wails desktop app with Svelte frontend

## High-Level Architecture

Event flow:

1. macOS observer emits `AppSwitchEvent`
2. `core.TimeKeeper` receives event through `PushEvent`
3. previous active event is closed when new one arrives
4. elapsed time is aggregated by app and category
5. optional raw event persistence (SQLite mode)
6. UI/CLI reads daily aggregates for reporting

Core components:

- `macos/observer.go`: subscribes to app activation/launch/terminate notifications
- `macos/browsers/tab_observer.go` + `tab_observer.m`: accessibility + AppleScript browser tab observer bridge
- `core/timekeeper.go`: orchestration, event handling, aggregation, exclusion handling
- `core/resolvers/default_category_resolver.go`: rule-based category resolution
- `internal/data/*`: storage abstraction and implementations (`sqlite`, `inmem`)

## Repository Map

- `cmd/cli/main.go` - CLI entrypoint, flags, standalone observer startup
- `core/` - domain logic (`TimeKeeper`, seeding, resolver, errors)
- `constants/` - known browser app names + additional data keys
- `datatypes/` - `DateOnly` custom type used in stores/queries
- `internal/models/` - domain models (`AppSwitchEvent`, `CategoryRule`, aggregations)
- `internal/data/interfaces/` - storage interfaces
- `internal/data/sqlite/` - SQLite-backed stores, table creation at startup
- `internal/data/inmem/` - in-memory stores for tests/dev
- `macos/` - macOS observation layer
- `ui/` - Wails Go app + bound API methods
- `ui/dtos/` - JSON-facing DTOs for Wails methods
- `ui/frontend/` - Svelte UI (dashboard, rules, categories)
- `ui/frontend/wailsjs/` - generated Wails bindings (do not hand-edit)

## Runtime Modes

### CLI mode (`cmd/cli/main.go`)

- Supports `--db sqlite|inmem`
- SQLite mode enables raw event persistence (`StoreEvents=true`)
- In-memory mode forces seed-only behavior
- Uses standalone macOS observer (`isStandalone=true`)

### Wails desktop mode (`ui/main.go`)

- Starts Wails app and binds methods from `ui/app.go`
- Uses non-standalone observer (`isStandalone=false`)
- Frontend calls Go methods through generated `wailsjs` bindings
- Emits and listens to `timekeeper:data-updated` events for refresh

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
- Supports plain substring match or regex (`IsRegex`)
- `Priority` determines ordering (higher first)
- `IsExclusion` rules short-circuit processing for matching events

Resolution process in `core/timekeeper.go`:

1. app-specific rules (`GetRulesByApp(event.AppName)`)
2. exclusion check
3. global rules (`GetRulesByApp(constants.ALL_APPS)`)
4. exclusion check
5. resolver selects first matching rule
6. fallback category is `models.UNDEFINED`

## Frontend/Backend Contract (Wails)

Primary bound methods used by UI:

- Tracking: `IsTrackingEnabled`, `EnableTracking`, `DisableTracking`
- Dashboard: `GetAppUsageData(date)`, `GetCategoryUsageData(date)`
- Rules CRUD: `GetRules`, `GetRule`, `AddRule`, `UpdateRule`, `DeleteRule`
- Categories CRUD: `GetCategories`, `GetCategory`, `AddCategory`, `UpdateCategory`, `DeleteCategory`

Files to inspect first:

- Go API: `ui/app.go`, `ui/api_rules.go`, `ui/api_categories.go`
- UI entry: `ui/frontend/src/App.svelte`
- Views: `ui/frontend/src/components/Dashboard.svelte`, `ui/frontend/src/components/rules/Rules.svelte`, `ui/frontend/src/components/categories/Categories.svelte`

## Commands

From repo root:

```bash
# Run tests
go test ./...

# Build CLI binary
go build -o timekeeper ./cmd/cli

# Run CLI with SQLite
./timekeeper --db sqlite --dbpath ./db/timekeeper.db
```

Wails desktop app:

```bash
cd ui
wails dev
```

Frontend-only (inside `ui/frontend`):

```bash
npm install
npm run dev
npm run build
```

## Known Quirks and Caveats

- macOS-only tracking stack today (darwinkit + AX observer + AppleScript)
- `ui/app.go` currently seeds data on startup and uses fixed SQLite path `../db/timekeeper-refactor.db`
- `wailsjs` files are generated artifacts; regenerate instead of manual edits
- Event loop treats near-identical events within 60 seconds as the same (`isSameEvent`)
- Existing tests are basic and not full end-to-end behavior coverage

## Fast Change Cookbook

- Add a new backend API method:
  1. implement in `ui/app.go` or split API file in `ui/`
  2. run Wails codegen/dev to refresh `ui/frontend/wailsjs`
  3. call from Svelte component/store

- Change category matching behavior:
  1. update `internal/models/rule.go` (match semantics)
  2. update resolver flow in `core/timekeeper.go` if needed
  3. add tests in `core/resolvers/resolvers_test.go` and `core/timekeeper_test.go`

- Change storage behavior/schema:
  1. update interface in `internal/data/interfaces/`
  2. implement in both `internal/data/sqlite/` and `internal/data/inmem/`
  3. update call sites in `core/` and UI API

- Update dashboard data shape:
  1. adjust Go return DTO/model in `ui/app.go`
  2. update generated binding usage in `ui/frontend`
  3. update chart/table components

## Recommended First Reads for New Agents

1. `AGENTS.md` (this file)
2. `README.md`
3. `core/timekeeper.go`
4. `ui/app.go`
5. `ui/frontend/src/App.svelte`
6. `internal/data/sqlite/storage.go`

## Agent Workflow Rules

- Do not create commits unless the user explicitly asks for a commit in the current conversation.
- Default behavior is to leave changes in the working tree and let the user commit manually.
- Do not push branches or open pull requests unless explicitly requested.
