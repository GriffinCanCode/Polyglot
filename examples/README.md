# Polyglot Framework Examples

This directory contains working examples demonstrating the framework's capabilities.

## Available Examples

### 01-hello-world ✅ VALIDATED

A comprehensive example demonstrating:
- Multi-language orchestration (Go + JavaScript + Python)
- Cross-language function calls
- Shared memory coordination
- Error handling and graceful degradation

**Status**: Fully tested and working  
**Performance**: 333k JavaScript calls/sec, 6M memory ops/sec  
**See**: [examples/01-hello-world/](./01-hello-world/)

## Quick Start

```bash
# Navigate to an example
cd 01-hello-world

# Build and run
go build -o dist/hello-world ./src/backend
./dist/hello-world

# Run tests
go test -v

# Run benchmarks
go test -bench=. -benchtime=3s
```

## Example Status

| Example | Status | Runtimes | Tests | Performance |
|---------|--------|----------|-------|-------------|
| 01-hello-world | ✅ Working | Go, JS, Python (stub) | ✅ Passing | ⚡ Validated |
| 02-ml-dashboard | 📋 Planned | Go, Python, React | - | - |
| 03-game-engine | 📋 Planned | Go, C++, Lua | - | - |

## Coming Soon

- **ML Dashboard**: Real-time data science with Python + React
- **Game Engine**: High-performance graphics with C++ + Lua scripting
- **File Converter**: Rust processing + TypeScript UI
- **API Gateway**: Multi-language microservices coordination

## Contributing Examples

Want to add an example? Follow this structure:

```
examples/XX-example-name/
├── README.md              # Setup and usage
├── VALIDATION.md         # Test results and performance
├── go.mod                # Go module config
├── main_test.go          # Unit tests
├── src/
│   ├── backend/          # Go orchestration
│   │   └── main.go
│   └── frontend/         # UI (if applicable)
│       └── index.html
└── dist/                 # Build output
```

See [01-hello-world](./01-hello-world/) as a reference.
