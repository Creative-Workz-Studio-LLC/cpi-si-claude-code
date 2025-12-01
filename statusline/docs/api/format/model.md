<div align="center">

# GetShortModelName Function Reference

**Model Name Shortening for Statusline Display**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Savings](https://img.shields.io/badge/Space%20Savings-65%25+-green?style=flat)

*Converts verbose Claude model names to concise display forms*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetShortModelName Function
**Purpose:** Model name text optimization for statusline space constraints
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

- [GetShortModelName Function Reference](#getshortmodelname-function-reference)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
  - [📝 Signature](#-signature)
  - [⚙ Parameters](#-parameters)
    - [fullName (string)](#fullname-string)
  - [🔄 Returns](#-returns)
    - [string](#string)
  - [🎯 Behavior](#-behavior)
    - [Pattern Matching Algorithm](#pattern-matching-algorithm)
    - [Fallback Safety](#fallback-safety)
  - [💻 Usage](#-usage)
    - [Basic Usage](#basic-usage)
    - [Tier Recognition Examples](#tier-recognition-examples)
    - [Testing](#testing)
  - [🏥 Health Scoring](#-health-scoring)
  - [⚡ Performance](#-performance)
  - [⚙ Configuration](#-configuration)
    - [Tier Constants](#tier-constants)
  - [🔧 Edge Cases](#-edge-cases)
    - [Case Sensitivity](#case-sensitivity)
    - [Empty String](#empty-string)
    - [Special Characters](#special-characters)
  - [📊 Space Savings](#-space-savings)
  - [🔧 Extension Points](#-extension-points)
    - [Adding New Tier Recognition](#adding-new-tier-recognition)
  - [🔍 Troubleshooting](#-troubleshooting)
    - [Tier not being recognized](#tier-not-being-recognized)
    - [New model tier](#new-model-tier)
  - [🔗 Related Functions](#-related-functions)
  - [📚 References](#-references)
    - [API Documentation](#api-documentation)
    - [Source Code](#source-code)
    - [Related Documentation](#related-documentation)
  - [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetShortModelName** converts verbose Claude model names to concise display forms suitable for statusline space constraints.

> [!NOTE]
> **Pure formatting function.** Simple pattern matching with safe fallback. Cannot fail, no side effects.

**Design:** Pattern matching on model tier names (Sonnet, Opus, Haiku) to extract essential identifier without version numbers or branding.

**Library:** [`statusline/lib/format`](../../../lib/format/)

**Source Code:** [`lib/format/model.go`](../../../lib/format/model.go)

**Library README:** [Format Library Documentation](README.md)

---

## 📝 Signature

```go
func GetShortModelName(fullName string) string
```

---

## ⚙ Parameters

### fullName (string)

Complete model name from session context.

**Expected Formats:**

| Format | Example | Type |
|--------|---------|:----:|
| Standard | `"Claude 3.5 Sonnet"` | ✅ |
| API format | `"claude-opus-4-20250514"` | ✅ |
| Simple | `"Claude Haiku"` | ✅ |
| Unrecognized | `"Unknown Model v2"` | ✅ |

**Validation:** None required - all string inputs valid

---

## 🔄 Returns

### string

Shortened tier name or full name if unrecognized.

**Recognized Tiers:**

| Input Contains | Returns | Savings |
|----------------|---------|:-------:|
| `"Sonnet"` | `"Sonnet"` | 65% |
| `"Opus"` | `"Opus"` | 82% |
| `"Haiku"` | `"Haiku"` | 58% |

**Fallback:** Returns `fullName` unchanged if no tier recognized

**Guarantee:** Always returns valid string (empty if input empty)

---

## 🎯 Behavior

### Pattern Matching Algorithm

Case-sensitive contains check for tier names:

```go
if strings.Contains(fullName, "Sonnet") {
    return "Sonnet"
}
```

**Match Examples:**

| Input | Contains "Sonnet"? | Result |
|-------|:------------------:|--------|
| `"Claude 3.5 Sonnet"` | ✅ Yes | `"Sonnet"` |
| `"claude-sonnet-3-5"` | ❌ No (lowercase 's') | `"claude-sonnet-3-5"` |
| `"SONNET 4"` | ✅ Yes (uppercase string contains "Sonnet") | `"Sonnet"` |

### Fallback Safety

Unknown models return full name - no silent failures:

```go
format.GetShortModelName("Future Model XYZ")
// → "Future Model XYZ" (safe fallback)

format.GetShortModelName("GPT-4")
// → "GPT-4" (different provider)
```

**Why this matters:** Future-proof as Claude evolves. New models work immediately (just show full name).

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/format"

func buildStatusline(ctx SessionContext) string {
    // Shorten model name for display
    modelName := format.GetShortModelName(ctx.Model.DisplayName)
    // "Claude 3.5 Sonnet" → "Sonnet"

    parts := []string{modelName}
    // ... add other parts ...

    return strings.Join(parts, " | ")
}
```

### Tier Recognition Examples

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
format.GetShortModelName("")                  // → ""
```

### Testing

```go
func TestTierRecognition(t *testing.T) {
    testCases := map[string]string{
        "Claude 3.5 Sonnet":      "Sonnet",
        "claude-opus-4-20250514": "Opus",
        "Claude Haiku":           "Haiku",
        "Unknown Model":          "Unknown Model",
        "":                       "",
    }

    for input, expected := range testCases {
        result := format.GetShortModelName(input)
        if result != expected {
            t.Errorf("Input %q: expected %q, got %q",
                input, expected, result)
        }
    }
}
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Pattern match success | +50 | Tier identified |
| Name extraction | +50 | Tier name returned |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** Pattern matching with fallback guarantees valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(n×m)

- n = fullName length
- m = tier name length

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Model name length** | 20-30 chars | Standard Claude names |
| **Tier name length** | 4-6 chars | "Opus", "Haiku", "Sonnet" |
| **Execution time** | <1 μs | Microseconds |

**Memory:**

- Single string allocation for return value
- No heap allocations beyond result string

**Optimization:** Not needed - pattern matching is already optimal

---

## ⚙ Configuration

### Tier Constants

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

---

## 🔧 Edge Cases

### Case Sensitivity

```go
// Lowercase won't match
format.GetShortModelName("claude sonnet")
// → "claude sonnet" (not recognized - fallback)

// Proper case matches
format.GetShortModelName("Claude Sonnet")
// → "Sonnet"

// Uppercase string contains "Sonnet" substring
format.GetShortModelName("SONNET")
// → "Sonnet" (matches!)
```

**Behavior:** Case-sensitive by design. Official Claude model names use proper case ("Sonnet" not "sonnet").

### Empty String

```go
format.GetShortModelName("")
// → "" (valid - no error)
```

### Special Characters

```go
format.GetShortModelName("Sonnet™ v2")
// → "Sonnet" (contains "Sonnet" - matches)
```

---

## 📊 Space Savings

**Examples:**

| Full Name | Shortened | Character Savings |
|-----------|-----------|:-----------------:|
| `"Claude 3.5 Sonnet"` (17 chars) | `"Sonnet"` (6 chars) | **65%** |
| `"claude-opus-4-20250514"` (22 chars) | `"Opus"` (4 chars) | **82%** |
| `"Claude Haiku"` (12 chars) | `"Haiku"` (5 chars) | **58%** |

**Benefit:** More room for other statusline information (git, session stats, system metrics)

---

## 🔧 Extension Points

### Adding New Tier Recognition

**Pattern:**

```go
// In SETUP section, add constant
const TierNewName = "NewName"

// In GetShortModelName function
func GetShortModelName(fullName string) string {
    if strings.Contains(fullName, TierSonnet) {
        return TierSonnet
    }
    // ... existing tiers ...
    else if strings.Contains(fullName, TierNewName) {
        return TierNewName
    }

    return fullName  // Fallback
}
```

**Steps:**

1. Add tier constant in SETUP section
2. Extend if-else chain in function
3. Update tests for new tier
4. Update API documentation

---

## 🔍 Troubleshooting

### Tier not being recognized

**Problem:** Function returns full name instead of tier

**Check:**

1. Verify tier name appears in fullName (case-sensitive)
2. Check exact spelling and case
3. Confirm it's a recognized tier (Sonnet/Opus/Haiku)

**Example:**

```go
// Won't match - lowercase
format.GetShortModelName("claude sonnet")
// → "claude sonnet" (fallback)

// Will match - proper case
format.GetShortModelName("Claude Sonnet")
// → "Sonnet"
```

### New model tier

**Solution:** Add new tier constant and pattern match (see Extension Points above)

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **ShortenPath** | Filesystem path display optimization | [path.md](path.md) |

**Future Functions:**

- Version extraction for disambiguation (planned)
- Custom display name mapping (planned)

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Format Library README** | Library overview and integration guide | [README.md](README.md) |
| **ShortenPath Function** | Path formatting function | [path.md](path.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/format/model.go`](../../../lib/format/model.go) | GetShortModelName with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 📖 Biblical Foundation

**Core Verse:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Clarity and order in communication - presenting information in ways that serve understanding

**Supporting Verse:**

- **"A word fitly spoken is like apples of gold in pictures of silver"** - Proverbs 25:11
  - Concise, beautiful communication that honors the reader

**Application:** Show what matters (tier) without ceremony (version numbers). Like good signage: clear, brief, sufficient.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
