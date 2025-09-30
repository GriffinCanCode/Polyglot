#!/bin/bash
# Smart build script that auto-detects available runtimes

set -e

BUILD_TAGS=""
VERBOSE=false

# Parse arguments
TARGET="${1:-.}"
OUTPUT="${2:-polyglot}"

if [[ "$3" == "-v" ]] || [[ "$3" == "--verbose" ]]; then
    VERBOSE=true
fi

echo "🔍 Detecting available runtimes..."

# Detect Python
if bash scripts/detect-python.sh &>/dev/null; then
    PYTHON_TAG=$(bash scripts/detect-python.sh)
    if [[ -n "$PYTHON_TAG" ]]; then
        BUILD_TAGS="$PYTHON_TAG"
        if [[ "$VERBOSE" == true ]]; then
            PYTHON_VERSION=$(pkg-config --modversion python3-embed 2>/dev/null || pkg-config --modversion python3 2>/dev/null || echo "unknown")
            echo "  ✓ Python detected (version: $PYTHON_VERSION)"
        else
            echo "  ✓ Python detected"
        fi
    fi
else
    echo "  ⊘ Python dev headers not found (will use stub)"
fi

# Build command
if [[ -n "$BUILD_TAGS" ]]; then
    echo ""
    echo "🔨 Building with tags: $BUILD_TAGS"
    go build -tags="$BUILD_TAGS" -o "$OUTPUT" "$TARGET"
else
    echo ""
    echo "🔨 Building with stub runtimes (no native runtimes detected)"
    go build -o "$OUTPUT" "$TARGET"
fi

echo "✅ Built: $OUTPUT"

# Show what was enabled
if [[ -n "$BUILD_TAGS" ]]; then
    echo ""
    echo "Enabled runtimes: Python (native)"
else
    echo ""
    echo "Enabled runtimes: Python (stub)"
    echo ""
    echo "💡 Tip: To enable native Python runtime, install development headers:"
    echo "   Run: make setup-python"
    echo "   Or: bash scripts/setup-python.sh"
fi
