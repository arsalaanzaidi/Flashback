# Tech Debt Design — Flashback Clipboard Manager

**Date:** 2026-04-15  
**Scope:** All identified tech debt across backend and frontend  
**Order:** Dependency-first (migrations → backend fixes → frontend polish)

---

## 1. Schema Migration System

### Problem
`db.go` uses a single `CREATE TABLE IF NOT EXISTS` block. Shipped apps need safe, incremental schema evolution. There is no way to add columns, drop indexes, or change constraints on existing user databases.

### Design

**Directory layout:**
```
internal/store/migrations/
    001_initial.sql    ← current schema extracted verbatim
    002_...sql         ← future changes land here
```

Embedded into the binary via `//go:embed migrations/*.sql` in `db.go`.

**New table** (always created first, before any migrations run):
```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    applied_at INTEGER NOT NULL
);
```

**Runner logic in `Open()`** (replaces the inline `schema` const):
1. Create `schema_migrations` table.
2. **Existing-install bootstrap:** if `schema_migrations` is empty but the `items` table already exists, record all known migration versions as applied with `applied_at = 0`. This prevents re-running migrations on existing user databases.
3. Read all embedded `*.sql` files, parse version from numeric prefix (e.g. `001`), sort ascending.
4. For each unapplied version: execute the SQL inside a transaction, then insert a row into `schema_migrations`. Fail fast on any error — do not skip.
5. Remove the old inline `schema` const entirely.

**Version extraction:** filename must start with a three-digit zero-padded number (`001`, `002`, …). The runner enforces this at startup and panics on malformed filenames so issues are caught in development, not production.

### Files changed
- `internal/store/db.go` — remove `schema` const, add embed + `runMigrations`
- `internal/store/migrations/001_initial.sql` — new file, current schema

---

## 2. `store.Upsert` — Return Full Item on Re-copy

### Problem
`Upsert` returns only `(id string, isNew bool, err error)`. On re-copy, the caller (`handleRaw`) uses its own freshly-constructed item struct to update the cache and emit `clipboard:new-item`. That struct has:
- `pinned = false` (default zero value)
- `type = TEXT` (pre-classification placeholder)
- no `created_at` from the original insert

Result: a pinned item that gets re-copied briefly appears unpinned in the frontend until the next full reload.

### Design

**New signature:**
```go
func (s *Store) Upsert(item Item) (result Item, isNew bool, err error)
```

**On re-copy (existing hash found):**
1. `UPDATE items SET copied_at = ? WHERE id = ?`
2. Fetch the full row via `GetByID(id)` — returns complete item including real `pinned`, `type`, `subtype`, `created_at`
3. Return that full row as `result`

**On new insert:**
1. Generate UUID, insert row
2. Return the caller's `item` struct with `ID` filled in as `result`

`handleRaw` in `app.go` updated to: `item, isNew, err = a.store.Upsert(item)` — one line change, all downstream uses (cache prepend, event emit) now use the authoritative DB state.

### Files changed
- `internal/store/queries.go` — `Upsert` signature + implementation
- `app.go` — update `handleRaw` to use returned `result`

---

## 3. `store.GetByID` + Fix `CopyToClipboard`

### Problem
`CopyToClipboard(id string)` calls `s.store.List(1000, 0)` and linearly scans for the matching ID. This is O(n) on every copy action, loads up to 1000 rows, and silently returns "not found" for any item beyond position 1000.

### Design

**New method:**
```go
func (s *Store) GetByID(id string) (Item, error)
```

Uses a `WHERE id = ?` point lookup — O(1), returns `sql.ErrNoRows` wrapped as a clear error if not found.

**`CopyToClipboard` rewrite:**
```go
func (a *App) CopyToClipboard(id string) error {
    it, err := a.store.GetByID(id)
    if err != nil {
        return fmt.Errorf("item %s: %w", id, err)
    }
    if it.ImagePath != "" {
        clipboard.WriteImageToClipboard(it.ImagePath)
    } else {
        clipboard.WriteTextToClipboard(it.Content)
    }
    return nil
}
```

`GetByID` is also used internally by `Upsert` on the re-copy path (section 2), so it earns its existence from two call sites.

### Files changed
- `internal/store/queries.go` — add `GetByID`
- `app.go` — rewrite `CopyToClipboard`

---

## 4. FTS5 Query Sanitization

