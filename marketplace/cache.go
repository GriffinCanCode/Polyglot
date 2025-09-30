package marketplace

import (
	"context"
	"fmt"
	"sync"
)

// MemoryCache implements an in-memory cache for testing
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string][]byte
	size  int64
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string][]byte),
	}
}

// Get retrieves a cached package
func (c *MemoryCache) Get(ctx context.Context, id, version string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%s@%s", id, version)
	data, ok := c.items[key]
	if !ok {
		return nil, fmt.Errorf("not found in cache")
	}

	return data, nil
}

// Put stores a package in cache
func (c *MemoryCache) Put(ctx context.Context, id, version string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s@%s", id, version)

	// Remove old entry if exists
	if old, ok := c.items[key]; ok {
		c.size -= int64(len(old))
	}

	c.items[key] = data
	c.size += int64(len(data))

	return nil
}

// Has checks if a package is cached
func (c *MemoryCache) Has(ctx context.Context, id, version string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%s@%s", id, version)
	_, ok := c.items[key]
	return ok
}

// Remove removes a package from cache
func (c *MemoryCache) Remove(ctx context.Context, id, version string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s@%s", id, version)
	if data, ok := c.items[key]; ok {
		c.size -= int64(len(data))
		delete(c.items, key)
	}

	return nil
}

// Clear clears all cached packages
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string][]byte)
	c.size = 0

	return nil
}

// Size returns cache size in bytes
func (c *MemoryCache) Size(ctx context.Context) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.size, nil
}
