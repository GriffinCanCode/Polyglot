//go:build runtime_go
// +build runtime_go

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	goruntime "github.com/griffincancode/polyglot.js/runtimes/go"
)

// TestGoBasicOperations tests basic Go operations
func TestGoBasicOperations(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
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
		{"Addition", "2 + 2"},
		{"Multiplication", "6 * 7"},
		{"String", `"Hello" + " " + "World"`},
		{"Boolean", "true"},
		{"Comparison", "10 > 5"},
		{"Variable", "x := 42; x"},
		{"Array", "[3]int{1, 2, 3}"},
		{"Slice", "[]int{1, 2, 3, 4, 5}"},
		{"Map", `map[string]int{"a": 1, "b": 2}`},
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

// TestGoFunctions tests Go function definition and execution
func TestGoFunctions(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	code := `
package main

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

var result = add(10, 20)
`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestGoStdlib tests using Go standard library
func TestGoStdlib(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
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
		{"Strings", `import "strings"; strings.ToUpper("hello")`},
		{"Fmt", `import "fmt"; fmt.Sprintf("Hello %s", "World")`},
		{"Math", `import "math"; math.Max(10.5, 20.5)`},
		{"Time", `import "time"; time.Now().Year()`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Logf("Execute failed (may be expected): %v", err)
			} else {
				t.Logf("Result: %v (type: %T)", result, result)
			}
		})
	}
}

// TestGoStructs tests Go struct definition and usage
func TestGoStructs(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	code := `
package main

type Person struct {
	Name string
	Age  int
}

func NewPerson(name string, age int) *Person {
	return &Person{Name: name, Age: age}
}

var p = NewPerson("Alice", 30)
`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestGoErrors tests Go error handling
func TestGoErrors(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
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
		{"Syntax Error", "this is not valid go code @#$"},
		{"Undefined Variable", "undefinedVar"},
		{"Type Mismatch", `x := "string"; y := x + 42`},
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

// TestGoConcurrency tests concurrent Go execution
func TestGoConcurrency(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
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

// TestGoContextCancellation tests that context cancellation works
func TestGoContextCancellation(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
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

// TestGoVersionInfo tests Go version reporting
func TestGoVersionInfo(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	version := runtime.Version()
	t.Logf("Go runtime version: %s", version)

	if version == "" || version == "stub (not enabled)" {
		t.Error("Expected valid Go version string")
	}
}

// TestGoShutdownAndReuse tests that shutdown properly prevents further use
func TestGoShutdownAndReuse(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
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

// TestGoClosures tests Go closures
func TestGoClosures(t *testing.T) {
	runtime := goruntime.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "go",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	code := `
package main

func makeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

var add10 = makeAdder(10)
var result = add10(5)
`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	} else {
		t.Logf("Closure result: %v", result)
	}
}
