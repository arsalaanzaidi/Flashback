// internal/store/queries_test.go
package store_test

import (
	"testing"
	"time"
	"clipboard-manager/internal/store"
)

func openTestStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestUpsert_NewItem(t *testing.T) {
	s := openTestStore(t)
	item := store.Item{
		Content:     "hello world",
		ContentHash: "abc123",
		Type:        store.TypeText,
		CopiedAt:    time.Now().UnixMilli(),
		CreatedAt:   time.Now().UnixMilli(),
		CharCount:   11,
	}
	id, isNew, err := s.Upsert(item)
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if !isNew {
		t.Fatal("expected isNew=true")
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
}

func TestUpsert_DeduplicatesOnHash(t *testing.T) {
	s := openTestStore(t)
	item := store.Item{
		Content: "dup", ContentHash: "dup123", Type: store.TypeText,
		CopiedAt: 1000, CreatedAt: 1000, CharCount: 3,
	}
	id1, isNew1, _ := s.Upsert(item)
	item.CopiedAt = 2000
	id2, isNew2, _ := s.Upsert(item)

	if !isNew1 {
		t.Fatal("first insert should be new")
	}
	if isNew2 {
		t.Fatal("second insert should not be new")
	}
	if id1 != id2 {
		t.Fatal("id should be stable across dedup upsert")
	}

	// Verify copied_at was bumped
	items, _ := s.List(10, 0)
	if len(items) != 1 || items[0].CopiedAt != 2000 {
		t.Fatalf("expected copied_at=2000, got %v", items)
	}
}

func TestList_OrderedByPinnedThenCopiedAt(t *testing.T) {
	s := openTestStore(t)
	for _, it := range []store.Item{
		{Content: "a", ContentHash: "h1", Type: store.TypeText, CopiedAt: 100, CreatedAt: 100},
		{Content: "b", ContentHash: "h2", Type: store.TypeText, CopiedAt: 200, CreatedAt: 200},
		{Content: "c", ContentHash: "h3", Type: store.TypeText, CopiedAt: 300, CreatedAt: 300},
	} {
		s.Upsert(it)
	}
	// pin the oldest
	items, _ := s.List(10, 0)
	s.Pin(items[len(items)-1].ID, true)

	items, _ = s.List(10, 0)
	if items[0].Content != "a" {
		t.Fatalf("pinned item should be first, got %q", items[0].Content)
	}
}

func TestSearch_FTS5(t *testing.T) {
	s := openTestStore(t)
	s.Upsert(store.Item{Content: "NSPasteboard changeCount", ContentHash: "h1", Type: store.TypeText, CopiedAt: 1, CreatedAt: 1, CharCount: 24})
	s.Upsert(store.Item{Content: "hello world",             ContentHash: "h2", Type: store.TypeText, CopiedAt: 2, CreatedAt: 2, CharCount: 11})

	results, err := s.Search("pasteboard")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 || results[0].Content != "NSPasteboard changeCount" {
		t.Fatalf("expected 1 result, got %v", results)
	}
}

func TestDelete(t *testing.T) {
	s := openTestStore(t)
	_, _, _ = s.Upsert(store.Item{Content: "del", ContentHash: "hd", Type: store.TypeText, CopiedAt: 1, CreatedAt: 1})
	items, _ := s.List(10, 0)
	if len(items) == 0 {
		t.Fatal("expected item before delete")
	}
	if err := s.Delete(items[0].ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	items, _ = s.List(10, 0)
	if len(items) != 0 {
		t.Fatal("expected empty after delete")
	}
}

func TestUpdateType(t *testing.T) {
	s := openTestStore(t)
	id, _, _ := s.Upsert(store.Item{Content: "fn main(){}", ContentHash: "hc", Type: store.TypeText, CopiedAt: 1, CreatedAt: 1})
	if err := s.UpdateType(id, store.TypeCode, "go"); err != nil {
		t.Fatalf("UpdateType: %v", err)
	}
	items, _ := s.List(1, 0)
	if items[0].Type != store.TypeCode || items[0].Subtype != "go" {
		t.Fatalf("UpdateType did not persist: %+v", items[0])
	}
}
