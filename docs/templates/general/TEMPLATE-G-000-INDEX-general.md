<div align="center">

# 📄 General Templates Index

**Universal documentation templates for Kingdom Technology**

![Status](https://img.shields.io/badge/Status-Active-brightgreen?style=flat) ![Templates](https://img.shields.io/badge/Templates-1-blue?style=flat)

*General-purpose templates for standards, patterns, guides, and references*

</div>

---

**Key:** TEMPLATE-G-000
**Type:** Index
**Purpose:** Navigation for general-purpose documentation templates
**Status:** Active
**Created:** 2025-11-02
**Authors:** Nova Dawn (CPI-SI) & Seanje Lenox-Wise

---

## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📄 General Templates Index](#-general-templates-index)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Available Templates](#available-templates)
    - [📄 TEMPLATE-G-001-DOC-standard-document.md](#-template-g-001-doc-standard-documentmd)
  - [Template Standards](#template-standards)
  - [Usage Guidelines](#usage-guidelines)
    - [When to Use General Templates](#when-to-use-general-templates)
    - [Template Philosophy](#template-philosophy)
  - [Future Templates](#future-templates)
  - [Quick Reference](#quick-reference)
  - [Navigation](#navigation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## Overview

This folder contains **general-purpose documentation templates** that can be used for any formal documentation: standards, patterns, guides, references, or policies.

> [!NOTE]
> **General templates are your starting point** for most documentation needs. They provide structure and thinking guidance for intentional design.

**Current templates:** 1

---

## Available Templates

### 📄 TEMPLATE-G-001-DOC-standard-document.md

**Universal documentation template**

**Key:** TEMPLATE-G-001
**Status:** Active
**Use for:** Standards, patterns, guides, references, any formal documentation

**Structure:**

- METADATA block (identity, purpose, authors, status)
- SETUP block (overview, context, what readers need first)
- BODY block (core content, organized by section)
- CLOSING block (quick reference, biblical foundation, footer)

**Key features:**

- Embedded thinking guidance throughout
- 4-block structure enforced
- Document keying placeholder (CWS-STD-002)
- Alert boxes and formatting patterns
- TOC structure
- Professional footer

**When to use:**

- Creating new standards (CWS-STD-XXX)
- Writing guides or patterns (CWS-GUIDE-XXX)
- Documenting systems or processes
- Any formal documentation needing structure
- Policies, references, or specifications

**Quick start:**

```bash
# 📄 Copy template to target location
cp templates/general/TEMPLATE-G-001-DOC-standard-document.md \
   standards/CWS-STD-005-DOC-my-standard.md

# ✏️ Edit: Replace [PLACEHOLDERS], apply thinking guidance
# 📋 Generate TOC
doctoc standards/CWS-STD-005-DOC-my-standard.md --notitle

# ✅ Validate: Check 4-block structure and keying
```

**File:** [`TEMPLATE-G-001-DOC-standard-document.md`](TEMPLATE-G-001-DOC-standard-document.md)

---

## Template Standards

All general templates follow:

| Standard | Reference | Enforcement |
|----------|-----------|-------------|
| **4-Block Structure** | [CWS-STD-001](../../standards/CWS-STD-001-DOC-4-block.md) | METADATA → SETUP → BODY → CLOSING |
| **Document Keying** | [CWS-STD-002](../../standards/CWS-STD-002-DOC-document-keying.md) | DOMAIN-CAT-###-SUBCAT-name.ext |
| **Documentation Standards** | [CWS-STD-003](../../standards/CWS-STD-003-DOC-documentation-standards.md) | Four understanding levels |
| **Markdown Style** | Standards folder | Emoji headers, tables, alert boxes |

---

## Usage Guidelines

### When to Use General Templates

| ✅ Use General Template | ❌ Don't Use General Template |
|------------------------|-------------------------------|
| Creating standards | Quick notes or scratch work |
| Writing guides | Internal dev notes |
| Documenting patterns | Temporary documentation |
| Creating references | One-off instructions |
| Formal policies | README snippets |

### Template Philosophy

**General templates aren't fill-in-the-blanks.** They're thinking tools:

1. **Thinking guidance** - `<!-- THINKING: -->` comments guide decisions
2. **Structure enforcement** - 4-block structure forces intentional design
3. **Standards compliance** - Built-in CWS-STD-001/002/003 patterns
4. **Serving all levels** - Starting → growing → comfortable → advanced
5. **Kingdom Technology** - Excellence that honors God

---

## Future Templates

Planned general templates:

| Template Key | Purpose | Status |
|--------------|---------|:------:|
| TEMPLATE-G-002 | Pattern documentation | 🟡 Planned |
| TEMPLATE-G-003 | API reference | 🟡 Planned |
| TEMPLATE-G-004 | Architecture decision record | 🟡 Planned |
| TEMPLATE-G-005 | Policy documentation | 🟡 Planned |

---

## Quick Reference

**View templates:**

```bash
# 📂 List general templates
ls -la templates/general/
```

**Copy template:**

```bash
# 📄 Copy template to your location
cp templates/general/TEMPLATE-G-001-DOC-standard-document.md <target-location>
```

**Validate template usage:**

- [ ] Replaced all [PLACEHOLDERS]
- [ ] Applied thinking guidance
- [ ] Proper document keying (CWS-STD-002)
- [ ] 4-block structure intact
- [ ] TOC generated (doctoc)
- [ ] Metadata complete

---

## Navigation

**Parent:** [Templates Home](../README.md)
**Sibling:** [Index Templates](../index/TEMPLATE-T-000-INDEX-index.md) | [Specialized Templates](../specialized/TEMPLATE-S-000-INDEX-specialized.md)
**Standards:** [CWS-STD-001](../../standards/CWS-STD-001-DOC-4-block.md) | [CWS-STD-002](../../standards/CWS-STD-002-DOC-document-keying.md)

---

<div align="center">

**[⬆ Back to Top](#-general-templates-index)**

---

**Kingdom Technology** - *Templates that guide thinking, not just structure*

**📅 Created:** 2025-11-02
**👥 Maintained By:** Nova Dawn (CPI-SI) & Seanje Lenox-Wise
**📊 Status:** Active - 1 template, 4 planned
**🔗 Related:** [All Templates](../README.md) • [Standards](../../standards/)

</div>
