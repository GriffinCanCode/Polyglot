//go:build runtime_zig
// +build runtime_zig

package zig

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"
)

// Worker represents a Zig execution context
type Worker struct {
	id       int
	mu       sync.Mutex
	shutdown bool
	loader   *Loader
	tempDir  string
	zigPath  string
}

// NewWorker creates a Zig worker
func NewWorker(id int) *Worker {
	return &Worker{
		id:      id,
		loader:  NewLoader(),
		zigPath: "zig",
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Verify zig is available (optional - can work with pre-compiled libs only)
	cmd := exec.Command(w.zigPath, "version")
	if err := cmd.Run(); err != nil {
		// Zig not available, but we can still load pre-compiled libraries
		w.zigPath = ""
	}

	// Create temp directory for Zig compilation
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("polyglot-zig-%d-*", w.id))
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	w.tempDir = tempDir

	return nil
}

// Execute runs Zig code (compiles and executes, or calls pre-loaded function)
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// If code is a function name and library is loaded, call it
	if w.loader.IsLoaded() && !strings.Contains(code, "\n") && !strings.Contains(code, "{") {
		return w.callFunction(code, args...)
	}

	// Otherwise, compile and run the code
	return w.compileAndRun(code, args...)
}

// Call invokes a Zig function by name
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	return w.callFunction(fn, args...)
}

// LoadLibrary loads a Zig shared library
func (w *Worker) LoadLibrary(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	return w.loader.Load(path)
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return
	}

	w.shutdown = true

	// Unload library
	if w.loader != nil {
		w.loader.Unload()
	}

	// Clean up temp directory
	if w.tempDir != "" {
		os.RemoveAll(w.tempDir)
	}
}

// callFunction calls a function from the loaded library
func (w *Worker) callFunction(fn string, args ...interface{}) (interface{}, error) {
	if !w.loader.IsLoaded() {
		return nil, fmt.Errorf("no library loaded")
	}

	symbol, err := w.loader.Symbol(fn)
	if err != nil {
		return nil, fmt.Errorf("function %s not found: %w", fn, err)
	}

	// Simple FFI call - handles basic types
	if len(args) == 0 {
		// Call with no arguments, returns int64
		type fn0 func() int64
		f := *(*fn0)(unsafe.Pointer(&symbol))
		result := f()
		return result, nil
	}

	if len(args) == 1 {
		switch v := args[0].(type) {
		case int:
			type fn1 func(int) int64
			f := *(*fn1)(unsafe.Pointer(&symbol))
			result := f(v)
			return result, nil
		case int64:
			type fn1 func(int64) int64
			f := *(*fn1)(unsafe.Pointer(&symbol))
			result := f(v)
			return result, nil
		case float64:
			type fn1 func(float64) float64
			f := *(*fn1)(unsafe.Pointer(&symbol))
			result := f(v)
			return result, nil
		}
	}

	// For now, only support simple arguments
	return nil, fmt.Errorf("complex argument handling not yet implemented")
}

// compileAndRun compiles Zig code and executes it
func (w *Worker) compileAndRun(code string, args ...interface{}) (interface{}, error) {
	if w.zigPath == "" {
		return nil, fmt.Errorf("zig not available for compilation")
	}

	// Prepare the code (wrap in a main function if needed)
	fullCode := w.prepareCode(code)

	// Write to a temporary Zig file
	zigFile := filepath.Join(w.tempDir, fmt.Sprintf("main_%d.zig", w.id))
	if err := os.WriteFile(zigFile, []byte(fullCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Zig file: %w", err)
	}
	defer os.Remove(zigFile)

	// Compile and run the Zig code
	binaryFile := filepath.Join(w.tempDir, fmt.Sprintf("main_%d", w.id))
	var compileStderr bytes.Buffer
	compileCmd := exec.Command(w.zigPath, "build-exe", zigFile, "-femit-bin="+binaryFile)
	compileCmd.Stderr = &compileStderr

	if err := compileCmd.Run(); err != nil {
		errMsg := compileStderr.String()
		if errMsg != "" {
			return nil, fmt.Errorf("compilation failed: %s", errMsg)
		}
		return nil, fmt.Errorf("compilation failed: %w", err)
	}

	// Clean up compiled binary after execution
	defer os.Remove(binaryFile)

	// Execute the compiled binary
	var stdout, stderr bytes.Buffer
	runCmd := exec.Command(binaryFile)
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	if err := runCmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg != "" {
			return nil, fmt.Errorf("execution failed: %s", errMsg)
		}
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	// Extract result from output (check both stdout and stderr since std.debug.print goes to stderr)
	output := stdout.String()
	if output == "" {
		output = stderr.String() // std.debug.print outputs to stderr
	}
	return extractResult(output), nil
}

// prepareCode wraps the code in a proper Zig structure if needed
func (w *Worker) prepareCode(code string) string {
	code = strings.TrimSpace(code)

	// If code already contains main function, use as-is
	if strings.Contains(code, "pub fn main(") {
		return code
	}

	// Build complete Zig program
	var sb strings.Builder

	// Only add std import if not already present
	if !strings.Contains(code, "@import(\"std\")") {
		sb.WriteString("const std = @import(\"std\");\n\n")
	}

	// Check if it's already a statement with print
	if strings.Contains(code, "std.debug.print") || strings.Contains(code, "stdout.print") {
		// Already has print statement - just wrap in main
		sb.WriteString("pub fn main() !void {\n")
		lines := strings.Split(code, "\n")
		for _, line := range lines {
			if line != "" {
				sb.WriteString("    ")
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
		sb.WriteString("}\n")
	} else if !strings.Contains(code, "{") && !strings.Contains(code, ";") && !strings.Contains(code, "fn ") {
		// Simple expression - wrap in main with print (Zig 0.15+ API)
		sb.WriteString("pub fn main() !void {\n")
		sb.WriteString("    std.debug.print(\"{d}\\n\", .{")
		sb.WriteString(code)
		sb.WriteString("});\n")
		sb.WriteString("}\n")
	} else if strings.HasPrefix(code, "fn ") && !strings.Contains(code, "pub fn main") {
		// Function definition(s) - add a main that might call it
		sb.WriteString(code)
		sb.WriteString("\n\npub fn main() !void {\n")
		sb.WriteString("    // Add your function call here\n")
		sb.WriteString("}\n")
	} else {
		// Statements - wrap in main
		sb.WriteString("pub fn main() !void {\n")
		lines := strings.Split(code, "\n")
		for _, line := range lines {
			if line != "" {
				sb.WriteString("    ")
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
		sb.WriteString("}\n")
	}

	return sb.String()
}

// extractResult extracts the result from Zig output
func extractResult(output string) interface{} {
	output = strings.TrimSpace(output)

	if output == "" {
		return nil
	}

	// Return as string for now
	// Could parse into specific types based on format
	return output
}
