package tests

import (
	"context"
	"testing"
	"time"

	"github.com/polyglot-framework/polyglot/core"
)

// MockRuntime implements a mock runtime for testing
type MockRuntime struct {
	name    string
	version string
	calls   int
}

func NewMockRuntime(name, version string) *MockRuntime {
	return &MockRuntime{
		name:    name,
		version: version,
	}
}

func (m *MockRuntime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	return nil
}

func (m *MockRuntime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	m.calls++
	return "executed: " + code, nil
}

func (m *MockRuntime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	m.calls++
	return "called: " + fn, nil
}

func (m *MockRuntime) Shutdown(ctx context.Context) error {
	return nil
}

func (m *MockRuntime) Name() string {
	return m.name
}

func (m *MockRuntime) Version() string {
	return m.version
}

func TestOrchestratorWithRuntime(t *testing.T) {
	config := core.DefaultConfig()
	config.EnableRuntime("mock", "1.0")

	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	// Register mock runtime
	mockRuntime := NewMockRuntime("mock", "1.0")
	if err := orch.RegisterRuntime(mockRuntime); err != nil {
		t.Fatalf("Failed to register runtime: %v", err)
	}

	// Initialize
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := orch.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Execute code
	result, err := orch.Execute(ctx, "mock", "test code")
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	if result != "executed: test code" {
		t.Errorf("Unexpected result: %v", result)
	}

	// Call function
	result, err = orch.Call(ctx, "mock", "test_func")
	if err != nil {
		t.Fatalf("Failed to call: %v", err)
	}

	if result != "called: test_func" {
		t.Errorf("Unexpected result: %v", result)
	}

	if mockRuntime.calls != 2 {
		t.Errorf("Expected 2 calls, got %d", mockRuntime.calls)
	}

	// Shutdown
	if err := orch.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to shutdown: %v", err)
	}
}

func TestMultipleRuntimes(t *testing.T) {
	config := core.DefaultConfig()
	config.EnableRuntime("python", "3.11")
	config.EnableRuntime("javascript", "latest")

	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	// Register runtimes
	orch.RegisterRuntime(NewMockRuntime("python", "3.11"))
	orch.RegisterRuntime(NewMockRuntime("javascript", "latest"))

	// Initialize
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := orch.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Check runtimes
	runtimes := orch.Runtimes()
	if len(runtimes) != 2 {
		t.Errorf("Expected 2 runtimes, got %d", len(runtimes))
	}

	// Execute in both
	_, err = orch.Execute(ctx, "python", "print('hello')")
	if err != nil {
		t.Fatalf("Failed to execute python: %v", err)
	}

	_, err = orch.Execute(ctx, "javascript", "console.log('hello')")
	if err != nil {
		t.Fatalf("Failed to execute javascript: %v", err)
	}

	orch.Shutdown(ctx)
}

func TestMemorySharing(t *testing.T) {
	config := core.DefaultConfig()

	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	mem := orch.Memory()

	// Allocate shared memory
	region, err := mem.Allocate("shared", 1024, core.TypeFloat64)
	if err != nil {
		t.Fatalf("Failed to allocate: %v", err)
	}

	// Write some data
	data := []byte{1, 2, 3, 4}
	copy(region.Data, data)

	// Read it back
	retrieved, err := mem.Get("shared")
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}

	if retrieved.Data[0] != 1 || retrieved.Data[1] != 2 {
		t.Error("Data mismatch")
	}

	mem.Free("shared")
}
