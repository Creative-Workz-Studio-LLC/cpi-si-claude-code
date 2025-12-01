#!/bin/bash
# ============================================================================
# METADATA
# ============================================================================
# build.sh - Builds CPI-SI framework hooks to native binaries
# Compiles Go hook files into executable binaries for optimal performance
# Non-blocking script - safe to extend and modify

# ============================================================================
# SETUP
# ============================================================================
set -e  # Exit on error

HOOKS_DIR="$(dirname "$0")"
cd "$HOOKS_DIR"

# ============================================================================
# BODY
# ============================================================================

build_session_hooks() {
    echo "Building session hooks..."
    echo "  → session/start"
    go build -o session/start ./session/cmd-start

    echo "  → session/end"
    go build -o session/end ./session/cmd-end

    echo "  → session/stop"
    go build -o session/stop ./session/cmd-stop

    echo "  → session/subagent-stop"
    go build -o session/subagent-stop ./session/cmd-subagent-stop

    echo "  → session/notification"
    go build -o session/notification ./session/cmd-notification

    echo "  → session/pre-compact"
    go build -o session/pre-compact ./session/cmd-pre-compact
}

build_tool_hooks() {
    echo "Building tool hooks..."
    echo "  → tool/post-use"
    go build -o tool/post-use ./tool/cmd-post-use

    echo "  → tool/pre-use"
    go build -o tool/pre-use ./tool/cmd-pre-use
}

build_prompt_hooks() {
    echo "Building prompt hooks..."
    echo "  → prompt/submit"
    go build -o prompt/submit ./prompt/cmd-submit
}

build_libraries() {
    # Libraries are compiled automatically when imported by hooks
    # Go doesn't build standalone library binaries like C
    # This function intentionally left as no-op for future expansion
    :
}

# ============================================================================
# CLOSING
# ============================================================================

echo "Building CPI-SI hooks..."
build_libraries
build_session_hooks
build_tool_hooks
build_prompt_hooks
echo "✓ All hooks and libraries built successfully"
