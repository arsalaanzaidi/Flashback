// internal/store/db.go
package store

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS items (
	id           TEXT PRIMARY KEY,
	content      TEXT,
	content_hash TEXT NOT NULL UNIQUE,
	type         TEXT NOT NULL DEFAULT 'TEXT',
	subtype      TEXT NOT NULL DEFAULT '',
	pinned       INTEGER NOT NULL DEFAULT 0,
	copied_at    INTEGER NOT NULL,
	created_at   INTEGER NOT NULL,
	char_count   INTEGER NOT NULL DEFAULT 0,
	image_path   TEXT NOT NULL DEFAULT '',
	thumb_blob   BLOB
);

CREATE INDEX IF NOT EXISTS idx_items_copied_at ON items(copied_at DESC);
CREATE INDEX IF NOT EXISTS idx_items_pinned    ON items(pinned, copied_at DESC);
CREATE INDEX IF NOT EXISTS idx_items_type      ON items(type, copied_at DESC);

CREATE VIRTUAL TABLE IF NOT EXISTS items_fts USING fts5(
	content,
	content='items',
	content_rowid='rowid',
	tokenize='trigram'
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

CREATE TABLE IF NOT EXISTS settings (
	id    INTEGER PRIMARY KEY CHECK(id = 1),
	value TEXT NOT NULL
);
`

type Store struct {
	db       *sql.DB
	imageDir string
}

func Open(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("store: mkdir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("store: open: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite: single writer

	for _, pragma := range []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous  = NORMAL",
		"PRAGMA cache_size   = -32000",
		"PRAGMA foreign_keys = ON",
	} {
		if _, err = db.Exec(pragma); err != nil {
			return nil, fmt.Errorf("store: pragma %q: %w", pragma, err)
		}
	}

	if _, err = db.Exec(schema); err != nil {
		return nil, fmt.Errorf("store: schema: %w", err)
	}

	imageDir := filepath.Join(
		os.Getenv("HOME"),
		"Library", "Application Support", "clipboard-manager", "images",
	)
	if err = os.MkdirAll(imageDir, 0755); err != nil {
		return nil, fmt.Errorf("store: image dir: %w", err)
	}

	return &Store{db: db, imageDir: imageDir}, nil
}

func (s *Store) DB() *sql.DB { return s.db }

func (s *Store) Close() error { return s.db.Close() }

// scanItems is the shared row scanner used by List, Search, GetPinned.
func scanItems(rows *sql.Rows) ([]Item, error) {
	var items []Item
	for rows.Next() {
		var it Item
		var thumb []byte
		err := rows.Scan(
			&it.ID, &it.Content, &it.ContentHash,
			&it.Type, &it.Subtype, &it.Pinned,
			&it.CopiedAt, &it.CreatedAt, &it.CharCount,
			&it.ImagePath, &thumb,
		)
		if err != nil {
			return nil, err
		}
		if thumb != nil {
			it.ThumbBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(thumb)
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func marshalJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func unmarshalJSON(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
