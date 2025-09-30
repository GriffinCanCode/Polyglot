# Hello World Example

A simple Polyglot application demonstrating multi-language orchestration.

## Features

- **Go**: Main orchestrator and coordination
- **Python**: Data processing and calculations
- **JavaScript**: UI logic and interactions
- **Webview**: Native UI with bidirectional communication

## Architecture

```
Go Orchestrator
├── Python Runtime (math operations)
├── JavaScript Runtime (UI logic)
└── Webview Bridge (frontend ↔ backend)
```

## Building

### Option 1: With Stub Runtimes (No Dependencies)

```bash
cd examples/01-hello-world
go build -o dist/hello-world ./src/backend
./dist/hello-world
```

### Option 2: With Real Python Runtime

Requires: Python 3.11+ with dev headers

```bash
# Install dependencies (macOS)
brew install python@3.11

# Build with Python support
cd examples/01-hello-world
go build -tags=runtime_python -o dist/hello-world ./src/backend
./dist/hello-world
```

### Option 3: With Multiple Real Runtimes

Requires: Python 3.11+, V8 libraries

```bash
# Build with multiple runtimes
go build -tags=runtime_python,runtime_javascript -o dist/hello-world ./src/backend
./dist/hello-world
```

## What It Demonstrates

1. **Multi-language orchestration**: Go coordinates Python and JavaScript
2. **Runtime registration**: Dynamic runtime loading
3. **Cross-language calls**: Go calling Python and JavaScript
4. **Error handling**: Graceful degradation when runtimes unavailable
5. **Shared memory**: Passing data between languages
6. **Webview integration**: Native UI with backend bridge

## Expected Output

```
Initializing Polyglot Hello World...
✓ Orchestrator created
✓ Runtimes registered: [python javascript]
✓ System initialized

Testing Python Runtime:
  Result: Python calculation completed (stub)
  
Testing JavaScript Runtime:
  Result: JavaScript execution completed (stub)

Testing Cross-Runtime Communication:
  Data passed between Go → Python → JavaScript
  Final result: Hello from Polyglot!

All systems operational!
Press Ctrl+C to exit...
```

## Notes

- Without build tags, stub implementations are used (safe fallback)
- Stubs return mock data but validate the architecture
- Real runtimes require native dependencies
- The framework gracefully handles missing runtimes
