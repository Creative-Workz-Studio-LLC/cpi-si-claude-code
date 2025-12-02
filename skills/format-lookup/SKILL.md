---
name: format-lookup
description: Quick lookup of OmniCode format mappings - how OmniCode sections translate to specific formats (AsciiDoc, Go, C, etc.). Use this skill when Nova Dawn needs to know how an OmniCode element maps to a target format.
version: "1.1.0"
allowed-tools: "Read,Grep,Bash"
---

# Format Lookup

> *"To every thing there is a season, and a time to every purpose under the heaven."* — Ecclesiastes 3:1

Quick lookup of how OmniCode maps to target formats.

---

## References

This skill has reference files in `~/.claude/skills/format-lookup/references/`:

| File | Contents |
|------|----------|
| **format-mappings.txt** | Complete OmniCode-to-format mapping tables |
| **linux-commands.txt** | Bash commands for format lookups |

---

## Quick Start

```bash
# === FIND FORMAT DEFINITIONS ===
grep -A 30 'format "adoc"' bereshit/word/omni/B-word-omni-syntax.omni

# === LOOKUP SPECIFIC MAPPING ===
grep -n "biblical" bereshit/word/seed/documentation/B-word-seed-doc-asciidoc-base.adoc

# === LIST ALL DEFINED FORMATS ===
grep -oP 'format "\K[^"]+' bereshit/word/omni/B-word-omni-syntax.omni

# === FIND TEMPLATE FOR FORMAT ===
find bereshit/word/seed/ -name "*.adoc" -type f
```

---

## When to Use

**Use this skill when:**
- Need to know how an OmniCode element maps to a format
- Writing format-aware OmniCode
- Checking correct translation
- Quick reference during template work

**Don't use for:**
- Full format analysis (use format-bridge agent)
- Validating structure (use validate-omni skill)
- Creating files (use create-from-template skill)

---

## Format Mapping Reference

### Documentation Formats

#### AsciiDoc (--adoc)

| OmniCode | AsciiDoc | Example |
|----------|----------|---------|
| `#!omni documentation --adoc` | `////` block with pragma | `////\n#!omni documentation --adoc\n////` |
| `this is a [document] for` | `:title:` | `:title: Document Name` |
| `keyed as` | `:key:` | `:key: B-WORD-DOC-001` |
| `following` | `:type:` | `:type: Reference` |
| `classifies as` | `:tags:` | `:tags: documentation, reference` |
| `derives from` | `:derives_from:` | `:derives_from: path/to/template.omni` |
| `strictness is` | `:strictness:` | `:strictness: tight` |
| `grounded in: scripture` | `:biblical_foundation:` | `:biblical_foundation: Genesis 1:1` |
| `authored by: architect` | `:author:` | `:author: Seanje Lenox-Wise` |
| `authored by: writer` | `:authors:` | `:authors: Nova Dawn` |
| `created` | `:created:` | `:created: 2025-11-30` |
| `version` | `:version:` | `:version: 1.0.0` |
| `modified` | `:updated:` | `:updated: 2025-12-01` |
| `exists to: purpose` | `:purpose:` | `:purpose: Define syntax` |
| `exists to: design` | `:description:` | `:description: Spec by example` |
| `presented as: title` | `= Title` | `= Document Title` |
| `presented as: tagline` | `_tagline_` | `_Kingdom Technology_` |
| `structure of this document` | Section organization | `== Section`, `=== Subsection` |
| `clarifications` | Admonition blocks | `NOTE:`, `TIP:`, `CAUTION:` |
| `closing note` | Final paragraph | Scripture reference |

#### Markdown (--md) [PLANNED]

| OmniCode | Markdown | Example |
|----------|----------|---------|
| `#!omni documentation --md` | HTML comment | `<!--\n#!omni documentation --md\n-->` |
| `this is a [document] for` | YAML `title:` | `title: Document Name` |
| `keyed as` | YAML `key:` | `key: B-WORD-DOC-001` |
| `grounded in: scripture` | YAML `biblical_foundation:` | `biblical_foundation: Genesis 1:1` |
| `presented as: title` | `# Title` | `# Document Title` |
| `structure of this document` | `##`, `###` headers | `## Section` |

### Code Formats

#### Go (--go)

| OmniCode | Go | Example |
|----------|-----|---------|
| `#!omni code --go` | Generated comment | `// Code generated from OmniCode.` |
| `this is a [code] for` | Package doc comment | `// Package name provides...` |
| `keyed as` | `// Key:` comment | `// Key: PKG-HEALTH-001` |
| `grounded in: scripture` | `// Biblical:` comment | `// Biblical: Proverbs 4:23` |
| `authored by` | `// Authors:` comment | `// Authors: Nova Dawn` |
| `version` | `// Version:` comment | `// Version: 1.0.0` |
| `exists to: purpose` | Package doc first line | First sentence of package doc |
| `inputs: requires: from stdlib (1)` | Standard library imports | `import "fmt"` |
| `inputs: requires: from internal (0)` | Internal package imports | `import "project/pkg/health"` |
| `inputs: requires: from external (-1)` | Third-party imports | `import "github.com/..."` |
| `definitions: constants` | `const` block | `const MaxScore = 100` |
| `definitions: variables` | `var` block | `var registry = ...` |
| `definitions: assemblies` | `type struct` | `type Config struct {...}` |
| `definitions: methods` | Method receivers | `func (c *Config) Validate()` |
| `implementation: helpers` | Unexported functions | `func helper() {...}` |
| `implementation: operations` | Core functions | `func calculate() {...}` |
| `interface: exposes: exported` | Exported (Capitalized) | `func Calculate() {...}` |
| `interface: exposes: internal` | Unexported (lowercase) | `func validate() {...}` |
| `operations: executed through: main` | `func main()` | `func main() {...}` |
| `operations: cleaned by` | `defer` statements | `defer file.Close()` |

