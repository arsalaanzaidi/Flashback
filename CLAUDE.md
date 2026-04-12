# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

**Flashback** is a macOS-native clipboard history manager built with Wails (Go backend + React/TypeScript frontend). It captures, classifies, and indexes clipboard history with image support, FTS, and a global hotkey (⌥Space) for instant access. macOS 10.13+ only; no dock icon, frameless UI.

## Commands

### Build & Run
```bash
# Full build (requires Wails CLI: go install github.com/wailsapp/wails/v2/cmd/wails@latest)
wails build -skipbindings          # Outputs: build/bin/flashback.app

# Dev: run the compiled app
open build/bin/flashback.app

# Package as DMG (requires: brew install create-dmg)
./scripts/build-dmg.sh
```

### Backend (Go)
```bash
go test ./internal/...             # All backend tests
go test ./internal/clipboard/...   # Single package
go test -run TestFunctionName ./internal/clipboard/...  # Single test
```

### Frontend (React/TypeScript)
```bash
cd frontend
npm install
npm run dev          # Vite dev server (for frontend-only work)
npm run build        # tsc + vite build (type-checks first)
npm test             # Vitest run (single pass)
npm run test:watch   # Vitest watch mode
```

## Architecture

### Clipboard Capture Pipeline

The core data flow is a channel-based pipeline in `app.go`:

```
NSPasteboard polling → rawCh (chan RawItem) → handleRaw()
  → SHA256/phash dedup → store.Upsert → cache.Prepend → EventsEmit("clipboard:new-item")
  → [async goroutine pool, 8 workers] → Classify → store.UpdateType → EventsEmit("clipboard:type-updated")
```

- **Watcher** (`internal/clipboard/watcher_darwin.go`): Uses cgo/Objective-C to poll `NSPasteboard.changeCount()` on an adaptive interval (250ms active, 2s after 5 min idle).
- **Classification** (`internal/clipboard/classifier.go`): Runs asynchronously after item is already stored. Four-tier: UTI mapping → regex patterns → structural heuristics → language detection. Produces 25+ types (`TEXT`, `URL`, `EMAIL`, `CODE`, `JSON`, `JWT`, `IMAGE`, etc.).
- **Dedup**: Text uses SHA256 of content; images use perceptual hash (phash) — re-copying the same item updates `copied_at` instead of inserting a new row.

### Storage Layer (`internal/store/`)

- **SQLite + FTS5**: `items` table with a shadow `items_fts` (trigram tokenizer) kept in sync by triggers. WAL mode, `synchronous=NORMAL`, 32MB cache.
- **LRU cache** (`cache.go`): 50-item in-memory cache seeded on startup. First `GetItems(limit, 0)` call is served from cache (zero-latency first paint); `offset > 0` hits SQLite.
- **Images** (`images.go`): Full images stored on disk at `~/Library/Application Support/clipboard-manager/images/<phash>.<ext>`; 128×128 PNG thumbnails stored as blobs in SQLite.
- **Data dir**: `~/Library/Application Support/clipboard-manager/`

### Wails IPC

`app.go` methods annotated as Go exports become callable from React via the auto-generated `frontend/wailsjs/go/main/App.js` bindings. The set of callable methods: `GetItems`, `SearchItems`, `CopyToClipboard`, `PinItem`, `DeleteItem`, `ClearHistory`, `GetImageBase64`, `GetSettings`, `SaveSettings`.

Real-time events emitted from Go → React: `clipboard:new-item`, `clipboard:type-updated`, `wails:window-show`.

### Frontend State

`frontend/src/hooks/useClipboardItems.ts` is the single state manager — it holds the item list, fires Wails IPC calls, and subscribes to `EventsOn()` for server-push updates. All keyboard shortcuts are handled in `App.tsx` via `keydown` listeners (arrow keys, `Cmd+K`, `Cmd+P`, `Cmd+Delete`, `Escape`).

### Global Hotkey

Registered in `internal/hotkey/manager.go` via `golang.design/x/hotkey`. On trigger: `WindowShow` → `WindowSetAlwaysOnTop` → emits `wails:window-show` so React can focus the search bar. On `Escape` or focus-loss, the frontend calls `runtime.WindowHide`.

## Key Constraints

- **macOS only**: `watcher_darwin.go` uses cgo with `#import <AppKit/AppKit.h>`. Will not compile on Linux/Windows.
- **Regenerating bindings**: If you add/rename a Go method bound to Wails, run `wails generate module` to update `frontend/wailsjs/go/main/App.{js,d.ts}`.
- **No Makefile**: There is no `make` target; use the commands above directly.
