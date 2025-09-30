package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const version = "0.1.0"

func handleInit(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: project name required")
		fmt.Println("Usage: polyglot init <project-name>")
		os.Exit(1)
	}

	projectName := args[0]
	fmt.Printf("Initializing Polyglot project: %s\n", projectName)

	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Create project structure
	dirs := []string{
		"src",
		"src/backend",
		"src/frontend",
		"dist",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("Error creating %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Create main.go
	mainGo := `package main

import (
	"context"
	"log"
	"time"
	
	"github.com/polyglot-framework/polyglot/core"
	"github.com/polyglot-framework/polyglot/runtimes/python"
	"github.com/polyglot-framework/polyglot/runtimes/javascript"
	"github.com/polyglot-framework/polyglot/webview"
)

func main() {
	// Create configuration
	config := core.DefaultConfig()
	config.App.Name = "` + projectName + `"
	config.EnableRuntime("python", "3.11")
	config.EnableRuntime("javascript", "latest")
	
	// Create orchestrator
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %v", err)
	}
	
	// Register runtimes
	orch.RegisterRuntime(python.NewRuntime())
	orch.RegisterRuntime(javascript.NewRuntime())
	
	// Initialize
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	if err := orch.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	
	// Create bridge
	bridge := core.NewBridge()
	
	// Register example function
	bridge.Register("greet", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		name := "World"
		if len(args) > 0 {
			if n, ok := args[0].(string); ok {
				name = n
			}
		}
		return "Hello, " + name + "!", nil
	})
	
	// Create webview
	wv := webview.New(config.Webview, bridge)
	if err := wv.Initialize(); err != nil {
		log.Fatalf("Failed to initialize webview: %v", err)
	}
	
	// Run
	log.Println("Starting application...")
	if err := wv.Run(); err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
	
	// Cleanup
	orch.Shutdown(context.Background())
}
`

	mainPath := filepath.Join(projectName, "src", "backend", "main.go")
	if err := os.WriteFile(mainPath, []byte(mainGo), 0644); err != nil {
		fmt.Printf("Error writing main.go: %v\n", err)
		os.Exit(1)
	}

	// Create index.html
	indexHTML := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>` + projectName + `</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		}
		.container {
			text-align: center;
			color: white;
		}
		h1 {
			font-size: 3rem;
			margin-bottom: 1rem;
		}
		button {
			padding: 1rem 2rem;
			font-size: 1.2rem;
			background: white;
			color: #667eea;
			border: none;
			border-radius: 8px;
			cursor: pointer;
			transition: transform 0.2s;
		}
		button:hover {
			transform: scale(1.05);
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Welcome to ` + projectName + `</h1>
		<p id="message">Click the button below</p>
		<button onclick="callBackend()">Call Backend</button>
	</div>
	
	<script>
		async function callBackend() {
			try {
				const result = await window.polyglot.call('greet', 'Polyglot');
				document.getElementById('message').textContent = result;
			} catch (err) {
				console.error('Error:', err);
				document.getElementById('message').textContent = 'Error: ' + err.message;
			}
		}
	</script>
</body>
</html>
`

	htmlPath := filepath.Join(projectName, "src", "frontend", "index.html")
	if err := os.WriteFile(htmlPath, []byte(indexHTML), 0644); err != nil {
		fmt.Printf("Error writing index.html: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nProject %s created successfully!\n", projectName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  polyglot dev")
}

func handleBuild(args []string) {
	fmt.Println("Building application...")

	cmd := exec.Command("go", "build", "-o", "dist/app", "./src/backend")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Build complete!")
}

func handleDev(args []string) {
	fmt.Println("Starting development mode...")
	fmt.Println("This will implement hot reload in the future")
	handleBuild(args)
}

func handleTest(args []string) {
	fmt.Println("Running tests...")

	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed!")
}

func handleVersion(args []string) {
	fmt.Printf("Polyglot CLI v%s\n", version)
}
