# Polyglot CLI

The official command-line interface for creating and managing Polyglot desktop applications.

## Overview

The Polyglot CLI provides a comprehensive project initialization wizard with:

- 🧙‍♂️ **Interactive Project Wizard** - Guided setup with smart defaults
- 📦 **Multiple Templates** - Web app, CLI tool, system utility, desktop app, or minimal
- 🌐 **Multi-Language Support** - 11 language runtimes to choose from
- 🔍 **Dependency Detection** - Automatic checking and installation guidance
- ⚙️ **Configuration Generation** - Project-specific config files
- 🔧 **Git Integration** - Automatic repository initialization
- 📄 **Smart Documentation** - Auto-generated README and LICENSE
- 🛠️ **Build System** - Complete Makefile and package management

## Installation

Build the CLI from source:

```bash
cd cli
go build -o polyglot
sudo mv polyglot /usr/local/bin/  # Optional: make globally available
```

Or install with Go:

```bash
go install github.com/griffincancode/polyglot.js/cli@latest
```

## Commands

### `polyglot init [project-name]`

Initialize a new Polyglot project.

**Interactive Mode:**
```bash
polyglot init
```

Launches the full interactive wizard that guides you through:
- Project name and metadata
- License selection (MIT, Apache-2.0, GPL-3.0, BSD-3-Clause, Unlicense)
- Template choice (webapp, CLI, system utility, desktop app, minimal)
- Language runtime selection (python, javascript, go, rust, cpp, java, ruby, php, lua, wasm, zig)
- Feature selection (webview, HMR, cloud, marketplace, security, signing)
- Git initialization
- Package manager preference
- Webview configuration

**Quick Mode:**
```bash
polyglot init my-app
```

Creates a project with sensible defaults:
- Template: Web Application
- Languages: Python + JavaScript
- Features: Webview + HMR
- License: MIT

### `polyglot build [--platform PLATFORM] [--arch ARCH]`

Build your Polyglot application.

```bash
# Build for current platform
polyglot build

# Build for specific platform
polyglot build --platform darwin --arch arm64
polyglot build --platform linux --arch amd64
polyglot build --platform windows --arch amd64
```

### `polyglot dev`

Start development mode with hot module reload.

```bash
polyglot dev
```

Features:
- Automatic recompilation on file changes
- Frontend hot reload (if HMR enabled)
- Real-time error reporting
- DevTools enabled

### `polyglot test`

Run project tests.

```bash
polyglot test
```

### `polyglot version`

Display CLI version information.

```bash
polyglot version
```

## Project Templates

### Web Application (Default)

Full-featured web application with native webview UI.

**Includes:**
- Go backend with orchestrator
- HTML/CSS/JS frontend
- WebView integration
- Bridge for frontend-backend communication
- Example API endpoints
- Modern, responsive UI

**Best for:** Desktop applications with rich UIs, productivity tools, content management

### CLI Tool

Command-line utility template.

**Includes:**
- Argument parsing with flags
- Help and version commands
- Structured command handling
- Multi-runtime support

**Best for:** Developer tools, automation scripts, system utilities

### System Utility

Background service or daemon template.

**Includes:**
- Signal handling for graceful shutdown
- Periodic task execution
- Service loop
- Configuration management

**Best for:** Background processors, monitoring tools, system services

### Desktop App

Cross-platform desktop application.

**Includes:**
- Full webview integration
- Native menus and dialogs
- System tray support
- Multi-window management

**Best for:** Native-feeling desktop applications, Electron alternatives

### Minimal

Bare-bones starter with essential structure only.

**Includes:**
- Basic Go entry point
- Runtime initialization
- Minimal configuration

**Best for:** Custom applications, learning, experimentation

## Supported Languages

| Language | Runtime Version | Use Cases |
|----------|----------------|-----------|
| Python | 3.8+ | Data processing, ML, scripting |
| JavaScript | Node 18+ | Web logic, async operations |
| Go | 1.21+ | Main application, performance |
| Rust | 1.70+ | Systems programming, performance |
| C++ | C++17 | Native integrations, legacy code |
| Java | JDK 11+ | Enterprise integrations |
| Ruby | 3.0+ | Scripting, DSLs |
| PHP | 8.0+ | Web logic, legacy systems |
| Lua | 5.4+ | Embedded scripting |
| WebAssembly | - | Portable code, sandboxing |
| Zig | 0.11+ | Systems programming |

## Features

### Core Features

- **Webview**: Native webview for cross-platform UI
- **HMR**: Hot module reload for rapid development
- **Cloud**: Cloud integration and sync
- **Marketplace**: Plugin marketplace support
- **Security**: Enhanced security sandbox
- **Signing**: Code signing for distribution

