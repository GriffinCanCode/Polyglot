package marketplace

import (
	"context"
	"fmt"
	"sync"
)

// DefaultClient implements the marketplace client
type DefaultClient struct {
	registry  Registry
	cache     Cache
	validator Validator
	mu        sync.RWMutex
	installed map[string]*Package
}

// NewClient creates a new marketplace client
func NewClient(registry Registry, cache Cache, validator Validator) *DefaultClient {
	return &DefaultClient{
		registry:  registry,
		cache:     cache,
		validator: validator,
		installed: make(map[string]*Package),
	}
}

// Registry returns the package registry
func (c *DefaultClient) Registry() Registry {
	return c.registry
}

// Cache returns the local cache
func (c *DefaultClient) Cache() Cache {
	return c.cache
}

// Validator returns the package validator
func (c *DefaultClient) Validator() Validator {
	return c.validator
}

// Install installs a package
func (c *DefaultClient) Install(ctx context.Context, id, version string) error {
	// Check if already installed
	c.mu.RLock()
	if pkg, ok := c.installed[id]; ok && pkg.Version == version {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	// Check cache first
	data, err := c.cache.Get(ctx, id, version)
	if err != nil || data == nil {
		// Fetch from registry
		pkg, err := c.registry.GetPackage(ctx, id, version)
		if err != nil {
			return fmt.Errorf("fetch package: %w", err)
		}

		// Download package data (simplified - would fetch actual binary)
		data = []byte(fmt.Sprintf("package-%s-%s", id, version))

		// Validate before installing
		if err := c.validator.ValidatePackage(ctx, pkg, data); err != nil {
			return fmt.Errorf("validate package: %w", err)
		}

		// Cache for future use
		if err := c.cache.Put(ctx, id, version, data); err != nil {
			return fmt.Errorf("cache package: %w", err)
		}
	}

	// Fetch metadata
	pkg, err := c.registry.GetPackage(ctx, id, version)
	if err != nil {
		return fmt.Errorf("fetch metadata: %w", err)
	}

	// Mark as installed
	c.mu.Lock()
	c.installed[id] = pkg
	c.mu.Unlock()

	return nil
}

// Uninstall removes a package
func (c *DefaultClient) Uninstall(ctx context.Context, id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.installed[id]; !ok {
		return fmt.Errorf("package not installed: %s", id)
	}

	delete(c.installed, id)
	return nil
}

// Update updates a package to the latest version
func (c *DefaultClient) Update(ctx context.Context, id string) error {
	c.mu.RLock()
	current, ok := c.installed[id]
	c.mu.RUnlock()

	if !ok {
		return fmt.Errorf("package not installed: %s", id)
	}

	// Search for latest version
	result, err := c.registry.Search(ctx, SearchQuery{
		Query: id,
		Limit: 1,
	})
	if err != nil {
		return fmt.Errorf("search for updates: %w", err)
	}

	if len(result.Packages) == 0 {
		return fmt.Errorf("package not found: %s", id)
	}

	latest := result.Packages[0]
	if latest.Version == current.Version {
		return nil // Already up to date
	}

	return c.Install(ctx, id, latest.Version)
}

// List lists installed packages
func (c *DefaultClient) List(ctx context.Context) ([]Package, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	packages := make([]Package, 0, len(c.installed))
	for _, pkg := range c.installed {
		packages = append(packages, *pkg)
	}

	return packages, nil
}

// InitFromTemplate initializes a project from a template
func (c *DefaultClient) InitFromTemplate(ctx context.Context, tmplID, targetDir string, vars map[string]string) error {
	// Fetch template
	tmpl, err := c.registry.GetTemplate(ctx, tmplID)
	if err != nil {
		return fmt.Errorf("fetch template: %w", err)
	}

	// Validate template
	if err := c.validator.ValidateTemplate(ctx, tmpl); err != nil {
		return fmt.Errorf("validate template: %w", err)
	}

	// Template expansion would happen here
	// For now, just verify the template is valid
	if len(tmpl.Files) == 0 {
		return fmt.Errorf("template has no files")
	}

	return nil
}
