//go:build !runtime_lua
// +build !runtime_lua

package lua

import (
	"context"
	"fmt"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements a stub Lua runtime when Lua is not enabled
type Runtime struct {
	config core.RuntimeConfig
}

// NewRuntime creates a stub Lua runtime
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize returns an error indicating Lua is not enabled
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return fmt.Errorf("Lua runtime not enabled (build with -tags runtime_lua)")
}

// Execute returns an error
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Lua runtime not enabled")
}

// Call returns an error
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Lua runtime not enabled")
}

// Shutdown does nothing
func (r *Runtime) Shutdown(ctx context.Context) error {
	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "lua"
}

// Version returns a stub version
func (r *Runtime) Version() string {
	return "stub (not enabled)"
}
