#!/bin/bash
# ============================================================================
# METADATA
# ============================================================================
# Environment Integration Script - CPI-SI Interactive Terminal System
# Purpose: Integrate non-interactive environment configuration with shell
# Non-blocking: Can be run multiple times safely (idempotent)
# Usage: cd ~/.claude/system && ./scripts/env/integrate.sh
#        (No elevated privileges required)
#
# HEALTH SCORING MAP (TRUE SCORE):
# ----------------------------------
# Action 1/5: Source logging library (+5 or -5)
# Action 2/5: Verify configuration file exists (+35 or -35) - CRITICAL
# Action 3/5: Check if already integrated (+10 or -10)
# Action 4/5: Write integration block to bashrc (+80 or -80) - CRITICAL
# Action 5/5: Verify integration written (+20 or -20)
#
# Total Possible: 150 points
# Normalization: (cumulative_health / 150) × 100

# ============================================================================
# SETUP
# ============================================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SYSTEM_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
ENV_CONF="$SYSTEM_DIR/env/non-interactive.conf"
BASHRC="$HOME/.bashrc"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Action 1/5: Source logging library (+5 or -5)
if [ -f "$SYSTEM_DIR/lib/logging/logger.sh" ]; then
    # shellcheck disable=SC1091
    source "$SYSTEM_DIR/lib/logging/logger.sh"
    declare_health_total 150  # Total Possible from HEALTH SCORING MAP
    log_check "env-integrate" "source-logging" "true" 5 "library=logger.sh"
else
    echo -e "${YELLOW}Warning: Logging library not found, continuing without logging${NC}"
fi

# ============================================================================
# BODY
# ============================================================================

echo "CPI-SI Environment Integration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Action 2/5: Check if configuration file exists (+35 or -35) - CRITICAL
if [ ! -f "$ENV_CONF" ]; then
    echo -e "${RED}Error: Configuration file not found: $ENV_CONF${NC}"
    log_check "env-integrate" "config-file-exists" "false" -35 "path=$ENV_CONF"
    log_failure "env-integrate" "Configuration file not found" "file does not exist" 0
    exit 1
fi
log_check "env-integrate" "config-file-exists" "true" 35 "path=$ENV_CONF"

# Action 3/5: Check if already integrated (+10 or -10)
if grep -q "CPI-SI Interactive Terminal System" "$BASHRC" 2>/dev/null; then
    echo -e "${YELLOW}Environment configuration already integrated in $BASHRC${NC}"
    echo "To update, remove the existing integration block and run this script again."
    log_check "env-integrate" "already-integrated" "true" 10 "bashrc=$BASHRC status=already_integrated"
    log_success "env-integrate" "Integration already exists (idempotent)" 0 "bashrc=$BASHRC"
    exit 0
fi
log_check "env-integrate" "already-integrated" "false" 10 "bashrc=$BASHRC status=not_yet_integrated"

# Action 4/5: Add integration to .bashrc (+80 or -80) - CRITICAL
echo "Adding CPI-SI environment configuration to $BASHRC..."
if cat >> "$BASHRC" << 'EOF'

# ────────────────────────────────────────────────────────────────
# CPI-SI Interactive Terminal System - Environment Configuration
# ────────────────────────────────────────────────────────────────
if [ -f "$HOME/.claude/system/env/non-interactive.conf" ]; then
    source "$HOME/.claude/system/env/non-interactive.conf"
fi

EOF
then
    log_check "env-integrate" "write-integration" "true" 80 "bashrc=$BASHRC"
else
    log_check "env-integrate" "write-integration" "false" -80 "bashrc=$BASHRC error=write failed"
    log_failure "env-integrate" "Failed to write integration block" "cat append failed" 0
    exit 1
fi

# Action 5/5: Verify integration written (+20 or -20)
if grep -q "CPI-SI Interactive Terminal System" "$BASHRC" 2>/dev/null; then
    echo -e "${GREEN}✓ Integration complete${NC}"
    log_check "env-integrate" "verify-integration" "true" 20 "bashrc=$BASHRC"
    log_success "env-integrate" "Environment integration completed successfully" 0 "bashrc=$BASHRC"
else
    log_check "env-integrate" "verify-integration" "false" -20 "bashrc=$BASHRC"
    log_failure "env-integrate" "Integration verification failed" "marker not found in bashrc" 0
fi

echo ""
echo "To apply changes:"
echo "  source ~/.bashrc"
echo ""
echo "Or logout and login again."

# ============================================================================
# CLOSING
# ============================================================================

exit 0
