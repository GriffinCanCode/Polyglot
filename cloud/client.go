package cloud

import (
	"context"
	"fmt"
	"sync"
)

// DefaultClient implements the cloud services client
type DefaultClient struct {
	builder Builder
	storage Storage
	auth    Auth
	creds   *Credentials
	mu      sync.RWMutex
}

// NewClient creates a new cloud client
func NewClient(builder Builder, storage Storage, auth Auth) *DefaultClient {
	return &DefaultClient{
		builder: builder,
		storage: storage,
		auth:    auth,
	}
}

// Builder returns the build orchestrator
func (c *DefaultClient) Builder() Builder {
	return c.builder
}

// Storage returns artifact storage
func (c *DefaultClient) Storage() Storage {
	return c.storage
}

// Auth returns authentication service
func (c *DefaultClient) Auth() Auth {
	return c.auth
}

// Authenticate authenticates with the cloud service
func (c *DefaultClient) Authenticate(ctx context.Context, apiKey, secretKey string) error {
	creds, err := c.auth.Authenticate(ctx, apiKey, secretKey)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	c.mu.Lock()
	c.creds = creds
	c.mu.Unlock()

	return nil
}

// CrossCompile performs cross-platform compilation
func (c *DefaultClient) CrossCompile(ctx context.Context, source []byte, platforms []Platform) ([]*BuildResult, error) {
	c.mu.RLock()
	if c.creds == nil {
		c.mu.RUnlock()
		return nil, fmt.Errorf("not authenticated")
	}
	creds := c.creds
	c.mu.RUnlock()

	// Validate credentials
	if err := c.auth.Validate(ctx, creds); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	results := make([]*BuildResult, 0, len(platforms))
	errChan := make(chan error, len(platforms))
	resultChan := make(chan *BuildResult, len(platforms))

	// Build for each platform in parallel
	var wg sync.WaitGroup
	for _, platform := range platforms {
		wg.Add(1)
		go func(p Platform) {
			defer wg.Done()

			req := &BuildRequest{
				ProjectID: creds.ProjectID,
				Platform:  p,
				Source:    source,
			}

			result, err := c.builder.Build(ctx, req)
			if err != nil {
				errChan <- fmt.Errorf("build for %s/%s failed: %w", p.OS, p.Arch, err)
				return
			}

			resultChan <- result
		}(platform)
	}

	// Wait for all builds to complete
	go func() {
		wg.Wait()
		close(errChan)
		close(resultChan)
	}()

	// Collect results
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	for res := range resultChan {
		results = append(results, res)
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("some builds failed: %v", errs)
	}

	return results, nil
}
