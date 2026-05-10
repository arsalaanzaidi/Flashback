// internal/store/cache_test.go
package store_test

import (
	"fmt"
	"testing"
	"clipboard-manager/internal/store"
)

func TestCache_PrependAndSnapshot(t *testing.T) {
	c := store.NewCache()
	c.Prepend(store.Item{ID: "a", Content: "first"})
	c.Prepend(store.Item{ID: "b", Content: "second"})

	snap := c.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 items, got %d", len(snap))
	}
	if snap[0].ID != "b" {
		t.Fatal("most recent item should be first")
	}
}

func TestCache_CapAt50(t *testing.T) {
	c := store.NewCache()
	for i := 0; i < 60; i++ {
		c.Prepend(store.Item{ID: fmt.Sprintf("%d", i), Content: "x"})
	}
	if len(c.Snapshot()) != 50 {
		t.Fatalf("cache should cap at 50")
	}
}

func TestCache_DeduplicatesOnPrepend(t *testing.T) {
	c := store.NewCache()
	c.Prepend(store.Item{ID: "x", Content: "original"})
	c.Prepend(store.Item{ID: "x", Content: "updated"})

	snap := c.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 item after dedup, got %d", len(snap))
	}
	if snap[0].Content != "updated" {
		t.Fatal("should contain updated version")
	}
}

func TestCache_UpdateType(t *testing.T) {
	c := store.NewCache()
	c.Prepend(store.Item{ID: "z", Type: store.TypeText})
	c.UpdateType("z", store.TypeURL, "")

	snap := c.Snapshot()
	if snap[0].Type != store.TypeURL {
		t.Fatal("UpdateType did not update cache")
	}
}

func TestCache_Remove(t *testing.T) {
	c := store.NewCache()
	c.Prepend(store.Item{ID: "r1"})
	c.Prepend(store.Item{ID: "r2"})
	c.Remove("r1")

	snap := c.Snapshot()
	if len(snap) != 1 || snap[0].ID != "r2" {
		t.Fatal("Remove failed")
	}
}
