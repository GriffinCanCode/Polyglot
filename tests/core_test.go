package tests

import (
	"context"
	"testing"
	"time"

	"github.com/polyglot-framework/polyglot/core"
)

func TestDefaultConfig(t *testing.T) {
	config := core.DefaultConfig()

	if config.App.Name != "polyglot-app" {
		t.Errorf("Expected app name 'polyglot-app', got '%s'", config.App.Name)
	}

	if config.Memory.MaxSharedMemory <= 0 {
		t.Error("Max shared memory should be positive")
	}
}

func TestConfigValidation(t *testing.T) {
	config := core.DefaultConfig()

	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid: %v", err)
	}

	// Test invalid config
	config.App.Name = ""
	if err := config.Validate(); err == nil {
		t.Error("Expected validation error for empty app name")
	}
}

func TestEnableRuntime(t *testing.T) {
	config := core.DefaultConfig()

	config.EnableRuntime("python", "3.11")

	if !config.IsRuntimeEnabled("python") {
		t.Error("Python runtime should be enabled")
	}

	if rtConfig := config.Languages["python"]; rtConfig.Version != "3.11" {
		t.Errorf("Expected version '3.11', got '%s'", rtConfig.Version)
	}
}

func TestDisableRuntime(t *testing.T) {
	config := core.DefaultConfig()

	config.EnableRuntime("python", "3.11")
	config.DisableRuntime("python")

	if config.IsRuntimeEnabled("python") {
		t.Error("Python runtime should be disabled")
	}
}

func TestOrchestrator(t *testing.T) {
	config := core.DefaultConfig()

	orch, err := core.NewOrchestrator(config)
	if err != nil {
		t.Fatalf("Failed to create orchestrator: %v", err)
	}

	if orch == nil {
		t.Fatal("Orchestrator should not be nil")
	}
}

func TestMemoryCoordinator(t *testing.T) {
	memConfig := core.MemoryConfig{
		MaxSharedMemory: 1024 * 1024, // 1MB
		EnableZeroCopy:  true,
		GCInterval:      time.Minute,
	}

	mem := core.NewMemoryCoordinator(memConfig)

	// Test allocation
	region, err := mem.Allocate("test", 1024, core.TypeBytes)
	if err != nil {
		t.Fatalf("Failed to allocate memory: %v", err)
	}

	if region.ID != "test" {
		t.Errorf("Expected ID 'test', got '%s'", region.ID)
	}

	if len(region.Data) != 1024 {
		t.Errorf("Expected size 1024, got %d", len(region.Data))
	}

	// Test retrieval
	retrieved, err := mem.Get("test")
	if err != nil {
		t.Fatalf("Failed to get memory: %v", err)
	}

	if retrieved.ID != region.ID {
		t.Error("Retrieved region should match allocated region")
	}

	// Test freeing
	if err := mem.Free("test"); err != nil {
		t.Fatalf("Failed to free memory: %v", err)
	}

	// Should not be able to get freed region
	if _, err := mem.Get("test"); err == nil {
		t.Error("Should not be able to get freed region")
	}
}

func TestMemoryReadWrite(t *testing.T) {
	memConfig := core.MemoryConfig{
		MaxSharedMemory: 1024 * 1024,
		EnableZeroCopy:  true,
		GCInterval:      time.Minute,
	}

	mem := core.NewMemoryCoordinator(memConfig)

	region, err := mem.Allocate("test", 1024, core.TypeBytes)
	if err != nil {
		t.Fatalf("Failed to allocate memory: %v", err)
	}

	// Acquire read
	if err := mem.AcquireRead("test"); err != nil {
		t.Fatalf("Failed to acquire read: %v", err)
	}

	if region.Readers != 1 {
		t.Errorf("Expected 1 reader, got %d", region.Readers)
	}

	// Release read
	if err := mem.ReleaseRead("test"); err != nil {
		t.Fatalf("Failed to release read: %v", err)
	}

	if region.Readers != 0 {
		t.Errorf("Expected 0 readers, got %d", region.Readers)
	}

	// Acquire write
	if err := mem.AcquireWrite("test"); err != nil {
		t.Fatalf("Failed to acquire write: %v", err)
	}

	if region.Writers != 1 {
		t.Errorf("Expected 1 writer, got %d", region.Writers)
	}

	// Can't acquire another write
	if err := mem.AcquireWrite("test"); err == nil {
		t.Error("Should not be able to acquire multiple writers")
	}

	mem.Free("test")
}

func TestBridge(t *testing.T) {
	bridge := core.NewBridge()

	// Register function
	err := bridge.Register("test", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		return "hello", nil
	})

	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// Call function
	result, err := bridge.Call(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	if result != "hello" {
		t.Errorf("Expected 'hello', got '%v'", result)
	}

	// Unregister function
	if err := bridge.Unregister("test"); err != nil {
		t.Fatalf("Failed to unregister function: %v", err)
	}

	// Should not be able to call unregistered function
	if _, err := bridge.Call(context.Background(), "test"); err == nil {
		t.Error("Should not be able to call unregistered function")
	}
}

func TestBridgeWithArgs(t *testing.T) {
	bridge := core.NewBridge()

	bridge.Register("add", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		a, ok1 := args[0].(float64)
		b, ok2 := args[1].(float64)
		if !ok1 || !ok2 {
			return nil, nil
		}
		return a + b, nil
	})

	result, err := bridge.Call(context.Background(), "add", 5.0, 3.0)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	if result != 8.0 {
		t.Errorf("Expected 8.0, got %v", result)
	}
}
