package javascript

import (
	"fmt"
	"sync"

	"rogchap.com/v8go"
)

// ContextPool manages V8 contexts
type ContextPool struct {
	contexts chan *v8go.Context
	isolate  *v8go.Isolate
	size     int
	mu       sync.Mutex
}

// NewContextPool creates a context pool
func NewContextPool(size int, isolate *v8go.Isolate) *ContextPool {
	return &ContextPool{
		contexts: make(chan *v8go.Context, size),
		isolate:  isolate,
		size:     size,
	}
}

// Initialize creates contexts
func (p *ContextPool) Initialize() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := 0; i < p.size; i++ {
		ctx := v8go.NewContext(p.isolate)
		if ctx == nil {
			return fmt.Errorf("failed to create context %d", i)
		}
		p.contexts <- ctx
	}

	return nil
}

// Acquire gets a context from the pool
func (p *ContextPool) Acquire() *v8go.Context {
	return <-p.contexts
}

// Release returns a context to the pool
func (p *ContextPool) Release(ctx *v8go.Context) {
	p.contexts <- ctx
}

// Close shuts down the pool
func (p *ContextPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.contexts)
	for ctx := range p.contexts {
		ctx.Close()
	}
}
