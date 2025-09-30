//go:build runtime_wasm
// +build runtime_wasm

package wasm

import (
	"fmt"
	"sync"
)

// Pool manages WASM execution workers
type Pool struct {
	workers chan *Worker
	size    int
	mu      sync.Mutex
	closed  bool
}

// NewPool creates a worker pool
func NewPool(size int) *Pool {
	return &Pool{
		size: size,
	}
}

// Initialize creates workers
func (p *Pool) Initialize() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return fmt.Errorf("pool is closed")
	}

	p.workers = make(chan *Worker, p.size)

	for i := 0; i < p.size; i++ {
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
	select {
	case p.workers <- worker:
	default:
		// Pool is closed, just cleanup the worker
		worker.Shutdown()
	}
}

// Close shuts down the pool
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.workers)

	// Drain and shutdown all workers
	for worker := range p.workers {
		worker.Shutdown()
	}
}
