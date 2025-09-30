package updates

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DefaultManager implements the update manager
type DefaultManager struct {
	mu          sync.RWMutex
	differ      Differ
	downloader  Downloader
	verifier    Verifier
	checkpoints map[string]*Checkpoint
	releases    map[string]*Release
}

// NewManager creates a new update manager
func NewManager(differ Differ, downloader Downloader, verifier Verifier) *DefaultManager {
	return &DefaultManager{
		differ:      differ,
		downloader:  downloader,
		verifier:    verifier,
		checkpoints: make(map[string]*Checkpoint),
		releases:    make(map[string]*Release),
	}
}

// Check checks for available updates
func (m *DefaultManager) Check(ctx context.Context, current Version, channel string) (*Update, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Find latest release for channel
	var latest *Release
	for _, release := range m.releases {
		if release.Channel != channel {
			continue
		}
		if latest == nil || compareVersions(release.Version, latest.Version) > 0 {
			latest = release
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no releases found for channel: %s", channel)
	}

	// Check if update is needed
	if compareVersions(latest.Version, current) <= 0 {
		return nil, nil // Already up to date
	}

	update := &Update{
		Current:   current,
		Available: latest.Version,
		Release:   latest,
		Mandatory: latest.Critical,
		Metadata:  make(map[string]string),
	}

	return update, nil
}

// Download downloads an update
func (m *DefaultManager) Download(ctx context.Context, update *Update, progress chan<- *DownloadProgress) ([]byte, error) {
	if update == nil || update.Release == nil {
		return nil, fmt.Errorf("invalid update")
	}

	// Download the update
	data, err := m.downloader.Download(ctx, update.Release.URL, progress)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}

	// Verify checksum
	if err := m.verifier.VerifyChecksum(ctx, data, update.Release.Checksum); err != nil {
		return nil, fmt.Errorf("checksum verification failed: %w", err)
	}

	// Verify signature
	if err := m.verifier.VerifySignature(ctx, data, update.Release.Signature); err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return data, nil
}

// Apply applies an update
func (m *DefaultManager) Apply(ctx context.Context, data []byte, update *Update) (*ApplyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to apply")
	}

	// Verify version compatibility
	if err := m.verifier.VerifyVersion(ctx, update.Current, update.Available); err != nil {
		return nil, fmt.Errorf("version verification failed: %w", err)
	}

	// Create checkpoint before applying
	checkpoint, err := m.CreateCheckpoint(ctx, update.Current, []byte("current-binary"))
	if err != nil {
		return nil, fmt.Errorf("failed to create checkpoint: %w", err)
	}

	// Apply the update (simplified - would actually replace binary)
	result := &ApplyResult{
		Success:    true,
		Version:    update.Available,
		RollbackID: checkpoint.ID,
		Metadata:   make(map[string]string),
		AppliedAt:  time.Now(),
	}

	return result, nil
}

// Rollback rolls back to a previous version
func (m *DefaultManager) Rollback(ctx context.Context, checkpointID string) error {
	m.mu.RLock()
	checkpoint, ok := m.checkpoints[checkpointID]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("checkpoint not found: %s", checkpointID)
	}

	// Restore from checkpoint (simplified)
	_ = checkpoint.Binary

	return nil
}

// CreateCheckpoint creates a restore point
func (m *DefaultManager) CreateCheckpoint(ctx context.Context, version Version, binary []byte) (*Checkpoint, error) {
	checkpoint := &Checkpoint{
		ID:        fmt.Sprintf("checkpoint-%d", time.Now().UnixNano()),
		Version:   version,
		Binary:    binary,
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	m.mu.Lock()
	m.checkpoints[checkpoint.ID] = checkpoint
	m.mu.Unlock()

	return checkpoint, nil
}

// ListCheckpoints lists available restore points
func (m *DefaultManager) ListCheckpoints(ctx context.Context) ([]*Checkpoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	checkpoints := make([]*Checkpoint, 0, len(m.checkpoints))
	for _, cp := range m.checkpoints {
		checkpoints = append(checkpoints, cp)
	}

	return checkpoints, nil
}

// AddRelease adds a release (for testing)
func (m *DefaultManager) AddRelease(release *Release) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s-%d.%d.%d", release.Channel, release.Version.Major, release.Version.Minor, release.Version.Patch)
	m.releases[key] = release
}

// Helper function to compare versions
func compareVersions(a, b Version) int {
	if a.Major != b.Major {
		return a.Major - b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor - b.Minor
	}
	if a.Patch != b.Patch {
		return a.Patch - b.Patch
	}
	return 0
}
