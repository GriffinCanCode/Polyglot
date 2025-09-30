//go:build runtime_cpp
// +build runtime_cpp

package cpp

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// Loader manages C++ library loading
type Loader struct {
	handle  unsafe.Pointer
	symbols map[string]unsafe.Pointer
	mu      sync.RWMutex
}

// NewLoader creates a library loader
func NewLoader() *Loader {
	return &Loader{
		symbols: make(map[string]unsafe.Pointer),
	}
}

// Load opens a shared library
func (l *Loader) Load(path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.handle != nil {
		return fmt.Errorf("library already loaded")
	}

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	handle := C.dlopen(cPath, C.RTLD_LAZY)
	if handle == nil {
		errStr := C.GoString(C.dlerror())
		return fmt.Errorf("dlopen failed: %s", errStr)
	}

	l.handle = handle
	return nil
}

// Symbol looks up a function symbol
func (l *Loader) Symbol(name string) (unsafe.Pointer, error) {
	l.mu.RLock()

	// Check cache
	if sym, exists := l.symbols[name]; exists {
		l.mu.RUnlock()
		return sym, nil
	}
	l.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.handle == nil {
		return nil, fmt.Errorf("no library loaded")
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	// Clear any existing error
	C.dlerror()

	symbol := C.dlsym(l.handle, cName)
	if symbol == nil {
		errStr := C.GoString(C.dlerror())
		return nil, fmt.Errorf("dlsym failed: %s", errStr)
	}

	l.symbols[name] = symbol
	return symbol, nil
}

// Invoke calls a C++ function
func (l *Loader) Invoke(name string, args ...interface{}) (interface{}, error) {
	symbol, err := l.Symbol(name)
	if err != nil {
		return nil, err
	}

	// Simplified invocation - actual implementation would handle type marshaling
	if len(args) == 0 {
		type fn0 func() int64
		f := *(*fn0)(unsafe.Pointer(&symbol))
		result := f()
		return result, nil
	}

	return nil, fmt.Errorf("complex argument handling not yet implemented")
}

// Unload closes the shared library
func (l *Loader) Unload() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.handle == nil {
		return nil
	}

	if C.dlclose(l.handle) != 0 {
		errStr := C.GoString(C.dlerror())
		return fmt.Errorf("dlclose failed: %s", errStr)
	}

	l.handle = nil
	l.symbols = make(map[string]unsafe.Pointer)
	return nil
}

// IsLoaded checks if a library is loaded
func (l *Loader) IsLoaded() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.handle != nil
}
