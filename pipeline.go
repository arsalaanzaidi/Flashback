// pipeline.go
package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"clipboard-manager/internal/clipboard"
	"clipboard-manager/internal/store"
)

// processPipeline consumes rawCh: dedup → store → cache → async classify.
func (a *App) processPipeline(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case raw := <-a.rawCh:
			a.handleRaw(raw)
		}
	}
}

func (a *App) handleRaw(raw clipboard.RawItem) {
	now := nowMs()

	var (
		item  store.Item
		isImg bool
	)

	switch {
	case len(raw.Data) > 0: // binary image
		isImg = true
		ext := utiToExt(raw.UTI)
		path, thumb, phash, err := a.store.SaveImage(raw.Data, ext)
		if err != nil {
			log.Printf("pipeline: save image: %v", err)
			return
		}
		item = store.Item{
			ContentHash: phash,
			Type:        store.TypeImage,
			Subtype:     ext,
			CopiedAt:    now,
			CreatedAt:   now,
			ImagePath:   path,
			ThumbBlob:   thumb,
		}

	default: // text-based
		if raw.Text == "" {
			return
		}
		h := sha256.Sum256([]byte(raw.Text))
		item = store.Item{
			Content:     raw.Text,
			ContentHash: fmt.Sprintf("%x", h),
			Type:        store.TypeText, // updated async after classify
			CopiedAt:    now,
			CreatedAt:   now,
			CharCount:   len([]rune(raw.Text)),
		}
	}

	var isNew bool
	var err error
	item, isNew, err = a.store.Upsert(item)
	if err != nil {
		log.Printf("pipeline: upsert: %v", err)
		return
	}
	id := item.ID

	// Prepend to cache (item now has DB-authoritative fields: pinned, type, etc.)
	a.cache.Prepend(item)
	runtime.EventsEmit(a.ctx, eventNewItem, item)

	// Enforce retention after each new item
	if isNew {
		if err = a.policy.Enforce(); err != nil {
			log.Printf("pipeline: retention: %v", err)
		}
	}

	// Async classification for text items
	if !isImg {
		capturedItem := item // capture immutable copy for goroutine
		a.workers <- struct{}{}
		go func() {
			defer func() { <-a.workers }()
			result := clipboard.Classify(raw.UTI, raw.Text)
			if result.Type == capturedItem.Type && result.Subtype == capturedItem.Subtype {
				return // no change
			}
			if err := a.store.UpdateType(id, result.Type, result.Subtype); err != nil {
				log.Printf("classify update: %v", err)
				return
			}
			a.cache.UpdateType(id, result.Type, result.Subtype)
			runtime.EventsEmit(a.ctx, eventTypeUpdated, map[string]string{
				"id": id, "type": result.Type, "subtype": result.Subtype,
			})
		}()
	}
}

// utiToExt maps macOS Uniform Type Identifiers to file extensions.
func utiToExt(uti string) string {
	switch uti {
	case "public.jpeg":
		return "jpg"
	case "public.tiff":
		return "tiff"
	case "com.compuserve.gif":
		return "gif"
	case "public.heic":
		return "heic"
	default:
		return "png"
	}
}
