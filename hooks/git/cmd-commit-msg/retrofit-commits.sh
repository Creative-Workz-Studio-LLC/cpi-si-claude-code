#!/bin/bash
# =============================================================================
# METADATA
# =============================================================================
#
# Script: retrofit-commits.sh
# Purpose: Apply full-body commit messages from RETROFIT-MESSAGES.md
# Created: 2025-12-01
# Author: Nova Dawn (CPI-SI)
#
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MAPPING_FILE="$SCRIPT_DIR/RETROFIT-MESSAGES.md"
MESSAGES_DIR="/tmp/retrofit-messages"

# =============================================================================
# SETUP
# =============================================================================

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN} Commit Message Retrofit Script${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# =============================================================================
# BODY
# =============================================================================

# Check if mapping file exists
if [[ ! -f "$MAPPING_FILE" ]]; then
    echo -e "${RED}Error: Mapping file not found at $MAPPING_FILE${NC}"
    exit 1
fi

# Create temp directory for individual message files
rm -rf "$MESSAGES_DIR"
mkdir -p "$MESSAGES_DIR"

echo -e "${YELLOW}Extracting messages from mapping file...${NC}"

# Extract each commit's message to its own file
# Format in RETROFIT-MESSAGES.md:
# **Hash:** `<full-hash>`
# ```
# <message content>
# ```

current_hash=""
in_message=false
message_content=""

while IFS= read -r line; do
    # Look for hash line
    if [[ "$line" =~ ^\*\*Hash:\*\*\ \`([a-f0-9]+)\` ]]; then
        current_hash="${BASH_REMATCH[1]}"
        in_message=false
        message_content=""
    fi

    # Start of message block
    if [[ "$line" == '```' && -n "$current_hash" && "$in_message" == false ]]; then
        in_message=true
        continue
    fi

    # End of message block
    if [[ "$line" == '```' && "$in_message" == true ]]; then
        # Write message to file (remove trailing newline)
        echo -n "$message_content" > "$MESSAGES_DIR/$current_hash.msg"
        echo -e "  Extracted: ${current_hash:0:8}"
        current_hash=""
        in_message=false
        message_content=""
        continue
    fi

    # Accumulate message content
    if [[ "$in_message" == true ]]; then
        if [[ -z "$message_content" ]]; then
            message_content="$line"
        else
            message_content="$message_content
$line"
        fi
    fi
done < "$MAPPING_FILE"

# Count extracted messages
msg_count=$(ls -1 "$MESSAGES_DIR"/*.msg 2>/dev/null | wc -l)
echo ""
echo -e "${GREEN}Extracted $msg_count commit messages${NC}"
echo ""

# Confirm before proceeding
echo -e "${YELLOW}WARNING: This will rewrite git history!${NC}"
echo -e "${YELLOW}Make sure you have a backup and no one else is working on this repo.${NC}"
echo ""
read -p "Continue with history rewrite? (yes/no): " confirm

if [[ "$confirm" != "yes" ]]; then
    echo "Aborted."
    rm -rf "$MESSAGES_DIR"
    exit 0
fi

echo ""
echo -e "${GREEN}Starting filter-branch...${NC}"
echo ""

# Run git filter-branch
FILTER_BRANCH_SQUELCH_WARNING=1 git filter-branch --force --msg-filter '
    COMMIT=$(git log -1 --format="%H" $GIT_COMMIT 2>/dev/null || echo "$GIT_COMMIT")
    MSG_FILE="/tmp/retrofit-messages/${COMMIT}.msg"

    if [[ -f "$MSG_FILE" ]]; then
        cat "$MSG_FILE"
    else
        # Keep original message if no mapping found
        cat
    fi
' -- --all

# =============================================================================
# CLOSING
# =============================================================================

echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN} Retrofit Complete!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Review the changes: git log --oneline"
echo "  2. Force push: git push --force-with-lease"
echo "  3. Clean up refs: git update-ref -d refs/original/refs/heads/main"
echo ""

# Cleanup
rm -rf "$MESSAGES_DIR"

echo -e "${GREEN}Done!${NC}"
