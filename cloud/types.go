package cloud

import (
	"context"
	"time"
)

// Platform represents a target compilation platform
type Platform struct {
	OS           string   `json:"os"`
	Arch         string   `json:"arch"`
	Variant      string   `json:"variant"`
	CGOEnabled   bool     `json:"cgo_enabled"`
	Tags         []string `json:"tags"`
	GOARM        string   `json:"goarm,omitempty"`
	MinOSVersion string   `json:"min_os_version,omitempty"`
}

// BuildRequest represents a remote build request
type BuildRequest struct {
	ID           string            `json:"id"`
	ProjectID    string            `json:"project_id"`
	Platform     Platform          `json:"platform"`
	Source       []byte            `json:"source"`
	Config       map[string]string `json:"config"`
	Languages    []string          `json:"languages"`
	Tags         []string          `json:"tags"`
	Optimization string            `json:"optimization"`
	CreatedAt    time.Time         `json:"created_at"`
}

// BuildResult represents a completed build
type BuildResult struct {
	ID          string            `json:"id"`
	RequestID   string            `json:"request_id"`
	Platform    Platform          `json:"platform"`
	Binary      []byte            `json:"binary"`
	Artifacts   map[string][]byte `json:"artifacts"`
	Logs        string            `json:"logs"`
	Status      string            `json:"status"`
	Error       string            `json:"error,omitempty"`
	Duration    time.Duration     `json:"duration"`
	BinarySize  int64             `json:"binary_size"`
	CompletedAt time.Time         `json:"completed_at"`
}

// Credentials represents authentication credentials
type Credentials struct {
	APIKey    string    `json:"api_key"`
	SecretKey string    `json:"secret_key"`
	ProjectID string    `json:"project_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// StorageObject represents a stored artifact
type StorageObject struct {
	Key         string            `json:"key"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Metadata    map[string]string `json:"metadata"`
	Checksum    string            `json:"checksum"`
	CreatedAt   time.Time         `json:"created_at"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
}

// Builder handles remote build orchestration
type Builder interface {
	// Build submits a build request
	Build(ctx context.Context, req *BuildRequest) (*BuildResult, error)

	// GetBuild retrieves build status
	GetBuild(ctx context.Context, buildID string) (*BuildResult, error)

	// CancelBuild cancels an in-progress build
	CancelBuild(ctx context.Context, buildID string) error

	// ListBuilds lists recent builds
	ListBuilds(ctx context.Context, projectID string, limit int) ([]*BuildResult, error)
}

// Storage handles artifact storage
type Storage interface {
	// Put stores an artifact
	Put(ctx context.Context, key string, data []byte, metadata map[string]string) error

	// Get retrieves an artifact
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete removes an artifact
	Delete(ctx context.Context, key string) error

	// List lists artifacts matching a prefix
	List(ctx context.Context, prefix string, limit int) ([]*StorageObject, error)

	// GetMetadata retrieves artifact metadata
	GetMetadata(ctx context.Context, key string) (*StorageObject, error)
}

// Auth handles authentication and authorization
type Auth interface {
	// Authenticate validates credentials
	Authenticate(ctx context.Context, apiKey, secretKey string) (*Credentials, error)

	// Refresh refreshes expired credentials
	Refresh(ctx context.Context, creds *Credentials) (*Credentials, error)

	// Validate validates credentials
	Validate(ctx context.Context, creds *Credentials) error

	// Revoke revokes credentials
	Revoke(ctx context.Context, creds *Credentials) error
}

// Client is the cloud services API client
type Client interface {
	// Builder returns the build orchestrator
	Builder() Builder

	// Storage returns artifact storage
	Storage() Storage

	// Auth returns authentication service
	Auth() Auth

	// CrossCompile performs cross-platform compilation
	CrossCompile(ctx context.Context, source []byte, platforms []Platform) ([]*BuildResult, error)
}
