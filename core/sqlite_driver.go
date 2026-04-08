//go:build !nosqlite

package core

// Register the pure-Go SQLite driver (modernc.org/sqlite).
// No C compiler required. Excluded when building with -tags nosqlite
// (e.g. the in-memory debug profile).
import _ "modernc.org/sqlite"
