package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/updates"
)

func TestUpdateManager(t *testing.T) {
	ctx := context.Background()
	differ := updates.NewDiffer()
	downloader := updates.NewDownloader()
	verifier := updates.NewVerifier()
	manager := updates.NewManager(differ, downloader, verifier)

	// Add a test release
	release := &updates.Release{
		Version:     updates.Version{Major: 2, Minor: 0, Patch: 0},
		Channel:     "stable",
		Platform:    "linux",
		Arch:        "amd64",
		URL:         "https://example.com/v2.0.0",
		Size:        2048,
		Checksum:    "sha256-2048",
		Signature:   []byte("signature"),
		ReleaseDate: time.Now(),
		Notes:       "Version 2.0.0 release",
	}
	manager.AddRelease(release)

	// Check for updates
	current := updates.Version{Major: 1, Minor: 0, Patch: 0}
	update, err := manager.Check(ctx, current, "stable")
	if err != nil {
		t.Fatalf("failed to check for updates: %v", err)
	}

	if update == nil {
		t.Fatal("expected update to be available")
	}

	if update.Available.Major != 2 {
		t.Errorf("expected version 2.0.0, got %d.%d.%d",
			update.Available.Major, update.Available.Minor, update.Available.Patch)
	}

	// Download update
	progressChan := make(chan *updates.DownloadProgress, 10)
	go func() {
		for range progressChan {
			// Consume progress updates
		}
	}()

	data, err := manager.Download(ctx, update, progressChan)
	if err != nil {
		t.Fatalf("failed to download update: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected downloaded data to not be empty")
	}

	// Apply update
	result, err := manager.Apply(ctx, data, update)
	if err != nil {
		t.Fatalf("failed to apply update: %v", err)
	}

	if !result.Success {
		t.Error("expected update to be applied successfully")
	}

	if result.RollbackID == "" {
		t.Error("expected rollback ID to be set")
	}

	// Test rollback
	err = manager.Rollback(ctx, result.RollbackID)
	if err != nil {
		t.Fatalf("failed to rollback: %v", err)
	}
}

func TestBinaryDiff(t *testing.T) {
	ctx := context.Background()
	differ := updates.NewDiffer()

	oldBinary := []byte("old binary content")
	newBinary := []byte("new binary content with more data")

	// Generate diff
	diff, err := differ.Generate(ctx, oldBinary, newBinary)
	if err != nil {
		t.Fatalf("failed to generate diff: %v", err)
	}

	if diff.Size == 0 {
		t.Error("expected diff size to be non-zero")
	}

	// Apply diff
	result, err := differ.Apply(ctx, oldBinary, diff)
	if err != nil {
		t.Fatalf("failed to apply diff: %v", err)
	}

	if string(result) != string(newBinary) {
		t.Errorf("expected %s, got %s", string(newBinary), string(result))
	}

	// Test compression
	compressed, err := differ.Compress(ctx, diff)
	if err != nil {
		t.Fatalf("failed to compress diff: %v", err)
	}

	if !compressed.Compressed {
		t.Error("expected diff to be marked as compressed")
	}

	// Test decompression
	decompressed, err := differ.Decompress(ctx, compressed)
	if err != nil {
		t.Fatalf("failed to decompress diff: %v", err)
	}

	if decompressed.Compressed {
		t.Error("expected diff to be marked as not compressed")
	}
}

func TestDownloader(t *testing.T) {
	ctx := context.Background()
	downloader := updates.NewDownloader()

	// Test download with progress
	progressChan := make(chan *updates.DownloadProgress, 20)
	progressReceived := false

	go func() {
		for p := range progressChan {
			progressReceived = true
			if p.Percentage < 0 || p.Percentage > 100 {
				t.Errorf("invalid percentage: %f", p.Percentage)
			}
		}
	}()

	data, err := downloader.Download(ctx, "https://example.com/file", progressChan)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected downloaded data to not be empty")
	}

	if !progressReceived {
		t.Error("expected to receive progress updates")
	}

	// Test verification
	checksum := fmt.Sprintf("sha256-%d", len(data))
	err = downloader.Verify(ctx, data, checksum)
	if err != nil {
		t.Fatalf("verification failed: %v", err)
	}
}

