//go:build runtime_java
// +build runtime_java

package java

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements Java runtime integration using CLI execution
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
	javaPath string
}

// NewRuntime creates a Java runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		javaPath: "java",
	}
}

// Initialize prepares the Java runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	// Check if Java is available
	if err := r.checkJavaAvailable(); err != nil {
		return fmt.Errorf("Java not available: %w", err)
	}

	r.config = config

	// Determine pool size
	poolSize := config.MaxConcurrency
	if poolSize <= 0 {
		poolSize = 4
	}

	// Initialize the pool
	r.pool = NewPool(poolSize)
	if err := r.pool.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	return nil
}

// checkJavaAvailable verifies Java is installed and accessible
func (r *Runtime) checkJavaAvailable() error {
	cmd := exec.Command(r.javaPath, "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("java binary not found or not executable: %w", err)
	}
	return nil
}

// Execute runs Java code
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	// Execute with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := worker.Execute(code, args...)
		resultChan <- result{value: res, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.value, res.err
	}
}

// Call invokes a Java method
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	// Call with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := worker.Call(fn, args...)
		resultChan <- result{value: res, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.value, res.err
	}
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true

	if r.pool != nil {
		r.pool.Close()
	}

	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "java"
}

// Version returns the Java version
func (r *Runtime) Version() string {
	cmd := exec.Command(r.javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	// Java -version outputs to stderr typically
	versionStr := string(output)
	lines := strings.Split(versionStr, "\n")
	if len(lines) > 0 {
		// Extract version from first line
		firstLine := lines[0]
		if strings.Contains(firstLine, "version") {
			// Extract version number
			start := strings.Index(firstLine, "\"")
			if start >= 0 {
				end := strings.Index(firstLine[start+1:], "\"")
				if end >= 0 {
					return firstLine[start+1 : start+1+end]
				}
			}
		}
		return strings.TrimSpace(firstLine)
	}

	return "unknown"
}

type result struct {
	value interface{}
	err   error
}
