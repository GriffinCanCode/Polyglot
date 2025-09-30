package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/marketplace"
)

func TestMarketplaceRegistry(t *testing.T) {
	ctx := context.Background()
	registry := marketplace.NewMemoryRegistry()

	// Test publishing a package
	pkg := &marketplace.Package{
		ID:          "test-package",
		Name:        "Test Package",
		Version:     "1.0.0",
		Author:      "Test Author",
		Description: "A test package",
		License:     "MIT",
		Languages:   []string{"python", "javascript"},
		Tags:        []string{"testing", "demo"},
		Checksum:    "abc123",
	}

	err := registry.Publish(ctx, pkg, []byte("package-data"))
	if err != nil {
		t.Fatalf("failed to publish package: %v", err)
	}

	// Test retrieving a package
	retrieved, err := registry.GetPackage(ctx, "test-package", "1.0.0")
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if retrieved.Name != pkg.Name {
		t.Errorf("expected name %s, got %s", pkg.Name, retrieved.Name)
	}

	// Test searching packages
	result, err := registry.Search(ctx, marketplace.SearchQuery{
		Query: "test-package",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}

	if len(result.Packages) == 0 {
		t.Error("expected at least one package in search results")
	}
}

func TestMarketplaceClient(t *testing.T) {
	ctx := context.Background()
	registry := marketplace.NewMemoryRegistry()
	cache := marketplace.NewMemoryCache()
	validator := marketplace.NewValidator()
	client := marketplace.NewClient(registry, cache, validator)

	// Publish a test package
	pkg := &marketplace.Package{
		ID:       "client-test",
		Name:     "Client Test",
		Version:  "1.0.0",
		Author:   "Test",
		Checksum: "checksum-1",
	}
	registry.Publish(ctx, pkg, []byte("test-data"))

	// Test installing a package
	err := client.Install(ctx, "client-test", "1.0.0")
	if err != nil {
		t.Fatalf("failed to install package: %v", err)
	}

	// Test listing installed packages
	installed, err := client.List(ctx)
	if err != nil {
		t.Fatalf("failed to list packages: %v", err)
	}

	if len(installed) != 1 {
		t.Errorf("expected 1 installed package, got %d", len(installed))
	}

	// Test uninstalling a package
	err = client.Uninstall(ctx, "client-test")
	if err != nil {
		t.Fatalf("failed to uninstall package: %v", err)
	}

	installed, _ = client.List(ctx)
	if len(installed) != 0 {
		t.Errorf("expected 0 installed packages after uninstall, got %d", len(installed))
	}
}

func TestMarketplaceCache(t *testing.T) {
	ctx := context.Background()
	cache := marketplace.NewMemoryCache()

	// Test putting data in cache
	data := []byte("test-package-data")
	err := cache.Put(ctx, "pkg1", "1.0.0", data)
	if err != nil {
		t.Fatalf("failed to put in cache: %v", err)
	}

	// Test checking if cached
	if !cache.Has(ctx, "pkg1", "1.0.0") {
		t.Error("expected package to be in cache")
	}

	// Test getting from cache
	retrieved, err := cache.Get(ctx, "pkg1", "1.0.0")
	if err != nil {
		t.Fatalf("failed to get from cache: %v", err)
	}

	if string(retrieved) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(retrieved))
	}

	// Test cache size
	size, err := cache.Size(ctx)
	if err != nil {
		t.Fatalf("failed to get cache size: %v", err)
	}

	if size != int64(len(data)) {
		t.Errorf("expected cache size %d, got %d", len(data), size)
	}

	// Test removing from cache
	err = cache.Remove(ctx, "pkg1", "1.0.0")
	if err != nil {
		t.Fatalf("failed to remove from cache: %v", err)
	}

	if cache.Has(ctx, "pkg1", "1.0.0") {
		t.Error("expected package to be removed from cache")
	}
}

func TestMarketplaceValidator(t *testing.T) {
	ctx := context.Background()
	validator := marketplace.NewValidator()

	// Test valid package
	validPkg := &marketplace.Package{
		ID:       "valid-pkg",
		Name:     "Valid Package",
		Version:  "1.0.0",
		Author:   "Author",
		Checksum: "checksum",
	}

	err := validator.ValidatePackage(ctx, validPkg, []byte("data"))
	if err != nil {
		t.Fatalf("validation failed for valid package: %v", err)
	}

	// Test invalid package (missing name)
	invalidPkg := &marketplace.Package{
		ID:      "invalid-pkg",
		Version: "1.0.0",
		Author:  "Author",
	}

	err = validator.ValidatePackage(ctx, invalidPkg, []byte("data"))
	if err == nil {
		t.Error("expected validation to fail for invalid package")
	}

	// Test template validation
	validTemplate := &marketplace.Template{
		ID:     "valid-tmpl",
		Name:   "Valid Template",
		Author: "Author",
		Files: []marketplace.TemplateFile{
			{Path: "main.go", Content: "package main"},
		},
	}

	err = validator.ValidateTemplate(ctx, validTemplate)
	if err != nil {
		t.Fatalf("validation failed for valid template: %v", err)
	}

	// Test security scanning
	issues, err := validator.ScanSecurity(ctx, make([]byte, 200*1024*1024))
	if err != nil {
		t.Fatalf("security scan failed: %v", err)
	}

	if len(issues) == 0 {
		t.Error("expected security issues for large package")
	}
}

func TestMarketplaceTemplate(t *testing.T) {
	ctx := context.Background()
	registry := marketplace.NewMemoryRegistry()
	cache := marketplace.NewMemoryCache()
	validator := marketplace.NewValidator()
	client := marketplace.NewClient(registry, cache, validator)

	// Publish a template
	tmpl := &marketplace.Template{
		ID:          "test-template",
		Name:        "Test Template",
		Description: "A test template",
		Author:      "Test Author",
		Languages:   []string{"go"},
		Category:    "web",
		Files: []marketplace.TemplateFile{
			{Path: "main.go", Content: "package main\n\nfunc main() {}\n"},
			{Path: "README.md", Content: "# Test Project\n"},
		},
		Variables: map[string]string{
			"project_name": "myproject",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := registry.PublishTemplate(ctx, tmpl)
	if err != nil {
		t.Fatalf("failed to publish template: %v", err)
	}

	// Test initializing from template
	err = client.InitFromTemplate(ctx, "test-template", "/tmp/test-project", tmpl.Variables)
	if err != nil {
		t.Fatalf("failed to init from template: %v", err)
	}

	// Test retrieving template
	retrieved, err := registry.GetTemplate(ctx, "test-template")
	if err != nil {
		t.Fatalf("failed to get template: %v", err)
	}

	if retrieved.Name != tmpl.Name {
		t.Errorf("expected name %s, got %s", tmpl.Name, retrieved.Name)
	}
}
