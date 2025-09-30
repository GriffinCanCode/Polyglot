//go:build runtime_python
// +build runtime_python

package python

// #include <Python.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Worker represents a Python interpreter instance
type Worker struct {
	id       int
	state    *C.PyThreadState
	globals  *C.PyObject
	locals   *C.PyObject
	shutdown bool
}

// NewWorker creates a new worker
func NewWorker(id int) *Worker {
	return &Worker{
		id: id,
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Create new thread state
	w.state = C.PyThreadState_Get()

	// Create global and local dictionaries
	w.globals = C.PyDict_New()
	w.locals = C.PyDict_New()

	// Add builtins to globals
	builtins := C.PyEval_GetBuiltins()
	cKey := C.CString("__builtins__")
	defer C.free(unsafe.Pointer(cKey))
	C.PyDict_SetItemString(w.globals, cKey, builtins)

	return nil
}

// Execute runs Python code
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	// Run code
	result := C.PyRun_String(cCode, C.Py_file_input, w.globals, w.locals)
	if result == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("execution failed")
	}
	defer C.Py_DecRef(result)

	return convertFromPython(result), nil
}

// Call invokes a Python function
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Get function object
	cFn := C.CString(fn)
	defer C.free(unsafe.Pointer(cFn))

	fnObj := C.PyDict_GetItemString(w.globals, cFn)
	if fnObj == nil {
		return nil, fmt.Errorf("function %s not found", fn)
	}

	// Convert arguments
	pyArgs := C.PyTuple_New(C.Py_ssize_t(len(args)))
	for i, arg := range args {
		pyArg := convertToPython(arg)
		C.PyTuple_SetItem(pyArgs, C.Py_ssize_t(i), pyArg)
	}
	defer C.Py_DecRef(pyArgs)

	// Call function
	result := C.PyObject_CallObject(fnObj, pyArgs)
	if result == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("call failed")
	}
	defer C.Py_DecRef(result)

	return convertFromPython(result), nil
}

// Shutdown cleans up the worker
func (w *Worker) Shutdown() {
	if w.shutdown {
		return
	}

	w.shutdown = true

	if w.globals != nil {
		C.Py_DecRef(w.globals)
	}

	if w.locals != nil {
		C.Py_DecRef(w.locals)
	}
}
