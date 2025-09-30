//go:build !runtime_java
// +build !runtime_java

package java

import (
	"context"
	"fmt"

	"github.com/polyglot-framework/polyglot/core"
)

// Runtime is a stub when Java is not enabled
type Runtime struct{}

// NewRuntime creates a stub runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Java is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("java runtime not enabled in build")
}

// Execute is not available
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("java runtime not enabled in build")
}

// Call is not available
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("java runtime not enabled in build")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "java"
}

// Version returns a placeholder
func (r *Runtime) Version() string {
	return "disabled"
}
