//go:build !runtime_python
// +build !runtime_python

package python

import (
	"context"
	"fmt"

	"github.com/polyglot-framework/polyglot/core"
)

// Runtime implements a stub Python runtime when Python is not enabled
type Runtime struct {
	config core.RuntimeConfig
}

// NewRuntime creates a stub Python runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Python is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("Python runtime not enabled (build with -tags runtime_python)")
}

// Execute returns an error
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Python runtime not enabled")
}

// Call returns an error
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Python runtime not enabled")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "python"
}

// Version returns a stub version
func (r *Runtime) Version() string {
	return "stub (not enabled)"
}
