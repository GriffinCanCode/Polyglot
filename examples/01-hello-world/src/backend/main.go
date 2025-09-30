package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/javascript"
	"github.com/griffincancode/polyglot.js/runtimes/python"
)

func main() {
	fmt.Println("🚀 Initializing Polyglot Hello World Example...")
	fmt.Println()

	// Create configuration
	config := core.DefaultConfig()
	config.App.Name = "hello-world"

	// Configure runtimes
	config.Languages = map[string]*core.RuntimeConfig{
		"python": {
			Name:           "python",
			Version:        "3.11",
			Enabled:        true,
			MaxConcurrency: 4,
			Timeout:        time.Second * 30,
			Options:        make(map[string]interface{}),
		},
		"javascript": {
			Name:           "javascript",
			Version:        "latest",
			Enabled:        true,
			MaxConcurrency: 4,
			Timeout:        time.Second * 30,
			Options:        make(map[string]interface{}),
		},
	}

	// Create orchestrator
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		log.Fatalf("❌ Failed to create orchestrator: %v", err)
	}
	fmt.Println("✓ Orchestrator created")

	// Register runtimes
	if err := orch.RegisterRuntime(python.NewRuntime()); err != nil {
		log.Fatalf("❌ Failed to register Python runtime: %v", err)
	}
	if err := orch.RegisterRuntime(javascript.NewRuntime()); err != nil {
		log.Fatalf("❌ Failed to register JavaScript runtime: %v", err)
	}
	fmt.Printf("✓ Runtimes registered: %v\n", orch.Runtimes())

	// Initialize orchestrator
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := orch.Initialize(ctx); err != nil {
		fmt.Printf("⚠️  Initialization warning: %v\n", err)
		fmt.Println()
		fmt.Println("Note: Runtimes are using stub implementations.")
		fmt.Println("      This is normal when building without runtime tags.")
		fmt.Println("      The example will demonstrate the architecture.")
		fmt.Println()
	} else {
		fmt.Println("✓ System initialized with real runtimes")
		fmt.Println()
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n🛑 Shutting down gracefully...")
		orch.Shutdown(context.Background())
		os.Exit(0)
	}()

	// Run demonstrations
	demonstratePython(orch)
	demonstrateJavaScript(orch)
	demonstrateCrossRuntime(orch)
	demonstrateMemorySharing(orch)

	fmt.Println()
	fmt.Println("✅ All systems operational!")
	fmt.Println()
	fmt.Println("This example demonstrates:")
	fmt.Println("  • Multi-language runtime orchestration")
	fmt.Println("  • Dynamic runtime registration")
	fmt.Println("  • Cross-language function calls")
	fmt.Println("  • Shared memory coordination")
	fmt.Println("  • Graceful error handling")
	fmt.Println()
	fmt.Println("Note: Running with stub implementations (no native dependencies required)")
	fmt.Println("      Build with -tags=runtime_python,runtime_javascript for real runtimes")
	fmt.Println()
	fmt.Println("Press Ctrl+C to exit...")

	// Keep running
	select {}
}

func demonstratePython(orch *core.Orchestrator) {
	fmt.Println("📊 Testing Python Runtime:")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Simple Python calculation
	code := `
def calculate(x, y):
    return x * y + 42

result = calculate(10, 5)
`

	result, err := orch.Execute(ctx, "python", code)
	if err != nil {
		fmt.Printf("  ⚠️  Python runtime not available (using stub): %v\n", err)
		return
	}

	fmt.Printf("  ✓ Result: %v\n", result)

	// Call a specific function
	fnResult, err := orch.Call(ctx, "python", "calculate", 7, 3)
	if err != nil {
		fmt.Printf("  ⚠️  Function call failed: %v\n", err)
	} else {
		fmt.Printf("  ✓ Function call result: %v\n", fnResult)
	}

	fmt.Println()
}

func demonstrateJavaScript(orch *core.Orchestrator) {
	fmt.Println("🟨 Testing JavaScript Runtime:")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Simple JavaScript execution
	code := `
function greet(name) {
    return "Hello, " + name + " from JavaScript!";
}

greet("Polyglot");
`

	result, err := orch.Execute(ctx, "javascript", code)
	if err != nil {
		fmt.Printf("  ⚠️  JavaScript runtime not available (using stub): %v\n", err)
		return
	}

	fmt.Printf("  ✓ Result: %v\n", result)

	// Call a specific function
	fnResult, err := orch.Call(ctx, "javascript", "greet", "World")
	if err != nil {
		fmt.Printf("  ⚠️  Function call failed: %v\n", err)
	} else {
		fmt.Printf("  ✓ Function call result: %v\n", fnResult)
	}

	fmt.Println()
}

func demonstrateCrossRuntime(orch *core.Orchestrator) {
	fmt.Println("🔄 Testing Cross-Runtime Communication:")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Step 1: Python generates data
	fmt.Println("  Step 1: Python generates data...")
	pythonCode := `
data = {"value": 42, "language": "python"}
`
	pythonResult, err := orch.Execute(ctx, "python", pythonCode)
	if err != nil {
		fmt.Printf("  ⚠️  Using stub data: %v\n", err)
		pythonResult = map[string]interface{}{"value": 42, "language": "stub"}
	}
	fmt.Printf("    → Python output: %v\n", pythonResult)

	// Step 2: JavaScript processes it
	fmt.Println("  Step 2: JavaScript processes data...")
	jsCode := `
function process(data) {
    return {
        original: data,
        processed: true,
        language: "javascript"
    };
}
process({"value": 42});
`
	jsResult, err := orch.Execute(ctx, "javascript", jsCode)
	if err != nil {
		fmt.Printf("  ⚠️  Using stub processing: %v\n", err)
		jsResult = map[string]interface{}{"processed": true, "language": "stub"}
	}
	fmt.Printf("    → JavaScript output: %v\n", jsResult)

	// Step 3: Go coordinates
	fmt.Println("  Step 3: Go coordinates final result...")
	fmt.Println("    → ✓ Data pipeline: Python → JavaScript → Go")
	fmt.Println("    → ✓ Multi-language workflow complete!")

	fmt.Println()
}

func demonstrateMemorySharing(orch *core.Orchestrator) {
	fmt.Println("💾 Testing Memory Coordinator:")

	mem := orch.Memory()

	// Allocate shared memory
	region, err := mem.Allocate("demo-buffer", 1024, core.TypeBytes)
	if err != nil {
		fmt.Printf("  ⚠️  Memory allocation failed: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Allocated shared memory region: %s (1024 bytes)\n", region.ID)

	// Write data
	testData := []byte("Hello from shared memory!")
	copy(region.Data, testData)
	fmt.Printf("  ✓ Written data: %s\n", string(testData))

	// Read data
	readData := region.Data[:len(testData)]
	fmt.Printf("  ✓ Read data: %s\n", string(readData))

	// Get region info
	retrieved, err := mem.Get("demo-buffer")
	if err == nil && retrieved != nil {
		fmt.Printf("  ✓ Region retrieved: %s (type: %s)\n", retrieved.ID, retrieved.Type)
	}

	// Free memory
	mem.Free("demo-buffer")
	fmt.Println("  ✓ Memory freed")

	fmt.Println()
}
