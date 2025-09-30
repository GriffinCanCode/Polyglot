package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MakefileGenerator generates Makefile for the project
type MakefileGenerator struct {
	config *ProjectConfig
}

// NewMakefileGenerator creates a new Makefile generator
func NewMakefileGenerator(config *ProjectConfig) *MakefileGenerator {
	return &MakefileGenerator{config: config}
}

// Generate creates the Makefile
func (m *MakefileGenerator) Generate() error {
	content := fmt.Sprintf(`.PHONY: all build run clean test dev install help

# Project configuration
PROJECT_NAME := %s
VERSION := %s
BUILD_DIR := dist
SRC_DIR := src/backend
BINARY := $(BUILD_DIR)/$(PROJECT_NAME)

# Go configuration
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BUILD_FLAGS := -v

# Default target
all: clean build

# Build the application
build:
	@echo "🔨 Building $(PROJECT_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY) ./$(SRC_DIR)
	@echo "✅ Build complete: $(BINARY)"

# Build for all platforms
build-all: build-darwin build-linux build-windows
	@echo "✅ Built for all platforms"

# Build for macOS
build-darwin:
	@echo "🍎 Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64 ./$(SRC_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64 ./$(SRC_DIR)

# Build for Linux
build-linux:
	@echo "🐧 Building for Linux...
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64 ./$(SRC_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-arm64 ./$(SRC_DIR)

# Build for Windows
build-windows:
	@echo "🪟 Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe ./$(SRC_DIR)

# Run the application
run: build
	@echo "🚀 Running $(PROJECT_NAME)..."
	@./$(BINARY)

# Development mode with auto-reload
dev:
	@echo "🔧 Starting development mode..."
%s

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Install dependencies
install:
	@echo "📦 Installing Go dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
%s

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
	@echo "✅ Clean complete"

# Show help
help:
	@echo "$(PROJECT_NAME) v$(VERSION) - Makefile commands"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all          - Clean and build (default)"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for all platforms"
	@echo "  build-darwin - Build for macOS"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows- Build for Windows"
	@echo "  run          - Build and run the application"
	@echo "  dev          - Start development mode"
	@echo "  test         - Run tests"
	@echo "  install      - Install dependencies"
	@echo "  clean        - Remove build artifacts"
	@echo "  help         - Show this help message"
`,
		m.config.Name,
		m.config.Version,
		m.generateDevTarget(),
		m.generateInstallTargets(),
	)

	path := filepath.Join(m.config.Name, "Makefile")
	return os.WriteFile(path, []byte(content), 0644)
}

func (m *MakefileGenerator) generateDevTarget() string {
	if contains(m.config.Features, "hmr") {
		return `	@which air > /dev/null || $(GOGET) -u github.com/cosmtrek/air
	@air`
	}
	return `	@$(GOBUILD) -o $(BINARY) ./$(SRC_DIR) && ./$(BINARY)`
}

func (m *MakefileGenerator) generateInstallTargets() string {
	var targets []string

	for _, lang := range m.config.Languages {
		switch lang {
		case "python":
			targets = append(targets, `	@if [ -f requirements.txt ]; then \
		echo "📦 Installing Python dependencies..."; \
		pip install -r requirements.txt; \
	fi`)
		case "javascript":
			pm := m.config.PackageManager
			if pm == "" {
				pm = "npm"
			}
			targets = append(targets, fmt.Sprintf(`	@if [ -f package.json ]; then \
		echo "📦 Installing JavaScript dependencies..."; \
		%s install; \
	fi`, pm))
		case "rust":
			targets = append(targets, `	@if [ -f src/rust/Cargo.toml ]; then \
		echo "📦 Installing Rust dependencies..."; \
		cd src/rust && cargo build; \
	fi`)
		}
	}

	if len(targets) > 0 {
		return "\n" + strings.Join(targets, "\n")
	}
	return ""
}

// Additional generation methods for templates.go

func (t *ProjectTemplate) generateGitignore() error {
	gitManager := NewGitManager(t.config.Name)
	return gitManager.GenerateGitignore(t.config.Languages)
}

func (t *ProjectTemplate) generateLicense() error {
	licenseManager := NewLicenseManager(t.config)
	return licenseManager.Generate()
}

func (t *ProjectTemplate) generateGoMod() error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/griffincancode/polyglot.js v0.1.0
)
`,
		t.config.Name,
	)

	path := filepath.Join(t.config.Name, "go.mod")
	return os.WriteFile(path, []byte(content), 0644)
}

func (t *ProjectTemplate) generatePackageFiles() error {
	depManager := NewDependencyManager(t.config)
	return depManager.GeneratePackageFiles()
}

func (t *ProjectTemplate) generateMakefile() error {
	makefileGen := NewMakefileGenerator(t.config)
	return makefileGen.Generate()
}
