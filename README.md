# TimeKeeper

TimeKeeper is a free and open source time tracking app focused on macOS.
It tracks active application usage and browser tab changes, then maps activity to categories using configurable rules.

For contributor and LLM-oriented project context, read `AGENTS.md` first.

## Current Scope

- Tracks app switching and time spent per app
- Tracks active tab metadata for supported browsers (Chrome, Brave, Safari)
- Resolves events into categories via rule matching
- Stores data in SQLite or in-memory backends
- Includes both a headless CLI and a Wails desktop app (Svelte frontend)

## Planned / Not Yet Implemented

- Session-oriented workflows
- Goal tracking
- Cross-platform tracking support (non-macOS)
- Multi-device sync/account model

## Architecture at a Glance

- macOS observer emits `AppSwitchEvent`
- `core.TimeKeeper` handles event transitions and elapsed-time aggregation
- Rule engine resolves category
- Storage layer persists app/category aggregations (and optional raw events)
- UI reads daily aggregates and manages rules/categories through Wails-bound APIs

## Project Entry Points

- CLI tracker: `cmd/cli/main.go`
- Desktop app: `ui/main.go`

## Quick Start

### Requirements

- Go `1.23.5+`
- macOS (for runtime tracking)
- SQLite (bundled through `github.com/mattn/go-sqlite3`)
- Wails CLI (for desktop app development)
- Node.js/npm (for frontend development)

### Clone

```bash
git clone https://github.com/mihn1/timekeeper.git
cd timekeeper
```

### Run Tests

```bash
go test ./...
```

### Build and Run CLI

```bash
# Build
go build -o timekeeper ./cmd/cli

# Run with SQLite storage
./timekeeper --db sqlite --dbpath ./db/timekeeper.db

# Run with in-memory storage
./timekeeper --db inmem
```

### Run Desktop App (Wails)

```bash
cd ui
wails dev
```

### Frontend-Only Workflow

```bash
cd ui/frontend
npm install
npm run dev
npm run build
```

## Data Storage

SQLite mode creates and uses these tables:

- `categories`
- `rules`
- `app_aggregations`
- `category_aggregations`
- `events`

## Notes for Contributors

- `ui/frontend/wailsjs/` files are generated; do not hand-edit
- Runtime tracking currently depends on macOS accessibility APIs
- Wails app currently initializes with a fixed DB path in `ui/app.go`
- For deep codebase orientation and change cookbook, use `AGENTS.md`
