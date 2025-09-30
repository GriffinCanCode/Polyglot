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
	"unsafe"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements Zig runtime integration via shared libraries
type Runtime struct {
	config   core.RuntimeConfig
	loader   *Loader
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a Zig runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		loader: NewLoader(),
	}
}

// Initialize prepares the Zig runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Load Zig shared library if specified
	if libPath, ok := config.Options["library_path"].(string); ok {
		if err := r.loader.Load(libPath); err != nil {
			return fmt.Errorf("failed to load library: %w", err)
		}
	}

	return nil
}

// Execute runs Zig code (via pre-compiled functions)
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	// For Zig, "code" is typically a function name in the loaded library
	return r.Call(ctx, code, args...)
}

// Call invokes a Zig function by symbol name
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	symbol, err := r.loader.Symbol(fn)
	if err != nil {
		return nil, fmt.Errorf("symbol %s not found: %w", fn, err)
	}

	return r.invoke(symbol, args...)
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true
	return r.loader.Unload()
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "zig"
}

// Version returns the Zig version
func (r *Runtime) Version() string {
	// Try to call a version function if available
	if symbol, err := r.loader.Symbol("zig_version"); err == nil {
		if result, err := r.invoke(symbol); err == nil {
			if version, ok := result.(string); ok {
				return version
			}
		}
	}
	return "unknown"
}

// invoke calls a Zig function through FFI
func (r *Runtime) invoke(symbol unsafe.Pointer, args ...interface{}) (interface{}, error) {
	// Type conversion and FFI call logic
	// Zig has excellent C ABI compatibility

	if len(args) == 0 {
		// Call with no arguments
		type fn0 func() int64
		f := *(*fn0)(unsafe.Pointer(&symbol))
		result := f()
		return result, nil
	}

	// Handle single argument cases
	if len(args) == 1 {
		switch v := args[0].(type) {
		case int:
			type fn1 func(int) int64
			f := *(*fn1)(unsafe.Pointer(&symbol))
			result := f(v)
			return result, nil
		case int64:
			type fn1 func(int64) int64
			f := *(*fn1)(unsafe.Pointer(&symbol))
			result := f(v)
			return result, nil
		case float64:
			type fn1 func(float64) float64
			f := *(*fn1)(unsafe.Pointer(&symbol))
			result := f(v)
			return result, nil
		}
	}

	return nil, fmt.Errorf("complex argument handling not yet implemented")
}
