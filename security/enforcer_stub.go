//go:build !linux && !darwin && !windows
// +build !linux,!darwin,!windows

package security

import "fmt"

// StubEnforcer is a no-op enforcer for unsupported platforms
type StubEnforcer struct {
	policy *Policy
}

// newPlatformEnforcer creates a stub enforcer
func newPlatformEnforcer() (Enforcer, error) {
	return &StubEnforcer{}, nil
}

// Enable does nothing on unsupported platforms
func (e *StubEnforcer) Enable(policy *Policy) error {
	e.policy = policy
	fmt.Println("Warning: Sandbox not supported on this platform")
	return nil
}

// Disable does nothing
func (e *StubEnforcer) Disable() error {
	return nil
}

// Check allows all operations
func (e *StubEnforcer) Check(op Operation) error {
	return nil
}

// Name returns the enforcer identifier
func (e *StubEnforcer) Name() string {
	return "stub (unsupported platform)"
}
