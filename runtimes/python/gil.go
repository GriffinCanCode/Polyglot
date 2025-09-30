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
	C.PyErr_NormalizeException(&ptype, &pvalue, &ptraceback)

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

	// Build comprehensive error message with type and value
	errorMsg := ""

	// Get exception type name
	if ptype != nil {
		cTypeName := C.CString("__name__")
		defer C.free(unsafe.Pointer(cTypeName))
		typeName := C.PyObject_GetAttrString(ptype, cTypeName)
		if typeName != nil {
			defer C.Py_DecRef(typeName)
			typeStr := C.PyUnicode_AsUTF8(typeName)
			if typeStr != nil {
				errorMsg = C.GoString(typeStr) + ": "
			}
		}
	}

	// Get exception message
	pyStr := C.PyObject_Str(pvalue)
	if pyStr == nil {
		return errorMsg + "error converting Python error to string"
	}
	defer C.Py_DecRef(pyStr)

	cStr := C.PyUnicode_AsUTF8(pyStr)
	if cStr == nil {
		return errorMsg + "error getting UTF-8 from Python string"
	}

	errorMsg += C.GoString(cStr)

	// Add traceback if available
	if ptraceback != nil {
		tb := getTraceback(ptype, pvalue, ptraceback)
		if tb != "" {
			errorMsg += "\n\nTraceback:\n" + tb
		}
	}

	return errorMsg
}

// getTraceback extracts formatted traceback from Python exception
func getTraceback(ptype, pvalue, ptraceback *C.PyObject) string {
	if ptraceback == nil {
		return ""
	}

	// Import traceback module
	cTraceback := C.CString("traceback")
	defer C.free(unsafe.Pointer(cTraceback))
	tbModule := C.PyImport_ImportModule(cTraceback)
	if tbModule == nil {
		return ""
	}
	defer C.Py_DecRef(tbModule)

	// Get format_exception function
	cFormatException := C.CString("format_exception")
	defer C.free(unsafe.Pointer(cFormatException))
	formatFunc := C.PyObject_GetAttrString(tbModule, cFormatException)
	if formatFunc == nil {
		return ""
	}
	defer C.Py_DecRef(formatFunc)

	// Create argument tuple
	args := C.PyTuple_New(3)
	defer C.Py_DecRef(args)

	C.Py_IncRef(ptype)
	C.Py_IncRef(pvalue)
	C.Py_IncRef(ptraceback)

	C.PyTuple_SetItem(args, 0, ptype)
	C.PyTuple_SetItem(args, 1, pvalue)
	C.PyTuple_SetItem(args, 2, ptraceback)

	// Call format_exception
	result := C.PyObject_CallObject(formatFunc, args)
	if result == nil {
		return ""
	}
	defer C.Py_DecRef(result)

	// Join the list of strings
	if C.py_is_list(result) == 0 {
		return ""
	}

	size := C.PyList_Size(result)
	if size <= 0 {
		return ""
	}

	traceback := ""
	for i := C.Py_ssize_t(0); i < size; i++ {
		item := C.PyList_GetItem(result, i)
		if item != nil && C.py_is_unicode(item) != 0 {
			str := C.PyUnicode_AsUTF8(item)
			if str != nil {
				traceback += C.GoString(str)
			}
		}
	}

	return traceback
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
