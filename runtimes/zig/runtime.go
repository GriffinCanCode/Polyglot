//go:build runtime_zig
// +build runtime_zig

package zig

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"fmt"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements Zig runtime integration via shared libraries
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a Zig runtime instance
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize prepares the Zig runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Initialize worker pool
	poolSize := config.MaxConcurrency
	if poolSize <= 0 {
		poolSize = 4
	}

	r.pool = NewPool(poolSize)
	if err := r.pool.Initialize(poolSize); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	// Load Zig shared library if specified
	if libPath, ok := config.Options["library_path"].(string); ok {
		// Load library in all workers
		for i := 0; i < r.pool.Size(); i++ {
			worker := r.pool.Acquire()
			if worker != nil {
				if err := worker.LoadLibrary(libPath); err != nil {
					r.pool.Release(worker)
					return fmt.Errorf("failed to load library in worker %d: %w", i, err)
				}
				r.pool.Release(worker)
			}
		}
	}

	return nil
}

// Execute runs Zig code (via pre-compiled functions or compilation)
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	if r.pool == nil {
		return nil, fmt.Errorf("runtime not initialized")
	}

	// Acquire a worker from the pool
	worker := r.pool.Acquire()
	if worker == nil {
		return nil, fmt.Errorf("failed to acquire worker")
	}
	defer r.pool.Release(worker)

	// Execute the code
	return worker.Execute(code, args...)
}

// Call invokes a Zig function by symbol name
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	if r.pool == nil {
		return nil, fmt.Errorf("runtime not initialized")
	}

	// Acquire a worker from the pool
	worker := r.pool.Acquire()
	if worker == nil {
		return nil, fmt.Errorf("failed to acquire worker")
	}
	defer r.pool.Release(worker)

	// Call the function
	return worker.Call(fn, args...)
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true

	if r.pool != nil {
		r.pool.Close()
	}

	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "zig"
}

// Version returns the Zig version
func (r *Runtime) Version() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pool == nil {
		return "unknown"
	}

	// Try to get zig version
	worker := r.pool.Acquire()
	if worker == nil {
		return "unknown"
	}
	defer r.pool.Release(worker)

	// Try to call a version function if available in loaded library
	result, err := worker.Call("zig_version")
	if err == nil {
		if version, ok := result.(string); ok {
			return version
		}
	}

	return "unknown"
}
