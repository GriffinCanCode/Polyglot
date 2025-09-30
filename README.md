# Polyglot Framework

A radical rethinking of desktop application development that treats multilingual programming as a first-class citizen.

## Project Structure

```
/polyglot/
├── /core/                    # Core Bun runtime integration
├── /runtimes/               # Language runtime integrations
│   ├── /python/             # PyO3 + pyo3-napi integration
│   ├── /rust/               # Neon bindings
│   ├── /go/                 # CGO + N-API wrapper
│   ├── /java/               # GraalVM integration
│   ├── /php/                # PHP interpreter embedding
│   ├── /cpp/                # node-addon-api bindings
│   └── /zig/                # C ABI compilation
├── /cli/                    # Rust CLI tool
├── /build-system/           # Build orchestration (Turborepo, esbuild, etc.)
├── /webview/                # Wry webview wrapper
├── /types/                  # TypeScript definition generation
├── /memory/                 # Unified memory management
├── /security/               # Sandboxing and code signing
├── /examples/               # Example applications
├── /docs/                   # Documentation
└── /tests/                  # Cross-language test suite
```

## Getting Started

1. Configure languages in `polyglot.config.js`
2. Run `polyglot init` to set up a new project
3. Use `polyglot build` to compile your polyglot application

## Development Status

🚧 **Under Development** - This is the initial project structure setup.

## Core Features

- **Multi-language Runtime**: Embed Python, Rust, Go, Java, PHP, C++, and Zig in a single Bun process
- **Zero-copy Memory Sharing**: Shared ArrayBuffers across language boundaries
- **Type-safe Interop**: Automatic TypeScript definition generation
- **Native Performance**: Microsecond-latency function calls between languages
- **Single Binary Output**: 35-80MB applications vs Electron's 150MB+

## License

MIT
