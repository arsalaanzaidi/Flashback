// internal/store/db_test.go
package store_test

import (
	"testing"
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
