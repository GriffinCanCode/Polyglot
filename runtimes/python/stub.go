//go:build !runtime_python
// +build !runtime_python

package python

import (
	"context"
	"errors"

	"github.com/griffincancode/polyglot.js/core"
)

var (
	errNotEnabled = errors.New("python runtime not enabled (build with -tags runtime_python)")
)

// Runtime implements a stub Python runtime when Python is not enabled
type Runtime struct{}

// NewRuntime creates a stub Python runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Python is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	_ = config // Explicitly mark as used
	return errNotEnabled
}

// Execute returns an error
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, errNotEnabled
}

// Call returns an error
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, errNotEnabled
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
