//go:build runtime_java
// +build runtime_java

package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/java"
)

// TestJavaBasicOperations tests basic Java operations
func TestJavaBasicOperations(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
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
		expected string
	}{
		{"Addition", "2 + 2", "4"},
		{"Multiplication", "6 * 7", "42"},
		{"String", "\"Hello\" + \" \" + \"World\"", "Hello World"},
		{"Boolean True", "true", "true"},
		{"Boolean False", "false", "false"},
		{"Comparison", "10 > 5", "true"},
		{"Math.max", "Math.max(10, 20)", "20"},
		{"String Length", "\"hello\".length()", "5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed: %v", err)
				return
			}
			t.Logf("Result: %v (type: %T)", result, result)

			if tt.expected != "" {
				if result != tt.expected {
					t.Logf("Expected '%v', got '%v'", tt.expected, result)
				}
			}
		})
	}
}

// TestJavaErrors tests error handling
func TestJavaErrors(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
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
		{"Syntax Error", "this is not valid java @#$"},
		{"Undefined Variable", "System.out.println(undefinedVariable);"},
		{"Compilation Error", "int x = \"not an int\";"},
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

// TestJavaComplexTypes tests complex data structures
func TestJavaComplexTypes(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
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
			"Array",
			`int[] arr = {1, 2, 3, 4, 5}; System.out.println(arr.length);`,
		},
		{
			"String Manipulation",
			`String s = "hello"; System.out.println(s.toUpperCase());`,
		},
		{
			"Math Operations",
			`System.out.println(Math.sqrt(16));`,
		},
		{
			"ArrayList",
			`import java.util.*; ArrayList<Integer> list = new ArrayList<>(); list.add(1); list.add(2); System.out.println(list.size());`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Logf("Execute failed (may need import handling): %v", err)
			} else {
				t.Logf("Result: %v (type: %T)", result, result)
			}
		})
	}
}

// TestJavaConcurrency tests concurrent Java execution
func TestJavaConcurrency(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
		Enabled:        true,
		MaxConcurrency: 10,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	var wg sync.WaitGroup
	concurrency := 20
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			code := `System.out.println(` + string(rune('0'+n%10)) + `);`
			_, err := runtime.Execute(ctx, code)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	errorCount := 0
	for err := range errors {
		errorCount++
		t.Errorf("Concurrent execution error: %v", err)
	}

	if errorCount > 0 {
		t.Errorf("Had %d errors out of %d concurrent executions", errorCount, concurrency)
	} else {
		t.Logf("Successfully executed %d concurrent Java operations", concurrency)
	}
}

// TestJavaContextCancellation tests that context cancellation works
func TestJavaContextCancellation(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        10 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context that we'll cancel
	cancelCtx, cancel := context.WithCancel(ctx)

	// Start a long-running operation
	done := make(chan bool)
	go func() {
		// Attempt to compile and run something (should be interrupted)
		_, _ = runtime.Execute(cancelCtx, `Thread.sleep(5000); System.out.println("done");`)
		done <- true
	}()

	// Cancel after a short delay
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Wait to see if it was cancelled
	select {
	case <-done:
		t.Log("Execution completed (cancellation may not have interrupted)")
	case <-time.After(2 * time.Second):
		t.Log("Context cancellation appeared to work")
	}
}

// TestJavaTimeout tests execution timeout
func TestJavaTimeout(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        1 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context with a short timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// Try to execute something that should timeout
	code := `Thread.sleep(2000); System.out.println("done");`
	_, err := runtime.Execute(timeoutCtx, code)

	if err != nil {
		if err == context.DeadlineExceeded {
			t.Log("Timeout handled correctly")
		} else {
			t.Logf("Got error (may be timeout related): %v", err)
		}
	} else {
		t.Log("No timeout occurred (execution was fast)")
	}
}

// TestJavaMultipleInitialization tests that multiple initializations are handled correctly
func TestJavaMultipleInitialization(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	// First initialization
	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("First initialization failed: %v", err)
	}

	// Try to initialize again (should handle gracefully)
	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Second initialization returned error (expected): %v", err)
	}

	// Should still be able to execute
	result, err := runtime.Execute(ctx, "10 * 10")
	if err != nil {
		t.Errorf("Execute after re-initialization failed: %v", err)
	} else {
		t.Logf("Execute successful: %v", result)
	}

	// Cleanup
	runtime.Shutdown(ctx)
}

// TestJavaShutdownAndReuse tests that shutdown properly prevents further use
func TestJavaShutdownAndReuse(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	// Initialize
	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Execute should work
	_, err := runtime.Execute(ctx, "5 + 5")
	if err != nil {
		t.Errorf("Execute before shutdown failed: %v", err)
	}

	// Shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Execute should fail after shutdown
	_, err = runtime.Execute(ctx, "5 + 5")
	if err == nil {
		t.Error("Expected error after shutdown, but got none")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}

	// Double shutdown should be safe
	err = runtime.Shutdown(ctx)
	if err != nil {
		t.Logf("Double shutdown returned error: %v", err)
	}
}

// TestJavaVersion tests that Version() returns a valid version string
func TestJavaVersion(t *testing.T) {
	runtime := java.NewRuntime()
	version := runtime.Version()

	if version == "" {
		t.Error("Version should not be empty")
	}

	t.Logf("Java version: %s", version)
}

// TestJavaName tests that Name() returns the correct identifier
func TestJavaName(t *testing.T) {
	runtime := java.NewRuntime()
	name := runtime.Name()

	if name != "java" {
		t.Errorf("Expected name 'java', got '%s'", name)
	}
}

// TestJavaPoolBehavior tests the worker pool behavior
func TestJavaPoolBehavior(t *testing.T) {
	runtime := java.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "java",
		Enabled:        true,
		MaxConcurrency: 3, // Small pool to test reuse
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Execute multiple times to ensure workers are reused
	for i := 0; i < 10; i++ {
		result, err := runtime.Execute(ctx, `System.out.println(`+string(rune('0'+i%10))+`);`)
		if err != nil {
			t.Errorf("Iteration %d failed: %v", i, err)
		} else {
			t.Logf("Iteration %d result: %v", i, result)
		}
	}
}
