#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

# ── Dependency checks ──────────────────────────────────────────────────────────
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

echo "Building Flashback v${VERSION}..."

# ── 1. Kill running instance ───────────────────────────────────────────────────
echo "→ Stopping running app..."
pkill -ix "flashback" 2>/dev/null || true
sleep 0.5

# ── 2. Rebuild ─────────────────────────────────────────────────────────────────
echo "→ Building app (this takes ~2 min)..."
"$HOME/go/bin/wails" build -skipbindings

echo "✓ Done: build/bin/flashback.app"
