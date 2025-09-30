//go:build !runtime_rust
// +build !runtime_rust

package rust

import (
	"context"
	"fmt"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime is a stub when Rust is not enabled
type Runtime struct{}

// NewRuntime creates a stub runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Rust is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("rust runtime not enabled in build")
}

// Execute is not available
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("rust runtime not enabled in build")
}

// Call is not available
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("rust runtime not enabled in build")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "rust"
}

// Version returns a placeholder
func (r *Runtime) Version() string {
	return "disabled"
}
