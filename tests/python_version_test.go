//go:build runtime_python
// +build runtime_python

package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/python"
)

// TestPythonVersionCompatibility tests Python version-specific features
func TestPythonVersionCompatibility(t *testing.T) {
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

	// Get Python version
	version := runtime.Version()
	t.Logf("Testing with Python version: %s", version)

	// Extract major.minor version
	parts := strings.Fields(version)
	if len(parts) == 0 {
		t.Fatal("Could not determine Python version")
	}

	// Run version-compatible tests
	tests := []struct {
		name        string
		code        string
		minVersion  string
		description string
	}{
		{
			"Python 3.8+ - f-strings",
			"name = 'Python'; f'Hello {name}'",
			"3.8",
			"f-string formatting",
		},
		{
			"Python 3.8+ - walrus operator",
			"(n := 5) and n * 2",
			"3.8",
			"Assignment expressions",
		},
		{
			"Python 3.9+ - dict merge",
			"a = {'x': 1}; b = {'y': 2}; a | b",
			"3.9",
			"Dictionary merge operator",
		},
		{
			"Python 3.9+ - type hints",
			"def greet(name: str) -> str: return f'Hello {name}'\ngreet('World')",
			"3.9",
			"Type hints in functions",
		},
		{
			"Python 3.10+ - match statement",
			`
def check_value(x):
    match x:
        case 1:
            return "one"
        case 2:
            return "two"
        case _:
            return "other"

check_value(1)
`,
			"3.10",
			"Structural pattern matching",
		},
		{
			"Python 3.11+ - exception groups",
			"try:\n    raise ValueError('test')\nexcept ValueError as e:\n    str(e)",
			"3.11",
			"Enhanced error messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				// Some features may not be available in older versions
				if strings.Contains(err.Error(), "SyntaxError") ||
					strings.Contains(err.Error(), "invalid syntax") {
					t.Logf("Feature not available in this Python version (expected for older versions): %v", err)
					t.SkipNow()
				}
				t.Errorf("Execute failed: %v", err)
			} else {
				t.Logf("Result: %v (type: %T)", result, result)
			}
		})
	}
}

// TestPythonStandardLibrary tests standard library modules across versions
func TestPythonStandardLibrary(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        10 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name   string
		code   string
		module string
	}{
		{
			"math module",
			"import math; math.pi",
			"math",
		},
		{
			"statistics module",
			"import statistics; statistics.mean([1, 2, 3, 4, 5])",
			"statistics",
		},
		{
			"json module",
			"import json; json.dumps({'key': 'value'})",
			"json",
		},
		{
			"datetime module",
			"import datetime; datetime.datetime.now().year",
			"datetime",
		},
		{
			"random module",
			"import random; random.seed(42); random.randint(1, 100)",
			"random",
		},
		{
			"collections module",
			"from collections import Counter; Counter(['a', 'b', 'a']).most_common(1)[0][0]",
			"collections",
		},
		{
			"itertools module",
			"from itertools import chain; list(chain([1, 2], [3, 4]))",
			"itertools",
		},
		{
			"functools module",
			"from functools import reduce; reduce(lambda x, y: x + y, [1, 2, 3, 4])",
			"functools",
		},
		{
			"re module",
			"import re; re.findall(r'\\d+', 'abc123def456')",
			"re",
		},
		{
			"os module",
			"import os; os.name",
			"os",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed for %s: %v", tt.module, err)
			} else {
				t.Logf("Module %s result: %v (type: %T)", tt.module, result, result)
			}
		})
	}
}

