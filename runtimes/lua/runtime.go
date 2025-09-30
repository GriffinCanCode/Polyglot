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
	"context"
	"fmt"
	"sync"
	"unsafe"

	"github.com/griffincancode/polyglot.js/core"
)

// Runtime implements Lua runtime integration
type Runtime struct {
	config   core.RuntimeConfig
	pool     *Pool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a Lua runtime instance
func NewRuntime() *Runtime {
	return &Runtime{
		pool: NewPool(10),
	}
}

// Initialize prepares the Lua runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Initialize the pool
	if err := r.pool.Initialize(config.MaxConcurrency); err != nil {
		return fmt.Errorf("failed to initialize pool: %w", err)
	}

	return nil
}

// Execute runs Lua code
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	// Execute with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := worker.Execute(code, args...)
		resultChan <- result{value: res, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.value, res.err
	}
}

// Call invokes a Lua function
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.shutdown {
		r.mu.RUnlock()
		return nil, fmt.Errorf("runtime is shutdown")
	}
	r.mu.RUnlock()

	worker := r.pool.Acquire()
	defer r.pool.Release(worker)

	// Call with context cancellation support
	resultChan := make(chan result, 1)
	go func() {
		res, err := worker.Call(fn, args...)
		resultChan <- result{value: res, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.value, res.err
	}
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true

	if r.pool != nil {
		r.pool.Close()
	}

	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "lua"
}

// Version returns the Lua version
func (r *Runtime) Version() string {
	return "Lua 5.4"
}

type result struct {
	value interface{}
	err   error
}

// pushToLua pushes a Go value onto the Lua stack
func pushToLua(L *C.lua_State, val interface{}) {
	if val == nil {
		C.lua_pushnil(L)
		return
	}

	switch v := val.(type) {
	case string:
		cStr := C.CString(v)
		defer C.free(unsafe.Pointer(cStr))
		C.lua_pushstring(L, cStr)
	case int:
		C.lua_pushinteger(L, C.lua_Integer(v))
	case int64:
		C.lua_pushinteger(L, C.lua_Integer(v))
	case float64:
		C.lua_pushnumber(L, C.lua_Number(v))
	case bool:
		if v {
			C.lua_pushboolean(L, 1)
		} else {
			C.lua_pushboolean(L, 0)
		}
	default:
		C.lua_pushnil(L)
	}
}

// popFromLua pops a value from the Lua stack and converts to Go
func popFromLua(L *C.lua_State, idx C.int) interface{} {
	luaType := C.lua_type(L, idx)

	switch luaType {
	case C.LUA_TNIL:
		return nil
	case C.LUA_TBOOLEAN:
		return C.lua_toboolean(L, idx) != 0
	case C.LUA_TNUMBER:
		return float64(C.luawrap_tonumber(L, idx))
	case C.LUA_TSTRING:
		return C.GoString(C.luawrap_tostring(L, idx))
	default:
		return nil
	}
}
