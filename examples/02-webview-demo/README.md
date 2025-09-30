# Webview Demo Example

This example demonstrates the native webview integration in Polyglot, showcasing how to create a desktop application with a web-based UI that communicates with Go backend code.

## Features Demonstrated

1. **Counter Demo**: Simple state management between Go and JavaScript
2. **Greeting Demo**: Passing strings between frontend and backend
3. **Todo List**: Full CRUD operations with complex data structures
4. **Random Number Generator**: Calling Go functions with parameters
5. **System Info**: Retrieving structured data from the backend

## Building

### Standard Build (Native Webview)

Native webview is enabled by default. Just build normally:

```bash
# Navigate to the example directory
cd examples/02-webview-demo

# Build (native webview is default)
go build -o webview-demo

# Run the application
./webview-demo
```

Or use the Makefile from the project root:

```bash
# From project root
make build-webview-demo

# Run
cd examples/02-webview-demo && ./webview-demo
```

### Platform-Specific Requirements

**macOS:**
```bash
# No additional dependencies needed (uses native WebKit)
go build -o webview-demo
```

**Linux (Ubuntu/Debian):**
```bash
# Install WebKitGTK first
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev

# Build
go build -o webview-demo
```

**Windows:**
```bash
# Install WebView2 Runtime from Microsoft
# Then build with CGO enabled
set CGO_ENABLED=1
go build -o webview-demo.exe
```

### With Stub Backend (Testing/CI)

For testing without native dependencies (e.g., in CI):

```bash
# Build with stub backend (no GUI)
go build -tags stub -o webview-demo

# This will use a stub backend that logs operations to console
./webview-demo
```

## Architecture

### Go Backend (main.go)

The backend implements a `DemoBridge` that provides functions callable from JavaScript:

- `increment()`: Increments and returns a counter
- `getCounter()`: Returns the current counter value
- `greet(name)`: Returns a personalized greeting
- `getTodos()`: Returns the list of todos
- `addTodo(title)`: Creates a new todo
- `toggleTodo(id)`: Toggles a todo's completed status
- `deleteTodo(id)`: Deletes a todo
- `getRandomNumber(min, max)`: Generates a random number
- `getSystemInfo()`: Returns system information

### Frontend (Embedded HTML)

The frontend is an HTML page embedded directly in the Go binary using a data URL. It features:

- Modern, responsive UI with gradient backgrounds
- Real-time communication with the Go backend
- Error handling and user feedback
- Keyboard shortcuts (Enter key support)

### Bridge Communication

The `polyglot` JavaScript object provides seamless communication:

```javascript
// Call Go function from JavaScript
const result = await window.polyglot.call('functionName', arg1, arg2);
```

## Code Structure

```
02-webview-demo/
├── main.go           # Main application with bridge implementation
└── README.md         # This file
```

## Customization

### Adding New Functions

1. **In Go**: Add a new case to the `DemoBridge.Call` method:

```go
case "myFunction":
    if len(args) != 1 {
        return nil, fmt.Errorf("myFunction requires 1 argument")
    }
    param := args[0].(string)
    return fmt.Sprintf("Result: %s", param), nil
```

2. **In JavaScript**: Call the function:

```javascript
async function myFunction() {
    const result = await window.polyglot.call('myFunction', 'param');
    console.log(result);
}
```

### Loading External HTML

Instead of embedding HTML, you can load it from a file or URL:

```go
config := core.WebviewConfig{
    // Load from file
    URL: "file:///path/to/index.html",
    
    // Or load from web server
    URL: "http://localhost:3000",
}
```

### Styling

The embedded HTML includes inline CSS. You can:
- Modify the styles directly in `generateDemoHTML()`
- Load external CSS files
- Use a frontend framework like React or Vue

## Advanced Usage

### Multi-Window Applications

You can create multiple webview instances:

```go
// Create main window
mainWindow := webview.New(mainConfig, bridge)
mainWindow.Initialize()

// Create secondary window
secondWindow := webview.New(secondConfig, bridge)
secondWindow.Initialize()

// Run main window (blocks)
go mainWindow.Run()

// Run second window
secondWindow.Run()
```

### Custom Window Sizing

```go
config := core.WebviewConfig{
    Width:     1920,
    Height:    1080,
    Resizable: false,  // Fixed size window
}
```

### Debug Mode

Enable debug mode to open developer tools:

```go
config := core.WebviewConfig{
    Debug: true,  // Enables DevTools (F12)
}
```

## Troubleshooting

### "Using stub webview backend"

If you see this message but want native webview, make sure you're not using `-tags stub` and have the required dependencies installed.

### "failed to create webview"

Make sure you have the required dependencies installed for your platform (see Platform-Specific Requirements).

### Window doesn't appear

The webview runs on the main thread. Make sure you're calling `Run()` on the main goroutine.

### JavaScript errors

Open the developer console (F12 when debug is enabled) to see JavaScript errors and logs.

## Next Steps

- Explore the [main webview documentation](../../webview/BUILD.md)
- Check out [production build instructions](../../BUILD.md)
- Learn about [cross-platform deployment](../../docs/DEPLOYMENT.md)

## License

This example is part of the Polyglot project and follows the same license.
