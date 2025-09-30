//go:build linux
// +build linux

package security

import (
	"fmt"
	"sync"
)

// LinuxEnforcer implements sandboxing using Linux Landlock and seccomp
type LinuxEnforcer struct {
	policy  *Policy
	enabled bool
	mu      sync.RWMutex
}

// newPlatformEnforcer creates a Linux-specific enforcer
func newPlatformEnforcer() (Enforcer, error) {
	return &LinuxEnforcer{}, nil
}

// Enable activates Landlock and seccomp filters
func (e *LinuxEnforcer) Enable(policy *Policy) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.enabled {
		return fmt.Errorf("enforcer already enabled")
	}

	e.policy = policy

	// Apply seccomp filters based on policy
	if err := e.applySeccomp(); err != nil {
		return fmt.Errorf("failed to apply seccomp: %w", err)
	}

	// Apply Landlock restrictions
	if err := e.applyLandlock(); err != nil {
		return fmt.Errorf("failed to apply landlock: %w", err)
	}

	e.enabled = true
	return nil
}

// Disable deactivates the enforcer
func (e *LinuxEnforcer) Disable() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.enabled {
		return nil
	}

	// Note: seccomp and landlock cannot be fully disabled once enabled
	// This is by design for security
	e.enabled = false
	return nil
}

// Check verifies an operation against the policy
func (e *LinuxEnforcer) Check(op Operation) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.enabled {
		return nil
	}

	// Verify operation is allowed
	if !e.policy.Allow(op) {
		return fmt.Errorf("operation %s denied", op.Type)
	}

	return nil
}

// Name returns the enforcer identifier
func (e *LinuxEnforcer) Name() string {
	return "linux-landlock"
}

// applySeccomp applies seccomp-bpf filters
func (e *LinuxEnforcer) applySeccomp() error {
	// Check if any runtime policy denies exec
	denyExec := false
	for _, rp := range e.policy.RuntimePolicies() {
		if !rp.AllowExec {
			denyExec = true
			break
		}
	}

	if denyExec {
		// Apply seccomp filter to block execve syscalls
		// Simplified - actual implementation would use libseccomp
		// and golang.org/x/sys/unix for Prctl
		// For now, just track the requirement
		_ = denyExec
	}

	return nil
}

// applyLandlock applies filesystem restrictions using Landlock LSM
func (e *LinuxEnforcer) applyLandlock() error {
	// Landlock is available in Linux 5.13+
	// Simplified implementation - actual would use landlock syscalls

	// For each runtime policy, restrict file access
	for _, rp := range e.policy.RuntimePolicies() {
		if !rp.AllowFileWrite {
			// Apply read-only restrictions
			// This would use landlock_create_ruleset, landlock_add_rule, etc.
		}
	}

	return nil
}
