//go:build runtime_python
// +build runtime_python

package python

// #include <Python.h>
import "C"

import (
	"errors"
	"sync"
)

// Common errors
var (
	ErrShutdown      = errors.New("runtime is shutdown")
	ErrWorkerBusy    = errors.New("worker is busy")
	ErrCompileFailed = errors.New("code compilation failed")
	ErrExecFailed    = errors.New("code execution failed")
	ErrCallFailed    = errors.New("function call failed")
	ErrNotFound      = errors.New("object not found")
)

// State represents a Python execution context
type State struct {
	id       int
	gil      *C.PyGILState_STATE
	globals  *C.PyObject
	locals   *C.PyObject
	busy     bool
	shutdown bool
	mu       sync.Mutex
}

// Result represents execution result
type Result struct {
	Value interface{}
	Err   error
}

// CallParams represents function call parameters
type CallParams struct {
	Name string
	Args []interface{}
}

// ExecParams represents execution parameters
type ExecParams struct {
	Code string
	Args []interface{}
}
