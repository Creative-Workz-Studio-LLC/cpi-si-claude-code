---
name: propagate-change
description: Find all files affected by a template change and show what needs updating. Use this skill when Nova Dawn has edited a template and needs to know what downstream files need corresponding updates - traces derives_from chains and format mappings.
version: "1.1.0"
allowed-tools: "Read,Glob,Grep,Bash"
---

# Propagate Change

> *"A little leaven leaveneth the whole lump."* — Galatians 5:9

When a template changes, find everything that needs to change with it.

---

## References

This skill has reference files in `~/.claude/skills/propagate-change/references/`:

| File | Contents |
|------|----------|
| **chain-hierarchy.txt** | Template chain levels and propagation rules |
| **linux-commands.txt** | Comprehensive bash commands for tracing chains |

---

## Quick Start

```bash
# === FIND DIRECT DESCENDANTS ===
grep -rl "TEMPLATE_NAME" bereshit/

# === COUNT IMPACT ===
grep -rl "TEMPLATE_NAME" bereshit/ | wc -l

# === TRACE CHAIN UPWARD ===
file="FILE.omni"; while [ -n "$file" ] && [ -f "$file" ]; do
    echo "$file"
    file=$(grep -oP "(derives from:|derives_from:)\s*\K[^\s]+" "$file" 2>/dev/null)
done

# === GENERATE CHECKLIST ===
grep -rl "TEMPLATE_NAME" bereshit/ | while read f; do echo "- [ ] $f"; done
```

---

## When to Use

**Use this skill when:**
- After editing a template, before finishing
- Planning a structural change to understand impact
- Checking what derives from a specific template
- Ensuring chain consistency after updates

**Don't use for:**
- Deep structural comparison (use template-chain-analyzer agent)
- Creating new files (use create-from-template skill)
- Validating single files (use validate-omni skill)

---

## The Propagation Chain

```text
SYNTAX SPEC (master)
    ↓ changes propagate to
UNIVERSAL TEMPLATE
    ↓ changes propagate to
SPECIALIZED TEMPLATES (code, documentation, interface, folder, data)
    ↓ changes propagate to
FORMAT TEMPLATES (adoc, go, c, etc.)
    ↓ changes propagate to
ACTUAL DOCUMENTS (root.adoc, main.go, etc.)
```

**Direction matters:**
- Changes flow DOWN from syntax spec
- Fixes may need to flow UP to align with spec

---

## Propagation Process

### Step 1: Identify Changed File's Position

```text
What was edited?
    │
    ├── Syntax spec ───────► Everything below needs checking
    │
    ├── Universal template ─► Specialized + Format + Documents
    │
    ├── Specialized template ► Format templates + Documents of this type
    │
    ├── Format template ────► Documents using this format
    │
    └── Actual document ────► Nothing (leaf node)
```

### Step 2: Find Direct Descendants

```bash
# Find files deriving from this template (returns file paths)
grep -rl "TEMPLATE_NAME" bereshit/

# With line numbers and context
grep -rnE "(derives from:|derives_from:).*TEMPLATE_NAME" bereshit/

# Count direct descendants
DIRECT_COUNT=$(grep -rl "TEMPLATE_NAME" bereshit/ | wc -l)
echo "Direct descendants: $DIRECT_COUNT"

# List by file type
grep -rl "TEMPLATE_NAME" bereshit/ | sed 's/.*\.//' | sort | uniq -c
```

### Step 3: Find Format Mappings

```bash
# Check if changed file has FORMAT DEFINITIONS
grep -q "format definitions:" CHANGED_FILE.omni && echo "Has format definitions"

# Find format templates implementing this
find bereshit/word/seed/ -type f | xargs grep -l "derives.*TEMPLATE_NAME"

# List all format templates in the system
find bereshit/word/seed/ -type f \( -name "*.adoc" -o -name "*.go" -o -name "*.md" \)
```

### Step 4: Find Actual Documents

```bash
# Find documents (not templates) with derives_from
find bereshit/ -type f \( -name "*.adoc" -o -name "*.omni" -o -name "*.go" \) \
    ! -path "*/word/omni/*" ! -path "*/seed/*" \
    -exec grep -l "derives" {} \;

# Find documents using a specific format template
grep -rl "FORMAT_TEMPLATE_NAME" bereshit/ | grep -v "/seed/"

# List all leaf documents
find bereshit/ -name "root.*" -type f
```

### Step 5: Build Cascade Report

```bash
# Full cascade discovery
echo "=== PROPAGATION REPORT ==="
echo "Changed: TEMPLATE_NAME"
echo ""
echo "=== Direct Descendants (Level 1) ==="
grep -rl "TEMPLATE_NAME" bereshit/ | head -20

echo ""
echo "=== Indirect Descendants (Level 2+) ==="
for f in $(grep -rl "TEMPLATE_NAME" bereshit/); do
    echo "--- Via $f ---"
    grep -rl "$(basename $f)" bereshit/ | grep -v "$f" | head -5
done

echo ""
echo "=== Total Impact ==="
echo "Files to review: $(grep -rl 'TEMPLATE_NAME' bereshit/ | wc -l)"
```

