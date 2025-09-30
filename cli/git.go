package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GitManager handles git initialization
type GitManager struct {
	projectPath string
}

// NewGitManager creates a new git manager
func NewGitManager(projectPath string) *GitManager {
	return &GitManager{projectPath: projectPath}
}

// Initialize initializes a git repository
func (g *GitManager) Initialize() error {
	fmt.Println("\n🔧 Initializing Git repository...")

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(g.projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Initialize git
	cmd := exec.Command("git", "init")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize git: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("  ✅ Git repository initialized")

	// Create initial commit
	cmd = exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Warning: failed to stage files")
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit from Polyglot CLI")
	if err := cmd.Run(); err != nil {
		fmt.Println("  ⚠️  Warning: failed to create initial commit")
	} else {
		fmt.Println("  ✅ Created initial commit")
	}

	return nil
}

// GenerateGitignore creates a .gitignore file
func (g *GitManager) GenerateGitignore(languages []string) error {
	content := `# Polyglot Project
.polyglot/
dist/
*.log

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`

	// Add language-specific ignores
	for _, lang := range languages {
		switch lang {
		case "python":
			content += `
# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
*.egg-info/
.installed.cfg
*.egg
venv/
env/
ENV/
.venv
`
		case "javascript":
			content += `
# JavaScript/Node
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.pnpm-debug.log*
.npm
.eslintcache
.node_repl_history
*.tgz
.yarn-integrity
.env.local
.env.development.local
.env.test.local
.env.production.local
`
		case "rust":
			content += `
# Rust
target/
Cargo.lock
**/*.rs.bk
`
		case "java":
			content += `
# Java
*.class
*.jar
*.war
*.ear
.gradle/
build/
.mtj.tmp/
hs_err_pid*
`
		}
	}

	path := filepath.Join(g.projectPath, ".gitignore")
	return os.WriteFile(path, []byte(content), 0644)
}
