package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const version = "0.1.0"

func handleInit(args []string) {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║                                                       ║")
	fmt.Println("║       POLYGLOT PROJECT INITIALIZATION WIZARD          ║")
	fmt.Println("║                                                       ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	fmt.Println()

	var config *ProjectConfig
	var err error

	// Check if non-interactive mode
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		// Quick mode: just provide name
		config = &ProjectConfig{
			Name:            args[0],
			Description:     "A Polyglot desktop application",
			Version:         "0.1.0",
			License:         "MIT",
			Template:        "webapp",
			Languages:       []string{"python", "javascript"},
			Features:        []string{"webview", "hmr"},
			WindowWidth:     1280,
			WindowHeight:    720,
			WindowResizable: true,
			DevTools:        true,
			GitInit:         true,
			PythonVersion:   "3.11",
			PackageManager:  "npm",
		}
		fmt.Printf("🚀 Quick initialization mode for project: %s\n", config.Name)
		fmt.Println("   (Use 'polyglot init' without arguments for interactive mode)")
		fmt.Println()
	} else {
		// Interactive mode
		wizard := NewWizard()
		config, err = wizard.Run()
		if err != nil {
			fmt.Printf("❌ Error running wizard: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if project already exists
	if _, err := os.Stat(config.Name); !os.IsNotExist(err) {
		fmt.Printf("\n❌ Error: Directory '%s' already exists!\n", config.Name)
		fmt.Println("   Please choose a different project name or remove the existing directory.")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("📋 Project Summary")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("  Name:        %s\n", config.Name)
	fmt.Printf("  Description: %s\n", config.Description)
	fmt.Printf("  Version:     %s\n", config.Version)
	fmt.Printf("  License:     %s\n", config.License)
	fmt.Printf("  Template:    %s\n", config.Template)
	fmt.Printf("  Languages:   %s\n", strings.Join(config.Languages, ", "))
	fmt.Printf("  Features:    %s\n", strings.Join(config.Features, ", "))
	if config.Author != "" {
		fmt.Printf("  Author:      %s\n", config.Author)
	}
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()

	// Check dependencies
	depManager := NewDependencyManager(config)
	if err := depManager.DetectAndGuide(); err != nil {
		fmt.Printf("\n❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Generate project
	fmt.Println()
	fmt.Println("🔨 Generating project structure...")
	fmt.Println()

	template := NewTemplate(config)
	if err := template.Generate(); err != nil {
		fmt.Printf("❌ Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("  ✅ Created directory structure")
	fmt.Println("  ✅ Generated configuration files")
	fmt.Println("  ✅ Created main application file")

	if contains(config.Features, "webview") {
		fmt.Println("  ✅ Generated frontend files")
	}

	fmt.Println("  ✅ Generated README.md")
	fmt.Println("  ✅ Generated .gitignore")

	if config.License != "Custom" {
		fmt.Println("  ✅ Generated LICENSE")
	}

	fmt.Println("  ✅ Generated go.mod")
	fmt.Println("  ✅ Generated Makefile")

	// Generate package files
	if err := depManager.GeneratePackageFiles(); err != nil {
		fmt.Printf("  ⚠️  Warning: failed to generate some package files: %v\n", err)
	} else {
		for _, lang := range config.Languages {
			switch lang {
			case "python":
				fmt.Println("  ✅ Generated requirements.txt")
			case "javascript":
				fmt.Println("  ✅ Generated package.json")
			case "rust":
				fmt.Println("  ✅ Generated Cargo.toml")
			}
		}
	}

	// Initialize Git
	if config.GitInit {
		gitManager := NewGitManager(config.Name)
		if err := gitManager.Initialize(); err != nil {
			fmt.Printf("  ⚠️  Warning: failed to initialize git: %v\n", err)
		}
	}

	// Initialize Go module
	fmt.Println()
	fmt.Println("📦 Initializing Go module...")
	originalDir, _ := os.Getwd()
	os.Chdir(config.Name)

	cmd := exec.Command("go", "mod", "init", config.Name)
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Warning: failed to initialize go module")
	} else {
		fmt.Println("  ✅ Go module initialized")
	}

	// Try to download dependencies (may fail if not online or polyglot not published yet)
	fmt.Println()
	fmt.Println("📥 Attempting to download dependencies...")
	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Could not download dependencies (you may need to run 'go mod tidy' manually)")
	} else {
		fmt.Println("  ✅ Dependencies downloaded")
	}

	os.Chdir(originalDir)

	// Success!
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║                                                       ║")
	fmt.Println("║              ✨ PROJECT CREATED SUCCESSFULLY! ✨       ║")
	fmt.Println("║                                                       ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("🎉 Your Polyglot project is ready!")
	fmt.Println()
	fmt.Println("📚 Next steps:")
	fmt.Println()
	fmt.Printf("   1. Navigate to your project:\n")
	fmt.Printf("      $ cd %s\n", config.Name)
	fmt.Println()
	fmt.Println("   2. Review and customize the configuration:")
	fmt.Println("      $ cat polyglot.config.json")
	fmt.Println()
	fmt.Println("   3. Install dependencies:")
	fmt.Println("      $ make install")
	fmt.Println()
	fmt.Println("   4. Build your application:")
	fmt.Println("      $ make build")
	fmt.Println()
	fmt.Println("   5. Run in development mode:")
	fmt.Println("      $ make dev")
	fmt.Println()
	fmt.Println("   Or use Polyglot CLI commands:")
	fmt.Println("      $ polyglot dev")
	fmt.Println("      $ polyglot build")
	fmt.Println()
	fmt.Println("📖 Documentation:")
	fmt.Println("   - README.md in your project directory")
	fmt.Println("   - https://github.com/griffincancode/polyglot.js")
	fmt.Println()
	fmt.Println("💡 Tip: Check out the generated README.md for detailed information!")
	fmt.Println()
}

func handleBuild(args []string) {
	fmt.Println("🔨 Building application...")

	// Check if we're in a project directory
	if _, err := os.Stat("polyglot.config.json"); os.IsNotExist(err) {
		fmt.Println("❌ Error: Not a Polyglot project directory")
		fmt.Println("   Run this command from your project root, or initialize a new project with 'polyglot init'")
		os.Exit(1)
	}

	// Parse arguments for platform and arch
	platform := ""
	arch := ""
	for i, arg := range args {
		if arg == "--platform" && i+1 < len(args) {
			platform = args[i+1]
		}
		if arg == "--arch" && i+1 < len(args) {
			arch = args[i+1]
		}
	}

	// Build command
	var cmd *exec.Cmd
	if platform != "" && arch != "" {
		fmt.Printf("Building for %s/%s...\n", platform, arch)
		cmd = exec.Command("make", fmt.Sprintf("build-%s", platform))
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOARCH=%s", arch))
	} else {
		cmd = exec.Command("make", "build")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Build complete!")
}

func handleDev(args []string) {
	fmt.Println("🔧 Starting development mode...")

	// Check if we're in a project directory
	if _, err := os.Stat("polyglot.config.json"); os.IsNotExist(err) {
		fmt.Println("❌ Error: Not a Polyglot project directory")
		os.Exit(1)
	}

	// Check for HMR feature
	if _, err := os.Stat("src/backend/main.go"); os.IsNotExist(err) {
		fmt.Println("❌ Error: Could not find src/backend/main.go")
		os.Exit(1)
	}

	fmt.Println("🔥 Hot reload enabled")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	cmd := exec.Command("make", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// If make dev fails, try direct go run
		fmt.Println("Falling back to direct execution...")
		cmd = exec.Command("go", "run", "src/backend/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Failed to start: %v\n", err)
			os.Exit(1)
		}
	}
}

func handleTest(args []string) {
	fmt.Println("🧪 Running tests...")

	if _, err := os.Stat("polyglot.config.json"); os.IsNotExist(err) {
		fmt.Println("❌ Error: Not a Polyglot project directory")
		os.Exit(1)
	}

	cmd := exec.Command("make", "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback to go test
		cmd = exec.Command("go", "test", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Tests failed: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("✅ All tests passed!")
}

func handleVersion(args []string) {
	fmt.Printf("Polyglot CLI v%s\n", version)
	fmt.Println()
	fmt.Println("A modern desktop application framework supporting multiple")
	fmt.Println("language runtimes with native UI and zero-copy memory sharing.")
	fmt.Println()
	fmt.Println("Documentation: https://github.com/griffincancode/polyglot.js")
}
