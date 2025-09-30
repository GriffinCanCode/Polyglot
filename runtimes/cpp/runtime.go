//go:build runtime_cpp
// +build runtime_cpp

package cpp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements C++ runtime integration using CLI execution
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
	cppPath  string
}

// NewRuntime creates a C++ runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		cppPath: "g++",
	}
}

// Initialize prepares the C++ runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	// Check if g++ is available
	if err := r.checkCppAvailable(); err != nil {
		return fmt.Errorf("C++ compiler not available: %w", err)
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

// checkCppAvailable verifies g++ is installed and accessible
func (r *Runtime) checkCppAvailable() error {
	cmd := exec.Command(r.cppPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("g++ binary not found or not executable: %w", err)
	}
	return nil
}

// Execute runs C++ code
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

// Call invokes a C++ function
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
	return "cpp"
}

// Version returns the C++ compiler version
func (r *Runtime) Version() string {
	cmd := exec.Command(r.cppPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	// g++ --version outputs to stdout
	versionStr := string(output)
	lines := strings.Split(versionStr, "\n")
	if len(lines) > 0 {
		// Extract version from first line
		firstLine := lines[0]
		if strings.Contains(firstLine, "g++") {
			// Extract version number
			parts := strings.Fields(firstLine)
			if len(parts) >= 3 {
				return parts[len(parts)-1]
			}
		}
		return strings.TrimSpace(firstLine)
	}

	return "c++17"
}

type result struct {
	value interface{}
	err   error
}