#### C (--c) [PLANNED]

| OmniCode | C | Example |
|----------|---|---------|
| `#!omni code --c` | Generated comment | `/* Code generated from OmniCode. */` |
| `this is a [code] for` | File header comment | `/* @file name.c ... */` |
| `inputs: requires` | `#include` statements | `#include <stdio.h>` |
| `definitions: constants` | `#define` macros | `#define MAX_SIZE 100` |
| `definitions: assemblies` | `typedef struct` | `typedef struct {...} Config;` |
| `implementation: helpers` | `static` functions | `static void helper() {...}` |
| `interface: exposes` | Non-static functions | `void public_func() {...}` |
| `operations: executed through: main` | `int main()` | `int main(int argc, char** argv)` |

#### Rust (--rs) [PLANNED]

| OmniCode | Rust | Example |
|----------|------|---------|
| `#!omni code --rs` | Doc comment | `//! Code generated from OmniCode.` |
| `this is a [code] for` | Module doc | `//! Module description` |
| `inputs: requires` | `use` statements | `use std::io;` |
| `definitions: assemblies` | `struct` | `pub struct Config {...}` |
| `implementation: helpers` | Private functions | `fn helper() {...}` |
| `interface: exposes` | `pub` functions | `pub fn calculate() {...}` |
| `operations: executed through: main` | `fn main()` | `fn main() {...}` |

---

## Lookup Process

### Step 1: Identify Source Element

What OmniCode element are you looking up?

```text
"grounded in: scripture"
"exists to: purpose"
"inputs: requires: from stdlib (1)"
```

### Step 2: Identify Target Format

What format are you mapping to?

```text
adoc, md, go, c, rs, etc.
```

### Step 3: Look Up Mapping

Check the tables above or search the syntax spec:

```bash
# Find format definitions in syntax spec
grep -A 50 "format definitions:" bereshit/word/omni/B-word-omni-syntax.omni

# Find specific format mapping
grep -A 30 'format "adoc"' bereshit/word/omni/B-word-omni-syntax.omni
```

### Step 4: Verify in Template

Check the format template for implementation:

```bash
# AsciiDoc template
cat bereshit/word/seed/documentation/B-word-seed-doc-asciidoc-base.adoc

# Go template
cat bereshit/word/seed/go/B-word-seed-go-executable.go
```

---

## Quick Reference Card

### Pragma Formats

| Type | Format Flag | Output |
|------|-------------|--------|
| documentation | `--adoc` | AsciiDoc |
| documentation | `--md` | Markdown |
| code | `--go` | Go |
| code | `--c` | C |
| code | `--rs` | Rust |
| data | `--json` | JSON |
| data | `--yaml` | YAML |
| data | `--toml` | TOML |

### Common Mappings Summary

| OmniCode Concept | AsciiDoc | Go | C |
|------------------|----------|-----|---|
| Title/Name | `:title:` | Package doc | File header |
| Biblical foundation | `:biblical_foundation:` | `// Biblical:` | `/* Biblical: */` |
| Purpose | `:purpose:` | First doc sentence | `@brief` |
| Imports/Includes | N/A | `import` | `#include` |
| Constants | N/A | `const` | `#define` |
| Types | N/A | `type struct` | `typedef struct` |
| Private functions | N/A | `func lower()` | `static void` |
| Public functions | N/A | `func Upper()` | `void func()` |
| Entry point | N/A | `func main()` | `int main()` |

---

## Source Files

```text
FORMAT DEFINITIONS location:
bereshit/word/omni/B-word-omni-syntax.omni
    → See "FORMAT DEFINITIONS (Transpiler Specification)" section

Format-specific templates:
bereshit/word/seed/documentation/   → AsciiDoc, Markdown bases
bereshit/word/seed/go/              → Go templates
bereshit/word/seed/c/               → C templates [PLANNED]
bereshit/word/seed/rust/            → Rust templates [PLANNED]
```

---

## Biblical Foundation

> *"To every thing there is a season, and a time to every purpose under the heaven."* — Ecclesiastes 3:1

Every format has its purpose. OmniCode doesn't replace formats - it bridges them. This skill helps you navigate that bridge, knowing how truth expressed in one form translates to another.

**This skill embodies:**
- **Translation** — Same truth, different expression
- **Wisdom** — Know the right form for the context
- **Unity** — One source, many presentations

---

> *"And the Word was made flesh, and dwelt among us."* — John 1:14
