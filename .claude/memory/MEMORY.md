# Memory Index

## Project
- [Project Overview](project_overview.md) — TimeKeeper: cross-platform Go time tracker, macOS + Windows, Wails desktop + CLI
- [Architecture](project_architecture.md) — event flow, core components, storage layers, platform observers
- [Windows Observer](project_windows_observer.md) — WinEvent hooks, sleep/lock detection, SYSTEM_PAUSED pattern
- [Timezone System](project_timezone.md) — UTC storage, LocalDayToUTCRange, preferences backend, dateUtils.js
- [UI Features](project_ui_features.md) — Dashboard date nav, Preferences page, EventLog, Menu layout
- [Rules & Categories](project_rules_categories.md) — rule matching, category resolution, CRUD API

## Feedback
- [Workflow & Commit Rules](feedback_workflow.md) — never commit unless explicitly asked; user commits manually
- [Caveman Mode Always Active](feedback_caveman.md) — terse caveman responses every session; off only on "stop caveman"
- [Storage Pattern](feedback_storage_pattern.md) — new stores must be implemented in both sqlite + inmem or build breaks
- [Svelte Reactive Init](feedback_svelte_reactive_init.md) — never init `let` vars from `$:` labels; use `$store` directly
- [No toISOString for dates](feedback_timezone_datejs.md) — use dateUtils.js helpers; toISOString returns UTC date (wrong for UTC+)
