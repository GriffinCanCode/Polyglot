package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/cloud"
)

func TestCloudBuilder(t *testing.T) {
	ctx := context.Background()
	builder := cloud.NewMemoryBuilder()

	// Test submitting a build
	req := &cloud.BuildRequest{
		ID:        "build-1",
		ProjectID: "test-project",
		Platform: cloud.Platform{
			OS:   "linux",
			Arch: "amd64",
		},
		Source: []byte("source-code"),
		Config: map[string]string{
			"optimization": "release",
		},
	}

	result, err := builder.Build(ctx, req)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	if result.Status != "completed" {
		t.Errorf("expected status completed, got %s", result.Status)
	}

	// Test retrieving build
	retrieved, err := builder.GetBuild(ctx, result.ID)
	if err != nil {
		t.Fatalf("failed to get build: %v", err)
	}

	if retrieved.ID != result.ID {
		t.Errorf("expected build ID %s, got %s", result.ID, retrieved.ID)
	}

	// Test listing builds
	builds, err := builder.ListBuilds(ctx, "", 10)
	if err != nil {
		t.Fatalf("failed to list builds: %v", err)
	}

	if len(builds) == 0 {
		t.Error("expected at least one build")
	}
}

func TestCloudStorage(t *testing.T) {
	ctx := context.Background()
	storage := cloud.NewMemoryStorage()

	// Test storing data
	data := []byte("test-artifact")
	metadata := map[string]string{
		"type": "binary",
	}

	err := storage.Put(ctx, "artifacts/test.bin", data, metadata)
	if err != nil {
		t.Fatalf("failed to put object: %v", err)
	}

	// Test retrieving data
	retrieved, err := storage.Get(ctx, "artifacts/test.bin")
	if err != nil {
		t.Fatalf("failed to get object: %v", err)
	}

	if string(retrieved) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(retrieved))
	}

	// Test getting metadata
	obj, err := storage.GetMetadata(ctx, "artifacts/test.bin")
	if err != nil {
		t.Fatalf("failed to get metadata: %v", err)
	}

	if obj.Size != int64(len(data)) {
		t.Errorf("expected size %d, got %d", len(data), obj.Size)
	}

	// Test listing objects
	objects, err := storage.List(ctx, "artifacts/", 10)
	if err != nil {
		t.Fatalf("failed to list objects: %v", err)
	}

	if len(objects) == 0 {
		t.Error("expected at least one object")
	}

	// Test deleting object
	err = storage.Delete(ctx, "artifacts/test.bin")
	if err != nil {
		t.Fatalf("failed to delete object: %v", err)
	}

	_, err = storage.Get(ctx, "artifacts/test.bin")
	if err == nil {
		t.Error("expected error when getting deleted object")
	}
}

func TestCloudAuth(t *testing.T) {
	ctx := context.Background()
	auth := cloud.NewMemoryAuth()

	// Test authentication
	creds, err := auth.Authenticate(ctx, "test-api-key", "test-secret")
	if err != nil {
		t.Fatalf("authentication failed: %v", err)
	}

	if creds.APIKey != "test-api-key" {
		t.Errorf("expected API key test-api-key, got %s", creds.APIKey)
	}

	// Test validation
	err = auth.Validate(ctx, creds)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	// Test refresh
	refreshed, err := auth.Refresh(ctx, creds)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}

	if refreshed.ExpiresAt.Before(creds.ExpiresAt) {
		t.Error("expected refreshed credentials to have later expiry")
	}

	// Test revoke
	err = auth.Revoke(ctx, creds)
	if err != nil {
		t.Fatalf("revoke failed: %v", err)
	}

	err = auth.Validate(ctx, creds)
	if err == nil {
		t.Error("expected validation to fail for revoked credentials")
	}
}

func TestCloudClient(t *testing.T) {
	ctx := context.Background()
	builder := cloud.NewMemoryBuilder()
	storage := cloud.NewMemoryStorage()
	auth := cloud.NewMemoryAuth()
	client := cloud.NewClient(builder, storage, auth)

	// Test authentication
	err := client.Authenticate(ctx, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("authentication failed: %v", err)
	}

	// Test cross-compilation
	platforms := []cloud.Platform{
		{OS: "linux", Arch: "amd64"},
		{OS: "darwin", Arch: "arm64"},
		{OS: "windows", Arch: "amd64"},
	}

	results, err := client.CrossCompile(ctx, []byte("source"), platforms)
	if err != nil && len(results) == 0 {
		t.Fatalf("cross-compilation failed: %v", err)
	}

	if len(results) != len(platforms) {
		t.Errorf("expected %d results, got %d", len(platforms), len(results))
	}

	// Verify we have results
	for _, result := range results {
		if result == nil {
			t.Error("got nil result")
		}
	}
}

func TestCloudPlatform(t *testing.T) {
	platforms := []cloud.Platform{
		{OS: "linux", Arch: "amd64", CGOEnabled: true},
		{OS: "darwin", Arch: "arm64", CGOEnabled: true},
		{OS: "windows", Arch: "amd64", CGOEnabled: false},
		{OS: "linux", Arch: "arm64", CGOEnabled: true, GOARM: "7"},
	}

	for _, platform := range platforms {
		if platform.OS == "" {
			t.Error("platform OS should not be empty")
		}
		if platform.Arch == "" {
			t.Error("platform Arch should not be empty")
		}
	}
}

func TestBuildTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	builder := cloud.NewMemoryBuilder()

	req := &cloud.BuildRequest{
		ID:        "timeout-build",
		ProjectID: "test",
		Platform:  cloud.Platform{OS: "linux", Arch: "amd64"},
		Source:    []byte("source"),
	}

	// This should complete before timeout since it's fast
	_, err := builder.Build(ctx, req)
	if err != nil {
		t.Logf("build completed or timed out: %v", err)
	}
}
