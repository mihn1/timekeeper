---
name: Always implement new stores in both sqlite and inmem
description: Any new store interface must be added to both sqlite and inmem implementations or build breaks
type: feedback
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
When adding a new `XxxStore` interface, implement it in BOTH:
- `internal/data/sqlite/xxx_store.go`
- `internal/data/inmem/xxx_store.go`

AND wire it in:
- `internal/data/interfaces/storage.go` (add to interface)
- `internal/data/sqlite/storage.go` (add field + accessor)
- `internal/data/inmem/storage.go` (add field + accessor)

**Why:** Go interface compliance is checked at compile time. Missing one causes build failure immediately. Happened with Goals store and Preferences store.

**How to apply:** Treat it as a 4-file atomic change whenever adding a new store.
