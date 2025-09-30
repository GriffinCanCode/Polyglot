//go:build runtime_python
// +build runtime_python

package python

// #cgo pkg-config: python3
// #cgo LDFLAGS: -lpython3.11
// #include <Python.h>
// #include <stdlib.h>
import "C"

import (
	"context"
	"fmt"
	"sync"
	"unsafe"

	"github.com/polyglot-framework/polyglot/core"
)

// Runtime implements Python runtime integration
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a Python runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		pool: NewPool(10),
	}
}

// Initialize prepares the Python runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Initialize Python interpreter
	if C.Py_IsInitialized() == 0 {
		C.Py_Initialize()
	}

	// Initialize the pool
	if err := r.pool.Initialize(config.MaxConcurrency); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	return nil
}

// Execute runs Python code
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	return worker.Execute(code, args...)
}

// Call invokes a Python function
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	return worker.Call(fn, args...)
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true
	r.pool.Close()

	// Finalize Python interpreter
	if C.Py_IsInitialized() != 0 {
		C.Py_Finalize()
	}

	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "python"
}

// Version returns the Python version
func (r *Runtime) Version() string {
	cVersion := C.Py_GetVersion()
	return C.GoString(cVersion)
}

// pyStringToGo converts Python string to Go string
func pyStringToGo(pyStr *C.PyObject) string {
	if pyStr == nil {
		return ""
	}

	cStr := C.PyUnicode_AsUTF8(pyStr)
	if cStr == nil {
		return ""
	}

	return C.GoString(cStr)
}

// goStringToPy converts Go string to Python string
func goStringToPy(s string) *C.PyObject {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	return C.PyUnicode_FromString(cStr)
}

// convertToPython converts Go value to Python object
func convertToPython(val interface{}) *C.PyObject {
	if val == nil {
		C.Py_IncRef(C.Py_None)
		return C.Py_None
	}

	switch v := val.(type) {
	case string:
		return goStringToPy(v)
	case int:
		return C.PyLong_FromLong(C.long(v))
	case int64:
		return C.PyLong_FromLongLong(C.longlong(v))
	case float64:
		return C.PyFloat_FromDouble(C.double(v))
	case bool:
		if v {
			C.Py_IncRef(C.Py_True)
			return C.Py_True
		}
		C.Py_IncRef(C.Py_False)
		return C.Py_False
	default:
		C.Py_IncRef(C.Py_None)
		return C.Py_None
	}
}

// convertFromPython converts Python object to Go value
func convertFromPython(obj *C.PyObject) interface{} {
	if obj == nil || obj == C.Py_None {
		return nil
	}

	if C.PyBool_Check(obj) != 0 {
		if obj == C.Py_True {
			return true
		}
		return false
	}

	if C.PyLong_Check(obj) != 0 {
		return int64(C.PyLong_AsLongLong(obj))
	}

	if C.PyFloat_Check(obj) != 0 {
		return float64(C.PyFloat_AsDouble(obj))
	}

	if C.PyUnicode_Check(obj) != 0 {
		return pyStringToGo(obj)
	}

	return nil
}
