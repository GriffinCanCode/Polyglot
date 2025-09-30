package marketplace

import (
	"context"
	"fmt"
)

// DefaultValidator implements package validation
type DefaultValidator struct{}

// NewValidator creates a new validator
func NewValidator() *DefaultValidator {
	return &DefaultValidator{}
}

// ValidatePackage validates a package
func (v *DefaultValidator) ValidatePackage(ctx context.Context, pkg *Package, data []byte) error {
	if pkg.ID == "" {
		return fmt.Errorf("package ID is required")
	}
	if pkg.Name == "" {
		return fmt.Errorf("package name is required")
	}
	if pkg.Version == "" {
		return fmt.Errorf("package version is required")
	}
	if pkg.Author == "" {
		return fmt.Errorf("package author is required")
	}
	if len(data) == 0 {
		return fmt.Errorf("package data is empty")
	}
	if pkg.Checksum == "" {
		return fmt.Errorf("package checksum is required")
	}

	return nil
}

// ValidateTemplate validates a template
func (v *DefaultValidator) ValidateTemplate(ctx context.Context, tmpl *Template) error {
	if tmpl.ID == "" {
		return fmt.Errorf("template ID is required")
	}
	if tmpl.Name == "" {
		return fmt.Errorf("template name is required")
	}
	if tmpl.Author == "" {
		return fmt.Errorf("template author is required")
	}
	if len(tmpl.Files) == 0 {
		return fmt.Errorf("template must have at least one file")
	}

	// Validate files
	for _, file := range tmpl.Files {
		if file.Path == "" {
			return fmt.Errorf("template file path is required")
		}
	}

	return nil
}

// CheckSignature verifies package signature
func (v *DefaultValidator) CheckSignature(ctx context.Context, pkg *Package, signature []byte) error {
	if len(signature) == 0 {
		return fmt.Errorf("signature is empty")
	}

	// Simplified signature validation
	// In production, this would use proper cryptographic verification
	return nil
}

// ScanSecurity performs security scanning
func (v *DefaultValidator) ScanSecurity(ctx context.Context, data []byte) ([]SecurityIssue, error) {
	issues := make([]SecurityIssue, 0)

	// Simplified security scanning
	// In production, this would perform comprehensive vulnerability checks
	if len(data) > 100*1024*1024 { // 100MB
		issues = append(issues, SecurityIssue{
			Severity:    "warning",
			Type:        "size",
			Description: "Package is very large",
		})
	}

	return issues, nil
}
