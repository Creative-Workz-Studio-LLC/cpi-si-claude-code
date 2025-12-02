---
name: validate-omni
description: Quick validation of OmniCode file structure against its parent template. Use this skill when Nova Dawn needs to verify an OmniCode file follows its template correctly - checks block presence, section completeness, and structural alignment.
version: "1.1.0"
allowed-tools: "Read,Glob,Grep,Bash"
---

# Validate OmniCode

> *"Prove all things; hold fast that which is good."* — 1 Thessalonians 5:21

Quick validation of OmniCode files against their parent templates.

---

## References

This skill has reference files in `~/.claude/skills/validate-omni/references/`:

| File | Contents |
|------|----------|
| **template-locations.txt** | All template paths and chain hierarchy |
| **validation-patterns.txt** | Regex patterns for pragma, blocks, sections |
| **linux-commands.txt** | Comprehensive bash commands for validation |

---

## Quick Start

```bash
# === PRAGMA CHECK ===
head -2 FILE.omni | grep "^#!omni"

# === BLOCK COUNT ===
grep -c "^// ═" FILE.omni

# === REQUIRED SECTIONS ===
grep -nE "(grounded in:|serves as|classifies as|authored by:|exists to:)" FILE.omni

# === DERIVES_FROM ===
grep -nE "(derives from:|derives_from:)" FILE.omni
```

---

## When to Use

**Use this skill when:**
- Quick check if an OmniCode file follows its template
- Before committing OmniCode changes
- After editing to verify structure intact
- Spot-checking files in a chain

**Don't use for:**
- Deep chain analysis (use template-chain-analyzer agent)
- Format mapping questions (use format-lookup skill)
- Creating new files (use create-from-template skill)

---

## Validation Checks

### 1. Pragma Validation

```text
✓ Pragma present on line 1 or 2
✓ Format: #!omni [type] or #!omni [type] --[format]
✓ Type is valid: template, code, interface, data, documentation, folder
```

### 2. Block Structure by Type

| Type | Expected Blocks | Check |
|------|-----------------|-------|
| **template** | METADATA, MIDDLE (varies), CLOSING | Universal pattern |
| **code** | METADATA, SETUP, BODY, CLOSING | 4-block |
| **interface** | METADATA, CONTENT, CLOSING | 3-block |
| **data** | METADATA, CONTENT/SCHEMA/VALUES, CLOSING | 3/4/5-block |
| **documentation** | PRAGMA, METADATA, HEADER, CONTEXT, CONTENT, FOOTER | 5-block |
| **folder** | METADATA, CONTENT, CLOSING | 3-block |

### 3. Required Sections (METADATA)

All types require in METADATA:
- [ ] `grounded in:` (Biblical foundation)
- [ ] `serves as` or `classifies as` (Identity)
- [ ] `authored by:` (Authorship)
- [ ] `exists to:` or purpose statement (Purpose)

### 4. derives_from Chain

```text
✓ derives_from attribute present (if not original)
✓ Referenced template exists
✓ Path is valid
```

---

## Validation Process

### Step 1: Extract and Validate Pragma

```bash
# Get pragma (must be line 1 or 2)
head -2 FILE.omni | grep "^#!omni"

# Extract type
TYPE=$(head -2 FILE.omni | grep -oP "^#!omni \K(template|code|interface|data|documentation|folder)")
echo "Type: $TYPE"

# Extract format (if present)
FORMAT=$(head -2 FILE.omni | grep -oP "#!omni [a-z]+ --\K[a-z]+")
echo "Format: ${FORMAT:-pure}"

# Validate pragma is well-formed
head -2 FILE.omni | grep -qE "^#!omni (template|code|interface|data|documentation|folder)" \
    && echo "PRAGMA: VALID" || echo "PRAGMA: INVALID"
```

### Step 2: Determine Expected Block Count

| Type | Blocks | Expected Block Names |
|------|:------:|----------------------|
| `code` | 4 | METADATA, SETUP, BODY, CLOSING |
| `documentation` | 5 | METADATA, HEADER, CONTEXT, CONTENT, FOOTER |
| `interface` | 3 | METADATA, CONTENT, CLOSING |
| `folder` | 3 | METADATA, CONTENT, CLOSING |
| `data` | 3-5 | Varies by complexity |
| `template` | varies | Depends on specialization |

### Step 3: Validate Block Structure

```bash
# Count block markers (each block has 2: start + end)
MARKERS=$(grep -c "^// ═" FILE.omni)
echo "Block markers: $MARKERS"

# List block names (excluding END markers)
grep -A1 "^// ═" FILE.omni | grep "^// [A-Z]" | grep -v "END" | sed 's/^\/\/ //'

# Verify START/END pairs match
echo "=== Block Pairs ==="
diff <(grep "^// [A-Z]" FILE.omni | grep -v "END" | sed 's/^\/\/ //') \
     <(grep "^// END" FILE.omni | sed 's/^\/\/ END //') \
     && echo "All blocks paired correctly" || echo "MISMATCH in block pairs"

# Get block names with line numbers
grep -n "^// [A-Z][A-Z ]" FILE.omni | grep -v "END"
```

