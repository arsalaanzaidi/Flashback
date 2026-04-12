#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

# ── Dependency checks ──────────────────────────────────────────────────────────
if ! command -v create-dmg &>/dev/null; then
  echo "Error: create-dmg not found. Install with: brew install create-dmg" >&2
  exit 1
fi

if ! [ -x "$HOME/go/bin/wails" ]; then
  echo "Error: wails not found at ~/go/bin/wails" >&2
  exit 1
fi

# ── Read version from Info.plist ───────────────────────────────────────────────
VERSION=$(python3 -c "
import plistlib
with open('build/darwin/Info.plist', 'rb') as f:
    p = plistlib.load(f)
print(p['CFBundleShortVersionString'])
")

DMG_NAME="Flashback ${VERSION}.dmg"
DMG_PATH="build/${DMG_NAME}"

echo "Building Flashback v${VERSION}..."

# ── 1. Kill running instance ───────────────────────────────────────────────────
echo "→ Stopping running app..."
pkill -x "flashback" 2>/dev/null || true
sleep 0.5

# ── 2. Rebuild ─────────────────────────────────────────────────────────────────
echo "→ Building app (this takes ~2 min)..."
"$HOME/go/bin/wails" build -skipbindings

# ── 3. Remove stale DMG ────────────────────────────────────────────────────────
if [ -f "$DMG_PATH" ]; then
  echo "→ Removing existing DMG..."
  rm "$DMG_PATH"
fi

# ── 4. Package ─────────────────────────────────────────────────────────────────
echo "→ Creating DMG..."
create-dmg \
  --volname "Flashback" \
  --window-pos 200 120 \
  --window-size 540 380 \
  --icon-size 128 \
  --icon "flashback.app" 160 180 \
  --hide-extension "flashback.app" \
  --app-drop-link 380 180 \
  "$DMG_PATH" \
  "build/bin/"

echo "✓ Done: ${DMG_PATH}"
