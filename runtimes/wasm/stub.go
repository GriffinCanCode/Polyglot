//go:build !runtime_wasm
// +build !runtime_wasm

package wasm

import (
	"context"
	"fmt"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements a stub WASM runtime when WASM is not enabled
type Runtime struct {
	config core.RuntimeConfig
}

// NewRuntime creates a stub WASM runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating WASM is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("WASM runtime not enabled (build with -tags runtime_wasm)")
}

// Execute returns an error
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("WASM runtime not enabled")
}

// Call returns an error
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("WASM runtime not enabled")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "wasm"
}

// Version returns a stub version
func (r *Runtime) Version() string {
	return "stub (not enabled)"
}
