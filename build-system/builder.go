package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/polyglot-framework/polyglot/core"
)

// Builder handles application compilation
type Builder struct {
	config core.BuildConfig
}

// New creates a new builder
func New(config core.BuildConfig) *Builder {
	return &Builder{config: config}
}

// Build compiles the application
func (b *Builder) Build() error {
	// Create output directory
	if err := os.MkdirAll(b.config.OutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build Go binary
	if err := b.buildGo(); err != nil {
		return fmt.Errorf("failed to build Go: %w", err)
	}

	// Build frontend if exists
	if err := b.buildFrontend(); err != nil {
		return fmt.Errorf("failed to build frontend: %w", err)
	}

	// Compress if requested
	if b.config.Compress {
		if err := b.compress(); err != nil {
			return fmt.Errorf("failed to compress: %w", err)
		}
	}

	return nil
}

// buildGo compiles the Go application
func (b *Builder) buildGo() error {
	args := []string{"build"}

	// Add optimization flags
	if b.config.Optimize {
		args = append(args, "-ldflags", "-s -w")
	}

	// Add build tags
	if len(b.config.Tags) > 0 {
		args = append(args, "-tags", strings.Join(b.config.Tags, ","))
	}

	// Set output path
	outputName := "app"
	if b.config.Platform == "windows" {
		outputName += ".exe"
	}
	outputPath := filepath.Join(b.config.OutputPath, outputName)
	args = append(args, "-o", outputPath)

	// Add source path
	args = append(args, "./src/backend")

	// Set environment variables
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		"GOOS="+b.config.Platform,
		"GOARCH="+b.config.Arch,
		"CGO_ENABLED=1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// buildFrontend bundles frontend assets
func (b *Builder) buildFrontend() error {
	frontendDir := "src/frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		return nil // No frontend to build
	}

	// Copy frontend files to output
	outputFrontend := filepath.Join(b.config.OutputPath, "frontend")
	if err := os.MkdirAll(outputFrontend, 0755); err != nil {
		return err
	}

	return copyDir(frontendDir, outputFrontend)
}

// compress applies UPX compression
func (b *Builder) compress() error {
	outputName := "app"
	if b.config.Platform == "windows" {
		outputName += ".exe"
	}
	outputPath := filepath.Join(b.config.OutputPath, outputName)

	cmd := exec.Command("upx", "--best", "--lzma", outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create destination path
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}
