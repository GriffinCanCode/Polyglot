//go:build !runtime_zig
// +build !runtime_zig

package zig

import (
	"context"
	"fmt"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime is a stub when Zig is not enabled
type Runtime struct{}

// NewRuntime creates a stub runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Zig is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("Zig runtime not enabled (build with -tags runtime_zig)")
}

// Execute is not available
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Zig runtime not enabled")
}

// Call is not available
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Zig runtime not enabled")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "zig"
}

// Version returns a placeholder
func (r *Runtime) Version() string {
	return "disabled"
}
