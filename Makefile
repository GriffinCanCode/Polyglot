.PHONY: help setup-python test test-python test-python-version test-webview test-all build build-force-python build-webview-demo build-python-demo run-python-demo clean install verify-python

# Default target
help:
	@echo "Polyglot Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build           - Build CLI tool (auto-detects Python, includes webview)"
	@echo "  test            - Run tests (auto-detects Python)"
	@echo "  setup-python    - Install Python development dependencies"
	@echo "  verify-python   - Check Python runtime setup"
	@echo "  install         - Install CLI tool to GOPATH/bin"
	@echo "  clean           - Remove build artifacts"
	@echo ""
	@echo "Webview targets:"
	@echo "  build-webview-demo - Build webview demo application"
	@echo "  build-python-demo  - Build Python + JS + Webview demo"
	@echo "  run-python-demo    - Build and run Python webview demo"
	@echo "  test-webview       - Run webview tests"
	@echo ""
	@echo "Advanced targets:"
	@echo "  test-python        - Run Python runtime tests"
	@echo "  test-python-version - Run Python version compatibility tests"
	@echo "  build-force-python - Force build with Python (fails if unavailable)"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. Run: make build    (automatically uses Python if available)"
	@echo "  2. Or: make setup-python && make build  (to ensure Python is available)"
	@echo "  3. Try: make run-python-demo  (Python + JS + Webview demo)"
	@echo "  4. Or: make build-webview-demo && ./examples/02-webview-demo/webview-demo"
	@echo ""

# Setup Python development dependencies
setup-python:
	@bash scripts/setup-python.sh

# Run standard tests
test:
	@echo "🧪 Running core tests..."
	@go test -v ./tests/core_test.go
	@go test -v ./tests/runtime_test.go
	@echo ""
	@echo "🐍 Testing Python runtime..."
	@if bash scripts/detect-python.sh &>/dev/null; then \
		echo "Python detected, running native tests..."; \
		go test -v -tags=runtime_python ./tests/python_advanced_test.go || true; \
	else \
		echo "Python dev headers not found, skipping native tests"; \
		echo "Tip: Run 'make setup-python' to enable native Python runtime"; \
	fi
	@echo "✅ Tests complete"

# Run Python runtime tests only
test-python:
	@if ! bash scripts/detect-python.sh &>/dev/null; then \
		echo "❌ Python development headers not found"; \
		echo "   Run: make setup-python"; \
		exit 1; \
	fi
	@echo "Running Python runtime tests..."
	go test -v -tags=runtime_python ./tests/python_advanced_test.go
	@echo "✅ Python runtime tests passed"

# Run Python version compatibility tests
test-python-version:
	@if ! bash scripts/detect-python.sh &>/dev/null; then \
		echo "❌ Python development headers not found"; \
		echo "   Run: make setup-python"; \
		exit 1; \
	fi
	@echo "Running Python version compatibility tests..."
	go test -v -tags=runtime_python ./tests/python_version_test.go
	@echo "✅ Python version tests passed"

# Build CLI tool with auto-detection
build:
	@bash scripts/build.sh ./cli polyglot

# Build with Python runtime support (fails if unavailable)
build-force-python:
	@echo "Building with Python runtime support..."
	@if ! bash scripts/detect-python.sh &>/dev/null; then \
		echo "❌ Python development headers not found"; \
		echo "   Run: make setup-python"; \
		exit 1; \
	fi
	go build -tags=runtime_python -o polyglot ./cli
	@echo "✅ Built: ./polyglot (with Python runtime)"

# Build hello-world example (auto-detects Python)
example:
	@echo "Building hello-world example..."
	@bash scripts/build.sh ./examples/01-hello-world/src/backend examples/01-hello-world/dist/hello-world
	@echo "Run: cd examples/01-hello-world && ./dist/hello-world"

# Build webview demo application
build-webview-demo:
	@echo "Building webview demo..."
	@cd examples/02-webview-demo && go build -o webview-demo
	@echo "✅ Built: examples/02-webview-demo/webview-demo"
	@echo "Run: cd examples/02-webview-demo && ./webview-demo"

# Build Python + JS + Webview demo
build-python-demo:
	@echo "Building Python + JS + Webview demo..."
	@if ! bash scripts/detect-python.sh &>/dev/null; then \
		echo "❌ Python development headers not found"; \
		echo "   Run: make setup-python"; \
		exit 1; \
	fi
	@cd examples/03-python-webview-demo && go build -tags=runtime_python -o python-demo
	@echo "✅ Built: examples/03-python-webview-demo/python-demo"
	@echo "Run: cd examples/03-python-webview-demo && ./python-demo"

# Build and run Python demo
run-python-demo: build-python-demo
	@echo "Starting Python + JS + Webview demo..."
	@cd examples/03-python-webview-demo && ./python-demo

# Test webview functionality
test-webview:
	@echo "🧪 Running webview tests..."
	@go test -v ./tests/webview_test.go
	@echo "✅ Webview tests complete"

# Install CLI tool (with auto-detection)
install:
	@echo "Installing polyglot CLI..."
	@if bash scripts/detect-python.sh &>/dev/null; then \
		echo "Installing with Python runtime support..."; \
		go install -tags=runtime_python ./cli; \
	else \
		echo "Installing with stub runtimes..."; \
		echo "Tip: Run 'make setup-python' for native Python support"; \
		go install ./cli; \
	fi
	@echo "✅ Installed to GOPATH/bin"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f polyglot
	rm -f examples/01-hello-world/dist/hello-world
	rm -f examples/01-hello-world/dist/hello-world-python
	rm -f examples/02-webview-demo/webview-demo
	rm -f examples/03-python-webview-demo/python-demo
	@echo "✅ Cleaned"

# Verify Python runtime setup
verify-python:
	@echo "🔍 Verifying Python runtime setup..."
	@echo ""
	@echo "Python:"
	@if command -v python3 &>/dev/null; then \
		python3 --version; \
	else \
		echo "❌ python3 not found"; \
	fi
	@echo ""
	@echo "pkg-config:"
	@if command -v pkg-config &>/dev/null; then \
		pkg-config --version; \
	else \
		echo "❌ pkg-config not found"; \
	fi
	@echo ""
	@echo "python3-embed:"
	@if pkg-config --exists python3-embed 2>/dev/null; then \
		echo "✅ python3-embed found"; \
		echo "Version: $$(pkg-config --modversion python3-embed)"; \
		echo "CFLAGS: $$(pkg-config --cflags python3-embed)"; \
		echo "LDFLAGS: $$(pkg-config --libs python3-embed)"; \
		echo ""; \
		echo "✅ Python runtime is ready to use!"; \
		echo "   Just run: make build"; \
	else \
		echo "❌ python3-embed not found"; \
		echo ""; \
		echo "To enable Python runtime:"; \
		echo "  1. Run: make setup-python"; \
		echo "  2. Follow the prompts"; \
		echo "  3. Run: make build"; \
	fi
