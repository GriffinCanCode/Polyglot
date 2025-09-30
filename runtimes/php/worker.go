//go:build runtime_php
// +build runtime_php

package php

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// Worker represents a PHP execution context
type Worker struct {
	id       int
	mu       sync.Mutex
	shutdown bool
	phpPath  string
}

// NewWorker creates a PHP worker
func NewWorker(id int) *Worker {
	return &Worker{
		id:      id,
		phpPath: "php",
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Verify PHP is available
	cmd := exec.Command(w.phpPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("php not available: %w", err)
	}

	return nil
}

// Execute runs PHP code
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Prepare the code
	code = prepareCode(code)

	// Execute PHP code using -r flag
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(w.phpPath, "-r", code)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
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

// Call invokes a PHP function
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Build function call string
	var argStrs []string
	for _, arg := range args {
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}

	code := fmt.Sprintf("echo %s(%s);", fn, strings.Join(argStrs, ", "))

	// Execute the function call
	return w.Execute(code)
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.shutdown = true
}
