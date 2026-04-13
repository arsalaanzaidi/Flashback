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
*Goal: stability before polish*

| Feature | Notes | Priority |
|---------|-------|----------|
| **Search overhaul** | Fix FTS5 debounce logic, empty-state handling, and result ranking | High |
| **Window sizing** | Replace hardcoded 900×560 with responsive defaults | Medium |
| **Pinned/Favourites tab** | Separate tab for pinned items instead of mixing with main list | Medium |

### v1.2 — Polish & Power Features
*Goal: delight and workflow speed*

| Feature | Notes | Priority |
|---------|-------|----------|
| **Preview pane** | Side-by-side preview with syntax highlight + rendered markdown | Medium |
| **Keyboard shortcut customisation** | Remap ⌥Space and in-app shortcuts from Settings | Medium |
| **Richer type badges** | Color swatches for hex/CSS colors; language icons for code snippets | Low |
| **Quick actions menu** | Right-click context — copy, pin, delete, open URL, copy as plain text | Medium |
| **Paste-and-dismiss** | Return key copies + pastes into frontmost app then hides Flashback | High |

### v1.3 — Intelligence & Sync
*Goal: smarter history and cross-device*

| Feature | Notes | Priority |
|---------|-------|----------|
| **Smart dedup UI** | Show "copied N times" merge history instead of just latest timestamp | Medium |
| **Filter search syntax** | `type:url`, `pinned:true`, `after:2d` operators | High |
| **iCloud sync** | Text history across Macs via CloudKit (text only, not images) | Low |
| **Ignore list** | Per-app suppression (1Password, Terminal, etc.) | Medium |

### v2.0 — Distribution & Website
*Goal: ship it publicly*

| Feature | Notes | Priority |
|---------|-------|----------|
| **Landing page** | Hero, features, download, screenshots, feedback form | Medium |
| **Auto-update** | Sparkle framework integration — in-app new release notifications | High |
| **Apple notarisation** | Sign with Apple Developer cert so Gatekeeper doesn't block first launch | High |
| **Homebrew cask** | `brew install --cask flashback` distribution | Low |

### Future / Stretch
- **iOS companion**: Share clipboard items to/from iPhone via iCloud
- **Plugin system**: Let users write classifiers or actions in JS/Lua
- **Menu bar mode**: Optional menu bar icon as an alternative entry point to ⌥Space

## Silent App / Menu Bar Refactor

> **Status:** Planned — not yet implemented.

**Goal:** Remove Flashback from the Dock and app switcher entirely. Run only as a menu bar status item (like Bartender, 1Password, etc.).

**Files that will be touched:**

| File | Change |
|------|--------|
| `build/darwin/Info.plist` | Add `LSUIElement` key |
| `app.go` | Call menubar init early; remove startup `WindowShow` |
| `internal/menubar/manager.go` | New file — cgo NSStatusItem + NSMenu bridge |
| `internal/clipboard/watcher_darwin.go` | Add `Pause()` / `Resume()` methods |
| `internal/hotkey/manager.go` | Expose `Toggle()` for icon left-click |

### Phase 1 — Suppress Dock & app switcher
- Add `LSUIElement = true` to `Info.plist`
- In new `internal/menubar/manager.go`, add a cgo helper `SetAccessoryPolicy()` that calls `[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory]`
- Call `menubar.SetAccessoryPolicy()` early in `app.go` startup

### Phase 2 — Add menu bar icon
- Create `NSStatusItem` via cgo in `internal/menubar/manager.go`
- Add an 18×18pt template image asset (`menubar-icon.pdf`) to `build/darwin/Assets.xcassets/` — `image.template = YES` handles light/dark automatically
- Wire status item left-click to the existing hotkey toggle (`hotkeyManager.Toggle()`)

### Phase 3 — Dropdown context menu
- Build an `NSMenu` on right-click with: Open Flashback, separator, Pause capturing (toggle), Clear history (with confirmation), separator, Settings, About, Quit
- Add `Pause()` / `Resume()` methods to `internal/clipboard/watcher_darwin.go` and call them from the menu bridge

### Phase 4 — Cleanup & launch behaviour
- Remove any `WindowShow` from startup — window hidden by default, appears only on ⌥Space or icon click
- Override `windowShouldClose` in a cgo delegate to call `WindowHide` instead of terminating the app
- Add "Launch at login" toggle using `SMAppService` (macOS 13+) or a `LaunchAgent` plist for older versions — surface as a checkbox in Settings

## Key Constraints

- **macOS only**: `watcher_darwin.go` uses cgo with `#import <AppKit/AppKit.h>`. Will not compile on Linux/Windows.
- **Icon**: Wails generates a stripped-down icns from `build/appicon.png` that macOS doesn't always pick up. The proper icns (generated via `iconutil`) lives at `/tmp/flashback.icns` and must be copied into the app bundle after each Wails build. `CFBundleIconFile` in `build/darwin/Info.plist` is required for macOS to display it.
- **No Makefile**: Use the commands above directly.
