#!/bin/bash
# ============================================================================
# METADATA
# ============================================================================
# build.sh - Builds CPI-SI framework statusline to native binary
# Compiles Go statusline to executable binary for optimal performance
# Non-blocking script - safe to extend and modify

# ============================================================================
# SETUP
# ============================================================================
set -e  # Exit on error

STATUSLINE_DIR="$(dirname "$0")"
cd "$STATUSLINE_DIR"

# ============================================================================
# BODY
# ============================================================================

verify_libraries() {
    echo "Verifying library packages compile..."
    echo "  → Checking all packages"
    if ! go build ./lib/...; then
        echo "❌ Library compilation failed"
        exit 1
    fi
    echo "  ✓ All libraries compile cleanly"
}

run_vet() {
    echo "Running go vet on all packages..."
    if ! go vet ./...; then
        echo "❌ go vet found issues"
        exit 1
    fi
    echo "  ✓ No vet warnings"
}

build_statusline() {
    echo "Building statusline executable..."
    echo "  → statusline/statusline"
    if ! go build -o statusline statusline.go; then
        echo "❌ Statusline build failed"
        exit 1
    fi
    echo "  ✓ Statusline built successfully"
}

# ============================================================================
# CLOSING
# ============================================================================

verify_libraries
run_vet
build_statusline

echo ""
echo "✓ Build complete: ./statusline"
