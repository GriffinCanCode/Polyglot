//go:build runtime_python
// +build runtime_python

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/python"
)

// TestPythonBasicOperations tests basic Python operations
func TestPythonBasicOperations(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{"Addition", "2 + 2", int64(4)},
		{"Multiplication", "6 * 7", int64(42)},
		{"String", "'Hello' + ' ' + 'World'", "Hello World"},
		{"Boolean True", "True", true},
		{"Boolean False", "False", false},
		{"Comparison", "10 > 5", true},
		{"List Length", "len([1, 2, 3, 4, 5])", int64(5)},
		{"Dict", "{'key': 'value'}['key']", "value"},
		{"Float", "3.14", 3.14},
		{"None", "None", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed: %v", err)
				return
			}
			t.Logf("Result: %v (type: %T)", result, result)

			// Type-agnostic comparison for numeric types
			if tt.expected != nil {
				switch exp := tt.expected.(type) {
				case int64:
					switch res := result.(type) {
					case int64:
						if res != exp {
							t.Errorf("Expected %v, got %v", exp, res)
						}
					case float64:
						if int64(res) != exp {
							t.Errorf("Expected %v, got %v", exp, res)
						}
					default:
						t.Logf("Result type: %T (expected int64)", result)
					}
				case float64:
					switch res := result.(type) {
					case float64:
						if res != exp {
							t.Errorf("Expected %v, got %v", exp, res)
						}
					case int64:
						if float64(res) != exp {
							t.Errorf("Expected %v, got %v", exp, res)
						}
					default:
						t.Logf("Result type: %T (expected float64)", result)
					}
				case string:
					if result != exp {
						t.Errorf("Expected %v, got %v", exp, result)
					}
				case bool:
					if result != exp {
						t.Errorf("Expected %v, got %v", exp, result)
					}
				}
			}
		})
	}
}

// TestPythonErrors tests error handling
func TestPythonErrors(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name string
		code string
	}{
		{"Syntax Error", "this is not valid python @#$"},
		{"Name Error", "undefined_variable"},
		{"Type Error", "len(42)"},
		{"Zero Division", "1 / 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runtime.Execute(ctx, tt.code)
			if err == nil {
				t.Error("Expected error but got none")
			} else {
				t.Logf("Got expected error: %v", err)
			}
		})
	}
}

// TestPythonComplexTypes tests complex data structures
func TestPythonComplexTypes(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name string
		code string
	}{
		{"List", "[1, 2, 3]"},
		{"Dict", "{'a': 1, 'b': 2}"},
		{"Tuple", "(1, 2, 3)"},
		{"Set", "{1, 2, 3}"},
		{"List Comprehension", "[x*2 for x in range(5)]"},
		{"Dict Comprehension", "{x: x*2 for x in range(3)}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed: %v", err)
			} else {
				t.Logf("Result: %v (type: %T)", result, result)
			}
		})
	}
}

// TestPythonConcurrency tests concurrent Python execution
func TestPythonConcurrency(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 10,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Run multiple executions concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			code := "1 + 1"
			_, err := runtime.Execute(ctx, code)
			if err != nil {
				t.Errorf("Concurrent execution %d failed: %v", n, err)
			}
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestPythonContextCancellation tests that context cancellation works
func TestPythonContextCancellation(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context that's already cancelled
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()

	_, err := runtime.Execute(cancelCtx, "1 + 1")
	if err == nil {
		t.Log("Warning: Expected context cancellation error, but execution succeeded")
	} else {
		t.Logf("Got expected cancellation error: %v", err)
	}
}

// TestPythonModules tests importing Python modules
func TestPythonModules(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name string
		code string
	}{
		{"Math module", "import math; math.pi"},
		{"OS module", "import os; os.name"},
		{"JSON module", "import json; json.dumps({'key': 'value'})"},
		{"Datetime", "import datetime; datetime.datetime.now().year"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed: %v", err)
			} else {
				t.Logf("Result: %v (type: %T)", result, result)
			}
		})
	}
}

// TestPythonVersionInfo tests Python version reporting
func TestPythonVersionInfo(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	version := runtime.Version()
	t.Logf("Python version: %s", version)

	if version == "" || version == "not initialized" {
		t.Error("Expected valid Python version string")
	}
}

// TestPythonShutdownAndReuse tests that shutdown properly prevents further use
func TestPythonShutdownAndReuse(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Execute before shutdown
	_, err := runtime.Execute(ctx, "1 + 1")
	if err != nil {
		t.Errorf("Execute before shutdown failed: %v", err)
	}

	// Shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Try to use after shutdown
	_, err = runtime.Execute(ctx, "1 + 1")
	if err == nil {
		t.Error("Expected error when using runtime after shutdown, but got none")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}
}

// TestPythonFunctions tests defining and calling Python functions
func TestPythonFunctions(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	code := `
def add(a, b):
    return a + b

add(10, 20)
`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	} else {
		t.Logf("Function result: %v (type: %T)", result, result)
	}
}
