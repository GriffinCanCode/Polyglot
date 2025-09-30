package cloud

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryAuth implements an in-memory auth service for testing
type MemoryAuth struct {
	mu    sync.RWMutex
	creds map[string]*Credentials
}

// NewMemoryAuth creates a new in-memory auth service
func NewMemoryAuth() *MemoryAuth {
	return &MemoryAuth{
		creds: make(map[string]*Credentials),
	}
}

// Authenticate validates credentials
func (a *MemoryAuth) Authenticate(ctx context.Context, apiKey, secretKey string) (*Credentials, error) {
	if apiKey == "" || secretKey == "" {
		return nil, fmt.Errorf("API key and secret key are required")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if credentials already exist
	key := fmt.Sprintf("%s:%s", apiKey, secretKey)
	if creds, ok := a.creds[key]; ok {
		if time.Now().Before(creds.ExpiresAt) {
			return creds, nil
		}
	}

	// Create new credentials
	creds := &Credentials{
		APIKey:    apiKey,
		SecretKey: secretKey,
		ProjectID: fmt.Sprintf("project-%s", apiKey[:8]),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	a.creds[key] = creds
	return creds, nil
}

// Refresh refreshes expired credentials
func (a *MemoryAuth) Refresh(ctx context.Context, creds *Credentials) (*Credentials, error) {
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	key := fmt.Sprintf("%s:%s", creds.APIKey, creds.SecretKey)
	existing, ok := a.creds[key]
	if !ok {
		return nil, fmt.Errorf("credentials not found")
	}

	existing.ExpiresAt = time.Now().Add(24 * time.Hour)
	return existing, nil
}

// Validate validates credentials
func (a *MemoryAuth) Validate(ctx context.Context, creds *Credentials) error {
	if creds == nil {
		return fmt.Errorf("credentials are required")
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", creds.APIKey, creds.SecretKey)
	existing, ok := a.creds[key]
	if !ok {
		return fmt.Errorf("credentials not found")
	}

	if time.Now().After(existing.ExpiresAt) {
		return fmt.Errorf("credentials expired")
	}

	return nil
}

// Revoke revokes credentials
func (a *MemoryAuth) Revoke(ctx context.Context, creds *Credentials) error {
	if creds == nil {
		return fmt.Errorf("credentials are required")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	key := fmt.Sprintf("%s:%s", creds.APIKey, creds.SecretKey)
	delete(a.creds, key)

	return nil
}
