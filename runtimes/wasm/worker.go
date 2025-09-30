//go:build runtime_wasm
// +build runtime_wasm

package wasm

import (
	"fmt"
	"sync"
)

// Worker represents a WASM execution context
type Worker struct {
	id       int
	engine   *Engine
	mu       sync.Mutex
	shutdown bool
}

// NewWorker creates a WASM worker
func NewWorker(id int) *Worker {
	return &Worker{
		id:     id,
		engine: NewEngine(),
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Initialize the WASM engine
	if err := w.engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize engine: %w", err)
	}

	return nil
}

// Execute runs WASM bytecode
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Load WASM module from code (assumed to be path or bytecode)
	module, err := w.engine.LoadModule([]byte(code))
	if err != nil {
		return nil, fmt.Errorf("failed to load module: %w", err)
	}

	// Execute the module's main or start function
	return w.engine.Execute(module, args...)
}

// Call invokes a WASM exported function
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	return w.engine.CallFunction(fn, args...)
}

// LoadModule loads a WASM module from bytes
func (w *Worker) LoadModule(bytecode []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	_, err := w.engine.LoadModule(bytecode)
	return err
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.shutdown && w.engine != nil {
		w.engine.Shutdown()
		w.engine = nil
	}

	w.shutdown = true
}
