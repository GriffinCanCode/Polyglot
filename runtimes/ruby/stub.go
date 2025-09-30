//go:build !runtime_ruby
// +build !runtime_ruby

package ruby

import (
	"context"
	"fmt"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements a stub Ruby runtime when Ruby is not enabled
type Runtime struct {
	config core.RuntimeConfig
}

// NewRuntime creates a stub Ruby runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Ruby is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("Ruby runtime not enabled (build with -tags runtime_ruby)")
}

// Execute returns an error
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Ruby runtime not enabled")
}

// Call returns an error
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Ruby runtime not enabled")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "ruby"
}

// Version returns a stub version
func (r *Runtime) Version() string {
	return "stub (not enabled)"
}
