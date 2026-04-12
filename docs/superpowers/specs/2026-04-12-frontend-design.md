---
title: Clipboard Manager — Frontend Design Spec
date: 2026-04-12
status: approved
---

# Clipboard Manager Frontend

## Overview

A macOS-native clipboard history panel built with React + TypeScript on Wails v2. Fixed 900×560, frameless, dark (#141414), starts hidden, shows on `option+space`. The entire frontend is a single view — no routing.

## Layout

Three conceptual zones rendered inside one 360px-wide panel (the Wails window is 900px wide but the panel is centred):

```
[Expand Tooltip] | [Main Panel 360px] | [Settings Panel — slides in over main]
```

The **expand tooltip** is not a separate DOM panel — it's an absolutely-positioned overlay that appears to the left of the selected row, showing the full content of the currently-selected item. It animates in on selection change and out when the panel is dismissed.

The **settings panel** slides in from the right over the main panel (same 360px width), triggered by `⌘,` or the ⚙ button in the footer.

## Main Panel Structure

```
┌─────────────────────────────────┐
│ SearchBar                       │  ← always visible, auto-focused on show
├─────────────────────────────────┤
│ PINNED section label            │
│ ClipboardItem (pinned=true) … │
│ ────────────────────────────── │
│ RECENT section label            │
│ ClipboardItem (selected)        │  ← highlighted, expand tooltip visible
│ ClipboardItem                   │
│ ClipboardItem                   │
│ …                               │
├─────────────────────────────────┤
│ Footer: count · ↑↓ ↵ · ⚙      │
└─────────────────────────────────┘
```

## Components

### `SearchBar`
- Single `<input>` — auto-focused on every window show event
- Debounced (150ms) → calls `SearchItems(query)` or falls back to cache snapshot on empty
- `⌘K` focuses it from anywhere; `Escape` clears text then dismisses window if already empty
- Displays a subtle spinner while fetching

### `ClipboardItem`
- Row: `[TypeBadge] [ContentPreview] [Timestamp]`
- **Unselected:** timestamp visible, actions hidden
- **Selected:** left accent border, actions (📌 pin, ✕ delete) appear; timestamp replaced by "✓ Copied" flash for 1.2s after `↵`
- Images show a 36×24 thumbnail inline (from `thumbBase64`)
- Colors show a 16×16 swatch inline
- Content preview is single-line, ellipsised. Monospace font for CODE/JSON/XML/YAML/SQL/JWT/BASE64/HASH/UUID

### `TypeBadge`
- Pill: icon + short label, colour-coded background/text per type
- 23 types mapped to 8 colour groups (see Colour System below)

### `ExpandTooltip`
- Absolutely positioned, appears to the left of the selected row
- Content varies by type:
  - **Text/Code/JSON/etc.** → scrollable `<pre>` with full content, char count + age in footer
  - **Image** → full-size thumbnail, filename + dimensions in footer
  - **Color** → large swatch + HEX/RGB/HSL table
- Animates in with a short fade+slide (100ms). Hidden when search is focused or no item selected.

### `SettingsPanel`
- Slides in from right, same dimensions as main panel
- Fields: retention mode (select: unlimited/count/days), retention value (number), launch at login (toggle)
- Save button calls `SaveSettings()`. Escape / ⌘, closes it.

### `App` (root)
- Holds all state: `items`, `selectedIndex`, `query`, `settingsOpen`
- Subscribes to `clipboard:new-item` and `clipboard:type-updated` Wails events on mount
- Listens for Wails `wails:window-show` to auto-focus search and reset selectedIndex to 0
- Keyboard handler attached to `window`: `↑↓`, `↵`, `⌘K`, `⌘P`, `⌘⌫`, `⌘,`, `Esc`

## Data Flow

```
Wails backend
  │  GetItems(50, 0)        ← initial load (from cache, zero-latency)
  │  SearchItems(query)     ← on debounced search input
  │  CopyToClipboard(id)    ← on ↵ (selected item)
  │  PinItem(id, pinned)    ← on ⌘P or pin button
  │  DeleteItem(id)         ← on ⌘⌫ or delete button
  │  GetSettings()          ← on settings panel open
  │  SaveSettings(cfg)      ← on settings save
  │
  │  Event: clipboard:new-item     → prepend to items list
  │  Event: clipboard:type-updated → patch type/subtype on matching item
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move selection |
| `↵` | Copy selected item to clipboard |
| `⌘K` | Focus search bar |
| `⌘P` | Pin / unpin selected item |
| `⌘⌫` | Delete selected item |
| `⌘,` | Open / close settings |
| `Esc` | Clear search → dismiss window |

## Colour System (TypeBadge)

| Group | Colour | Types |
|-------|--------|-------|
| Code | green `#4ade80` | CODE, JSON, XML, YAML, SQL |
| Link | blue `#60a5fa` | URL |
| Secret | amber `#fbbf24` | API_KEY, JWT, SSH_KEY, HASH |
| Identity | red `#f87171` | EMAIL, PHONE, IP |
| Media | purple `#c084fc` | IMAGE, PDF, RTF, HTML |
| Color | yellow `#facc15` | COLOR, COLOR_CODE |
| File | teal `#34d399` | FILE_REF, FILE_PATH |
| Text | grey `#888` | TEXT, UUID, BASE64, MARKDOWN |

## Styling

- Background: `#141414` (panel), `#161616` (header/footer/search)
- Border: `1px solid #2a2a2a`
- Selected row accent: `2px solid` (type colour)
- Font: `-apple-system` for UI, `SF Mono / monospace` for code previews
- No external CSS framework — plain CSS modules or inline styles
- Transitions: 100ms ease for selection, tooltip fade, settings slide

## Dependencies to Add

- None required. React 18 + TypeScript already present. Plain CSS is sufficient.
- Optional: `framer-motion` for the settings panel slide if desired (decide at implementation time)

## Out of Scope

- Type filter shortcuts (⌘1–8) — deferred
- Clear history (⌘⇧⌫) — deferred
- Customisable global shortcut UI — settings shows current value read-only
