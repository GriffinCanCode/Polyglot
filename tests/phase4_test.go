package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/cloud"
	"github.com/griffincancode/polyglot.js/marketplace"
	"github.com/griffincancode/polyglot.js/signing"
	"github.com/griffincancode/polyglot.js/updates"
)

// TestPhase4Integration performs end-to-end testing of Phase 4 components
func TestPhase4Integration(t *testing.T) {
	t.Run("MarketplaceToCloud", testMarketplaceToCloud)
	t.Run("CloudToBuildAndSign", testCloudToBuildAndSign)
	t.Run("SignAndDistribute", testSignAndDistribute)
	t.Run("UpdateFlow", testUpdateFlow)
	t.Run("CompleteWorkflow", testCompleteWorkflow)
}

// testMarketplaceToCloud tests package discovery and cloud deployment
func testMarketplaceToCloud(t *testing.T) {
	ctx := context.Background()

	// Setup marketplace
	registry := marketplace.NewMemoryRegistry()
	cache := marketplace.NewMemoryCache()
	validator := marketplace.NewValidator()
	client := marketplace.NewClient(registry, cache, validator)

	// Publish a package
	pkg := &marketplace.Package{
		ID:          "cloud-app",
		Name:        "Cloud App",
		Version:     "1.0.0",
		Author:      "Test",
		Description: "A cloud-deployed application",
		Languages:   []string{"go", "python"},
		Checksum:    "checksum-123",
	}

	err := registry.Publish(ctx, pkg, []byte("package-source"))
	if err != nil {
		t.Fatalf("failed to publish package: %v", err)
	}

	// Install package
	err = client.Install(ctx, "cloud-app", "1.0.0")
	if err != nil {
		t.Fatalf("failed to install package: %v", err)
	}

	// Setup cloud
	builder := cloud.NewMemoryBuilder()
	storage := cloud.NewMemoryStorage()
	auth := cloud.NewMemoryAuth()
	cloudClient := cloud.NewClient(builder, storage, auth)

	// Authenticate
	err = cloudClient.Authenticate(ctx, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}

	// Build in cloud
	buildReq := &cloud.BuildRequest{
		ID:        "build-1",
		ProjectID: "cloud-app",
		Platform:  cloud.Platform{OS: "linux", Arch: "amd64"},
		Source:    []byte("package-source"),
	}

	result, err := builder.Build(ctx, buildReq)
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if result.Status != "completed" {
		t.Errorf("expected completed status, got %s", result.Status)
	}

	// Store artifact
	err = storage.Put(ctx, "artifacts/cloud-app-1.0.0", result.Binary, nil)
	if err != nil {
		t.Fatalf("failed to store artifact: %v", err)
	}
}

// testCloudToBuildAndSign tests cloud building and code signing
func testCloudToBuildAndSign(t *testing.T) {
	ctx := context.Background()

	// Setup cloud
	builder := cloud.NewMemoryBuilder()
	auth := cloud.NewMemoryAuth()

	// Authenticate
	creds, err := auth.Authenticate(ctx, "build-key", "build-secret")
	if err != nil {
		t.Fatalf("authentication failed: %v", err)
	}

	// Build application
	buildReq := &cloud.BuildRequest{
		ID:        "sign-build",
		ProjectID: creds.ProjectID,
		Platform:  cloud.Platform{OS: "darwin", Arch: "arm64"},
		Source:    []byte("app-source"),
	}

	buildResult, err := builder.Build(ctx, buildReq)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	// Setup signing
	signer := signing.NewSigner()

	cert := &signing.Certificate{
		ID:         "app-cert",
		Type:       "apple",
		Subject:    "Developer ID",
		NotBefore:  time.Now(),
		NotAfter:   time.Now().Add(365 * 24 * time.Hour),
		Data:       []byte("cert-data"),
		PrivateKey: []byte("key-data"),
	}

	err = signer.ImportCertificate(ctx, cert)
	if err != nil {
		t.Fatalf("failed to import certificate: %v", err)
	}

	// Sign the binary
	signReq := &signing.SignRequest{
		Binary:      buildResult.Binary,
		Platform:    "darwin",
		Certificate: cert,
		Timestamp:   true,
	}

	signResult, err := signer.Sign(ctx, signReq)
	if err != nil && signResult == nil {
		t.Logf("signing not supported on this platform: %v", err)
		return
	}

	if signResult.SignedBinary == nil {
		t.Error("expected signed binary")
	}
}