func TestVerifier(t *testing.T) {
	ctx := context.Background()
	verifier := updates.NewVerifier()

	data := []byte("test data")
	checksum := "sha256-9"
	signature := []byte("signature")

	// Test checksum verification
	err := verifier.VerifyChecksum(ctx, data, checksum)
	if err != nil {
		t.Fatalf("checksum verification failed: %v", err)
	}

	// Test signature verification
	err = verifier.VerifySignature(ctx, data, signature)
	if err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}

	// Test version compatibility
	current := updates.Version{Major: 1, Minor: 5, Patch: 0}
	target := updates.Version{Major: 2, Minor: 0, Patch: 0}

	err = verifier.VerifyVersion(ctx, current, target)
	if err != nil {
		t.Fatalf("version verification failed: %v", err)
	}

	// Test downgrade prevention
	downgrade := updates.Version{Major: 0, Minor: 9, Patch: 0}
	err = verifier.VerifyVersion(ctx, current, downgrade)
	if err == nil {
		t.Error("expected error when downgrading major version")
	}
}

func TestCheckpoints(t *testing.T) {
	ctx := context.Background()
	manager := updates.NewManager(
		updates.NewDiffer(),
		updates.NewDownloader(),
		updates.NewVerifier(),
	)

	// Create checkpoint
	version := updates.Version{Major: 1, Minor: 0, Patch: 0}
	binary := []byte("current binary")

	checkpoint, err := manager.CreateCheckpoint(ctx, version, binary)
	if err != nil {
		t.Fatalf("failed to create checkpoint: %v", err)
	}

	if checkpoint.ID == "" {
		t.Error("expected checkpoint ID to be set")
	}

	// List checkpoints
	checkpoints, err := manager.ListCheckpoints(ctx)
	if err != nil {
		t.Fatalf("failed to list checkpoints: %v", err)
	}

	if len(checkpoints) == 0 {
		t.Error("expected at least one checkpoint")
	}

	// Rollback to checkpoint
	err = manager.Rollback(ctx, checkpoint.ID)
	if err != nil {
		t.Fatalf("failed to rollback: %v", err)
	}
}

func TestVersionComparison(t *testing.T) {
	v1 := updates.Version{Major: 1, Minor: 0, Patch: 0}
	v2 := updates.Version{Major: 2, Minor: 0, Patch: 0}
	v1_5 := updates.Version{Major: 1, Minor: 5, Patch: 0}
	v1_5_1 := updates.Version{Major: 1, Minor: 5, Patch: 1}

	tests := []struct {
		a        updates.Version
		b        updates.Version
		expected string
	}{
		{v1, v2, "less"},
		{v2, v1, "greater"},
		{v1, v1, "equal"},
		{v1, v1_5, "less"},
		{v1_5, v1_5_1, "less"},
		{v1_5_1, v1_5, "greater"},
	}

	for _, tt := range tests {
		var result string
		if tt.a.Major < tt.b.Major {
			result = "less"
		} else if tt.a.Major > tt.b.Major {
			result = "greater"
		} else if tt.a.Minor < tt.b.Minor {
			result = "less"
		} else if tt.a.Minor > tt.b.Minor {
			result = "greater"
		} else if tt.a.Patch < tt.b.Patch {
			result = "less"
		} else if tt.a.Patch > tt.b.Patch {
			result = "greater"
		} else {
			result = "equal"
		}

		if result != tt.expected {
			t.Errorf("comparing %v and %v: expected %s, got %s",
				tt.a, tt.b, tt.expected, result)
		}
	}
}
