package core

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// MemoryCoordinator manages shared memory across runtimes
type MemoryCoordinator struct {
	config  MemoryConfig
	regions map[string]*MemoryRegion
	usage   int64
	mu      sync.RWMutex
}

// NewMemoryCoordinator creates a memory coordinator
func NewMemoryCoordinator(config MemoryConfig) *MemoryCoordinator {
	return &MemoryCoordinator{
		config:  config,
		regions: make(map[string]*MemoryRegion),
		usage:   0,
	}
}

// Allocate creates a new shared memory region
func (m *MemoryCoordinator) Allocate(id string, size int, memType MemoryType) (*MemoryRegion, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.regions[id]; exists {
		return nil, fmt.Errorf("region %s already exists", id)
	}

	newUsage := atomic.AddInt64(&m.usage, int64(size))
	if newUsage > m.config.MaxSharedMemory {
		atomic.AddInt64(&m.usage, -int64(size))
		return nil, fmt.Errorf("memory limit exceeded")
	}

	region := &MemoryRegion{
		ID:   id,
		Data: make([]byte, size),
		Type: memType,
	}

	m.regions[id] = region
	return region, nil
}

// Get retrieves a memory region by ID
func (m *MemoryCoordinator) Get(id string) (*MemoryRegion, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	region, exists := m.regions[id]
	if !exists {
		return nil, fmt.Errorf("region %s not found", id)
	}

	return region, nil
}

// Free releases a memory region
func (m *MemoryCoordinator) Free(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	region, exists := m.regions[id]
	if !exists {
		return fmt.Errorf("region %s not found", id)
	}

	if region.Readers > 0 || region.Writers > 0 {
		return fmt.Errorf("region %s still has active users", id)
	}

	atomic.AddInt64(&m.usage, -int64(len(region.Data)))
	delete(m.regions, id)

	return nil
}

// AcquireRead marks a region as being read
func (m *MemoryCoordinator) AcquireRead(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	region, exists := m.regions[id]
	if !exists {
		return fmt.Errorf("region %s not found", id)
	}

	region.Readers++
	return nil
}

// ReleaseRead marks a region read as complete
func (m *MemoryCoordinator) ReleaseRead(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	region, exists := m.regions[id]
	if !exists {
		return fmt.Errorf("region %s not found", id)
	}

	if region.Readers > 0 {
		region.Readers--
	}

	return nil
}

// AcquireWrite marks a region as being written
func (m *MemoryCoordinator) AcquireWrite(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	region, exists := m.regions[id]
	if !exists {
		return fmt.Errorf("region %s not found", id)
	}

	if region.Writers > 0 {
		return fmt.Errorf("region %s already has a writer", id)
	}

	region.Writers++
	return nil
}

// ReleaseWrite marks a region write as complete
func (m *MemoryCoordinator) ReleaseWrite(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	region, exists := m.regions[id]
	if !exists {
		return fmt.Errorf("region %s not found", id)
	}

	if region.Writers > 0 {
		region.Writers--
	}

	return nil
}

// Usage returns current memory usage in bytes
func (m *MemoryCoordinator) Usage() int64 {
	return atomic.LoadInt64(&m.usage)
}

// Stats returns memory statistics
func (m *MemoryCoordinator) Stats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"regions":     len(m.regions),
		"usage":       m.Usage(),
		"limit":       m.config.MaxSharedMemory,
		"utilization": float64(m.Usage()) / float64(m.config.MaxSharedMemory),
	}
}
