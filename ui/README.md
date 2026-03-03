# UI (Wails Desktop App)

This directory contains the desktop application runtime built with Wails.

- Go desktop entrypoint: `ui/main.go`
- Bound backend app/API: `ui/app.go`, `ui/api_rules.go`, `ui/api_categories.go`
- Frontend source: `ui/frontend/`
- Generated bridge files: `ui/frontend/wailsjs/` (do not hand-edit)

For full project architecture and agent context, read `../AGENTS.md`.

## Run in Development

From this directory:

```bash
wails dev
```

Notes:

- Wails starts the desktop app and a Vite dev server for hot reload.
- Frontend calls Go methods through generated Wails bindings.
- Tracking lifecycle and storage initialization are handled in `ui/app.go`.

## Build Desktop App

From this directory:

```bash
wails build
```

## Important Project-Specific Behavior

- The app currently initializes SQLite with a fixed path in `ui/app.go` (`../db/timekeeper-refactor.db`).
- `core.SeedData(...)` is called on startup.
- Tracking uses the non-standalone macOS observer mode.

## Contributor / Agent Rules

- Do not edit files under `ui/frontend/wailsjs/` manually.
- Regenerate bindings using Wails workflows when backend API signatures change.
- Do not create commits unless explicitly requested by the user.
