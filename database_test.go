package zen_test

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	zen "github.com/phenixrizen/zen-go/v2"
)

// The one thing this repository had no test for: a databaseNode evaluating through the
// SQLite handler that the vendored libzen_ffi.a already contains. The fixture .db is built
// at test time from test-data/database.sql — a .db is never committed.
func buildCatalog(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	db, err := sql.Open("sqlite", filepath.Join(root, "catalog.db"))
	require.NoError(t, err)
	defer db.Close()
	ddl, err := os.ReadFile(filepath.Join("test-data", "database.sql"))
	require.NoError(t, err)
	_, err = db.Exec(string(ddl))
	require.NoError(t, err)
	return root
}

func TestDatabaseNode_EvaluatesThroughTheSqliteHandler(t *testing.T) {
	root := buildCatalog(t)
	cfg, _ := json.Marshal(map[string]any{"root": root})
	// SqliteConfig needs a Loader CALLBACK: the FilesystemLoader form is a LoaderConfig,
	// which the FFI cannot combine with a handler config (engine.go:131).
	engine := zen.NewEngine(zen.EngineConfig{
		Loader:       zen.Loader(readTestFile),
		SqliteConfig: string(cfg),
	})
	// Dispose drops the handler and its connection pool. Without it the engine leaks
	// (the memory_test job) and, on Windows, the still-open catalog.db makes TempDir
	// cleanup fail — the file cannot be deleted while a handle is held.
	defer engine.Dispose()

	hit := result(t, engine, "36415")
	assert.Equal(t, true, hit["covered"], "a code in the table is covered")

	miss := result(t, engine, "00000")
	assert.Equal(t, false, miss["covered"], "a code absent from the table is not")

	// $params is reserved and stripped: it must never leak into a result
	_, leaked := hit["$params"]
	assert.False(t, leaked)
}

func result(t *testing.T, engine zen.Engine, code string) map[string]any {
	t.Helper()
	resp, err := engine.Evaluate("database.json", map[string]any{"procedureCode": code})
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, json.Unmarshal(resp.Result, &out))
	return out
}

func TestDatabaseNode_WithoutAHandlerIsAHardError(t *testing.T) {
	// No SqliteConfig: the node reports "Database handler not provided" and the whole
	// evaluation errors. A host that drafts database-backed graphs must install one.
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.FilesystemLoader{Path: "test-data"}})
	defer engine.Dispose()
	_, err := engine.Evaluate("database.json",
		map[string]any{"procedureCode": "36415"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "handler not provided")
}
