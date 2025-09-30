# Python + JavaScript + Webview Demo

A comprehensive demonstration of Polyglot's seamless integration between Python, JavaScript, and Go through a webview interface.

## Features

### 🐍 Python Integration
- **Real-time Python Execution**: Execute Python code directly from JavaScript
- **Mathematical Operations**: Access Python's `math` module and advanced calculations
- **Data Processing**: Use Python's powerful list comprehensions and data structures
- **Statistical Analysis**: Leverage Python's `statistics` module for data analysis
- **Text Processing**: Analyze text using Python's string manipulation capabilities
- **Error Handling**: Full traceback and detailed error messages

### 🎨 Interactive UI Components
- **Python Calculator**: Execute arbitrary Python expressions
- **Fibonacci Generator**: Generate Fibonacci sequences using Python algorithms
- **Statistical Analysis**: Calculate mean, median, standard deviation, and more
- **Text Analyzer**: Analyze character count, word count, sentence count, etc.
- **Data Transformer**: Transform arrays using Python list operations
- **List Processing**: Demonstrate Python's comprehension capabilities

### 📋 Task Manager
- **CRUD Operations**: Create, Read, Update, Delete tasks
- **State Management**: Go backend manages application state
- **Filtering**: Filter tasks by status and priority
- **Priority Levels**: High, medium, and low priority support
- **Status Tracking**: Pending, in-progress, and completed states

### 💻 System Integration
- **Runtime Information**: View Go, Python, and system details
- **Performance Metrics**: Monitor CPU, goroutines, and uptime
- **Cross-Language Communication**: Seamless bridge between languages

## Prerequisites

### Required
- **Go 1.21+**: Main runtime
- **Python 3.9+** with development headers: For Python runtime
- **pkg-config**: For Python detection

### Installation

#### macOS
```bash
# Install Python with Homebrew
brew install python@3.11 pkg-config

# Configure pkg-config (add to ~/.zshrc or ~/.bashrc)
export PKG_CONFIG_PATH="$(brew --prefix python@3.11)/lib/pkgconfig:${PKG_CONFIG_PATH}"
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install python3-dev pkg-config
```

#### Linux (Fedora/RHEL)
```bash
sudo dnf install python3-devel pkgconfig
```

#### Linux (Arch)
```bash
sudo pacman -S python pkgconf
```