### Step 4: Validate Required METADATA Sections

```bash
# All required sections in one check
echo "=== Required Sections ==="

# Biblical foundation
grep -q "grounded in:" FILE.omni \
    && echo "[x] Biblical foundation (grounded in:)" \
    || echo "[ ] Biblical foundation (grounded in:) - MISSING"

# Identity
grep -qE "(serves as|classifies as)" FILE.omni \
    && echo "[x] Identity (serves as / classifies as)" \
    || echo "[ ] Identity - MISSING"

# Authorship
grep -q "authored by:" FILE.omni \
    && echo "[x] Authorship (authored by:)" \
    || echo "[ ] Authorship - MISSING"

# Purpose
grep -qE "(exists to:|purpose:)" FILE.omni \
    && echo "[x] Purpose (exists to: / purpose:)" \
    || echo "[ ] Purpose - MISSING"

# Health (recommended)
grep -q "health:" FILE.omni \
    && echo "[x] Health scoring" \
    || echo "[ ] Health scoring (recommended)"
```

### Step 5: Validate derives_from Chain

```bash
# Find derives_from declaration
DERIVES=$(grep -oP "(derives from:|derives_from:)\s*\K[^\s]+" FILE.omni)

if [ -n "$DERIVES" ]; then
    echo "derives_from: $DERIVES"
    # Check if template exists (relative to workspace root)
    if [ -f "$DERIVES" ]; then
        echo "Template EXISTS"
    else
        echo "Template MISSING - check path"
    fi
else
    echo "No derives_from (may be original/root template)"
fi
```

### Step 6: Full Validation Summary

```bash
# One-shot validation
echo "=========================================="
echo "VALIDATING: FILE.omni"
echo "=========================================="

# Pragma
echo -n "Pragma: "
head -2 FILE.omni | grep -oP "^#!omni [a-z]+( --[a-z]+)?" || echo "MISSING"

# Blocks
echo -n "Blocks: "
echo "$(grep -c '^// ═' FILE.omni) markers"

# Required sections count
FOUND=$(grep -cE "(grounded in:|serves as|classifies as|authored by:|exists to:)" FILE.omni)
echo "Required sections: $FOUND/4 found"

# Status
if head -2 FILE.omni | grep -qE "^#!omni" && [ "$FOUND" -ge 4 ]; then
    echo "STATUS: VALID"
else
    echo "STATUS: NEEDS ATTENTION"
fi
```

---

## Output Format

```markdown
# OmniCode Validation: [filename]

## Status: [VALID | WARNINGS | INVALID]

## Pragma
- Type: [type]
- Format: [format or "pure"]
- Line: [line number]

## Block Structure
- Expected: [X]-block ([type])
- Found: [list blocks found]
- Status: [Complete | Missing: list]

## Required Sections
- [x] Biblical Foundation (grounded in)
- [x] Identity (serves as / classifies as)
- [ ] Authorship (authored by) ← MISSING
- [x] Purpose (exists to)

## Template Chain
- derives_from: [path or "original"]
- Template exists: [Yes | No | N/A]

## Issues Found
1. [Issue description]
2. [Issue description]

## Recommendations
- [Action to fix]
```

---

## Quick Reference

### Valid Pragma Types

```text
#!omni template          # Universal/seed templates
#!omni code              # 4-block executable/library
#!omni interface         # 3-block contracts/APIs
#!omni data              # 3/4/5-block data files
#!omni documentation     # 5-block documents
#!omni folder            # 3-block folder metadata
```

### Block Markers Pattern

```text
// ═══════════════════════════════════════════════════════════════════════════
// BLOCK NAME
// ═══════════════════════════════════════════════════════════════════════════

[block content]

// ═══════════════════════════════════════════════════════════════════════════
// END BLOCK NAME
// ═══════════════════════════════════════════════════════════════════════════
```

### Template Locations

```text
bereshit/word/omni/seed/
├── B-word-omni-seed-universal.omni      # Universal template
├── B-word-omni-seed-code.omni           # 4-block code
├── B-word-omni-seed-documentation.omni  # 5-block docs
├── B-word-omni-seed-interface.omni      # 3-block interface
├── B-word-omni-seed-folder.omni         # 3-block folder
└── B-word-omni-seed-data-*.omni         # Data variants
```

---

## Biblical Foundation

> *"Prove all things; hold fast that which is good."* — 1 Thessalonians 5:21

Validation is not skepticism - it's stewardship. We verify structure not from doubt but from care. Just as the Bereans searched the scriptures to verify truth, we validate our files to ensure they align with the revealed pattern.

**This skill embodies:**
- **Diligence** — Check before assuming
- **Stewardship** — Care for the codebase
- **Faithfulness** — Maintain alignment with the pattern

---

> *"Examine yourselves, whether ye be in the faith; prove your own selves."* — 2 Corinthians 13:5
