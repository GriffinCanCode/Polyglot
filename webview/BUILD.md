# Webview Build Instructions

This document provides platform-specific instructions for building applications with native webview support.

## Overview

The polyglot webview system **includes native webview support by default**. The webview/webview library is automatically used when you build applications.

- **Default**: Native webview is enabled (no build tags needed)
- **Stub mode**: Use `-tags stub` to build with a stub backend for testing/CI without GUI dependencies

## Platform Requirements

### macOS

**System Requirements:**
- macOS 10.12 (Sierra) or later
- Xcode Command Line Tools

**Setup:**
```bash
# Install Xcode Command Line Tools if not already installed
xcode-select --install
```

**Build:**
```bash
# Build with native webview
go build -tags webview_enabled -o myapp

# The webview uses WebKit (native to macOS)
# No additional dependencies required
```

### Linux

**System Requirements:**
- GTK 3 and WebKitGTK 4.0 or 4.1

**Setup on Ubuntu/Debian:**
```bash
# Install WebKitGTK development libraries
sudo apt-get update
sudo apt-get install -y \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    pkg-config
```

**Setup on Fedora:**
```bash
# Install WebKitGTK development libraries
sudo dnf install -y \
    gtk3-devel \
    webkit2gtk4.0-devel \
    pkg-config
```

**Setup on Arch Linux:**
```bash
# Install WebKitGTK development libraries
sudo pacman -S webkit2gtk pkg-config
```

**Build:**
```bash
# Build with native webview
go build -tags webview_enabled -o myapp

# Set CGO flags if needed
CGO_ENABLED=1 go build -tags webview_enabled -o myapp
```

### Windows

**System Requirements:**
- Windows 7 or later
- WebView2 Runtime (pre-installed on Windows 11, required for older versions)

**Setup:**
1. Install WebView2 Runtime from [Microsoft's website](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)
2. Install MinGW-w64 or use TDM-GCC for CGO support

**Build:**
```bash
# Build with native webview (Windows uses Edge WebView2)
go build -tags webview_enabled -o myapp.exe

# Or with CGO explicitly enabled
CGO_ENABLED=1 go build -tags webview_enabled -o myapp.exe
```

## Cross-Compilation

Cross-compilation with CGO can be complex. Here are some tips:

### macOS to Linux
```bash
# Install cross-compilation toolchain
brew install FiloSottile/musl-cross/musl-cross

# Build for Linux
CC=x86_64-linux-musl-gcc \
CXX=x86_64-linux-musl-g++ \
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=amd64 \
go build -tags webview_enabled -o myapp-linux
```

### Using Docker for consistent builds
```bash
# Build for Linux in Docker
docker run --rm -v "$PWD":/app -w /app \
    -e CGO_ENABLED=1 \
    golang:1.21 \
    bash -c "apt-get update && \
             apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev && \
             go build -tags webview_enabled -o myapp"
```

## Testing Without Native Webview

For CI/CD pipelines or headless servers, you can build with the stub backend:

```bash
# Build with stub backend (no native dependencies)
go build -tags stub -o myapp
```

The stub backend will log webview operations to stdout instead of creating windows. This is useful for:
- CI/CD pipelines without display servers
- Automated testing
- Headless environments
- Quick testing without installing GUI dependencies

## Common Issues

### macOS: "ld: framework not found WebKit"
**Solution:** Ensure Xcode Command Line Tools are installed:
```bash
xcode-select --install
xcode-select -p  # Should show: /Library/Developer/CommandLineTools
```

### Linux: "Package webkit2gtk-4.0 was not found"
**Solution:** Install WebKitGTK development packages (see Linux setup above)

### Windows: "WebView2Loader.dll not found"
**Solution:** Install WebView2 Runtime or distribute it with your application

### CGO errors
**Solution:** Ensure CGO is enabled:
```bash
export CGO_ENABLED=1
go build -tags webview_enabled
```

## Distribution

When distributing applications:

- **macOS**: Bundle as `.app` with proper codesigning
- **Linux**: Include WebKitGTK dependencies or use static linking
- **Windows**: Include WebView2 runtime installer or use the Evergreen runtime

## Environment Variables

- `CGO_ENABLED=1`: Required for native webview (default on most platforms)
- `PKG_CONFIG_PATH`: May need to be set to find WebKitGTK on Linux
- `GODEBUG=cgocheck=0`: Can help with some CGO pointer issues (use with caution)

## Further Reading

- [webview/webview GitHub](https://github.com/webview/webview)
- [WebView2 Documentation](https://docs.microsoft.com/en-us/microsoft-edge/webview2/)
- [WebKitGTK Documentation](https://webkitgtk.org/)
