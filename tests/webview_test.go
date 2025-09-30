package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/webview"
)

// Test basic webview creation and configuration
func TestWebview_Creation(t *testing.T) {
	config := core.WebviewConfig{
		Title:     "Test Window",
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     false,
		URL:       "https://example.com",
	}

	wv := webview.New(config, nil)

	if wv == nil {
		t.Fatal("Failed to create webview")
	}
}

// Test webview initialization
func TestWebview_Initialize(t *testing.T) {
	// Skip in CI or when DISPLAY is not set (headless)
	if testing.Short() {
		t.Skip("Skipping webview test in short mode")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	config := core.WebviewConfig{
		Title:  "Test Initialize",
		Width:  1024,
		Height: 768,
		Debug:  false, // Disable debug in tests
		URL:    "data:text/html,<html><body><h1>Test</h1></body></html>",
	}

	wv := webview.New(config, nil)

	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize webview: %v", err)
	}

	// Clean up
	defer wv.Terminate()
}

// Test JavaScript evaluation (using stub backend)
func TestWebview_Eval(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping webview test in short mode")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	config := core.WebviewConfig{
		Title:  "Test Eval",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	wv := webview.New(config, nil)
	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer wv.Terminate()

	// Test eval (with stub backend, this just logs)
	err = wv.Eval("console.log('Hello from Go')")
	if err != nil {
		t.Errorf("Eval failed: %v", err)
	}
}

// Test function binding
func TestWebview_Bind(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping webview test in short mode")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	config := core.WebviewConfig{
		Title:  "Test Bind",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	wv := webview.New(config, nil)
	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer wv.Terminate()

	// Bind a simple function
	err = wv.Bind("testFunc", func() string {
		return "Hello from Go"
	})

	if err != nil {
		t.Errorf("Bind failed: %v", err)
	}

	// Note: With stub backend, the function won't actually be called
	// This test mainly verifies the binding interface works
}

// Test bridge integration
func TestWebview_BridgeIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping webview test in short mode")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	config := core.WebviewConfig{
		Title:  "Test Bridge",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	bridge := core.NewBridge()

	// Register test functions
	bridge.Register("add", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("add requires 2 arguments")
		}

		a, ok1 := args[0].(float64)
		b, ok2 := args[1].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("arguments must be numbers")
		}

		return a + b, nil
	})

	bridge.Register("greet", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("greet requires 1 argument")
		}

		name, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("name must be a string")
		}

		return fmt.Sprintf("Hello, %s!", name), nil
	})

	wv := webview.New(config, bridge)
	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer wv.Terminate()

	// Test bridge calls directly
	result, err := bridge.Call(nil, "add", 5.0, 3.0)
	if err != nil {
		t.Errorf("Bridge call failed: %v", err)
	}
	if result != 8.0 {
		t.Errorf("Expected 8.0, got %v", result)
	}

	result, err = bridge.Call(nil, "greet", "World")
	if err != nil {
		t.Errorf("Bridge call failed: %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %v", result)
	}
}

