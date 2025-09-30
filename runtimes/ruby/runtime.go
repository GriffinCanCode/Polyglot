//go:build runtime_ruby
// +build runtime_ruby

package ruby

/*
#cgo pkg-config: ruby
#cgo LDFLAGS: -lruby
#include <ruby.h>
#include <ruby/encoding.h>
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"fmt"
	"sync"
	"unsafe"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements Ruby runtime integration
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a Ruby runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		pool: NewPool(10),
	}
}

// Initialize prepares the Ruby runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Initialize Ruby interpreter
	C.ruby_init()
	C.ruby_init_loadpath()

	// Initialize the pool
	if err := r.pool.Initialize(config.MaxConcurrency); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	return nil
}

// Execute runs Ruby code
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

// Call invokes a Ruby method
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

	// Finalize Ruby interpreter
	C.ruby_finalize()

	return nil
}

type result struct {
	value interface{}
	err   error
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "ruby"
}

// Version returns the Ruby version
func (r *Runtime) Version() string {
	version := C.ruby_version
	return C.GoString(version)
}

// rubyStringToGo converts Ruby VALUE to Go string
func rubyStringToGo(val C.VALUE) string {
	cStr := C.StringValueCStr(val)
	return C.GoString(cStr)
}

// goStringToRuby converts Go string to Ruby VALUE
func goStringToRuby(s string) C.VALUE {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	return C.rb_str_new_cstr(cStr)
}

// convertToRuby converts Go value to Ruby VALUE
func convertToRuby(val interface{}) C.VALUE {
	if val == nil {
		return C.Qnil
	}

	switch v := val.(type) {
	case string:
		return goStringToRuby(v)
	case int:
		return C.INT2NUM(C.int(v))
	case int64:
		return C.LL2NUM(C.longlong(v))
	case float64:
		return C.DBL2NUM(C.double(v))
	case bool:
		if v {
			return C.Qtrue
		}
		return C.Qfalse
	default:
		return C.Qnil
	}
}

// convertFromRuby converts Ruby VALUE to Go value
func convertFromRuby(val C.VALUE) interface{} {
	switch C.TYPE(val) {
	case C.T_NIL:
		return nil
	case C.T_TRUE:
		return true
	case C.T_FALSE:
		return false
	case C.T_FIXNUM:
		return int64(C.NUM2LL(val))
	case C.T_FLOAT:
		return float64(C.NUM2DBL(val))
	case C.T_STRING:
		return rubyStringToGo(val)
	default:
		return nil
	}
}
