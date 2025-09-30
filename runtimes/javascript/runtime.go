package javascript

import (
	"context"
	"fmt"
	"sync"

	"github.com/griffincancode/polyglot.js/core"
	"rogchap.com/v8go"
)

// Runtime implements JavaScript runtime integration using V8
type Runtime struct {
	config   core.RuntimeConfig
	isolate  *v8go.Isolate
	contexts *ContextPool
	mu       sync.RWMutex
	shutdown bool
}

// NewRuntime creates a JavaScript runtime instance
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Initialize prepares the JavaScript runtime
func (r *Runtime) Initialize(ctx context.Context, config core.RuntimeConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return fmt.Errorf("runtime is shutdown")
	}

	r.config = config

	// Create V8 isolate
	isolate := v8go.NewIsolate()
	if isolate == nil {
		return fmt.Errorf("failed to create V8 isolate")
	}
	r.isolate = isolate

	// Initialize context pool
	r.contexts = NewContextPool(config.MaxConcurrency, isolate)
	if err := r.contexts.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize context pool: %w", err)
	}

	return nil
}

// Execute runs JavaScript code
func (r *Runtime) Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	jsCtx := r.contexts.Acquire()
	defer r.contexts.Release(jsCtx)

	// Execute code
	val, err := jsCtx.RunScript(code, "execute.js")
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	return convertFromV8(val), nil
}

// Call invokes a JavaScript function
func (r *Runtime) Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.shutdown {
		return nil, fmt.Errorf("runtime is shutdown")
	}

	jsCtx := r.contexts.Acquire()
	defer r.contexts.Release(jsCtx)

	// Get function
	global := jsCtx.Global()
	fnVal, err := global.Get(fn)
	if err != nil {
		return nil, fmt.Errorf("function %s not found: %w", fn, err)
	}

	if !fnVal.IsFunction() {
		return nil, fmt.Errorf("%s is not a function", fn)
	}

	// Convert arguments to Valuers
	v8Args := make([]v8go.Valuer, len(args))
	for i, arg := range args {
		v8Args[i] = convertToV8(jsCtx, arg)
	}

	// Call function
	fnObj, _ := fnVal.AsFunction()
	result, err := fnObj.Call(global, v8Args...)
	if err != nil {
		return nil, fmt.Errorf("call failed: %w", err)
	}

	return convertFromV8(result), nil
}

// Shutdown stops the runtime
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shutdown {
		return nil
	}

	r.shutdown = true

	if r.contexts != nil {
		r.contexts.Close()
	}

	if r.isolate != nil {
		r.isolate.Dispose()
	}

	return nil
}

// Name returns the runtime identifier
func (r *Runtime) Name() string {
	return "javascript"
}

// Version returns the V8 version
func (r *Runtime) Version() string {
	return v8go.Version()
}

// convertToV8 converts Go value to V8 value
func convertToV8(ctx *v8go.Context, val interface{}) *v8go.Value {
	if val == nil {
		return v8go.Null(ctx.Isolate())
	}

	switch v := val.(type) {
	case string:
		val, _ := v8go.NewValue(ctx.Isolate(), v)
		return val
	case int:
		val, _ := v8go.NewValue(ctx.Isolate(), int32(v))
		return val
	case int32:
		val, _ := v8go.NewValue(ctx.Isolate(), v)
		return val
	case int64:
		val, _ := v8go.NewValue(ctx.Isolate(), float64(v))
		return val
	case float64:
		val, _ := v8go.NewValue(ctx.Isolate(), v)
		return val
	case bool:
		val, _ := v8go.NewValue(ctx.Isolate(), v)
		return val
	default:
		return v8go.Null(ctx.Isolate())
	}
}

// convertFromV8 converts V8 value to Go value
func convertFromV8(val *v8go.Value) interface{} {
	if val == nil || val.IsNull() || val.IsUndefined() {
		return nil
	}

	if val.IsBoolean() {
		return val.Boolean()
	}

	if val.IsNumber() {
		if val.IsInt32() {
			return val.Int32()
		}
		return val.Number()
	}

	if val.IsString() {
		return val.String()
	}

	return nil
}
