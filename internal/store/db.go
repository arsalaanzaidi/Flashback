// internal/store/db.go
package store

import (
	"database/sql"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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
	db.SetMaxOpenConns(1)

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

	files, err := loadMigrations()
	if err != nil {
		return nil, fmt.Errorf("store: load migrations: %w", err)
	}
	if err = runMigrations(db, files); err != nil {
		return nil, fmt.Errorf("store: run migrations: %w", err)
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

// scanItems is the shared row scanner used by List, Search, GetByID, GetPinned.
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

// ── migration helpers ────────────────────────────────────────────────────────

func loadMigrations() (map[int]string, error) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}
	files := make(map[int]string, len(entries))
	for _, e := range entries {
		name := e.Name()
		underscoreIdx := strings.Index(name, "_")
		if underscoreIdx < 1 {
			return nil, fmt.Errorf("invalid migration filename %q: must start with NNN_", name)
		}
		version, err := strconv.Atoi(name[:underscoreIdx])
		if err != nil {
			return nil, fmt.Errorf("invalid migration filename %q: %w", name, err)
		}
		data, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return nil, err
		}
		files[version] = string(data)
	}
	return files, nil
}

func runMigrations(db *sql.DB, files map[int]string) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    INTEGER PRIMARY KEY,
		applied_at INTEGER NOT NULL
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	// Existing-install bootstrap: if tracking table is empty but items exists,
	// this is a pre-migration DB — mark all migrations as already applied.
	var migrationCount int
	db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&migrationCount)
	if migrationCount == 0 {
		var itemsExists int
		db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='items'").Scan(&itemsExists)
		if itemsExists > 0 {
			for version := range files {
				db.Exec("INSERT OR IGNORE INTO schema_migrations (version, applied_at) VALUES (?, 0)", version)
			}
			return nil
		}
	}

	versions := sortedVersions(files)
	for _, v := range versions {
		var applied int
		db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", v).Scan(&applied)
		if applied > 0 {
			continue
		}
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for migration %d: %w", v, err)
		}
		if _, err = tx.Exec(files[v]); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %d: %w", v, err)
		}
		if _, err = tx.Exec(
			"INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)",
			v, time.Now().UnixMilli(),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %d: %w", v, err)
		}
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", v, err)
		}
	}
	return nil
}

func sortedVersions(files map[int]string) []int {
	versions := make([]int, 0, len(files))
	for v := range files {
		versions = append(versions, v)
	}
	sort.Ints(versions)
	return versions
}