### Problem
`store.Search` passes user input directly as `query+"*"` to FTS5. Characters like `"`, `-`, `(`, `)`, `OR`, `AND`, `NOT`, `:`, `^` are FTS5 operators. A query like `foo-bar` or `(test` causes SQLite to return an error, which is silently swallowed by the caller returning `nil, err` — the frontend sees an empty result with no feedback.

### Design

New private function `sanitizeFTSQuery(q string) string` in `queries.go`:
- Replace FTS5 special characters (`"`, `'`, `(`, `)`, `-`, `+`, `*`, `:`, `^`, `{`, `}`, `[`, `]`) with spaces
- Collapse whitespace, trim
- Append `*` for prefix matching
- Return `""` if nothing remains after sanitization

`store.Search` calls `sanitizeFTSQuery` before constructing the query. If the sanitized result is empty, return `nil, nil` (same behaviour as an empty query — callers already handle empty results).

### Files changed
- `internal/store/queries.go` — add `sanitizeFTSQuery`, call it in `Search`

---

## 5. Frontend Error Handling

### Problem
All Wails IPC calls (`CopyToClipboard`, `PinItem`, `DeleteItem`) have return types that include errors, but `useClipboardItems.ts` never inspects them. Failures are completely silent — the user sees no feedback.

### Design

**`useClipboardItems` additions:**
- `error: string | null` — current error message, null when none
- `clearError: () => void` — dismiss manually

Each action (`copyItem`, `pinItem`, `deleteItem`) wrapped in try/catch. On error, `setError(err.message ?? 'Action failed')`. Auto-clear via `setTimeout(clearError, 3000)`.

**New `ErrorToast` component** (`frontend/src/components/ErrorToast.tsx`):
- Renders a small bar above the panel footer
- CSS opacity transition: visible when `error != null`, hidden otherwise
- No library dependency — pure CSS + React

**`ClipboardState` interface** gains `error` and `clearError` fields.

**`App.tsx`:** destructure `error` and `clearError` from the hook, pass to `<ErrorToast>`.

### Files changed
- `frontend/src/hooks/useClipboardItems.ts` — add error state + try/catch wrappers
- `frontend/src/components/ErrorToast.tsx` — new component
- `frontend/src/App.tsx` — wire up ErrorToast
- `frontend/src/App.css` — add `.error-toast` styles

---

## 6. `app.go` File Split

### Problem
`app.go` is 300 lines mixing startup/shutdown, Wails IPC handlers, the clipboard pipeline goroutines, and utility helpers. It will grow as features are added.

### Design

Split into two files in the `main` package (no new packages, just file-level organization):

**`app.go`** — keeps:
- `App` struct + `NewApp`
- `startup` / `shutdown`
- All Wails-bound IPC methods (`GetItems`, `SearchItems`, `CopyToClipboard`, `PinItem`, `DeleteItem`, `ClearHistory`, `GetImageBase64`, `GetSettings`, `SaveSettings`)
- `nowMs` helper

**`pipeline.go`** — extracted to:
- `processPipeline`
- `handleRaw`
- `utiToExt`

Both files remain `package main`. No interface changes, no new types — purely a readability split done while touching `app.go` for the `handleRaw` fix in section 2.

### Files changed
- `app.go` — remove pipeline methods + `utiToExt`
- `pipeline.go` — new file with extracted methods

---

## Testing Strategy

Each change has corresponding test coverage:

| Change | Test approach |
|--------|--------------|
| Migration runner | `db_test.go`: fresh DB runs all migrations; existing-install bootstrap path tested by pre-creating `items` table without `schema_migrations` |
| `Upsert` signature | `queries_test.go`: assert re-copy returns full item with original `pinned=true` preserved |
| `GetByID` | `queries_test.go`: happy path + not-found returns wrapped `sql.ErrNoRows` |
| FTS sanitization | `queries_test.go`: special-char queries return results (or empty) instead of erroring |
| Frontend error handling | `useClipboardItems.test.ts`: mock IPC call to reject, assert `error` is set and auto-clears |
| `ErrorToast` | Covered by hook test; no separate component test needed |

---

## Execution Order

1. Migration system (foundational — everything else may add migrations)
2. `GetByID` method (needed by Upsert fix and CopyToClipboard fix)
3. `Upsert` signature change + `handleRaw` update
4. `CopyToClipboard` fix
5. FTS5 sanitization
6. `app.go` / `pipeline.go` split
7. Frontend error handling