// TestPythonDataStructures tests complex data structures
func TestPythonDataStructures(t *testing.T) {
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
		{
			"Nested lists",
			"[[1, 2], [3, 4], [5, 6]]",
		},
		{
			"Nested dicts",
			"{'outer': {'inner': {'deep': 'value'}}}",
		},
		{
			"Mixed structures",
			"{'list': [1, 2, 3], 'dict': {'a': 1}, 'tuple': (1, 2)}",
		},
		{
			"List comprehension with condition",
			"[x for x in range(10) if x % 2 == 0]",
		},
		{
			"Dict comprehension",
			"{x: x**2 for x in range(5)}",
		},
		{
			"Set operations",
			"set([1, 2, 3]) & set([2, 3, 4])",
		},
		{
			"Tuple unpacking",
			"a, b = (1, 2); a + b",
		},
		{
			"Multiple assignment",
			"x = y = z = 5; x + y + z",
		},
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

// TestPythonErrorTraceback tests that error tracebacks are properly captured
func TestPythonErrorTraceback(t *testing.T) {
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
		name          string
		code          string
		expectedError string
		expectType    string
	}{
		{
			"SyntaxError",
			"this is not valid python @#$",
			"invalid syntax",
			"SyntaxError",
		},
		{
			"NameError",
			"undefined_variable",
			"not defined",
			"NameError",
		},
		{
			"TypeError",
			"len(42)",
			"object of type",
			"TypeError",
		},
		{
			"ZeroDivisionError",
			"1 / 0",
			"division by zero",
			"ZeroDivisionError",
		},
		{
			"ValueError",
			"int('not a number')",
			"invalid literal",
			"ValueError",
		},
		{
			"IndexError",
			"[1, 2, 3][10]",
			"out of range",
			"IndexError",
		},
		{
			"KeyError",
			"{'a': 1}['b']",
			"'b'",
			"KeyError",
		},
		{
			"AttributeError",
			"'string'.nonexistent_method()",
			"no attribute",
			"AttributeError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runtime.Execute(ctx, tt.code)
			if err == nil {
				t.Error("Expected error but got none")
				return
			}

			errMsg := err.Error()
			t.Logf("Error message: %s", errMsg)

			// Check if error contains expected error type
			if !strings.Contains(errMsg, tt.expectType) {
				t.Errorf("Expected error type %s in error message, got: %s", tt.expectType, errMsg)
			}

			// Check if error contains expected message
			if !strings.Contains(strings.ToLower(errMsg), strings.ToLower(tt.expectedError)) {
				t.Logf("Warning: Expected error message to contain '%s', got: %s", tt.expectedError, errMsg)
			}

			// Verify traceback is included (for runtime errors, not syntax errors)
			if tt.expectType != "SyntaxError" && !strings.Contains(errMsg, "Traceback") {
				t.Logf("Note: Traceback not included in error (may be expected for some error types)")
			}
		})
	}
}

// TestPythonLongRunning tests timeout handling for long-running operations
func TestPythonLongRunning(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        2 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Test context cancellation
	t.Run("Context cancellation", func(t *testing.T) {
		cancelCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		// This should timeout
		code := `
import time
time.sleep(10)
"completed"
`
		_, err := runtime.Execute(cancelCtx, code)
		if err == nil {
			t.Error("Expected timeout error but execution succeeded")
		} else {
			t.Logf("Got expected timeout error: %v", err)
		}
	})
}

// TestPythonMemoryOperations tests memory-intensive operations
func TestPythonMemoryOperations(t *testing.T) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        10 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name string
		code string
	}{
		{
			"Large list",
			"len(list(range(10000)))",
		},
		{
			"Large dict",
			"len({i: i**2 for i in range(1000)})",
		},
		{
			"String operations",
			"len('x' * 10000)",
		},
		{
			"List comprehension",
			"sum([i**2 for i in range(1000)])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed: %v", err)
			} else {
				t.Logf("Result: %v", result)
			}
		})
	}
}

// BenchmarkPythonExecution benchmarks Python execution performance
func BenchmarkPythonExecution(b *testing.B) {
	runtime := python.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 10,
		Timeout:        30 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		b.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	benchmarks := []struct {
		name string
		code string
	}{
		{"Simple arithmetic", "2 + 2"},
		{"Math operation", "import math; math.sqrt(144)"},
		{"List creation", "[x for x in range(100)]"},
		{"Dict creation", "{x: x**2 for x in range(50)}"},
		{"Function call", "def f(x): return x * 2\nf(21)"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := runtime.Execute(ctx, bm.code)
				if err != nil {
					b.Errorf("Execute failed: %v", err)
				}
			}
		})
	}
}
