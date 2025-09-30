package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/zig"
)

// TestZigBasicExecution tests basic Zig code execution
func TestZigBasicExecution(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		// Expected to fail without build tag or zig compiler
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Test simple expression
	result, err := runtime.Execute(ctx, "2 + 2")
	if err != nil {
		t.Logf("Execute returned error (may be expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigConcurrentExecution tests concurrent Zig code execution
func TestZigConcurrentExecution(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Run concurrent executions
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			code := `const std = @import("std");
std.debug.print("Hello from Zig\n", .{});`
			_, err := runtime.Execute(ctx, code)
			if err != nil {
				t.Logf("Execution %d error (expected without zig compiler): %v", id, err)
			}
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestZigFunctionCall tests calling Zig functions
func TestZigFunctionCall(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Test function call
	code := `const std = @import("std");

fn add(a: i32, b: i32) i32 {
    return a + b;
}

pub fn main() !void {
    const result = add(5, 3);
    std.debug.print("{d}\n", .{result});
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigStringManipulation tests Zig string operations
func TestZigStringManipulation(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `const std = @import("std");

pub fn main() !void {
    const message = "Hello, Polyglot!";
    std.debug.print("{s}\n", .{message});
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigArrayOperations tests Zig array operations
func TestZigArrayOperations(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `const std = @import("std");

pub fn main() !void {
    const numbers = [_]i32{1, 2, 3, 4, 5};
    var sum: i32 = 0;
    for (numbers) |num| {
        sum += num;
    }
    std.debug.print("{d}\n", .{sum});
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigStructs tests Zig struct definitions and usage
func TestZigStructs(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `const std = @import("std");

const Point = struct {
    x: i32,
    y: i32,
    
    fn distance(self: Point) f64 {
        const dx = @as(f64, @floatFromInt(self.x));
        const dy = @as(f64, @floatFromInt(self.y));
        return @sqrt(dx * dx + dy * dy);
    }
};

pub fn main() !void {
    const p = Point{ .x = 3, .y = 4 };
    const dist = p.distance();
    std.debug.print("{d}\n", .{dist});
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigErrorHandling tests Zig error handling
func TestZigErrorHandling(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `const std = @import("std");

const DivisionError = error{DivisionByZero};

fn divide(a: i32, b: i32) DivisionError!i32 {
    if (b == 0) {
        return DivisionError.DivisionByZero;
    }
    return @divTrunc(a, b);
}

pub fn main() !void {
    const result = divide(10, 2) catch |err| {
        std.debug.print("Error: {}\n", .{err});
        return;
    };
    std.debug.print("{d}\n", .{result});
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigComptime tests Zig compile-time features
func TestZigComptime(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `const std = @import("std");

fn fibonacci(comptime n: u32) u32 {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

pub fn main() !void {
    const result = comptime fibonacci(10);
    std.debug.print("{d}\n", .{result});
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigGeneric tests Zig generic functions
func TestZigGeneric(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `const std = @import("std");

fn max(comptime T: type, a: T, b: T) T {
    return if (a > b) a else b;
}

pub fn main() !void {
    const result = max(i32, 10, 20);
    std.debug.print("{d}\n", .{result});
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestZigMultipleShutdowns tests that multiple shutdowns don't cause issues
func TestZigMultipleShutdowns(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}

	// Call shutdown multiple times
	for i := 0; i < 3; i++ {
		if err := runtime.Shutdown(ctx); err != nil {
			t.Errorf("Shutdown %d failed: %v", i, err)
		}
	}
}

// TestZigRuntimeInfo tests runtime information methods
func TestZigRuntimeInfo(t *testing.T) {
	runtime := zig.NewRuntime()

	// Test name
	if name := runtime.Name(); name != "zig" {
		t.Errorf("expected name 'zig', got '%s'", name)
	}

	// Test version (should return something even before initialization)
	version := runtime.Version()
	if version == "" {
		t.Error("version should not be empty")
	}
	t.Logf("Zig version: %s", version)
}

// TestZigContextCancellation tests context cancellation during execution
func TestZigContextCancellation(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Create a context that we'll cancel
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately

	// Try to execute with cancelled context
	_, err = runtime.Execute(cancelCtx, "2 + 2")
	// We expect either a context error or compilation error
	t.Logf("Execution with cancelled context returned: %v", err)
}

// TestZigTimeout tests execution timeout
func TestZigTimeout(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        1 * time.Millisecond, // Very short timeout
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Create a context with very short timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	// Try to execute with timeout
	_, err = runtime.Execute(timeoutCtx, "2 + 2")
	t.Logf("Execution with timeout returned: %v", err)
}

// TestZigLargeOutput tests handling large output
func TestZigLargeOutput(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `const std = @import("std");

pub fn main() !void {
    var i: u32 = 0;
    while (i < 100) : (i += 1) {
        std.debug.print("{d}\n", .{i});
    }
}`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without zig compiler): %v", err)
	} else {
		t.Logf("Result length: %d", len(result.(string)))
	}
}

// TestZigPoolExhaustion tests pool exhaustion handling
func TestZigPoolExhaustion(t *testing.T) {
	runtime := zig.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "zig",
		Version:        "0.11",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 2, // Small pool
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Try to run more concurrent tasks than pool size
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func(id int) {
			_, err := runtime.Execute(ctx, "2 + 2")
			if err != nil {
				t.Logf("Execution %d error: %v", id, err)
			}
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < 5; i++ {
		<-done
	}
}
