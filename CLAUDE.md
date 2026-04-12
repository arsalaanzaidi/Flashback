# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

**Flashback** is a macOS-native clipboard history manager built with Wails (Go backend + React/TypeScript frontend). It captures, classifies, and indexes clipboard history with image support, full-text search, and a global hotkey (⌥Space) for instant access. macOS 10.13+ only; no dock icon, frameless UI.

## Commands

### Build & Run
```bash
# Full build (requires Wails CLI: go install github.com/wailsapp/wails/v2/cmd/wails@latest)
wails build -skipbindings          # Outputs: build/bin/flashback.app

# Run the compiled app
open build/bin/flashback.app

# Package as DMG (requires: brew install create-dmg)
./scripts/build-dmg.sh
```

> After every `wails build`, replace the Wails-generated icns with the proper one before packaging:
> `cp /tmp/flashback.icns build/bin/flashback.app/Contents/Resources/iconfile.icns`

### Backend (Go)
```bash
go test ./internal/...                                          # All backend tests
go test ./internal/clipboard/...                               # Single package
go test -run TestFunctionName ./internal/clipboard/...         # Single test
```

### Frontend (React/TypeScript)
```bash
cd frontend
npm install
npm run dev          # Vite dev server (frontend-only work)
npm run build        # tsc + vite build (type-checks first)
npm test             # Vitest single pass
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
- **Classification** (`internal/clipboard/classifier.go`): Runs async after item is already stored. Four-tier strategy: UTI mapping → regex patterns → structural heuristics → language detection. Produces 25+ types (`TEXT`, `URL`, `EMAIL`, `CODE`, `JSON`, `JWT`, `IMAGE`, etc.).
- **Dedup**: Text uses SHA256 of content; images use perceptual hash (phash). Re-copying the same item updates `copied_at` rather than inserting a new row.

### Storage Layer (`internal/store/`)

- **SQLite + FTS5**: `items` table with a shadow `items_fts` (trigram tokenizer) kept in sync by triggers. WAL mode, `synchronous=NORMAL`, 32MB cache.
- **LRU cache** (`cache.go`): 50-item in-memory cache seeded on startup. `GetItems(limit, 0)` is served from cache (zero-latency first paint); `offset > 0` hits SQLite.
- **Images** (`images.go`): Full images stored on disk at `~/Library/Application Support/clipboard-manager/images/<phash>.<ext>`; 128×128 PNG thumbnails stored as blobs in SQLite.
- **Data dir**: `~/Library/Application Support/clipboard-manager/`

### Wails IPC

`app.go` methods become callable from React via auto-generated bindings in `frontend/wailsjs/go/main/App.{js,d.ts}`. If you add or rename a bound method, regenerate with `wails generate module`.

Callable methods: `GetItems`, `SearchItems`, `CopyToClipboard`, `PinItem`, `DeleteItem`, `ClearHistory`, `GetImageBase64`, `GetSettings`, `SaveSettings`.

Real-time events Go → React: `clipboard:new-item`, `clipboard:type-updated`, `wails:window-show`.

### Frontend State

`frontend/src/hooks/useClipboardItems.ts` is the single state manager — holds the item list, fires Wails IPC calls, and subscribes to `EventsOn()` for server-push updates. All keyboard shortcuts live in `App.tsx` via `keydown` listeners (arrow keys, `Cmd+K`, `Cmd+P`, `Cmd+Delete`, `Escape`).

### Global Hotkey

Registered in `internal/hotkey/manager.go` via `golang.design/x/hotkey`. On trigger: `WindowShow` → `WindowSetAlwaysOnTop` → emits `wails:window-show` so React can focus the search bar. On `Escape` or focus-loss, the frontend calls `runtime.WindowHide`.

## Roadmap / Backlog

### v1.1 — Core UX Fixes
- **Search overhaul**: FTS5 trigram search is wired but broken in the UI — debounce logic, empty-state handling, and result ranking all need work
- **Window sizing**: App dimensions are hardcoded (900×560); needs responsive behaviour or at minimum better defaults for different screen sizes
- **Pinned/Favourites tab**: Separate tab for pinned items instead of mixing them at the top of the main list

### v1.2 — Polish & Power Features
- **Preview pane**: Side-by-side preview for selected item (full text, rendered markdown, syntax-highlighted code) instead of the expand overlay
- **Keyboard shortcut customisation**: Let users remap ⌥Space and in-app shortcuts from Settings
- **Richer type badges**: Color swatch inline for hex/CSS colors; language icon for code snippets
- **Quick actions**: Right-click / long-press context menu — copy, pin, delete, open URL, copy as plain text
- **Paste-and-dismiss**: Pressing Return copies the item AND pastes it into the frontmost app, then hides Flashback

### v1.3 — Intelligence & Sync
- **Smart deduplication UI**: Show merge history ("copied 5 times") rather than just the latest timestamp
- **Regex / filter search**: `type:url`, `type:code`, `pinned:true`, `after:2d` filter syntax
- **iCloud sync**: Sync history across Macs via CloudKit (text only, not images)
- **Ignore list**: Per-app suppression (e.g. don't capture 1Password, Terminal)

### v2.0 — Website & Distribution
- **Flashback landing page**: Static page within a personal projects site — hero, feature list, download button (links to GitHub releases), screenshots, feedback form
- **Auto-update**: Sparkle framework integration so users get notified of new releases in-app
- **Apple notarisation**: Sign and notarise with an Apple Developer certificate so Gatekeeper doesn't block first launch
- **Homebrew cask**: `brew install --cask flashback` distribution

### Future / Stretch
- **iOS companion**: Share clipboard items to/from iPhone via iCloud
- **Plugin system**: Let users write classifiers or actions in JS/Lua
- **Menu bar mode**: Optional menu bar icon as an alternative entry point to ⌥Space

## Key Constraints

- **macOS only**: `watcher_darwin.go` uses cgo with `#import <AppKit/AppKit.h>`. Will not compile on Linux/Windows.
- **Icon**: Wails generates a stripped-down icns from `build/appicon.png` that macOS doesn't always pick up. The proper icns (generated via `iconutil`) lives at `/tmp/flashback.icns` and must be copied into the app bundle after each Wails build. `CFBundleIconFile` in `build/darwin/Info.plist` is required for macOS to display it.
- **No Makefile**: Use the commands above directly.
