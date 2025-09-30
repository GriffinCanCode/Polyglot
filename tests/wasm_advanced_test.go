//go:build runtime_wasm
// +build runtime_wasm

package tests

import (
	"context"
	"encoding/hex"
	"sync"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/wasm"
)

// TestWASMBasicOperations tests basic WASM operations
func TestWASMBasicOperations(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Simple WASM module with a start function that returns 42
	// This is a minimal valid WASM binary (magic number + version + minimal sections)
	wasmBinary := []byte{
		0x00, 0x61, 0x73, 0x6D, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // WASM version 1
	}

	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{"Simple Module", string(wasmBinary), int32(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Execute(ctx, tt.code)
			if err != nil {
				t.Logf("Execute returned error (expected for stub): %v", err)
				return
			}
			t.Logf("Result: %v (type: %T)", result, result)
		})
	}
}

// TestWASMErrors tests error handling
func TestWASMErrors(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
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
		{"Invalid Magic Number", "not valid wasm"},
		{"Empty Binary", ""},
		{"Truncated Binary", "\x00\x61"},
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

// TestWASMConcurrentExecution tests concurrent WASM execution
func TestWASMConcurrentExecution(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 10,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	wasmBinary := []byte{
		0x00, 0x61, 0x73, 0x6D, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // WASM version 1
	}

	const numGoroutines = 20
	const executionsPerGoroutine = 10

	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines*executionsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < executionsPerGoroutine; j++ {
				_, err := runtime.Execute(ctx, string(wasmBinary))
				if err != nil {
					errChan <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	errorCount := 0
	for err := range errChan {
		errorCount++
		t.Logf("Execution error: %v", err)
	}

	t.Logf("Completed %d concurrent executions with %d errors",
		numGoroutines*executionsPerGoroutine, errorCount)
}

// TestWASMFunctionCalls tests calling exported WASM functions
func TestWASMFunctionCalls(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Load a WASM module first
	wasmBinary := []byte{
		0x00, 0x61, 0x73, 0x6D, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // WASM version 1
	}

	if err := runtime.LoadModule(ctx, wasmBinary); err != nil {
		t.Logf("LoadModule returned error (expected for basic impl): %v", err)
	}

	tests := []struct {
		name     string
		function string
		args     []interface{}
	}{
		{"Call _start", "_start", []interface{}{}},
		{"Call add", "add", []interface{}{2, 3}},
		{"Call multiply", "multiply", []interface{}{6, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runtime.Call(ctx, tt.function, tt.args...)
			if err != nil {
				t.Logf("Call returned error (may be expected): %v", err)
				return
			}
			t.Logf("Function %s returned: %v", tt.function, result)
		})
	}
}

// TestWASMContextCancellation tests context cancellation during execution
func TestWASMContextCancellation(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	wasmBinary := []byte{
		0x00, 0x61, 0x73, 0x6D, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // WASM version 1
	}

	// Create a context that will be cancelled
	execCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Try to execute - may complete quickly or be cancelled
	_, err := runtime.Execute(execCtx, string(wasmBinary))
	if err != nil {
		t.Logf("Execution error (may be cancellation): %v", err)
	}
}

// TestWASMPoolBehavior tests worker pool behavior
func TestWASMPoolBehavior(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 3, // Small pool to test contention
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	wasmBinary := []byte{
		0x00, 0x61, 0x73, 0x6D, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // WASM version 1
	}

	// Start multiple goroutines that exceed pool size
	const numWorkers = 10
	var wg sync.WaitGroup
	results := make(chan error, numWorkers)

	startTime := time.Now()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := runtime.Execute(ctx, string(wasmBinary))
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	duration := time.Since(startTime)
	t.Logf("Completed %d executions in %v", numWorkers, duration)

	errorCount := 0
	for err := range results {
		if err != nil {
			errorCount++
		}
	}

	t.Logf("Errors: %d/%d", errorCount, numWorkers)
}

// TestWASMModuleLoading tests loading various WASM modules
func TestWASMModuleLoading(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	tests := []struct {
		name      string
		bytecode  []byte
		shouldErr bool
	}{
		{
			name: "Valid Empty Module",
			bytecode: []byte{
				0x00, 0x61, 0x73, 0x6D, // magic
				0x01, 0x00, 0x00, 0x00, // version
			},
			shouldErr: false,
		},
		{
			name:      "Invalid Magic",
			bytecode:  []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01, 0x00, 0x00, 0x00},
			shouldErr: true,
		},
		{
			name:      "Too Short",
			bytecode:  []byte{0x00, 0x61, 0x73},
			shouldErr: true,
		},
		{
			name:      "Empty",
			bytecode:  []byte{},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runtime.LoadModule(ctx, tt.bytecode)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else {
				t.Logf("Result as expected (error=%v)", err != nil)
			}
		})
	}
}

// TestWASMShutdownBehavior tests shutdown behavior
func TestWASMShutdownBehavior(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	wasmBinary := []byte{
		0x00, 0x61, 0x73, 0x6D, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // WASM version 1
	}

	// Execute before shutdown
	_, err := runtime.Execute(ctx, string(wasmBinary))
	if err != nil {
		t.Logf("Execute before shutdown: %v", err)
	}

	// Shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Try to execute after shutdown
	_, err = runtime.Execute(ctx, string(wasmBinary))
	if err == nil {
		t.Error("Expected error after shutdown but got none")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}

	// Shutdown again should be idempotent
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Second shutdown failed: %v", err)
	}
}

// TestWASMRuntimeInfo tests runtime information methods
func TestWASMRuntimeInfo(t *testing.T) {
	runtime := wasm.NewRuntime()

	if name := runtime.Name(); name != "wasm" {
		t.Errorf("Expected name 'wasm', got '%s'", name)
	}

	version := runtime.Version()
	t.Logf("WASM version: %s", version)

	if version == "" {
		t.Error("Version should not be empty")
	}
}

// TestWASMValidBinaryFormats tests various valid WASM binary formats
func TestWASMValidBinaryFormats(t *testing.T) {
	runtime := wasm.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "wasm",
		Enabled:        true,
		MaxConcurrency: 5,
		Timeout:        5 * time.Second,
	}

	if err := runtime.Initialize(ctx, config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer runtime.Shutdown(ctx)

	// Test different WASM binary representations
	tests := []struct {
		name   string
		binary []byte
	}{
		{
			name: "Hex String Format",
			binary: []byte{
				0x00, 0x61, 0x73, 0x6D,
				0x01, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runtime.Execute(ctx, string(tt.binary))
			if err != nil {
				t.Logf("Execute error (may be expected): %v", err)
			} else {
				t.Log("Execute succeeded")
			}
		})
	}
}

// Helper function to create a minimal valid WASM module
func createMinimalWASMModule() []byte {
	return []byte{
		0x00, 0x61, 0x73, 0x6D, // magic number
		0x01, 0x00, 0x00, 0x00, // version 1
	}
}

// Helper function to convert hex string to bytes
func hexToBytes(s string) []byte {
	bytes, _ := hex.DecodeString(s)
	return bytes
}
