//go:build runtime_cpp
// +build runtime_cpp

package cpp

/*
#cgo CXXFLAGS: -std=c++17
#cgo LDFLAGS: -lstdc++
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"fmt"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements C++ runtime integration via CGO
type Runtime struct {
	config   core.RuntimeConfig
	loader   *Loader
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a C++ runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		loader: NewLoader(),
	}
}

// Initialize prepares the C++ runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Load C++ shared library if specified
	if libPath, ok := config.Options["library_path"].(string); ok {
		if err := r.loader.Load(libPath); err != nil {
			return fmt.Errorf("failed to load library: %w", err)
		}
	}

	// Call initialization function if specified
	if initFn, ok := config.Options["init_function"].(string); ok {
		if _, err := r.Call(ctx, initFn); err != nil {
			return fmt.Errorf("initialization failed: %w", err)
		}
	}

	return nil
}

// Execute runs C++ code (via pre-compiled functions)
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	// For C++, "code" is typically a function name
	return r.Call(ctx, code, args...)
}

// Call invokes a C++ function by symbol name
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	return r.loader.Invoke(fn, args...)
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true

	// Call cleanup function if specified
	if cleanupFn, ok := r.config.Options["cleanup_function"].(string); ok {
		r.loader.Invoke(cleanupFn)
	}

	return r.loader.Unload()
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "cpp"
}

// Version returns the C++ standard version
func (r *Runtime) Version() string {
	// Try to call version function if available
	if result, err := r.loader.Invoke("cpp_version"); err == nil {
		if version, ok := result.(string); ok {
			return version
		}
	}
	return "c++17"
}
