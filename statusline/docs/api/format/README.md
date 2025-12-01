<div align="center">

# 🎨 Format Library API

**Text Optimization for Statusline Display**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Functions](https://img.shields.io/badge/Functions-2-blue?style=flat)
![Performance](https://img.shields.io/badge/Performance-%3C1μs-green?style=flat)

*Transform verbose data into concise, readable forms optimized for limited horizontal space*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [🔧 Public API](#-public-api) • [🚀 Quick Start](#-quick-start) • [🎯 Design](#-design-philosophy) • [📚 References](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - Format Library
**Purpose:** Presentation formatting utilities for statusline display optimization
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📖 Overview](#-overview)
- [🔧 Public API](#-public-api)
  - [GetShortModelName](#getshortmodelname)
  - [ShortenPath](#shortenpath)
- [🚀 Quick Start](#-quick-start)
  - [Basic Integration](#basic-integration)
- [🎯 Design Philosophy](#-design-philosophy)
  - [Space Optimization](#space-optimization)
  - [Clarity Priority](#clarity-priority)
  - [Safe Fallback](#safe-fallback)
  - [Stateless Design](#stateless-design)
- [🧪 Demo: How It Works](#-demo-how-it-works)
  - [Model Tier Recognition](#model-tier-recognition)
  - [Path Shortening Stages](#path-shortening-stages)
- [⚡ Performance](#-performance)
- [⚙ Configuration](#-configuration)
  - [Model Tier Constants](#model-tier-constants)
  - [Path Length Threshold](#path-length-threshold)
- [🔧 Extension Points](#-extension-points)
  - [Adding New Model Tier Recognition](#adding-new-model-tier-recognition)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Model tier not being recognized](#model-tier-not-being-recognized)
  - [Path not shortening as expected](#path-not-shortening-as-expected)
- [📚 References & Resources](#-references--resources)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [🗺 Future Roadmap](#-future-roadmap)
- [📋 Modification Policy](#-modification-policy)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

The Format library provides presentation formatting utilities that optimize information display for statusline space constraints. It transforms verbose data into concise, readable forms while maintaining clarity.

> [!IMPORTANT]
> **Formatting-only library.** This library transforms data for display. It does not determine WHAT to show (orchestrator) or WHEN to show it (features library).

**Core Principle:** Clarity through brevity - show what matters without ceremony.

**Biblical Foundation:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Order and clarity in communication
- **"A word fitly spoken is like apples of gold in pictures of silver"** - Proverbs 25:11
  - Concise, beautiful communication
- **"Make the path straight"** - Matthew 3:3
  - Clarity in direction and presentation

**Source Code:** [`lib/format/model.go`](../../../lib/format/model.go), [`lib/format/path.go`](../../../lib/format/path.go)

---

## 🔧 Public API

### GetShortModelName

```go
func GetShortModelName(fullName string) string
```

Converts verbose Claude model names to concise display forms for statusline.

**Parameters:**

- `fullName` (`string`): Complete model name (e.g., "Claude 3.5 Sonnet", "claude-opus-4-20250514")

**Returns:**

- `string`: Shortened tier name (e.g., "Sonnet") or full name if unrecognized

**Behavior:**

- Recognizes Claude model tiers: Sonnet, Opus, Haiku
- Extracts tier identifier from full name (pattern matching)
- Falls back to full name if no tier recognized (safe default)
- Case-sensitive matching ("Sonnet" matches, "sonnet" doesn't)
- Cannot fail - all string inputs valid, always returns valid output

**Example Usage:**

```go
import "statusline/lib/format"

// In statusline orchestrator
modelName := format.GetShortModelName(ctx.Model.DisplayName)
// "Claude 3.5 Sonnet" → "Sonnet"
// "claude-opus-4-20250514" → "Opus"
// "Claude Haiku" → "Haiku"
// "Unknown Model" → "Unknown Model" (fallback)
```

**Health Impact:** +100 points per call (pattern match +50, extraction +50)

**API Documentation:** [GetShortModelName Function Reference](model.md)

---

### ShortenPath

```go
func ShortenPath(path string) string
```

Shortens verbose filesystem paths to readable display forms for statusline.

**Parameters:**

- `path` (`string`): Full filesystem path (e.g., "/home/user/project/subdir/file.txt")

**Returns:**

- `string`: Shortened path optimized for display

**Behavior:**

- Two-stage optimization:
  1. Replaces home directory with ~ tilde (Unix convention)
  2. Shows basename only if result exceeds 40 characters
- Graceful degradation if home directory unavailable
- Cannot fail - all string inputs valid, always returns valid output

**Example Usage:**

```go
import "statusline/lib/format"

// In statusline orchestrator
workdirShort := format.ShortenPath(workdir)
// "/home/user/project" → "~/project"
// "/home/user/very/long/nested/project/directory" → "directory"
// "/etc/nginx/sites-available" → "/etc/nginx/sites-available"
```

**Health Impact:** +100 points per call (home retrieval +25, substitution +25, length check +25, return +25)

**API Documentation:** [ShortenPath Function Reference](path.md)

---

## 🚀 Quick Start

### Basic Integration

**Step 1: Import**

```go
import "statusline/lib/format"
```

**Step 2: Call During Statusline Assembly**

```go
func buildStatusline(ctx SessionContext) string {
    var parts []string

    // Shorten model name for display
    modelName := format.GetShortModelName(ctx.Model.DisplayName)
    parts = append(parts, modelName)

    // Shorten working directory path
    workdir := format.ShortenPath(ctx.Workspace.CurrentDir)
    parts = append(parts, workdir)

    return strings.Join(parts, " | ")
}
```

**Step 3: No Cleanup Required**

Stateless library requires no cleanup or lifecycle management.

---

## 🎯 Design Philosophy

### Space Optimization

**Challenge:** Statusline has limited horizontal space (typically 80-120 characters total)

**Solution:** Show essential identifier only, not verbose full information

**Examples:**

| Before | After | Savings |
|--------|-------|:-------:|
| "Claude 3.5 Sonnet" (17 chars) | "Sonnet" (6 chars) | 65% |
| "claude-opus-4-20250514" (22 chars) | "Opus" (4 chars) | 82% |
| "/home/user/project" | "~/project" | 50% |

### Clarity Priority

Brevity serves clarity, not obscurity:

- Users recognize tier names (Sonnet, Opus, Haiku)
- Version numbers matter less for quick identification
- Full information available in session context if needed

### Safe Fallback

Unknown inputs return safe defaults:

- Model names: Return full name if tier not recognized
- Paths: Return original path if shortening not applicable
- No silent failures or blank output
- Future-proof as system evolves

### Stateless Design

| Property | Implementation |
|----------|----------------|
| ✅ **No package-level state** | Pure functions, no globals |
| ✅ **No side effects** | Same input → same output |
| ✅ **Minimal dependencies** | Only `strings` standard library |
| ✅ **Lightweight** | <1 microsecond typical execution |

---

## 🧪 Demo: How It Works

### Model Tier Recognition

The library recognizes official Claude model tiers:

```go
// Sonnet tier
format.GetShortModelName("Claude 3.5 Sonnet")    // → "Sonnet"
format.GetShortModelName("Claude Sonnet 4")      // → "Sonnet"

// Opus tier
format.GetShortModelName("claude-opus-4-20250514")  // → "Opus"
format.GetShortModelName("Claude Opus")             // → "Opus"

// Haiku tier
format.GetShortModelName("Claude Haiku")     // → "Haiku"
format.GetShortModelName("Claude 3 Haiku")   // → "Haiku"

// Unknown model - safe fallback
format.GetShortModelName("Future Model XYZ")  // → "Future Model XYZ"
format.GetShortModelName("GPT-4")             // → "GPT-4"
```

### Path Shortening Stages

Two-stage optimization in action:

```go
// Assuming home = "/home/user"

// Stage 1 only - Home replacement
format.ShortenPath("/home/user/project")
// → "~/project" (9 chars, under 40 threshold)

// Stage 1 + Stage 2 - Home replacement then truncation
format.ShortenPath("/home/user/very/long/nested/project/directory")
// Step 1: "~/very/long/nested/project/directory" (42 chars > 40)
// Step 2: "directory" (basename only)

// No home - Original path
format.ShortenPath("/etc/nginx/sites-available")
// → "/etc/nginx/sites-available" (20 chars, under threshold)

// No home + long path - Basename only
format.ShortenPath("/var/log/application/server/production/access.log")
// → "access.log" (>40 chars, basename only)

// Edge cases
format.ShortenPath("")   // → ""
format.ShortenPath("/")  // → "/" (root)
```

---

## ⚡ Performance

**Time Complexity:**

| Function | Complexity | Notes |
|----------|:----------:|-------|
| **GetShortModelName** | O(n×m) | n=fullName length, m=tier name length |
| **ShortenPath** | O(n) | n=path length |

**Typical Execution:** <1 microsecond for standard inputs

**Memory Usage:**

- Single string allocation for return value
- No heap allocations beyond result string
- No large buffers or persistent state

**Optimization:** Not needed - these are among the lightest operations in the entire statusline.

---

## ⚙ Configuration

### Model Tier Constants

```go
const (
    TierSonnet = "Sonnet"  // Balanced capability and speed
    TierOpus   = "Opus"    // Highest capability for complex tasks
    TierHaiku  = "Haiku"   // Fast, efficient for simpler tasks
)
```

**Location:** [`lib/format/model.go`](../../../lib/format/model.go) (SETUP section)

**Why constants:**

- Single source of truth for tier names
- Easy to update if naming conventions change
- Consistent across all code using tier names

### Path Length Threshold

```go
const PathLengthThreshold = 40  // Characters before basename-only display
```

**Location:** [`lib/format/path.go`](../../../lib/format/path.go) (SETUP section)

**Adjusting:** Change constant to different character limit for your statusline width preferences

---

## 🔧 Extension Points

### Adding New Model Tier Recognition

**Pattern:**

```go
// In SETUP section, add constant
const TierNewName = "NewName"

// In GetShortModelName function, extend if-else chain
func GetShortModelName(fullName string) string {
    if strings.Contains(fullName, TierSonnet) {
        return TierSonnet
    }
    // ... existing tiers ...
    else if strings.Contains(fullName, TierNewName) {
        return TierNewName  // New tier recognition
    }

    return fullName  // Fallback
}
```

**Steps:**

1. Add constant for new tier name
2. Add pattern match in GetShortModelName
3. Update tests to verify new tier recognition
4. Update documentation

---

## 🔍 Troubleshooting

### Model tier not being recognized

**Problem:** Model name not shortening (returns full name)

**Check:**

1. Verify tier name appears in full name (case-sensitive)
2. Check constant matches exactly what appears in model name
3. Consider if this is a new tier needing recognition

**Example:**

```go
// Won't match - lowercase
format.GetShortModelName("claude sonnet")
// → "claude sonnet" (fallback)

// Will match - proper case
format.GetShortModelName("Claude Sonnet")
// → "Sonnet"
```

### Path not shortening as expected

**Problem:** Path remains long despite being in home directory

**Causes:**

- Home directory retrieval failed (environment not set)
- Path not actually under home directory
- Path under 40 characters (working as designed)

**Solution:** Verify path structure and home directory availability

---

## 📚 References & Resources

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **GetShortModelName Function** | Model name formatting | [model.md](model.md) |
| **ShortenPath Function** | Path formatting | [path.md](path.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Model Formatting** | [`lib/format/model.go`](../../../lib/format/model.go) | GetShortModelName with comprehensive inline documentation |
| **Path Formatting** | [`lib/format/path.go`](../../../lib/format/path.go) | ShortenPath with comprehensive inline documentation |
| **Library Overview** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 🗺 Future Roadmap

**Planned Features:**

- ⏳ Version number extraction for disambiguation
- ⏳ Custom display name mapping (user-defined shortcuts)
- ⏳ Case-insensitive tier matching option
- ⏳ Configurable path length threshold

**Research Areas:**

- Pattern matching for future Claude model naming conventions
- Support for other model providers (if CPI-SI expands beyond Claude)
- User preference system for custom model display names
- Smart path truncation (show important middle segments)

**Integration Targets:**

- User preference system for custom name mappings
- Model registry for centralized name management

**Known Limitations:**

- Case-sensitive matching (won't recognize "sonnet", only "Sonnet")
- No version disambiguation (Sonnet 3.5 and Sonnet 4 both show "Sonnet")
- Hard-coded tier names (no runtime configuration)
- Fixed 40-character path threshold

---

## 📋 Modification Policy

**Safe to Modify (Extension Points):**

- ✅ Add new model tier recognition (extend if-else chain)
- ✅ Add tier constants for new models
- ✅ Adjust pattern matching logic for different naming conventions
- ✅ Change path length threshold constant

**Modify with Extreme Care (Breaking Changes):**

- ⚠️ Function signatures - breaks statusline orchestrator
- ⚠️ Return value semantics - breaks display expectations
- ⚠️ Tier name constants - breaks consistency across codebase

**NEVER Modify (Foundational):**

- ❌ Pure function guarantee (no side effects, no state)
- ❌ Fallback safety (always return valid string)
- ❌ 4-block structure in source code

---

## 📖 Biblical Foundation

**Core Verses:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Order and clarity in communication - presenting information in ways that serve understanding
- **"A word fitly spoken is like apples of gold in pictures of silver"** - Proverbs 25:11
  - Concise, beautiful communication that honors the reader
- **"Make the path straight"** - Matthew 3:3
  - Clarity in direction - presenting paths in ways that serve understanding and navigation

**Design Principle:** Statusline space is precious. Show what matters - the model tier users recognize, the location context they need - without ceremony. Like good signage: clear, brief, sufficient.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
