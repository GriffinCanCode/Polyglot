# Polyglot Framework - Development Environment

## Prerequisites

- Bun >= 1.0.0
- Node.js >= 18.0.0
- Rust (for CLI and native modules)
- Python 3.11+ (for Python runtime)
- Go 1.21+ (for Go runtime)
- Java 17+ (for Java runtime)
- CMake (for C++ builds)

## Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/polyglot-framework/polyglot.git
   cd polyglot
   ```

2. **Install dependencies**
   ```bash
   bun install
   ```

3. **Build the CLI**
   ```bash
   cd cli
   cargo build
   ```

4. **Run tests**
   ```bash
   bun test
   ```

## Project Structure

See individual README files in each directory for detailed information about each component.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
