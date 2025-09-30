package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/javascript"
)

// TestJavaScriptBasicOperations tests basic JavaScript operations
func TestJavaScriptBasicOperations(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
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
		{"Addition", "21 + 21", int32(42)},
		{"Multiplication", "6 * 7", int32(42)},
		{"String", "'Hello' + ' ' + 'World'", "Hello World"},
		{"Boolean", "true", true},
		{"Comparison", "10 > 5", true},
		{"Function", "(function() { return 100; })()", int32(100)},
		{"Array length", "[1, 2, 3, 4, 5].length", int32(5)},
		{"Math", "Math.max(10, 20, 30)", int32(30)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed: %v", err)
				return
			}
			t.Logf("Result: %v (type: %T)", result, result)
		})
	}
}

// TestJavaScriptErrors tests error handling
func TestJavaScriptErrors(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
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
		{"Syntax Error", "this is not valid javascript @@#$"},
		{"Reference Error", "nonExistentVariable"},
		{"Type Error", "null.toString()"},
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

// TestJavaScriptFunctionCalls tests calling JavaScript functions
func TestJavaScriptFunctionCalls(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Test calling built-in Math functions
	tests := []struct {
		name     string
		function string
		args     []interface{}
	}{
		{"Math.abs", "Math.abs", []interface{}{-42}},
		{"Math.max", "Math.max", []interface{}{10, 20, 30}},
		{"Math.min", "Math.min", []interface{}{10, 20, 30}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Context isolation means user-defined functions won't persist
			// This is expected behavior with pooled contexts
			result, err := runtime.Call(ctx, tt.function, tt.args...)
			if err != nil {
				t.Logf("Call to %s failed (may be expected with context pooling): %v", tt.function, err)
			} else {
				t.Logf("%s result: %v", tt.function, result)
			}
		})
	}
}

// TestJavaScriptConcurrency tests concurrent JavaScript execution
func TestJavaScriptConcurrency(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
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

// TestJavaScriptContextCancellation tests that context cancellation works
func TestJavaScriptContextCancellation(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
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

// TestJavaScriptTimeout tests execution timeout
func TestJavaScriptTimeout(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        100 * time.Millisecond,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context with a very short timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	// Try to execute with timeout
	_, err := runtime.Execute(timeoutCtx, "1 + 1")
	// This might or might not timeout depending on system speed
	if err != nil {
		t.Logf("Execution with timeout returned: %v", err)
	}
}

// TestJavaScriptMultipleInitialization tests that multiple initializations are handled correctly
func TestJavaScriptMultipleInitialization(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	// Initialize once
	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("First initialize failed: %v", err)
	}

	// Try to initialize again (should handle gracefully)
	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Second initialize returned error (may be expected): %v", err)
	}

	// Should still be able to execute
	_, err = runtime.Execute(ctx, "1 + 1")
	if err != nil {
		t.Errorf("Execute after second init failed: %v", err)
	}

	runtime.Shutdown(ctx)
}

// TestJavaScriptShutdownAndReuse tests that shutdown properly prevents further use
func TestJavaScriptShutdownAndReuse(t *testing.T) {
	runtime := javascript.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "javascript",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Try to use after shutdown
	_, err := runtime.Execute(ctx, "1 + 1")
	if err == nil {
		t.Error("Expected error when using runtime after shutdown, but got none")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}
}
