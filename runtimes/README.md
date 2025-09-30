# Language Runtime Integrations

This directory contains the individual language runtime integrations for Polyglot.

## Supported Languages

### Python (`/python/`)
- **Integration**: PyO3 + pyo3-napi
- **Runtime**: Full CPython interpreter embedded in Bun process
- **Ecosystem**: NumPy, PyTorch, scikit-learn, pandas with full C extension support
- **Size**: ~20MB

### Rust (`/rust/`)
- **Integration**: Neon bindings
- **Compilation**: Direct to native Node modules
- **Features**: Zero-cost abstractions, memory safety, Tokio async support
- **Size**: ~5MB (compiled to native)

### Go (`/go/`)
- **Integration**: CGO + thin C wrapper + N-API
- **Features**: Goroutine support, channels, concurrency model
- **Size**: ~8MB (compiled to native)

### Java (`/java/`)
- **Integration**: GraalVM Native Image
- **Options**: Native compilation (~15MB) or embedded JVM (~50MB)
- **Compatibility**: Most Java libraries supported

### PHP (`/php/`)
- **Integration**: Embedded PHP interpreter
- **Features**: Most PHP features and extensions
- **Size**: ~10MB

### C/C++ (`/cpp/`)
- **Integration**: node-addon-api
- **Use Cases**: Graphics programming, game engines, existing C++ codebases
- **Libraries**: OpenCV, FFmpeg, custom rendering engines
- **Size**: ~3MB (compiled to native)

### Zig (`/zig/`)
- **Integration**: Compiles to C ABI with perfect C interop
- **Features**: Compile-time code execution, zero-overhead abstractions
- **Size**: ~3MB (compiled to native)

## Implementation Status

🚧 **Planning Phase** - Language integration architecture design in progress.
