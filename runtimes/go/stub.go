//go:build !runtime_go
// +build !runtime_go

package goruntime

import (
	"context"
	"fmt"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements a stub Go runtime when Go interpreter is not enabled
type Runtime struct{}

// NewRuntime creates a stub Go runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Go runtime is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("Go runtime not enabled (build with -tags runtime_go)")
}

// Execute returns an error
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Go runtime not enabled")
}

// Call returns an error
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Go runtime not enabled")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "go"
}

// Version returns a stub version
func (r *Runtime) Version() string {
	return "stub (not enabled)"
}
