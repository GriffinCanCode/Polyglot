package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectTemplate handles project template generation
type ProjectTemplate struct {
	config *ProjectConfig
}

// NewTemplate creates a new template generator
func NewTemplate(config *ProjectConfig) *ProjectTemplate {
	return &ProjectTemplate{config: config}
}

// Generate creates the project structure and files
func (t *ProjectTemplate) Generate() error {
	// Create base directories
	if err := t.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate configuration files
	if err := t.generateConfig(); err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Generate main application file
	if err := t.generateMain(); err != nil {
		return fmt.Errorf("failed to generate main: %w", err)
	}

	// Generate frontend files
	if contains(t.config.Features, "webview") {
		if err := t.generateFrontend(); err != nil {
			return fmt.Errorf("failed to generate frontend: %w", err)
		}
	}

	// Generate README
	if err := t.generateReadme(); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	// Generate .gitignore
	if err := t.generateGitignore(); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	// Generate LICENSE
	if err := t.generateLicense(); err != nil {
		return fmt.Errorf("failed to generate LICENSE: %w", err)
	}

	// Generate Go module files
	if err := t.generateGoMod(); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	// Generate package files for languages
	if err := t.generatePackageFiles(); err != nil {
		return fmt.Errorf("failed to generate package files: %w", err)
	}

	// Generate Makefile
	if err := t.generateMakefile(); err != nil {
		return fmt.Errorf("failed to generate Makefile: %w", err)
	}

	return nil
}

