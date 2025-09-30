//go:build runtime_java
// +build runtime_java

package java

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Worker represents a Java execution context
type Worker struct {
	id       int
	mu       sync.Mutex
	shutdown bool
	javaPath string
	tempDir  string
}

// NewWorker creates a Java worker
func NewWorker(id int) *Worker {
	return &Worker{
		id:       id,
		javaPath: "java",
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Verify Java is available
	cmd := exec.Command(w.javaPath, "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("java not available: %w", err)
	}

	// Create temp directory for Java class files
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("polyglot-java-%d-*", w.id))
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	w.tempDir = tempDir

	return nil
}

// Execute runs Java code
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Prepare the code (wrap in a class if needed)
	className, fullCode := w.prepareCode(code)

	// Write to a temporary Java file
	javaFile := filepath.Join(w.tempDir, className+".java")
	if err := os.WriteFile(javaFile, []byte(fullCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Java file: %w", err)
	}
	defer os.Remove(javaFile)

	// Compile the Java code
	var compileStderr bytes.Buffer
	compileCmd := exec.Command("javac", javaFile)
	compileCmd.Stderr = &compileStderr
	compileCmd.Dir = w.tempDir

	if err := compileCmd.Run(); err != nil {
		errMsg := compileStderr.String()
		if errMsg != "" {
			return nil, fmt.Errorf("compilation failed: %s", errMsg)
		}
		return nil, fmt.Errorf("compilation failed: %w", err)
	}

	// Clean up compiled class file
	classFile := filepath.Join(w.tempDir, className+".class")
	defer os.Remove(classFile)

	// Execute the compiled Java class
	var stdout, stderr bytes.Buffer
	runCmd := exec.Command(w.javaPath, "-cp", w.tempDir, className)
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

// Call invokes a Java method
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Parse function name (could be "ClassName.methodName" or just "methodName")
	parts := strings.Split(fn, ".")
	var className, methodName string

	if len(parts) == 2 {
		className = parts[0]
		methodName = parts[1]
	} else {
		className = "PolyglotRunner"
		methodName = fn
	}

	// Build Java code that calls the method
	var argStrs []string
	for _, arg := range args {
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}

	code := fmt.Sprintf(`
public class %s {
    public static void main(String[] args) {
        System.out.println(%s(%s));
    }
    
    public static Object %s(%s) {
        // Placeholder - this would need actual method implementation
        return null;
    }
}
`, className, methodName, strings.Join(argStrs, ", "), methodName, buildParamList(len(args)))

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

// prepareCode wraps the code in a proper Java class structure if needed
func (w *Worker) prepareCode(code string) (string, string) {
	code = strings.TrimSpace(code)

	// If code already contains a public class, extract its name
	if strings.Contains(code, "public class") {
		// Extract class name
		start := strings.Index(code, "public class") + 13
		end := strings.Index(code[start:], "{")
		if end > 0 {
			className := strings.TrimSpace(code[start : start+end])
			return className, code
		}
	}

	// Otherwise, wrap the code in a main method
	className := fmt.Sprintf("PolyglotRunner_%d", w.id)

	// Check if code is just an expression or statement
	fullCode := fmt.Sprintf(`
public class %s {
    public static void main(String[] args) {
        %s
    }
}
`, className, ensureStatement(code))

	return className, fullCode
}

// ensureStatement ensures the code is a valid statement
func ensureStatement(code string) string {
	code = strings.TrimSpace(code)

	// If it's an expression (doesn't end with ;), wrap it in System.out.println
	if !strings.HasSuffix(code, ";") && !strings.HasSuffix(code, "}") {
		return fmt.Sprintf("System.out.println(%s);", code)
	}

	return code
}

// buildParamList builds a parameter list for a method
func buildParamList(count int) string {
	if count == 0 {
		return ""
	}

	var params []string
	for i := 0; i < count; i++ {
		params = append(params, fmt.Sprintf("Object arg%d", i))
	}
	return strings.Join(params, ", ")
}

// extractResult extracts the result from Java output
func extractResult(output string) interface{} {
	output = strings.TrimSpace(output)

	if output == "" {
		return nil
	}

	// Try to parse as different types
	// For now, return as string
	return output
}
