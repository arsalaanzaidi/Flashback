// internal/retention/policy_test.go
package retention_test

import (
	"fmt"
	"testing"
	"time"
	"clipboard-manager/internal/retention"
	"clipboard-manager/internal/store"
)

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(t.TempDir() + "/ret.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func seedItems(t *testing.T, s *store.Store, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		s.Upsert(store.Item{
			Content:     fmt.Sprintf("item %d", i),
			ContentHash: fmt.Sprintf("hash%d", i),
			Type:        store.TypeText,
			CopiedAt:    int64(i+1) * 1000,
			CreatedAt:   int64(i+1) * 1000,
			CharCount:   6,
		})
	}
}

func TestEnforce_CountLimit(t *testing.T) {
	s := makeStore(t)
	seedItems(t, s, 10)
	// pin one item so it survives
	items, _ := s.List(10, 0)
	s.Pin(items[0].ID, true)

	p := retention.New(s, "count", 5)
	if err := p.Enforce(); err != nil {
		t.Fatalf("Enforce: %v", err)
	}

	remaining, _ := s.List(20, 0)
	// 5 non-pinned + 1 pinned = 6
	if len(remaining) != 6 {
		t.Fatalf("expected 6 items after enforcing count=5, got %d", len(remaining))
	}
	// Pinned item must survive
	found := false
	for _, it := range remaining {
		if it.Pinned {
			found = true
		}
	}
	if !found {
		t.Fatal("pinned item was evicted")
	}
}

func TestEnforce_DaysLimit(t *testing.T) {
	s := makeStore(t)
	now := time.Now().UnixMilli()
	// Insert old item (40 days ago) and recent item (1 day ago)
	s.Upsert(store.Item{Content: "old", ContentHash: "hold", Type: store.TypeText,
		CopiedAt: now - int64(40*24*time.Hour/time.Millisecond), CreatedAt: 1})
	s.Upsert(store.Item{Content: "recent", ContentHash: "hrec", Type: store.TypeText,
		CopiedAt: now - int64(1*24*time.Hour/time.Millisecond), CreatedAt: 2})

	p := retention.New(s, "days", 30)
	if err := p.Enforce(); err != nil {
		t.Fatalf("Enforce: %v", err)
	}

	items, _ := s.List(10, 0)
	if len(items) != 1 || items[0].Content != "recent" {
		t.Fatalf("expected only 'recent', got %v", items)
	}
}

func TestEnforce_Unlimited(t *testing.T) {
	s := makeStore(t)
	seedItems(t, s, 20)
	p := retention.New(s, "unlimited", 0)
	if err := p.Enforce(); err != nil {
		t.Fatalf("Enforce: %v", err)
	}
	items, _ := s.List(30, 0)
	if len(items) != 20 {
		t.Fatalf("unlimited should not evict anything, got %d", len(items))
	}
}
