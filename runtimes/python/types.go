//go:build runtime_python
// +build runtime_python

package python

// #include <Python.h>
import "C"

import (
	"errors"
	"fmt"
	"sync"
)

// Common errors
var (
	ErrShutdown         = errors.New("runtime is shutdown")
	ErrWorkerBusy       = errors.New("worker is busy")
	ErrCompileFailed    = errors.New("code compilation failed")
	ErrExecFailed       = errors.New("code execution failed")
	ErrCallFailed       = errors.New("function call failed")
	ErrNotFound         = errors.New("object not found")
	ErrTypeConversion   = errors.New("type conversion failed")
	ErrImportFailed     = errors.New("module import failed")
	ErrAttributeError   = errors.New("attribute error")
	ErrInitFailed       = errors.New("initialization failed")
	ErrInvalidArguments = errors.New("invalid arguments")
)

// PythonError represents a detailed Python error with context
type PythonError struct {
	Type      string // Python exception type (e.g., "ValueError", "TypeError")
	Message   string // Error message
	Traceback string // Full traceback if available
	Code      string // Code that caused the error
	Line      int    // Line number where error occurred
}

func (e *PythonError) Error() string {
	msg := e.Type + ": " + e.Message
	if e.Code != "" {
		msg += "\nCode: " + e.Code
	}
	if e.Line > 0 {
		msg += fmt.Sprintf("\nLine: %d", e.Line)
	}
	if e.Traceback != "" {
		msg += "\n" + e.Traceback
	}
	return msg
}

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
