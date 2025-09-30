//go:build runtime_php
// +build runtime_php

package php

/*
#include <sapi/embed/php_embed.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// Worker represents a PHP execution context
type Worker struct {
	id       int
	mu       sync.Mutex
	shutdown bool
}

// NewWorker creates a PHP worker
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

	// Worker-specific initialization if needed
	return nil
}

// Execute runs PHP code
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	// Execute PHP code
	result := C.zend_eval_string(cCode, nil, (*C.char)(unsafe.Pointer(C.CString("polyglot"))))
	if result == nil {
		return nil, fmt.Errorf("execution failed")
	}

	return nil, nil
}

// Call invokes a PHP function
func (w *Worker) Call(fn string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	// Build function call string
	code := fmt.Sprintf("%s()", fn)
	return w.Execute(code, args...)
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.shutdown = true
}
