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
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a WASM runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		pool: NewPool(10),
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

	// Initialize the pool
	if err := r.pool.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	return nil
}

// Execute runs WASM bytecode
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	// Execute with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := worker.Execute(code, args...)
		resultChan <- result{value: res, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.value, res.err
	}
}

// Call invokes a WASM exported function
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	// Call with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := worker.Call(fn, args...)
		resultChan <- result{value: res, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.value, res.err
	}
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
	return "wasm"
}

// Version returns the WASM version
func (r *Runtime) Version() string {
	return "1.0 (MVP)"
}

// LoadModule loads a WASM module from bytes
func (r *Runtime) LoadModule(ctx context.Context, bytecode []byte) error {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	return worker.LoadModule(bytecode)
}

type result struct {
	value interface{}
	err   error
}
