//go:build runtime_lua
// +build runtime_lua

package lua

/*
#cgo CFLAGS: -I/opt/homebrew/include/lua
#cgo LDFLAGS: -L/opt/homebrew/lib -llua -lm
#include "luawrap.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// Worker represents a Lua state
type Worker struct {
	id       int
	state    *C.lua_State
	mu       sync.Mutex
	shutdown bool
}

// NewWorker creates a Lua worker
func NewWorker(id int) *Worker {
	return &Worker{
		id: id,
	}
}

// Initialize prepares the worker
func (w *Worker) Initialize() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return fmt.Errorf("worker is shutdown")
	}

	// Create new Lua state
	w.state = C.luaL_newstate()
	if w.state == nil {
		return fmt.Errorf("failed to create Lua state")
	}

	// Open standard libraries
	C.luaL_openlibs(w.state)

	return nil
}

// Execute runs Lua code
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	// Load and execute the code
	if C.luaL_loadstring(w.state, cCode) != 0 {
		err := C.GoString(C.luawrap_tostring(w.state, -1))
		C.luawrap_pop(w.state, 1)
		return nil, fmt.Errorf("lua load error: %s", err)
	}

	if C.luawrap_pcall(w.state, 0, 1, 0) != 0 {
		err := C.GoString(C.luawrap_tostring(w.state, -1))
		C.luawrap_pop(w.state, 1)
		return nil, fmt.Errorf("lua execution error: %s", err)
	}

	// Get result from stack
	result := popFromLua(w.state, -1)
	C.luawrap_pop(w.state, 1)

	return result, nil
}

// Call invokes a Lua function
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	cFn := C.CString(fn)
	defer C.free(unsafe.Pointer(cFn))

	// Get the function
	C.lua_getglobal(w.state, cFn)

	if C.luawrap_isfunction(w.state, -1) == 0 {
		C.luawrap_pop(w.state, 1)
		return nil, fmt.Errorf("function %s not found", fn)
	}

	// Push arguments
	for _, arg := range args {
		pushToLua(w.state, arg)
	}

	// Call the function
	nArgs := C.int(len(args))
	if C.luawrap_pcall(w.state, nArgs, 1, 0) != 0 {
		err := C.GoString(C.luawrap_tostring(w.state, -1))
		C.luawrap_pop(w.state, 1)
		return nil, fmt.Errorf("lua call error: %s", err)
	}

	// Get result
	result := popFromLua(w.state, -1)
	C.luawrap_pop(w.state, 1)

	return result, nil
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.shutdown && w.state != nil {
		C.lua_close(w.state)
		w.state = nil
	}

	w.shutdown = true
}