### Generated Files

Every project includes:

```
my-app/
├── src/
│   ├── backend/
│   │   └── main.go           # Main application
│   └── frontend/             # UI files (if webview enabled)
│       ├── index.html
│       ├── styles/
│       │   └── main.css
│       └── scripts/
│           └── main.js
├── dist/                     # Build outputs
├── .polyglot/               # Internal state
├── polyglot.config.json     # Project configuration
├── go.mod                   # Go dependencies
├── Makefile                 # Build automation
├── .gitignore              # Git ignore rules
├── LICENSE                 # License file
└── README.md               # Project documentation
```

### Language-Specific Files

**Python:**
- `requirements.txt` - Python package dependencies

**JavaScript:**
- `package.json` - NPM/Yarn/PNPM configuration

**Rust:**
- `Cargo.toml` - Rust package manifest

## Dependency Detection

The CLI automatically detects and reports on:

- ✅ Go installation and version
- ⚠️ Language runtime availability
- ⚠️ Package manager installation
- ⚠️ Git availability

Missing optional dependencies are reported with installation guidance but won't block project creation.

## Configuration

### polyglot.config.json

The generated configuration file controls all aspects of your application:

```json
{
  "name": "my-app",
  "version": "0.1.0",
  "description": "My Polyglot app",
  "languages": ["python", "javascript"],
  "features": ["webview", "hmr"],
  "webview": {
    "width": 1280,
    "height": 720,
    "resizable": true,
    "devTools": true
  },
  "memory": {
    "maxSharedMemory": 1073741824,
    "enableZeroCopy": true,
    "gcInterval": "5m"
  },
  "runtimes": {
    "python": {
      "version": "3.11",
      "maxConcurrency": 10,
      "timeout": "30s"
    }
  }
}
```

## Examples

### Create a Python + JavaScript Web App

```bash
polyglot init

# Follow prompts:
# - Name: calculator-app
# - Template: Web Application
# - Languages: python, javascript
# - Features: webview, hmr
# - Git: yes

cd calculator-app
make install
make dev
```

### Create a Multi-Language CLI Tool

```bash
polyglot init

# Follow prompts:
# - Name: data-processor
# - Template: CLI Tool
# - Languages: python, rust, go
# - Features: security

cd data-processor
make build
./dist/data-processor --help
```

### Quick Start with Defaults

```bash
polyglot init my-quick-app
cd my-quick-app
make build
make run
```

## Makefile Commands

Generated projects include these Makefile targets:

- `make build` - Build application
- `make build-all` - Build for all platforms
- `make build-darwin` - Build for macOS
- `make build-linux` - Build for Linux
- `make build-windows` - Build for Windows
- `make run` - Build and run
- `make dev` - Development mode
- `make test` - Run tests
- `make install` - Install dependencies
- `make clean` - Clean build artifacts
- `make help` - Show help

## Troubleshooting

### "command not found: polyglot"

Ensure the binary is in your PATH:

```bash
export PATH=$PATH:/usr/local/bin
# Or add to ~/.bashrc or ~/.zshrc
```

### "Not a Polyglot project directory"

Commands like `build`, `dev`, and `test` must be run from a project root directory (containing `polyglot.config.json`).

### Dependency Detection Fails

If dependency detection hangs or fails:

1. Ensure you have network connectivity
2. Check that language runtimes are properly installed
3. Verify versions with `python --version`, `node --version`, etc.

### Build Errors

If builds fail:

1. Run `go mod tidy` to resolve Go dependencies
2. Check that all required runtimes are installed
3. Review error messages for missing packages
4. Ensure you're using compatible versions

## Development

### Project Structure

```
cli/
├── main.go           # CLI entry point
├── commands.go       # Command handlers
├── wizard.go         # Interactive wizard
├── templates.go      # Template generation
├── generator.go      # File generators
├── dependencies.go   # Dependency detection
├── git.go           # Git operations
├── licenses.go      # License generation
├── makefile_gen.go  # Makefile generation
└── types.go         # Shared types
```

### Adding a New Template

1. Add template type to `wizard.go` `parseTemplate()`
2. Create generation function in `templates.go`
3. Add template-specific logic in `generateMain()`
4. Update README generation for new template

### Adding a New Language

1. Add language to `wizard.go` `parseLanguages()`
2. Add import in `templates.go` `generateRuntimeImports()`
3. Add registration in `generateRuntimeRegistrations()`
4. Add dependency check in `dependencies.go`
5. Add package file generation if needed

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

- GitHub Issues: https://github.com/griffincancode/polyglot.js/issues
- Documentation: https://github.com/griffincancode/polyglot.js
- Discussions: https://github.com/griffincancode/polyglot.js/discussions
