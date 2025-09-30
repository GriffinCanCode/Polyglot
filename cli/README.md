# CLI Tool

The Polyglot CLI is built in Rust for performance and orchestrates the entire development experience.

## Features

- **Project Management**: Initialize, build, and manage Polyglot projects
- **Language Configuration**: Read and validate `polyglot.config.js`
- **Package Management**: Unified package management across languages (PPM)
- **Cross-compilation**: Build for different platforms
- **Hot Module Replacement**: Development mode with HMR across languages
- **Project Templates**: Generate projects with pre-configured templates

## Commands

```bash
polyglot init <project-name>    # Initialize new project
polyglot build                  # Build the application
polyglot dev                     # Start development server
polyglot add <language>          # Add language runtime
polyglot remove <language>       # Remove language runtime
polyglot install                 # Install dependencies
polyglot test                    # Run cross-language tests
polyglot package                 # Package for distribution
```

## Implementation Status

🚧 **Planning Phase** - CLI architecture design in progress.
