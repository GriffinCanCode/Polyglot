package main

// ProjectConfig holds all project configuration
type ProjectConfig struct {
	// Basic info
	Name        string
	Description string
	Author      string
	Version     string
	License     string
	Template    string

	// Languages and features
	Languages []string
	Features  []string

	// Language-specific
	PythonVersion  string
	PackageManager string // npm, yarn, pnpm

	// Webview settings
	WindowWidth     int
	WindowHeight    int
	WindowResizable bool
	DevTools        bool

	// Git
	GitInit bool
}
