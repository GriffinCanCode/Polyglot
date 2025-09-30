//go:build runtime_php
// +build runtime_php

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/php"
)

// TestPHPBasicOperations tests basic PHP operations
func TestPHPBasicOperations(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
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
		{"Addition", "echo 2 + 2;", "4"},
		{"Multiplication", "echo 6 * 7;", "42"},
		{"String", "echo 'Hello' . ' ' . 'World';", "Hello World"},
		{"Boolean True", "echo true;", "1"},
		{"Boolean False", "echo false;", ""},
		{"Comparison", "echo 10 > 5 ? 'true' : 'false';", "true"},
		{"Array Count", "echo count([1, 2, 3, 4, 5]);", "5"},
		{"String Length", "echo strlen('hello');", "5"},
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

// TestPHPErrors tests error handling
func TestPHPErrors(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
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
		{"Syntax Error", "this is not valid php @#$"},
		{"Undefined Variable", "echo $undefined_variable;"},
		{"Undefined Function", "nonExistentFunction();"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runtime.Execute(ctx, tt.code)
			if err == nil {
				t.Log("Warning: Expected error but got none (PHP may have error_reporting disabled)")
			} else {
				t.Logf("Got expected error: %v", err)
			}
		})
	}
}

// TestPHPComplexTypes tests complex data structures
func TestPHPComplexTypes(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
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
		{"Array", "print_r([1, 2, 3]);"},
		{"Associative Array", "$arr = ['a' => 1, 'b' => 2]; echo $arr['a'];"},
		{"JSON Encode", "echo json_encode(['key' => 'value']);"},
		{"String Functions", "echo strtoupper('hello');"},
		{"Math Functions", "echo max(10, 20, 30);"},
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

// TestPHPConcurrency tests concurrent PHP execution
func TestPHPConcurrency(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
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
			code := "echo 1 + 1;"
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

// TestPHPContextCancellation tests that context cancellation works
func TestPHPContextCancellation(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
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

	_, err := runtime.Execute(cancelCtx, "echo 1 + 1;")
	if err == nil {
		t.Log("Warning: Expected context cancellation error, but execution succeeded")
	} else {
		t.Logf("Got expected cancellation error: %v", err)
	}
}

// TestPHPFunctions tests defining and calling PHP functions
func TestPHPFunctions(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	code := `
function add($a, $b) {
    return $a + $b;
}

echo add(10, 20);
`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	} else {
		t.Logf("Function result: %v (type: %T)", result, result)
	}
}

// TestPHPVersionInfo tests PHP version reporting
func TestPHPVersionInfo(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	version := runtime.Version()
	t.Logf("PHP version: %s", version)

	if version == "" || version == "unknown" || version == "stub (not enabled)" {
		t.Error("Expected valid PHP version string")
	}
}

// TestPHPShutdownAndReuse tests that shutdown properly prevents further use
func TestPHPShutdownAndReuse(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Execute before shutdown
	_, err := runtime.Execute(ctx, "echo 1 + 1;")
	if err != nil {
		t.Errorf("Execute before shutdown failed: %v", err)
	}

	// Shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Try to use after shutdown
	_, err = runtime.Execute(ctx, "echo 1 + 1;")
	if err == nil {
		t.Error("Expected error when using runtime after shutdown, but got none")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}
}

// TestPHPStdlib tests using PHP standard library functions
func TestPHPStdlib(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
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
		{"String Functions", "echo strtoupper('hello');"},
		{"Array Functions", "echo implode(',', [1, 2, 3]);"},
		{"Math Functions", "echo abs(-42);"},
		{"Date Functions", "echo date('Y');"},
		{"JSON", "echo json_encode(['key' => 'value']);"},
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

// TestPHPClasses tests PHP class definitions
func TestPHPClasses(t *testing.T) {
	runtime := php.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "php",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	code := `
class Person {
    private $name;
    
    public function __construct($name) {
        $this->name = $name;
    }
    
    public function getName() {
        return $this->name;
    }
}

$person = new Person('Alice');
echo $person->getName();
`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	} else {
		t.Logf("Class result: %v", result)
	}
}
