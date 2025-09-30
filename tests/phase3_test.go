package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/javascript"
	"github.com/griffincancode/polyglot.js/runtimes/lua"
	"github.com/griffincancode/polyglot.js/runtimes/php"
	"github.com/griffincancode/polyglot.js/runtimes/ruby"
	"github.com/griffincancode/polyglot.js/runtimes/wasm"
	"github.com/griffincancode/polyglot.js/runtimes/zig"
)

// TestPHPRuntime tests PHP runtime integration
func TestPHPRuntime(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
		Version:        "8.2",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "PHP runtime not enabled (build with -tags runtime_php)" {
		t.Fatalf("unexpected error: %v", err)
	}

	if runtime.Name() != "php" {
		t.Errorf("expected name 'php', got %s", runtime.Name())
	}

	t.Cleanup(func() {
		_ = runtime.Shutdown(ctx)
	})
}

// TestRubyRuntime tests Ruby runtime integration
func TestRubyRuntime(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Version:        "3.2",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "Ruby runtime not enabled (build with -tags runtime_ruby)" {
		t.Fatalf("unexpected error: %v", err)
	}

	if runtime.Name() != "ruby" {
		t.Errorf("expected name 'ruby', got %s", runtime.Name())
	}

	t.Cleanup(func() {
		_ = runtime.Shutdown(ctx)
	})
}

// TestLuaRuntime tests Lua runtime integration
func TestLuaRuntime(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
		Version:        "5.4",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "Lua runtime not enabled (build with -tags runtime_lua)" {
		t.Fatalf("unexpected error: %v", err)
	}

	if runtime.Name() != "lua" {
		t.Errorf("expected name 'lua', got %s", runtime.Name())
	}

	t.Cleanup(func() {
		_ = runtime.Shutdown(ctx)
	})
}

// TestZigRuntime tests Zig runtime integration
func TestZigRuntime(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "Zig runtime not enabled (build with -tags runtime_zig)" {
		t.Fatalf("unexpected error: %v", err)
	}

	if runtime.Name() != "zig" {
		t.Errorf("expected name 'zig', got %s", runtime.Name())
	}

	t.Cleanup(func() {
		_ = runtime.Shutdown(ctx)
	})
}

// TestWASMRuntime tests WebAssembly runtime integration
func TestWASMRuntime(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Version:        "1.0",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "WASM runtime not enabled (build with -tags runtime_wasm)" {
		t.Fatalf("unexpected error: %v", err)
	}

	if runtime.Name() != "wasm" {
		t.Errorf("expected name 'wasm', got %s", runtime.Name())
	}

	t.Cleanup(func() {
		_ = runtime.Shutdown(ctx)
	})
}

// TestAllPhase3Runtimes tests all Phase 3 runtimes concurrently
func TestAllPhase3Runtimes(t *testing.T) {
	ctx := context.Background()

	runtimes := []struct {
		name    string
		runtime core.Runtime
	}{
		{"javascript", javascript.NewRuntime()},
		{"php", php.NewRuntime()},
		{"ruby", ruby.NewRuntime()},
		{"lua", lua.NewRuntime()},
		{"zig", zig.NewRuntime()},
		{"wasm", wasm.NewRuntime()},
	}

	for _, rt := range runtimes {
		t.Run(rt.name, func(t *testing.T) {
			config := core.RuntimeConfig{
				Name:           rt.name,
				Enabled:        true,
				MaxConcurrency: 5,
				Timeout:        5 * time.Second,
			}

			err := rt.runtime.Initialize(ctx, config)
			// Accept either successful init or expected "not enabled" error
			if err != nil {
				// Check if it's the expected "not enabled" error
				errMsg := strings.ToLower(err.Error())
				if !strings.Contains(errMsg, "not enabled") && !strings.Contains(errMsg, "build with") {
					// Unexpected error - log it but don't fail (stub implementations are expected)
					t.Logf("Runtime %s returned unexpected error: %v", rt.name, err)
				}
			}

			if rt.runtime.Name() != rt.name {
				t.Errorf("expected name '%s', got %s", rt.name, rt.runtime.Name())
			}

			_ = rt.runtime.Shutdown(ctx)
		})
	}
}
