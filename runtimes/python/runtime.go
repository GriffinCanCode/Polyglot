//go:build runtime_python
// +build runtime_python

package python

// #cgo pkg-config: python3-embed
// #include <Python.h>
// #include <stdlib.h>
import "C"

import (
	"context"
	"fmt"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

var (
	// Global initialization lock
	initMu sync.Mutex
	// Track if Python is initialized
	initialized bool
)

// Runtime implements Python runtime integration with proper GIL management
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a Python runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		pool:     NewPool(4),
		shutdown: false,
	}
}

// Initialize prepares the Python runtime with proper threading support
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return ErrShutdown
	}

	initMu.Lock()
	defer initMu.Unlock()

	// Initialize Python interpreter once
	if !initialized {
		if C.Py_IsInitialized() == 0 {
			C.Py_Initialize()

			// Note: In Python 3.7+, threading is automatically initialized
			// PyEval_InitThreads() is deprecated and removed in Python 3.9+
			// The GIL is created automatically when Py_Initialize() is called

			// Save main thread state and release GIL
			// This allows other threads to acquire it
			mainThreadState := C.PyEval_SaveThread()
			if mainThreadState == nil {
				return fmt.Errorf("failed to save main thread state")
			}
		}
		initialized = true
	}

	r.config = config

	// Determine pool size
	poolSize := config.MaxConcurrency
	if poolSize <= 0 {
		poolSize = 4
	}

	// Initialize the state pool
	if err := r.pool.Initialize(poolSize); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	return nil
}

// Execute runs Python code with proper GIL management
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, ErrShutdown
	}
	r.mu.RUnlock()

	state := r.pool.Acquire()
	if state == nil {
		return nil, fmt.Errorf("failed to acquire state")
	}
	defer r.pool.Release(state)

	// Execute with context cancellation support
	resultChan := make(chan Result, 1)
	go func() {
		result, err := state.Execute(code, args...)
		resultChan <- Result{Value: result, Err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.Value, res.Err
	}
}

// Call invokes a Python function with proper GIL management
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, ErrShutdown
	}
	r.mu.RUnlock()

	state := r.pool.Acquire()
	if state == nil {
		return nil, fmt.Errorf("failed to acquire state")
	}
	defer r.pool.Release(state)

	// Call with context cancellation support
	resultChan := make(chan Result, 1)
	go func() {
		result, err := state.Call(fn, args...)
		resultChan <- Result{Value: result, Err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.Value, res.Err
	}
}

// Shutdown stops the runtime and cleans up resources
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true

	// Close pool and cleanup all states
	r.pool.Close()

	initMu.Lock()
	defer initMu.Unlock()

	// Only finalize if we initialized
	if initialized {
		// Note: In a multi-threaded environment, calling Py_Finalize() can cause
		// crashes if any Python objects are still referenced. It's safer to just
		// not call Py_Finalize() and let the process cleanup handle it.
		// The pool cleanup above already releases all our thread states.

		// Mark as uninitialized
		initialized = false
	}

	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "python"
}

// Version returns the Python version
func (r *Runtime) Version() string {
	if C.Py_IsInitialized() == 0 {
		return "not initialized"
	}

	gil := AcquireGIL()
	defer gil.Release()

	cVersion := C.Py_GetVersion()
	return C.GoString(cVersion)
}
