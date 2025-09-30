package signing

import (
	"context"
	"fmt"
	"sync"
)

// DefaultSigner implements the main signing orchestrator
type DefaultSigner struct {
	mu      sync.RWMutex
	certs   map[string]*Certificate
	signers map[string]PlatformSigner
}

// NewSigner creates a new signer
func NewSigner() *DefaultSigner {
	s := &DefaultSigner{
		certs:   make(map[string]*Certificate),
		signers: make(map[string]PlatformSigner),
	}

	// Register platform signers
	s.RegisterPlatformSigner(NewDarwinSigner())
	s.RegisterPlatformSigner(NewWindowsSigner())
	s.RegisterPlatformSigner(NewLinuxSigner())

	return s
}

// RegisterPlatformSigner registers a platform-specific signer
func (s *DefaultSigner) RegisterPlatformSigner(ps PlatformSigner) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signers[ps.Platform()] = ps
}

// Sign signs a binary
func (s *DefaultSigner) Sign(ctx context.Context, req *SignRequest) (*SignResult, error) {
	if len(req.Binary) == 0 {
		return nil, fmt.Errorf("binary is required")
	}
	if req.Certificate == nil {
		return nil, fmt.Errorf("certificate is required")
	}

	s.mu.RLock()
	signer, ok := s.signers[req.Platform]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", req.Platform)
	}

	if !signer.Supported() {
		return nil, fmt.Errorf("signing not supported on this system for platform: %s", req.Platform)
	}

	return signer.Sign(ctx, req)
}

// Verify verifies a signed binary
func (s *DefaultSigner) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResult, error) {
	if len(req.Binary) == 0 {
		return nil, fmt.Errorf("binary is required")
	}

	s.mu.RLock()
	signer, ok := s.signers[req.Platform]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", req.Platform)
	}

	return signer.Verify(ctx, req)
}

// GetCertificate retrieves a certificate
func (s *DefaultSigner) GetCertificate(ctx context.Context, id string) (*Certificate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cert, ok := s.certs[id]
	if !ok {
		return nil, fmt.Errorf("certificate not found: %s", id)
	}

	return cert, nil
}

// ListCertificates lists available certificates
func (s *DefaultSigner) ListCertificates(ctx context.Context, certType string) ([]*Certificate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var certs []*Certificate
	for _, cert := range s.certs {
		if certType == "" || cert.Type == certType {
			certs = append(certs, cert)
		}
	}

	return certs, nil
}

// ImportCertificate imports a certificate
func (s *DefaultSigner) ImportCertificate(ctx context.Context, cert *Certificate) error {
	if cert.ID == "" {
		return fmt.Errorf("certificate ID is required")
	}
	if cert.Type == "" {
		return fmt.Errorf("certificate type is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.certs[cert.ID] = cert
	return nil
}

// DeleteCertificate removes a certificate
func (s *DefaultSigner) DeleteCertificate(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.certs[id]; !ok {
		return fmt.Errorf("certificate not found: %s", id)
	}

	delete(s.certs, id)
	return nil
}
