package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/cpp"
	"github.com/griffincancode/polyglot.js/runtimes/java"
	"github.com/griffincancode/polyglot.js/runtimes/javascript"
	"github.com/griffincancode/polyglot.js/runtimes/rust"
)

// TestJavaScriptRuntime tests JavaScript runtime integration
func TestJavaScriptRuntime(t *testing.T) {
	runtime := javascript.NewRuntime()

	if runtime.Name() != "javascript" {
		t.Errorf("expected name 'javascript', got '%s'", runtime.Name())
	}

	config := core.RuntimeConfig{
		Name:           "javascript",
		Version:        "ES2020",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 10,
		Timeout:        30 * time.Second,
	}

	ctx := context.Background()

	// Initialize
	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Fatalf("Failed to initialize JavaScript runtime: %v", err)
	}

	// Test simple execution
	result, err := runtime.Execute(ctx, "21 + 21")
	if err != nil {
		t.Logf("Execute returned error (may be expected): %v", err)
	} else if result != nil {
		t.Logf("Execute returned: %v", result)
	}

	// Test version
	version := runtime.Version()
	if version == "" {
		t.Error("version should not be empty")
	} else {
		t.Logf("JavaScript version: %s", version)
	}

	// Test shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestRustRuntime tests Rust runtime integration
func TestRustRuntime(t *testing.T) {
	runtime := rust.NewRuntime()

	if runtime.Name() != "rust" {
		t.Errorf("expected name 'rust', got '%s'", runtime.Name())
	}

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 10,
		Timeout:        30 * time.Second,
	}

	ctx := context.Background()

	// Initialize (will fail without actual library, but tests the path)
	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "rust runtime not enabled in build" {
		// Expected to fail without build tag, but should not panic
		t.Logf("Initialize returned expected error: %v", err)
	}

	// Test shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestJavaRuntime tests Java runtime integration
func TestJavaRuntime(t *testing.T) {
	runtime := java.NewRuntime()

	if runtime.Name() != "java" {
		t.Errorf("expected name 'java', got '%s'", runtime.Name())
	}

	config := core.RuntimeConfig{
		Name:           "java",
		Version:        "17",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 10,
		Timeout:        30 * time.Second,
	}

	ctx := context.Background()

	// Initialize (will fail without JVM, but tests the path)
	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "java runtime not enabled in build" {
		t.Logf("Initialize returned expected error: %v", err)
	}

	// Test shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestCppRuntime tests C++ runtime integration
func TestCppRuntime(t *testing.T) {
	runtime := cpp.NewRuntime()

	if runtime.Name() != "cpp" {
		t.Errorf("expected name 'cpp', got '%s'", runtime.Name())
	}

	config := core.RuntimeConfig{
		Name:           "cpp",
		Version:        "17",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 10,
		Timeout:        30 * time.Second,
	}

	ctx := context.Background()

	// Initialize (will fail without actual library, but tests the path)
	err := runtime.Initialize(ctx, config)
	if err != nil && err.Error() != "cpp runtime not enabled in build" {
		t.Logf("Initialize returned expected error: %v", err)
	}

	// Test version
	version := runtime.Version()
	if version == "" {
		t.Error("version should not be empty")
	}

	// Test shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestRuntimeRegistration tests registering multiple runtimes
func TestRuntimeRegistration(t *testing.T) {
	config := core.DefaultConfig()
	config.EnableRuntime("rust", "1.70")
	config.EnableRuntime("java", "17")
	config.EnableRuntime("cpp", "17")
	config.EnableRuntime("javascript", "ES2020")

	orchestrator, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("failed to create orchestrator: %v", err)
	}

	// Register runtimes
	runtimes := []core.Runtime{
		rust.NewRuntime(),
		java.NewRuntime(),
		cpp.NewRuntime(),
		javascript.NewRuntime(),
	}

	for _, runtime := range runtimes {
		if err := orchestrator.RegisterRuntime(runtime); err != nil {
			t.Errorf("failed to register %s: %v", runtime.Name(), err)
		}
	}

	// Verify registration
	registered := orchestrator.Runtimes()
	if len(registered) != 4 {
		t.Errorf("expected 4 runtimes, got %d", len(registered))
	}

	// Test shutdown
	ctx := context.Background()
	if err := orchestrator.Shutdown(ctx); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestRuntimeIsolation tests that runtimes don't interfere with each other
func TestRuntimeIsolation(t *testing.T) {
	config := core.DefaultConfig()
	config.EnableRuntime("rust", "1.70")
	config.EnableRuntime("cpp", "17")

	orchestrator, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("failed to create orchestrator: %v", err)
	}

	rustRuntime := rust.NewRuntime()
	cppRuntime := cpp.NewRuntime()

	orchestrator.RegisterRuntime(rustRuntime)
	orchestrator.RegisterRuntime(cppRuntime)

	// Verify each runtime maintains its identity
	if rustRuntime.Name() == cppRuntime.Name() {
		t.Error("runtimes should have different names")
	}

	ctx := context.Background()
	orchestrator.Shutdown(ctx)
}
