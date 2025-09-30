#!/bin/bash
# Setup script for Python runtime development dependencies

set -e

echo "🐍 Setting up Python runtime dependencies..."
echo ""

# Detect OS
OS="$(uname -s)"

case "${OS}" in
    Darwin*)
        echo "📦 Detected macOS"
        echo ""
        
        # Check if Python is already installed
        if command -v python3 &> /dev/null; then
            PYTHON_VERSION=$(python3 --version 2>&1 | cut -d' ' -f2)
            echo "Found Python ${PYTHON_VERSION}"
        fi
        
        echo "Installing Python development headers..."
        echo ""
        echo "Option 1 (Recommended): Homebrew"
        echo "  brew install python@3.11"
        echo ""
        echo "Option 2: pyenv (if you manage Python versions with pyenv)"
        echo "  pyenv install 3.11"
        echo "  pyenv global 3.11"
        echo ""
        
        if ! command -v brew &> /dev/null; then
            echo "⚠️  Homebrew not found. Install from: https://brew.sh"
            echo "   Or ensure Python dev headers are installed manually"
        else
            read -p "Install Python 3.11 via Homebrew? (y/N) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                brew install python@3.11 || brew upgrade python@3.11
                export PKG_CONFIG_PATH="$(brew --prefix python@3.11)/lib/pkgconfig:${PKG_CONFIG_PATH}"
                echo "✅ Python 3.11 installed"
            fi
        fi
        
        echo ""
        echo "To build with Python support, add to your shell profile:"
        echo "  export PKG_CONFIG_PATH=\"\$(brew --prefix python@3.11)/lib/pkgconfig:\${PKG_CONFIG_PATH}\""
        ;;
        
    Linux*)
        echo "📦 Detected Linux"
        echo ""
        
        # Check if Python is already installed
        if command -v python3 &> /dev/null; then
            PYTHON_VERSION=$(python3 --version 2>&1 | cut -d' ' -f2)
            echo "Found Python ${PYTHON_VERSION}"
        fi
        
        echo "Installing Python development headers..."
        echo ""
        
        # Detect package manager
        if command -v apt-get &> /dev/null; then
            echo "Using apt (Debian/Ubuntu)..."
            echo "  sudo apt-get install python3-dev pkg-config"
            read -p "Install now? (y/N) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                sudo apt-get update
                sudo apt-get install -y python3-dev pkg-config
            fi
            
        elif command -v dnf &> /dev/null; then
            echo "Using dnf (Fedora/RHEL)..."
            echo "  sudo dnf install python3-devel pkgconfig"
            read -p "Install now? (y/N) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                sudo dnf install -y python3-devel pkgconfig
            fi
            
        elif command -v yum &> /dev/null; then
            echo "Using yum (CentOS/RHEL)..."
            echo "  sudo yum install python3-devel pkgconfig"
            read -p "Install now? (y/N) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                sudo yum install -y python3-devel pkgconfig
            fi
            
        elif command -v pacman &> /dev/null; then
            echo "Using pacman (Arch)..."
            echo "  sudo pacman -S python pkgconf"
            read -p "Install now? (y/N) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                sudo pacman -S --noconfirm python pkgconf
            fi
            
        elif command -v apk &> /dev/null; then
            echo "Using apk (Alpine)..."
            echo "  sudo apk add python3-dev pkgconfig"
            read -p "Install now? (y/N) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                sudo apk add python3-dev pkgconfig
            fi
            
        else
            echo "❌ Could not detect package manager"
            echo "   Please install python3-dev (or python3-devel) and pkg-config manually"
            exit 1
        fi
        
        echo "✅ Python dev dependencies ready"
        ;;
        
    MINGW*|MSYS*|CYGWIN*)
        echo "📦 Detected Windows"
        echo ""
        echo "Python development setup for Windows:"
        echo ""
        echo "1. Install Python 3.11+ from https://www.python.org/downloads/"
        echo "   ✓ Check 'Install for all users'"
        echo "   ✓ Check 'Add Python to PATH'"
        echo "   ✓ Check 'Include development headers (py.h, etc.)'"
        echo ""
        echo "2. Install pkg-config support (choose one):"
        echo "   Option A: MSYS2 (recommended)"
        echo "     - Download from https://www.msys2.org/"
        echo "     - Run: pacman -S mingw-w64-x86_64-pkgconf"
        echo ""
        echo "   Option B: Use Go's CGO with direct paths"
        echo "     - Set CGO_CFLAGS and CGO_LDFLAGS manually"
        echo ""
        echo "3. Verify installation:"
        echo "   python --version"
        echo "   pkg-config --version"
        ;;
        
    *)
        echo "❌ Unsupported operating system: ${OS}"
        exit 1
        ;;
esac

# Verify installation
echo ""
echo "🔍 Verifying installation..."

if command -v python3 &> /dev/null; then
    PYTHON_VERSION=$(python3 --version 2>&1)
    echo "✓ Python: ${PYTHON_VERSION}"
else
    echo "❌ python3 not found in PATH"
    exit 1
fi

if command -v pkg-config &> /dev/null; then
    echo "✓ pkg-config: $(pkg-config --version)"
    
    if pkg-config --exists python3-embed; then
        PYTHON_EMBED_VERSION=$(pkg-config --modversion python3-embed)
        echo "✓ python3-embed: ${PYTHON_EMBED_VERSION}"
    else
        echo "⚠️  python3-embed not found by pkg-config"
        echo "   This may cause build issues. Ensure Python dev headers are installed."
    fi
else
    echo "❌ pkg-config not found"
    exit 1
fi

echo ""
echo "✅ Python runtime dependencies are ready!"
echo ""
echo "Next steps:"
echo "  1. Run tests: go test -tags=runtime_python ./tests/python_advanced_test.go"
echo "  2. Build example: cd examples/01-hello-world && go build -tags=runtime_python ./src/backend"
