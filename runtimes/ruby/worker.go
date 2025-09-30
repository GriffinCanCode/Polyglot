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

	// Build the method call with proper argument formatting
	var code string
	if len(args) == 0 {
		code = fmt.Sprintf("%s()", fn)
	} else {
		// Convert arguments to Ruby literal syntax
		argStrings := make([]string, len(args))
		for i, arg := range args {
			argStrings[i] = formatRubyArgument(arg)
		}
		argList := ""
		for i, argStr := range argStrings {
			if i > 0 {
				argList += ", "
			}
			argList += argStr
		}
		code = fmt.Sprintf("%s(%s)", fn, argList)
	}

	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	var state C.int
	result := C.rb_eval_string_protect(cCode, &state)

	if state != 0 {
		// Exception occurred
		errVal := C.rb_errinfo()
		errMsg := C.rb_obj_as_string(errVal)
		return nil, fmt.Errorf("ruby call error: %s", rubyStringToGo(errMsg))
	}

	return convertFromRuby(result), nil
}

// formatRubyArgument converts a Go value to Ruby literal syntax
func formatRubyArgument(arg interface{}) string {
	if arg == nil {
		return "nil"
	}

	switch v := arg.(type) {
	case string:
		// Escape quotes in strings
		escaped := ""
		for _, ch := range v {
			if ch == '"' {
				escaped += "\\\""
			} else if ch == '\\' {
				escaped += "\\\\"
			} else {
				escaped += string(ch)
			}
		}
		return fmt.Sprintf("\"%s\"", escaped)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case float32:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		// For other types, attempt string conversion
		return fmt.Sprintf("\"%v\"", v)
	}
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.shutdown = true
}
