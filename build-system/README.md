# Build System

The build system orchestrates compilation across multiple languages and produces a single executable.

## Components

- **Turborepo**: Monorepo management and caching
- **esbuild**: TypeScript bundling and optimization
- **Maturin**: Python extension building
- **Cargo**: Rust compilation
- **GraalVM**: Java native compilation
- **CMake**: C/C++ project building
- **Go Build**: Go compilation
- **Zig Build**: Zig compilation

## Build Pipeline

1. **Configuration**: Read `polyglot.config.js` to determine enabled languages
2. **Dependency Resolution**: Resolve cross-language dependencies
3. **Selective Compilation**: Only compile enabled language runtimes
4. **Native Extensions**: Build language-specific native modules
5. **Bundling**: Bundle application code and assets
6. **Optimization**: Dead code elimination, tree shaking, compression
7. **Packaging**: Create single executable with embedded runtimes

## Output

- **Single Binary**: Self-contained executable
- **Size Range**: 35-80MB depending on enabled languages
- **Platform Support**: macOS universal binaries, Windows code signing, Linux AppImage/Flatpak

## Implementation Status

🚧 **Planning Phase** - Build system architecture design in progress.
