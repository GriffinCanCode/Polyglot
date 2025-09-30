//go:build runtime_cpp
// +build runtime_cpp

package cpp

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Worker represents a C++ execution context
type Worker struct {
	id       int
	mu       sync.Mutex
	shutdown bool
	cppPath  string
	tempDir  string
}

// NewWorker creates a C++ worker
func NewWorker(id int) *Worker {
	return &Worker{
		id:      id,
		cppPath: "g++",
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Verify g++ is available
	cmd := exec.Command(w.cppPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("g++ not available: %w", err)
	}

	// Create temp directory for C++ binaries
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("polyglot-cpp-%d-*", w.id))
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	w.tempDir = tempDir

	return nil
}

// Execute runs C++ code
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Prepare the code (wrap in a main function if needed)
	fullCode := w.prepareCode(code)

	// Write to a temporary C++ file
	cppFile := filepath.Join(w.tempDir, fmt.Sprintf("main_%d.cpp", w.id))
	if err := os.WriteFile(cppFile, []byte(fullCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write C++ file: %w", err)
	}
	defer os.Remove(cppFile)

	// Compile the C++ code
	binaryFile := filepath.Join(w.tempDir, fmt.Sprintf("main_%d", w.id))
	var compileStderr bytes.Buffer
	compileCmd := exec.Command(w.cppPath, "-std=c++17", "-o", binaryFile, cppFile)
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

// Call invokes a C++ function (by compiling and running)
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Build C++ code that calls the function
	var argStrs []string
	for _, arg := range args {
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}

	code := fmt.Sprintf(`
#include <iostream>

int %s(%s) {
    // Placeholder - this would need actual function implementation
    return 0;
}

int main() {
    std::cout << %s(%s) << std::endl;
    return 0;
}
`, fn, buildParamList(len(args)), fn, strings.Join(argStrs, ", "))

	// Execute the code
	return w.Execute(code)
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.shutdown = true

	// Clean up temp directory
	if w.tempDir != "" {
		os.RemoveAll(w.tempDir)
	}
}

// prepareCode wraps the code in a proper C++ structure if needed
func (w *Worker) prepareCode(code string) string {
	code = strings.TrimSpace(code)

	// If code already contains main function, use as-is
	if strings.Contains(code, "int main(") || strings.Contains(code, "int main (") {
		return code
	}

	// Check if code has includes
	hasIncludes := strings.Contains(code, "#include")

	// Build complete C++ program
	var sb strings.Builder

	// Add iostream if not already included
	if !strings.Contains(code, "#include <iostream>") {
		sb.WriteString("#include <iostream>\n")
	}

	// Add other standard includes
	if !strings.Contains(code, "#include <string>") {
		sb.WriteString("#include <string>\n")
	}
	if !strings.Contains(code, "#include <vector>") {
		sb.WriteString("#include <vector>\n")
	}

	// Add using namespace std for convenience
	sb.WriteString("using namespace std;\n\n")

	// If code has includes, extract and place them at the top
	if hasIncludes {
		lines := strings.Split(code, "\n")
		var includeLines, codeLines []string
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "#include") {
				includeLines = append(includeLines, line)
			} else {
				codeLines = append(codeLines, line)
			}
		}
		// Insert includes at the beginning
		code = strings.Join(codeLines, "\n")
	}

	// If it's just an expression, wrap it in main with cout
	if !strings.Contains(code, "{") && !strings.Contains(code, ";") {
		sb.WriteString("int main() {\n")
		sb.WriteString("    cout << ")
		sb.WriteString(code)
		sb.WriteString(" << endl;\n")
		sb.WriteString("    return 0;\n")
		sb.WriteString("}\n")
	} else if !strings.Contains(code, "int main") {
		// It's statements/functions, wrap in main
		sb.WriteString("int main() {\n")
		sb.WriteString(ensureIndented(code))
		sb.WriteString("\n    return 0;\n")
		sb.WriteString("}\n")
	} else {
		// Already has structure, just add it
		sb.WriteString(code)
	}

	return sb.String()
}

// ensureIndented adds indentation to code
func ensureIndented(code string) string {
	lines := strings.Split(code, "\n")
	var sb strings.Builder
	for _, line := range lines {
		if line != "" {
			sb.WriteString("    ")
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// buildParamList builds a parameter list for a function
func buildParamList(count int) string {
	if count == 0 {
		return ""
	}

	var params []string
	for i := 0; i < count; i++ {
		params = append(params, fmt.Sprintf("int arg%d", i))
	}
	return strings.Join(params, ", ")
}

// extractResult extracts the result from C++ output
func extractResult(output string) interface{} {
	output = strings.TrimSpace(output)

	if output == "" {
		return nil
	}

	// Try to parse as different types
	// For now, return as string
	return output
}
