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

### 02-webview-demo ✅ WORKING

An interactive webview demo showcasing:
- Native webview integration
- JavaScript ↔ Go bridge communication
- Real-time UI updates
- Task management with state persistence
- System information display

**Status**: Fully functional  
**Features**: Cross-platform webview, bidirectional communication  
**See**: [examples/02-webview-demo/](./02-webview-demo/)

### 03-python-webview-demo ✅ NEW

A comprehensive Python + JS + Webview integration demonstrating:
- **Real-time Python Execution**: Execute Python code from JavaScript
- **Mathematical Operations**: Access Python's math module
- **Statistical Analysis**: Use Python's statistics module
- **Text Processing**: Analyze text using Python
- **Data Transformation**: Python list comprehensions
- **Task Management**: Full CRUD with Go backend
- **Error Handling**: Detailed tracebacks and error messages

**Status**: Production-ready  
**Features**: Python runtime, webview UI, comprehensive examples  
**See**: [examples/03-python-webview-demo/](./03-python-webview-demo/)

## Quick Start

### Hello World Example
```bash
# Navigate to example
cd 01-hello-world

# Build and run
go build -o dist/hello-world ./src/backend
./dist/hello-world

# Run tests
go test -v

# Run benchmarks
go test -bench=. -benchtime=3s
```

### Webview Demo
```bash
# From project root
make build-webview-demo

# Or manually
cd examples/02-webview-demo
go build -o webview-demo
./webview-demo
```

### Python + Webview Demo (Recommended!)
```bash
# Ensure Python runtime is available
make setup-python  # If needed

# Build and run (from project root)
make run-python-demo

# Or manually
cd examples/03-python-webview-demo
go build -tags=runtime_python -o python-demo
./python-demo
```

## Example Status

| Example | Status | Runtimes | Tests | Performance |
|---------|--------|----------|-------|-------------|
| 01-hello-world | ✅ Working | Go, JS, Python (stub) | ✅ Passing | ⚡ Validated |
| 02-webview-demo | ✅ Working | Go, Webview, JS | ✅ Passing | ⚡ Fast |
| 03-python-webview-demo | ✅ Working | Go, Python, JS, Webview | ✅ Comprehensive | ⚡ Optimized |
| 04-ml-dashboard | 📋 Planned | Go, Python, React | - | - |
| 05-game-engine | 📋 Planned | Go, C++, Lua | - | - |

## Highlights

### Python + Webview Demo (03-python-webview-demo)
This is our most comprehensive example! It demonstrates:

**Python Features**:
- Execute arbitrary Python code from JavaScript
- Real-time calculations (Fibonacci, statistics, math operations)
- Text analysis and data transformation
- List comprehensions and data processing
- Full error handling with tracebacks

**UI Features**:
- Beautiful, modern interface
- Interactive Python calculator
- Statistical analysis with charts
- Text analyzer
- Data transformation tools
- Task management system

**Perfect for**:
- Learning Python integration
- Building data analysis tools
- Creating educational apps
- Prototyping scientific applications

Run it with: `make run-python-demo`

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
