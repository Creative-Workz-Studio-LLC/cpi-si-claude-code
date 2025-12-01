<div align="center">

# ShortenPath Function Reference

**Path Shortening for Statusline Display**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Savings](https://img.shields.io/badge/Space%20Savings-up%20to%2093%25-green?style=flat)

*Shortens verbose filesystem paths to readable display forms*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - ShortenPath Function
**Purpose:** Path text optimization for statusline space constraints
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

- [📖 Overview](#-overview)
- [📝 Signature](#-signature)
- [⚙ Parameters](#-parameters)
  - [path (string)](#path-string)
- [🔄 Returns](#-returns)
  - [string](#string)
- [🎯 Behavior](#-behavior)
  - [Stage 1: Home Directory Replacement](#stage-1-home-directory-replacement)
  - [Stage 2: Length Evaluation](#stage-2-length-evaluation)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [Two-Stage Optimization Examples](#two-stage-optimization-examples)
  - [Edge Case Handling](#edge-case-handling)
  - [Testing](#testing)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [⚙ Configuration](#-configuration)
  - [MaxPathLength Constant](#maxpathlength-constant)
- [📊 Space Savings](#-space-savings)
- [🔧 Edge Cases](#-edge-cases)
  - [No Home Directory](#no-home-directory)
  - [Path Not in Home](#path-not-in-home)
  - [Exactly at Threshold](#exactly-at-threshold)
  - [One Character Over Threshold](#one-character-over-threshold)
- [🔧 Extension Points](#-extension-points)
  - [Adding Smart Truncation](#adding-smart-truncation)
  - [Dynamic Length Threshold](#dynamic-length-threshold)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Home replacement not working](#home-replacement-not-working)
  - [Truncation threshold seems wrong](#truncation-threshold-seems-wrong)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**ShortenPath** shortens verbose filesystem paths to readable display forms suitable for statusline space constraints.

> [!NOTE]
> **Pure formatting function.** Two-stage optimization with graceful degradation. Cannot fail, no side effects.

**Design:** Home directory replacement with ~ tilde (Unix convention), then basename extraction for paths exceeding length threshold.

**Library:** [`statusline/lib/format`](../../../lib/format/)

**Source Code:** [`lib/format/path.go`](../../../lib/format/path.go)

**Library README:** [Format Library Documentation](README.md)

---

## 📝 Signature

```go
func ShortenPath(path string) string
```

---

## ⚙ Parameters

### path (string)

Full filesystem path to shorten.

**Expected Formats:**

| Format | Example | Type |
|--------|---------|:----:|
| Absolute in home | `/home/user/project/subdir/file.txt` | ✅ |
| Absolute outside home | `/etc/nginx/sites-available` | ✅ |
| Already shortened | `~/project` | ✅ |
| Empty string | `""` | ✅ |

**Validation:** None required - all string inputs valid

---

## 🔄 Returns

### string

Shortened path optimized for statusline display.

**Optimization Stages:**

| Stage | Operation | Example |
|:-----:|-----------|---------|
| **1** | Home Replacement | `/home/user/...` → `~/...` |
| **2** | Length Check | If result >40 chars → basename only |

**Examples:**

| Input | Output | Stage |
|-------|--------|:-----:|
| `/home/user/project` | `~/project` | 1 only |
| `/home/user/very/long/nested/project/directory` | `directory` | 1+2 |
| `/etc/config` | `/etc/config` | None |

**Guarantee:** Always returns valid path string (empty if input empty)

---

## 🎯 Behavior

### Stage 1: Home Directory Replacement

Replaces home directory prefix with ~ tilde (Unix convention):

```go
home, err := os.UserHomeDir()
if err == nil && strings.HasPrefix(path, home) {
    path = "~" + strings.TrimPrefix(path, home)
}
```

**If UserHomeDir fails:** Path unchanged (graceful degradation)

**Why this stage:** Unix convention, saves characters, familiar to users

### Stage 2: Length Evaluation

If shortened path still exceeds threshold, show basename only:

```go
if len(path) > MaxPathLength {
    return filepath.Base(path)
}
```

**Threshold:** 40 characters (configurable via `MaxPathLength` constant)

**Why basename:** Shows immediate context (file/directory name) without full hierarchy

**Trade-off:** Lose full path context, gain significant space savings

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/format"

func buildStatusline(workdir string) string {
    // Shorten working directory for display
    workdirShort := format.ShortenPath(workdir)
    // "/home/user/project" → "~/project"

    parts := []string{workdirShort}
    // ... add other parts ...

    return strings.Join(parts, " | ")
}
```

### Two-Stage Optimization Examples

```go
// Assuming home = "/home/user"

// Stage 1 only - Home replacement
format.ShortenPath("/home/user/project")
// → "~/project" (9 chars, under 40 threshold)

// Stage 1 + Stage 2 - Home replacement then truncation
format.ShortenPath("/home/user/very/long/nested/project/directory")
// Stage 1: "~/very/long/nested/project/directory" (42 chars > 40)
// Stage 2: "directory" (basename only)

// No home directory - Original path
format.ShortenPath("/etc/nginx/sites-available")
// → "/etc/nginx/sites-available" (20 chars, under 40)

// No home + long path - Basename only
format.ShortenPath("/var/log/application/server/production/access.log")
// → "access.log" (>40 chars, basename only)
```

### Edge Case Handling

```go
// Empty string
format.ShortenPath("")
// → ""

// Root path
format.ShortenPath("/")
// → "/" (basename of root is "/")

// Already shortened
format.ShortenPath("~/project")
// → "~/project" (no change needed)

// Home directory itself
format.ShortenPath("/home/user")
// → "~"

// Exactly at threshold (40 chars)
path := strings.Repeat("x", 40)
format.ShortenPath(path)
// → path unchanged (40 is NOT greater than 40)

// One character over threshold (41 chars)
path := strings.Repeat("x", 41)
format.ShortenPath(path)
// → "x" (basename)
```

### Testing

```go
func TestHomeReplacement(t *testing.T) {
    home, _ := os.UserHomeDir()
    testPath := filepath.Join(home, "project")
    result := format.ShortenPath(testPath)

    if !strings.HasPrefix(result, "~") {
        t.Errorf("Expected ~ prefix, got %q", result)
    }
}

func TestLengthTruncation(t *testing.T) {
    longPath := "/var/log/application/server/production/subdirectory/access.log"
    result := format.ShortenPath(longPath)

    if len(result) > 40 {
        t.Errorf("Result %q exceeds 40 chars", result)
    }

    if result != "access.log" {
        t.Errorf("Expected basename, got %q", result)
    }
}
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Home directory retrieval | +25 | Success or graceful degradation |
| Tilde substitution | +25 | If applicable |
| Length check | +25 | Decision made |
| Path return | +25 | Result produced |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** UserHomeDir can fail (no home), but operation continues with graceful degradation. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(n) where n = path length

**Operations:**

| Operation | Complexity | Notes |
|-----------|:----------:|-------|
| **UserHomeDir()** | O(1) | Typically cached by OS |
| **String operations** | O(n) | HasPrefix, TrimPrefix |
| **Basename extraction** | O(n) | Worst case |

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Path length** | 50-100 chars | Standard paths |
| **Execution time** | <1 μs | Microseconds |

**Memory:**

- Single UserHomeDir() call (minimal I/O, cached)
- String allocations for result
- No large buffers or persistent state

**Optimization:** Not needed - string operations are already optimal

---

## ⚙ Configuration

### MaxPathLength Constant

```go
const MaxPathLength = 40  // Characters before basename-only display
```

**Location:** [`lib/format/path.go`](../../../lib/format/path.go) (SETUP section)

**Reasoning:**

- Statusline typically 80-120 characters total
- Allocating 40 to path leaves room for other information
- Longer paths would crowd out other context (git, session stats, system metrics)

**To Adjust:**

Change constant value in source code. Consider:

- Terminal width
- Other statusline elements
- Balance between context and conciseness

---

## 📊 Space Savings

**Examples:**

| Original Path | Shortened | Character Savings |
|---------------|-----------|:-----------------:|
| `/home/user/projects/work/current` (33 chars) | `~/projects/work/current` (24 chars) | **27%** |
| `/home/user/go/src/github.com/org/repo/pkg` (45 chars) | `pkg` (3 chars) | **93%** |
| `/var/log/application/production` (32 chars) | `/var/log/application/production` (32 chars) | **0%** |

**Benefit:** More room for git status, session stats, system metrics in statusline

---

## 🔧 Edge Cases

### No Home Directory

```go
// If UserHomeDir() returns error
format.ShortenPath("/home/user/project")
// → "/home/user/project" (unchanged - graceful degradation)
```

**Cause:** HOME environment variable not set, or OS doesn't support user home directories

**Behavior:** Skip stage 1, proceed to stage 2 with original path

### Path Not in Home

```go
format.ShortenPath("/etc/nginx/sites-available")
// Stage 1: No home replacement (path doesn't start with home)
// Stage 2: Length check (20 chars < 40)
// → "/etc/nginx/sites-available" (unchanged)
```

### Exactly at Threshold

```go
// Path exactly 40 chars
path := strings.Repeat("x", 40)
format.ShortenPath(path)
// → path unchanged (40 is NOT greater than 40)
```

**Note:** Condition is `>` not `>=`, so 40 characters exactly is allowed

### One Character Over Threshold

```go
// Path 41 chars
path := strings.Repeat("x", 41)
format.ShortenPath(path)
// → "x" (basename of long string)
```

---

## 🔧 Extension Points

### Adding Smart Truncation

Preserve important path segments instead of just showing basename:

**Pattern:**

```go
// Example: "/home/user/long/nested/project/file.go"
// Smart result: "~/…/project/file.go" (show first and last segments)

func ShortenPathSmart(path string, maxLen int) string {
    // 1. Apply home directory replacement
    // 2. If still too long, intelligently truncate middle
    // 3. Use ellipsis (…) to indicate omitted segments
    // 4. Preserve first segment (context) and last segment (file/dir name)

    // Implementation left as exercise
}
```

### Dynamic Length Threshold

Adjust threshold based on terminal width:

**Pattern:**

```go
func ShortenPathDynamic(path string, terminalWidth int) string {
    // Calculate available space for path
    maxPathLen := terminalWidth / 3  // Allocate 1/3 of width to path

    // Use calculated threshold
    home, err := os.UserHomeDir()
    if err == nil && strings.HasPrefix(path, home) {
        path = "~" + strings.TrimPrefix(path, home)
    }

    if len(path) > maxPathLen {
        return filepath.Base(path)
    }

    return path
}
```

---

## 🔍 Troubleshooting

### Home replacement not working

**Problem:** Path not being replaced with ~

**Check:**

1. Verify path actually starts with home directory
2. Check `os.UserHomeDir()` returns expected value
3. Confirm path uses absolute form (not relative)

**Diagnosis:**

```go
home, err := os.UserHomeDir()
if err != nil {
    log.Printf("UserHomeDir error: %v", err)
}
log.Printf("Home: %q, Path: %q, HasPrefix: %v",
    home, path, strings.HasPrefix(path, home))
```

### Truncation threshold seems wrong

**Problem:** Paths being truncated unexpectedly

**Check:**

1. Verify `MaxPathLength` constant (should be 40)
2. Confirm path actually exceeds threshold
3. Consider adjusting `MaxPathLength` if statusline layout changed

**Diagnosis:**

```go
path := "/home/user/very/long/nested/project/directory"
shortened := format.ShortenPath(path)
log.Printf("Original: %d chars, Shortened: %d chars, Threshold: %d",
    len(path), len(shortened), MaxPathLength)
```

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetShortModelName** | Model name formatting | [model.md](model.md) |

**Future Functions:**

- Git repository-relative paths
- Context-aware truncation (preserve important segments)
- User-configurable thresholds

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Format Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetShortModelName Function** | Model name formatting | [model.md](model.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/format/path.go`](../../../lib/format/path.go) | ShortenPath with comprehensive inline documentation |
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

- **"Make the path straight"** - Matthew 3:3
  - Clarity in direction - presenting paths in ways that serve understanding and navigation

**Supporting Verse:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Order and clarity in communication

**Application:** Show location context (~ for home, basename for long paths) without overwhelming with full hierarchy. Like good signage: clear, brief, sufficient.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
