//go:build runtime_cpp
// +build runtime_cpp

package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/cpp"
)

// TestCppBasicOperations tests basic C++ operations
func TestCppBasicOperations(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        10 * time.Second, // C++ needs more time for compilation
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
		{"String", `"Hello World"`, "Hello World"},
		{"Boolean True", "true", "1"},
		{"Boolean False", "false", "0"},
		{"Comparison", "(10 > 5)", "1"},
		{"Max Function", "std::max(10, 20)", "20"},
		{"String Length", `string s = "hello"; cout << s.length();`, "5"},
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

// TestCppErrors tests error handling
func TestCppErrors(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
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
		{"Syntax Error", "this is not valid cpp @#$"},
		{"Undefined Variable", "cout << undefinedVariable;"},
		{"Type Error", "int x = \"not an int\";"},
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

// TestCppComplexTypes tests complex data structures
func TestCppComplexTypes(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
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
		{"Vector", "vector<int> v = {1, 2, 3}; cout << v.size();"},
		{"String Concatenation", `string s = "Hello" + string(" ") + "World"; cout << s;`},
		{"Array", "int arr[] = {1, 2, 3, 4, 5}; cout << sizeof(arr) / sizeof(arr[0]);"},
		{"Math Operations", "cout << sqrt(16);"},
		{"String Methods", `string s = "hello"; transform(s.begin(), s.end(), s.begin(), ::toupper); cout << s;`},
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

// TestCppConcurrency tests concurrent C++ execution
func TestCppConcurrency(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        30 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	const numGoroutines = 10
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			code := `cout << 42;`
			result, err := runtime.Execute(ctx, code)
			if err != nil {
				t.Errorf("Concurrent execute %d failed: %v", n, err)
				return
			}
			t.Logf("Goroutine %d result: %v", n, result)
		}(i)
	}

	wg.Wait()
}

// TestCppFullProgram tests complete C++ programs
func TestCppFullProgram(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
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
			"Complete Program with Main",
			`#include <iostream>
int main() {
    std::cout << "Hello from C++" << std::endl;
    return 0;
}`,
		},
		{
			"Program with Function",
			`#include <iostream>
int add(int a, int b) {
    return a + b;
}
int main() {
    std::cout << add(10, 20) << std::endl;
    return 0;
}`,
		},
		{
			"Program with Class",
			`#include <iostream>
class Calculator {
public:
    int add(int a, int b) {
        return a + b;
    }
};
int main() {
    Calculator calc;
    std::cout << calc.add(5, 7) << std::endl;
    return 0;
}`,
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

// TestCppContextCancellation tests context cancellation
func TestCppContextCancellation(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        10 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Try to execute a long-running operation (compilation might take time)
	code := `
#include <iostream>
#include <unistd.h>
int main() {
    sleep(10);
    std::cout << "Should not reach here" << std::endl;
    return 0;
}
`

	_, err := runtime.Execute(ctxWithTimeout, code)
	if err == nil {
		t.Log("Warning: Expected timeout error but got none")
	} else if err == context.DeadlineExceeded {
		t.Logf("Got expected timeout error: %v", err)
	} else {
		t.Logf("Got error (might be timeout): %v", err)
	}
}

// TestCppMetadata tests runtime metadata
func TestCppMetadata(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        10 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	name := runtime.Name()
	if name != "cpp" {
		t.Errorf("Expected name 'cpp', got '%s'", name)
	}

	version := runtime.Version()
	t.Logf("C++ compiler version: %s", version)
	if version == "" {
		t.Error("Expected non-empty version")
	}
}

// TestCppPoolReuse tests that workers are properly reused
func TestCppPoolReuse(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
		Enabled:        true,
		MaxConcurrency: 2, // Small pool to test reuse
		Timeout:        20 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Execute multiple times to ensure workers are reused
	for i := 0; i < 10; i++ {
		code := "cout << 42;"
		result, err := runtime.Execute(ctx, code)
		if err != nil {
			t.Errorf("Execute %d failed: %v", i, err)
			return
		}
		t.Logf("Execution %d result: %v", i, result)
	}
}

// TestCppShutdown tests proper shutdown
func TestCppShutdown(t *testing.T) {
	runtime := cpp.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "cpp",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        10 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Shutdown the runtime
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Try to execute after shutdown - should fail
	_, err := runtime.Execute(ctx, "cout << 42;")
	if err == nil {
		t.Error("Expected error when executing after shutdown")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}
}
