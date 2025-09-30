package core

import (
	"context"
	"fmt"
	"sync"
)

// Orchestrator coordinates all language runtimes
type Orchestrator struct {
	config   *Config
	runtimes map[string]Runtime
	memory   *MemoryCoordinator
	bridge   Bridge
	mu       sync.RWMutex
	shutdown chan struct{}
}

// NewOrchestrator creates a new orchestrator instance
func NewOrchestrator(config *Config) (*Orchestrator, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &Orchestrator{
		config:   config,
		runtimes: make(map[string]Runtime),
		memory:   NewMemoryCoordinator(config.Memory),
		shutdown: make(chan struct{}),
	}, nil
}

// RegisterRuntime adds a runtime to the orchestrator
func (o *Orchestrator) RegisterRuntime(runtime Runtime) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	name := runtime.Name()
	if _, exists := o.runtimes[name]; exists {
		return fmt.Errorf("runtime %s already registered", name)
	}

	o.runtimes[name] = runtime
	return nil
}

// Initialize starts all enabled runtimes
func (o *Orchestrator) Initialize(ctx context.Context) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for name, cfg := range o.config.Languages {
		if !cfg.Enabled {
			continue
		}

		runtime, exists := o.runtimes[name]
		if !exists {
			return fmt.Errorf("runtime %s not registered", name)
		}

		if err := runtime.Initialize(ctx, *cfg); err != nil {
			return fmt.Errorf("failed to initialize %s: %w", name, err)
		}
	}

	return nil
}

// Execute runs code in a specific runtime
func (o *Orchestrator) Execute(ctx context.Context, runtime string, code string, args ...interface{}) (interface{}, error) {
	o.mu.RLock()
	rt, exists := o.runtimes[runtime]
	o.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("runtime %s not found", runtime)
	}

	return rt.Execute(ctx, code, args...)
}

// Call invokes a function in a specific runtime
func (o *Orchestrator) Call(ctx context.Context, runtime string, fn string, args ...interface{}) (interface{}, error) {
	o.mu.RLock()
	rt, exists := o.runtimes[runtime]
	o.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("runtime %s not found", runtime)
	}

	return rt.Call(ctx, fn, args...)
}

// Memory returns the memory coordinator
func (o *Orchestrator) Memory() *MemoryCoordinator {
	return o.memory
}

// SetBridge configures the webview bridge
func (o *Orchestrator) SetBridge(bridge Bridge) {
	o.bridge = bridge
}

// Shutdown gracefully stops all runtimes
func (o *Orchestrator) Shutdown(ctx context.Context) error {
	close(o.shutdown)

	o.mu.RLock()
	defer o.mu.RUnlock()

	var errs []error
	for name, runtime := range o.runtimes {
		if err := runtime.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// Runtimes returns a list of registered runtime names
func (o *Orchestrator) Runtimes() []string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	names := make([]string, 0, len(o.runtimes))
	for name := range o.runtimes {
		names = append(names, name)
	}
	return names
}
