//go:build !nosqlite

package core

// Register the SQLite3 CGo driver. Excluded when building with -tags nosqlite
// (e.g. the in-memory debug profile) so the binary can be built without a C compiler.
import _ "github.com/mattn/go-sqlite3"
