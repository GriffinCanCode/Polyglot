package updates

import (
	"context"
	"time"
)

// Version represents a software version
type Version struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	Prerelease string `json:"prerelease,omitempty"`
	Build      string `json:"build,omitempty"`
}

// Release represents a software release
type Release struct {
	Version     Version           `json:"version"`
	Channel     string            `json:"channel"`
	Platform    string            `json:"platform"`
	Arch        string            `json:"arch"`
	URL         string            `json:"url"`
	Size        int64             `json:"size"`
	Checksum    string            `json:"checksum"`
	Signature   []byte            `json:"signature"`
	ReleaseDate time.Time         `json:"release_date"`
	Notes       string            `json:"notes"`
	Critical    bool              `json:"critical"`
	Metadata    map[string]string `json:"metadata"`
}

// Update represents an available update
type Update struct {
	Current   Version           `json:"current"`
	Available Version           `json:"available"`
	Release   *Release          `json:"release"`
	Diff      *Diff             `json:"diff,omitempty"`
	Mandatory bool              `json:"mandatory"`
	Metadata  map[string]string `json:"metadata"`
}

// Diff represents a binary diff patch
type Diff struct {
	FromVersion Version `json:"from_version"`
	ToVersion   Version `json:"to_version"`
	Size        int64   `json:"size"`
	Checksum    string  `json:"checksum"`
	Data        []byte  `json:"data"`
	Algorithm   string  `json:"algorithm"`
	Compressed  bool    `json:"compressed"`
}

// DownloadProgress represents download progress
type DownloadProgress struct {
	BytesDownloaded int64     `json:"bytes_downloaded"`
	TotalBytes      int64     `json:"total_bytes"`
	Percentage      float64   `json:"percentage"`
	Speed           int64     `json:"speed"`          // bytes per second
	TimeRemaining   int64     `json:"time_remaining"` // seconds
	StartedAt       time.Time `json:"started_at"`
}

// ApplyResult represents the result of applying an update
type ApplyResult struct {
	Success    bool              `json:"success"`
	Version    Version           `json:"version"`
	Error      string            `json:"error,omitempty"`
	RollbackID string            `json:"rollback_id,omitempty"`
	Metadata   map[string]string `json:"metadata"`
	AppliedAt  time.Time         `json:"applied_at"`
}

// Checkpoint represents a restore point
type Checkpoint struct {
	ID        string            `json:"id"`
	Version   Version           `json:"version"`
	Binary    []byte            `json:"binary"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
}

// Manager handles update orchestration
type Manager interface {
	// Check checks for available updates
	Check(ctx context.Context, current Version, channel string) (*Update, error)

	// Download downloads an update
	Download(ctx context.Context, update *Update, progress chan<- *DownloadProgress) ([]byte, error)

	// Apply applies an update
	Apply(ctx context.Context, data []byte, update *Update) (*ApplyResult, error)

	// Rollback rolls back to a previous version
	Rollback(ctx context.Context, checkpointID string) error

	// CreateCheckpoint creates a restore point
	CreateCheckpoint(ctx context.Context, version Version, binary []byte) (*Checkpoint, error)

	// ListCheckpoints lists available restore points
	ListCheckpoints(ctx context.Context) ([]*Checkpoint, error)
}

// Differ generates binary diffs
type Differ interface {
	// Generate generates a diff between two binaries
	Generate(ctx context.Context, oldBinary, newBinary []byte) (*Diff, error)

	// Apply applies a diff to a binary
	Apply(ctx context.Context, oldBinary []byte, diff *Diff) ([]byte, error)

	// Compress compresses a diff
	Compress(ctx context.Context, diff *Diff) (*Diff, error)

	// Decompress decompresses a diff
	Decompress(ctx context.Context, diff *Diff) (*Diff, error)
}

// Downloader handles update downloads
type Downloader interface {
	// Download downloads data from a URL
	Download(ctx context.Context, url string, progress chan<- *DownloadProgress) ([]byte, error)

	// Verify verifies downloaded data
	Verify(ctx context.Context, data []byte, checksum string) error

	// Resume resumes a partial download
	Resume(ctx context.Context, url string, offset int64, progress chan<- *DownloadProgress) ([]byte, error)
}

// Verifier handles update verification
type Verifier interface {
	// VerifySignature verifies update signature
	VerifySignature(ctx context.Context, data []byte, signature []byte) error

	// VerifyChecksum verifies data checksum
	VerifyChecksum(ctx context.Context, data []byte, checksum string) error

	// VerifyVersion verifies version compatibility
	VerifyVersion(ctx context.Context, current, target Version) error
}
