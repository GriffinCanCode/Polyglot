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

// NewState creates a new Python execution state
func NewState(id int) *State {
	return &State{
		id:       id,
		busy:     false,
		shutdown: false,
	}
}

// Initialize prepares the state with its own dictionaries
func (s *State) Initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.shutdown {
		return ErrShutdown
	}

	gil := AcquireGIL()
	defer gil.Release()

	// Create global and local dictionaries
	s.globals = C.PyDict_New()
	s.locals = C.PyDict_New()

	if s.globals == nil || s.locals == nil {
		return fmt.Errorf("failed to create dictionaries")
	}

	// Add builtins to globals
	builtins := C.PyEval_GetBuiltins()
	if builtins != nil {
		cKey := C.CString("__builtins__")
		C.PyDict_SetItemString(s.globals, cKey, builtins)
		C.free(unsafe.Pointer(cKey))
	}

	return nil
}

// Execute runs Python code and returns result
func (s *State) Execute(code string, args ...interface{}) (interface{}, error) {
	s.mu.Lock()
	if s.shutdown {
		s.mu.Unlock()
		return nil, ErrShutdown
	}
	if s.busy {
		s.mu.Unlock()
		return nil, ErrWorkerBusy
	}
	s.busy = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.busy = false
		s.mu.Unlock()
	}()

	gil := AcquireGIL()
	defer gil.Release()

	// Clear any previous errors
	ClearError()

	// Set arguments if provided
	if len(args) > 0 {
		for i, arg := range args {
			argName := fmt.Sprintf("arg%d", i)
			cArgName := C.CString(argName)
			pyArg := ToPython(arg)
			C.PyDict_SetItemString(s.locals, cArgName, pyArg)
			C.Py_DecRef(pyArg)
			C.free(unsafe.Pointer(cArgName))
		}
	}

	// Try eval mode first for expressions
	cCode := C.CString(code)
	cFilename := C.CString("<string>")

	// First try as an expression (eval mode) - this returns the value
	compiled := C.Py_CompileString(cCode, cFilename, C.Py_eval_input)

	// If compilation fails, try as statements (exec mode)
	if compiled == nil {
		ClearError()
		compiled = C.Py_CompileString(cCode, cFilename, C.Py_file_input)
	}

	C.free(unsafe.Pointer(cCode))
	C.free(unsafe.Pointer(cFilename))

	if compiled == nil {
		return nil, fmt.Errorf("%w: %s", ErrCompileFailed, GetError())
	}
	defer C.Py_DecRef(compiled)

	// Execute compiled code
	result := C.PyEval_EvalCode(compiled, s.globals, s.locals)
	if result == nil {
		return nil, fmt.Errorf("%w: %s", ErrExecFailed, GetError())
	}
	defer C.Py_DecRef(result)

	// If result is None, code was probably exec mode (statements)
	// In that case, return nil
	if result == C.Py_None {
		return nil, nil
	}

	return FromPython(result), nil
}

// Call invokes a Python function by name
func (s *State) Call(fn string, args ...interface{}) (interface{}, error) {
	s.mu.Lock()
	if s.shutdown {
		s.mu.Unlock()
		return nil, ErrShutdown
	}
	if s.busy {
		s.mu.Unlock()
		return nil, ErrWorkerBusy
	}
	s.busy = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.busy = false
		s.mu.Unlock()
	}()

	gil := AcquireGIL()
	defer gil.Release()

	// Clear any previous errors
	ClearError()

	// Get function object
	cFn := C.CString(fn)
	fnObj := C.PyDict_GetItemString(s.globals, cFn)
	C.free(unsafe.Pointer(cFn))

	if fnObj == nil {
		return nil, fmt.Errorf("%w: function '%s'", ErrNotFound, fn)
	}

	// Check if callable
	if C.PyCallable_Check(fnObj) == 0 {
		return nil, fmt.Errorf("'%s' is not callable", fn)
	}

	// Convert arguments to Python tuple
	pyArgs := C.PyTuple_New(C.Py_ssize_t(len(args)))
	defer C.Py_DecRef(pyArgs)

	for i, arg := range args {
		pyArg := ToPython(arg)
		C.PyTuple_SetItem(pyArgs, C.Py_ssize_t(i), pyArg)
	}

	// Call function
	result := C.PyObject_CallObject(fnObj, pyArgs)
	if result == nil {
		return nil, fmt.Errorf("%w: %s", ErrCallFailed, GetError())
	}
	defer C.Py_DecRef(result)

	return FromPython(result), nil
}

// Shutdown cleans up the state
func (s *State) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.shutdown {
		return
	}

	s.shutdown = true

	gil := AcquireGIL()
	defer gil.Release()

	if s.globals != nil {
		C.Py_DecRef(s.globals)
		s.globals = nil
	}

	if s.locals != nil {
		C.Py_DecRef(s.locals)
		s.locals = nil
	}
}

// ID returns the state identifier
func (s *State) ID() int {
	return s.id
}

// IsBusy returns whether state is currently executing
func (s *State) IsBusy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.busy
}
