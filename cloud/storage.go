package cloud

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// MemoryStorage implements an in-memory storage for testing
type MemoryStorage struct {
	mu      sync.RWMutex
	objects map[string]*StorageObject
	data    map[string][]byte
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		objects: make(map[string]*StorageObject),
		data:    make(map[string][]byte),
	}
}

// Put stores an artifact
func (s *MemoryStorage) Put(ctx context.Context, key string, data []byte, metadata map[string]string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}
	if data == nil {
		return fmt.Errorf("data is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	obj := &StorageObject{
		Key:         key,
		Size:        int64(len(data)),
		ContentType: "application/octet-stream",
		Metadata:    metadata,
		Checksum:    fmt.Sprintf("sha256-%d", len(data)),
		CreatedAt:   time.Now(),
	}

	s.objects[key] = obj
	s.data[key] = data

	return nil
}

// Get retrieves an artifact
func (s *MemoryStorage) Get(ctx context.Context, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("object not found: %s", key)
	}

	return data, nil
}

// Delete removes an artifact
func (s *MemoryStorage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.objects, key)
	delete(s.data, key)

	return nil
}

// List lists artifacts matching a prefix
func (s *MemoryStorage) List(ctx context.Context, prefix string, limit int) ([]*StorageObject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]*StorageObject, 0)
	for key, obj := range s.objects {
		if strings.HasPrefix(key, prefix) {
			results = append(results, obj)
			if len(results) >= limit && limit > 0 {
				break
			}
		}
	}

	return results, nil
}

// GetMetadata retrieves artifact metadata
func (s *MemoryStorage) GetMetadata(ctx context.Context, key string) (*StorageObject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	obj, ok := s.objects[key]
	if !ok {
		return nil, fmt.Errorf("object not found: %s", key)
	}

	return obj, nil
}
