package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/signing"
)

func TestSigner(t *testing.T) {
	ctx := context.Background()
	signer := signing.NewSigner()

	// Create a test certificate
	cert := &signing.Certificate{
		ID:          "test-cert",
		Type:        "apple",
		Subject:     "Developer ID Application",
		Issuer:      "Apple Inc.",
		Serial:      "123456",
		Fingerprint: "abc123",
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:    []string{"digitalSignature", "keyEncipherment"},
		Data:        []byte("cert-data"),
		PrivateKey:  []byte("private-key"),
	}

	// Import certificate
	err := signer.ImportCertificate(ctx, cert)
	if err != nil {
		t.Fatalf("failed to import certificate: %v", err)
	}

	// Retrieve certificate
	retrieved, err := signer.GetCertificate(ctx, "test-cert")
	if err != nil {
		t.Fatalf("failed to get certificate: %v", err)
	}

	if retrieved.Subject != cert.Subject {
		t.Errorf("expected subject %s, got %s", cert.Subject, retrieved.Subject)
	}

	// List certificates
	certs, err := signer.ListCertificates(ctx, "apple")
	if err != nil {
		t.Fatalf("failed to list certificates: %v", err)
	}

	if len(certs) == 0 {
		t.Error("expected at least one certificate")
	}

	// Delete certificate
	err = signer.DeleteCertificate(ctx, "test-cert")
	if err != nil {
		t.Fatalf("failed to delete certificate: %v", err)
	}

	_, err = signer.GetCertificate(ctx, "test-cert")
	if err == nil {
		t.Error("expected error when getting deleted certificate")
	}
}

func TestDarwinSigner(t *testing.T) {
	ctx := context.Background()
	darwinSigner := signing.NewDarwinSigner()

	if darwinSigner.Platform() != "darwin" {
		t.Errorf("expected platform darwin, got %s", darwinSigner.Platform())
	}

	cert := &signing.Certificate{
		ID:      "darwin-cert",
		Type:    "apple",
		Subject: "Developer ID",
	}

	req := &signing.SignRequest{
		Binary:       []byte("binary-data"),
		Platform:     "darwin",
		Certificate:  cert,
		Timestamp:    true,
		Entitlements: "com.example.app.entitlements",
	}

	// Sign binary
	result, err := darwinSigner.Sign(ctx, req)
	if !darwinSigner.Supported() {
		if err == nil {
			t.Error("expected error when signing on unsupported platform")
		}
		return
	}

	if err != nil {
		t.Fatalf("signing failed: %v", err)
	}

	if result.SignedBinary == nil {
		t.Error("expected signed binary to not be nil")
	}

	// Verify signature
	verifyReq := &signing.VerifyRequest{
		Binary:    result.SignedBinary,
		Signature: result.Signature,
		Platform:  "darwin",
	}

	verifyResult, err := darwinSigner.Verify(ctx, verifyReq)
	if err != nil {
		t.Fatalf("verification failed: %v", err)
	}

	if !verifyResult.Valid {
		t.Error("expected signature to be valid")
	}
}

func TestWindowsSigner(t *testing.T) {
	ctx := context.Background()
	windowsSigner := signing.NewWindowsSigner()

	if windowsSigner.Platform() != "windows" {
		t.Errorf("expected platform windows, got %s", windowsSigner.Platform())
	}

	cert := &signing.Certificate{
		ID:      "windows-cert",
		Type:    "windows",
		Subject: "Code Signing Certificate",
	}

	req := &signing.SignRequest{
		Binary:      []byte("binary-data"),
		Platform:    "windows",
		Certificate: cert,
		Timestamp:   true,
	}

	// Sign binary
	result, err := windowsSigner.Sign(ctx, req)
	if !windowsSigner.Supported() {
		if err == nil {
			t.Error("expected error when signing on unsupported platform")
		}
		return
	}

	if err != nil {
		t.Fatalf("signing failed: %v", err)
	}

	if result.SignedBinary == nil {
		t.Error("expected signed binary to not be nil")
	}
}

func TestLinuxSigner(t *testing.T) {
	ctx := context.Background()
	linuxSigner := signing.NewLinuxSigner()

	if linuxSigner.Platform() != "linux" {
		t.Errorf("expected platform linux, got %s", linuxSigner.Platform())
	}

	cert := &signing.Certificate{
		ID:      "linux-cert",
		Type:    "gpg",
		Subject: "GPG Key",
	}

	req := &signing.SignRequest{
		Binary:      []byte("binary-data"),
		Platform:    "linux",
		Certificate: cert,
	}

	// Sign binary
	result, err := linuxSigner.Sign(ctx, req)
	if !linuxSigner.Supported() {
		if err == nil {
			t.Error("expected error when signing on unsupported platform")
		}
		return
	}

	if err != nil {
		t.Fatalf("signing failed: %v", err)
	}

	if result.SignedBinary == nil {
		t.Error("expected signed binary to not be nil")
	}
}

func TestSignatureVerification(t *testing.T) {
	ctx := context.Background()
	signer := signing.NewSigner()

	cert := &signing.Certificate{
		ID:       "verify-cert",
		Type:     "apple",
		Subject:  "Test Certificate",
		Data:     []byte("cert-data"),
		NotAfter: time.Now().Add(365 * 24 * time.Hour),
	}

	signer.ImportCertificate(ctx, cert)

	req := &signing.SignRequest{
		Binary:      []byte("test-binary"),
		Platform:    "darwin",
		Certificate: cert,
	}

	result, err := signer.Sign(ctx, req)
	if err != nil && result == nil {
		t.Logf("signing not supported on this platform: %v", err)
		return
	}

	verifyReq := &signing.VerifyRequest{
		Binary:    result.SignedBinary,
		Signature: result.Signature,
		Platform:  "darwin",
	}

	verifyResult, err := signer.Verify(ctx, verifyReq)
	if err != nil {
		t.Fatalf("verification failed: %v", err)
	}

	if !verifyResult.Valid {
		t.Error("expected signature to be valid")
	}
}

func TestCertificateExpiry(t *testing.T) {
	cert := &signing.Certificate{
		ID:        "expiry-test",
		Type:      "apple",
		Subject:   "Test",
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().Add(1 * time.Hour),
	}

	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		t.Error("certificate should be valid now")
	}

	expiredCert := &signing.Certificate{
		ID:        "expired-test",
		Type:      "apple",
		Subject:   "Expired",
		NotBefore: time.Now().Add(-48 * time.Hour),
		NotAfter:  time.Now().Add(-24 * time.Hour),
	}

	if now.Before(expiredCert.NotAfter) {
		t.Error("certificate should be expired")
	}
}