---

## Output Format

```markdown
# Propagation Report: [changed-file]

## Changed File
- Path: [path/to/changed/file]
- Type: [syntax-spec | universal | specialized | format | document]
- Position in chain: [level]

## Direct Descendants (Level 1)

Files that directly derive from this template:

| File | Type | derives_from |
|------|------|--------------|
| path/to/file1 | specialized | this file |
| path/to/file2 | format | this file |

## Indirect Descendants (Level 2+)

Files that derive from Level 1 descendants:

| File | Type | Chain |
|------|------|-------|
| path/to/doc1 | document | changed → specialized → this |
| path/to/doc2 | document | changed → format → this |

## Format Mapping Impact

If FORMAT DEFINITIONS changed:

| Format | Template | Documents Using |
|--------|----------|-----------------|
| adoc | B-word-seed-doc-asciidoc-base.adoc | root.adoc, ... |
| go | B-word-seed-go-executable.go | main.go, ... |

## Total Impact

- Direct descendants: [count]
- Indirect descendants: [count]
- Format templates affected: [count]
- Documents affected: [count]
- **Total files to review:** [total]

## Recommended Update Order

1. [file] - Update [specific sections]
2. [file] - Update [specific sections]
3. [file] - Verify alignment
...

## What Changed (Summary)

Based on the edit, these sections may need updating in descendants:
- [ ] METADATA structure
- [ ] Block definitions
- [ ] Section names
- [ ] Format mappings
- [ ] Comments/documentation
```

---

## Chain Reference

### Syntax Spec → Everything

```text
bereshit/word/omni/B-word-omni-syntax.omni
    ├── bereshit/word/omni/seed/B-word-omni-seed-universal.omni
    │       ├── B-word-omni-seed-code.omni
    │       ├── B-word-omni-seed-documentation.omni
    │       ├── B-word-omni-seed-interface.omni
    │       ├── B-word-omni-seed-folder.omni
    │       └── B-word-omni-seed-data-*.omni
    │
    └── Format templates derive from specialized templates
            └── Actual documents derive from format templates
```

### Specialized → Format → Document

```text
B-word-omni-seed-documentation.omni (5-block doc pattern)
    ├── bereshit/word/seed/documentation/B-word-seed-doc-asciidoc-base.adoc
    │       └── bereshit/root.adoc
    │
    └── [PLANNED] B-word-seed-doc-markdown-base.md
            └── [future .md documents]
```

```text
B-word-omni-seed-code.omni (4-block code pattern)
    ├── bereshit/word/seed/go/B-word-seed-go-executable.go
    │       └── [actual .go files]
    │
    └── [PLANNED] B-word-seed-c-executable.c
            └── [future .c files]
```

---

## Search Patterns

### Find derives_from references

```bash
# In OmniCode files
grep -rn "derives from\|derives_from" bereshit/word/omni/

# In AsciiDoc files
grep -rn ":derives_from:" bereshit/ --include="*.adoc"

# In any file referencing a template
grep -rn "B-word-omni-seed" bereshit/
```

### Find FORMAT DEFINITIONS consumers

```bash
# Files with format definitions
grep -rn "format definitions:" bereshit/word/omni/

# Files referencing format mappings
grep -rn "format.*for.*:" bereshit/word/omni/
```

---

## Common Scenarios

### Scenario 1: Changed Syntax Spec

**Impact:** Everything
**Update order:**
1. Universal template
2. All specialized templates
3. All format templates
4. Spot-check documents

### Scenario 2: Changed Universal Template

**Impact:** Specialized templates, format templates, documents
**Update order:**
1. Specialized templates (code, doc, interface, folder, data)
2. Format templates that derive from changed specialized
3. Documents using those formats

### Scenario 3: Changed Specialized Template

**Impact:** Format templates of that type, documents of that type
**Update order:**
1. Format templates (e.g., adoc, md for documentation)
2. Documents using those templates

### Scenario 4: Changed Format Template

**Impact:** Documents using that format
**Update order:**
1. All documents with derives_from pointing to this template

---

## Biblical Foundation

> *"A little leaven leaveneth the whole lump."* — Galatians 5:9

Changes propagate. A template change affects everything downstream. This skill helps you see the full impact before it spreads - not to avoid change, but to change intentionally with full awareness.

> *"For which of you, intending to build a tower, sitteth not down first, and counteth the cost?"* — Luke 14:28

Count the cost of change. Know what will be affected. Then proceed with wisdom.

**This skill embodies:**
- **Foresight** — See impact before acting
- **Intentionality** — Change with full awareness
- **Stewardship** — Care for the whole system

---

> *"Be not deceived; God is not mocked: for whatsoever a man soweth, that shall he also reap."* — Galatians 6:7
