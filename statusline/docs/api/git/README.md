<div align="center">

# 🌿 Git Library API

**Git Repository Display Formatting**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Function-1-blue?style=flat)
![Performance](https://img.shields.io/badge/Performance-%3C1μs-green?style=flat)

*Transform git assessment data into formatted display strings for statusline*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [🔧 Public API](#-public-api) • [🚀 Quick Start](#-quick-start) • [🎯 Design](#-design-philosophy) • [📚 References](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - Git Library
**Purpose:** Presentation layer for git repository status
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

- [🌿 Git Library API](#-git-library-api)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
  - [🔧 Public API](#-public-api)
    - [GetGitDisplay](#getgitdisplay)
  - [🚀 Quick Start](#-quick-start)
    - [Basic Integration](#basic-integration)
  - [🎯 Design Philosophy](#-design-philosophy)
    - [Assessment → Presentation Pattern](#assessment--presentation-pattern)
    - [Visual Clarity](#visual-clarity)
    - [Space Optimization](#space-optimization)
    - [Zero-State Handling](#zero-state-handling)
  - [🧪 Demo: How It Works](#-demo-how-it-works)
    - [Repository States](#repository-states)
    - [Statusline Integration Example](#statusline-integration-example)
  - [⚡ Performance](#-performance)
  - [🏗 Architecture](#-architecture)
  - [🔧 Extension Points](#-extension-points)
    - [Adding New Git Indicators](#adding-new-git-indicators)
    - [Custom Visual Elements](#custom-visual-elements)
  - [🔍 Troubleshooting](#-troubleshooting)
    - [Expected Behaviors](#expected-behaviors)
    - [If Unexpected Results Occur](#if-unexpected-results-occur)
  - [📚 References \& Resources](#-references--resources)
    - [API Documentation](#api-documentation)
    - [Source Code](#source-code)
    - [Related Documentation](#related-documentation)
  - [🗺 Future Roadmap](#-future-roadmap)
  - [📋 Modification Policy](#-modification-policy)
  - [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

The Git library transforms git repository assessment data into formatted display strings suitable for statusline space constraints. It receives assessment from `system/lib/git`, outputs formatted display structures with visual indicators (colors, icons, symbols).

> [!IMPORTANT]
> **Pure presentation layer.** This library formats data for display. Assessment logic (git commands, branch detection, dirty checks) handled by `system/lib/git`.

**Core Principle:** Assessment and presentation separation - system lib handles git logic, statusline lib handles display formatting.

**Biblical Foundation:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Clarity and order in communication - presenting information in ways that serve understanding

**Source Code:** [`lib/git/git.go`](../../../lib/git/git.go)

---

## 🔧 Public API

### GetGitDisplay

```go
func GetGitDisplay(workdir string) GitDisplay
```

Formats git repository information for statusline display.

**Parameters:**

- `workdir` (`string`): Working directory path to check for git repository

**Returns:**

- `GitDisplay` struct with formatted display data

**Display Structure:**

```go
type GitDisplay struct {
    DisplayString string  // Formatted status (e.g., "main*↑2")
    Color         string  // Color for terminal rendering
    Icon          string  // Visual icon (🌿)
    HasInfo       bool    // True if git repository, false otherwise
}
```

**Display Examples:**

| Repository State | DisplayString | Meaning |
|------------------|---------------|---------|
| Clean synced branch | `main` | No changes, synced with remote |
| Dirty branch | `main*` | Uncommitted changes |
| Ahead of remote | `main↑2` | 2 commits ahead |
| Behind remote | `main↓1` | 1 commit behind |
| Complex state | `main*↑2↓1` | Dirty, 2 ahead, 1 behind |

**Behavior:**

- Shows current branch name or short commit hash (detached HEAD)
- Adds `*` indicator for uncommitted changes
- Adds `↑N` for commits ahead of remote
- Adds `↓N` for commits behind remote
- Returns `HasInfo: false` for non-repository directories
- Cannot fail - graceful degradation for all edge cases

**Health Impact:** No health scoring - pure presentation library. All operations guaranteed to succeed through graceful degradation.

**API Documentation:** [GetGitDisplay Function Reference](git.md)

---

## 🚀 Quick Start

### Basic Integration

**Step 1: Import**

```go
import "statusline/lib/git"
```

**Step 2: Call During Statusline Assembly**

```go
func buildStatusline(workdir string) string {
    var parts []string

    // Get git display information
    gitDisplay := git.GetGitDisplay(workdir)

    // Add git info if available
    if gitDisplay.HasInfo {
        gitPart := fmt.Sprintf("%s %s", gitDisplay.Icon, gitDisplay.DisplayString)
        parts = append(parts, gitPart)
    }

    // ... add other statusline parts ...

    return strings.Join(parts, " | ")
}
```

**Step 3: No Cleanup Required**

Stateless library requires no cleanup or lifecycle management.

---

## 🎯 Design Philosophy

### Assessment → Presentation Pattern

**Clear Separation:**

```bash
system/lib/git                    statusline/lib/git
(Assessment)                      (Presentation)
     ↓                                  ↓
GetInfo(workdir)              GetGitDisplay(workdir)
     ↓                                  ↓
Info{                           GitDisplay{
  Branch: "main",                 DisplayString: "main*↑2",
  Dirty: true,                    Color: "green",
  Ahead: 2,                       Icon: "🌿",
  Behind: 0                       HasInfo: true
}                               }
```

**System lib handles:**

- Git command execution
- Branch detection
- Dirty status checks
- Ahead/behind calculation

**Statusline lib handles:**

- String formatting
- Symbol selection
- Color assignment
- Icon choice

**Why this matters:** Never duplicate assessment logic. Always use `system/lib/git` for data, format it here for display.

### Visual Clarity

| Element | Purpose | Example |
|---------|---------|---------|
| **Branch name** | Context identification | `main`, `feature/auth` |
| **Dirty indicator (*)** | Uncommitted changes present | `main*` |
| **Ahead arrow (↑)** | Commits ahead of remote | `main↑2` |
| **Behind arrow (↓)** | Commits behind remote | `main↓1` |
| **Color** | Quick visual recognition | Green (success) |
| **Icon** | Categorical identification | 🌿 (git branch) |

### Space Optimization

**Compact format for statusline constraints:**

- Branch name only (not full ref)
- Symbols instead of words (`*` not "dirty")
- Count only when non-zero (`↑2` not `↑0`)
- Combined indicators (`main*↑2↓1` not separate fields)

### Zero-State Handling

**Graceful degradation for edge cases:**

- Non-repository directory → `HasInfo: false`
- Missing git data → Omits corresponding indicators
- Invalid paths → Returns zero-state display

**Philosophy:** This library cannot fail. Always returns valid display structures.

---

## 🧪 Demo: How It Works

### Repository States

```go
import "statusline/lib/git"

// Clean repository, synced with remote
gitDisplay := git.GetGitDisplay("/home/user/project")
// → DisplayString: "main", HasInfo: true

// Dirty repository (uncommitted changes)
gitDisplay = git.GetGitDisplay("/home/user/dirty-project")
// → DisplayString: "main*", HasInfo: true

// Ahead of remote
gitDisplay = git.GetGitDisplay("/home/user/ahead-project")
// → DisplayString: "main↑2", HasInfo: true

// Behind remote
gitDisplay = git.GetGitDisplay("/home/user/behind-project")
// → DisplayString: "main↓1", HasInfo: true

// Complex state (dirty + ahead + behind)
gitDisplay = git.GetGitDisplay("/home/user/complex-project")
// → DisplayString: "main*↑2↓1", HasInfo: true

// Not a git repository
gitDisplay = git.GetGitDisplay("/etc/config")
// → DisplayString: "", HasInfo: false
```

### Statusline Integration Example

```go
func buildStatusline(ctx SessionContext) string {
    var parts []string

    // Model name
    parts = append(parts, format.GetShortModelName(ctx.Model.DisplayName))

    // Git status (if repository)
    gitDisplay := git.GetGitDisplay(ctx.Workspace.CurrentDir)
    if gitDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%s %s", gitDisplay.Icon, gitDisplay.DisplayString))
    }

    // Working directory
    parts = append(parts, format.ShortenPath(ctx.Workspace.CurrentDir))

    return strings.Join(parts, " | ")
    // → "Sonnet | 🌿 main*↑2 | ~/project"
}
```

---

## ⚡ Performance

**Function Performance:**

- `GetGitDisplay()`: O(1) formatting, ~50 byte allocation, <1 microsecond
- No I/O operations (assessment handled by `system/lib/git`)
- Single struct allocation per call
- Garbage collected automatically

**Memory:**

- GitDisplay struct (~100 bytes)
- DisplayString allocation (variable, typically 5-30 bytes)
- No persistent state or large buffers

**Optimization:** This library needs no optimization. String formatting and struct allocation are already optimal for this use case.

---

## 🏗 Architecture

**Component Type:** Ladder (Library - Middle Rung)

**Role:** Presentation layer for git repository status

**Dependencies:**

- `system/lib/git` - Assessment data (lower rung)
- `system/lib/display` - Color constants

**Used By:**

- `statusline` - Main orchestrator (higher rung)

**Design Principle:** Assessment logic stays in system lib, presentation logic stays here. Never duplicate assessment - always use system lib for data, format it here for display.

---

## 🔧 Extension Points

### Adding New Git Indicators

**Pattern:**

1. Ensure `system/lib/git.Info` struct provides the data field
2. In `GetGitDisplay()`, after ahead/behind section, add new indicator logic
3. Follow pattern: `if gitInfo.[Field] [condition] { gitDisplayStr += "[symbol][value]" }`
4. Update API documentation with new indicator examples
5. Update tests to verify new indicator

**Example - Adding Stash Indicator:**

```go
// In GetGitDisplay(), after ahead/behind section
if gitInfo.Stashes > 0 {
    gitDisplayStr += fmt.Sprintf("$%d", gitInfo.Stashes)
}
// Result: "main*↑2$3" (dirty, 2 ahead, 3 stashes)
```

### Custom Visual Elements

**Adjusting Colors/Icons:**

```go
// Current implementation
return GitDisplay{
    DisplayString: gitDisplayStr,
    Color:         display.Green,  // Modify color choice
    Icon:          "🌿",            // Modify icon choice
    HasInfo:       true,
}
```

**Future Enhancement:** Color/icon selection based on git state (dirty=yellow, conflicts=red)

---

## 🔍 Troubleshooting

**This library has no common failure modes - all operations guaranteed to succeed.**

### Expected Behaviors

**Non-repository returns HasInfo: false**

- This is correct, not an error
- Zero-state response for directories without `.git/`

**No remote tracking omits ahead/behind**

- This is correct, not an error
- Ahead/behind requires remote tracking branch
- Shows branch and dirty status only

**Empty DisplayString when HasInfo: false**

- This is correct zero-state
- Calling code checks HasInfo to determine display

### If Unexpected Results Occur

**Problem:** HasInfo true but DisplayString empty

- **Cause:** `system/lib/git.GetInfo()` returned Info with Branch set but empty
- **Solution:** Check `system/lib/git` implementation - should return `Branch: ""`

**Problem:** Ahead/behind not showing when expected

- **Cause:** Repository has no remote tracking branch
- **Solution:** Expected behavior - verify with `git branch -vv`
- **Note:** `system/lib/git` returns Ahead: 0, Behind: 0 when no remote

---

## 📚 References & Resources

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **GetGitDisplay Function** | Complete function documentation | [git.md](git.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Display Formatting** | [`lib/git/git.go`](../../../lib/git/git.go) | GetGitDisplay with comprehensive inline documentation |
| **Library Overview** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |
| **Assessment Logic** | [`system/lib/git`](~/.claude/cpi-si/system/lib/git/) | Git data assessment (provides Info struct) |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 🗺 Future Roadmap

**Planned Features:**

- ✅ Branch name display - COMPLETED
- ✅ Dirty status indicator - COMPLETED
- ✅ Ahead/behind tracking - COMPLETED
- ⏳ Stash indicator (if `system/lib/git` provides data)
- ⏳ Conflict indicator (if `system/lib/git` provides data)
- ⏳ Multiple remote tracking support
- ⏳ Custom indicator symbols via configuration

**Research Areas:**

- Color coding based on git state (dirty=yellow, conflicts=red)
- Abbreviated branch names for very long branch names
- User-configurable display format strings
- Icon variations based on git state

**Known Limitations:**

- Only shows single remote tracking branch (most repos have one remote)
- No visual distinction between different dirty states (staged vs unstaged)
- Fixed color (green) regardless of git state
- No configuration options (format, symbols, colors all hardcoded)

---

## 📋 Modification Policy

**Safe to Modify (Extension Points):**

- ✅ Add new visual indicators (stashes, conflicts, etc.)
- ✅ Adjust color or icon choices for different git states
- ✅ Extend DisplayString format with additional symbols
- ✅ Add helper functions for complex formatting logic
- ✅ Create additional display structure types

**Modify with Extreme Care (Breaking Changes):**

- ⚠️ `GitDisplay` struct fields - breaks statusline orchestrator
- ⚠️ `GetGitDisplay()` signature - breaks all calling code
- ⚠️ DisplayString format - breaks expectations in consumers
- ⚠️ HasInfo semantics - breaks zero-state handling

**NEVER Modify (Foundational):**

- ❌ Pure function guarantee (stateless, no side effects)
- ❌ Graceful degradation (always return valid display)
- ❌ Assessment vs Presentation separation (system lib does assessment)
- ❌ Non-blocking guarantee (no I/O operations)

---

## 📖 Biblical Foundation

**Core Verse:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Clarity and order in communication - presenting information in ways that serve understanding

**Application:** Show repository context without overwhelming - branch, changes, sync status at a glance. Assessment and presentation working together to create clarity.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
