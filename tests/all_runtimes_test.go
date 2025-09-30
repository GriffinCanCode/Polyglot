package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/cpp"
	goruntime "github.com/griffincancode/polyglot.js/runtimes/go"
	"github.com/griffincancode/polyglot.js/runtimes/java"
	"github.com/griffincancode/polyglot.js/runtimes/javascript"
	"github.com/griffincancode/polyglot.js/runtimes/lua"
	"github.com/griffincancode/polyglot.js/runtimes/php"
	"github.com/griffincancode/polyglot.js/runtimes/python"
	"github.com/griffincancode/polyglot.js/runtimes/ruby"
	"github.com/griffincancode/polyglot.js/runtimes/rust"
	"github.com/griffincancode/polyglot.js/runtimes/wasm"
	"github.com/griffincancode/polyglot.js/runtimes/zig"
)

// TestAllRuntimes tests all supported runtimes
func TestAllRuntimes(t *testing.T) {
	ctx := context.Background()

	runtimes := []struct {
		name     string
		runtime  core.Runtime
		hasImpl  bool // Whether the runtime has a full implementation (not just stub)
		testCode string
	}{
		{"python", python.NewRuntime(), true, "2 + 2"},
		{"javascript", javascript.NewRuntime(), true, "2 + 2"},
		{"go", goruntime.NewRuntime(), true, "2 + 2"},
		{"php", php.NewRuntime(), true, "echo 2 + 2;"},
		{"java", java.NewRuntime(), true, "2 + 2"},
		{"cpp", cpp.NewRuntime(), true, "2 + 2"},
		{"ruby", ruby.NewRuntime(), true, "2 + 2"},
		{"lua", lua.NewRuntime(), true, "return 2 + 2"},
		{"zig", zig.NewRuntime(), false, ""},
		{"wasm", wasm.NewRuntime(), false, ""},
		{"rust", rust.NewRuntime(), true, "println!(\"4\")"},
	}

	for _, rt := range runtimes {
		rt := rt // capture loop variable
		t.Run(rt.name, func(t *testing.T) {
			// Verify name
			if rt.runtime.Name() != rt.name {
				t.Errorf("expected name '%s', got '%s'", rt.name, rt.runtime.Name())
			}

			// Verify version returns something
			version := rt.runtime.Version()
			if version == "" {
				t.Error("version should not be empty")
			}
			t.Logf("%s version: %s", rt.name, version)

			config := core.RuntimeConfig{
				Name:           rt.name,
				Enabled:        true,
				MaxConcurrency: 5,
				Timeout:        5 * time.Second,
			}

			// Test initialization
			err := rt.runtime.Initialize(ctx, config)
			if err != nil {
				// Check if it's the expected "not enabled" error
				errMsg := strings.ToLower(err.Error())
				if strings.Contains(errMsg, "not enabled") || strings.Contains(errMsg, "build with") {
					t.Logf("Runtime %s not enabled (expected for stub): %v", rt.name, err)
					// Skip execution tests for disabled runtimes
					_ = rt.runtime.Shutdown(ctx)
					return
				} else if rt.hasImpl {
					t.Fatalf("Runtime %s with full implementation returned unexpected error: %v", rt.name, err)
				} else {
					t.Logf("Runtime %s (stub) returned error: %v", rt.name, err)
				}
			} else if rt.hasImpl {
				t.Logf("Runtime %s initialized successfully", rt.name)
			}

			// Test execution if we have test code and runtime is initialized
			if rt.testCode != "" && rt.hasImpl {
				result, execErr := rt.runtime.Execute(ctx, rt.testCode)
				if execErr != nil {
					t.Errorf("Execute failed: %v", execErr)
				} else {
					t.Logf("Execute result: %v (type: %T)", result, result)
				}
			}

			// Test shutdown
			shutdownErr := rt.runtime.Shutdown(ctx)
			if shutdownErr != nil {
				t.Errorf("shutdown failed: %v", shutdownErr)
			}
		})
	}
}

// TestAllRuntimesConcurrency tests all runtimes can be created and shut down concurrently
func TestAllRuntimesConcurrency(t *testing.T) {
	runtimes := []core.Runtime{
		python.NewRuntime(),
		javascript.NewRuntime(),
		goruntime.NewRuntime(),
		php.NewRuntime(),
		ruby.NewRuntime(),
		lua.NewRuntime(),
		zig.NewRuntime(),
		wasm.NewRuntime(),
		rust.NewRuntime(),
		java.NewRuntime(),
		cpp.NewRuntime(),
	}

	// Create all runtimes concurrently
	for _, runtime := range runtimes {
		runtime := runtime // capture loop variable
		t.Run(runtime.Name()+"_concurrent", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			config := core.RuntimeConfig{
				Name:           runtime.Name(),
				Enabled:        true,
				MaxConcurrency: 5,
				Timeout:        5 * time.Second,
			}

			// Initialize
			err := runtime.Initialize(ctx, config)
			if err != nil {
				errMsg := strings.ToLower(err.Error())
				if !strings.Contains(errMsg, "not enabled") && !strings.Contains(errMsg, "build with") {
					t.Logf("Runtime %s returned error: %v", runtime.Name(), err)
				}
			}

			// Shutdown
			_ = runtime.Shutdown(ctx)
		})
	}
}

// TestRuntimeInterface ensures all runtimes implement the core.Runtime interface correctly
func TestRuntimeInterface(t *testing.T) {
	var _ core.Runtime = python.NewRuntime()
	var _ core.Runtime = javascript.NewRuntime()
	var _ core.Runtime = goruntime.NewRuntime()
	var _ core.Runtime = php.NewRuntime()
	var _ core.Runtime = ruby.NewRuntime()
	var _ core.Runtime = lua.NewRuntime()
	var _ core.Runtime = zig.NewRuntime()
	var _ core.Runtime = wasm.NewRuntime()
	var _ core.Runtime = rust.NewRuntime()
	var _ core.Runtime = java.NewRuntime()
	var _ core.Runtime = cpp.NewRuntime()
}

// TestRuntimeNames ensures all runtime names are unique and correct
func TestRuntimeNames(t *testing.T) {
	expectedNames := map[string]core.Runtime{
		"python":     python.NewRuntime(),
		"javascript": javascript.NewRuntime(),
		"go":         goruntime.NewRuntime(),
		"php":        php.NewRuntime(),
		"ruby":       ruby.NewRuntime(),
		"lua":        lua.NewRuntime(),
		"zig":        zig.NewRuntime(),
		"wasm":       wasm.NewRuntime(),
		"rust":       rust.NewRuntime(),
		"java":       java.NewRuntime(),
		"cpp":        cpp.NewRuntime(),
	}

	// Verify each runtime has the correct name
	for expectedName, runtime := range expectedNames {
		actualName := runtime.Name()
		if actualName != expectedName {
			t.Errorf("expected runtime name '%s', got '%s'", expectedName, actualName)
		}
	}

	// Verify all names are unique
	names := make(map[string]bool)
	for _, runtime := range expectedNames {
		name := runtime.Name()
		if names[name] {
			t.Errorf("duplicate runtime name: %s", name)
		}
		names[name] = true
	}

	// Verify we have all expected runtimes
	if len(names) != 11 {
		t.Errorf("expected 11 unique runtime names, got %d", len(names))
	}
}
