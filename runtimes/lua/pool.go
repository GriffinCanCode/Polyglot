//go:build runtime_lua
// +build runtime_lua

package lua

/*
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
*/
import "C"

import (
	"fmt"
	"sync"
)

// Pool manages Lua state workers
type Pool struct {
	workers chan *Worker
	size    int
	mu      sync.Mutex
}

// NewPool creates a worker pool
func NewPool(size int) *Pool {
	return &Pool{
		workers: make(chan *Worker, size),
		size:    size,
	}
}

// Initialize creates workers
func (p *Pool) Initialize(size int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.size = size
	p.workers = make(chan *Worker, size)

	for i := 0; i < size; i++ {
		worker := NewWorker(i)
		if err := worker.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize worker %d: %w", i, err)
		}
		p.workers <- worker
	}

	return nil
}

// Acquire gets a worker from the pool
func (p *Pool) Acquire() *Worker {
	return <-p.workers
}

// Release returns a worker to the pool
func (p *Pool) Release(worker *Worker) {
	p.workers <- worker
}

// Close shuts down the pool
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.workers)
	for worker := range p.workers {
		worker.Shutdown()
	}
}
