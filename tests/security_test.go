package tests

import (
	"context"
	"regexp"
	"testing"

	"github.com/griffincancode/polyglot.js/security"
)

// TestPolicyCreation tests security policy creation
func TestPolicyCreation(t *testing.T) {
	policy := security.NewPolicy("test")

	if policy.Name() != "test" {
		t.Errorf("expected name 'test', got %s", policy.Name())
	}

	if len(policy.Rules()) != 0 {
		t.Errorf("expected 0 rules, got %d", len(policy.Rules()))
	}
}

// TestDefaultPolicy tests default policy behavior
func TestDefaultPolicy(t *testing.T) {
	policy := security.DefaultPolicy()

	if policy.Name() != "default" {
		t.Errorf("expected name 'default', got %s", policy.Name())
	}

	// Should deny exec by default
	op := security.Operation{
		Type:    security.OpExec,
		Target:  "/bin/sh",
		Runtime: "test",
	}

	if policy.Allow(op) {
		t.Error("default policy should deny exec operations")
	}
}

// TestPermissivePolicy tests permissive policy
func TestPermissivePolicy(t *testing.T) {
	policy := security.PermissivePolicy()

	if policy.Name() != "permissive" {
		t.Errorf("expected name 'permissive', got %s", policy.Name())
	}

	// Should allow most operations
	op := security.Operation{
		Type:    security.OpFileRead,
		Target:  "/tmp/test.txt",
		Runtime: "test",
	}

	if !policy.Allow(op) {
		t.Error("permissive policy should allow file read")
	}
}

// TestRuntimePolicy tests runtime-specific policies
func TestRuntimePolicy(t *testing.T) {
	policy := security.NewPolicy("runtime-test")

	rp := &security.RuntimePolicy{
		Runtime:        "python",
		AllowNetwork:   false,
		AllowFileRead:  true,
		AllowFileWrite: false,
		AllowExec:      false,
		MaxMemory:      1024 * 1024 * 100, // 100MB
	}

	policy.SetRuntimePolicy("python", rp)

	// Test network deny
	netOp := security.Operation{
		Type:    security.OpNetConnect,
		Target:  "example.com:443",
		Runtime: "python",
	}

	if policy.Allow(netOp) {
		t.Error("should deny network operations for python")
	}

	// Test file read allow
	readOp := security.Operation{
		Type:    security.OpFileRead,
		Target:  "/tmp/test.txt",
		Runtime: "python",
	}

	if !policy.Allow(readOp) {
		t.Error("should allow file read for python")
	}
}

// TestRuleMatching tests pattern-based rule matching
func TestRuleMatching(t *testing.T) {
	policy := security.NewPolicy("pattern-test")

	// Allow reads from /tmp
	rule := security.MustCompileRule(
		security.OpFileRead,
		"^/tmp/.*",
		security.ActionAllow,
		100,
	)
	policy.AddRule(rule)

	// Test matching path
	op1 := security.Operation{
		Type:   security.OpFileRead,
		Target: "/tmp/test.txt",
	}

	if !policy.Allow(op1) {
		t.Error("should allow /tmp/test.txt")
	}

	// Test non-matching path (but should still allow by default)
	op2 := security.Operation{
		Type:   security.OpFileRead,
		Target: "/var/log/test.log",
	}

	// No deny rule, so should allow
	if !policy.Allow(op2) {
		t.Error("should allow by default when no deny rule")
	}
}

// TestPathRestrictions tests allowed/denied path patterns
func TestPathRestrictions(t *testing.T) {
	policy := security.NewPolicy("path-test")

	rp := &security.RuntimePolicy{
		Runtime:       "rust",
		AllowFileRead: true,
		AllowedPaths: []*regexp.Regexp{
			regexp.MustCompile("^/home/user/.*"),
		},
		DeniedPaths: []*regexp.Regexp{
			regexp.MustCompile("^/home/user/secrets/.*"),
		},
	}

	policy.SetRuntimePolicy("rust", rp)

	// Test allowed path
	op1 := security.Operation{
		Type:    security.OpFileRead,
		Target:  "/home/user/data.txt",
		Runtime: "rust",
	}

	if !policy.Allow(op1) {
		t.Error("should allow /home/user/data.txt")
	}

	// Test denied path
	op2 := security.Operation{
		Type:    security.OpFileRead,
		Target:  "/home/user/secrets/key.txt",
		Runtime: "rust",
	}

	if policy.Allow(op2) {
		t.Error("should deny /home/user/secrets/key.txt")
	}
}

// TestSandbox tests sandbox creation and lifecycle
func TestSandbox(t *testing.T) {
	policy := security.PermissivePolicy()
	sandbox, err := security.NewSandbox(policy)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}

	ctx := context.Background()

	// Enable sandbox
	err = sandbox.Enable(ctx)
	if err != nil {
		t.Logf("sandbox enable returned error (may be expected): %v", err)
	}

	if !sandbox.IsActive() && err == nil {
		t.Error("sandbox should be active after enable")
	}

	// Check operation
	op := security.Operation{
		Type:   security.OpFileRead,
		Target: "/tmp/test.txt",
	}

	err = sandbox.Check(op)
	if err != nil {
		t.Errorf("check failed: %v", err)
	}

	// Disable sandbox
	err = sandbox.Disable(ctx)
	if err != nil {
		t.Errorf("disable failed: %v", err)
	}
}

// TestSandboxPolicyUpdate tests updating sandbox policy
func TestSandboxPolicyUpdate(t *testing.T) {
	policy1 := security.PermissivePolicy()
	sandbox, err := security.NewSandbox(policy1)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}

	// Update to default policy
	policy2 := security.DefaultPolicy()
	err = sandbox.UpdatePolicy(policy2)
	if err != nil {
		t.Errorf("policy update failed: %v", err)
	}

	if sandbox.Policy().Name() != "default" {
		t.Error("policy should be updated to default")
	}
}

// TestEnforcerNames tests platform-specific enforcer names
func TestEnforcerNames(t *testing.T) {
	policy := security.NewPolicy("test")
	sandbox, err := security.NewSandbox(policy)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}

	name := sandbox.EnforcerName()
	if name == "" {
		t.Error("enforcer name should not be empty")
	}

	t.Logf("Platform enforcer: %s", name)
}
