// internal/store/db_test.go
package store_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
	"clipboard-manager/internal/store"
)

func TestOpen_CreatesSchema(t *testing.T) {
	path := t.TempDir() + "/test.db"
	s, err := store.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	// items table must exist
	var name string
	err = s.DB().QueryRow(
		"SELECT name FROM sqlite_master WHERE type='table' AND name='items'",
	).Scan(&name)
	if err != nil || name != "items" {
		t.Fatal("items table not created")
	}

	// FTS5 virtual table must exist
	err = s.DB().QueryRow(
		"SELECT name FROM sqlite_master WHERE type='table' AND name='items_fts'",
	).Scan(&name)
	if err != nil || name != "items_fts" {
		t.Fatal("items_fts table not created")
	}
}

func TestOpen_WALEnabled(t *testing.T) {
	path := t.TempDir() + "/test.db"
	s, err := store.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	var mode string
	s.DB().QueryRow("PRAGMA journal_mode").Scan(&mode)
	if mode != "wal" {
		t.Fatalf("expected WAL, got %s", mode)
	}
}

func TestOpen_MigrationsTableCreated(t *testing.T) {
	path := t.TempDir() + "/test.db"
	s, err := store.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	var name string
	err = s.DB().QueryRow(
		"SELECT name FROM sqlite_master WHERE type='table' AND name='schema_migrations'",
	).Scan(&name)
	if err != nil || name != "schema_migrations" {
		t.Fatal("schema_migrations table not created")
	}
}

func TestOpen_ExistingInstallBootstrap(t *testing.T) {
	// Simulate an existing install: items table exists but schema_migrations does not.
	// After Open(), schema_migrations should have entries and items data must survive.
	path := t.TempDir() + "/legacy.db"

	// Bootstrap a pre-migration DB by opening with the old schema inline.
	legacySchema := `
		CREATE TABLE IF NOT EXISTS items (
			id TEXT PRIMARY KEY, content TEXT, content_hash TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL DEFAULT 'TEXT', subtype TEXT NOT NULL DEFAULT '',
			pinned INTEGER NOT NULL DEFAULT 0, copied_at INTEGER NOT NULL,
			created_at INTEGER NOT NULL, char_count INTEGER NOT NULL DEFAULT 0,
			image_path TEXT NOT NULL DEFAULT '', thumb_blob BLOB
		);
		CREATE INDEX IF NOT EXISTS idx_items_copied_at ON items(copied_at DESC);
		CREATE INDEX IF NOT EXISTS idx_items_pinned    ON items(pinned, copied_at DESC);
		CREATE INDEX IF NOT EXISTS idx_items_type      ON items(type, copied_at DESC);
		CREATE VIRTUAL TABLE IF NOT EXISTS items_fts USING fts5(
			content, content='items', content_rowid='rowid', tokenize='trigram'
		);
		CREATE TRIGGER IF NOT EXISTS items_ai AFTER INSERT ON items BEGIN
			INSERT INTO items_fts(rowid, content) VALUES (new.rowid, new.content);
		END;
		CREATE TRIGGER IF NOT EXISTS items_ad AFTER DELETE ON items BEGIN
			INSERT INTO items_fts(items_fts, rowid, content) VALUES ('delete', old.rowid, old.content);
		END;
		CREATE TRIGGER IF NOT EXISTS items_au AFTER UPDATE OF content ON items BEGIN
			INSERT INTO items_fts(items_fts, rowid, content) VALUES ('delete', old.rowid, old.content);
			INSERT INTO items_fts(rowid, content) VALUES (new.rowid, new.content);
		END;
		CREATE TABLE IF NOT EXISTS settings (id INTEGER PRIMARY KEY CHECK(id = 1), value TEXT NOT NULL);
	`
	legacyDB, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("bootstrap open: %v", err)
	}
	if _, err = legacyDB.Exec(legacySchema); err != nil {
		legacyDB.Close()
		t.Fatalf("bootstrap schema: %v", err)
	}
	if _, err = legacyDB.Exec(
		`INSERT INTO items (id,content,content_hash,type,copied_at,created_at) VALUES ('seed-1','hello','h1','TEXT',1,1)`,
	); err != nil {
		legacyDB.Close()
		t.Fatalf("bootstrap seed: %v", err)
	}
	legacyDB.Close()

	// Now open with the migration system.
	s, err := store.Open(path)
	if err != nil {
		t.Fatalf("Open on existing DB: %v", err)
	}
	defer s.Close()

	// schema_migrations must be populated (bootstrap marked migrations as applied)
	var count int
	s.DB().QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = 1").Scan(&count)
	if count != 1 {
		t.Fatal("expected schema_migrations to be bootstrapped for existing install")
	}

	// Existing data must survive
	var content string
	s.DB().QueryRow("SELECT content FROM items WHERE id = 'seed-1'").Scan(&content)
	if content != "hello" {
		t.Fatalf("existing item lost after migration bootstrap, got %q", content)
	}
}
