package security

import (
	"fmt"
	"regexp"
	"sync"
)

// Policy defines security rules for runtime execution
type Policy struct {
	name         string
	rules        []*Rule
	runtimeRules map[string]*RuntimePolicy
	mu           sync.RWMutex
}

// Rule defines a security rule
type Rule struct {
	Operation OperationType
	Target    *regexp.Regexp
	Action    Action
	Priority  int
}

// Action defines what to do with an operation
type Action int

const (
	ActionAllow Action = iota
	ActionDeny
	ActionAudit
)

// RuntimePolicy defines per-runtime security settings
type RuntimePolicy struct {
	Runtime        string
	AllowNetwork   bool
	AllowFileRead  bool
	AllowFileWrite bool
	AllowExec      bool
	AllowedPaths   []*regexp.Regexp
	DeniedPaths    []*regexp.Regexp
	MaxMemory      int64
}

// NewPolicy creates a new security policy
func NewPolicy(name string) *Policy {
	return &Policy{
		name:         name,
		rules:        make([]*Rule, 0),
		runtimeRules: make(map[string]*RuntimePolicy),
	}
}

// DefaultPolicy creates a restrictive default policy
func DefaultPolicy() *Policy {
	policy := NewPolicy("default")

	// Deny all by default
	policy.AddRule(&Rule{
		Operation: OpExec,
		Action:    ActionDeny,
		Priority:  100,
	})

	policy.AddRule(&Rule{
		Operation: OpNetConnect,
		Action:    ActionAudit,
		Priority:  50,
	})

	return policy
}

// PermissivePolicy creates a permissive policy for development
func PermissivePolicy() *Policy {
	policy := NewPolicy("permissive")

	// Allow everything but audit
	policy.AddRule(&Rule{
		Operation: OpFileWrite,
		Action:    ActionAudit,
		Priority:  10,
	})

	policy.AddRule(&Rule{
		Operation: OpNetConnect,
		Action:    ActionAudit,
		Priority:  10,
	})

	return policy
}

// AddRule adds a security rule
func (p *Policy) AddRule(rule *Rule) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.rules = append(p.rules, rule)
}

// SetRuntimePolicy sets security policy for a specific runtime
func (p *Policy) SetRuntimePolicy(runtime string, rp *RuntimePolicy) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rp.Runtime = runtime
	p.runtimeRules[runtime] = rp
}

// Allow checks if an operation is allowed by the policy
func (p *Policy) Allow(op Operation) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Check runtime-specific policy first
	if rp, exists := p.runtimeRules[op.Runtime]; exists {
		if !p.checkRuntimePolicy(op, rp) {
			return false
		}
	}

	// Check general rules
	var bestMatch *Rule
	for _, rule := range p.rules {
		if rule.Operation != op.Type {
			continue
		}

		if rule.Target != nil && op.Target != "" {
			if !rule.Target.MatchString(op.Target) {
				continue
			}
		}

		if bestMatch == nil || rule.Priority > bestMatch.Priority {
			bestMatch = rule
		}
	}

	if bestMatch == nil {
		return true // No rule found, allow by default
	}

	return bestMatch.Action == ActionAllow || bestMatch.Action == ActionAudit
}

// checkRuntimePolicy checks runtime-specific policies
func (p *Policy) checkRuntimePolicy(op Operation, rp *RuntimePolicy) bool {
	switch op.Type {
	case OpNetConnect, OpNetListen:
		return rp.AllowNetwork
	case OpFileRead:
		return rp.AllowFileRead && p.checkPaths(op.Target, rp)
	case OpFileWrite:
		return rp.AllowFileWrite && p.checkPaths(op.Target, rp)
	case OpExec:
		return rp.AllowExec
	default:
		return true
	}
}

// checkPaths verifies if a path is allowed
func (p *Policy) checkPaths(target string, rp *RuntimePolicy) bool {
	// Check denied paths first
	for _, pattern := range rp.DeniedPaths {
		if pattern.MatchString(target) {
			return false
		}
	}

	// If no allowed paths specified, allow all (except denied)
	if len(rp.AllowedPaths) == 0 {
		return true
	}

	// Check allowed paths
	for _, pattern := range rp.AllowedPaths {
		if pattern.MatchString(target) {
			return true
		}
	}

	return false
}

// Name returns the policy name
func (p *Policy) Name() string {
	return p.name
}

// Rules returns all rules
func (p *Policy) Rules() []*Rule {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rules := make([]*Rule, len(p.rules))
	copy(rules, p.rules)
	return rules
}

// RuntimePolicies returns all runtime-specific policies
func (p *Policy) RuntimePolicies() map[string]*RuntimePolicy {
	p.mu.RLock()
	defer p.mu.RUnlock()

	policies := make(map[string]*RuntimePolicy)
	for k, v := range p.runtimeRules {
		policies[k] = v
	}
	return policies
}

// MustCompileRule creates a rule with a compiled regex pattern
func MustCompileRule(op OperationType, pattern string, action Action, priority int) *Rule {
	var target *regexp.Regexp
	if pattern != "" {
		target = regexp.MustCompile(pattern)
	}

	return &Rule{
		Operation: op,
		Target:    target,
		Action:    action,
		Priority:  priority,
	}
}

// String returns a string representation of the action
func (a Action) String() string {
	switch a {
	case ActionAllow:
		return "allow"
	case ActionDeny:
		return "deny"
	case ActionAudit:
		return "audit"
	default:
		return fmt.Sprintf("unknown(%d)", a)
	}
}
