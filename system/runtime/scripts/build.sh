#!/bin/bash
# ============================================================================
# METADATA
# ============================================================================
# Build Script - CPI-SI Interactive Terminal System
# Purpose: Compile all system binaries from source
# Non-blocking: Build process with clear error reporting
# Usage: ./scripts/build.sh (from system directory)
#
# HEALTH SCORING MAP (TRUE SCORE):
# ----------------------------------
# Setup actions:
#   Action 1/3: Source logging library (+5 or -5)
#   Action 2/3: Create bin/ directory (+3 or -3)
#   Action 3/3: Initialize operation (+2)
#
# Per binary (7 binaries, 148 points each = 1036 total):
#   Core Commands (4): validate, test, status, diagnose
#   Utility Commands (2): debugger, unix-safe
#   Demo Programs (1): rails-demo
#
#   Action 1/6: Check source exists (+15 or -15)
#   Action 2/6: Go compiler available (+20 or -20)
#   Action 3/6: Compilation starts (+5)
#   Action 4/6: Compilation completes (+90 or -90) - CRITICAL
#   Action 5/6: Binary written (+10 or -10)
#   Action 6/6: Binary executable (+8 or -8)
#
# Verification (3 actions = 10 points):
#   Action 1/3: Count results (+3)
#   Action 2/3: Log final state (+5 or -5)
#   Action 3/3: Display output (+2 or -2)
#
# Total Possible: 1056 points
# Normalization: (cumulative_health / 1056) × 100

# ============================================================================
# SETUP
# ============================================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SYSTEM_DIR="$(dirname "$SCRIPT_DIR")"
WORKSPACE_ROOT="$(cd "$SYSTEM_DIR/../../.." && pwd)"  # Navigate to .claude/ (workspace root)

# Change to workspace root for go.work-aware builds
cd "$WORKSPACE_ROOT"

# Store relative paths from workspace root
SYSTEM_RELATIVE="cpi-si/system/runtime"
BIN_DIR="$SYSTEM_DIR/bin"

# Setup Action 1/3: Source logging library (+5 or -5)
# shellcheck disable=SC1091
if source "$SYSTEM_DIR/lib/logging/logger.sh" 2>/dev/null; then
    declare_health_total 1056  # Total Possible from HEALTH SCORING MAP
    log_check "build" "source-logging" "true" 5 "library: logger.sh"
else
    echo "ERROR: Failed to source logging library"
    exit 1
fi

# Colors
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Setup Action 2/3: Create bin/ directory (+3 or -3)
if mkdir -p "$BIN_DIR" 2>/dev/null; then
    log_check "build" "create-bin-directory" "true" 3 "directory: $BIN_DIR"
else
    log_check "build" "create-bin-directory" "false" -3 "directory: $BIN_DIR"
    exit 1
fi

# Setup Action 3/3: Initialize operation (+2)
log_operation "build" "build-all-binaries" 2 "validate test status diagnose debugger unix-safe rails-demo"
log_snapshot "build" "before-build" 0

# ============================================================================
# BODY
# ============================================================================

echo -e "${YELLOW}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║     CPI-SI Interactive Terminal System - Build Script        ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Build each command
# Core commands in cmd/, demo programs in demo/
COMMANDS=("validate" "test" "status" "diagnose" "debugger" "unix-safe" "rails-demo")
SUCCESS=0
FAILED=0
FAILED_COMMANDS=()

