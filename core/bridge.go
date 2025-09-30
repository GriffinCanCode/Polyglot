package core

import (
	"context"
	"fmt"
	"sync"
)

// SimpleBridge implements a basic bridge for frontend-backend communication
type SimpleBridge struct {
	functions map[string]BridgeFunc
	mu        sync.RWMutex
}

// NewBridge creates a new bridge instance
func NewBridge() *SimpleBridge {
	return &SimpleBridge{
		functions: make(map[string]BridgeFunc),
	}
}

// Register adds a callable function to the bridge
func (b *SimpleBridge) Register(name string, fn BridgeFunc) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, exists := b.functions[name]; exists {
		return fmt.Errorf("function %s already registered", name)
	}
	
	b.functions[name] = fn
	return nil
}

// Unregister removes a callable function
func (b *SimpleBridge) Unregister(name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if _, exists := b.functions[name]; !exists {
		return fmt.Errorf("function %s not found", name)
	}
	
	delete(b.functions, name)
	return nil
}

// Call invokes a registered function
func (b *SimpleBridge) Call(ctx context.Context, name string, args ...interface{}) (interface{}, error) {
	b.mu.RLock()
	fn, exists := b.functions[name]
	b.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}
	
	return fn(ctx, args...)
}

// Functions returns a list of registered function names
func (b *SimpleBridge) Functions() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	names := make([]string, 0, len(b.functions))
	for name := range b.functions {
		names = append(names, name)
	}
	return names
}
