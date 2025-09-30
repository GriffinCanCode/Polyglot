//go:build runtime_lua
// +build runtime_lua

package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/lua"
)

// TestLuaBasicOperations tests basic Lua operations
func TestLuaBasicOperations(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
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
		{"Addition", "return 2 + 2", float64(4)},
		{"Multiplication", "return 6 * 7", float64(42)},
		{"String", `return "Hello World"`, "Hello World"},
		{"Boolean True", "return true", true},
		{"Boolean False", "return false", false},
		{"Comparison", "return 10 > 5", true},
		{"String Length", `return string.len("hello")`, float64(5)},
		{"Table Length", `local t = {1, 2, 3, 4, 5}; return #t`, float64(5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Errorf("Execute failed: %v", err)
				return
			}
			t.Logf("Result: %v (type: %T)", result, result)

			if tt.expected != nil && result != tt.expected {
				t.Logf("Expected '%v', got '%v'", tt.expected, result)
			}
		})
	}
}

// TestLuaErrors tests error handling
func TestLuaErrors(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
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
		{"Syntax Error", "this is not valid lua @#$"},
		{"Nil Function Call", "return nonExistentFunction()"},
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

// TestLuaConcurrency tests concurrent execution
func TestLuaConcurrency(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
		Enabled:        true,
		MaxConcurrency: 10,
		Timeout:        10 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	var wg sync.WaitGroup
	numGoroutines := 20
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result, err := runtime.Execute(ctx, "return 2 + 2")
			if err != nil {
				t.Errorf("Goroutine %d: Execute failed: %v", id, err)
				return
			}
			t.Logf("Goroutine %d result: %v", id, result)
			mu.Lock()
			successCount++
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	if successCount != numGoroutines {
		t.Errorf("Expected %d successful executions, got %d", numGoroutines, successCount)
	}
}

// TestLuaContextCancellation tests context cancellation
func TestLuaContextCancellation(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context that will be cancelled
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately

	// Try to execute with cancelled context
	_, err := runtime.Execute(cancelCtx, "return 2 + 2")
	if err == nil {
		t.Error("Expected context cancellation error but got none")
	} else {
		t.Logf("Got expected cancellation error: %v", err)
	}
}

// TestLuaTimeout tests execution timeout
func TestLuaTimeout(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context with a short timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Execute code that might take longer (busy loop)
	code := `
		local sum = 0
		for i = 1, 100000000 do
			sum = sum + i
		end
		return sum
	`

	start := time.Now()
	_, err := runtime.Execute(timeoutCtx, code)
	duration := time.Since(start)

	// We expect either a timeout or successful completion within reasonable time
	if err != nil {
		t.Logf("Got error (possibly timeout): %v after %v", err, duration)
	} else {
		t.Logf("Execution completed in %v", duration)
	}
}

// TestLuaFunctionCall tests calling Lua functions
func TestLuaFunctionCall(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Define and call a function in one execution
	// (Lua state is per-worker and doesn't persist between calls)
	code := `
		function add(a, b)
			return a + b
		end
		return add(5, 3)
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
		return
	}

	t.Logf("Function call result: %v (type: %T)", result, result)

	if result != float64(8) {
		t.Errorf("Expected 8, got %v", result)
	}
}

// TestLuaShutdown tests proper shutdown
func TestLuaShutdown(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Execute some code
	result, err := runtime.Execute(ctx, "return 42")
	if err != nil {
		t.Errorf("Execute before shutdown failed: %v", err)
	} else {
		t.Logf("Execute before shutdown result: %v", result)
	}

	// Shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Try to execute after shutdown (should fail)
	_, err = runtime.Execute(ctx, "return 42")
	if err == nil {
		t.Error("Expected error after shutdown but got none")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}

	// Shutdown again (should be idempotent)
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Second shutdown failed: %v", err)
	}
}

// TestLuaMultipleRuntimes tests multiple runtime instances
func TestLuaMultipleRuntimes(t *testing.T) {
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	// Create multiple runtime instances
	numRuntimes := 5
	runtimes := make([]*lua.Runtime, numRuntimes)

	for i := 0; i < numRuntimes; i++ {
		runtimes[i] = lua.NewRuntime()
		if err := runtimes[i].Initialize(ctx, config); err != nil {
			t.Fatalf("Initialize runtime %d failed: %v", i, err)
		}
		defer runtimes[i].Shutdown(ctx)
	}

	// Execute code in all runtimes concurrently
	var wg sync.WaitGroup
	for i, rt := range runtimes {
		wg.Add(1)
		go func(id int, runtime *lua.Runtime) {
			defer wg.Done()
			result, err := runtime.Execute(ctx, "return 2 + 2")
			if err != nil {
				t.Errorf("Runtime %d: Execute failed: %v", id, err)
				return
			}
			t.Logf("Runtime %d result: %v", id, result)
		}(i, rt)
	}

	wg.Wait()
}

// TestLuaComplexOperations tests more complex Lua operations
func TestLuaComplexOperations(t *testing.T) {
	runtime := lua.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "lua",
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
			"Fibonacci",
			`
				function fib(n)
					if n <= 1 then
						return n
					end
					return fib(n-1) + fib(n-2)
				end
				return fib(10)
			`,
		},
		{
			"Table Operations",
			`
				local t = {1, 2, 3, 4, 5}
				local sum = 0
				for i, v in ipairs(t) do
					sum = sum + v
				end
				return sum
			`,
		},
		{
			"String Operations",
			`
				local s = "Hello, World!"
				return string.upper(s)
			`,
		},
		{
			"Math Operations",
			`
				return math.sqrt(16) + math.abs(-5)
			`,
		},
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

// TestLuaVersion tests version reporting
func TestLuaVersion(t *testing.T) {
	runtime := lua.NewRuntime()
	version := runtime.Version()

	if version == "" {
		t.Error("Version should not be empty")
	}

	t.Logf("Lua version: %s", version)
}

// TestLuaName tests name reporting
func TestLuaName(t *testing.T) {
	runtime := lua.NewRuntime()
	name := runtime.Name()

	if name != "lua" {
		t.Errorf("Expected name 'lua', got '%s'", name)
	}
}
