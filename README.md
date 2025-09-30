# Polyglot Framework

> A radical rethinking of desktop application development that treats multilingual programming as a first-class citizen.

## Overview

Polyglot enables developers to use the absolute best language for each component while maintaining a cohesive, type-safe development experience. Built with Go as the orchestrator, it embeds language runtimes directly into a single process with true parallel execution via goroutines.

## Architecture

```
polyglot/
├── core/           # Core orchestrator, config, types, memory, profiler, HMR
├── runtimes/       # Language runtime integrations
│   ├── python/     # Python runtime with CGO bindings (Phase 1)
│   ├── javascript/ # JavaScript/TypeScript runtime with V8 (Phase 1)
│   ├── rust/       # Rust integration with shared library loading (Phase 2)
│   ├── java/       # Java integration with JNI bindings (Phase 2)
│   ├── cpp/        # C++ integration with CGO bindings (Phase 2)
│   ├── php/        # PHP integration with embedded interpreter (Phase 3)
│   ├── ruby/       # Ruby integration with CGO bindings (Phase 3)
│   ├── lua/        # Lua integration with state management (Phase 3)
│   ├── zig/        # Zig integration with C ABI (Phase 3)
│   └── wasm/       # WebAssembly runtime (Phase 3)
├── webview/        # Native webview integration
├── build-system/   # Build tooling, selective compilation, and binding generation
├── cli/            # CLI tool for project management
├── marketplace/    # Package registry and template management (Phase 4)
├── cloud/          # Cloud build infrastructure (Phase 4)
├── signing/        # Code signing for all platforms (Phase 4)
├── updates/        # Differential update system (Phase 4)
├── tests/          # Comprehensive test suite
├── types/          # Shared type definitions
├── security/       # Security and sandboxing
└── examples/       # Example applications

```

## Features

### Phase 1 (✅ Complete)

- **Core Orchestrator**: Go-based runtime coordinator with goroutine pooling
- **Python Integration**: Full CPython embedding with CGO bindings
- **JavaScript/TypeScript**: V8 runtime integration
- **Memory Coordinator**: Zero-copy shared memory architecture
- **Webview**: Native webview with bidirectional bridge
- **CLI Tool**: Project initialization, build, dev, and test commands
- **Build System**: Selective compilation with build tags
- **Test Suite**: Comprehensive unit and integration tests

### Phase 2 (✅ Complete)

- **Rust Integration**: Shared library loading with dlopen/dlsym
- **Java Integration**: JNI bindings with JVM management
- **C++ Integration**: Direct CGO bindings with dynamic loading
- **Binding Generator**: Automatic type definitions for TypeScript/Python/Rust
- **Profiler**: Cross-runtime performance tracking with detailed metrics
- **Hot Module Replacement**: File watching with runtime-specific reload handlers

### Phase 3 (✅ Complete)

- **PHP Integration**: Embedded PHP interpreter with SAPI
- **Ruby Integration**: CGO bindings to libruby with worker pools
- **Lua Integration**: Lightweight Lua state management
- **Zig Integration**: C ABI compatibility with dynamic loading
- **WASM Runtime**: WebAssembly bytecode execution engine
- **Security Sandboxing**: Platform-specific enforcers (Landlock, App Sandbox, AppContainer)

### Phase 4 (✅ Complete)

- **Marketplace**: Package registry with search, caching, and validation
- **Cloud Services**: Remote build infrastructure with authentication
- **Cross-Platform Compilation**: Parallel builds for multiple platforms
- **Code Signing**: Platform-specific signing (macOS, Windows, Linux)
- **Update System**: Differential patching with rollback support

### Phase 5 (Planned)

- Mobile runtime exploration (iOS/Android)
- Embedded systems support
- Plugin architecture for custom runtimes

## Quick Start

### Installation

```bash
go install github.com/griffincancode/polyglot.js/cli@latest
```

### Create a New Project

```bash
polyglot init myapp
cd myapp
polyglot dev
```

### Project Structure

```
myapp/
├── src/
│   ├── backend/
│   │   └── main.go      # Go orchestrator
│   └── frontend/
│       └── index.html   # Frontend UI
└── dist/                # Build output
```

## Configuration

Configure your application in `main.go`:

```go
config := core.DefaultConfig()
config.App.Name = "myapp"
config.EnableRuntime("python", "3.11")
config.EnableRuntime("javascript", "latest")
```

## Performance

- **Startup**: Sub-10ms with multiple runtimes
- **Memory**: ~30MB minimum (Python-only), ~70MB full-stack
- **Inter-language calls**: 0.05-0.5 microseconds
- **True parallelism**: Genuine multi-core utilization via goroutines

## Development

### Build from Source

```bash
git clone https://github.com/griffincancode/polyglot.js.git
cd polyglot.js
go build -o polyglot ./cli
```

### Run Tests

```bash
go test ./...
```

### Build Examples

```bash
cd examples/01-hello-world
go build -o dist/hello-world ./src/backend
./dist/hello-world
```

Or run tests:
```bash
cd examples/01-hello-world
go test -v
```

## Design Principles

- **Extensible**: Modular runtime architecture
- **Testable**: Comprehensive test coverage with mocks
- **Compact**: One-word file names, short functions
- **Type-safe**: Strong typing throughout
- **Zero-debt**: Clean architecture, readable code

## License

MIT License - See LICENSE file for details.

## Contributing

See CONTRIBUTING.md for guidelines.

## Links

- [Documentation](https://polyglot.dev/docs)
- [Examples](./examples)
- [Plan](./plan.md)
