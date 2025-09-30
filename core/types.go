package core

import (
	"context"
	"time"
)

// Runtime represents a language runtime interface
type Runtime interface {
	// Initialize prepares the runtime for execution
	Initialize(ctx context.Context, config RuntimeConfig) error

	// Execute runs code in the runtime
	Execute(ctx context.Context, code string, args ...interface{}) (interface{}, error)

	// Call invokes a function by name
	Call(ctx context.Context, fn string, args ...interface{}) (interface{}, error)

	// Shutdown cleanly stops the runtime
	Shutdown(ctx context.Context) error

	// Name returns the runtime identifier
	Name() string

	// Version returns the runtime version
	Version() string
}

// RuntimeConfig holds runtime-specific configuration
type RuntimeConfig struct {
	// Name of the runtime (python, javascript, rust, etc.)
	Name string

	// Version constraint (e.g., "3.11", ">=1.70")
	Version string

	// Enabled determines if this runtime should be initialized
	Enabled bool

	// Options for runtime-specific settings
	Options map[string]interface{}

	// MaxConcurrency limits parallel executions
	MaxConcurrency int

	// Timeout for initialization
	Timeout time.Duration
}

// MemoryRegion represents shared memory accessible across runtimes
type MemoryRegion struct {
	// ID uniquely identifies this memory region
	ID string

	// Data is the underlying byte slice
	Data []byte

	// Type describes the data structure
	Type MemoryType

	// Readers tracks active readers for sync
	Readers int

	// Writers tracks active writers for sync
	Writers int
}

// MemoryType describes the structure of shared memory
type MemoryType string

const (
	TypeBytes   MemoryType = "bytes"
	TypeInt32   MemoryType = "int32"
	TypeInt64   MemoryType = "int64"
	TypeFloat32 MemoryType = "float32"
	TypeFloat64 MemoryType = "float64"
	TypeString  MemoryType = "string"
	TypeStruct  MemoryType = "struct"
)

// Bridge connects the frontend webview to backend runtimes
type Bridge interface {
	// Register adds a callable function to the bridge
	Register(name string, fn BridgeFunc) error

	// Unregister removes a callable function
	Unregister(name string) error

	// Call invokes a registered function
	Call(ctx context.Context, name string, args ...interface{}) (interface{}, error)
}

// BridgeFunc is a function callable from the frontend
type BridgeFunc func(ctx context.Context, args ...interface{}) (interface{}, error)

// Message represents cross-language communication
type Message struct {
	// ID for request-response correlation
	ID string

	// Type of message (call, response, error, event)
	Type MessageType

	// Target runtime or function
	Target string

	// Payload contains the message data
	Payload interface{}

	// Error contains error information if applicable
	Error error

	// Timestamp of message creation
	Timestamp time.Time
}

// MessageType categorizes messages
type MessageType string

const (
	TypeCall     MessageType = "call"
	TypeResponse MessageType = "response"
	TypeError    MessageType = "error"
	TypeEvent    MessageType = "event"
)

// BuildConfig defines compilation settings
type BuildConfig struct {
	// OutputPath for the final binary
	OutputPath string

	// Platform target (darwin, linux, windows)
	Platform string

	// Arch target (amd64, arm64)
	Arch string

	// Optimize enables optimizations
	Optimize bool

	// Compress applies UPX compression
	Compress bool

	// Tags are build tags to apply
	Tags []string
}
