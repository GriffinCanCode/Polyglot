//go:build runtime_python
// +build runtime_python

package python

// #include <Python.h>
import "C"

import (
	"errors"
	"unsafe"
)

// GILGuard manages GIL state for safe Python C API access
type GILGuard struct {
	state C.PyGILState_STATE
}

// AcquireGIL acquires the GIL for the current thread
func AcquireGIL() *GILGuard {
	state := C.PyGILState_Ensure()
	return &GILGuard{state: state}
}

// Release releases the GIL
func (g *GILGuard) Release() {
	C.PyGILState_Release(g.state)
}

// SafeDecRef safely decrements Python object reference count
func SafeDecRef(obj *C.PyObject) {
	if obj != nil {
		gil := AcquireGIL()
		defer gil.Release()
		C.Py_DecRef(obj)
	}
}

// SafeIncRef safely increments Python object reference count
func SafeIncRef(obj *C.PyObject) {
	if obj != nil {
		gil := AcquireGIL()
		defer gil.Release()
		C.Py_IncRef(obj)
	}
}

// GetError retrieves Python error if any
func GetError() string {
	if C.PyErr_Occurred() == nil {
		return ""
	}

	var ptype, pvalue, ptraceback *C.PyObject
	C.PyErr_Fetch(&ptype, &pvalue, &ptraceback)
	defer func() {
		if ptype != nil {
			C.Py_DecRef(ptype)
		}
		if pvalue != nil {
			C.Py_DecRef(pvalue)
		}
		if ptraceback != nil {
			C.Py_DecRef(ptraceback)
		}
	}()

	if pvalue == nil {
		return "unknown Python error"
	}

	pyStr := C.PyObject_Str(pvalue)
	if pyStr == nil {
		return "error converting Python error to string"
	}
	defer C.Py_DecRef(pyStr)

	cStr := C.PyUnicode_AsUTF8(pyStr)
	if cStr == nil {
		return "error getting UTF-8 from Python string"
	}

	return C.GoString(cStr)
}

// ClearError clears any Python error state
func ClearError() {
	C.PyErr_Clear()
}

// CheckError checks if Python error occurred and returns it
func CheckError() error {
	if C.PyErr_Occurred() == nil {
		return nil
	}
	errMsg := GetError()
	ClearError()
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return errors.New("unknown Python error")
}

// SafeString converts C string to Go string and frees it
func SafeString(cStr *C.char) string {
	if cStr == nil {
		return ""
	}
	goStr := C.GoString(cStr)
	C.free(unsafe.Pointer(cStr))
	return goStr
}
