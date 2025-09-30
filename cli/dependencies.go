package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DependencyManager handles dependency detection and installation
type DependencyManager struct {
	config *ProjectConfig
}

// NewDependencyManager creates a new dependency manager
func NewDependencyManager(config *ProjectConfig) *DependencyManager {
	return &DependencyManager{config: config}
}

// DetectAndGuide detects installed dependencies and provides guidance
func (d *DependencyManager) DetectAndGuide() error {
	fmt.Println("\n📦 Checking dependencies...")
	fmt.Println("=" + strings.Repeat("=", 40))

	missing := []string{}
	warnings := []string{}

	// Check Go
	if !d.checkGo() {
		missing = append(missing, "Go 1.21+")
	}

	// Check language-specific dependencies
	for _, lang := range d.config.Languages {
		switch lang {
		case "python":
			if !d.checkPython() {
				warnings = append(warnings, fmt.Sprintf("Python %s not found", d.config.PythonVersion))
			}
		case "javascript":
			if !d.checkNode() {
				warnings = append(warnings, "Node.js not found")
			}
			if !d.checkPackageManager() {
				warnings = append(warnings, fmt.Sprintf("%s not found", d.config.PackageManager))
			}
		case "rust":
			if !d.checkRust() {
				warnings = append(warnings, "Rust not found")
			}
		case "java":
			if !d.checkJava() {
				warnings = append(warnings, "Java JDK not found")
			}
		case "ruby":
			if !d.checkRuby() {
				warnings = append(warnings, "Ruby not found")
			}
		case "php":
			if !d.checkPHP() {
				warnings = append(warnings, "PHP not found")
			}
		}
	}

	// Check Git if needed
	if d.config.GitInit {
		if !d.checkGit() {
			warnings = append(warnings, "Git not found")
		}
	}

	// Print results
	if len(missing) == 0 && len(warnings) == 0 {
		fmt.Println("✅ All dependencies are installed!")
		return nil
	}

	if len(missing) > 0 {
		fmt.Println("\n❌ Required dependencies missing:")
		for _, dep := range missing {
			fmt.Printf("  - %s\n", dep)
		}
		fmt.Println("\nPlease install these dependencies before continuing.")
		return fmt.Errorf("missing required dependencies")
	}

	if len(warnings) > 0 {
		fmt.Println("\n⚠️  Optional dependencies missing:")
		for _, warning := range warnings {
			fmt.Printf("  - %s\n", warning)
		}
		fmt.Println("\nThe project will be created, but some features may not work.")
		fmt.Println("See the README for installation instructions.")
	}

	return nil
}

func (d *DependencyManager) checkGo() bool {
	cmd := exec.Command("go", "version")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ❌ Go not found")
		fmt.Println("     Install from: https://golang.org/dl/")
		return false
	}
	output, _ := cmd.CombinedOutput()
	fmt.Printf("  ✅ Go found: %s", string(output))
	return true
}

func (d *DependencyManager) checkPython() bool {
	// Try python3 first, then python
	for _, cmd := range []string{"python3", "python"} {
		if out, err := exec.Command(cmd, "--version").CombinedOutput(); err == nil {
			version := strings.TrimSpace(string(out))
			fmt.Printf("  ✅ Python found: %s\n", version)
			return true
		}
	}
	fmt.Println("  ⚠️  Python not found")
	fmt.Println("     Install from: https://www.python.org/downloads/")
	return false
}

func (d *DependencyManager) checkNode() bool {
	cmd := exec.Command("node", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Node.js not found")
		fmt.Println("     Install from: https://nodejs.org/")
		return false
	}
	output, _ := cmd.CombinedOutput()
	fmt.Printf("  ✅ Node.js found: %s", string(output))
	return true
}

func (d *DependencyManager) checkPackageManager() bool {
	pm := d.config.PackageManager
	if pm == "" {
		pm = "npm"
	}

	cmd := exec.Command(pm, "--version")
	if err := cmd.Run(); err != nil {
		fmt.Printf("  ⚠️  %s not found\n", pm)
		if pm != "npm" {
			fmt.Printf("     Install from: https://classic.yarnpkg.com/ or https://pnpm.io/\n")
		}
		return false
	}
	output, _ := cmd.CombinedOutput()
	fmt.Printf("  ✅ %s found: %s", pm, string(output))
	return true
}

func (d *DependencyManager) checkRust() bool {
	cmd := exec.Command("rustc", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Rust not found")
		fmt.Println("     Install from: https://rustup.rs/")
		return false
	}
	output, _ := cmd.CombinedOutput()
	fmt.Printf("  ✅ Rust found: %s", string(output))
	return true
}

func (d *DependencyManager) checkJava() bool {
	cmd := exec.Command("java", "-version")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Java not found")
		fmt.Println("     Install from: https://adoptium.net/")
		return false
	}
	output, _ := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		fmt.Printf("  ✅ Java found: %s\n", lines[0])
	}
	return true
}

func (d *DependencyManager) checkRuby() bool {
	cmd := exec.Command("ruby", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Ruby not found")
		fmt.Println("     Install from: https://www.ruby-lang.org/")
		return false
	}
	output, _ := cmd.CombinedOutput()
	fmt.Printf("  ✅ Ruby found: %s", string(output))
	return true
}

func (d *DependencyManager) checkPHP() bool {
	cmd := exec.Command("php", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  PHP not found")
		fmt.Println("     Install from: https://www.php.net/downloads")
		return false
	}
	output, _ := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		fmt.Printf("  ✅ PHP found: %s\n", lines[0])
	}
	return true
}

func (d *DependencyManager) checkGit() bool {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Git not found")
		fmt.Println("     Install from: https://git-scm.com/downloads")
		return false
	}
	output, _ := cmd.CombinedOutput()
	fmt.Printf("  ✅ Git found: %s", string(output))
	return true
}

// GeneratePackageFiles creates language-specific package management files
func (d *DependencyManager) GeneratePackageFiles() error {
	for _, lang := range d.config.Languages {
		switch lang {
		case "python":
			if err := d.generatePythonRequirements(); err != nil {
				return err
			}
		case "javascript":
			if err := d.generatePackageJSON(); err != nil {
				return err
			}
		case "rust":
			if err := d.generateCargoToml(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DependencyManager) generatePythonRequirements() error {
	content := `# Python dependencies for ` + d.config.Name + `
# Add your packages here, for example:
# requests>=2.28.0
# numpy>=1.24.0
`

	path := filepath.Join(d.config.Name, "requirements.txt")
	return os.WriteFile(path, []byte(content), 0644)
}

func (d *DependencyManager) generatePackageJSON() error {
	scripts := `"scripts": {
    "dev": "polyglot dev",
    "build": "polyglot build",
    "test": "polyglot test"
  }`

	content := fmt.Sprintf(`{
  "name": "%s",
  "version": "%s",
  "description": "%s",
  "author": "%s",
  "license": "%s",
  %s,
  "dependencies": {},
  "devDependencies": {}
}
`,
		d.config.Name,
		d.config.Version,
		d.config.Description,
		d.config.Author,
		d.config.License,
		scripts,
	)

	path := filepath.Join(d.config.Name, "package.json")
	return os.WriteFile(path, []byte(content), 0644)
}

func (d *DependencyManager) generateCargoToml() error {
	content := fmt.Sprintf(`[package]
name = "%s"
version = "%s"
edition = "2021"

[dependencies]
`,
		strings.ReplaceAll(d.config.Name, "-", "_"),
		d.config.Version,
	)

	path := filepath.Join(d.config.Name, "src", "rust", "Cargo.toml")
	return os.WriteFile(path, []byte(content), 0644)
}
