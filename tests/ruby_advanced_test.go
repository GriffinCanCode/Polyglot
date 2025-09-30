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

// TestRubyCallWithArguments tests calling Ruby functions with arguments
func TestRubyCallWithArguments(t *testing.T) {
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

	// Define functions that take arguments
	_, err := runtime.Execute(ctx, `
		def add(a, b)
			a + b
		end

		def multiply(a, b)
			a * b
		end

		def greet_person(name)
			"Hello, #{name}!"
		end

		def calculate(x, y, z)
			x * y + z
		end
	`)
	if err != nil {
		t.Fatalf("Failed to define functions: %v", err)
	}

	tests := []struct {
		name     string
		fn       string
		args     []interface{}
		expected interface{}
	}{
		{
			"add integers",
			"add",
			[]interface{}{5, 3},
			int64(8),
		},
		{
			"multiply integers",
			"multiply",
			[]interface{}{7, 6},
			int64(42),
		},
		{
			"greet with string",
			"greet_person",
			[]interface{}{"Ruby"},
			"Hello, Ruby!",
		},
		{
			"three arguments",
			"calculate",
			[]interface{}{10, 5, 2},
			int64(52),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Call(ctx, tt.fn, tt.args...)
			if err != nil {
				t.Fatalf("Call failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

// TestRubyCallWithDifferentTypes tests calling with various argument types
func TestRubyCallWithDifferentTypes(t *testing.T) {
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

	// Define test functions
	_, err := runtime.Execute(ctx, `
		def test_float(x)
			x * 2.0
		end

		def test_boolean(flag)
			if flag
				"yes"
			else
				"no"
			end
		end

		def test_nil(value)
			value.nil?
		end

		def test_mixed(str, num, flag)
			if flag
				"#{str}: #{num}"
			else
				"none"
			end
		end
	`)
	if err != nil {
		t.Fatalf("Failed to define functions: %v", err)
	}

	tests := []struct {
		name     string
		fn       string
		args     []interface{}
		validate func(interface{}) bool
	}{
		{
			"float argument",
			"test_float",
			[]interface{}{3.5},
			func(result interface{}) bool {
				if f, ok := result.(float64); ok {
					return f == 7.0
				}
				return false
			},
		},
		{
			"boolean true",
			"test_boolean",
			[]interface{}{true},
			func(result interface{}) bool {
				return result == "yes"
			},
		},
		{
			"boolean false",
			"test_boolean",
			[]interface{}{false},
			func(result interface{}) bool {
				return result == "no"
			},
		},
		{
			"nil argument",
			"test_nil",
			[]interface{}{nil},
			func(result interface{}) bool {
				return result == true
			},
		},
		{
			"mixed types",
			"test_mixed",
			[]interface{}{"Answer", 42, true},
			func(result interface{}) bool {
				return result == "Answer: 42"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Call(ctx, tt.fn, tt.args...)
			if err != nil {
				t.Fatalf("Call failed: %v", err)
			}

			if !tt.validate(result) {
				t.Errorf("Validation failed for result: %v (%T)", result, result)
			}
		})
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

// TestRubyWorkerPoolReuse tests that workers are properly reused
func TestRubyWorkerPoolReuse(t *testing.T) {
	runtime := ruby.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "ruby",
		Enabled:        true,
		MaxConcurrency: 2, // Small pool to force reuse
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Define a variable in the first execution
	_, err := runtime.Execute(ctx, "$test_counter = 0")
	if err != nil {
		t.Fatalf("Failed to set variable: %v", err)
	}

	// Increment the counter multiple times
	for i := 0; i < 10; i++ {
		result, err := runtime.Execute(ctx, "$test_counter += 1")
		if err != nil {
			t.Fatalf("Iteration %d failed: %v", i, err)
		}
		t.Logf("Counter value: %v", result)
	}
}

// TestRubyClassesAndObjects tests Ruby class definition and object creation
func TestRubyClassesAndObjects(t *testing.T) {
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

	// Define a simple class
	_, err := runtime.Execute(ctx, `
		class Calculator
			def add(a, b)
				a + b
			end

			def multiply(a, b)
				a * b
			end
		end

		$calc = Calculator.new
	`)
	if err != nil {
		t.Fatalf("Failed to define class: %v", err)
	}

	// Use the class instance
	result, err := runtime.Execute(ctx, "$calc.add(10, 20)")
	if err != nil {
		t.Fatalf("Failed to call method: %v", err)
	}

	if result != int64(30) {
		t.Errorf("Expected 30, got %v", result)
	}

	result, err = runtime.Execute(ctx, "$calc.multiply(7, 8)")
	if err != nil {
		t.Fatalf("Failed to call method: %v", err)
	}

	if result != int64(56) {
		t.Errorf("Expected 56, got %v", result)
	}
}

// TestRubyBlocks tests Ruby blocks and iterators
func TestRubyBlocks(t *testing.T) {
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
		{
			"map with block",
			"[1, 2, 3, 4].map { |x| x * 2 }.sum",
			int64(20),
		},
		{
			"select with block",
			"[1, 2, 3, 4, 5].select { |x| x > 2 }.length",
			int64(3),
		},
		{
			"times iterator",
			"result = 0; 5.times { |i| result += i }; result",
			int64(10),
		},
		{
			"each with accumulator",
			"sum = 0; [10, 20, 30].each { |n| sum += n }; sum",
			int64(60),
		},
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

// TestRubyHashOperations tests Ruby hash/dictionary operations
func TestRubyHashOperations(t *testing.T) {
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

	// Test hash operations
	result, err := runtime.Execute(ctx, "{ a: 1, b: 2, c: 3 }.length")
	if err != nil {
		t.Fatalf("Failed to get hash length: %v", err)
	}
	t.Logf("Hash length: %v", result)

	// Test hash access
	result, err = runtime.Execute(ctx, "h = { name: 'Ruby', version: 3 }; h[:name]")
	if err != nil {
		t.Fatalf("Failed to access hash: %v", err)
	}
	t.Logf("Hash access result: %v", result)
}

// TestRubyRegularExpressions tests Ruby regex support
func TestRubyRegularExpressions(t *testing.T) {
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
		{
			"regex match",
			"'hello world'.match?(/world/)",
		},
		{
			"regex substitution",
			"'hello world'.gsub(/world/, 'Ruby')",
		},
		{
			"regex split",
			"'one,two,three'.split(/,/).length",
		},
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

// TestRubyExceptionHandling tests Ruby exception handling
func TestRubyExceptionHandling(t *testing.T) {
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

	// Test proper exception handling
	code := `
		begin
			result = 10 / 2
		rescue => e
			result = -1
		end
		result
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result != int64(5) {
		t.Errorf("Expected 5, got %v", result)
	}
}

// TestRubyStringEscaping tests proper string escaping in arguments
func TestRubyStringEscaping(t *testing.T) {
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

	// Define a function that returns the string
	_, err := runtime.Execute(ctx, `
		def echo(str)
			str
		end
	`)
	if err != nil {
		t.Fatalf("Failed to define function: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", "hello"},
		{"string with quotes", "say \"hello\"", "say \"hello\""},
		{"string with backslash", "path\\to\\file", "path\\to\\file"},
		{"string with newline chars", "line one\\nline two", "line one\\nline two"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Call(ctx, "echo", tt.input)
			if err != nil {
				t.Fatalf("Call failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
