# Build Directory

- `bin/` — compiled app output (`flashback.app`)
- `darwin/` — macOS-specific files (`Info.plist`, icon assets)
- `windows/` — Windows manifest and installer files (unused; this app is macOS-only)

To rebuild: `wails build -skipbindings` from the repo root.
