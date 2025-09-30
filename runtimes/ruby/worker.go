//go:build runtime_ruby
// +build runtime_ruby

package ruby

/*
#include <ruby.h>
#include <stdlib.h>

// Helper function to protect rb_eval_string
static VALUE protected_eval(VALUE code_str) {
    return rb_eval_string(StringValueCStr(code_str));
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// Worker represents a Ruby execution context
type Worker struct {
	id       int
	mu       sync.Mutex
	shutdown bool
}

// NewWorker creates a Ruby worker
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

// Execute runs Ruby code
func (w *Worker) Execute(code string, args ...interface{}) (interface{}, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shutdown {
		return nil, fmt.Errorf("worker is shutdown")
	}

	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	var state C.int
	result := C.rb_eval_string_protect(cCode, &state)

	if state != 0 {
		// Exception occurred
		errVal := C.rb_errinfo()
		errMsg := C.rb_obj_as_string(errVal)
		return nil, fmt.Errorf("ruby error: %s", rubyStringToGo(errMsg))
	}

	return convertFromRuby(result), nil
}

// Call invokes a Ruby method
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
