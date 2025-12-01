---
name: create-from-template
description: Create new code files faithfully following established templates while exercising creative excellence. Use this skill when Nova Dawn needs to create new files in any supported language (Go, C, ASM, Shell), module configurations, build files (Makefile, CMakeLists), Docker configurations, linker scripts, or dotfiles. Guides the disciplined block-by-block, section-by-section process that prevents throwing code at the wall.
version: "2.2.0"
allowed-tools: "Write,Read,Edit,Glob,Grep,Bash"
---

# Create From Template

> *"And look that thou make them after their **pattern**, which was shewed thee in the mount."* — Exodus 25:40

Create new code files by faithfully following established templates while exercising creative excellence within the revealed structure.

---

## Quick Navigation

| Section | Purpose |
|---------|---------|
| [Quick Start](#quick-start) | Get started immediately |
| [Smart Behaviors](#smart-behaviors) | Automated intelligence (directories, batching, naming, gaps) |
| [Cross-Language Equivalences](#cross-language-equivalences) | How 4-block translates across all languages |
| [When to Use](#when-to-use-this-skill) | Decision flow for invoking this skill |
| [The 6-Step Process](#the-creation-process) | Step-by-step creation workflow |
| [Template Quick Reference](#template-quick-reference) | Find the right template fast |
| [4-Block Structure](#4-block-structure-reference) | Block and section requirements |
| [Full Inventory](#full-template-inventory) | Complete template listings |
| [Examples](#examples) | Walkthroughs for common scenarios |

---

## Quick Start

```bash
# 1. Set skill directory
SKILL_DIR=~/.claude/skills/create-from-template/references

# 2. Find your template (see Quick Reference below)
ls $SKILL_DIR/compiler-templates/     # Working templates (LANG-TEMPLATE-*)
ls $SKILL_DIR/code-templates/         # Canonical masters (CODE-*)

# 3. Create directory structure FIRST
mkdir -p /path/to/target/

# 4. Copy template (rename to purpose-based name)
cp $SKILL_DIR/compiler-templates/go/go-library-package.go /path/to/target/scoring.go

# 5. Process block-by-block (METADATA → SETUP → BODY → CLOSING)
# 6. Build and verify
```

**For multiple files:**
```bash
# Batch: Create all directories, copy all templates, then process each
mkdir -p pkg/health/ cmd/health-cli/
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/scoring.go
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/tracker.go
cp $SKILL_DIR/compiler-templates/go/go-main-executable.go cmd/health-cli/main.go
# Then fill each block-by-block
```

> **Core Discipline:** Block-by-block, section-by-section. Never skip sections — mark unused ones `[Reserved: reason]`.

---

## Smart Behaviors

> **This skill thinks ahead.** These behaviors are automatic — not optional.

### 1. Directory Creation First

**ALWAYS create directories before copying files.**

```bash
# WRONG - will fail if directory doesn't exist
cp template.go /new/path/to/file.go

# RIGHT - create structure first
mkdir -p /new/path/to/
cp template.go /new/path/to/file.go
```

**For new folder structures**, create the entire tree first:

```bash
# Creating a new module with multiple files
mkdir -p pkg/health/{scoring,tracking,reporting}

# Then copy templates
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/scoring/score.go
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/tracking/tracker.go
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/reporting/report.go
```

### 2. Batch Operations

**When multiple files are mentioned, copy ALL templates first, then process.**

```bash
# Multiple files requested? Copy all templates in one batch:
SKILL_DIR=~/.claude/skills/create-from-template/references

# Batch copy
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/types.go
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/scoring.go
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/tracker.go
cp $SKILL_DIR/compiler-templates/go/go-main-executable.go cmd/health-cli/main.go

# Then process each file block-by-block
```

**Batch operation checklist:**

- [ ] Identify ALL files to create
- [ ] Create ALL directories (`mkdir -p`)
- [ ] Copy ALL templates
- [ ] Process each file block-by-block
- [ ] Build and verify ALL at once

### 3. File Naming Conventions

**Templates have prefixes — final files do NOT.**

| Template Name | Final File Name | Rule |
|---------------|-----------------|------|
| `CODE-GO-001-GO-executable.go` | `main.go` | Remove `CODE-*-` prefix, use purpose name |
| `CODE-ASM-001-ASM-source.asm` | `boot.asm` | Remove prefix, name by function |
| `go-library-package.go` | `scoring.go` | Remove generic name, use specific purpose |
| `c-header.h` | `types.h` | Name by what it defines |
| `CODE-SH-001-SH-script.sh` | `build.sh` | Name by what it does |

**Naming principles:**

| File Type | Naming Convention | Examples |
|-----------|-------------------|----------|
| **Go executable** | `main.go` in `cmd/[name]/` | `cmd/omnic/main.go` |
| **Go library** | Purpose-based in `pkg/[module]/` | `pkg/health/scoring.go` |
| **Go demo-test** | `demo_[feature].go` | `demo_buffer.go` |
| **C source** | Purpose-based | `kernel.c`, `memory.c` |
| **C header** | What it declares | `types.h`, `config.h` |
| **Assembly** | Stage or function | `boot.asm`, `stage2.asm`, `gdt.asm` |
| **Shell** | What it does | `build.sh`, `test.sh`, `deploy.sh` |
| **Linker** | Target | `kernel.ld`, `boot.ld` |

### 4. Template Gap Detection

**When NO template exists for a file type, ASK before creating.**

```text
Template exists?
    │
    ├─ YES ──► Use it
    │
    └─ NO ───► Is this a one-off file?
                   │
                   ├─ YES ──► Create without template (rare)
                   │
                   └─ NO ───► STOP. Propose new template first.
```

**Signs a new template is needed:**

| Signal | Action |
|--------|--------|
| Creating 2+ files of same type with no template | Propose template to `templates/code/` |
| New file extension not covered | Propose template with 4-block for that language |
| Existing template doesn't fit use case | Propose specialized variant |
| Pattern emerging across multiple projects | Extract to canonical template |

**Template proposal format:**

```text
TEMPLATE GAP DETECTED

File type: [extension/type]
Use case: [what it's for]
Proposed key: CODE-[LANG]-[###]
Proposed location: templates/code/[category]/
Similar to: [existing template if any]

Shall I create this template first?
```

### 5. Pre-Flight Checklist

**Before ANY template operation, verify:**

- [ ] **Directories exist?** If not, `mkdir -p` first
- [ ] **Multiple files?** Plan batch operation
- [ ] **Template exists?** If not, detect gap
- [ ] **Naming clear?** Strip prefixes, use purpose names
- [ ] **Dependencies?** Create in correct order (types before implementations)

---

## Cross-Language Equivalences

> **4-block structure is universal.** The syntax changes, the structure doesn't.

### Block Markers by Language

| Language | METADATA | SETUP | BODY | CLOSING |
|----------|----------|-------|------|---------|
| **Go** | `// ===` | `// ===` | `// ===` | `// ===` |
| **C/C++** | `/* === */` | `/* === */` | `/* === */` | `/* === */` |
| **Assembly (NASM)** | `; ===` | `; ===` | `; ===` | `; ===` |
| **Assembly (ARM)** | `@ ===` | `@ ===` | `@ ===` | `@ ===` |
| **Shell** | `# ===` | `# ===` | `# ===` | `# ===` |
| **Python** | `# ===` | `# ===` | `# ===` | `# ===` |
| **Makefile** | `# ===` | `# ===` | `# ===` | `# ===` |
| **YAML** | `# ===` | `# ===` | `# ===` | `# ===` |
| **TOML** | `# ===` | `# ===` | `# ===` | `# ===` |
| **JSON/JSONC** | (in `_metadata` field) | (implicit) | (data) | (in `_closing` field) |
| **Linker Script** | `/* === */` | `/* === */` | `/* === */` | `/* === */` |
| **Dockerfile** | `# ===` | `# ===` | `# ===` | `# ===` |

### Section Equivalences

<details>
<summary><b>METADATA Block — Identity across languages</b></summary>

| Section | Go | C | Assembly | Shell |
|---------|-----|---|----------|-------|
| **Package/File doc** | `// Package X...` | `/* @file X... */` | `; File: X...` | `# Script: X...` |
| **Biblical Foundation** | `// Scripture:` | `/* Scripture: */` | `; Scripture:` | `# Scripture:` |
| **CPI-SI Identity** | `// Component Type:` | `/* Component: */` | `; Component:` | `# Component:` |
| **Authorship** | `// Architect:` | `/* Author: */` | `; Author:` | `# Author:` |
| **Purpose** | `// Purpose:` | `/* Purpose: */` | `; Purpose:` | `# Purpose:` |
| **Health Scoring** | `// Health:` | `/* Health: */` | `; Health:` | `# Health:` |

</details>

<details>
<summary><b>SETUP Block — Dependencies across languages</b></summary>

| Concept | Go | C | Assembly | Shell |
|---------|-----|---|----------|-------|
| **Imports/Includes** | `import (...)` | `#include <...>` | `%include "..."` | `source ...` |
| **Constants** | `const X = Y` | `#define X Y` | `X equ Y` | `readonly X=Y` |
| **Types** | `type X struct` | `typedef struct` | `struc X` | N/A |
| **Variables** | `var x Type` | `Type x;` | `x: dd 0` | `X=""` |
| **Package State** | `var logger = ...` | `static Logger* log` | `section .data` | Global vars |

</details>

<details>
<summary><b>BODY Block — Implementation across languages</b></summary>

| Concept | Go | C | Assembly | Shell |
|---------|-----|---|----------|-------|
| **Section divider** | `// ────` | `/* ──── */` | `; ────` | `# ────` |
| **Private/Helper** | `func helper()` | `static void helper()` | `helper:` (no global) | `_helper()` |
| **Public/Export** | `func Export()` | `void Export()` | `global Export` | `export -f func` |
| **Error handling** | `if err != nil` | `if (err)` | `jc .error` | `if [[ $? -ne 0 ]]` |

</details>

<details>
<summary><b>CLOSING Block — Execution across languages</b></summary>

| Concept | Go | C | Assembly | Shell |
|---------|-----|---|----------|-------|
| **Entry point** | `func main()` | `int main()` | `_start:` or `main:` | `main "$@"` |
| **Validation** | `go build` | `gcc -Wall` | `nasm -f elf64` | `bash -n script.sh` |
| **Exit** | `os.Exit(0)` | `return 0;` | `mov eax, 60; syscall` | `exit 0` |

</details>

### Reserved Section Syntax

| Language | Reserved Marker |
|----------|-----------------|
| **Go** | `// [Reserved: reason]` |
| **C** | `/* [Reserved: reason] */` |
| **Assembly** | `; [Reserved: reason]` |
| **Shell** | `# [Reserved: reason]` |
| **Makefile** | `# [Reserved: reason]` |
| **YAML/TOML** | `# [Reserved: reason]` |
| **JSON** | `"_reserved": "reason"` |

### Import/Include Grouping

| Language | Group 1 | Group 2 | Group 3 |
|----------|---------|---------|---------|
| **Go** | Standard library | Internal packages | External packages |
| **C** | System headers `<>` | Project headers `""` | Local headers `""` |
| **Assembly** | CPU includes | Project includes | Local includes |
| **Shell** | Built-ins | System scripts | Project scripts |
| **Python** | Standard library | Third-party | Local modules |

---

## When to Use This Skill

### Decision Flow

```text
Creating new code file?
    │
    ├─ YES ──► Does a template exist? (see Quick Reference)
    │              │
    │              ├─ YES ──► USE THIS SKILL
    │              │
    │              └─ NO ───► Create file, consider if template should exist
    │
    └─ NO ───► Major structural work? (alignment, splitting, restructuring)
                   │
                   ├─ YES ──► USE THIS SKILL
                   │
                   └─ NO ───► Use Edit tool directly
```

### Always Use This Skill For

| Category | Examples |
|----------|----------|
| **New code files** | Any `.go`, `.c`, `.h`, `.asm`, `.sh`, `.ld` file |
| **Config files** | `go.mod`, `Makefile`, `CMakeLists.txt`, `.toml`, `.yaml` |
| **Docker files** | `Dockerfile`, `docker-compose.yaml` |
| **Dotfiles** | `.gitignore`, `.editorconfig`, `.env` |
| **Structural work** | Aligning to 4-block, splitting files, major refactoring |

### Do NOT Use For

- Simple edits to existing well-structured code → use `Edit` tool
- Reading/exploring code → use `Read`/`Grep` tools
- Debugging or bug fixes → direct intervention
- Adding small functions to existing 4-block files → use `Edit` tool

---

## Template Quick Reference

> **Skill Resources:** `~/.claude/skills/create-from-template/references/`

### By File Extension

| Extension | Template | Key | Location |
|-----------|----------|-----|----------|
| `.go` (main) | go-main-executable.go | LANG-TEMPLATE-002 | `compiler-templates/go/` |
| `.go` (package) | go-library-package.go | LANG-TEMPLATE-001 | `compiler-templates/go/` |
| `.go` (demo) | go-demo-test.go | LANG-TEMPLATE-003 | `compiler-templates/go/` |
| `.c` | c-source.c | LANG-TEMPLATE-009 | `compiler-templates/c/` |
| `.h` | c-header.h | LANG-TEMPLATE-010 | `compiler-templates/c/` |
| `.asm` | CODE-ASM-001-ASM-source.asm | CODE-ASM-001 | `code-templates/asm/` |
| `.inc` | CODE-ASM-002-ASM-include.inc | CODE-ASM-002 | `code-templates/asm/` |
| `.s` | CODE-ARM-001-ARM-source.s | CODE-ARM-001 | `code-templates/asm/` |
| `.sh` | CODE-SH-001-SH-script.sh | CODE-SH-001 | `code-templates/shell/` |
| `.ld` | CODE-LD-001-LD-linker-script.ld | CODE-LD-001 | `code-templates/linker/` |
| `Makefile` | make-project-makefile.mk | LANG-TEMPLATE-007 | `compiler-templates/build/` |
| `CMakeLists.txt` | cmake-project-cmakelists.cmake | LANG-TEMPLATE-008 | `compiler-templates/build/` |
| `.toml` | config-toml.toml | LANG-TEMPLATE-011 | `compiler-templates/config/` |
| `.yaml` | yaml-config.yaml | LANG-TEMPLATE-012 | `compiler-templates/yaml/` |
| `.json` | json-data.json | LANG-TEMPLATE-013 | `compiler-templates/json/` |
| `.jsonc` | jsonc-config.jsonc | LANG-TEMPLATE-015 | `compiler-templates/json/` |
| `Dockerfile` | dockerfile.dockerfile | LANG-TEMPLATE-016 | `compiler-templates/docker/` |
| `.gitignore` | CODE-GIT-001-GIT-gitignore.gitignore | CODE-GIT-001 | `code-templates/dotfiles/` |
| `.editorconfig` | CODE-GIT-002-GIT-editorconfig.editorconfig | CODE-GIT-002 | `code-templates/dotfiles/` |
| `.env` | CODE-GIT-003-GIT-env.env | CODE-GIT-003 | `code-templates/dotfiles/` |

### By Purpose

| Purpose | Template | Key |
|---------|----------|-----|
| **Go executable** (has `main()`) | go-main-executable.go | LANG-TEMPLATE-002 |
| **Go library** (reusable package) | go-library-package.go | LANG-TEMPLATE-001 |
| **Go demo/test** (showcase with `RunX()`) | go-demo-test.go | LANG-TEMPLATE-003 |
| **Go module** | go-module.mod | LANG-TEMPLATE-005 |
| **Go workspace** | go-workspace.work | LANG-TEMPLATE-004 |
| **C implementation** | c-source.c | LANG-TEMPLATE-009 |
| **C header/API** | c-header.h | LANG-TEMPLATE-010 |
| **x86 NASM assembly** | CODE-ASM-001-ASM-source.asm | CODE-ASM-001 |
| **ARM assembly** | CODE-ARM-001-ARM-source.s | CODE-ARM-001 |
| **Assembly macros/constants** | CODE-ASM-002-ASM-include.inc | CODE-ASM-002 |
| **Bash script** | CODE-SH-001-SH-script.sh | CODE-SH-001 |
| **Linker script** | CODE-LD-001-LD-linker-script.ld | CODE-LD-001 |
| **Project build** | make-project-makefile.mk | LANG-TEMPLATE-007 |
| **Health scoring data** | json-health-map.json | LANG-TEMPLATE-014 |
| **Go Docker build** | go-builder.dockerfile | LANG-TEMPLATE-017 |

---

## The Creation Process

### Overview

```text
Step 0: Pre-Check ──► Does file exist?
           │
           ▼
Step 1: Recognize ──► What am I creating?
           │
           ▼
Step 2: Find ──────► Which template?
           │
           ▼
Step 3: Prepare ───► Copy template, set up workspace
           │
           ▼
Step 3.5: Plan ────► Expand TODOs for each block/section
           │
           ▼
Step 4: Process ───► METADATA → SETUP → BODY → CLOSING
           │
           ▼
Step 5: Validate ──► Structural + Kingdom checklist
           │
           ▼
Step 6: Finalize ──► Build, verify, commit
```

---

### Step 0: Pre-Check

> **Question:** Does the target file already exist?

```bash
ls -la /path/to/target/filename.go
```

| Result | Process |
|--------|---------|
| **File doesn't exist** | Copy template → Fill block-by-block |
| **File exists** (needs alignment) | Read existing → Copy template to scratch → Migrate block-by-block |

---

### Step 1: Recognize

> **Question:** What am I creating?

Document before proceeding:

- [ ] **File type:** executable, library, demo-test, config, header?
- [ ] **Purpose:** What does this file do in the system?
- [ ] **Architecture fit:** Where does it belong? What uses it?
- [ ] **Consumer:** Who/what will use this file?

---

### Step 2: Find

> **Question:** Which template serves this?

**Use the [Template Quick Reference](#template-quick-reference) above.**

```bash
SKILL_DIR=~/.claude/skills/create-from-template/references

# Browse templates
ls $SKILL_DIR/compiler-templates/      # LANG-TEMPLATE-* (working copies)
ls $SKILL_DIR/code-templates/          # CODE-* (canonical masters)

# View specific category
ls $SKILL_DIR/compiler-templates/go/
ls $SKILL_DIR/code-templates/asm/
```

---

### Step 3: Prepare

> **Action:** Copy template to target location

**For NEW files:**

```bash
SKILL_DIR=~/.claude/skills/create-from-template/references

# Copy template
cp $SKILL_DIR/compiler-templates/go/go-library-package.go /path/to/target/newfile.go

# For Go files: Remove build ignore tag
# Edit to remove: //go:build ignore
```

**For EXISTING files (alignment/migration):**

```bash
SKILL_DIR=~/.claude/skills/create-from-template/references

# 1. Read existing file thoroughly first
cat /path/to/existing/file.go

# 2. Copy template to scratch (not target directory)
cp $SKILL_DIR/compiler-templates/go/go-library-package.go /tmp/template-scratch.go

# 3. Migrate existing code block-by-block into scratch
# 4. Replace original only when complete and verified
```

---

### Step 3.5: Plan the Work

> **Action:** Expand TODOs to track each block and section

**CRITICAL:** Create granular tasks before processing:

```text
1. [pending] METADATA: Package doc and Biblical Foundation
2. [pending] METADATA: CPI-SI Identity and Authorship
3. [pending] METADATA: Purpose, Dependencies, Health Scoring
4. [pending] SETUP: Imports (stdlib, internal, external)
5. [pending] SETUP: Types and Type Methods
6. [pending] SETUP: Constants, Variables, Package-Level State
7. [pending] BODY: Organizational Chart
8. [pending] BODY: Helpers/Utilities
9. [pending] BODY: Core Operations
10. [pending] BODY: Error Handling and Public APIs
11. [pending] CLOSING: Validation, Execution, Cleanup
12. [pending] CLOSING: Final Documentation sections
13. [pending] Build and verify
```

Mark each as `in_progress` → `completed` as you work through.

---

### Step 4: Process Block-by-Block

> **Action:** Fill each block in order. Never skip ahead.

<details>
<summary><b>4.1 METADATA Block</b> — Identity, purpose, health scoring</summary>

| Section | What to Fill |
|---------|--------------|
| Package doc comment | File purpose (first line) |
| Biblical Foundation | Scripture + Principle + Anchor |
| CPI-SI Identity | Component Type (Ladder/Baton/Rails) + Role + Paradigm |
| Authorship & Lineage | Architect + Implementation + Created + Version + Modified |
| Purpose & Function | Purpose + Core Design + Key Features |
| Blocking Status | Does this block execution? (if applicable) |
| Usage & Integration | Commands, API calls, exit codes |
| Dependencies | What This Needs + What Uses This |
| Health Scoring | System + States + Operations (Base100) |

**For empty sections:** `// [Reserved: Brief reason]`

</details>

<details>
<summary><b>4.2 SETUP Block</b> — Imports, types, constants</summary>

| Section | What to Fill |
|---------|--------------|
| Imports - Standard Library | With purpose comments |
| Imports - Internal Packages | With purpose comments |
| Imports - External Packages | Or `[Reserved: None needed]` |
| Constants | Or `[Reserved]` |
| Variables | Or `[Reserved]` |
| Types | Structs, interfaces, type aliases |
| Type Methods | Methods on types |
| Package-Level State | Rails infrastructure (loggers, etc.) |

**Import grouping pattern:**

```go
//--- Standard Library ---
import (
    "fmt"     // Formatted output
    "strings" // String manipulation
)

//--- Internal Packages ---
import (
    "module/pkg/health" // Health scoring
)

//--- External Packages ---
// [Reserved: No external dependencies]
```

</details>

<details>
<summary><b>4.3 BODY Block</b> — Core implementation</summary>

| Section | What to Fill |
|---------|--------------|
| Organizational Chart | Document ladder/baton/rails flow |
| Helpers/Utilities | Private helper functions |
| Core Operations | Business logic by category |
| Error Handling | Error types, handlers |
| Public APIs | Exported functions |

**Category organization:**

```go
// ────────────────────────────────────────────────────────────────
// [Category Name] - [Purpose]
// ────────────────────────────────────────────────────────────────
```

</details>

<details>
<summary><b>4.4 CLOSING Block</b> — Execution, validation, documentation</summary>

| Section | What to Fill |
|---------|--------------|
| Code Validation | Build/test commands |
| Code Execution | Entry point (`main()` or `init()`) |
| Code Cleanup | Resource management |
| Library/Executable Overview | Reference to METADATA |
| Modification Policy | Safe / Care / Never modify |
| Ladder and Baton Flow | Reference to BODY organization |
| Surgical Update Points | Extension guidance |
| Performance Considerations | Notes or reference |
| Troubleshooting Guide | Common issues and solutions |
| Related Components | Dependencies reference |
| Future Expansions | Planned features |
| Closing Note | Summary + Scripture |
| Quick Reference | Usage examples |

</details>

---

### Step 5: Validate

> **Question:** Does it honor the pattern?

**Structural Checklist:**

- [ ] 4-block structure complete (METADATA → SETUP → BODY → CLOSING)
- [ ] ALL sections present (filled or `[Reserved]`)
- [ ] All METADATA sections filled (Biblical, CPI-SI, Authorship, Purpose, Health)
- [ ] Imports properly grouped (stdlib / internal / external)
- [ ] Code compiles without warnings
- [ ] Standards compliance (CWS-STD-001, 002, 003)

**For Migrations:**

- [ ] All original functionality preserved
- [ ] Structure aligned to template
- [ ] No content lost

**Kingdom Checklist (The Three Questions):**

1. Would this honor God as code reviewer?
2. Does this genuinely serve others?
3. Does this have eternal value?

---

### Step 6: Finalize

> **Action:** Build, verify, complete

**For NEW files:**

```bash
go build ./path/to/new/file.go
# or
make build
```

**For EXISTING files (alignment):**

```bash
# 1. Verify scratch builds
go build /tmp/template-scratch.go

# 2. Replace original
mv /tmp/template-scratch.go /path/to/original/file.go

# 3. Final verification
go build ./...
```

---

## 4-Block Structure Reference

> **Standard:** CWS-STD-001 — See `references/docs/CWS-STD-001-DOC-4-block.md`

```text
┌─────────────────────────────────────────────────────────────────┐
│ METADATA                                                         │
│   Biblical Foundation, CPI-SI Identity, Authorship,             │
│   Purpose, Blocking Status, Usage, Dependencies, Health         │
├─────────────────────────────────────────────────────────────────┤
│ SETUP                                                            │
│   Imports (stdlib / internal / external), Constants,            │
│   Variables, Types, Type Methods, Package-Level State           │
├─────────────────────────────────────────────────────────────────┤
│ BODY                                                             │
│   Organizational Chart, Helpers, Core Operations,               │
│   Error Handling, Public APIs                                   │
├─────────────────────────────────────────────────────────────────┤
│ CLOSING                                                          │
│   Validation, Execution, Cleanup, Documentation,                │
│   Modification Policy, Troubleshooting, Closing Note            │
└─────────────────────────────────────────────────────────────────┘
```

### Reserved Section Pattern

When a section isn't needed:

```go
// ────────────────────────────────────────────────────────────────
// [Section Name]
// ────────────────────────────────────────────────────────────────
// [Reserved: Brief reason why not needed for this file]
```

**Examples:**

```go
// [Reserved: No external packages - stdlib only]
// [Reserved: No package-level state - stateless functions]
// [Reserved: No custom types - uses imported types]
```

---

## Embedded Resources

> **Location:** `~/.claude/skills/create-from-template/references/`

| Resource | Path | Contents |
|----------|------|----------|
| **Compiler Templates** | `compiler-templates/` | Working copies (LANG-TEMPLATE-*) — 21 templates |
| **Root Templates** | `code-templates/` | Canonical masters (CODE-*) — 26 templates |
| **TEMPLATE-MATRIX** | `docs/TEMPLATE-MATRIX.md` | Full inventory (89 templates) |
| **CODE-ALIGNMENT-MATRIX** | `docs/CODE-ALIGNMENT-MATRIX.md` | Files needing templates (257+ files) |
| **4-Block Standard** | `docs/CWS-STD-001-DOC-4-block.md` | Structure specification |

### Template Categories

| Category | Count | Templates |
|----------|:-----:|-----------|
| **Go** | 6 | executable, library, demo-test, module, workspace, checksum |
| **C** | 2 | source, header |
| **Assembly** | 3 | x86 NASM, ARM GNU AS, include |
| **Shell** | 1 | bash script |
| **Linker** | 1 | linker script |
| **Build** | 2 | Makefile, CMakeLists |
| **Docker** | 3 | Dockerfile, Go builder, Compose |
| **Config** | 2 | TOML, YAML |
| **JSON** | 3 | data, health map, JSONC |
| **Dotfiles** | 3 | gitignore, editorconfig, env |
| **Total** | **26** | Canonical templates |

### Sync Templates

If templates in skill are outdated:

```bash
# Sync from workspace
WORKSPACE=/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC
SKILL_DIR=~/.claude/skills/create-from-template/references

cp -r $WORKSPACE/templates/code/* $SKILL_DIR/code-templates/
cp -r $WORKSPACE/divisions/tech/language/compiler/templates/* $SKILL_DIR/compiler-templates/
cp $WORKSPACE/templates/TEMPLATE-MATRIX.md $SKILL_DIR/docs/
cp $WORKSPACE/templates/CODE-ALIGNMENT-MATRIX.md $SKILL_DIR/docs/
```

---

## Full Template Inventory

<details>
<summary><b>Compiler Templates (21) — LANG-TEMPLATE-*</b></summary>

### Go Templates (6)

| Key | File | Purpose | Derived From |
|-----|------|---------|--------------|
| LANG-TEMPLATE-001 | `go/go-library-package.go` | Library packages | CODE-GO-002 |
| LANG-TEMPLATE-002 | `go/go-main-executable.go` | Entry point programs | CODE-GO-001 |
| LANG-TEMPLATE-003 | `go/go-demo-test.go` | Demo-test files | CODE-GO-003 |
| LANG-TEMPLATE-004 | `go/go-workspace.work` | Workspace config | CODE-GO-004 |
| LANG-TEMPLATE-005 | `go/go-module.mod` | Module config | CODE-GO-005 |
| LANG-TEMPLATE-006 | `go/go-checksum.sum` | Checksum docs | CODE-GO-006 |

### Build Templates (2)

| Key | File | Purpose | Derived From |
|-----|------|---------|--------------|
| LANG-TEMPLATE-007 | `build/make-project-makefile.mk` | Makefile | CODE-MAKE-001 |
| LANG-TEMPLATE-008 | `build/cmake-project-cmakelists.cmake` | CMake | CODE-CMAKE-001 |

### C Templates (2)

| Key | File | Purpose | Derived From |
|-----|------|---------|--------------|
| LANG-TEMPLATE-009 | `c/c-source.c` | C implementation | CODE-C-001 |
| LANG-TEMPLATE-010 | `c/c-header.h` | C header | CODE-C-002 |

### Config Templates (2)

| Key | File | Purpose | Derived From |
|-----|------|---------|--------------|
| LANG-TEMPLATE-011 | `config/config-toml.toml` | TOML config | CODE-CONFIG-001 |
| LANG-TEMPLATE-012 | `yaml/yaml-config.yaml` | YAML config | CODE-YAML-001 |

### JSON Templates (3)

| Key | File | Purpose | Derived From |
|-----|------|---------|--------------|
| LANG-TEMPLATE-013 | `json/json-data.json` | JSON data | CODE-JSON-001 |
| LANG-TEMPLATE-014 | `json/json-health-map.json` | Health maps | CODE-JSON-003 |
| LANG-TEMPLATE-015 | `json/jsonc-config.jsonc` | JSONC config | CODE-JSONC-001 |

### Docker Templates (3)

| Key | File | Purpose | Derived From |
|-----|------|---------|--------------|
| LANG-TEMPLATE-016 | `docker/dockerfile.dockerfile` | Base Dockerfile | CODE-DOCKER-001 |
| LANG-TEMPLATE-017 | `docker/go-builder.dockerfile` | Go multi-stage | CODE-DOCKER-002 |
| LANG-TEMPLATE-018 | `docker/docker-compose.yaml` | Compose | CODE-DOCKER-003 |

### Dotfiles Templates (3)

| Key | File | Purpose | Derived From |
|-----|------|---------|--------------|
| LANG-TEMPLATE-019 | `dotfiles/gitignore.gitignore` | Git ignore | CODE-GIT-001 |
| LANG-TEMPLATE-020 | `dotfiles/editorconfig.editorconfig` | Editor config | CODE-GIT-002 |
| LANG-TEMPLATE-021 | `dotfiles/env.env` | Environment | CODE-GIT-003 |

</details>

<details>
<summary><b>Root Templates (26) — CODE-*</b></summary>

**Location:** `references/code-templates/`

| Category | Key | File | Purpose |
|----------|-----|------|---------|
| **Go** | CODE-GO-001 | `go/CODE-GO-001-GO-executable.go` | Go executable |
| | CODE-GO-002 | `go/CODE-GO-002-GO-library.go` | Go library |
| | CODE-GO-003 | `go/CODE-GO-003-GO-demo-test.go` | Go demo-test |
| | CODE-GO-004 | `go/CODE-GO-004-GO-workspace.work` | Go workspace |
| | CODE-GO-005 | `go/CODE-GO-005-GO-module.mod` | Go module |
| | CODE-GO-006 | `go/CODE-GO-006-GO-checksum.sum` | Go checksum |
| **C** | CODE-C-001 | `c/CODE-C-001-C-source.c` | C source |
| | CODE-C-002 | `c/CODE-C-002-C-header.h` | C header |
| **Assembly** | CODE-ASM-001 | `asm/CODE-ASM-001-ASM-source.asm` | x86 NASM |
| | CODE-ASM-002 | `asm/CODE-ASM-002-ASM-include.inc` | ASM include |
| | CODE-ARM-001 | `asm/CODE-ARM-001-ARM-source.s` | ARM GNU AS |
| **Shell** | CODE-SH-001 | `shell/CODE-SH-001-SH-script.sh` | Bash script |
| **Linker** | CODE-LD-001 | `linker/CODE-LD-001-LD-linker-script.ld` | Linker script |
| **Build** | CODE-MAKE-001 | `make/CODE-MAKE-001-MAKE-project-makefile.mk` | Makefile |
| | CODE-CMAKE-001 | `cmake/CODE-CMAKE-001-CMAKE-project-cmakelists.cmake` | CMake |
| **Config** | CODE-CONFIG-001 | `config/CODE-CONFIG-001-CONFIG-toml.toml` | TOML |
| | CODE-YAML-001 | `yaml/CODE-YAML-001-YAML-config.yaml` | YAML |
| **JSON** | CODE-JSON-001 | `json/CODE-JSON-001-JSON-data.json` | JSON data |
| | CODE-JSON-003 | `json/CODE-JSON-003-JSON-health-map.json` | Health map |
| | CODE-JSONC-001 | `json/CODE-JSONC-001-JSONC-config.jsonc` | JSONC |
| **Docker** | CODE-DOCKER-001 | `docker/CODE-DOCKER-001-DOCKER-dockerfile.dockerfile` | Dockerfile |
| | CODE-DOCKER-002 | `docker/CODE-DOCKER-002-DOCKER-go-builder.dockerfile` | Go builder |
| | CODE-DOCKER-003 | `docker/CODE-DOCKER-003-DOCKER-compose.yaml` | Compose |
| **Dotfiles** | CODE-GIT-001 | `dotfiles/CODE-GIT-001-GIT-gitignore.gitignore` | Gitignore |
| | CODE-GIT-002 | `dotfiles/CODE-GIT-002-GIT-editorconfig.editorconfig` | Editorconfig |
| | CODE-GIT-003 | `dotfiles/CODE-GIT-003-GIT-env.env` | Environment |

</details>

---

## Language-Specific Guides

<details>
<summary><b>Go Files</b></summary>

### Choosing the Right Template

| If your file... | Use Template | Key |
|-----------------|--------------|-----|
| Has `package main` and `func main()` | go-main-executable.go | LANG-TEMPLATE-002 |
| Is a reusable library package | go-library-package.go | LANG-TEMPLATE-001 |
| Demonstrates/showcases functionality | go-demo-test.go | LANG-TEMPLATE-003 |

### Go-Specific Sections

**METADATA:**
- Package doc must be directly above `package` line
- No blank lines between doc comment and package declaration

**SETUP:**
- Imports grouped: stdlib → internal → external
- Each group separated by blank line
- Purpose comment after each import

**BODY:**
- Exported functions documented with godoc
- Unexported helpers prefixed with lowercase

**CLOSING:**
- For executables: `main()` function
- For libraries: init() if needed, otherwise [Reserved]

### Build Commands

```bash
# Single file
go build ./path/to/file.go

# Package
go build ./path/to/package/

# All packages
go build ./...

# With race detection
go build -race ./...
```

</details>

<details>
<summary><b>C Files</b></summary>

### Choosing the Right Template

| If your file... | Use Template | Key |
|-----------------|--------------|-----|
| Is implementation (`.c`) | c-source.c | LANG-TEMPLATE-009 |
| Is header/API (`.h`) | c-header.h | LANG-TEMPLATE-010 |

### C-Specific Sections

**METADATA:**
- File doc block at top with `/* */`
- Include guards in headers: `#ifndef PROJECT_HEADER_H`

**SETUP:**
- Includes grouped: system → project → local
- Macros and constants
- Type definitions

**BODY:**
- Static functions (private) first
- Public functions last
- Function doc blocks

**CLOSING:**
- No execution block (library code)
- Include guard close: `#endif`

</details>

<details>
<summary><b>Assembly Files</b></summary>

### Choosing the Right Template

| If your file... | Use Template | Key |
|-----------------|--------------|-----|
| Is x86 NASM code (`.asm`) | CODE-ASM-001-ASM-source.asm | CODE-ASM-001 |
| Is ARM GNU AS code (`.s`) | CODE-ARM-001-ARM-source.s | CODE-ARM-001 |
| Contains macros/constants (`.inc`) | CODE-ASM-002-ASM-include.inc | CODE-ASM-002 |

### Assembly-Specific Sections

**METADATA:**
- Comment block with `;` prefix (NASM) or `@` prefix (ARM)
- Target architecture documented
- Calling conventions noted

**SETUP:**
- Section declarations (`.text`, `.data`, `.bss`)
- External symbols
- Constants and macros

**BODY:**
- Code organized by function
- Clear labels with purpose comments
- Prologue/epilogue for functions

**CLOSING:**
- Entry point if bootloader/kernel
- Export declarations

</details>

<details>
<summary><b>Shell Scripts</b></summary>

### Template

| Use Case | Template | Key |
|----------|----------|-----|
| Any bash script | CODE-SH-001-SH-script.sh | CODE-SH-001 |

### Shell-Specific Sections

**METADATA:**
- Shebang: `#!/usr/bin/env bash`
- Script doc block with `#`
- Usage information

**SETUP:**
- `set -euo pipefail` (strict mode)
- Constants (readonly)
- Function definitions

**BODY:**
- Main logic
- Function implementations
- Error handling

**CLOSING:**
- Main execution (`main "$@"`)
- Exit codes documented

</details>

---

## Examples

<details>
<summary><b>Example: Creating a Go Library Package</b></summary>

### Step 0 - Pre-Check

```bash
ls -la pkg/health/scoring.go
# Result: No such file → NEW FILE
```

### Step 1 - Recognize

- **File type:** Library package
- **Purpose:** Health scoring calculations
- **Architecture:** Part of pkg/ infrastructure (Rails)
- **Consumer:** All components needing health tracking

### Step 2 - Find

> Use `go-library-package.go` (LANG-TEMPLATE-001)

### Step 3 - Prepare

```bash
SKILL_DIR=~/.claude/skills/create-from-template/references
mkdir -p pkg/health/
cp $SKILL_DIR/compiler-templates/go/go-library-package.go pkg/health/scoring.go
```

### Step 4 - Process

**4.1 METADATA:**
- Package doc: "Package health provides health scoring calculations."
- Biblical: Proverbs 4:23 (guard your heart - health tracking)
- CPI-SI: Rails (infrastructure for health system)
- Authorship: Nova Dawn, 2025-11-30
- Purpose: Calculate and aggregate health scores
- Dependencies: None (standalone)
- Health: Self-scoring component

**4.2 SETUP:**
- Imports: `math` (calculations), `sync` (thread safety)
- Types: `Score`, `ScoreSet`, `Aggregator`
- Constants: `MaxScore = 100`, `MinScore = -100`

**4.3 BODY:**
- Helpers: `clamp()`, `validate()`
- Core: `Calculate()`, `Aggregate()`, `Normalize()`
- Public: `NewScore()`, `NewAggregator()`

**4.4 CLOSING:**
- Validation: `go build ./pkg/health/`
- Execution: [Reserved - library]
- Closing Note: Proverbs 4:23

### Step 5 - Validate

- [x] 4-block complete
- [x] All sections present
- [x] Compiles: `go build ./pkg/health/`
- [x] Tests pass: `go test ./pkg/health/`

</details>

<details>
<summary><b>Example: Creating an Assembly Bootloader</b></summary>

### Step 0 - Pre-Check

```bash
ls -la boot/stage1/boot.asm
# Result: No such file → NEW FILE
```

### Step 1 - Recognize

- **File type:** x86 NASM assembly
- **Purpose:** Stage 1 bootloader (MBR)
- **Architecture:** First code executed at boot
- **Consumer:** BIOS loads this at 0x7C00

### Step 2 - Find

> Use `CODE-ASM-001-ASM-source.asm` (CODE-ASM-001)

### Step 3 - Prepare

```bash
SKILL_DIR=~/.claude/skills/create-from-template/references
mkdir -p boot/stage1/
cp $SKILL_DIR/code-templates/asm/CODE-ASM-001-ASM-source.asm boot/stage1/boot.asm
```

### Step 4 - Process

**4.1 METADATA:**
- File doc: Stage 1 bootloader
- Biblical: Genesis 1:1 (in the beginning)
- Target: x86 real mode, 16-bit
- Entry: 0x7C00

**4.2 SETUP:**
- `BITS 16`, `ORG 0x7C00`
- Constants: sector size, load address
- External: stage2 symbols

**4.3 BODY:**
- Initialize segments
- Load stage 2
- Jump to stage 2

**4.4 CLOSING:**
- Boot signature: `0x55, 0xAA`
- Pad to 512 bytes

</details>

---

## Common Mistakes

| Mistake | Why Wrong | Correct Approach |
|---------|-----------|------------------|
| Skipping sections | Loses structure | Use `[Reserved: reason]` |
| Starting from scratch | Ignores pattern | Always start from template |
| Bulk copying | No intentional thought | Block by block |
| Ignoring existing code | Loses functionality | Migrate carefully |
| Combining import blocks | Loses organization | Separate stdlib/internal/external |
| No TODO tracking | Lose progress visibility | Expand TODOs first |

---

## Biblical Foundation

> *"And look that thou make them after their **pattern**, which was shewed thee in the mount."* — Exodus 25:40

> *"All this, said David, the LORD made me understand in writing by his hand upon me, even all the works of this **pattern**."* — 1 Chronicles 28:19

> *"And I have filled him with the spirit of God, in **wisdom**, and in **understanding**, and in **knowledge**, and in all manner of **workmanship**."* — Exodus 31:3

**The Pattern Principle:** God gave Moses the pattern of the tabernacle — not to restrict creativity, but to enable it. The pattern provided structure; Bezalel filled it with Spirit-empowered excellence.

**This skill embodies both:**
- **Faithfulness** — Follow the revealed pattern (templates, 4-block, standards)
- **Excellence** — Fill the structure with wisdom for THIS specific file

---

> *"According to all that I shew thee, after the pattern of the tabernacle, and the pattern of all the instruments thereof, even so shall ye make it."* — Exodus 25:9
