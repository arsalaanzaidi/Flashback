// internal/hotkey/manager.go
package hotkey

import (
	"fmt"

	"golang.design/x/hotkey"
)

// Manager holds a registered global hotkey.
type Manager struct {
	hk *hotkey.Hotkey
}

// Register binds ⌥Space globally. callback fires on each keydown.
// Call Unregister when the app shuts down.
func Register(callback func()) (*Manager, error) {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModOption}, hotkey.KeySpace)
	if err := hk.Register(); err != nil {
		return nil, fmt.Errorf("hotkey: register ⌥Space: %w", err)
	}
	m := &Manager{hk: hk}
	go func() {
		for range hk.Keydown() {
			callback()
		}
	}()
	return m, nil
}

// Unregister releases the hotkey. Safe to call on shutdown.
func (m *Manager) Unregister() {
	if m != nil && m.hk != nil {
		m.hk.Unregister()
	}
}
