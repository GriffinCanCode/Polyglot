# Native Webview Integration

This package provides native desktop window functionality for Polyglot applications using the [webview/webview](https://github.com/webview/webview) library.

## Overview

The webview package enables you to create cross-platform desktop applications with web-based UIs. It uses native webview components on each platform:

- **macOS**: WebKit (Cocoa/WebKit API)
- **Linux**: WebKitGTK
- **Windows**: Microsoft Edge WebView2

## Features

✅ **Cross-platform**: Works on macOS, Linux, and Windows  
✅ **Native performance**: Uses platform-native webview components  
✅ **Bidirectional communication**: JavaScript ↔ Go bridge  
✅ **Zero dependencies**: No external browser or runtime needed  
✅ **Development tools**: Built-in DevTools support  
✅ **Flexible**: Load local HTML, remote URLs, or embedded content  
✅ **Type-safe**: Strong typing with Go structs and interfaces  
✅ **Testing support**: Stub backend for CI/CD and testing  

## Quick Start

### 1. Install Platform Dependencies

**macOS:**
```bash
# No dependencies needed - uses native WebKit
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
```

**Windows:**
- Install [Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)

### 2. Create Your Application

```go
package main

import (
    "log"
    "github.com/griffincancode/polyglot.js/core"
    "github.com/griffincancode/polyglot.js/webview"
)

func main() {
    // Configure the window
    config := core.WebviewConfig{
        Title:     "My App",
        Width:     1280,
        Height:    720,
        Resizable: true,
        Debug:     true, // Enable DevTools
        URL:       "https://example.com",
    }

    // Create webview (bridge is optional)
    wv := webview.New(config, nil)
    
    // Initialize and run
    if err := wv.Initialize(); err != nil {
        log.Fatal(err)
    }
    defer wv.Terminate()

    // Run blocks until window is closed
    if err := wv.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### 3. Build and Run

```bash
# Build (native webview is enabled by default)
go build -o myapp

# Run
./myapp
```

## Bridge Communication

The bridge enables seamless communication between JavaScript and Go:

### Go Side: Implement Bridge Interface

```go
type MyBridge struct {
    counter int
}

func (b *MyBridge) Call(_ interface{}, name string, args ...interface{}) (interface{}, error) {
    switch name {
    case "increment":
        b.counter++
        return b.counter, nil
    
    case "greet":
        name := args[0].(string)
        return fmt.Sprintf("Hello, %s!", name), nil
    
    default:
        return nil, fmt.Errorf("unknown function: %s", name)
    }
}
```

### JavaScript Side: Call Go Functions

The bridge automatically injects a `window.polyglot` object:

```javascript
// Call Go functions from JavaScript
async function increment() {
    const result = await window.polyglot.call('increment');
    console.log('Counter:', result);
}

async function greetUser() {
    const greeting = await window.polyglot.call('greet', 'Alice');
    console.log(greeting); // "Hello, Alice!"
}
```

## Architecture

### Component Structure

```
┌─────────────────────────────────────────────┐
│           Application Code                   │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          Webview (webview.go)                │
│  - Window management                         │
│  - Bridge coordination                       │
│  - Lifecycle management                      │
└──────────────────┬──────────────────────────┘
                   │
       ┌───────────┴───────────┐
       │                       │
┌──────▼──────┐       ┌───────▼────────┐
│NativeBackend│       │  StubBackend   │
│(native.go)  │       │   (stub.go)    │
│             │       │                │
│Uses webview/│       │For testing/CI  │
│webview lib  │       │                │
└─────────────┘       └────────────────┘
```

### Build Tags

- **Default**: Native webview is enabled (uses `native.go`)
- **Stub mode**: Use `-tags stub` to enable stub backend (uses `stub.go`)

The build tags use a negative condition pattern:
- `native.go`: `//go:build !stub` (default when no tags specified)
- `stub.go`: `//go:build stub` (only with `-tags stub`)

## API Reference

### Webview

```go
type Webview struct {
    // ...
}

// Create new webview
func New(config core.WebviewConfig, bridge core.Bridge) *Webview

// Initialize creates the window
func (w *Webview) Initialize() error

// Run starts the event loop (blocks)
func (w *Webview) Run() error

// Eval executes JavaScript
func (w *Webview) Eval(script string) error

// Bind adds a Go function callable from JavaScript
func (w *Webview) Bind(name string, fn interface{}) error

// Terminate closes the window
func (w *Webview) Terminate() error
```

### WebviewConfig

```go
type WebviewConfig struct {
    Title     string  // Window title
    Width     int     // Width in pixels
    Height    int     // Height in pixels
    Resizable bool    // Allow window resizing
    Debug     bool    // Enable DevTools
    URL       string  // URL to load
}
```

### Bridge Interface

```go
type Bridge interface {
    Call(context interface{}, name string, args ...interface{}) (interface{}, error)
}
```

## Examples

### Example 1: Simple Window

```go
config := core.WebviewConfig{
    Title:  "Simple App",
    Width:  800,
    Height: 600,
    URL:    "data:text/html,<h1>Hello World</h1>",
}

wv := webview.New(config, nil)
wv.Initialize()
defer wv.Terminate()
wv.Run()
```

### Example 2: With Bridge

See the complete example in `examples/02-webview-demo/` which includes:
- Todo list management
- State synchronization
- Complex data structures
- Error handling
- Modern UI

### Example 3: Loading Local Files

```go
config := core.WebviewConfig{
    Title: "Local App",
    Width: 1024,
    Height: 768,
    URL:   "file:///path/to/your/index.html",
}
```

### Example 4: Multiple Windows

```go
// Main window
mainWV := webview.New(mainConfig, bridge)
mainWV.Initialize()

// Settings window
settingsWV := webview.New(settingsConfig, bridge)
settingsWV.Initialize()

// Run in separate goroutines
go mainWV.Run()
settingsWV.Run()
```

## Testing

### Unit Tests

```bash
# Run webview tests
make test-webview

# Or directly
go test ./tests/webview_test.go
```

### CI/CD Integration

For headless CI environments, use the stub backend:

```bash
# Build with stub (no GUI dependencies)
go build -tags stub -o myapp

# Tests with stub
go test -tags stub ./...
```

### GitHub Actions Example

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      # Test with stub backend (no GUI)
      - name: Run tests
        run: go test -tags stub ./tests/webview_test.go
```

## Platform-Specific Notes

### macOS

- Uses native WebKit framework
- No external dependencies
- Supports arm64 and amd64
- DevTools available via right-click context menu

### Linux

- Requires WebKitGTK installation
- GTK+ 3.0 or later
- Works on X11 and Wayland
- DevTools available via F12 or right-click

### Windows

- Requires WebView2 Runtime
- Pre-installed on Windows 11
- For Windows 10, distribute runtime with app
- DevTools available via F12

## Deployment

### macOS App Bundle

```bash
# Create .app structure
mkdir -p MyApp.app/Contents/MacOS
mkdir -p MyApp.app/Contents/Resources

# Copy binary
cp myapp MyApp.app/Contents/MacOS/

# Create Info.plist
# ... (see BUILD.md for details)

# Codesign
codesign --force --deep --sign - MyApp.app
```

### Linux Distribution

```bash
# Static build (if possible)
CGO_ENABLED=1 go build -o myapp

# Or provide installation script for dependencies
sudo apt-get install libwebkit2gtk-4.0-37
```

### Windows Distribution

Include WebView2 runtime installer or use the Evergreen runtime.

## Troubleshooting

### Build Issues

**Error: "package github.com/webview/webview: no Go files"**
- The library uses CGO. Ensure `CGO_ENABLED=1`

**Error: "ld: framework not found WebKit" (macOS)**
- Install Xcode Command Line Tools: `xcode-select --install`

**Error: "Package webkit2gtk-4.0 was not found" (Linux)**
- Install dependencies: `sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev`

### Runtime Issues

**Window doesn't appear**
- Ensure you're calling `Run()` on the main thread
- Check that `Initialize()` succeeded
- Try with `Debug: true` to see console errors

**JavaScript errors**
- Open DevTools (F12 when `Debug: true`)
- Check browser console for error messages
- Verify the bridge is injected: `console.log(window.polyglot)`

**Bridge calls fail**
- Ensure bridge is provided to `New()`
- Check function names match exactly (case-sensitive)
- Verify argument types are JSON-serializable

## Best Practices

1. **Error Handling**: Always check errors from `Initialize()` and `Run()`
2. **Resource Cleanup**: Use `defer wv.Terminate()` after `Initialize()`
3. **Thread Safety**: Bridge calls may come from different threads - use mutexes
4. **Data Serialization**: Bridge arguments are JSON-serialized - keep them simple
5. **Debug Mode**: Enable for development, disable for production
6. **Testing**: Use stub backend for automated tests

## Performance Tips

- **Minimize Bridge Calls**: Batch operations when possible
- **Async Operations**: Use goroutines for long-running operations
- **Memory Management**: Clean up resources in `Terminate()`
- **HTML Optimization**: Minify assets, lazy load content
- **Native Operations**: Move heavy computation to Go side

## Security Considerations

- **Content Security**: Validate and sanitize user input
- **HTTPS**: Use HTTPS for remote content
- **Bridge Exposure**: Only expose necessary functions
- **Input Validation**: Validate all data from JavaScript
- **Update Runtime**: Keep WebView2/WebKitGTK updated

## Further Reading

- [webview/webview Documentation](https://github.com/webview/webview)
- [Build Instructions](BUILD.md)
- [Example Application](../examples/02-webview-demo/)
- [Test Suite](../tests/webview_test.go)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

This package is part of the Polyglot project and follows the same license.