func (t *ProjectTemplate) createDirectories() error {
	dirs := []string{
		t.config.Name,
		filepath.Join(t.config.Name, "src"),
		filepath.Join(t.config.Name, "src", "backend"),
		filepath.Join(t.config.Name, "dist"),
		filepath.Join(t.config.Name, ".polyglot"),
	}

	if contains(t.config.Features, "webview") {
		dirs = append(dirs,
			filepath.Join(t.config.Name, "src", "frontend"),
			filepath.Join(t.config.Name, "src", "frontend", "assets"),
			filepath.Join(t.config.Name, "src", "frontend", "styles"),
			filepath.Join(t.config.Name, "src", "frontend", "scripts"),
		)
	}

	// Language-specific directories
	for _, lang := range t.config.Languages {
		switch lang {
		case "python":
			dirs = append(dirs, filepath.Join(t.config.Name, "src", "python"))
		case "javascript":
			dirs = append(dirs, filepath.Join(t.config.Name, "src", "js"))
		case "rust":
			dirs = append(dirs, filepath.Join(t.config.Name, "src", "rust"))
		}
	}

	// Template-specific directories
	switch t.config.Template {
	case "cli":
		dirs = append(dirs, filepath.Join(t.config.Name, "cmd"))
	case "system":
		dirs = append(dirs,
			filepath.Join(t.config.Name, "services"),
			filepath.Join(t.config.Name, "config"),
		)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (t *ProjectTemplate) generateConfig() error {
	content := fmt.Sprintf(`{
  "name": "%s",
  "version": "%s",
  "description": "%s",
  "author": "%s",
  "license": "%s",
  "template": "%s",
  "languages": [%s],
  "features": [%s],
  "webview": {
    "width": %d,
    "height": %d,
    "resizable": %t,
    "devTools": %t,
    "title": "%s"
  },
  "memory": {
    "maxSharedMemory": 1073741824,
    "enableZeroCopy": true,
    "gcInterval": "5m"
  },
  "build": {
    "outputPath": "./dist",
    "optimize": true,
    "compress": true,
    "platforms": ["darwin", "linux", "windows"]
  },
  "runtimes": {
%s
  }
}`,
		t.config.Name,
		t.config.Version,
		t.config.Description,
		t.config.Author,
		t.config.License,
		t.config.Template,
		t.formatStringArray(t.config.Languages),
		t.formatStringArray(t.config.Features),
		t.config.WindowWidth,
		t.config.WindowHeight,
		t.config.WindowResizable,
		t.config.DevTools,
		t.config.Name,
		t.generateRuntimesConfig(),
	)

	path := filepath.Join(t.config.Name, "polyglot.config.json")
	return os.WriteFile(path, []byte(content), 0644)
}

func (t *ProjectTemplate) generateRuntimesConfig() string {
	var configs []string

	for _, lang := range t.config.Languages {
		var version string
		switch lang {
		case "python":
			version = t.config.PythonVersion
			if version == "" {
				version = "3.11"
			}
		case "javascript":
			version = "latest"
		default:
			version = "latest"
		}

		config := fmt.Sprintf(`    "%s": {
      "enabled": true,
      "version": "%s",
      "maxConcurrency": 10,
      "timeout": "30s"
    }`, lang, version)
		configs = append(configs, config)
	}

	return strings.Join(configs, ",\n")
}

func (t *ProjectTemplate) generateMain() error {
	var content string

	switch t.config.Template {
	case "webapp", "desktop":
		content = t.generateWebAppMain()
	case "cli":
		content = t.generateCLIMain()
	case "system":
		content = t.generateSystemMain()
	default:
		content = t.generateMinimalMain()
	}

	path := filepath.Join(t.config.Name, "src", "backend", "main.go")
	return os.WriteFile(path, []byte(content), 0644)
}

func (t *ProjectTemplate) generateWebAppMain() string {
	runtimeImports := t.generateRuntimeImports()
	runtimeRegistrations := t.generateRuntimeRegistrations()
	bridgeFunctions := t.generateBridgeFunctions()

	return fmt.Sprintf(`package main

import (
	"context"
	"log"
	"time"
	
	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/webview"
%s
)

func main() {
	// Load configuration from polyglot.config.json
	config := core.DefaultConfig()
	config.App.Name = "%s"
	config.App.Version = "%s"
	
	// Configure runtimes
%s
	
	// Create orchestrator
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %%v", err)
	}
	
	// Register runtimes
%s
	
	// Initialize
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	if err := orch.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %%v", err)
	}
	defer orch.Shutdown(context.Background())
	
	// Create bridge for frontend-backend communication
	bridge := core.NewBridge()
	
	// Register API functions
%s
	
	// Configure webview
	config.Webview.Title = "%s"
	config.Webview.Width = %d
	config.Webview.Height = %d
	config.Webview.Resizable = %t
	config.Webview.Debug = %t
	
	// Create and initialize webview
	wv := webview.New(config.Webview, bridge)
	if err := wv.Initialize(); err != nil {
		log.Fatalf("Failed to initialize webview: %%v", err)
	}
	
	log.Println("Starting %s...")
	
	// Run application
	if err := wv.Run(); err != nil {
		log.Fatalf("Failed to run application: %%v", err)
	}
}
`,
		runtimeImports,
		t.config.Name,
		t.config.Version,
		t.generateRuntimeConfigs(),
		runtimeRegistrations,
		bridgeFunctions,
		t.config.Name,
		t.config.WindowWidth,
		t.config.WindowHeight,
		t.config.WindowResizable,
		t.config.DevTools,
		t.config.Name,
	)
}

func (t *ProjectTemplate) generateCLIMain() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/griffincancode/polyglot.js/core"
%s
)

var (
	version = "%s"
)

func main() {
	// Define CLI flags
	versionFlag := flag.Bool("version", false, "Show version information")
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.Parse()
	
	if *versionFlag {
		fmt.Printf("%s v%%s\n", version)
		os.Exit(0)
	}
	
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}
	
	// Create configuration
	config := core.DefaultConfig()
	config.App.Name = "%s"
	config.App.Version = "%s"
	
%s
	
	// Create orchestrator
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %%v", err)
	}
	
%s
	
	// Initialize
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	if err := orch.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %%v", err)
	}
	defer orch.Shutdown(context.Background())
	
	// Your CLI logic here
	args := flag.Args()
	if len(args) == 0 {
		printHelp()
		os.Exit(1)
	}
	
	command := args[0]
	switch command {
	case "run":
		handleRun(orch, args[1:])
	default:
		fmt.Printf("Unknown command: %%s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf("%s - %s\n\n", "%s", "%s")
	fmt.Println("Usage:")
	fmt.Printf("  %s [command] [arguments]\n\n", "%s")
	fmt.Println("Commands:")
	fmt.Println("  run      Execute the main task")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --version  Show version information")
	fmt.Println("  --help     Show this help message")
}

func handleRun(orch *core.Orchestrator, args []string) {
	fmt.Println("Running command...")
	// Implement your logic here
}
`,
		t.generateRuntimeImports(),
		t.config.Version,
		t.config.Name,
		t.config.Version,
		t.generateRuntimeConfigs(),
		t.generateRuntimeRegistrations(),
		t.config.Name,
		t.config.Description,
		t.config.Name,
	)
}

func (t *ProjectTemplate) generateSystemMain() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/griffincancode/polyglot.js/core"
%s
)

func main() {
	// Create configuration
	config := core.DefaultConfig()
	config.App.Name = "%s"
	config.App.Version = "%s"
	
%s
	
	// Create orchestrator
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %%v", err)
	}
	
%s
	
	// Initialize
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	if err := orch.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %%v", err)
	}
	defer orch.Shutdown(context.Background())
	
	log.Println("Starting %s system service...")
	
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Main service loop
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			// Periodic task execution
			handleTick(orch)
		case <-sigChan:
			log.Println("Received shutdown signal, cleaning up...")
			return
		}
	}
}

func handleTick(orch *core.Orchestrator) {
	// Implement your periodic task here
	log.Println("Tick...")
}
`,
		t.generateRuntimeImports(),
		t.config.Name,
		t.config.Version,
		t.generateRuntimeConfigs(),
		t.generateRuntimeRegistrations(),
		t.config.Name,
	)
}