for cmd in "${COMMANDS[@]}"; do
    echo -n "Building $cmd... "

    # Determine source directory (dev/demo/ for rails-demo, cmd/ for everything else)
    # Paths are workspace-relative from .claude/ root
    if [ "$cmd" = "rails-demo" ]; then
        SRC_DIR="cpi-si/system/dev/demo"
        SRC_FILE="rails_correlation_demo.go"
    else
        SRC_DIR="$SYSTEM_RELATIVE/cmd/$cmd"
        SRC_FILE=""  # Will use directory path
    fi

    # Per-binary Action 1/6: Check source exists (+15 or -15)
    if [ "$cmd" = "rails-demo" ]; then
        if [ -f "$SRC_DIR/$SRC_FILE" ]; then
            log_check "build" "source-exists-$cmd" "true" 15 "source: $SRC_DIR/$SRC_FILE"
        else
            log_check "build" "source-exists-$cmd" "false" -15 "source: $SRC_DIR/$SRC_FILE"
            echo -e "${RED}✗ Source not found${NC}"
            FAILED=$((FAILED + 1))
            FAILED_COMMANDS+=("$cmd")
            continue
        fi
    else
        if [ -d "$SRC_DIR" ] && compgen -G "$SRC_DIR/*.go" > /dev/null; then
            log_check "build" "source-exists-$cmd" "true" 15 "source: $SRC_DIR"
        else
            log_check "build" "source-exists-$cmd" "false" -15 "source: $SRC_DIR"
            echo -e "${RED}✗ Source not found${NC}"
            FAILED=$((FAILED + 1))
            FAILED_COMMANDS+=("$cmd")
            continue
        fi
    fi

    # Per-binary Action 2/6: Go compiler available (+20 or -20)
    if command -v go &> /dev/null; then
        log_check "build" "go-compiler-$cmd" "true" 20 "compiler: $(go version | awk '{print $3}')"
    else
        log_check "build" "go-compiler-$cmd" "false" -20 "compiler: not found"
        echo -e "${RED}✗ Go compiler not found${NC}"
        FAILED=$((FAILED + 1))
        FAILED_COMMANDS+=("$cmd")
        continue
    fi

    # Per-binary Action 3/6: Compilation starts (+5)
    log_check "build" "compilation-start-$cmd" "true" 5 "command: go build"

    # Per-binary Action 4/6: Compilation completes (+90 or -90) - CRITICAL
    if [ "$cmd" = "rails-demo" ]; then
        BUILD_SRC="./$SRC_DIR/$SRC_FILE"
    else
        BUILD_SRC="./$SRC_DIR"
    fi

    # Build from workspace root with go.work - outputs to system/runtime/bin/
    if go build -o "$BIN_DIR/$cmd" "$BUILD_SRC" 2>/dev/null; then
        log_check "build" "compilation-complete-$cmd" "true" 90 "binary: $BIN_DIR/$cmd"
    else
        log_check "build" "compilation-complete-$cmd" "false" -90 "binary: $BIN_DIR/$cmd"
        echo -e "${RED}✗ Compilation failed${NC}"
        FAILED=$((FAILED + 1))
        FAILED_COMMANDS+=("$cmd")
        continue
    fi

    # Per-binary Action 5/6: Binary written (+10 or -10)
    if [ -f "$BIN_DIR/$cmd" ]; then
        log_check "build" "binary-written-$cmd" "true" 10 "file: $BIN_DIR/$cmd"
    else
        log_check "build" "binary-written-$cmd" "false" -10 "file: $BIN_DIR/$cmd"
        echo -e "${RED}✗ Binary not written${NC}"
        FAILED=$((FAILED + 1))
        FAILED_COMMANDS+=("$cmd")
        continue
    fi

    # Per-binary Action 6/6: Binary executable (+8 or -8)
    if [ -x "$BIN_DIR/$cmd" ]; then
        log_check "build" "binary-executable-$cmd" "true" 8 "permissions: executable"
        echo -e "${GREEN}✓ Success${NC}"
        SUCCESS=$((SUCCESS + 1))
    else
        log_check "build" "binary-executable-$cmd" "false" -8 "permissions: not executable"
        echo -e "${RED}✗ Not executable${NC}"
        FAILED=$((FAILED + 1))
        FAILED_COMMANDS+=("$cmd")
    fi
done

echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "  Build Results: ${GREEN}$SUCCESS passed${NC}, ${RED}$FAILED failed${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Verification Action 1/3: Count results (+3)
log_check "build" "count-results" "true" 3 "success: $SUCCESS, failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ All binaries built successfully${NC}"
    echo ""
    echo "Core Commands:"
    echo "  ./bin/status      - Quick health check"
    echo "  ./bin/validate    - Validate installation"
    echo "  ./bin/test        - Test operations"
    echo "  ./bin/diagnose    - Detailed diagnostics"
    echo ""
    echo "Utility Commands:"
    echo "  ./bin/debugger    - Log analysis and health assessment"
    echo "  ./bin/unix-safe   - Unix line ending converter (CRLF → LF)"
    echo ""
    echo "Demo Programs:"
    echo "  ./bin/rails-demo  - Rails correlation demonstration"
    echo ""
    echo "Try: ./bin/status"

    # Verification Action 2/3: Log final state (+5)
    log_success "build" "Build completed successfully" 5 \
        "total: ${#COMMANDS[@]}" \
        "passed: $SUCCESS" \
        "failed: $FAILED"

    # Verification Action 3/3: Display output (+2)
    log_check "build" "display-output" "true" 2 "output: success message shown"

    exit 0
else
    echo ""
    echo -e "${RED}✗ Build failed for $FAILED binaries${NC}"

    # Verification Action 2/3: Log final state (-5)
    log_failure "build" "Build failed" "compilation errors" -5 \
        "total: ${#COMMANDS[@]}" \
        "passed: $SUCCESS" \
        "failed: $FAILED" \
        "failed_commands: ${FAILED_COMMANDS[*]}"

    # Verification Action 3/3: Display output (-2)
    log_check "build" "display-output" "false" -2 "output: failure message shown"

    exit 1
fi

# ============================================================================
# CLOSING
# ============================================================================
