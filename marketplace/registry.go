package marketplace

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryRegistry implements an in-memory registry for testing
type MemoryRegistry struct {
	mu        sync.RWMutex
	packages  map[string]map[string]*Package // id -> version -> package
	templates map[string]*Template           // id -> template
	data      map[string][]byte              // packageID-version -> data
}

// NewMemoryRegistry creates a new in-memory registry
func NewMemoryRegistry() *MemoryRegistry {
	return &MemoryRegistry{
		packages:  make(map[string]map[string]*Package),
		templates: make(map[string]*Template),
		data:      make(map[string][]byte),
	}
}

// Search searches for packages and templates
func (r *MemoryRegistry) Search(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := &SearchResult{
		Packages:  make([]Package, 0),
		Templates: make([]Template, 0),
	}

	// Search packages
	for _, versions := range r.packages {
		for _, pkg := range versions {
			if matchesQuery(pkg, query) {
				result.Packages = append(result.Packages, *pkg)
			}
		}
	}

	// Search templates
	for _, tmpl := range r.templates {
		if matchesTemplateQuery(tmpl, query) {
			result.Templates = append(result.Templates, *tmpl)
		}
	}

	result.Total = len(result.Packages) + len(result.Templates)
	result.HasMore = false

	return result, nil
}

// GetPackage retrieves a specific package
func (r *MemoryRegistry) GetPackage(ctx context.Context, id, version string) (*Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, ok := r.packages[id]
	if !ok {
		return nil, fmt.Errorf("package not found: %s", id)
	}

	pkg, ok := versions[version]
	if !ok {
		return nil, fmt.Errorf("version not found: %s@%s", id, version)
	}

	return pkg, nil
}

// GetTemplate retrieves a specific template
func (r *MemoryRegistry) GetTemplate(ctx context.Context, id string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tmpl, ok := r.templates[id]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", id)
	}

	return tmpl, nil
}

// Publish publishes a new package
func (r *MemoryRegistry) Publish(ctx context.Context, pkg *Package, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.packages[pkg.ID]; !ok {
		r.packages[pkg.ID] = make(map[string]*Package)
	}

	pkg.CreatedAt = time.Now()
	pkg.UpdatedAt = time.Now()
	r.packages[pkg.ID][pkg.Version] = pkg
	r.data[fmt.Sprintf("%s-%s", pkg.ID, pkg.Version)] = data

	return nil
}

// PublishTemplate publishes a new template
func (r *MemoryRegistry) PublishTemplate(ctx context.Context, tmpl *Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tmpl.CreatedAt = time.Now()
	tmpl.UpdatedAt = time.Now()
	r.templates[tmpl.ID] = tmpl

	return nil
}

// UpdatePackage updates an existing package
func (r *MemoryRegistry) UpdatePackage(ctx context.Context, pkg *Package) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	versions, ok := r.packages[pkg.ID]
	if !ok {
		return fmt.Errorf("package not found: %s", pkg.ID)
	}

	if _, ok := versions[pkg.Version]; !ok {
		return fmt.Errorf("version not found: %s@%s", pkg.ID, pkg.Version)
	}

	pkg.UpdatedAt = time.Now()
	versions[pkg.Version] = pkg

	return nil
}

// DeletePackage removes a package
func (r *MemoryRegistry) DeletePackage(ctx context.Context, id, version string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	versions, ok := r.packages[id]
	if !ok {
		return fmt.Errorf("package not found: %s", id)
	}

	delete(versions, version)
	delete(r.data, fmt.Sprintf("%s-%s", id, version))

	if len(versions) == 0 {
		delete(r.packages, id)
	}

	return nil
}

// Helper functions
func matchesQuery(pkg *Package, query SearchQuery) bool {
	if query.Query != "" && pkg.Name != query.Query && pkg.ID != query.Query {
		return false
	}
	if query.Author != "" && pkg.Author != query.Author {
		return false
	}
	if len(query.Languages) > 0 && !containsAny(pkg.Languages, query.Languages) {
		return false
	}
	if len(query.Tags) > 0 && !containsAny(pkg.Tags, query.Tags) {
		return false
	}
	return true
}

func matchesTemplateQuery(tmpl *Template, query SearchQuery) bool {
	if query.Query != "" && tmpl.Name != query.Query && tmpl.ID != query.Query {
		return false
	}
	if query.Author != "" && tmpl.Author != query.Author {
		return false
	}
	if len(query.Languages) > 0 && !containsAny(tmpl.Languages, query.Languages) {
		return false
	}
	if len(query.Tags) > 0 && !containsAny(tmpl.Tags, query.Tags) {
		return false
	}
	return true
}

func containsAny(haystack, needles []string) bool {
	for _, needle := range needles {
		for _, h := range haystack {
			if h == needle {
				return true
			}
		}
	}
	return false
}
