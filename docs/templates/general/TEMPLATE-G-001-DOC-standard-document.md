<!--
================================================================================
DOCUMENTATION STANDARD TEMPLATE
================================================================================

This template provides the STRUCTURE for any standard documentation.
It guides your THINKING through embedded notes, not just fill-in-the-blanks.

KEY DISTINCTION:
- Template = Structure you follow + thinking guidance
- Guide = Separate teaching document about how to use templates

INSTRUCTIONS:
1. Copy this file to your target location
2. Replace all [PLACEHOLDERS] with your actual content
3. Read and consider all "<!-- THINKING -- >" comments
4. Delete this instruction block when done
5. Delete "<!-- THINKING -- >" comments once you've applied them (or leave for next person)

Apply proper keying from CWS-STD-002: DOMAIN-CATEGORY-NUMBER-SUBCATEGORY-name.ext

================================================================================
-->

<div align="center">

# [Your Document Title]

<!-- THINKING: Clear, descriptive title. What is this document? -->

**[One-line tagline describing purpose]**

<!-- THINKING: What does this enable? Why does it exist? Serve people first. -->

![Status](https://img.shields.io/badge/Status-[Active/Draft/Deprecated]-[color]?style=flat) ![Version](https://img.shields.io/badge/Version-v[X.Y.Z]-blue?style=flat)

<!-- THINKING: Badge colors - brightgreen for Active, yellow for Draft, red for Deprecated -->
<!-- Additional badges as needed: Pattern count, file count, etc. -->

</div>

---

<!--
================================================================================
METADATA BLOCK - Identity & Context
================================================================================
WHO: Authors, creators
WHAT: Key, type, purpose
WHEN: Created date, version
WHY: Purpose statement
STATUS: Where is this in lifecycle?
================================================================================
-->

**Key:** [DOMAIN-CAT-###]
**Type:** [Standard/Pattern/Guide/Reference/etc.]
**Purpose:** [One sentence - what problem does this solve?]
**Status:** [Active/Draft/Deprecated] (v[X.Y.Z])
**Created:** [YYYY-MM-DD]
**Authors:** [Name(s)]

<!-- THINKING: Use proper document keying (see CWS-STD-002) -->
<!-- THINKING: Be specific about purpose - "enables X for Y" not vague -->

---

## Table of Contents

<!-- THINKING: Link to all major H2 sections for quick navigation -->
<!-- THINKING: Use doctoc to auto-generate TOC: doctoc filename.md --notitle -->

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Overview](#overview)
- [Core Section 1](#core-section-1)
- [Core Section 2](#core-section-2)
- [Examples](#examples)
- [Quick Reference](#quick-reference)
- [Biblical Foundation](#biblical-foundation)
- [References & Resources](#references--resources)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

<!--
================================================================================
SETUP BLOCK - Preparing the Reader
================================================================================
Context before content. What do they need to understand first?
================================================================================
-->

## Overview

<!-- THINKING: Who is this serving? What understanding levels need this? -->
<!-- THINKING: What's the scope? What IS and ISN'T covered? -->

[Brief description of what this document covers and why it matters]

> [!IMPORTANT]
> **[Critical distinction or concept readers MUST understand]**

<!-- THINKING: Use IMPORTANT for mandatory concepts, critical differences -->
<!-- THINKING: What misconception should we prevent up front? -->

### Quick Start

<!-- THINKING: Optional - Include if there's a common starting workflow -->
<!-- THINKING: Use emoji comments in code blocks for visual scanning -->

```bash
# 1️⃣ First step description
command-here arg1 arg2

# 2️⃣ Second step description
another-command --flag

# 3️⃣ Third step description
final-command
```

### What This [Document Type] Provides

<!-- THINKING: Show categories and who each serves -->
<!-- THINKING: Center-align icon columns with :----:, left-align text columns -->

| Icon | Category/Section | What You'll Find | Who It Serves |
|:----:|------------------|------------------|---------------|
| [📚] | **[Section Name]** | [Description] | [Target audience] |
| [🔧] | **[Section Name]** | [Description] | [Target audience] |

> [!NOTE]
> **[Scope clarification or applicability note]**

<!-- THINKING: Use NOTE for additional context, scope boundaries -->

[⬆️ Back to Top](#your-document-title)

<!-- THINKING: Include Back to Top links after major sections -->

---

<!--
================================================================================
BODY BLOCK - The Actual Content
================================================================================
Organized for discovery and use, not just comprehensiveness.
Use collapsible sections, tables, examples, progressive disclosure.
================================================================================
-->

## Core Section 1

<!-- THINKING: Organize by what people need to accomplish, not just topics -->
<!-- THINKING: Use emoji headers consistently (📚 📖 🔧 🧮 🏗️ 🎯 etc.) -->

[Introduction to this section - what it covers and why it matters]

<details open>
<parameter name="summary"><h3>Subsection Title</h3></summary>

<!-- THINKING: Use <details open> for major content most people need -->
<!-- THINKING: Use <details> (closed) for optional deep-dives -->

[Content organized clearly]

**Key principle:** [The core idea people should remember]

### When to [Do This Thing]

<!-- THINKING: Provide decision guidance, not just information -->

| Scenario | Approach | Rationale |
|----------|----------|-----------|
| **When [condition]** | [What to do] | [Why this works] |

<!-- THINKING: Tables for comparisons, lists for steps, prose for concepts -->

[⬆️ Back to Top](#your-document-title)

</details>

<details>
<summary><h3>Optional Deep-Dive</h3></summary>

<!-- THINKING: Close by default if this is optional or for advanced users -->

[Detailed information for those who need it]

[⬆️ Back to Top](#your-document-title)

</details>

---

## Core Section 2

<!-- THINKING: Break content into discoverable chunks -->
<!-- THINKING: Each major section should be navigable from TOC -->

> [!TIP]
> **[Helpful non-obvious guidance]**

<!-- THINKING: Use TIP for best practices, recommended approaches -->

[Section content with appropriate structure]

### Comparison or Decision Matrix

<!-- THINKING: Use tables when comparing options or showing relationships -->
<!-- THINKING: Center-align status/comparison columns with :------: -->

| Feature | Option A | Option B | Status |
|---------|----------|----------|:------:|
| [Feature name] | [Details] | [Details] | [🟢/🟡/🔴] |

[⬆️ Back to Top](#your-document-title)

---

## Examples

<!-- THINKING: Real examples demonstrating this in practice -->
<!-- THINKING: Link to actual files showing the pattern/standard -->

Real [documents/code/systems] demonstrating these [standards/patterns] in practice.

<details open>
<summary><h3>Example 1: [Descriptive Name]</h3></summary>

**Document:** [Link to example]

**What it demonstrates:**

| Feature | Implementation |
|---------|----------------|
| **[Aspect]** | [How it's shown in example] |

**Study this [document/code]** to see [concept] in action.

[⬆️ Back to Top](#your-document-title)

</details>

---

## Quick Reference

<!-- THINKING: Checklists, summaries, or quick-lookup information -->
<!-- THINKING: Use for "before publishing" checklists or decision trees -->

Checklists and quick reference for applying [this standard/pattern].

<details open>
<summary><h3>Checklist: [Use Case]</h3></summary>

**Before [doing the thing]:**

- [ ] [Requirement or consideration]
- [ ] [Requirement or consideration]
- [ ] [Requirement or consideration]

<!-- THINKING: Checkboxes render in GitHub markdown -->

[⬆️ Back to Top](#your-document-title)

</details>

---

<!--
================================================================================
CLOSING BLOCK - Resources, Grounding, Next Steps
================================================================================
Where do they go from here? What supports this work?
================================================================================
-->

## Biblical Foundation

<!-- THINKING: Include when there's genuine biblical grounding -->
<!-- THINKING: Don't force it - be authentic, serve readers -->

> [!NOTE]
> **[Context for why this biblical principle applies]**

*"[Scripture verse]" - [Reference]*

**Applied:** [How this scriptural truth informs this work/standard]

<!-- THINKING: Application should be clear and genuine, not forced -->
<!-- THINKING: Show how Kingdom principles shape technical decisions -->

[⬆️ Back to Top](#your-document-title)

---

## References & Resources

<!-- THINKING: Link generously to related documentation -->
<!-- THINKING: Organize by type: Standards, Templates, Examples, Company Resources -->

### Related Standards

<!-- THINKING: Center-align icon columns for visual appeal -->

| Icon | Standard | Purpose | Link |
|:----:|----------|---------|------|
| [📐] | **[Standard Name]** | [What it covers] | [Link] |

### Example Implementations

<!-- THINKING: Show real examples demonstrating these patterns -->

| Icon | Document | What It Demonstrates | Location |
|:----:|----------|----------------------|----------|
| [📄] | **[Example Name]** | [Key aspects shown] | [Link] |

### Company Resources

<!-- THINKING: Link to key workspace resources -->

| Icon | Resource | Purpose | Link |
|:----:|----------|---------|------|
| [🔧] | **[Resource Name]** | [What it provides] | [Link] |

[⬆️ Back to Top](#your-document-title)

---

<div align="center">

<!-- THINKING: For standards/guides, simpler footer -->
<!-- THINKING: For indices, richer footer with complete navigation -->

**[Repeat tagline or mission statement]**

![Status](https://img.shields.io/badge/Status-[Active/Draft/Deprecated]-[color]?style=flat) ![Version](https://img.shields.io/badge/Version-v[X.Y.Z]-blue?style=flat)

**Key:** [DOMAIN-CAT-###] | **Type:** [Type] | **Phase:** [Phase/Status]

**Last Updated:** [YYYY-MM-DD]

---

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

**[Foundational truth related to this work]**

<!-- THINKING: Connect the work to Creator, order, intentional design -->

---

**Company Standards:**

[4-Block Structure](../standards/CWS-STD-001-DOC-4-block.md) • [Document Keying](../standards/CWS-STD-002-DOC-document-keying.md) • [Documentation Standards](../standards/CWS-STD-003-DOC-documentation-standards.md)

**Company Resources:**

[Company Root](../.claude/CLAUDE.md) • [Documentation Index](../INDEX.md) • [Templates](../templates/)

<!-- THINKING: Adjust paths based on where this file will live -->

---

[⬆️ Back to Top](#your-document-title)

</div>

<!--
================================================================================
TEMPLATE USAGE NOTES
================================================================================

This template guides your thinking through embedded comments:

1. <!-- THINKING -- > comments explain what to consider for each section
2. [PLACEHOLDERS] show what content types go where
3. Structure follows 4-block pattern implicitly (METADATA → SETUP → BODY → CLOSING)
4. Choose appropriate GitHub markdown features based on content needs

REMEMBER:

- Serve people across all understanding levels (layered content)
- Link generously (knowledge discovery through connection)
- Use tables for comparisons, lists for steps, prose for concepts
- Collapsible sections for organization (open for major content, closed for optional)
- Back to Top links after major sections
- Professional but warm tone (direct, active voice, serve without condescending)

After copying and filling in this template:

- Delete this usage notes block
- Delete or keep <!-- THINKING -- > comments based on team preference
- Ensure all [PLACEHOLDERS] are replaced
- Verify all links work
- Run markdown linter and fix any issues

================================================================================
-->