// Test concurrent webview operations
func TestWebview_Concurrency(t *testing.T) {
	config := core.WebviewConfig{
		Title:  "Test Concurrent",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	wv := webview.New(config, nil)
	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer wv.Terminate()

	// Execute multiple eval operations concurrently
	var wg sync.WaitGroup
	errorChan := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			script := fmt.Sprintf("console.log('Message %d')", n)
			if err := wv.Eval(script); err != nil {
				errorChan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	// Check for errors
	for err := range errorChan {
		t.Errorf("Concurrent eval error: %v", err)
	}
}

// Test webview lifecycle
func TestWebview_Lifecycle(t *testing.T) {
	config := core.WebviewConfig{
		Title:  "Test Lifecycle",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	wv := webview.New(config, nil)

	// Initialize
	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Perform operations
	err = wv.Eval("console.log('test')")
	if err != nil {
		t.Errorf("Eval failed: %v", err)
	}

	// Terminate
	err = wv.Terminate()
	if err != nil {
		t.Errorf("Terminate failed: %v", err)
	}

	// Operations after terminate should fail
	err = wv.Eval("console.log('after terminate')")
	if err == nil {
		t.Error("Expected error after terminate, got nil")
	}
}

// Test error handling
func TestWebview_ErrorHandling(t *testing.T) {
	config := core.WebviewConfig{
		Title:  "Test Errors",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	wv := webview.New(config, nil)

	// Eval before initialize should fail
	err := wv.Eval("console.log('test')")
	if err == nil {
		t.Error("Expected error before initialize, got nil")
	}

	// Initialize
	err = wv.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer wv.Terminate()

	// Double initialize should fail
	err = wv.Initialize()
	if err == nil {
		t.Error("Expected error on double initialize, got nil")
	}
}

// Test stub backend functionality
func TestWebview_StubBackend(t *testing.T) {
	// Test webview with nil bridge (uses stub backend when built without native support)
	config := core.WebviewConfig{
		Title:  "Stub Test",
		Width:  1024,
		Height: 768,
		URL:    "https://example.com",
	}

	wv := webview.New(config, nil)
	if wv == nil {
		t.Fatal("Failed to create webview")
	}

	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Perform operations
	err = wv.Eval("console.log('test')")
	if err != nil {
		t.Errorf("Eval failed: %v", err)
	}

	err = wv.Bind("testFunc", func() string { return "test" })
	if err != nil {
		t.Errorf("Bind failed: %v", err)
	}

	// Clean up
	err = wv.Terminate()
	if err != nil {
		t.Errorf("Terminate failed: %v", err)
	}

	// Should not panic or error
}

// Test webview with data URL
func TestWebview_DataURL(t *testing.T) {
	htmlContent := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
    <script>
        window.addEventListener('DOMContentLoaded', () => {
            console.log('Page loaded');
        });
    </script>
</head>
<body>
    <h1>Test Page</h1>
    <p>This is a test page loaded from a data URL</p>
</body>
</html>
`

	config := core.WebviewConfig{
		Title:  "Data URL Test",
		Width:  800,
		Height: 600,
		URL:    "data:text/html," + htmlContent,
	}

	wv := webview.New(config, nil)
	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize with data URL: %v", err)
	}
	defer wv.Terminate()
}

// Test JSON serialization in bridge
func TestWebview_JSONSerialization(t *testing.T) {
	bridge := core.NewBridge()

	// Register function that returns complex data
	bridge.Register("getUser", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		return map[string]interface{}{
			"id":   123,
			"name": "John Doe",
			"tags": []string{"admin", "developer"},
		}, nil
	})

	// Call and verify serialization
	result, err := bridge.Call(nil, "getUser")
	if err != nil {
		t.Fatalf("Bridge call failed: %v", err)
	}

	// Verify result can be serialized to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonData, &parsed)
	if err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	if parsed["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%v'", parsed["name"])
	}
}

// Benchmark webview initialization
func BenchmarkWebview_Initialize(b *testing.B) {
	config := core.WebviewConfig{
		Title:  "Benchmark",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Benchmark</body></html>",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wv := webview.New(config, nil)
		_ = wv.Initialize()
		_ = wv.Terminate()
	}
}

// Benchmark eval operations
func BenchmarkWebview_Eval(b *testing.B) {
	config := core.WebviewConfig{
		Title:  "Benchmark Eval",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	wv := webview.New(config, nil)
	_ = wv.Initialize()
	defer wv.Terminate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wv.Eval("console.log('test')")
	}
}

// Benchmark bridge calls
func BenchmarkWebview_BridgeCalls(b *testing.B) {
	bridge := core.NewBridge()
	bridge.Register("test", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		return "result", nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bridge.Call(nil, "test")
	}
}

// Example test for documentation
func ExampleWebview() {
	config := core.WebviewConfig{
		Title:     "Example App",
		Width:     1280,
		Height:    720,
		Resizable: true,
		Debug:     false,
		URL:       "https://example.com",
	}

	wv := webview.New(config, nil)
	defer wv.Terminate()

	if err := wv.Initialize(); err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		return
	}

	// In a real application, you would call wv.Run() which blocks
	// For this example, we just demonstrate the setup
	fmt.Println("Webview initialized successfully")
}

// Test running webview in goroutine
func TestWebview_Goroutine(t *testing.T) {
	config := core.WebviewConfig{
		Title:  "Goroutine Test",
		Width:  800,
		Height: 600,
		URL:    "data:text/html,<html><body>Test</body></html>",
	}

	wv := webview.New(config, nil)
	err := wv.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	done := make(chan bool)

	go func() {
		// Simulate some work
		time.Sleep(100 * time.Millisecond)

		err := wv.Eval("console.log('From goroutine')")
		if err != nil {
			t.Errorf("Eval from goroutine failed: %v", err)
		}

		done <- true
	}()

	// Wait for goroutine with timeout
	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Error("Goroutine test timed out")
	}

	wv.Terminate()
}
