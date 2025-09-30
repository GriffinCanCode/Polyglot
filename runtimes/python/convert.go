//go:build runtime_python
// +build runtime_python

package python

// #include <Python.h>
// #include <stdlib.h>
//
// // Helper functions for type checking
// static int py_is_bool(PyObject *obj) {
//     return PyBool_Check(obj);
// }
// static int py_is_long(PyObject *obj) {
//     return PyLong_Check(obj);
// }
// static int py_is_float(PyObject *obj) {
//     return PyFloat_Check(obj);
// }
// static int py_is_unicode(PyObject *obj) {
//     return PyUnicode_Check(obj);
// }
// static int py_is_true(PyObject *obj) {
//     return obj == Py_True;
// }
// static int py_is_list(PyObject *obj) {
//     return PyList_Check(obj);
// }
// static int py_is_dict(PyObject *obj) {
//     return PyDict_Check(obj);
// }
// static int py_is_tuple(PyObject *obj) {
//     return PyTuple_Check(obj);
// }
import "C"

import "unsafe"

// ToPython converts Go value to Python object (caller must hold GIL)
func ToPython(val interface{}) *C.PyObject {
	if val == nil {
		C.Py_IncRef(C.Py_None)
		return C.Py_None
	}

	switch v := val.(type) {
	case string:
		return stringToPy(v)
	case int:
		return C.PyLong_FromLong(C.long(v))
	case int64:
		return C.PyLong_FromLongLong(C.longlong(v))
	case float64:
		return C.PyFloat_FromDouble(C.double(v))
	case bool:
		if v {
			C.Py_IncRef(C.Py_True)
			return C.Py_True
		}
		C.Py_IncRef(C.Py_False)
		return C.Py_False
	case []interface{}:
		return sliceToPy(v)
	case map[string]interface{}:
		return mapToPy(v)
	default:
		C.Py_IncRef(C.Py_None)
		return C.Py_None
	}
}

// FromPython converts Python object to Go value (caller must hold GIL)
func FromPython(obj *C.PyObject) interface{} {
	if obj == nil || obj == C.Py_None {
		return nil
	}

	// Check bool first (before long, since bools are longs in Python 3)
	if C.py_is_bool(obj) != 0 {
		return C.py_is_true(obj) != 0
	}

	// Check numeric types
	if C.py_is_long(obj) != 0 {
		return int64(C.PyLong_AsLongLong(obj))
	}

	if C.py_is_float(obj) != 0 {
		return float64(C.PyFloat_AsDouble(obj))
	}

	// Check string
	if C.py_is_unicode(obj) != 0 {
		return pyToString(obj)
	}

	// Check list
	if C.py_is_list(obj) != 0 {
		return pyToSlice(obj)
	}

	// Check dict
	if C.py_is_dict(obj) != 0 {
		return pyToMap(obj)
	}

	// Check tuple
	if C.py_is_tuple(obj) != 0 {
		return pyToSlice(obj)
	}

	// Fallback: try to convert to string representation
	return nil
}

// stringToPy converts Go string to Python string
func stringToPy(s string) *C.PyObject {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	return C.PyUnicode_FromString(cStr)
}

// pyToString converts Python string to Go string
func pyToString(pyStr *C.PyObject) string {
	if pyStr == nil {
		return ""
	}
	cStr := C.PyUnicode_AsUTF8(pyStr)
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}

// sliceToPy converts Go slice to Python list
func sliceToPy(slice []interface{}) *C.PyObject {
	pyList := C.PyList_New(C.Py_ssize_t(len(slice)))
	if pyList == nil {
		return C.Py_None
	}

	for i, item := range slice {
		pyItem := ToPython(item)
		C.PyList_SetItem(pyList, C.Py_ssize_t(i), pyItem)
	}

	return pyList
}

// pyToSlice converts Python list or tuple to Go slice
func pyToSlice(pyObj *C.PyObject) []interface{} {
	var size C.Py_ssize_t

	// Check if it's a list or tuple and get appropriate size
	if C.py_is_list(pyObj) != 0 {
		size = C.PyList_Size(pyObj)
	} else if C.py_is_tuple(pyObj) != 0 {
		size = C.PyTuple_Size(pyObj)
	} else {
		return []interface{}{}
	}

	// Safety check for size
	if size < 0 {
		return []interface{}{}
	}

	result := make([]interface{}, int(size))
	for i := 0; i < int(size); i++ {
		var item *C.PyObject
		if C.py_is_list(pyObj) != 0 {
			item = C.PyList_GetItem(pyObj, C.Py_ssize_t(i))
		} else {
			item = C.PyTuple_GetItem(pyObj, C.Py_ssize_t(i))
		}
		result[i] = FromPython(item)
	}
	return result
}

// mapToPy converts Go map to Python dict
func mapToPy(m map[string]interface{}) *C.PyObject {
	pyDict := C.PyDict_New()
	if pyDict == nil {
		return C.Py_None
	}

	for k, v := range m {
		pyKey := stringToPy(k)
		pyVal := ToPython(v)
		C.PyDict_SetItem(pyDict, pyKey, pyVal)
		C.Py_DecRef(pyKey)
		C.Py_DecRef(pyVal)
	}

	return pyDict
}

// pyToMap converts Python dict to Go map
func pyToMap(pyDict *C.PyObject) map[string]interface{} {
	result := make(map[string]interface{})

	var pos C.Py_ssize_t
	var key, value *C.PyObject

	for C.PyDict_Next(pyDict, &pos, &key, &value) != 0 {
		if C.py_is_unicode(key) != 0 {
			goKey := pyToString(key)
			result[goKey] = FromPython(value)
		}
	}

	return result
}
