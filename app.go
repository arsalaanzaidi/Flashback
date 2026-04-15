// app.go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"clipboard-manager/internal/clipboard"
	"clipboard-manager/internal/hotkey"
	"clipboard-manager/internal/retention"
	"clipboard-manager/internal/store"
)

const (
	eventNewItem     = "clipboard:new-item"
	eventTypeUpdated = "clipboard:type-updated"
)

type App struct {
	ctx     context.Context
	store   *store.Store
	cache   *store.Cache
	watcher *clipboard.Watcher
	hk      *hotkey.Manager
	policy  *retention.Policy
	rawCh   chan clipboard.RawItem
	workers chan struct{} // semaphore for classifier goroutine pool (8)
}

func NewApp() *App {
	return &App{
		rawCh:   make(chan clipboard.RawItem, 64),
		workers: make(chan struct{}, 8),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Open store
	dbPath := filepath.Join(
		os.Getenv("HOME"),
		"Library", "Application Support", "clipboard-manager", "data.db",
	)
	s, err := store.Open(dbPath)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	a.store = s
	a.cache = store.NewCache()

	// Seed cache from DB
	items, _ := s.List(50, 0)
	for i := len(items) - 1; i >= 0; i-- {
		a.cache.Prepend(items[i])
	}

	// Retention policy
	cfg := s.GetSettings()
	a.policy = retention.New(s, cfg.RetentionMode, cfg.RetentionValue)

	// Clipboard watcher
	a.watcher = clipboard.NewWatcher()
	go a.watcher.Start(ctx, a.rawCh)
	go a.processPipeline(ctx)

	// Global hotkey
	a.hk, err = hotkey.Register(func() {
		runtime.WindowShow(ctx)
		runtime.WindowSetAlwaysOnTop(ctx, true)
		runtime.EventsEmit(ctx, "wails:window-show")
	})
	if err != nil {
		log.Printf("hotkey: %v (continuing without global shortcut)", err)
	}
}

func (a *App) shutdown(_ context.Context) {
	if a.hk != nil {
		a.hk.Unregister()
	}
	if a.store != nil {
		a.store.Close()
	}
}


// ─── Wails-bound methods (callable from React frontend) ─────────────────────

// GetItems returns up to limit items starting at offset.
// offset=0 is served from the hot cache for zero-latency first paint.
func (a *App) GetItems(limit, offset int) []store.Item {
	if offset == 0 {
		snap := a.cache.Snapshot()
		if limit > 0 && len(snap) > limit {
			return snap[:limit]
		}
		return snap
	}
	items, _ := a.store.List(limit, offset)
	return items
}

// SearchItems queries FTS5. Empty query returns the cache snapshot.
func (a *App) SearchItems(query string) []store.Item {
	if strings.TrimSpace(query) == "" {
		return a.cache.Snapshot()
	}
	items, _ := a.store.Search(query)
	return items
}

// CopyToClipboard writes the item identified by id back to NSPasteboard.
func (a *App) CopyToClipboard(id string) error {
	it, err := a.store.GetByID(id)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	if it.ImagePath != "" {
		clipboard.WriteImageToClipboard(it.ImagePath)
	} else {
		clipboard.WriteTextToClipboard(it.Content)
	}
	return nil
}

// PinItem pins or unpins an item.
func (a *App) PinItem(id string, pinned bool) error {
	return a.store.Pin(id, pinned)
}

// DeleteItem removes an item from history.
func (a *App) DeleteItem(id string) error {
	a.cache.Remove(id)
	return a.store.Delete(id)
}

// ClearHistory removes all non-pinned items.
func (a *App) ClearHistory() error {
	return a.store.ClearNonPinned()
}

// GetImageBase64 reads the full image from disk and returns a base64 data URL.
func (a *App) GetImageBase64(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("image read: %w", err)
	}
	ext := strings.ToLower(filepath.Ext(imagePath))
	mime := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".gif":
		mime = "image/gif"
	case ".tiff":
		mime = "image/tiff"
	case ".heic":
		mime = "image/heic"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

// GetSettings returns the current settings.
func (a *App) GetSettings() store.Settings {
	return a.store.GetSettings()
}

// SaveSettings persists settings and updates the live retention policy.
func (a *App) SaveSettings(cfg store.Settings) error {
	if err := a.store.SaveSettings(cfg); err != nil {
		return err
	}
	a.policy.Update(cfg.RetentionMode, cfg.RetentionValue)
	return nil
}

// ─── helpers ────────────────────────────────────────────────────────────────

func nowMs() int64 {
	return time.Now().UnixMilli()
}

