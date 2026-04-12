// internal/retention/policy.go
package retention

import (
	"time"
	"clipboard-manager/internal/store"
)

type Policy struct {
	store *store.Store
	mode  string
	value int
}

// New creates a Policy. mode is "unlimited" | "count" | "days". value is N.
func New(s *store.Store, mode string, value int) *Policy {
	return &Policy{store: s, mode: mode, value: value}
}

// Update applies new settings to the live policy (call after SaveSettings).
func (p *Policy) Update(mode string, value int) {
	p.mode = mode
	p.value = value
}

// Enforce evicts items that exceed the current retention policy.
// Pinned items are never evicted.
func (p *Policy) Enforce() error {
	switch p.mode {
	case "count":
		if p.value <= 0 {
			return nil
		}
		n, err := p.store.CountNonPinned()
		if err != nil {
			return err
		}
		if n > p.value {
			return p.store.DeleteOldestNonPinned(n - p.value)
		}

	case "days":
		if p.value <= 0 {
			return nil
		}
		cutoff := time.Now().AddDate(0, 0, -p.value).UnixMilli()
		return p.store.DeleteOlderThan(cutoff)
	}
	// "unlimited" — no-op
	return nil
}
