# Execution Plan (Engineering)

This plan tracks backend/frontend reliability work prioritized for developer velocity and correctness.

## Phase 0: Guardrails Before Changes

- Define commit boundaries by recommendation area to keep review focused.
- Add/extend tests near high-risk logic before broad refactors.
- Keep generated Wails bindings untouched unless regenerated intentionally.

## Phase 1: Correctness Fixes (Recommendation 1)

Status: completed

- Fix event dedup logic to use absolute time delta.
- Fix combined rule sorting so app-specific and global rules are ordered together.
- Fix SQLite event retrieval to actually return scanned rows.
- Fix frontend category DTO casing mismatch (`id/name` vs `Id/Name`).

## Phase 2: Lifecycle Safety (Recommendation 2)

Status: completed

- Make `StartTracking()` idempotent.
- Ensure observer/event-loop startup happens once.
- Add close guard (`sync.Once`) and close signal channel.
- Add lifecycle tests for repeated enable/disable/close calls.

## Phase 3: Startup Configuration (Recommendation 3)

Status: completed

- Move desktop app runtime config to env-driven config loader (`TIMEKEEPER_DB`, `TIMEKEEPER_DB_PATH`, `TIMEKEEPER_SEED_MODE`).
- Replace hardcoded DB path.
- Seed only by explicit policy (`always`, `never`, `if-empty`).
- Ensure SQLite parent directory creation attempt before startup.

## Phase 4: Error Handling and API Contracts (Recommendation 4)

Status: completed

- Return explicit errors for dashboard data methods on invalid input/storage failure.
- Add payload validation for category/rule create and update operations.
- Improve category usage enrichment behavior by skipping missing category rows with warning logs.

## Phase 5: Test Coverage Expansion (Recommendation 5)

Status: completed

- Add tests for exclusion rule behavior.
- Add tests for priority conflict resolution.
- Add tests for disabled tracking event drop behavior.
- Add/keep lifecycle-focused tests in `core/timekeeper_test.go`.

## Phase 6: Maintainability Refactors (Recommendation 6)

Status: completed

- Add regex compilation cache for rule evaluation.
- Standardize `macos/observer.go` logging to `slog`-based logger.
- Add model-level tests for regex matching behavior.

## Follow-Up Validation Checklist

- Run backend test targets after each phase: `go test ./core ./internal/data/sqlite ./internal/models`.
- Manually validate Wails desktop flows for dashboard/category/rule views.
- Regenerate Wails bindings if method signature updates are consumed by frontend build workflows.
