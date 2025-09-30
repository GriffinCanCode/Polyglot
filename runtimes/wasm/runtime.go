//go:build runtime_wasm
// +build runtime_wasm

package wasm

import (
	"context"
	"fmt"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements WebAssembly runtime integration
type Runtime struct {
	config   core.RuntimeConfig
	engine   *Engine
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a WASM runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		engine: NewEngine(),
	}
}

// Initialize prepares the WASM runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Initialize the WASM engine
	if err := r.engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize engine: %w", err)
	}

	return nil
}

// Execute runs WASM bytecode
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	// Load WASM module from code (assumed to be path or bytecode)
	module, err := r.engine.LoadModule([]byte(code))
	if err != nil {
		return nil, fmt.Errorf("failed to load module: %w", err)
	}

	// Execute the module's main or start function
	return r.engine.Execute(module, args...)
}

// Call invokes a WASM exported function
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	return r.engine.CallFunction(fn, args...)
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true
	return r.engine.Shutdown()
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "wasm"
}

// Version returns the WASM version
func (r *Runtime) Version() string {
	return "1.0 (MVP)"
}

// LoadModule loads a WASM module from bytes
func (r *Runtime) LoadModule(ctx context.Context, bytecode []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	_, err := r.engine.LoadModule(bytecode)
	return err
}
