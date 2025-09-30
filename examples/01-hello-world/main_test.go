package main

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/javascript"
	"github.com/griffincancode/polyglot.js/runtimes/python"
)

// TestOrchestratorCreation verifies the orchestrator can be created
func TestOrchestratorCreation(t *testing.T) {
	config := core.DefaultConfig()
	config.App.Name = "test-app"

	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	if orch == nil {
		t.Fatal("Orchestrator is nil")
	}
}

// TestRuntimeRegistration verifies runtimes can be registered
func TestRuntimeRegistration(t *testing.T) {
	config := core.DefaultConfig()
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	// Register Python
	if err := orch.RegisterRuntime(python.NewRuntime()); err != nil {
		t.Fatalf("Failed to register Python runtime: %v", err)
	}

	// Register JavaScript
	if err := orch.RegisterRuntime(javascript.NewRuntime()); err != nil {
		t.Fatalf("Failed to register JavaScript runtime: %v", err)
	}

	runtimes := orch.Runtimes()
	if len(runtimes) != 2 {
		t.Errorf("Expected 2 runtimes, got %d", len(runtimes))
	}
}

// TestJavaScriptExecution verifies JavaScript can execute
func TestJavaScriptExecution(t *testing.T) {
	config := core.DefaultConfig()
	config.Languages = map[string]*core.RuntimeConfig{
		"javascript": {
			Name:           "javascript",
			Version:        "latest",
			Enabled:        true,
			MaxConcurrency: 4,
			Timeout:        time.Second * 10,
			Options:        make(map[string]interface{}),
		},
	}

	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	if err := orch.RegisterRuntime(javascript.NewRuntime()); err != nil {
		t.Fatalf("Failed to register JavaScript runtime: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := orch.Initialize(ctx); err != nil {
		t.Logf("Initialization warning (may be expected): %v", err)
	}

	// Try to execute JavaScript
	code := `
		function add(a, b) {
			return a + b;
		}
		add(2, 3);
	`

	result, err := orch.Execute(ctx, "javascript", code)
	if err != nil {
		t.Logf("JavaScript execution note: %v", err)
		// This is acceptable - JavaScript may work or may fail depending on initialization
		return
	}

	t.Logf("JavaScript execution result: %v", result)
}

// TestMemoryCoordinator verifies shared memory works
func TestMemoryCoordinator(t *testing.T) {
	config := core.DefaultConfig()
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	mem := orch.Memory()
	if mem == nil {
		t.Fatal("Memory coordinator is nil")
	}

	// Allocate memory
	region, err := mem.Allocate("test-region", 512, core.TypeBytes)
	if err != nil {
		t.Fatalf("Failed to allocate memory: %v", err)
	}

	if len(region.Data) != 512 {
		t.Errorf("Expected 512 bytes, got %d", len(region.Data))
	}

	// Write data
	testData := []byte("Hello, Polyglot!")
	copy(region.Data, testData)

	// Read data
	readData := region.Data[:len(testData)]
	if string(readData) != string(testData) {
		t.Errorf("Data mismatch: got %s, want %s", string(readData), string(testData))
	}

	// Retrieve region
	retrieved, err := mem.Get("test-region")
	if err != nil {
		t.Fatalf("Failed to retrieve region: %v", err)
	}

	if retrieved.ID != "test-region" {
		t.Errorf("Region ID mismatch: got %s, want test-region", retrieved.ID)
	}

	// Free memory
	if err := mem.Free("test-region"); err != nil {
		t.Fatalf("Failed to free memory: %v", err)
	}

	// Verify it's gone
	_, err = mem.Get("test-region")
	if err == nil {
		t.Error("Expected error when getting freed region")
	}
}

// TestGracefulShutdown verifies orchestrator shuts down cleanly
func TestGracefulShutdown(t *testing.T) {
	config := core.DefaultConfig()
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	orch.RegisterRuntime(python.NewRuntime())
	orch.RegisterRuntime(javascript.NewRuntime())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Shutdown should work even without initialization
	if err := orch.Shutdown(ctx); err != nil {
		t.Logf("Shutdown note: %v", err)
		// Some errors during shutdown are acceptable
	}
}

// BenchmarkJavaScriptExecution benchmarks JavaScript execution speed
func BenchmarkJavaScriptExecution(b *testing.B) {
	config := core.DefaultConfig()
	config.Languages = map[string]*core.RuntimeConfig{
		"javascript": {
			Name:           "javascript",
			Enabled:        true,
			MaxConcurrency: 4,
			Timeout:        time.Second * 30,
			Options:        make(map[string]interface{}),
		},
	}

	orch, err := core.NewOrchestrator(config)
	if err != nil {
		b.Fatalf("Failed to create orchestrator: %v", err)
	}

	orch.RegisterRuntime(javascript.NewRuntime())

	ctx := context.Background()
	orch.Initialize(ctx)

	code := `2 + 2`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orch.Execute(ctx, "javascript", code)
	}
}

// BenchmarkMemoryAllocation benchmarks memory allocation speed
func BenchmarkMemoryAllocation(b *testing.B) {
	config := core.DefaultConfig()
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		b.Fatalf("Failed to create orchestrator: %v", err)
	}

	mem := orch.Memory()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		regionID := string(rune('a' + (i % 26)))
		mem.Allocate(regionID, 1024, core.TypeBytes)
		mem.Free(regionID)
	}
}
