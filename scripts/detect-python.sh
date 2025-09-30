#!/bin/bash
# Auto-detect if Python runtime can be enabled

# Exit with 0 if Python is available, 1 otherwise
# Outputs "runtime_python" if available, empty string otherwise

if command -v pkg-config &> /dev/null; then
    if pkg-config --exists python3-embed 2>/dev/null; then
        echo "runtime_python"
        exit 0
    elif pkg-config --exists python-3.11 2>/dev/null; then
        echo "runtime_python"
        exit 0
    elif pkg-config --exists python3 2>/dev/null; then
        echo "runtime_python"
        exit 0
    fi
fi

# Fallback: check for Python installation directly
if command -v python3-config &> /dev/null; then
    echo "runtime_python"
    exit 0
fi

# Python not available for CGO
exit 1