func (t *ProjectTemplate) generateMinimalMain() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"log"
	"time"
	
	"github.com/griffincancode/polyglot.js/core"
%s
)

func main() {
	config := core.DefaultConfig()
	config.App.Name = "%s"
	
%s
	
	orch, err := core.NewOrchestrator(config)
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %%v", err)
	}
	
%s
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	if err := orch.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %%v", err)
	}
	defer orch.Shutdown(context.Background())
	
	log.Println("Application started successfully")
	
	// Your code here
}
`,
		t.generateRuntimeImports(),
		t.config.Name,
		t.generateRuntimeConfigs(),
		t.generateRuntimeRegistrations(),
	)
}

func (t *ProjectTemplate) generateRuntimeImports() string {
	var imports []string
	for _, lang := range t.config.Languages {
		switch lang {
		case "python":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/python"`)
		case "javascript":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/javascript"`)
		case "rust":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/rust"`)
		case "cpp":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/cpp"`)
		case "java":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/java"`)
		case "ruby":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/ruby"`)
		case "php":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/php"`)
		case "lua":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/lua"`)
		case "wasm":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/wasm"`)
		case "zig":
			imports = append(imports, `	"github.com/griffincancode/polyglot.js/runtimes/zig"`)
		}
	}
	return strings.Join(imports, "\n")
}

func (t *ProjectTemplate) generateRuntimeConfigs() string {
	var configs []string
	for _, lang := range t.config.Languages {
		version := "latest"
		if lang == "python" && t.config.PythonVersion != "" {
			version = t.config.PythonVersion
		}
		configs = append(configs, fmt.Sprintf(`	config.EnableRuntime("%s", "%s")`, lang, version))
	}
	return strings.Join(configs, "\n")
}

func (t *ProjectTemplate) generateRuntimeRegistrations() string {
	var registrations []string
	for _, lang := range t.config.Languages {
		switch lang {
		case "python":
			registrations = append(registrations, `	orch.RegisterRuntime(python.NewRuntime())`)
		case "javascript":
			registrations = append(registrations, `	orch.RegisterRuntime(javascript.NewRuntime())`)
		case "rust":
			registrations = append(registrations, `	orch.RegisterRuntime(rust.NewRuntime())`)
		case "cpp":
			registrations = append(registrations, `	orch.RegisterRuntime(cpp.NewRuntime())`)
		case "java":
			registrations = append(registrations, `	orch.RegisterRuntime(java.NewRuntime())`)
		case "ruby":
			registrations = append(registrations, `	orch.RegisterRuntime(ruby.NewRuntime())`)
		case "php":
			registrations = append(registrations, `	orch.RegisterRuntime(php.NewRuntime())`)
		case "lua":
			registrations = append(registrations, `	orch.RegisterRuntime(lua.NewRuntime())`)
		case "wasm":
			registrations = append(registrations, `	orch.RegisterRuntime(wasm.NewRuntime())`)
		case "zig":
			registrations = append(registrations, `	orch.RegisterRuntime(zig.NewRuntime())`)
		}
	}
	return strings.Join(registrations, "\n")
}

func (t *ProjectTemplate) generateBridgeFunctions() string {
	return `	// Example: Greet function
	bridge.Register("greet", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		name := "World"
		if len(args) > 0 {
			if n, ok := args[0].(string); ok {
				name = n
			}
		}
		return map[string]interface{}{
			"message": "Hello, " + name + "!",
			"timestamp": time.Now().Unix(),
		}, nil
	})
	
	// Example: Get app info
	bridge.Register("getAppInfo", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		return map[string]interface{}{
			"name": config.App.Name,
			"version": config.App.Version,
			"description": config.App.Description,
		}, nil
	})`
}

func (t *ProjectTemplate) formatStringArray(items []string) string {
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf(`"%s"`, item)
	}
	return strings.Join(quoted, ", ")
}
