---
name: TimeKeeper project overview
description: What TimeKeeper is, how it's structured, and its two runnable surfaces
type: project
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
TimeKeeper is a cross-platform (Windows + macOS) desktop time-tracking app.

**Why:** Tracks active application usage, maps activity to categories via rules, and stores daily aggregations in SQLite or in-memory.

**Two runnable surfaces:**
- `cmd/cli/main.go` — headless CLI tracker
- `ui/main.go` — Wails v2 desktop app with Svelte frontend

**Tech stack:** Go backend, Wails v2 bindings, Svelte frontend, SQLite (or inmem for dev/tests), Chart.js via svelte-chartjs.

**Event flow:**
1. Platform observer (Windows: WinEvent hooks; macOS: NSWorkspace + AX) emits `AppSwitchEvent`
2. `core.TimeKeeper.PushEvent` receives event
3. Previous event closed, elapsed time computed
4. Aggregated by app and category into SQLite tables
5. Raw events optionally stored (`StoreEvents=true` in SQLite mode)
6. UI reads aggregations + events for dashboard

**Key domain models:** `AppSwitchEvent`, `CategoryRule`, `AppAggregation`, `CategoryAggregation`, `CategoryGoal`, `UserPreferences`

**How to apply:** Use this as the starting mental model for any new feature or bug.
