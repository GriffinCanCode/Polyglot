//go:build runtime_rust
// +build runtime_rust

package rust

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Worker represents a Rust execution context
type Worker struct {
	id        int
	mu        sync.Mutex
	shutdown  bool
	loader    *Loader
	tempDir   string
	rustcPath string
}

// NewWorker creates a Rust worker
func NewWorker(id int) *Worker {
	return &Worker{
		id:        id,
		loader:    NewLoader(),
		rustcPath: "rustc",
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Verify rustc is available (optional - can work with pre-compiled libs only)
	cmd := exec.Command(w.rustcPath, "--version")
	if err := cmd.Run(); err != nil {
		// Rustc not available, but we can still load pre-compiled libraries
		w.rustcPath = ""
	}

	// Create temp directory for Rust compilation
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("polyglot-rust-%d-*", w.id))
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	w.tempDir = tempDir

	return nil
}

// Execute runs Rust code (compiles and executes, or calls pre-loaded function)
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

// Call invokes a Rust function by name
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	return w.callFunction(fn, args...)
}

// LoadLibrary loads a Rust shared library
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
		f := *(*fn0)(symbol)
		result := f()
		return result, nil
	}

	// For now, only support simple integer arguments
	// More complex type handling can be added as needed
	return nil, fmt.Errorf("complex argument handling not yet implemented")
}

// compileAndRun compiles Rust code and executes it
func (w *Worker) compileAndRun(code string, args ...interface{}) (interface{}, error) {
	if w.rustcPath == "" {
		return nil, fmt.Errorf("rustc not available for compilation")
	}

	// Prepare the code (wrap in a main function if needed)
	fullCode := w.prepareCode(code)

	// Write to a temporary Rust file
	rustFile := filepath.Join(w.tempDir, fmt.Sprintf("main_%d.rs", w.id))
	if err := os.WriteFile(rustFile, []byte(fullCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Rust file: %w", err)
	}
	defer os.Remove(rustFile)

	// Compile the Rust code
	binaryFile := filepath.Join(w.tempDir, fmt.Sprintf("main_%d", w.id))
	var compileStderr bytes.Buffer
	compileCmd := exec.Command(w.rustcPath, "-o", binaryFile, rustFile)
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

	// Extract result from output
	output := stdout.String()
	return extractResult(output), nil
}

// prepareCode wraps the code in a proper Rust structure if needed
func (w *Worker) prepareCode(code string) string {
	code = strings.TrimSpace(code)

	// If code already contains main function, use as-is
	if strings.Contains(code, "fn main(") {
		return code
	}

	// Build complete Rust program
	var sb strings.Builder

	// Check if it's already a statement (like println!, let, etc.)
	if strings.Contains(code, "println!") || strings.Contains(code, "print!") {
		// Already has print statement - just wrap in main
		sb.WriteString("fn main() {\n")
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
		// Simple expression - wrap in main with println
		sb.WriteString("fn main() {\n")
		sb.WriteString("    println!(\"{}\", ")
		sb.WriteString(code)
		sb.WriteString(");\n")
		sb.WriteString("}\n")
	} else if strings.HasPrefix(code, "fn ") && !strings.Contains(code, "fn main") {
		// Function definition(s) - add a main that calls it
		sb.WriteString(code)
		sb.WriteString("\n\nfn main() {\n")
		sb.WriteString("    // Add your function call here\n")
		sb.WriteString("}\n")
	} else {
		// Statements - wrap in main
		sb.WriteString("fn main() {\n")
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

// extractResult extracts the result from Rust output
func extractResult(output string) interface{} {
	output = strings.TrimSpace(output)

	if output == "" {
		return nil
	}

	// Return as string for now
	// Could parse into specific types based on format
	return output
}
