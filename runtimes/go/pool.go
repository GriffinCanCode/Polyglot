//go:build runtime_go
// +build runtime_go

package goruntime

import (
	"fmt"
	"sync"
)

// InterpreterPool manages a pool of Go interpreters
type InterpreterPool struct {
	interpreters chan *Interpreter
	size         int
	mu           sync.Mutex
	closed       bool
}

// NewInterpreterPool creates a new interpreter pool
func NewInterpreterPool(size int) *InterpreterPool {
	if size <= 0 {
		size = 4
	}

	return &InterpreterPool{
		interpreters: make(chan *Interpreter, size),
		size:         size,
	}
}

// Initialize creates the interpreter pool
func (p *InterpreterPool) Initialize() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return fmt.Errorf("pool is closed")
	}

	// Create interpreters
	for i := 0; i < p.size; i++ {
		interpreter, err := NewInterpreter()
		if err != nil {
			// Clean up any created interpreters
			close(p.interpreters)
			return fmt.Errorf("failed to create interpreter %d: %w", i, err)
		}
		p.interpreters <- interpreter
	}

	return nil
}

// Acquire gets an interpreter from the pool
func (p *InterpreterPool) Acquire() *Interpreter {
	return <-p.interpreters
}

// Release returns an interpreter to the pool
func (p *InterpreterPool) Release(interpreter *Interpreter) {
	if p.closed {
		return
	}

	select {
	case p.interpreters <- interpreter:
	default:
		// Pool is full or closed, discard the interpreter
	}
}

// Close shuts down the pool
func (p *InterpreterPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.interpreters)

	// Drain the channel
	for range p.interpreters {
		// Interpreters will be garbage collected
	}
}
