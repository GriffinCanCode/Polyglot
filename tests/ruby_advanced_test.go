//go:build runtime_ruby
// +build runtime_ruby

package tests

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/ruby"
)

// TestRubyBasicExecution tests basic Ruby code execution
func TestRubyBasicExecution(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{"simple addition", "2 + 2", int64(4)},
		{"string concat", "'hello' + ' world'", "hello world"},
		{"multiplication", "5 * 3", int64(15)},
		{"float division", "10.0 / 4.0", 2.5},
		{"boolean true", "true", true},
		{"boolean false", "false", false},
		{"nil value", "nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

// TestRubyStringOperations tests Ruby string operations
func TestRubyStringOperations(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name     string
		code     string
		validate func(interface{}) bool
	}{
		{
			"string upcase",
			"'hello'.upcase",
			func(result interface{}) bool {
				return result == "HELLO"
			},
		},
		{
			"string downcase",
			"'WORLD'.downcase",
			func(result interface{}) bool {
				return result == "world"
			},
		},
		{
			"string length",
			"'polyglot'.length",
			func(result interface{}) bool {
				return result == int64(8)
			},
		},
		{
			"string reverse",
			"'ruby'.reverse",
			func(result interface{}) bool {
				return result == "ybur"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}

			if !tt.validate(result) {
				t.Errorf("Validation failed for result: %v", result)
			}
		})
	}
}

// TestRubyArrayOperations tests Ruby array operations
func TestRubyArrayOperations(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name string
		code string
	}{
		{"array creation", "[1, 2, 3, 4, 5]"},
		{"array first", "[1, 2, 3].first"},
		{"array last", "[1, 2, 3].last"},
		{"array length", "[1, 2, 3, 4].length"},
		{"array sum", "[1, 2, 3, 4, 5].sum"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}
			t.Logf("Result: %v (type: %T)", result, result)
		})
	}
}

// TestRubyMathOperations tests Ruby mathematical operations
func TestRubyMathOperations(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{"power", "2 ** 8", int64(256)},
		{"modulo", "10 % 3", int64(1)},
		{"absolute value", "-42.abs", int64(42)},
		{"max", "[5, 2, 8, 1].max", int64(8)},
		{"min", "[5, 2, 8, 1].min", int64(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestRubyFunctionDefinition tests defining and calling Ruby functions
func TestRubyFunctionDefinition(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Define a function
	_, err := runtime.Execute(ctx, `
		def greet
			"Hello from Ruby!"
		end
	`)
	if err != nil {
		t.Fatalf("Failed to define function: %v", err)
	}

	// Call the function
	result, err := runtime.Call(ctx, "greet")
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	expected := "Hello from Ruby!"
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestRubyComplexFunction tests more complex Ruby functions
func TestRubyComplexFunction(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Define a function that calculates factorial
	_, err := runtime.Execute(ctx, `
		def factorial(n)
			return 1 if n <= 1
			n * factorial(n - 1)
		end
	`)
	if err != nil {
		t.Fatalf("Failed to define function: %v", err)
	}

	// Call the function (without arguments for now)
	result, err := runtime.Execute(ctx, "factorial(5)")
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	expected := int64(120)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestRubyConcurrency tests concurrent Ruby execution
func TestRubyConcurrency(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 10,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	const numGoroutines = 20
	const numIterations = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numIterations)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				code := fmt.Sprintf("%d + %d", id, j)
				result, err := runtime.Execute(ctx, code)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d iteration %d: %w", id, j, err)
					return
				}

				expected := int64(id + j)
				if result != expected {
					errors <- fmt.Errorf("goroutine %d iteration %d: expected %v, got %v", id, j, expected, result)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	var errorList []error
	for err := range errors {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		t.Errorf("Encountered %d errors during concurrent execution:", len(errorList))
		for _, err := range errorList {
			t.Errorf("  - %v", err)
		}
	}
}

// TestRubyContextCancellation tests context cancellation
func TestRubyContextCancellation(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context that will be cancelled
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately

	// Try to execute with cancelled context
	_, err := runtime.Execute(cancelCtx, "2 + 2")
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected 'context canceled' error, got: %v", err)
	}
}

// TestRubyTimeout tests execution timeout
func TestRubyTimeout(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context with very short timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Ensure timeout has passed

	// Try to execute with timed-out context
	_, err := runtime.Execute(timeoutCtx, "2 + 2")
	if err == nil {
		t.Error("Expected error due to timeout, got nil")
	}

	if !strings.Contains(err.Error(), "deadline exceeded") && !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// TestRubyErrorHandling tests Ruby error handling
func TestRubyErrorHandling(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name string
		code string
	}{
		{"syntax error", "def broken"},
		{"undefined variable", "undefined_variable"},
		{"undefined method", "nonexistent_method()"},
		{"division by zero", "1 / 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runtime.Execute(ctx, tt.code)
			if err == nil {
				t.Error("Expected error, got nil")
			} else {
				t.Logf("Got expected error: %v", err)
			}
		})
	}
}

// TestRubyVersion tests Ruby version reporting
func TestRubyVersion(t *testing.T) {
	runtime := ruby.NewRuntime()

	version := runtime.Version()
	if version == "" {
		t.Error("Version should not be empty")
	}

	t.Logf("Ruby version: %s", version)

	// Version should contain a number
	hasNumber := false
	for _, char := range version {
		if char >= '0' && char <= '9' {
			hasNumber = true
			break
		}
	}

	if !hasNumber {
		t.Errorf("Version string should contain a number: %s", version)
	}
}

// TestRubyRuntimeName tests runtime name
func TestRubyRuntimeName(t *testing.T) {
	runtime := ruby.NewRuntime()

	name := runtime.Name()
	if name != "ruby" {
		t.Errorf("Expected name 'ruby', got '%s'", name)
	}
}

// TestRubyShutdownIdempotency tests that shutdown can be called multiple times
func TestRubyShutdownIdempotency(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Call shutdown multiple times
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("First shutdown failed: %v", err)
	}

	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Second shutdown failed: %v", err)
	}

	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Third shutdown failed: %v", err)
	}
}

// TestRubyExecuteAfterShutdown tests execution after shutdown
func TestRubyExecuteAfterShutdown(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Shutdown the runtime
	if err := runtime.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Try to execute after shutdown
	_, err := runtime.Execute(ctx, "2 + 2")
	if err == nil {
		t.Error("Expected error when executing after shutdown, got nil")
	}

	if !strings.Contains(err.Error(), "shutdown") {
		t.Errorf("Expected 'shutdown' in error message, got: %v", err)
	}
}