#### Windows
1. Install Python from [python.org](https://www.python.org/downloads/)
   - Check "Add Python to PATH"
   - Check "Include development headers"
2. Install MSYS2 from [msys2.org](https://www.msys2.org/)
3. Run: `pacman -S mingw-w64-x86_64-pkgconf`

## Building and Running

### Quick Start
```bash
# From this directory
go mod download
go build -tags=runtime_python -o python-demo
./python-demo
```

### Using Makefile (from project root)
```bash
# Build the example
make example-python

# Run the example
make run-python-example
```

### Manual Build
```bash
# With Python runtime (native)
go build -tags=runtime_python -o python-demo

# Without Python runtime (stub - limited functionality)
go build -o python-demo
```

## Usage Guide

### Python Calculator
1. Enter any valid Python expression in the input field
2. Examples:
   - `2**10` - Calculate 2 to the power of 10
   - `import math; math.sqrt(144)` - Square root
   - `sum(range(1, 101))` - Sum of numbers 1-100
3. Click "Calculate" to execute

### Fibonacci Generator
1. Enter a number (1-50)
2. Click "Generate" to calculate the Fibonacci number
3. Uses Python's efficient algorithm

### Statistical Analysis
1. Enter comma-separated numbers (e.g., `1,2,3,4,5`)
2. Click "Analyze" to calculate statistics
3. Or click "Random Data" to generate random numbers
4. View mean, median, standard deviation, min, max, sum, and count

### Text Analysis
1. Enter or paste text in the textarea
2. Click "Analyze" to process
3. View character count, word count, sentence count, average word length, longest word, and unique words

### Data Transformation
1. Enter comma-separated numbers
2. Select an operation:
   - **Square**: Square each number
   - **Double**: Double each number
   - **Reverse**: Reverse the array
   - **Sort**: Sort the array
3. Click "Transform" to apply

### List Processing
1. Enter a size (1-20)
2. Click "Process" to generate:
   - Range of numbers
   - Squares
   - Even numbers
   - Odd numbers
   - Sum of squares
   - Fibonacci sequence

### Task Manager
1. **View Tasks**: See all tasks in the Tasks tab
2. **Add Task**: 
   - Switch to "Add Task" tab
   - Fill in title, description, and priority
   - Click "Add Task"
3. **Update Tasks**:
   - Click "✓ Complete" to mark as done
   - Click "↺ Reopen" to mark as pending
4. **Delete Tasks**: Click "🗑 Delete"
5. **Filter Tasks**:
   - Switch to "Filter" tab
   - Select field and enter value
   - Click "Filter"

### System Information
1. Click "Load Info" to view:
   - Platform and architecture
   - Go version
   - Python version
   - CPU cores
   - Active goroutines
   - Application uptime

## Architecture

### Component Flow
```
┌─────────────────────────────────────────────┐
│            JavaScript (Frontend)             │
│  - User Interface                            │
│  - Event Handling                            │
│  - Data Presentation                         │
└─────────────────┬───────────────────────────┘
                  │ window.polyglot.call()
                  ↓
┌─────────────────────────────────────────────┐
│              Go Bridge (Core)                │
│  - Function Registration                     │
│  - Type Conversion                           │
│  - Error Handling                            │
└──────────┬──────────────────────┬────────────┘
           │                      │
           ↓                      ↓
┌────────────────────┐  ┌────────────────────┐
│  Go Functions      │  │  Python Runtime    │
│  - State Mgmt      │  │  - Code Execution  │
│  - Task CRUD       │  │  - Math Operations │
│  - System Info     │  │  - Data Processing │
└────────────────────┘  └────────────────────┘
```

### Data Flow Example
```javascript
// JavaScript calls Python through Go bridge
const result = await window.polyglot.call('pythonCalculate', 'math.pi * 2');

// Go Bridge receives call
func pythonCalculate(ctx context.Context, args ...interface{}) (interface{}, error) {
    // Execute in Python runtime
    result, err := appState.pythonRuntime.Execute(ctx, expr)
    return result, err
}

// Python executes code
import math
math.pi * 2  # Returns: 6.283185307179586
```

## Python Code Examples

All Python code in this demo is executed through the bridge. Here are some examples:

### Simple Expression
```python
2 + 2  # Returns: 4
```

### Using Modules
```python
import math
math.sqrt(144)  # Returns: 12.0
```

### List Comprehension
```python
[x*x for x in range(10)]  # Returns: [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
```

### Function Definition
```python
def fibonacci(n):
    if n <= 1:
        return n
    a, b = 0, 1
    for _ in range(n - 1):
        a, b = b, a + b
    return b

fibonacci(10)  # Returns: 55
```

### Dictionary Creation
```python
{
    'mean': statistics.mean(data),
    'median': statistics.median(data),
    'stdev': statistics.stdev(data)
}
```

## Error Handling

The demo includes comprehensive error handling:

### Python Errors
```python
# Syntax Error
this is not valid python  # Shows: SyntaxError with traceback

# Name Error
undefined_variable  # Shows: NameError: name 'undefined_variable' is not defined

# Type Error
len(42)  # Shows: TypeError: object of type 'int' has no len()

# Division by Zero
1 / 0  # Shows: ZeroDivisionError: division by zero
```

All Python errors include:
- Exception type (e.g., `ValueError`, `TypeError`)
- Error message
- Full traceback with line numbers
- Code context

### Bridge Errors
- Invalid arguments
- Type mismatches
- Runtime unavailable
- Timeout errors

## Performance

### Benchmarks (on M1 Mac)
- **Python Initialization**: ~5ms
- **Simple Calculation**: <1ms
- **Complex Operation**: 1-5ms
- **Bridge Overhead**: <0.1ms
- **UI Response**: <10ms

### Optimization Tips
1. **Batch Operations**: Group multiple calculations when possible
2. **Reuse State**: Python runtime reuses interpreter state
3. **Async Operations**: UI remains responsive during calculations
4. **Error Recovery**: Errors don't crash the application

## Troubleshooting

### "Python runtime not available"
**Problem**: Python runtime initialization failed

**Solutions**:
1. Verify Python installation: `python3 --version`
2. Check pkg-config: `pkg-config --exists python3-embed && echo "OK"`
3. Rebuild with Python support: `go build -tags=runtime_python`
4. See [Python Runtime README](../../runtimes/python/README.md) for detailed setup

### "Calculation failed"
**Problem**: Python execution error

**Solutions**:
1. Check Python syntax
2. Verify module availability (e.g., `import math`)
3. Look at error traceback for details
4. Test code in Python REPL first

### Build Errors
**Problem**: Compilation fails

**Solutions**:
1. Ensure Go 1.21+: `go version`
2. Update dependencies: `go mod download`
3. Check Python dev headers: `pkg-config --modversion python3-embed`
4. Set PKG_CONFIG_PATH (macOS): See installation instructions above

### Performance Issues
**Problem**: Slow execution

**Solutions**:
1. Reduce data size for analysis
2. Simplify Python code
3. Check system resources
4. Profile using system tools

## Examples of What You Can Build

This demo showcases patterns for building:

1. **Data Analysis Tools**: Use Python's scientific libraries
2. **Text Processors**: Leverage Python's string manipulation
3. **Calculation Engines**: Access Python's math capabilities
4. **Educational Apps**: Teach Python interactively
5. **Automation Tools**: Combine Go's performance with Python's flexibility
6. **Scientific Applications**: Use NumPy, SciPy (when installed)
7. **Machine Learning Interfaces**: Access TensorFlow, PyTorch (when installed)

## Extending the Demo

### Add New Python Functions
```go
// In setupBridge()
bridge.Register("myPythonFunc", func(ctx context.Context, args ...interface{}) (interface{}, error) {
    code := ` + "`" + `
    # Your Python code here
    result = some_calculation()
    result
    ` + "`" + `
    return appState.pythonRuntime.Execute(ctx, code)
})
```

### Add UI Components
```html
<!-- In generateDemoHTML() -->
<div class="demo-section">
    <h2>My Feature</h2>
    <button onclick="myFunction()">Execute</button>
    <div id="myResult" class="result-box"></div>
</div>
```

### Add JavaScript Handlers
```javascript
async function myFunction() {
    try {
        const result = await window.polyglot.call('myPythonFunc', arg1, arg2);
        showResult('myResult', result);
    } catch (error) {
        showMessage('Error: ' + error, 'error', 'myResult');
    }
}
```

## Best Practices

1. **Error Handling**: Always wrap Python calls in try-catch
2. **Input Validation**: Validate user input before sending to Python
3. **Type Checking**: Handle type conversions carefully
4. **User Feedback**: Show loading states for long operations
5. **Resource Cleanup**: Let Go manage Python runtime lifecycle
6. **Security**: Validate and sanitize user-provided Python code in production
7. **Documentation**: Comment complex Python algorithms
8. **Testing**: Test edge cases and error conditions

## Known Limitations

1. **Python Packages**: Only stdlib modules available by default
   - Solution: Install packages in Python environment
2. **Long Operations**: Very long calculations may timeout
   - Solution: Increase timeout in RuntimeConfig
3. **Memory**: Large data structures use more memory
   - Solution: Process data in chunks
4. **Platform**: Some Python features are platform-specific
   - Solution: Test on target platforms

## Security Considerations

⚠️ **Important**: This demo executes arbitrary Python code for demonstration purposes.

In production applications:
1. **Never** execute untrusted user-provided Python code directly
2. Use sandboxing for user code execution
3. Implement allowlists for permitted operations
4. Add input validation and sanitization
5. Set strict timeouts
6. Monitor resource usage
7. Implement audit logging

## Resources

- [Polyglot Main README](../../README.md)
- [Python Runtime Documentation](../../runtimes/python/README.md)
- [Webview Documentation](../../webview/README.md)
- [Core Bridge Documentation](../../core/README.md)

## Contributing

Found a bug or have a feature request? Please open an issue or submit a pull request!

## License

MIT License - See [LICENSE](../../LICENSE) file for details.
