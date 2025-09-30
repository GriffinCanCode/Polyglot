//go:build runtime_php
// +build runtime_php

package php

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements PHP runtime integration using CLI execution
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
	phpPath  string
}

// NewRuntime creates a PHP runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		phpPath: "php", // Default PHP binary path
	}
}

// Initialize prepares the PHP runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	// Check if PHP is available
	if err := r.checkPHPAvailable(); err != nil {
		return fmt.Errorf("PHP not available: %w", err)
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

// checkPHPAvailable verifies PHP is installed and accessible
func (r *Runtime) checkPHPAvailable() error {
	cmd := exec.Command(r.phpPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("php binary not found or not executable: %w", err)
	}
	return nil
}

// Execute runs PHP code
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

// Call invokes a PHP function
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
	return "php"
}

// Version returns the PHP version
func (r *Runtime) Version() string {
	cmd := exec.Command(r.phpPath, "-r", "echo PHP_VERSION;")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

type result struct {
	value interface{}
	err   error
}

// prepareCode prepares PHP code for execution with php -r
func prepareCode(code string) string {
	// Remove any existing <?php tags since php -r doesn't need them
	code = strings.TrimSpace(code)

	// Remove <?php opening tag if present
	if strings.HasPrefix(code, "<?php") {
		code = strings.TrimPrefix(code, "<?php")
		code = strings.TrimSpace(code)
	}

	// Remove short opening tag if present
	if strings.HasPrefix(code, "<?") {
		code = strings.TrimPrefix(code, "<?")
		code = strings.TrimSpace(code)
	}

	// Remove closing tag if present
	if strings.HasSuffix(code, "?>") {
		code = strings.TrimSuffix(code, "?>")
		code = strings.TrimSpace(code)
	}

	return code
}

// extractResult extracts the result from PHP output
func extractResult(output string) interface{} {
	output = strings.TrimSpace(output)

	if output == "" {
		return nil
	}

	// Try to parse as different types
	// For now, return as string
	return output
}
