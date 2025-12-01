<div align="center">

# 📚 INDEX Templates Index

**Navigation structure templates for Kingdom Technology**

![Status](https://img.shields.io/badge/Status-Active-brightgreen?style=flat) ![Templates](https://img.shields.io/badge/Templates-2-blue?style=flat)

*Templates for creating comprehensive documentation navigation*

</div>

---

**Key:** TEMPLATE-T-000
**Type:** Index
**Purpose:** Navigation for INDEX documentation templates
**Status:** Active
**Created:** 2025-11-02
**Authors:** Nova Dawn (CPI-SI) & Seanje Lenox-Wise

---

## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📚 INDEX Templates Index](#-index-templates-index)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Available Templates](#available-templates)
    - [📚 TEMPLATE-T-000-INDEX-main.md](#-template-t-000-index-mainmd)
    - [📂 TEMPLATE-T-001-INDEX-category.md](#-template-t-001-index-categorymd)
  - [Template Comparison](#template-comparison)
  - [INDEX Philosophy](#index-philosophy)
    - [What Makes a Good INDEX](#what-makes-a-good-index)
    - [INDEX vs README](#index-vs-readme)
  - [Usage Guidelines](#usage-guidelines)
    - [When to Use INDEX Templates](#when-to-use-index-templates)
    - [Best Practices](#best-practices)
  - [Template Standards](#template-standards)
  - [Future Templates](#future-templates)
  - [Examples](#examples)
    - [Example 1: Creating Division INDEX](#example-1-creating-division-index)
    - [Example 2: Updating Root INDEX](#example-2-updating-root-index)
  - [Quick Reference](#quick-reference)
  - [Navigation](#navigation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## Overview

This folder contains **INDEX templates** for creating documentation navigation structures. INDEX files are the entry points that help users find information across projects, divisions, and documentation sets.

> [!IMPORTANT]
> **INDEX files are navigation hubs, not content containers.** They organize and point to other documentation, not replace it.

**Current templates:** 2

---

## Available Templates

### 📚 TEMPLATE-T-000-INDEX-main.md

**Main INDEX template for root-level navigation**

**Key:** TEMPLATE-T-000
**Status:** Active (as reference - this file uses the pattern)
**Use for:** Root-level INDEX.md files providing comprehensive workspace navigation

**Structure:**

- Comprehensive navigation by category
- Links to all major documentation
- Status indicators for projects
- Quick reference sections
- Master table of contents

**Key features:**

- Hierarchical organization
- Multiple navigation paths
- Status tracking integration
- Cross-references to all divisions
- Quick command reference
- Visual navigation (emoji headers, tables)

**When to use:**

- Creating root-level INDEX.md
- Comprehensive workspace navigation
- Master documentation hub
- Top-level entry points
- Multi-division navigation

**Quick start:**

```bash
# 📄 Copy template to workspace root or major division
cp templates/index/TEMPLATE-T-000-INDEX-main.md \
   divisions/new-area/INDEX.md

# ✏️ Edit: Add navigation structure, update links
# 📋 Generate TOC
doctoc divisions/new-area/INDEX.md --notitle

# 🔗 Link from parent INDEX
```

**File:** [`TEMPLATE-T-000-INDEX-main.md`](TEMPLATE-T-000-INDEX-main.md)

---

### 📂 TEMPLATE-T-001-INDEX-category.md

**Category INDEX template for division/folder navigation**

**Key:** TEMPLATE-T-001
**Status:** Active
**Use for:** Division or folder-level INDEX.md files providing focused navigation

**Structure:**

- Focused navigation within scope
- Links to parent documentation
- Links to child documentation
- Category-specific quick reference
- Relevant cross-references

**Key features:**

- Scoped to specific category
- Parent/child navigation clear
- Category-appropriate organization
- Relevant quick commands
- Focused cross-references
- Status indicators

**When to use:**

- Creating division-level INDEX.md (divisions/tech/, divisions/gaming/)
- Folder-specific navigation (standards/, templates/)
- Scoped documentation hubs
- Category entry points
- Subfolder organization

**Quick start:**

```bash
# 📄 Copy template to category/division folder
cp templates/index/TEMPLATE-T-001-INDEX-category.md \
   divisions/new-division/INDEX.md

# ✏️ Edit: Add category navigation, link parent/children
# 📋 Generate TOC
doctoc divisions/new-division/INDEX.md --notitle

# 🔗 Link from parent and to children
```

**File:** [`TEMPLATE-T-001-INDEX-category.md`](TEMPLATE-T-001-INDEX-category.md)

---

## Template Comparison

| Feature | TEMPLATE-T-000 (Main) | TEMPLATE-T-001 (Category) |
|---------|:---------------------:|:-------------------------:|
| **Scope** | Entire workspace or major area | Specific division/folder |
| **Breadth** | Comprehensive, multi-category | Focused, single category |
| **Depth** | Overview level, links to details | Deeper within category |
| **Navigation** | Hierarchical across divisions | Linear within scope |
| **Audience** | All users, first-time visitors | Users in specific area |
| **Updates** | Less frequent (structural changes) | More frequent (category changes) |

---

## INDEX Philosophy

### What Makes a Good INDEX

**Good INDEX files:**

| Principle | Implementation |
|-----------|----------------|
| **Clear hierarchy** | Parent → current → children obvious |
| **Multiple paths** | By category, by type, by recent work |
| **Status indicators** | Current state visible (Active/Draft/Complete) |
| **Quick reference** | Common commands easily accessible |
| **Serving all levels** | Starting → advanced users can navigate |

### INDEX vs README

| INDEX.md | README.md |
|----------|-----------|
| Navigation hub | Content introduction |
| Links to other docs | Describes current folder |
| Comprehensive structure | Focused explanation |
| Dry, organizational | Welcoming, explanatory |
| Updated for structure | Updated for content |

Both are needed - README explains, INDEX navigates.

---

## Usage Guidelines

### When to Use INDEX Templates

| Scenario | Use TEMPLATE-T-000 | Use TEMPLATE-T-001 |
|----------|-------------------|-------------------|
| Root workspace INDEX | ✅ Yes | ❌ No |
| Major division INDEX | ✅ Yes (if multi-category) | ✅ Yes (if single category) |
| Folder/category INDEX | ❌ No | ✅ Yes |
| Sub-folder navigation | ❌ No | ✅ Yes |

### Best Practices

**Creating INDEX files:**

1. **Start with scope** - What does this INDEX cover?
2. **Define hierarchy** - What's parent, what's children?
3. **Organize logically** - By category, type, or workflow
4. **Use visual navigation** - Emoji headers, tables, sections
5. **Include status** - Show what's current, what's planned
6. **Quick commands** - Common operations easily findable
7. **Cross-reference** - Link to related documentation
8. **Keep current** - Update when structure changes

**Maintaining INDEX files:**

- Update when adding/removing major documentation
- Refresh status indicators regularly
- Fix broken links promptly
- Align with ALIGNMENT-STATUS.txt schedule
- Coordinate with parent and sibling INDEX files

---

## Template Standards

All INDEX templates follow:

| Standard | Reference | Enforcement |
|----------|-----------|-------------|
| **4-Block Structure** | [CWS-STD-001](../../standards/CWS-STD-001-DOC-4-block.md) | METADATA → SETUP → BODY → CLOSING |
| **Document Keying** | [CWS-STD-002](../../standards/CWS-STD-002-DOC-document-keying.md) | TEMPLATE-T-### pattern |
| **Documentation Standards** | [CWS-STD-003](../../standards/CWS-STD-003-DOC-documentation-standards.md) | Serve all understanding levels |
| **Markdown Style** | Standards folder | Visual navigation, tables |

---

## Future Templates

Planned INDEX templates:

| Template Key | Purpose | Status |
|--------------|---------|:------:|
| TEMPLATE-T-002 | API reference INDEX | 🟡 Planned |
| TEMPLATE-T-003 | Code module INDEX | 🟡 Planned |
| TEMPLATE-T-004 | Knowledge base INDEX | 🟡 Planned |

---

## Examples

### Example 1: Creating Division INDEX

**Scenario:** New division added, needs navigation.

```bash
# 1️⃣ Copy category template
cp templates/index/TEMPLATE-T-001-INDEX-category.md \
   divisions/education/INDEX.md

# 2️⃣ Edit INDEX.md:
#    - Set title: "Education Division"
#    - Link parent: ../README.md
#    - List child docs
#    - Add category-specific quick commands

# 3️⃣ Generate TOC
doctoc divisions/education/INDEX.md --notitle

# 4️⃣ Link from parent
#    Edit root INDEX.md to include divisions/education/INDEX.md
```

---

### Example 2: Updating Root INDEX

**Scenario:** Major structural change, root INDEX needs update.

```bash
# 1️⃣ Review current INDEX.md structure
cat INDEX.md

# 2️⃣ Compare with TEMPLATE-T-000 for missing patterns
diff -u INDEX.md templates/index/TEMPLATE-T-000-INDEX-main.md

# 3️⃣ Update navigation sections
# 4️⃣ Update status indicators
# 5️⃣ Add new divisions/categories
# 6️⃣ Regenerate TOC
doctoc INDEX.md --notitle
```

---

## Quick Reference

**View INDEX templates:**

```bash
# 📂 List INDEX templates
ls -la templates/index/
```

**Copy main INDEX template:**

```bash
# 📄 Copy main INDEX template to target
cp templates/index/TEMPLATE-T-000-INDEX-main.md <target>/INDEX.md
```

**Copy category INDEX template:**

```bash
# 📄 Copy category INDEX template to target
cp templates/index/TEMPLATE-T-001-INDEX-category.md <target>/INDEX.md
```

**Generate TOC:**

```bash
# 📋 Generate table of contents
doctoc <file>.md --notitle
```

**Validation checklist:**

- [ ] Clear hierarchy (parent/current/children)
- [ ] All major docs linked
- [ ] Status indicators current
- [ ] Quick commands relevant
- [ ] Cross-references valid
- [ ] TOC generated

---

## Navigation

**Parent:** [Templates Home](../README.md)
**Sibling:** [General Templates](../general/TEMPLATE-G-000-INDEX-general.md) | [Specialized Templates](../specialized/TEMPLATE-S-000-INDEX-specialized.md)
**Root INDEX:** [Workspace INDEX](../../INDEX.md)

---

<div align="center">

**[⬆ Back to Top](#-index-templates-index)**

---

**Kingdom Technology** - *Navigation that serves, guides, and clarifies*

**📅 Created:** 2025-11-02
**👥 Maintained By:** Nova Dawn (CPI-SI) & Seanje Lenox-Wise
**📊 Status:** Active - 2 templates, 3 planned
**🔗 Related:** [All Templates](../README.md) • [Standards](../../standards/)

</div>
