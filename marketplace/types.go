package marketplace

import (
	"context"
	"time"
)

// Package represents a Polyglot package in the marketplace
type Package struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Description string            `json:"description"`
	License     string            `json:"license"`
	Homepage    string            `json:"homepage"`
	Repository  string            `json:"repository"`
	Languages   []string          `json:"languages"`
	Tags        []string          `json:"tags"`
	Downloads   int64             `json:"downloads"`
	Rating      float64           `json:"rating"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
	Checksum    string            `json:"checksum"`
}

// Template represents a project template
type Template struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Languages   []string          `json:"languages"`
	Category    string            `json:"category"`
	Tags        []string          `json:"tags"`
	Downloads   int64             `json:"downloads"`
	Rating      float64           `json:"rating"`
	Files       []TemplateFile    `json:"files"`
	Variables   map[string]string `json:"variables"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// TemplateFile represents a file in a template
type TemplateFile struct {
	Path       string `json:"path"`
	Content    string `json:"content"`
	Executable bool   `json:"executable"`
	Templated  bool   `json:"templated"`
}

// SearchQuery represents marketplace search parameters
type SearchQuery struct {
	Query     string   `json:"query"`
	Languages []string `json:"languages"`
	Tags      []string `json:"tags"`
	Author    string   `json:"author"`
	SortBy    string   `json:"sort_by"`
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
}

// SearchResult represents search results
type SearchResult struct {
	Packages  []Package  `json:"packages"`
	Templates []Template `json:"templates"`
	Total     int        `json:"total"`
	HasMore   bool       `json:"has_more"`
}

// Registry manages package registration and discovery
type Registry interface {
	// Search searches for packages and templates
	Search(ctx context.Context, query SearchQuery) (*SearchResult, error)

	// GetPackage retrieves a specific package
	GetPackage(ctx context.Context, id, version string) (*Package, error)

	// GetTemplate retrieves a specific template
	GetTemplate(ctx context.Context, id string) (*Template, error)

	// Publish publishes a new package
	Publish(ctx context.Context, pkg *Package, data []byte) error

	// PublishTemplate publishes a new template
	PublishTemplate(ctx context.Context, tmpl *Template) error

	// UpdatePackage updates an existing package
	UpdatePackage(ctx context.Context, pkg *Package) error

	// DeletePackage removes a package
	DeletePackage(ctx context.Context, id, version string) error
}

// Cache manages local package caching
type Cache interface {
	// Get retrieves a cached package
	Get(ctx context.Context, id, version string) ([]byte, error)

	// Put stores a package in cache
	Put(ctx context.Context, id, version string, data []byte) error

	// Has checks if a package is cached
	Has(ctx context.Context, id, version string) bool

	// Remove removes a package from cache
	Remove(ctx context.Context, id, version string) error

	// Clear clears all cached packages
	Clear(ctx context.Context) error

	// Size returns cache size in bytes
	Size(ctx context.Context) (int64, error)
}

// Validator validates packages for security and correctness
type Validator interface {
	// ValidatePackage validates a package
	ValidatePackage(ctx context.Context, pkg *Package, data []byte) error

	// ValidateTemplate validates a template
	ValidateTemplate(ctx context.Context, tmpl *Template) error

	// CheckSignature verifies package signature
	CheckSignature(ctx context.Context, pkg *Package, signature []byte) error

	// ScanSecurity performs security scanning
	ScanSecurity(ctx context.Context, data []byte) ([]SecurityIssue, error)
}

// SecurityIssue represents a security vulnerability
type SecurityIssue struct {
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Fix         string `json:"fix"`
}

// Client is the marketplace API client
type Client interface {
	// Registry returns the package registry
	Registry() Registry

	// Cache returns the local cache
	Cache() Cache

	// Validator returns the package validator
	Validator() Validator

	// Install installs a package
	Install(ctx context.Context, id, version string) error

	// Uninstall removes a package
	Uninstall(ctx context.Context, id string) error

	// Update updates a package
	Update(ctx context.Context, id string) error

	// List lists installed packages
	List(ctx context.Context) ([]Package, error)

	// InitFromTemplate initializes a project from template
	InitFromTemplate(ctx context.Context, tmplID, targetDir string, vars map[string]string) error
}
