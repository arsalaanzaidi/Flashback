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
