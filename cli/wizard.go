package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ProjectWizard handles interactive project creation
type ProjectWizard struct {
	scanner *bufio.Scanner
}

// NewWizard creates a new project wizard
func NewWizard() *ProjectWizard {
	return &ProjectWizard{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Run executes the interactive wizard
func (w *ProjectWizard) Run() (*ProjectConfig, error) {
	fmt.Println("🚀 Welcome to Polyglot Project Wizard!")
	fmt.Println("========================================")
	fmt.Println()

	config := &ProjectConfig{
		Languages: []string{},
		Features:  []string{},
	}

	// Project name
	config.Name = w.prompt("Project name", "my-polyglot-app")

	// Description
	config.Description = w.prompt("Description", "A Polyglot desktop application")

	// Author
	config.Author = w.prompt("Author", "")

	// Version
	config.Version = w.prompt("Version", "0.1.0")

	// License
	fmt.Println()
	fmt.Println("Available licenses:")
	fmt.Println("  1. MIT")
	fmt.Println("  2. Apache-2.0")
	fmt.Println("  3. GPL-3.0")
	fmt.Println("  4. BSD-3-Clause")
	fmt.Println("  5. Unlicense")
	fmt.Println("  6. Custom/None")
	licenseChoice := w.prompt("Choose license (1-6)", "1")
	config.License = w.parseLicense(licenseChoice)

	// Template selection
	fmt.Println()
	fmt.Println("Project templates:")
	fmt.Println("  1. Web Application    - Full-featured web app with UI")
	fmt.Println("  2. CLI Tool          - Command-line utility")
	fmt.Println("  3. System Utility    - Background service/daemon")
	fmt.Println("  4. Desktop App       - Cross-platform desktop application")
	fmt.Println("  5. Minimal           - Bare-bones starter")
	templateChoice := w.prompt("Choose template (1-5)", "1")
	config.Template = w.parseTemplate(templateChoice)

	// Language selection
	fmt.Println()
	fmt.Println("Select languages to enable (comma-separated):")
	fmt.Println("  python, javascript, go, rust, cpp, java, ruby, php, lua, wasm, zig")
	langInput := w.prompt("Languages", "python,javascript")
	config.Languages = w.parseLanguages(langInput)

	// Features
	fmt.Println()
	fmt.Println("Additional features (comma-separated):")
	fmt.Println("  webview      - Native webview for UI")
	fmt.Println("  hmr          - Hot module reload")
	fmt.Println("  cloud        - Cloud integration")
	fmt.Println("  marketplace  - Plugin marketplace")
	fmt.Println("  security     - Enhanced security sandbox")
	fmt.Println("  signing      - Code signing")
	featuresInput := w.prompt("Features", "webview,hmr")
	config.Features = w.parseFeatures(featuresInput)

	// Git initialization
	fmt.Println()
	gitInit := w.promptBool("Initialize Git repository?", true)
	config.GitInit = gitInit

	// Package manager
	if contains(config.Languages, "javascript") {
		fmt.Println()
		fmt.Println("JavaScript package manager:")
		fmt.Println("  1. npm")
		fmt.Println("  2. yarn")
		fmt.Println("  3. pnpm")
		pmChoice := w.prompt("Choose (1-3)", "1")
		config.PackageManager = w.parsePackageManager(pmChoice)
	}

	// Python version
	if contains(config.Languages, "python") {
		fmt.Println()
		config.PythonVersion = w.prompt("Python version", "3.11")
	}

	// Webview settings
	if contains(config.Features, "webview") {
		fmt.Println()
		fmt.Println("Webview configuration:")
		config.WindowWidth = w.promptInt("Window width", 1280)
		config.WindowHeight = w.promptInt("Window height", 720)
		config.WindowResizable = w.promptBool("Resizable window?", true)
		config.DevTools = w.promptBool("Enable DevTools?", true)
	}

	return config, nil
}

func (w *ProjectWizard) prompt(question, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", question, defaultValue)
	} else {
		fmt.Printf("%s: ", question)
	}

	w.scanner.Scan()
	input := strings.TrimSpace(w.scanner.Text())

	if input == "" {
		return defaultValue
	}
	return input
}

func (w *ProjectWizard) promptBool(question string, defaultValue bool) bool {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	fmt.Printf("%s [%s]: ", question, defaultStr)
	w.scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(w.scanner.Text()))

	if input == "" {
		return defaultValue
	}

	return input == "y" || input == "yes"
}

func (w *ProjectWizard) promptInt(question string, defaultValue int) int {
	for {
		input := w.prompt(question, fmt.Sprintf("%d", defaultValue))
		var value int
		_, err := fmt.Sscanf(input, "%d", &value)
		if err == nil && value > 0 {
			return value
		}
		fmt.Println("Please enter a valid positive number")
	}
}

func (w *ProjectWizard) parseLicense(choice string) string {
	licenses := map[string]string{
		"1": "MIT",
		"2": "Apache-2.0",
		"3": "GPL-3.0",
		"4": "BSD-3-Clause",
		"5": "Unlicense",
		"6": "Custom",
	}
	if license, ok := licenses[choice]; ok {
		return license
	}
	return "MIT"
}

func (w *ProjectWizard) parseTemplate(choice string) string {
	templates := map[string]string{
		"1": "webapp",
		"2": "cli",
		"3": "system",
		"4": "desktop",
		"5": "minimal",
	}
	if template, ok := templates[choice]; ok {
		return template
	}
	return "webapp"
}

func (w *ProjectWizard) parseLanguages(input string) []string {
	if input == "" {
		return []string{"python", "javascript"}
	}

	validLangs := map[string]bool{
		"python": true, "javascript": true, "go": true, "rust": true,
		"cpp": true, "java": true, "ruby": true, "php": true,
		"lua": true, "wasm": true, "zig": true,
	}

	parts := strings.Split(input, ",")
	langs := []string{}
	for _, part := range parts {
		lang := strings.TrimSpace(strings.ToLower(part))
		if validLangs[lang] {
			langs = append(langs, lang)
		}
	}

	if len(langs) == 0 {
		return []string{"python", "javascript"}
	}

	return langs
}

func (w *ProjectWizard) parseFeatures(input string) []string {
	if input == "" {
		return []string{}
	}

	validFeatures := map[string]bool{
		"webview": true, "hmr": true, "cloud": true,
		"marketplace": true, "security": true, "signing": true,
	}

	parts := strings.Split(input, ",")
	features := []string{}
	for _, part := range parts {
		feature := strings.TrimSpace(strings.ToLower(part))
		if validFeatures[feature] {
			features = append(features, feature)
		}
	}

	return features
}

func (w *ProjectWizard) parsePackageManager(choice string) string {
	managers := map[string]string{
		"1": "npm",
		"2": "yarn",
		"3": "pnpm",
	}
	if pm, ok := managers[choice]; ok {
		return pm
	}
	return "npm"
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
