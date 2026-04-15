// internal/store/queries.go
package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (s *Store) Upsert(item Item) (result Item, isNew bool, err error) {
	// Check for existing hash.
	var existingID string
	err = s.db.QueryRow(
		"SELECT id FROM items WHERE content_hash = ?", item.ContentHash,
	).Scan(&existingID)

	if err == nil {
		// Exists — bump copied_at, then return the full stored row (preserves
		// pinned status, original type, created_at, etc.)
		if _, err = s.db.Exec(
			"UPDATE items SET copied_at = ? WHERE id = ?", item.CopiedAt, existingID,
		); err != nil {
			return Item{}, false, fmt.Errorf("upsert bump: %w", err)
		}
		full, err := s.GetByID(existingID)
		return full, false, err
	}
	if err != sql.ErrNoRows {
		return Item{}, false, fmt.Errorf("upsert lookup: %w", err)
	}

	// New item — generate ID, insert, return the item with ID set.
	id := uuid.New().String()
	_, err = s.db.Exec(`
		INSERT INTO items
			(id, content, content_hash, type, subtype, pinned, copied_at, created_at, char_count, image_path, thumb_blob)
		VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?)`,
		id, item.Content, item.ContentHash, item.Type, item.Subtype,
		item.CopiedAt, item.CreatedAt, item.CharCount, item.ImagePath, nilBlob(item.ThumbBlob),
	)
	if err != nil {
		return Item{}, false, err
	}
	item.ID = id
	return item, true, nil
}

func (s *Store) List(limit, offset int) ([]Item, error) {
	rows, err := s.db.Query(`
		SELECT id, content, content_hash, type, subtype, pinned, copied_at, created_at, char_count, image_path, thumb_blob
		FROM items
		ORDER BY pinned DESC, copied_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanItems(rows)
}

func (s *Store) Search(query string) ([]Item, error) {
	q := sanitizeFTSQuery(query)
	if q == "" {
		return nil, nil
	}
	rows, err := s.db.Query(`
		SELECT i.id, i.content, i.content_hash, i.type, i.subtype, i.pinned,
		       i.copied_at, i.created_at, i.char_count, i.image_path, i.thumb_blob
		FROM items_fts f
		JOIN items i ON i.rowid = f.rowid
		WHERE items_fts MATCH ?
		ORDER BY i.pinned DESC, i.copied_at DESC
		LIMIT 200`, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanItems(rows)
}

// GetByID fetches a single item by its primary key.
// Returns a wrapped sql.ErrNoRows if the id is not found.
func (s *Store) GetByID(id string) (Item, error) {
	rows, err := s.db.Query(`
		SELECT id, content, content_hash, type, subtype, pinned, copied_at, created_at, char_count, image_path, thumb_blob
		FROM items WHERE id = ? LIMIT 1`, id)
	if err != nil {
		return Item{}, err
	}
	defer rows.Close()
	items, err := scanItems(rows)
	if err != nil {
		return Item{}, err
	}
	if len(items) == 0 {
		return Item{}, fmt.Errorf("item %s: %w", id, sql.ErrNoRows)
	}
	return items[0], nil
}

func (s *Store) Pin(id string, pinned bool) error {
	p := 0
	if pinned {
		p = 1
	}
	_, err := s.db.Exec("UPDATE items SET pinned = ? WHERE id = ?", p, id)
	return err
}

func (s *Store) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}

func (s *Store) DeleteByImagePath(path string) error {
	_, err := s.db.Exec("DELETE FROM items WHERE image_path = ?", path)
	return err
}

func (s *Store) UpdateType(id, typ, subtype string) error {
	_, err := s.db.Exec(
		"UPDATE items SET type = ?, subtype = ? WHERE id = ?", typ, subtype, id,
	)
	return err
}

func (s *Store) CountNonPinned() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM items WHERE pinned = 0").Scan(&n)
	return n, err
}

func (s *Store) DeleteOldestNonPinned(n int) error {
	_, err := s.db.Exec(`
		DELETE FROM items WHERE id IN (
			SELECT id FROM items WHERE pinned = 0 ORDER BY copied_at ASC LIMIT ?
		)`, n)
	return err
}

func (s *Store) DeleteOlderThan(cutoffMs int64) error {
	_, err := s.db.Exec(
		"DELETE FROM items WHERE pinned = 0 AND copied_at < ?", cutoffMs,
	)
	return err
}

func (s *Store) ClearNonPinned() error {
	_, err := s.db.Exec("DELETE FROM items WHERE pinned = 0")
	return err
}

// Settings

func (s *Store) GetSettings() Settings {
	var raw string
	err := s.db.QueryRow("SELECT value FROM settings WHERE id = 1").Scan(&raw)
	if err != nil {
		return DefaultSettings()
	}
	var cfg Settings
	if err = unmarshalJSON(raw, &cfg); err != nil {
		return DefaultSettings()
	}
	return cfg
}

func (s *Store) SaveSettings(cfg Settings) error {
	raw, err := marshalJSON(cfg)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`
		INSERT INTO settings (id, value) VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET value = excluded.value`, raw)
	return err
}

// helpers

func nilBlob(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	return b
}

// sanitizeFTSQuery strips FTS5 operator characters from user input and
// returns a safe prefix query. Each token is quoted so FTS5 treats it as a
// literal (not an operator). The wildcard * is appended to the last token
// for prefix matching.
func sanitizeFTSQuery(q string) string {
	var b strings.Builder
	for _, r := range q {
		switch r {
		case '"', '\'', '(', ')', '-', '+', '*', ':', '^', '{', '}', '[', ']':
			b.WriteRune(' ')
		default:
			b.WriteRune(r)
		}
	}
	tokens := strings.Fields(b.String()) // splits on any whitespace run; no loop needed
	if len(tokens) == 0 {
		return ""
	}
	// Wrap each token in double-quotes (disables operator interpretation).
	// Append * outside the last closing quote for prefix matching.
	parts := make([]string, len(tokens))
	for i, tok := range tokens {
		if i == len(tokens)-1 {
			parts[i] = `"` + tok + `"*`
		} else {
			parts[i] = `"` + tok + `"`
		}
	}
	return strings.Join(parts, " ")
}
