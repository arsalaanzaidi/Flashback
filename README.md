# Clipboard Manager

A fast, hotkey-triggered clipboard history manager for macOS.

![macOS](https://img.shields.io/badge/macOS-10.13%2B-blue) ![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8) ![React](https://img.shields.io/badge/React-19-61DAFB)

---

## What it does

Press **⌥Space** from anywhere to instantly access your clipboard history — no dock icon, no menubar clutter.

- **History** — keeps your last 50 copied items: text, code, URLs, images, and colors
- **Search** — full-text search across everything you've ever copied
- **Pin** — pin items to keep them forever, regardless of retention settings
- **Image previews** — hover any image entry to see a full-resolution preview
- **Color detection** — recognizes hex codes and CSS colors with a live swatch
- **Retention** — configurable by count or by days; pinned items are always exempt
- **Launch at login** — always ready when you need it

---

## Install

1. Download **`Clipboard Manager.dmg`** from the [Releases](../../releases) page
2. Open the DMG and drag **Clipboard Manager** into **Applications**
3. Open **Applications**, right-click **Clipboard Manager** → **Open**
4. Click **Open** in the security dialog *(required once — see note below)*
5. Press **⌥Space** to open the panel

> **Why right-click → Open?**
> This app is not signed with an Apple Developer certificate. macOS Gatekeeper blocks unsigned apps on first launch — right-clicking and choosing **Open** bypasses this check. You only need to do it once.

---

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| ⌥Space | Open / close panel |
| ↑ ↓ | Navigate items |
| Return or click | Copy selected item |
| ⌘⌫ | Delete selected item |
| ⌘P | Pin / unpin selected item |
| ← | Expand full preview (images, long text) |
| ⌘, | Open settings |
| Esc | Close panel |

---

## Build from source

**Prerequisites**

- Go 1.21+
- Node.js 18+
- Wails v2: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

**Build**

```bash
git clone <your-repo-url>
cd clipboard-manager
~/go/bin/wails build -skipbindings
```

**Run**

```bash
open build/bin/clipboard-manager.app
```

**Package as DMG** *(requires `brew install create-dmg`)*

```bash
./scripts/build-dmg.sh
```

---

## Known limitations

- **macOS only** — requires macOS 10.13+
- **Unsigned app** — Gatekeeper will warn on first launch; right-click → Open to bypass (once only)
- Clipboard access may require a permissions prompt on macOS 14+
