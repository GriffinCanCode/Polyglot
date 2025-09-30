# Webview Integration

The webview wrapper provides the native display surface for Polyglot applications.

## Technology

- **Library**: Wry (from Tauri project)
- **Integration**: Wrapped as Bun native module
- **Size**: ~15MB vs Electron's 100MB+ Chromium
- **Communication**: High-performance IPC bridge

## Features

- **Native Performance**: Direct integration with Bun runtime
- **Cross-platform**: macOS, Windows, Linux support
- **Modern APIs**: Web standards with native capabilities
- **Security**: Sandboxed execution environment
- **Customization**: Configurable window properties and behaviors

## Architecture

The webview communicates with the Bun runtime through a high-performance IPC bridge that makes cross-language function calls feel like direct calls to developers. This eliminates the serialization overhead of traditional IPC while maintaining security boundaries.

## Implementation Status

🚧 **Planning Phase** - Webview integration architecture design in progress.
