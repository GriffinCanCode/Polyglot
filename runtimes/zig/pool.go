//go:build runtime_zig
// +build runtime_zig

package zig

import (
	"fmt"
	"sync"
)

// Pool manages Zig worker instances
type Pool struct {
	workers chan *Worker
	all     []*Worker
	size    int
	mu      sync.RWMutex
	closed  bool
}

// NewPool creates a worker pool
func NewPool(size int) *Pool {
	if size <= 0 {
		size = 4
	}
	return &Pool{
		workers: make(chan *Worker, size),
		all:     make([]*Worker, 0, size),
		size:    size,
		closed:  false,
	}
}

// Initialize creates workers
func (p *Pool) Initialize(size int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return fmt.Errorf("pool is closed")
	}

	if size <= 0 {
		size = 4
	}

	p.size = size
	p.workers = make(chan *Worker, size)
	p.all = make([]*Worker, 0, size)

	for i := 0; i < size; i++ {
		worker := NewWorker(i)
		if err := worker.Initialize(); err != nil {
			// Clean up already created workers
			for _, w := range p.all {
				w.Shutdown()
			}
			return fmt.Errorf("failed to initialize worker %d: %w", i, err)
		}
		p.all = append(p.all, worker)
		p.workers <- worker
	}

	return nil
}

// Acquire gets a worker from the pool (blocks until available)
func (p *Pool) Acquire() *Worker {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil
	}
	p.mu.RUnlock()

	return <-p.workers
}

// Release returns a worker to the pool
func (p *Pool) Release(worker *Worker) {
	if worker == nil {
		return
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.closed {
		p.workers <- worker
	}
}

// Size returns the pool size
func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.size
}

// Close shuts down the pool
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	// Drain all workers from channel
	for len(p.workers) > 0 {
		<-p.workers
	}

	// Shutdown all workers
	for _, worker := range p.all {
		worker.Shutdown()
	}

	p.all = nil
}
