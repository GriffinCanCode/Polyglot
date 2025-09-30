//go:build runtime_python
// +build runtime_python

package python

import (
	"fmt"
	"sync"
)

// Pool manages Python execution states
type Pool struct {
	states chan *State
	all    []*State
	size   int
	mu     sync.RWMutex
	closed bool
}

// NewPool creates a state pool
func NewPool(size int) *Pool {
	if size <= 0 {
		size = 4
	}
	return &Pool{
		states: make(chan *State, size),
		all:    make([]*State, 0, size),
		size:   size,
		closed: false,
	}
}

// Initialize creates states
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
	p.states = make(chan *State, size)
	p.all = make([]*State, 0, size)

	for i := 0; i < size; i++ {
		state := NewState(i)
		if err := state.Initialize(); err != nil {
			// Clean up already created states
			for _, s := range p.all {
				s.Shutdown()
			}
			return fmt.Errorf("failed to initialize state %d: %w", i, err)
		}
		p.all = append(p.all, state)
		p.states <- state
	}

	return nil
}

// Acquire gets a state from the pool (blocks until available)
func (p *Pool) Acquire() *State {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil
	}
	p.mu.RUnlock()

	return <-p.states
}

// Release returns a state to the pool
func (p *Pool) Release(state *State) {
	if state == nil {
		return
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.closed {
		p.states <- state
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

	// Don't close the channel - just mark as closed
	// Drain all states from channel
	for len(p.states) > 0 {
		<-p.states
	}

	// Shutdown all states
	for _, state := range p.all {
		state.Shutdown()
	}

	p.all = nil
}