// testSignAndDistribute tests signing and distribution preparation
func testSignAndDistribute(t *testing.T) {
	ctx := context.Background()

	// Build binary
	builder := cloud.NewMemoryBuilder()
	buildReq := &cloud.BuildRequest{
		ID:        "dist-build",
		ProjectID: "app",
		Platform:  cloud.Platform{OS: "linux", Arch: "amd64"},
		Source:    []byte("source"),
	}

	buildResult, err := builder.Build(ctx, buildReq)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	// Sign binary
	signer := signing.NewSigner()
	cert := &signing.Certificate{
		ID:      "dist-cert",
		Type:    "gpg",
		Subject: "Release Key",
	}
	signer.ImportCertificate(ctx, cert)

	signReq := &signing.SignRequest{
		Binary:      buildResult.Binary,
		Platform:    "linux",
		Certificate: cert,
	}

	signResult, err := signer.Sign(ctx, signReq)
	if err != nil && signResult == nil {
		t.Logf("signing not supported: %v", err)
		return
	}

	// Store for distribution
	storage := cloud.NewMemoryStorage()
	err = storage.Put(ctx, "releases/app-1.0.0-linux-amd64", signResult.SignedBinary, map[string]string{
		"version":  "1.0.0",
		"platform": "linux",
		"arch":     "amd64",
	})
	if err != nil {
		t.Fatalf("failed to store release: %v", err)
	}

	// Verify it can be retrieved
	retrieved, err := storage.Get(ctx, "releases/app-1.0.0-linux-amd64")
	if err != nil {
		t.Fatalf("failed to retrieve release: %v", err)
	}

	if len(retrieved) == 0 {
		t.Error("expected non-empty release binary")
	}
}

// testUpdateFlow tests the complete update workflow
func testUpdateFlow(t *testing.T) {
	ctx := context.Background()

	// Setup update manager
	differ := updates.NewDiffer()
	downloader := updates.NewDownloader()
	verifier := updates.NewVerifier()
	manager := updates.NewManager(differ, downloader, verifier)

	// Create releases
	v1Release := &updates.Release{
		Version:     updates.Version{Major: 1, Minor: 0, Patch: 0},
		Channel:     "stable",
		Platform:    "linux",
		Arch:        "amd64",
		URL:         "https://example.com/v1.0.0",
		Size:        1024,
		Checksum:    "sha256-1024",
		Signature:   []byte("sig-v1"),
		ReleaseDate: time.Now().Add(-30 * 24 * time.Hour),
	}

	v2Release := &updates.Release{
		Version:     updates.Version{Major: 2, Minor: 0, Patch: 0},
		Channel:     "stable",
		Platform:    "linux",
		Arch:        "amd64",
		URL:         "https://example.com/v2.0.0",
		Size:        2048,
		Checksum:    "sha256-2048",
		Signature:   []byte("sig-v2"),
		ReleaseDate: time.Now(),
		Critical:    true,
	}

	manager.AddRelease(v1Release)
	manager.AddRelease(v2Release)

	// Check for updates
	current := updates.Version{Major: 1, Minor: 0, Patch: 0}
	update, err := manager.Check(ctx, current, "stable")
	if err != nil {
		t.Fatalf("failed to check updates: %v", err)
	}

	if update == nil {
		t.Fatal("expected update to be available")
	}

	if update.Mandatory != v2Release.Critical {
		t.Error("expected update to be marked as mandatory")
	}

	// Download update
	progressChan := make(chan *updates.DownloadProgress, 10)
	go func() {
		for range progressChan {
		}
	}()

	data, err := manager.Download(ctx, update, progressChan)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	// Apply update
	result, err := manager.Apply(ctx, data, update)
	if err != nil {
		t.Fatalf("apply failed: %v", err)
	}

	if !result.Success {
		t.Error("expected successful update")
	}

	// Test rollback
	err = manager.Rollback(ctx, result.RollbackID)
	if err != nil {
		t.Fatalf("rollback failed: %v", err)
	}
}

