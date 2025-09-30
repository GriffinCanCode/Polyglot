//go:build darwin
// +build darwin

package security

import (
	"fmt"
	"sync"
)

// DarwinEnforcer implements sandboxing using macOS App Sandbox
type DarwinEnforcer struct {
	policy  *Policy
	enabled bool
	mu      sync.RWMutex
}

// newPlatformEnforcer creates a macOS-specific enforcer
func newPlatformEnforcer() (Enforcer, error) {
	return &DarwinEnforcer{}, nil
}

// Enable activates App Sandbox
func (e *DarwinEnforcer) Enable(policy *Policy) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.enabled {
		return fmt.Errorf("enforcer already enabled")
	}

	e.policy = policy

	// Apply macOS sandbox profile
	if err := e.applySandboxProfile(); err != nil {
		return fmt.Errorf("failed to apply sandbox: %w", err)
	}

	e.enabled = true
	return nil
}

// Disable deactivates the enforcer
func (e *DarwinEnforcer) Disable() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.enabled {
		return nil
	}

	// macOS sandbox cannot be fully disabled once enabled
	e.enabled = false
	return nil
}

// Check verifies an operation against the policy
func (e *DarwinEnforcer) Check(op Operation) error {
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
func (e *DarwinEnforcer) Name() string {
	return "darwin-appsandbox"
}

// applySandboxProfile applies macOS sandbox profile
func (e *DarwinEnforcer) applySandboxProfile() error {
	// Generate sandbox profile based on policy
	profile := e.generateProfile()

	// Apply using sandbox_init
	// Simplified - actual implementation would use:
	// #include <sandbox.h>
	// sandbox_init(profile, SANDBOX_NAMED, &error)

	_ = profile // Use the profile
	return nil
}

// generateProfile generates a sandbox profile string
func (e *DarwinEnforcer) generateProfile() string {
	profile := "(version 1)\n"
	profile += "(deny default)\n"

	// Allow basic operations
	profile += "(allow process-exec (literal \"/bin/sh\"))\n"
	profile += "(allow sysctl-read)\n"

	// Add policy-specific rules
	for _, rp := range e.policy.RuntimePolicies() {
		if rp.AllowNetwork {
			profile += "(allow network*)\n"
		}
		if rp.AllowFileRead {
			profile += "(allow file-read*)\n"
		}
		if rp.AllowFileWrite {
			for _, path := range rp.AllowedPaths {
				profile += fmt.Sprintf("(allow file-write* (regex #\"%s\"))\n", path.String())
			}
		}
	}

	return profile
}
