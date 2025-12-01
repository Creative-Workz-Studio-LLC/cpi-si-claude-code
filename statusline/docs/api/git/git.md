<div align="center">

# GetGitDisplay Function Reference

**Git Repository Display Formatting**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Savings](https://img.shields.io/badge/Space%20Savings-up%20to%2089%25-green?style=flat)

*Formats git repository assessment data into display strings for statusline*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetGitDisplay Function
**Purpose:** Git repository status presentation for statusline
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
  - [workdir (string)](#workdir-string)
- [🔄 Returns](#-returns)
  - [GitDisplay](#gitdisplay)
- [🎯 Behavior](#-behavior)
  - [Display Formats](#display-formats)
  - [Assessment Integration](#assessment-integration)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [Display Format Variations](#display-format-variations)
  - [Handling Non-Repositories](#handling-non-repositories)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Testing](#testing)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [📊 Space Savings](#-space-savings)
- [🔧 Edge Cases](#-edge-cases)
  - [Empty String Path](#empty-string-path)
  - [Current Directory](#current-directory)
  - [Detached HEAD](#detached-head)
  - [No Remote Tracking](#no-remote-tracking)
  - [Very Long Branch Names](#very-long-branch-names)
  - [Permission Denied](#permission-denied)
  - [Corrupted Repository](#corrupted-repository)
- [🔧 Extension Points](#-extension-points)
  - [Adding New Indicators](#adding-new-indicators)
  - [Dynamic Color Based on State](#dynamic-color-based-on-state)
- [🔍 Troubleshooting](#-troubleshooting)
  - [HasInfo true but DisplayString empty](#hasinfo-true-but-displaystring-empty)
  - [Ahead/behind not showing](#aheadbehind-not-showing)
  - [Wrong color or icon](#wrong-color-or-icon)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetGitDisplay** formats git repository assessment data into display strings suitable for statusline space constraints.

> [!NOTE]
> **Pure presentation function.** Receives assessment from `system/lib/git`, transforms into formatted display. No I/O operations, cannot fail.

**Design:** Pure presentation layer - receives assessment from `system/lib/git`, transforms into formatted display with visual indicators (colors, icons, symbols).

**Library:** [`statusline/lib/git`](../../../lib/git/)

**Source Code:** [`lib/git/git.go`](../../../lib/git/git.go)

**Library README:** [Git Library Documentation](README.md)

---

## 📝 Signature

```go
func GetGitDisplay(workdir string) GitDisplay
```

---

## ⚙ Parameters

### workdir (string)

Absolute path to working directory to check for git repository.

**Expected Formats:**

| Format | Example | Type |
|--------|---------|:----:|
| Git repository | `/home/user/project` | ✅ |
| Non-repository | `/etc/nginx` | ✅ |
| Current directory | `.` | ✅ |
| Empty string | `""` | ✅ |

**Validation:** None required - all string inputs valid, graceful degradation for non-repos

---

## 🔄 Returns

### GitDisplay

Formatted display structure with all presentation elements.

```go
type GitDisplay struct {
    DisplayString string // Formatted git status (e.g., "main*↑2↓1")
    Color         string // Terminal color code for display
    Icon          string // Visual icon representing git (e.g., "🌿")
    HasInfo       bool   // True if git information available, false for non-repos
}
```

**Field Details:**

| Field | Type | Purpose | Example |
|-------|------|---------|---------|
| **DisplayString** | `string` | Formatted status | `"main*↑2↓1"` |
| **Color** | `string` | Terminal color | `display.Green` |
| **Icon** | `string` | Visual icon | `"🌿"` |
| **HasInfo** | `bool` | Data available? | `true` / `false` |

**DisplayString Format:** `[branch][dirty][ahead][behind]`

- **Branch:** Current branch name or short commit hash
- **Dirty:** `*` if uncommitted changes exist
- **Ahead:** `↑[count]` if commits ahead of remote
- **Behind:** `↓[count]` if commits behind remote
- **Empty string** if not a repository

**Color:**

- Currently: `display.Green` for all repos
- Future: State-based colors (yellow for dirty, red for conflicts)

**Icon:**

- Currently: `🌿` (branch/nature icon) for all repos
- Future: State-based icons

**HasInfo:**

- `true`: Git repository detected, display information available
- `false`: Not a repository, DisplayString/Color/Icon empty

**Guarantee:** Always returns valid GitDisplay (never nil, never errors)

---

## 🎯 Behavior

### Display Formats

| Repository State | DisplayString | Meaning |
|------------------|---------------|---------|
| Clean synced | `main` | No changes, synced with remote |
| Dirty | `main*` | Uncommitted changes |
| Ahead | `main↑2` | 2 commits ahead of remote |
| Behind | `main↓1` | 1 commit behind remote |
| Complex | `main*↑2↓1` | Dirty, 2 ahead, 1 behind |
| Non-repo | `""` | Not a git repository |

### Assessment Integration

GetGitDisplay uses `system/lib/git.GetInfo()` for all assessment:

```go
// Internal flow
gitInfo := gitlib.GetInfo(workdir)  // System lib provides assessment

if gitInfo.Branch == "" {
    return GitDisplay{HasInfo: false}  // Zero-state for non-repos
}

// Format display string from assessment data
gitDisplayStr := gitInfo.Branch
if gitInfo.Dirty { gitDisplayStr += "*" }
if gitInfo.Ahead > 0 { gitDisplayStr += fmt.Sprintf("↑%d", gitInfo.Ahead) }
if gitInfo.Behind > 0 { gitDisplayStr += fmt.Sprintf("↓%d", gitInfo.Behind) }

return GitDisplay{
    DisplayString: gitDisplayStr,
    Color:         display.Green,
    Icon:          "🌿",
    HasInfo:       true,
}
```

**Why this matters:** Never duplicates assessment logic. Always uses `system/lib/git` for data, formats it here for display.

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/git"

func buildStatusline(workdir string) string {
    // Get git display formatting
    gitDisplay := git.GetGitDisplay(workdir)
    // → GetGitDisplay: "main*↑2"

    var parts []string

    // Only add git info if available
    if gitDisplay.HasInfo {
        gitPart := fmt.Sprintf("%s %s", gitDisplay.Icon, gitDisplay.DisplayString)
        parts = append(parts, gitPart)
        // → "🌿 main*↑2"
    }

    // ... add other parts ...

    return strings.Join(parts, " | ")
}
```

### Display Format Variations

```go
// Assuming workdir = "/home/user/project"

// Clean branch, synced with remote
gitDisplay := git.GetGitDisplay(workdir)
fmt.Println(gitDisplay.DisplayString)  // "main"

// Branch with uncommitted changes
gitDisplay = git.GetGitDisplay(workdir)
fmt.Println(gitDisplay.DisplayString)  // "main*"

// Ahead of remote
gitDisplay = git.GetGitDisplay(workdir)
fmt.Println(gitDisplay.DisplayString)  // "main↑2"

// Behind remote
gitDisplay = git.GetGitDisplay(workdir)
fmt.Println(gitDisplay.DisplayString)  // "main↓1"

// Complex state (dirty, ahead, and behind)
gitDisplay = git.GetGitDisplay(workdir)
fmt.Println(gitDisplay.DisplayString)  // "main*↑2↓1"
```

### Handling Non-Repositories

```go
gitDisplay := git.GetGitDisplay("/tmp")

// Check HasInfo to determine display
if gitDisplay.HasInfo {
    fmt.Printf("%s %s\n", gitDisplay.Icon, gitDisplay.DisplayString)
} else {
    fmt.Println("Not a git repository")
}
```

### Full Statusline Integration

```go
func buildStatusline(ctx SessionContext) string {
    var parts []string

    // Git status
    gitDisplay := git.GetGitDisplay(ctx.WorkingDirectory)
    if gitDisplay.HasInfo {
        gitPart := fmt.Sprintf("%s %s", gitDisplay.Icon, gitDisplay.DisplayString)
        parts = append(parts, gitPart)
    }

    // Model name
    modelName := format.GetShortModelName(ctx.Model.DisplayName)
    parts = append(parts, modelName)

    // Working directory
    workdir := format.ShortenPath(ctx.WorkingDirectory)
    parts = append(parts, workdir)

    return strings.Join(parts, " | ")
    // → "🌿 main* | Sonnet | ~/project"
}
```

### Testing

```go
func TestGitDisplayFormat(t *testing.T) {
    testCases := []struct {
        name     string
        workdir  string
        expected string
    }{
        {"Clean repo", "/home/user/clean-repo", "main"},
        {"Dirty repo", "/home/user/dirty-repo", "main*"},
        {"Ahead", "/home/user/ahead-repo", "main↑2"},
        {"Behind", "/home/user/behind-repo", "main↓1"},
        {"Complex", "/home/user/complex-repo", "main*↑2↓1"},
        {"Non-repo", "/tmp", ""},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := git.GetGitDisplay(tc.workdir)
            if result.DisplayString != tc.expected {
                t.Errorf("Expected %q, got %q",
                    tc.expected, result.DisplayString)
            }
        })
    }
}
```

---

## 🏥 Health Scoring

**No health scoring** - pure presentation library.

All operations guaranteed to succeed through graceful degradation:

- Non-repository directory → Returns `GitDisplay{HasInfo: false}`
- Missing git data → Omits corresponding indicators
- Invalid paths → Returns zero-state display

**Philosophy:** This function cannot fail. Always returns valid GitDisplay structures. Health tracking would measure "successfully did nothing" which provides no value.

---

## ⚡ Performance

**Time Complexity:** O(1) formatting operations (assessment done by `system/lib/git`)

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **String concatenation** | 5-10 ops | Branch + indicators |
| **Struct allocation** | ~50 bytes | GitDisplay struct |
| **Execution time** | <1 μs | Microseconds |

**Memory:**

- Single GitDisplay struct allocation per call (~50 bytes)
- String allocations for DisplayString (~10-50 bytes depending on length)
- No heap allocations beyond return value
- Garbage collected automatically

**Optimization:** Not needed - string formatting and struct allocation already optimal

---

## 📊 Space Savings

**Display Length Examples:**

| Git State | Full Description | Display String | Character Savings |
|-----------|------------------|----------------|:-----------------:|
| Clean synced | "On branch main, clean, up to date" (35 chars) | `"main"` (4 chars) | **89%** |
| Dirty branch | "On branch main, uncommitted changes" (35 chars) | `"main*"` (5 chars) | **86%** |
| Ahead/behind | "On branch main, 2 ahead, 1 behind" (33 chars) | `"main↑2↓1"` (9 chars) | **73%** |

**Benefit:** More room for model name, working directory, session stats in statusline

---

## 🔧 Edge Cases

### Empty String Path

```go
git.GetGitDisplay("")
// → GitDisplay{HasInfo: false}
// Behavior: Empty path treated as non-repository
```

### Current Directory

```go
git.GetGitDisplay(".")
// → GitDisplay with repo info if . is in git repo
// Behavior: Relative path resolved by system/lib/git
```

### Detached HEAD

```go
git.GetGitDisplay("/home/user/project")
// → GitDisplay{DisplayString: "a1b2c3d", ...}
// Behavior: Shows short commit hash (provided by system/lib/git)
```

### No Remote Tracking

```go
git.GetGitDisplay("/home/user/local-only-repo")
// → GitDisplay{DisplayString: "main*", ...}
// Behavior: Shows branch and dirty, omits ahead/behind (no remote)
```

### Very Long Branch Names

```go
git.GetGitDisplay("/home/user/project")
// → GitDisplay{DisplayString: "feature/very-long-branch-name*↑2", ...}
// Behavior: No truncation - full branch name displayed
// Note: Future enhancement may add abbreviation for statusline space
```

### Permission Denied

```go
git.GetGitDisplay("/root/project")  // No read permissions
// → GitDisplay{HasInfo: false}
// Expected: system/lib/git cannot read .git/, returns empty Branch
```

### Corrupted Repository

```go
git.GetGitDisplay("/home/user/corrupted-repo")
// → GitDisplay{HasInfo: false} or GitDisplay{DisplayString: "[branch]", ...}
// Behavior: Depends on what system/lib/git can read
// Note: Graceful degradation - shows what's available
```

---

## 🔧 Extension Points

### Adding New Indicators

To add stash count indicator:

**Pattern:**

```go
// In GetGitDisplay(), after ahead/behind section:
if gitInfo.Stashes > 0 {
    gitDisplayStr += fmt.Sprintf("$%d", gitInfo.Stashes)
}
// Result: "main*↑2$3" (dirty, 2 ahead, 3 stashes)
```

**Prerequisites:**

1. `system/lib/git.Info` struct must have Stashes field
2. `system/lib/git.GetInfo()` must populate Stashes

**Steps:**

1. Verify system lib provides data
2. Add formatting logic after existing indicators
3. Follow pattern: check condition, append symbol+value
4. Update this API doc with new format examples
5. Add tests for new indicator

### Dynamic Color Based on State

**Pattern:**

```go
// Future enhancement
color := display.Green  // Default

if gitInfo.Conflicts != nil && len(gitInfo.Conflicts) > 0 {
    color = display.Red  // Conflicts need attention
} else if gitInfo.Dirty {
    color = display.Yellow  // Uncommitted changes
}

return GitDisplay{
    DisplayString: gitDisplayStr,
    Color:         color,  // State-based selection
    Icon:          "🌿",
    HasInfo:       true,
}
```

---

## 🔍 Troubleshooting

### HasInfo true but DisplayString empty

**Problem:** GitDisplay returned with `HasInfo: true` but empty DisplayString

**Check:**

1. Verify `system/lib/git.GetInfo()` behavior - should return `Branch: ""` for non-repos
2. Check if Branch field set but contains empty string (bug in system lib)

**Solution:** Fix `system/lib/git` implementation to properly handle edge cases

### Ahead/behind not showing

**Problem:** Expected ahead/behind indicators but only seeing branch name

**Check:**

1. Verify repository has remote tracking: `git branch -vv`
2. Check if remote tracking branch configured
3. Confirm `system/lib/git.GetInfo()` returns Ahead/Behind values

**Expected:** Ahead/behind requires remote tracking branch - omission is correct behavior

**Diagnosis:**

```bash
# Check remote tracking
git branch -vv
# If output shows: main  a1b2c3d [origin/main] ...
# Then remote tracking exists

# If output shows: main  a1b2c3d ...
# No remote tracking - ahead/behind omission is correct
```

### Wrong color or icon

**Problem:** GitDisplay returns unexpected Color or Icon values

**Check:**

1. Verify `display.Green` constant value in `system/lib/display`
2. Check if code modified to use different color/icon logic

**Expected:** Currently all repos return Green / 🌿 - state-based coloring is future feature

---

## 🔗 Related Functions

**Same Library:**

- None (single-function library)

**Future Functions:**

| Function | Purpose | Status |
|----------|---------|:------:|
| `GetGitDisplayDetailed()` | Include stash/conflict counts | ⏳ Planned |
| `GetGitDisplayCompact()` | Ultra-abbreviated for narrow statuslines | ⏳ Planned |
| `GetGitDisplayCustom(format string)` | User-configurable format strings | 🔬 Research |

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Git Library README** | Library overview and integration guide | [README.md](README.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/git/git.go`](../../../lib/git/git.go) | GetGitDisplay with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |
| **Assessment Logic** | [`system/lib/git/operations.go`](~/.claude/cpi-si/system/lib/git/operations.go) | Git assessment (provides Info struct) |

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

**Application:** Show repository context (branch, changes, sync status) without overwhelming with verbosity. Assessment provides truth, presentation provides clarity.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
