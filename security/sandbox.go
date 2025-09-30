package security

import (
	"context"
	"fmt"
	"sync"
)

// Sandbox manages runtime isolation and security policies
type Sandbox struct {
	policy   *Policy
	enforcer Enforcer
	mu       sync.RWMutex
	active   bool
}

// Enforcer implements platform-specific sandboxing
type Enforcer interface {
	// Enable activates the sandbox with the given policy
	Enable(policy *Policy) error

	// Disable deactivates the sandbox
	Disable() error

	// Check verifies if an operation is allowed
	Check(op Operation) error

	// Name returns the enforcer identifier
	Name() string
}

// Operation represents a sandboxed operation
type Operation struct {
	Type     OperationType
	Target   string
	Runtime  string
	Metadata map[string]interface{}
}

// OperationType categorizes operations
type OperationType string

const (
	OpFileRead    OperationType = "file_read"
	OpFileWrite   OperationType = "file_write"
	OpNetConnect  OperationType = "net_connect"
	OpNetListen   OperationType = "net_listen"
	OpExec        OperationType = "exec"
	OpMemAlloc    OperationType = "mem_alloc"
	OpSyscall     OperationType = "syscall"
	OpRuntimeCall OperationType = "runtime_call"
)

// NewSandbox creates a sandbox with the given policy
func NewSandbox(policy *Policy) (*Sandbox, error) {
	enforcer, err := newPlatformEnforcer()
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	return &Sandbox{
		policy:   policy,
		enforcer: enforcer,
	}, nil
}

// Enable activates the sandbox
func (s *Sandbox) Enable(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.active {
		return fmt.Errorf("sandbox already active")
	}

	if err := s.enforcer.Enable(s.policy); err != nil {
		return fmt.Errorf("failed to enable enforcer: %w", err)
	}

	s.active = true
	return nil
}

// Disable deactivates the sandbox
func (s *Sandbox) Disable(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return nil
	}

	if err := s.enforcer.Disable(); err != nil {
		return fmt.Errorf("failed to disable enforcer: %w", err)
	}

	s.active = false
	return nil
}

// Check verifies if an operation is allowed
func (s *Sandbox) Check(op Operation) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.active {
		return nil // Sandbox not active, allow all
	}

	// Check policy first
	if !s.policy.Allow(op) {
		return fmt.Errorf("operation %s denied by policy", op.Type)
	}

	// Check with enforcer
	return s.enforcer.Check(op)
}

// UpdatePolicy updates the sandbox policy
func (s *Sandbox) UpdatePolicy(policy *Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wasActive := s.active

	// Disable if active
	if wasActive {
		if err := s.enforcer.Disable(); err != nil {
			return fmt.Errorf("failed to disable for update: %w", err)
		}
	}

	s.policy = policy

	// Re-enable if was active
	if wasActive {
		if err := s.enforcer.Enable(s.policy); err != nil {
			return fmt.Errorf("failed to re-enable after update: %w", err)
		}
	}

	return nil
}

// Policy returns the current policy
func (s *Sandbox) Policy() *Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.policy
}

// IsActive returns whether the sandbox is active
func (s *Sandbox) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.active
}

// EnforcerName returns the name of the active enforcer
func (s *Sandbox) EnforcerName() string {
	return s.enforcer.Name()
}
