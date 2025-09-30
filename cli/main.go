package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "init":
		handleInit(args)
	case "build":
		handleBuild(args)
	case "dev":
		handleDev(args)
	case "test":
		handleTest(args)
	case "version":
		handleVersion(args)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Polyglot CLI - Modern Desktop Application Framework")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  polyglot <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init     Initialize a new Polyglot project")
	fmt.Println("  build    Build the application")
	fmt.Println("  dev      Start development mode")
	fmt.Println("  test     Run tests")
	fmt.Println("  version  Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  polyglot init myapp")
	fmt.Println("  polyglot build --platform darwin --arch arm64")
	fmt.Println("  polyglot dev --port 3000")
	fmt.Println()
}
