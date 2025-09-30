package cloud

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryBuilder implements an in-memory builder for testing
type MemoryBuilder struct {
	mu     sync.RWMutex
	builds map[string]*BuildResult
}

// NewMemoryBuilder creates a new in-memory builder
func NewMemoryBuilder() *MemoryBuilder {
	return &MemoryBuilder{
		builds: make(map[string]*BuildResult),
	}
}

// Build submits a build request
func (b *MemoryBuilder) Build(ctx context.Context, req *BuildRequest) (*BuildResult, error) {
	if req.ProjectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if len(req.Source) == 0 {
		return nil, fmt.Errorf("source is required")
	}

	buildID := fmt.Sprintf("build-%d", time.Now().UnixNano())

	result := &BuildResult{
		ID:          buildID,
		RequestID:   req.ID,
		Platform:    req.Platform,
		Binary:      []byte(fmt.Sprintf("binary-%s-%s-%s", req.ProjectID, req.Platform.OS, req.Platform.Arch)),
		Artifacts:   make(map[string][]byte),
		Logs:        fmt.Sprintf("Building for %s/%s\nBuild completed successfully", req.Platform.OS, req.Platform.Arch),
		Status:      "completed",
		Duration:    time.Second * 30,
		BinarySize:  int64(len(req.Source) * 2),
		CompletedAt: time.Now(),
	}

	b.mu.Lock()
	b.builds[buildID] = result
	b.mu.Unlock()

	return result, nil
}

// GetBuild retrieves build status
func (b *MemoryBuilder) GetBuild(ctx context.Context, buildID string) (*BuildResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result, ok := b.builds[buildID]
	if !ok {
		return nil, fmt.Errorf("build not found: %s", buildID)
	}

	return result, nil
}

// CancelBuild cancels an in-progress build
func (b *MemoryBuilder) CancelBuild(ctx context.Context, buildID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	result, ok := b.builds[buildID]
	if !ok {
		return fmt.Errorf("build not found: %s", buildID)
	}

	if result.Status == "completed" || result.Status == "failed" {
		return fmt.Errorf("build already finished")
	}

	result.Status = "cancelled"
	return nil
}

// ListBuilds lists recent builds
func (b *MemoryBuilder) ListBuilds(ctx context.Context, projectID string, limit int) ([]*BuildResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	results := make([]*BuildResult, 0)
	for _, result := range b.builds {
		if result.RequestID == projectID || projectID == "" {
			results = append(results, result)
			if len(results) >= limit && limit > 0 {
				break
			}
		}
	}

	return results, nil
}
