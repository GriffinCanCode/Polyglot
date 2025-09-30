package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Frontend file generation methods

func (t *ProjectTemplate) generateFrontend() error {
	// Generate index.html
	if err := t.generateIndexHTML(); err != nil {
		return err
	}

	// Generate CSS
	if err := t.generateCSS(); err != nil {
		return err
	}

	// Generate JavaScript
	if err := t.generateJavaScript(); err != nil {
		return err
	}

	return nil
}

func (t *ProjectTemplate) generateIndexHTML() error {
	content := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<link rel="stylesheet" href="styles/main.css">
</head>
<body>
	<div class="app">
		<header class="app-header">
			<h1>%s</h1>
			<p class="subtitle">%s</p>
		</header>
		
		<main class="app-main">
			<div class="card">
				<h2>Welcome to Polyglot</h2>
				<p>Your application is running successfully!</p>
				
				<div class="actions">
					<button id="greetBtn" class="btn btn-primary">Say Hello</button>
					<button id="infoBtn" class="btn btn-secondary">Get Info</button>
				</div>
				
				<div id="output" class="output"></div>
			</div>
			
			<div class="info-grid">
				<div class="info-card">
					<h3>🚀 Multi-Language</h3>
					<p>Enabled: %s</p>
				</div>
				<div class="info-card">
					<h3>⚡ Features</h3>
					<p>%s</p>
				</div>
				<div class="info-card">
					<h3>📦 Version</h3>
					<p>%s</p>
				</div>
			</div>
		</main>
		
		<footer class="app-footer">
			<p>Powered by Polyglot Framework</p>
		</footer>
	</div>
	
	<script src="scripts/main.js"></script>
</body>
</html>`,
		t.config.Name,
		t.config.Name,
		t.config.Description,
		strings.Join(t.config.Languages, ", "),
		strings.Join(t.config.Features, ", "),
		t.config.Version,
	)

	path := filepath.Join(t.config.Name, "src", "frontend", "index.html")
	return os.WriteFile(path, []byte(content), 0644)
}

func (t *ProjectTemplate) generateCSS() error {
	content := `* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

body {
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 
		'Helvetica Neue', Arial, sans-serif;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: #333;
	line-height: 1.6;
}

.app {
	min-height: 100vh;
	display: flex;
	flex-direction: column;
}

.app-header {
	background: rgba(255, 255, 255, 0.95);
	padding: 2rem;
	text-align: center;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.app-header h1 {
	font-size: 2.5rem;
	color: #667eea;
	margin-bottom: 0.5rem;
}

.subtitle {
	color: #666;
	font-size: 1.1rem;
}

.app-main {
	flex: 1;
	padding: 2rem;
	max-width: 1200px;
	width: 100%;
	margin: 0 auto;
}

.card {
	background: white;
	border-radius: 12px;
	padding: 2rem;
	margin-bottom: 2rem;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.card h2 {
	color: #667eea;
	margin-bottom: 1rem;
}

.actions {
	display: flex;
	gap: 1rem;
	margin: 1.5rem 0;
}

.btn {
	padding: 0.75rem 1.5rem;
	font-size: 1rem;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	transition: all 0.3s ease;
	font-weight: 500;
}

.btn-primary {
	background: #667eea;
	color: white;
}

.btn-primary:hover {
	background: #5568d3;
	transform: translateY(-2px);
	box-shadow: 0 4px 8px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
	background: #764ba2;
	color: white;
}

.btn-secondary:hover {
	background: #63408a;
	transform: translateY(-2px);
	box-shadow: 0 4px 8px rgba(118, 75, 162, 0.4);
}

.output {
	margin-top: 1.5rem;
	padding: 1rem;
	background: #f5f5f5;
	border-radius: 8px;
	min-height: 60px;
	font-family: 'Courier New', monospace;
	white-space: pre-wrap;
	word-break: break-word;
}

.output:empty::before {
	content: 'Output will appear here...';
	color: #999;
}

.info-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
	gap: 1rem;
}

.info-card {
	background: white;
	border-radius: 12px;
	padding: 1.5rem;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	transition: transform 0.3s ease;
}

.info-card:hover {
	transform: translateY(-4px);
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.info-card h3 {
	color: #667eea;
	margin-bottom: 0.5rem;
	font-size: 1.2rem;
}

.info-card p {
	color: #666;
}

.app-footer {
	background: rgba(255, 255, 255, 0.95);
	padding: 1rem;
	text-align: center;
	color: #666;
	font-size: 0.9rem;
}

/* Loading animation */
@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}

.loading::after {
	content: '';
	display: inline-block;
	width: 1rem;
	height: 1rem;
	border: 2px solid #667eea;
	border-top-color: transparent;
	border-radius: 50%;
	animation: spin 0.6s linear infinite;
	margin-left: 0.5rem;
}
`

	path := filepath.Join(t.config.Name, "src", "frontend", "styles", "main.css")
	return os.WriteFile(path, []byte(content), 0644)
}

func (t *ProjectTemplate) generateJavaScript() error {
	content := `// Polyglot Bridge API
// Access backend functions via window.polyglot.call(functionName, ...args)

document.addEventListener('DOMContentLoaded', () => {
	setupEventListeners();
	checkPolyglotAPI();
});

function setupEventListeners() {
	const greetBtn = document.getElementById('greetBtn');
	const infoBtn = document.getElementById('infoBtn');
	
	if (greetBtn) {
		greetBtn.addEventListener('click', handleGreet);
	}
	
	if (infoBtn) {
		infoBtn.addEventListener('click', handleGetInfo);
	}
}

function checkPolyglotAPI() {
	if (typeof window.polyglot === 'undefined') {
		displayError('Polyglot API not available. Make sure the app is running in the native webview.');
		return false;
	}
	return true;
}

async function handleGreet() {
	if (!checkPolyglotAPI()) return;
	
	const output = document.getElementById('output');
	output.textContent = 'Calling backend...';
	output.classList.add('loading');
	
	try {
		const result = await window.polyglot.call('greet', 'Polyglot User');
		output.classList.remove('loading');
		output.textContent = JSON.stringify(result, null, 2);
	} catch (error) {
		output.classList.remove('loading');
		displayError('Error calling greet: ' + error.message);
	}
}

async function handleGetInfo() {
	if (!checkPolyglotAPI()) return;
	
	const output = document.getElementById('output');
	output.textContent = 'Getting app info...';
	output.classList.add('loading');
	
	try {
		const result = await window.polyglot.call('getAppInfo');
		output.classList.remove('loading');
		output.textContent = JSON.stringify(result, null, 2);
	} catch (error) {
		output.classList.remove('loading');
		displayError('Error getting info: ' + error.message);
	}
}

function displayError(message) {
	const output = document.getElementById('output');
	output.textContent = '❌ ' + message;
	output.style.color = '#e74c3c';
}

// Example: Listen for events from backend
if (typeof window.polyglot !== 'undefined' && window.polyglot.on) {
	window.polyglot.on('notification', (data) => {
		console.log('Received notification:', data);
		// Handle notification
	});
}
`

	path := filepath.Join(t.config.Name, "src", "frontend", "scripts", "main.js")
	return os.WriteFile(path, []byte(content), 0644)
}

func (t *ProjectTemplate) generateReadme() error {
	languageSetup := t.generateLanguageSetup()
	buildInstructions := t.generateBuildInstructions()

	content := fmt.Sprintf(`# %s

%s

## Overview

This is a Polyglot desktop application built with multiple language runtimes.

**Version:** %s  
**License:** %s  
**Author:** %s  
**Template:** %s

## Features

- 🚀 Multi-language support: %s
- ⚡ Enabled features: %s
- 🎨 Native webview UI
- 🔧 Hot module reload support
- 📦 Zero-copy memory sharing

## Prerequisites

%s

## Installation

1. Clone or navigate to this project:
   `+"```"+`bash
   cd %s
   `+"```"+`

2. Install Go dependencies:
   `+"```"+`bash
   go mod download
   `+"```"+`

%s

## Building

%s

## Development

Start the development server with hot reload:

`+"```"+`bash
polyglot dev
`+"```"+`

Or manually:

`+"```"+`bash
go run src/backend/main.go
`+"```"+`

## Project Structure

`+"```"+`
%s/
├── src/
│   ├── backend/          # Go backend code
│   │   └── main.go       # Main application entry
│   ├── frontend/         # Web UI files
│   │   ├── index.html
│   │   ├── styles/
│   │   └── scripts/
%s
├── dist/                 # Build outputs
├── polyglot.config.json  # Project configuration
├── go.mod                # Go dependencies
└── README.md
`+"```"+`

## Configuration

Edit `+"`"+`polyglot.config.json`+"`"+` to customize:

- Runtime settings and versions
- Webview dimensions and behavior
- Memory limits and optimization
- Build targets and platforms

## API Reference

### Frontend → Backend Communication

Call backend functions from JavaScript:

`+"```"+`javascript
const result = await window.polyglot.call('functionName', arg1, arg2);
`+"```"+`

### Backend → Frontend Communication

Send events to frontend from Go:

`+"```"+`go
bridge.Emit("eventName", data)
`+"```"+`

## Available Scripts

- `+"`"+`polyglot dev`+"`"+` - Start development mode
- `+"`"+`polyglot build`+"`"+` - Build production binary
- `+"`"+`polyglot test`+"`"+` - Run tests
- `+"`"+`make build`+"`"+` - Build using Makefile
- `+"`"+`make run`+"`"+` - Run the application
- `+"`"+`make clean`+"`"+` - Clean build artifacts

## Deployment

Build for your target platform:

`+"```"+`bash
polyglot build --platform darwin --arch arm64
polyglot build --platform linux --arch amd64
polyglot build --platform windows --arch amd64
`+"```"+`

Binaries will be created in the `+"`"+`dist/`+"`"+` directory.

## Troubleshooting

### Runtime Issues

If you encounter runtime initialization errors:

1. Verify all required language runtimes are installed
2. Check version compatibility in `+"`"+`polyglot.config.json`+"`"+`
3. Review runtime-specific logs in `+"`"+`.polyglot/logs/`+"`"+`

### Webview Issues

If the webview doesn't appear:

1. Ensure webview dependencies are installed for your OS
2. Check DevTools console for JavaScript errors (if enabled)
3. Verify frontend files are in the correct location

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the %s License - see the LICENSE file for details.

## Resources

- [Polyglot Documentation](https://github.com/griffincancode/polyglot.js)
- [Examples](https://github.com/griffincancode/polyglot.js/tree/main/examples)
- [API Reference](https://github.com/griffincancode/polyglot.js/blob/main/docs/API.md)

## Support

For issues and questions:

- GitHub Issues: https://github.com/griffincancode/polyglot.js/issues
- Discussions: https://github.com/griffincancode/polyglot.js/discussions

---

Built with ❤️ using [Polyglot Framework](https://github.com/griffincancode/polyglot.js)
`,
		t.config.Name,
		t.config.Description,
		t.config.Version,
		t.config.License,
		t.config.Author,
		t.config.Template,
		strings.Join(t.config.Languages, ", "),
		strings.Join(t.config.Features, ", "),
		t.generatePrerequisites(),
		t.config.Name,
		languageSetup,
		buildInstructions,
		t.config.Name,
		t.generateProjectStructureExtras(),
		t.config.License,
	)

	path := filepath.Join(t.config.Name, "README.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func (t *ProjectTemplate) generatePrerequisites() string {
	prereqs := []string{
		"- Go 1.21 or higher",
	}

	for _, lang := range t.config.Languages {
		switch lang {
		case "python":
			version := t.config.PythonVersion
			if version == "" {
				version = "3.11"
			}
			prereqs = append(prereqs, fmt.Sprintf("- Python %s or higher", version))
		case "javascript":
			prereqs = append(prereqs, "- Node.js 18+ (for JavaScript runtime)")
		case "rust":
			prereqs = append(prereqs, "- Rust 1.70+ (for Rust runtime)")
		case "java":
			prereqs = append(prereqs, "- Java JDK 11+ (for Java runtime)")
		case "ruby":
			prereqs = append(prereqs, "- Ruby 3.0+ (for Ruby runtime)")
		case "php":
			prereqs = append(prereqs, "- PHP 8.0+ (for PHP runtime)")
		}
	}

	return strings.Join(prereqs, "\n")
}

func (t *ProjectTemplate) generateLanguageSetup() string {
	var setup []string

	for _, lang := range t.config.Languages {
		switch lang {
		case "python":
			setup = append(setup, `
3. Install Python dependencies (if any):
   `+"```"+`bash
   pip install -r requirements.txt  # If you add Python packages
   `+"```")
		case "javascript":
			pm := t.config.PackageManager
			if pm == "" {
				pm = "npm"
			}
			setup = append(setup, fmt.Sprintf(`
3. Install JavaScript dependencies:
   `+"```"+`bash
   %s install
   `+"```", pm))
		}
	}

	return strings.Join(setup, "\n")
}

func (t *ProjectTemplate) generateBuildInstructions() string {
	return `Build the application:

` + "```" + `bash
polyglot build
` + "```" + `

Or use the Makefile:

` + "```" + `bash
make build
` + "```" + `

The binary will be created in the ` + "`dist/`" + ` directory.`
}

func (t *ProjectTemplate) generateProjectStructureExtras() string {
	var extras []string

	for _, lang := range t.config.Languages {
		switch lang {
		case "python":
			extras = append(extras, "│   ├── python/           # Python modules")
		case "javascript":
			extras = append(extras, "│   ├── js/               # JavaScript modules")
		case "rust":
			extras = append(extras, "│   ├── rust/             # Rust crates")
		}
	}

	if len(extras) > 0 {
		return strings.Join(extras, "\n")
	}
	return ""
}
