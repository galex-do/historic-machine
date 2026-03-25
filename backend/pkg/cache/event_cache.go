package cache

import (
	"sync"
	"time"
)

const event_cache_ttl = 60 * time.Second

type cached_payload struct {
	data       []byte
	expires_at time.Time
}

// EventCache is a TTL in-memory cache for the full events list, keyed by locale.
// Call Invalidate() on any write (create/update/delete/import).
type EventCache struct {
	mu      sync.RWMutex
	entries map[string]cached_payload
}

func NewEventCache() *EventCache {
	return &EventCache{
		entries: make(map[string]cached_payload),
	}
}

// Get returns the cached JSON bytes for the given locale and true if the entry
// is present and has not expired. Returns nil, false otherwise.
func (c *EventCache) Get(locale string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[locale]
	if !ok || time.Now().After(entry.expires_at) {
		return nil, false
	}
	return entry.data, true
}

// Set stores pre-serialised JSON bytes for the given locale.
func (c *EventCache) Set(locale string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[locale] = cached_payload{
		data:       data,
		expires_at: time.Now().Add(event_cache_ttl),
	}
}

// Invalidate removes all cached entries so the next request re-queries the DB.
func (c *EventCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]cached_payload)
}
