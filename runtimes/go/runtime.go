//go:build runtime_go
// +build runtime_go

package goruntime

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Runtime implements Go runtime integration using Yaegi interpreter
type Runtime struct {
	config   core.RuntimeConfig
	pool     *InterpreterPool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a Go runtime instance
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize prepares the Go runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Determine pool size
	poolSize := config.MaxConcurrency
	if poolSize <= 0 {
		poolSize = 4
	}

	// Initialize the interpreter pool
	r.pool = NewInterpreterPool(poolSize)
	if err := r.pool.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	return nil
}

// Execute runs Go code
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	interpreter := r.pool.Acquire()
	defer r.pool.Release(interpreter)

	// Prepare code for execution
	code = prepareCode(code)

	// Execute with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := interpreter.Execute(code, args...)
		resultChan <- result{value: res, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.value, res.err
	}
}

// Call invokes a Go function by name
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	interpreter := r.pool.Acquire()
	defer r.pool.Release(interpreter)

	// Call with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := interpreter.Call(fn, args...)
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
	return "go"
}

// Version returns the Go version
func (r *Runtime) Version() string {
	return "1.24 (Yaegi interpreter)"
}

type result struct {
	value interface{}
	err   error
}

// hasPackageDeclaration checks if code has a package declaration
func hasPackageDeclaration(code string) bool {
	// Simple check for "package " at the beginning
	for i := 0; i < len(code) && i < 100; i++ {
		if i+8 <= len(code) && code[i:i+8] == "package " {
			return true
		}
		// Skip whitespace and comments
		if code[i] == ' ' || code[i] == '\t' || code[i] == '\n' || code[i] == '\r' {
			continue
		}
		if code[i] == '/' && i+1 < len(code) && code[i+1] == '/' {
			// Line comment, skip to end of line
			for i < len(code) && code[i] != '\n' {
				i++
			}
			continue
		}
		// Found non-whitespace, non-comment before package
		return false
	}
	return false
}

// hasDeclaration checks if code has a declaration (func, var, type, const)
func hasDeclaration(code string) bool {
	// Look for common declaration keywords
	for i := 0; i < len(code); i++ {
		// Skip whitespace
		for i < len(code) && (code[i] == ' ' || code[i] == '\t' || code[i] == '\n' || code[i] == '\r') {
			i++
		}
		if i >= len(code) {
			break
		}

		// Check for declaration keywords
		if i+4 <= len(code) && code[i:i+4] == "func" {
			return true
		}
		if i+3 <= len(code) && code[i:i+3] == "var" {
			return true
		}
		if i+4 <= len(code) && code[i:i+4] == "type" {
			return true
		}
		if i+5 <= len(code) && code[i:i+5] == "const" {
			return true
		}
		break
	}
	return false
}

// hasShortVarDecl checks if code contains short variable declaration (:=)
func hasShortVarDecl(code string) bool {
	// Look for := in the code
	for i := 0; i < len(code)-1; i++ {
		if code[i] == ':' && code[i+1] == '=' {
			return true
		}
	}
	return false
}

// hasImport checks if code starts with import statement
func hasImport(code string) bool {
	// Skip leading whitespace
	i := 0
	for i < len(code) && (code[i] == ' ' || code[i] == '\t' || code[i] == '\n' || code[i] == '\r') {
		i++
	}
	if i+6 <= len(code) && code[i:i+6] == "import" {
		return true
	}
	return false
}

// prepareCode prepares Go code for execution
func prepareCode(code string) string {
	// If code already has package declaration, use it as-is
	if hasPackageDeclaration(code) {
		return code
	}

	// If code has import statement, wrap in package with import
	if hasImport(code) {
		return "package main\n\n" + code + "\n\nvar _ = 0\n"
	}

	// If code has declarations (func, var, type, const), wrap in package
	if hasDeclaration(code) {
		return "package main\n\n" + code
	}

	// If code has short variable declaration, wrap in a function
	if hasShortVarDecl(code) {
		return "package main\n\nfunc init() {\n\t" + code + "\n}\n"
	}

	// For simple expressions, wrap in a variable assignment at package level
	return fmt.Sprintf("package main\n\nvar _ = %s\n", code)
}

// Interpreter wraps a Yaegi interpreter instance
type Interpreter struct {
	interp *interp.Interpreter
	mu     sync.Mutex
	busy   bool
}

// NewInterpreter creates a new Yaegi interpreter
func NewInterpreter() (*Interpreter, error) {
	i := interp.New(interp.Options{})

	// Use stdlib symbols
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, fmt.Errorf("failed to load stdlib: %w", err)
	}

	return &Interpreter{
		interp: i,
	}, nil
}

// Execute runs Go code in the interpreter
func (i *Interpreter) Execute(code string, args ...interface{}) (interface{}, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Evaluate the code
	result, err := i.interp.Eval(code)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	// Convert reflect.Value to interface{}, dereferencing pointers
	if result.IsValid() {
		// Dereference pointers to get actual values
		for result.Kind() == reflect.Ptr {
			if result.IsNil() {
				return nil, nil
			}
			result = result.Elem()
		}

		if result.CanInterface() {
			return result.Interface(), nil
		}
	}

	return nil, nil
}

// Call invokes a function in the interpreter
func (i *Interpreter) Call(fn string, args ...interface{}) (interface{}, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Get the function value
	fnVal, err := i.interp.Eval(fn)
	if err != nil {
		return nil, fmt.Errorf("function %s not found: %w", fn, err)
	}

	if !fnVal.IsValid() || fnVal.Kind() != reflect.Func {
		return nil, fmt.Errorf("%s is not a function", fn)
	}

	// Convert args to reflect.Value
	in := make([]reflect.Value, len(args))
	for idx, arg := range args {
		in[idx] = reflect.ValueOf(arg)
	}

	// Call the function
	out := fnVal.Call(in)

	// Return the first return value if any
	if len(out) > 0 && out[0].IsValid() && out[0].CanInterface() {
		return out[0].Interface(), nil
	}

	return nil, nil
}
