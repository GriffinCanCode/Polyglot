package core

import (
	"fmt"
	"time"
)

// Config is the main application configuration
type Config struct {
	// App contains application metadata
	App AppConfig

	// Languages maps runtime names to their configurations
	Languages map[string]*RuntimeConfig

	// Memory configures the memory coordinator
	Memory MemoryConfig

	// Webview configures the frontend
	Webview WebviewConfig

	// Build configures compilation
	Build BuildConfig
}

// AppConfig holds application metadata
type AppConfig struct {
	// Name of the application
	Name string

	// Version of the application
	Version string

	// Description of the application
	Description string

	// Author information
	Author string

	// License identifier
	License string
}

// MemoryConfig configures memory management
type MemoryConfig struct {
	// MaxSharedMemory in bytes
	MaxSharedMemory int64

	// EnableZeroCopy enables zero-copy optimizations
	EnableZeroCopy bool

	// GCInterval for memory cleanup
	GCInterval time.Duration
}

// WebviewConfig configures the frontend webview
type WebviewConfig struct {
	// Title of the window
	Title string

	// Width in pixels
	Width int

	// Height in pixels
	Height int

	// Resizable allows window resizing
	Resizable bool

	// Debug enables devtools
	Debug bool

	// URL to load
	URL string
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "polyglot-app",
			Version: "0.1.0",
			License: "MIT",
		},
		Languages: map[string]*RuntimeConfig{},
		Memory: MemoryConfig{
			MaxSharedMemory: 1024 * 1024 * 1024, // 1GB
			EnableZeroCopy:  true,
			GCInterval:      time.Minute * 5,
		},
		Webview: WebviewConfig{
			Title:     "Polyglot Application",
			Width:     1280,
			Height:    720,
			Resizable: true,
			Debug:     false,
			URL:       "http://localhost:3000",
		},
		Build: BuildConfig{
			OutputPath: "./dist",
			Optimize:   true,
			Compress:   false,
		},
	}
}

// Validate checks configuration validity
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if c.Memory.MaxSharedMemory <= 0 {
		return fmt.Errorf("max shared memory must be positive")
	}

	if c.Webview.Width <= 0 || c.Webview.Height <= 0 {
		return fmt.Errorf("webview dimensions must be positive")
	}

	return nil
}

// EnableRuntime enables a language runtime with default settings
func (c *Config) EnableRuntime(name string, version string) {
	c.Languages[name] = &RuntimeConfig{
		Name:           name,
		Version:        version,
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 10,
		Timeout:        time.Second * 30,
	}
}

// DisableRuntime disables a language runtime
func (c *Config) DisableRuntime(name string) {
	if cfg, exists := c.Languages[name]; exists {
		cfg.Enabled = false
	}
}

// IsRuntimeEnabled checks if a runtime is enabled
func (c *Config) IsRuntimeEnabled(name string) bool {
	if cfg, exists := c.Languages[name]; exists {
		return cfg.Enabled
	}
	return false
}
