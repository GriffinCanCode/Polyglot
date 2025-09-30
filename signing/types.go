package signing

import (
	"context"
	"time"
)

// Certificate represents a code signing certificate
type Certificate struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"` // "apple", "windows", "linux"
	Subject     string            `json:"subject"`
	Issuer      string            `json:"issuer"`
	Serial      string            `json:"serial"`
	Fingerprint string            `json:"fingerprint"`
	NotBefore   time.Time         `json:"not_before"`
	NotAfter    time.Time         `json:"not_after"`
	KeyUsage    []string          `json:"key_usage"`
	Data        []byte            `json:"data"`
	PrivateKey  []byte            `json:"private_key"`
	Metadata    map[string]string `json:"metadata"`
}

// SignRequest represents a signing request
type SignRequest struct {
	Binary       []byte            `json:"binary"`
	Platform     string            `json:"platform"`
	Certificate  *Certificate      `json:"certificate"`
	Options      map[string]string `json:"options"`
	Timestamp    bool              `json:"timestamp"`
	Entitlements string            `json:"entitlements,omitempty"`
}

// SignResult represents a signing result
type SignResult struct {
	SignedBinary []byte            `json:"signed_binary"`
	Signature    []byte            `json:"signature"`
	Certificate  *Certificate      `json:"certificate"`
	Metadata     map[string]string `json:"metadata"`
	SignedAt     time.Time         `json:"signed_at"`
}

// VerifyRequest represents a verification request
type VerifyRequest struct {
	Binary    []byte       `json:"binary"`
	Signature []byte       `json:"signature"`
	Platform  string       `json:"platform"`
	TrustRoot *Certificate `json:"trust_root,omitempty"`
}

// VerifyResult represents a verification result
type VerifyResult struct {
	Valid       bool              `json:"valid"`
	Certificate *Certificate      `json:"certificate"`
	Chain       []*Certificate    `json:"chain"`
	Errors      []string          `json:"errors"`
	Warnings    []string          `json:"warnings"`
	Metadata    map[string]string `json:"metadata"`
	VerifiedAt  time.Time         `json:"verified_at"`
}

// Signer handles code signing operations
type Signer interface {
	// Sign signs a binary
	Sign(ctx context.Context, req *SignRequest) (*SignResult, error)

	// Verify verifies a signed binary
	Verify(ctx context.Context, req *VerifyRequest) (*VerifyResult, error)

	// GetCertificate retrieves a certificate
	GetCertificate(ctx context.Context, id string) (*Certificate, error)

	// ListCertificates lists available certificates
	ListCertificates(ctx context.Context, certType string) ([]*Certificate, error)

	// ImportCertificate imports a certificate
	ImportCertificate(ctx context.Context, cert *Certificate) error

	// DeleteCertificate removes a certificate
	DeleteCertificate(ctx context.Context, id string) error
}

// PlatformSigner handles platform-specific signing
type PlatformSigner interface {
	// Platform returns the platform identifier
	Platform() string

	// Sign signs a binary for this platform
	Sign(ctx context.Context, req *SignRequest) (*SignResult, error)

	// Verify verifies a signed binary for this platform
	Verify(ctx context.Context, req *VerifyRequest) (*VerifyResult, error)

	// Supported checks if signing is supported on current system
	Supported() bool
}
