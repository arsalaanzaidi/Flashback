// internal/store/cache.go
package store

import "sync"

const cacheCapacity = 50

type Cache struct {
	mu    sync.RWMutex
	items []Item
}

func NewCache() *Cache {
	return &Cache{items: make([]Item, 0, cacheCapacity)}
}

// Prepend adds item to the front. If an item with the same ID already exists
// it is removed first (handles re-copy / copied_at bump).
func (c *Cache) Prepend(item Item) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.removeByID(item.ID)
	c.items = append([]Item{item}, c.items...)
	if len(c.items) > cacheCapacity {
		c.items = c.items[:cacheCapacity]
	}
}

// UpdateType patches the Type and Subtype of a cached item in place.
func (c *Cache) UpdateType(id, typ, subtype string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.items {
		if c.items[i].ID == id {
			c.items[i].Type = typ
			c.items[i].Subtype = subtype
			return
		}
	}
}

// Remove deletes an item from the cache by ID.
func (c *Cache) Remove(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.removeByID(id)
}

// Snapshot returns a copy of the current cache contents (newest first).
func (c *Cache) Snapshot() []Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cp := make([]Item, len(c.items))
	copy(cp, c.items)
	return cp
}

func (c *Cache) removeByID(id string) {
	for i, item := range c.items {
		if item.ID == id {
			c.items = append(c.items[:i], c.items[i+1:]...)
			return
		}
	}
}
