# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Snapshot

TimeKeeper is a cross-platform (macOS + Windows) time tracking app written in Go.
It tracks active applications and browser tab changes, maps activity to categories using rules, and stores daily aggregations in SQLite (or in-memory for tests/dev).

There are two runnable surfaces:

- `cmd/cli/main.go`: headless tracker process
- `ui/main.go`: Wails desktop app with Svelte frontend

## Commands

From repo root:

```bash
# Run all tests (excluding ui which needs frontend built)
go test ./core/... ./platforms/... ./internal/...

# Run a single test
go test ./core/... -run TestName

# Build CLI binary (macOS)
go build -o timekeeper ./cmd/cli

# Build CLI binary (Windows)
GOARCH=amd64 GOOS=windows go build -o timekeeper.exe ./cmd/cli

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

1. Platform observer emits `AppSwitchEvent`
2. `core.TimeKeeper` receives event through `PushEvent`
3. Previous active event is closed when new one arrives
4. Elapsed time is aggregated by app and category
5. Optional raw event persistence (SQLite mode)
6. UI/CLI reads daily aggregates for reporting

Core components:

- `platforms/observer_darwin.go` / `platforms/observer_windows.go`: platform factory that returns the correct observer
- `platforms/darwin/observer.go`: subscribes to app activation/launch/terminate notifications (darwin)
- `platforms/darwin/browsers/tab_observer.go` + `tab_observer.m`: accessibility + AppleScript browser tab observer bridge (darwin)
- `platforms/windows/observer.go`: WinEvent hook-based observer (windows)
- `platforms/windows/browsers/url_extractor.go`: browser URL extraction via Win32 API + session file fallback (windows)
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
- `platforms/` — all platform-specific code
  - `platforms/observer_darwin.go` — factory: wraps `platforms/darwin.NewObserver`
  - `platforms/observer_windows.go` — factory: wraps `platforms/windows.NewObserver`
  - `platforms/darwin/observer.go` — macOS observer (darwinkit + NSWorkspace notifications)
  - `platforms/darwin/browsers/` — AppleScript tab observer (Cgo + Objective-C)
  - `platforms/windows/observer.go` — WinEvent hook observer
  - `platforms/windows/browsers/url_extractor.go` — Win32 + PowerShell URL extraction
- `ui/` — Wails Go app + bound API methods
- `ui/dtos/` — JSON-facing DTOs for Wails methods
- `ui/frontend/` — Svelte UI (dashboard, rules, categories)
- `ui/frontend/wailsjs/` — **generated** Wails bindings — never hand-edit, regenerate via `wails dev`

## Runtime Modes

### CLI mode (`cmd/cli/main.go`)

- Supports `--db sqlite|inmem`
- SQLite mode enables raw event persistence (`StoreEvents=true`)
- In-memory mode forces seed-only behavior
- Uses standalone observer (`isStandalone=true`) via `platforms.NewPlatformObserver`

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

**Add a new platform observer:** implement `Start()`/`Stop()` satisfying `core.Observer` → add a `platforms/observer_GOOS.go` file with the appropriate build tag that returns it from `NewPlatformObserver`.

**Change rule/category matching:** update `internal/models/rule.go` → adjust resolver flow in `core/timekeeper.go` if needed → add tests in `core/resolvers/resolvers_test.go` and `core/timekeeper_test.go`.

**Change storage schema:** update interface in `internal/data/interfaces/` → implement in both `internal/data/sqlite/` and `internal/data/inmem/` → update call sites in `core/` and `ui/`.

**Update dashboard data shape:** adjust Go return DTO/model in `ui/app.go` → update generated binding usage in `ui/frontend` → update chart/table components.

## Windows Observer

The Windows observer (`platforms/windows/observer.go`) uses Win32 WinEvent hooks:

- `EVENT_SYSTEM_FOREGROUND` — fires when the active window changes
- `EVENT_OBJECT_NAMECHANGE` — fires when the foreground window's title changes (catches browser tab switches)
- Both hooks use `WINEVENT_OUTOFCONTEXT` which requires a Win32 message pump (`GetMessageW`/`DispatchMessageW`) on the hook thread
- `core.StartTracking` calls `go obs.Start()` so the blocking message loop runs in its own goroutine
- `Stop()` posts `WM_QUIT` to the message loop thread; a `readyCh` channel prevents a race when `Stop()` is called before `o.tid` is set

### Browser URL extraction (Windows)

`platforms/windows/browsers/url_extractor.go` tries two approaches in order:

1. **Win32 child-window enumeration** — walks the browser window's child HWND tree looking for `Chrome_OmniboxView` / `Edit` class controls and reads their text via `WM_GETTEXT`
2. **PowerShell UI Automation fallback** — targets the specific window handle via `AutomationElement::FromHandle`

Both fail gracefully if the browser window is not in focus or the HWND is invalid.

### App name normalization

`normalizeAppName` maps `.exe` base names to display names:

| Executable   | App name        |
|--------------|-----------------|
| `chrome.exe` | `Google Chrome` |
| `brave.exe`  | `Brave Browser` |
| `msedge.exe` | `Microsoft Edge`|
| `firefox.exe`| `Firefox`       |
| others       | stem (lower-case, no `.exe`) |

## Known Quirks

- macOS tracking relies on darwinkit, AX observer, and AppleScript; Windows tracking uses WinEvent hooks.
- Events within 60 seconds of an identical prior event are deduplicated (`isSameEvent` in `core/timekeeper.go`).
- `wailsjs/` files are generated artifacts; do not edit manually.
- `GetCategoryUsageData` returns `[]map[string]any` (not a typed DTO) — enriches aggregations with category names inline.
- darwinkit (`github.com/progrium/darwinkit`) is listed in `go.mod` but is only needed on darwin; `go mod tidy` on Windows will remove it — add it back manually or run tidy on macOS.

## Recommended First Reads

1. `core/timekeeper.go`
2. `ui/app.go`
3. `ui/frontend/src/App.svelte`
4. `internal/data/sqlite/storage.go`

## Workflow Rules

- Do not create commits unless the user explicitly asks for one in the current conversation.
- Default behavior is to leave changes in the working tree and let the user commit manually.
- Do not push branches or open pull requests unless explicitly requested.