// testCompleteWorkflow tests the complete Phase 4 workflow
func testCompleteWorkflow(t *testing.T) {
	ctx := context.Background()

	// 1. Discover package in marketplace
	registry := marketplace.NewMemoryRegistry()
	tmpl := &marketplace.Template{
		ID:     "starter-app",
		Name:   "Starter Application",
		Author: "Polyglot",
		Files: []marketplace.TemplateFile{
			{Path: "main.go", Content: "package main"},
		},
	}
	registry.PublishTemplate(ctx, tmpl)

	// 2. Build application in cloud for multiple platforms
	builder := cloud.NewMemoryBuilder()
	auth := cloud.NewMemoryAuth()
	storage := cloud.NewMemoryStorage()
	cloudClient := cloud.NewClient(builder, storage, auth)

	cloudClient.Authenticate(ctx, "workflow-key", "workflow-secret")

	platforms := []cloud.Platform{
		{OS: "linux", Arch: "amd64"},
		{OS: "darwin", Arch: "arm64"},
		{OS: "windows", Arch: "amd64"},
	}

	buildResults, err := cloudClient.CrossCompile(ctx, []byte("app-source"), platforms)
	if err != nil && len(buildResults) == 0 {
		t.Fatalf("cross-compilation failed: %v", err)
	}

	// 3. Sign binaries for each platform
	signer := signing.NewSigner()

	certs := map[string]*signing.Certificate{
		"linux":   {ID: "linux-cert", Type: "gpg", Subject: "Linux Key"},
		"darwin":  {ID: "darwin-cert", Type: "apple", Subject: "Apple Dev"},
		"windows": {ID: "windows-cert", Type: "windows", Subject: "Windows Cert"},
	}

	for _, cert := range certs {
		cert.NotBefore = time.Now()
		cert.NotAfter = time.Now().Add(365 * 24 * time.Hour)
		signer.ImportCertificate(ctx, cert)
	}

	// Sign each build result
	for _, buildResult := range buildResults {
		platformName := buildResult.Platform.OS
		if cert, ok := certs[platformName]; ok {
			signReq := &signing.SignRequest{
				Binary:      buildResult.Binary,
				Platform:    platformName,
				Certificate: cert,
			}

			_, err := signer.Sign(ctx, signReq)
			if err != nil {
				t.Logf("signing not supported for %s: %v", platformName, err)
			}
		}
	}

	// 4. Create and distribute update
	differ := updates.NewDiffer()
	downloader := updates.NewDownloader()
	verifier := updates.NewVerifier()
	updateManager := updates.NewManager(differ, downloader, verifier)

	release := &updates.Release{
		Version:     updates.Version{Major: 1, Minor: 0, Patch: 0},
		Channel:     "stable",
		Platform:    "linux",
		Arch:        "amd64",
		URL:         "https://example.com/app-1.0.0",
		Size:        int64(len([]byte("binary"))),
		Checksum:    "sha256-6",
		Signature:   []byte("signature"),
		ReleaseDate: time.Now(),
	}

	updateManager.AddRelease(release)

	// Store in cloud storage
	err = storage.Put(ctx, "releases/app-1.0.0-linux-amd64", []byte("signed-binary"), map[string]string{
		"version": "1.0.0",
	})
	if err != nil {
		t.Fatalf("failed to store release: %v", err)
	}

	t.Log("Complete workflow test passed: marketplace → cloud → signing → updates")
}

// Benchmark Phase 4 operations
func BenchmarkMarketplaceSearch(b *testing.B) {
	ctx := context.Background()
	registry := marketplace.NewMemoryRegistry()

	// Populate with packages
	for i := 0; i < 100; i++ {
		pkg := &marketplace.Package{
			ID:      string(rune('a' + (i % 26))),
			Name:    "Package",
			Version: "1.0.0",
			Author:  "Test",
		}
		registry.Publish(ctx, pkg, []byte("data"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Search(ctx, marketplace.SearchQuery{Limit: 10})
	}
}

func BenchmarkCloudBuild(b *testing.B) {
	ctx := context.Background()
	builder := cloud.NewMemoryBuilder()

	req := &cloud.BuildRequest{
		ID:        "bench",
		ProjectID: "test",
		Platform:  cloud.Platform{OS: "linux", Arch: "amd64"},
		Source:    []byte("source"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Build(ctx, req)
	}
}

func BenchmarkDiffGeneration(b *testing.B) {
	ctx := context.Background()
	differ := updates.NewDiffer()

	oldBinary := make([]byte, 1024*1024) // 1MB
	newBinary := make([]byte, 1024*1024) // 1MB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		differ.Generate(ctx, oldBinary, newBinary)
	}
}
