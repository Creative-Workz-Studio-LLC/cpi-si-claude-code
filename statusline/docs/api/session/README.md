<div align="center">

# 📊 Session Library API

**Session Statistics Display Formatting**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Functions](https://img.shields.io/badge/Functions-4-blue?style=flat)
![Performance](https://img.shields.io/badge/Performance-%3C10μs-green?style=flat)

*Transform session statistics into formatted display strings for statusline*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [🔧 Public API](#-public-api) • [🚀 Quick Start](#-quick-start) • [🎯 Design](#-design-philosophy) • [📚 References](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - Session Library
**Purpose:** Presentation layer for session statistics
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

- [� Session Library API](#-session-library-api)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
  - [🔧 Public API](#-public-api)
    - [GetLinesModifiedDisplay](#getlinesmodifieddisplay)
    - [GetDurationDisplay](#getdurationdisplay)
    - [GetCostDisplay](#getcostdisplay)
    - [GetFormattedCost](#getformattedcost)
  - [🚀 Quick Start](#-quick-start)
    - [Basic Integration](#basic-integration)
  - [🎯 Design Philosophy](#-design-philosophy)
    - [Data → Presentation Pattern](#data--presentation-pattern)
    - [Visual Clarity](#visual-clarity)
    - [Space Optimization](#space-optimization)
    - [Zero-State Handling](#zero-state-handling)
  - [🧪 Demo: How It Works](#-demo-how-it-works)
    - [Session Statistics](#session-statistics)
    - [Statusline Integration Example](#statusline-integration-example)
  - [⚡ Performance](#-performance)
  - [🏗 Architecture](#-architecture)
  - [🔧 Extension Points](#-extension-points)
    - [Adding New Session Statistics](#adding-new-session-statistics)
    - [Custom Formatting Options](#custom-formatting-options)
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

The Session library transforms session statistics (lines modified, duration, cost) into formatted display strings suitable for statusline space constraints. It receives session data from SessionContext, outputs formatted display structures with visual indicators (colors, icons).

> [!IMPORTANT]
> **Pure presentation layer.** This library formats data for display. Data collection (lines tracking, duration tracking, cost calculation) handled elsewhere in SessionContext population.

**Core Principle:** Data vs Presentation separation - SessionContext provides data, session lib handles display formatting.

**Biblical Foundation:**

- **"So teach us to number our days, that we may apply our hearts unto wisdom"** - Psalm 90:12
  - Awareness of time and stewardship - tracking work to understand patterns and use time wisely
- **"Redeeming the time, because the days are evil"** - Ephesians 5:16
  - Wise stewardship of creative work

**Source Code:** [`lib/session/session.go`](../../../lib/session/session.go)

---

## 🔧 Public API

### GetLinesModifiedDisplay

```go
func GetLinesModifiedDisplay(ctx SessionContext) LinesModifiedDisplay
```

Formats lines changed (added + removed) for statusline display.

**Parameters:**

- `ctx` (`SessionContext`): Session context containing Cost.TotalLinesAdded and Cost.TotalLinesRemoved

**Returns:**

- `LinesModifiedDisplay` struct with formatted line count data

**Display Structure:**

```go
type LinesModifiedDisplay struct {
    TotalLines int    // Sum of lines added + removed
    Color      string // Color for terminal rendering
    Icon       string // Visual icon (📝)
    HasInfo    bool   // True if lines modified > 0
}
```

**Display Example:** `153 lines` (100 added + 53 removed)

**Behavior:**

- Sums TotalLinesAdded + TotalLinesRemoved for total activity
- Returns `HasInfo: false` when zero lines modified
- Cannot fail - graceful degradation for all edge cases

**Health Impact:** No health scoring - pure presentation library

**API Documentation:** [GetLinesModifiedDisplay Function Reference](lines.md)

---

### GetDurationDisplay

```go
func GetDurationDisplay(ctx SessionContext) DurationDisplay
```

Formats session duration in human-readable time format for statusline display.

**Parameters:**

- `ctx` (`SessionContext`): Session context containing Cost.TotalDurationMS

**Returns:**

- `DurationDisplay` struct with formatted duration string

**Display Structure:**

```go
type DurationDisplay struct {
    Duration string  // Human-readable duration (e.g., "1h 23m")
    Color    string  // Color for terminal rendering
    Icon     string  // Visual icon (⏱)
    HasInfo  bool    // True if duration > 0
}
```

**Display Examples:**

- `45s` (45 seconds)
- `1m 23s` (1 minute 23 seconds)
- `1h 23m` (1 hour 23 minutes)

**Behavior:**

- Uses `system/lib/sessiontime.FormatDuration()` for formatting
- Returns `HasInfo: false` when zero duration
- Cannot fail - graceful degradation for all edge cases

**Health Impact:** No health scoring - pure presentation library

**API Documentation:** [GetDurationDisplay Function Reference](duration.md)

---

### GetCostDisplay

```go
func GetCostDisplay(ctx SessionContext) CostDisplay
```

Formats session cost (USD with 4 decimal precision) for statusline display.

**Parameters:**

- `ctx` (`SessionContext`): Session context containing Cost.TotalCostUSD

**Returns:**

- `CostDisplay` struct with cost value

**Display Structure:**

```go
type CostDisplay struct {
    Cost    float64 // USD cost value
    Color   string  // Color for terminal rendering
    Icon    string  // Visual icon (💰)
    HasInfo bool    // True if cost > 0
}
```

**Display Example:** `$0.0123` (formatted with GetFormattedCost)

**Behavior:**

- Extracts cost value from SessionContext
- Returns `HasInfo: false` when zero cost
- Cannot fail - graceful degradation for all edge cases

**Health Impact:** No health scoring - pure presentation library

**API Documentation:** [GetCostDisplay Function Reference](cost.md)

---

### GetFormattedCost

```go
func GetFormattedCost(cost float64) string
```

Helper function to format cost as USD string with 4 decimal precision.

**Parameters:**

- `cost` (`float64`): Cost value in USD

**Returns:**

- `string`: Formatted cost string (e.g., "$0.0123")

**Format:** `$%.4f` (always 4 decimal places)

**Examples:**

- `0.0` → `"$0.0000"`
- `0.0123` → `"$0.0123"`
- `1.5` → `"$1.5000"`

**API Documentation:** [GetFormattedCost Function Reference](cost.md)

---

## 🚀 Quick Start

### Basic Integration

**Step 1: Import**

```go
import "statusline/lib/session"
```

**Step 2: Call During Statusline Assembly**

```go
func buildStatusline(ctx SessionContext) string {
    var parts []string

    // Lines modified
    linesDisplay := session.GetLinesModifiedDisplay(ctx)
    if linesDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%d lines", linesDisplay.TotalLines))
    }

    // Duration
    durationDisplay := session.GetDurationDisplay(ctx)
    if durationDisplay.HasInfo {
        parts = append(parts, durationDisplay.Duration)
    }

    // Cost
    costDisplay := session.GetCostDisplay(ctx)
    if costDisplay.HasInfo {
        parts = append(parts, session.GetFormattedCost(costDisplay.Cost))
    }

    // ... add other statusline parts ...

    return strings.Join(parts, " | ")
    // → "150 lines | 1h 30m | $0.0123"
}
```

**Step 3: No Cleanup Required**

Stateless library requires no cleanup or lifecycle management.

---

## 🎯 Design Philosophy

### Data → Presentation Pattern

**Clear Separation:**

```bash
SessionContext                    Display Functions
(Data Source)                     (Presentation)
     ↓                                  ↓
Cost.TotalLinesAdded: 100        GetLinesModifiedDisplay()
Cost.TotalLinesRemoved: 50             ↓
     ↓                            LinesModifiedDisplay{
                                    TotalLines: 150,
                                    Color: "yellow",
                                    Icon: "📝",
                                    HasInfo: true
                                  }

Cost.TotalDurationMS: 5400000    GetDurationDisplay()
     ↓                                  ↓
                                  DurationDisplay{
                                    Duration: "1h 30m",
                                    Color: "gray",
                                    Icon: "⏱",
                                    HasInfo: true
                                  }

Cost.TotalCostUSD: 0.0123        GetCostDisplay()
     ↓                                  ↓
                                  CostDisplay{
                                    Cost: 0.0123,
                                    Color: "yellow",
                                    Icon: "💰",
                                    HasInfo: true
                                  }
```

**SessionContext provides data, session lib formats for display.**

**Why this matters:** Never duplicates data tracking logic. Always uses SessionContext for data, formats it here for display.

### Visual Clarity

| Statistic | Icon | Purpose | Example |
|-----------|:----:|---------|---------|
| **Lines Modified** | 📝 | Code editing activity | `150 lines` |
| **Duration** | ⏱ | Time spent | `1h 23m` |
| **Cost** | 💰 | Token cost in USD | `$0.0123` |

### Space Optimization

**Human-readable formats optimized for statusline:**

- Lines: Total count only (not breakdown of added vs removed)
- Duration: Human-readable (`1h 23m` not `5400000ms`)
- Cost: 4 decimal USD (`$0.0123` not fractions of cents)

### Zero-State Handling

**Graceful degradation for edge cases:**

- Zero lines modified → `HasInfo: false`
- Zero duration → `HasInfo: false`
- Zero cost → `HasInfo: false`

**Philosophy:** This library cannot fail. Always returns valid display structures.

---

## 🧪 Demo: How It Works

### Session Statistics

```go
import "statusline/lib/session"
import "statusline/lib/types"

ctx := types.SessionContext{
    Cost: types.CostInfo{
        TotalLinesAdded: 100,
        TotalLinesRemoved: 50,
        TotalDurationMS: 5400000,  // 1h 30m
        TotalCostUSD: 0.0123,
    },
}

// Get session statistics displays
linesDisplay := session.GetLinesModifiedDisplay(ctx)
// → LinesModifiedDisplay{TotalLines: 150, HasInfo: true}

durationDisplay := session.GetDurationDisplay(ctx)
// → DurationDisplay{Duration: "1h 30m", HasInfo: true}

costDisplay := session.GetCostDisplay(ctx)
// → CostDisplay{Cost: 0.0123, HasInfo: true}

// Format for display
fmt.Printf("%s %d lines\n", linesDisplay.Icon, linesDisplay.TotalLines)
// → "📝 150 lines"

fmt.Printf("%s %s\n", durationDisplay.Icon, durationDisplay.Duration)
// → "⏱ 1h 30m"

fmt.Printf("%s %s\n", costDisplay.Icon, session.GetFormattedCost(costDisplay.Cost))
// → "💰 $0.0123"
```

### Statusline Integration Example

```go
func buildStatusline(ctx SessionContext) string {
    var parts []string

    // Model name
    parts = append(parts, format.GetShortModelName(ctx.Model.DisplayName))

    // Session statistics
    linesDisplay := session.GetLinesModifiedDisplay(ctx)
    if linesDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%d lines", linesDisplay.TotalLines))
    }

    durationDisplay := session.GetDurationDisplay(ctx)
    if durationDisplay.HasInfo {
        parts = append(parts, durationDisplay.Duration)
    }

    costDisplay := session.GetCostDisplay(ctx)
    if costDisplay.HasInfo {
        parts = append(parts, session.GetFormattedCost(costDisplay.Cost))
    }

    return strings.Join(parts, " | ")
    // → "Sonnet | 150 lines | 1h 30m | $0.0123"
}
```

---

## ⚡ Performance

**Function Performance:**

| Function | Complexity | Memory | Time |
|----------|:----------:|:------:|:----:|
| **GetLinesModifiedDisplay** | O(1) | ~50 bytes | <1μs |
| **GetDurationDisplay** | O(1) | ~50 bytes | <10μs |
| **GetCostDisplay** | O(1) | ~50 bytes | <1μs |
| **GetFormattedCost** | O(1) | ~20 bytes | <1μs |

**Total Memory:** Three struct allocations per statusline render (~150 bytes total)

**Optimization:** This library needs no optimization. Struct allocation and string formatting are already optimal for this use case.

---

## 🏗 Architecture

**Component Type:** Ladder (Library - Middle Rung)

**Role:** Presentation layer for session statistics

**Dependencies:**

- `system/lib/sessiontime` - Duration formatting (lower rung)
- `system/lib/display` - Color constants
- `statusline/lib/types` - SessionContext struct

**Used By:**

- `statusline` - Main orchestrator (higher rung)

**Design Principle:** Data collection happens elsewhere (SessionContext population), presentation happens here. Never duplicate data tracking logic.

---

## 🔧 Extension Points

### Adding New Session Statistics

**Pattern:**

1. Define new display structure in SETUP Types
2. Create `Get[Statistic]Display()` function in BODY
3. Extract data from SessionContext (or extend SessionContext if needed)
4. Follow existing patterns: check zero-state, assign color/icon
5. Update API documentation with new statistic
6. Add tests for new display function

**Example - Adding Commits Display:**

```go
// In SETUP Types:
type CommitsDisplay struct {
    CommitCount int
    Color       string
    Icon        string
    HasInfo     bool
}

// In BODY:
func GetCommitsDisplay(ctx types.SessionContext) CommitsDisplay {
    if ctx.Git.CommitCount <= 0 {
        return CommitsDisplay{HasInfo: false}
    }

    return CommitsDisplay{
        CommitCount: ctx.Git.CommitCount,
        Color:       display.Green,
        Icon:        "📦",
        HasInfo:     true,
    }
}
```

### Custom Formatting Options

**Future Enhancement:** Configurable display formats (compact, verbose, detailed)

```go
// Potential future API
func GetLinesModifiedDisplayCompact(ctx) LinesModifiedDisplay  // "150L"
func GetLinesModifiedDisplayVerbose(ctx) LinesModifiedDisplay  // "100 added, 50 removed"
```

---

## 🔍 Troubleshooting

**This library has no common failure modes - all operations guaranteed to succeed.**

### Expected Behaviors

**Zero values return HasInfo: false**

- Zero lines modified → `HasInfo: false` (correct, not an error)
- Zero duration → `HasInfo: false` (correct, not an error)
- Zero cost → `HasInfo: false` (correct, not an error)

**Empty display data when HasInfo: false**

- This is correct zero-state
- Calling code checks HasInfo to determine display

### If Unexpected Results Occur

**Problem:** HasInfo true but display data seems wrong

- **Cause:** SessionContext populated incorrectly
- **Solution:** Check where SessionContext gets populated
- **Verify:** `Cost.TotalLinesAdded`, `Cost.TotalLinesRemoved`, `Cost.TotalDurationMS`, `Cost.TotalCostUSD`

**Problem:** Duration format looks wrong

- **Cause:** `system/lib/sessiontime.FormatDuration()` behavior changed
- **Solution:** Check `system/lib/sessiontime` implementation
- **Expected formats:** `"45s"`, `"1m 23s"`, `"1h 23m"`, `"2h 15m 30s"`

**Problem:** Cost precision incorrect (not 4 decimals)

- **Cause:** `GetFormattedCost()` format string modified
- **Solution:** Verify `fmt.Sprintf("$%.4f", cost)` unchanged
- **Expected:** Always 4 decimal places (`$0.0000`, `$0.0123`, `$1.5000`)

**Problem:** Lines calculation wrong

- **Cause:** `TotalLinesAdded + TotalLinesRemoved` not matching expectations
- **Solution:** Check how lines are tracked (what counts as "added" vs "removed")
- **Note:** Simple addition - any calculation error is upstream in data collection

---

## 📚 References & Resources

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **GetLinesModifiedDisplay** | Complete function documentation | [lines.md](lines.md) |
| **GetDurationDisplay** | Complete function documentation | [duration.md](duration.md) |
| **GetCostDisplay & GetFormattedCost** | Complete function documentation | [cost.md](cost.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Session Formatting** | [`lib/session/session.go`](../../../lib/session/session.go) | All session display functions with comprehensive inline documentation |
| **Library Overview** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |
| **Duration Formatting** | [`system/lib/sessiontime`](~/.claude/cpi-si/system/lib/sessiontime/) | Duration formatting utility |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 🗺 Future Roadmap

**Planned Features:**

- ✅ Lines modified display - COMPLETED
- ✅ Session duration display - COMPLETED
- ✅ Cost tracking display - COMPLETED
- ⏳ Commit count display
- ⏳ Files changed display
- ⏳ Token usage display (input/output breakdown)
- ⏳ Custom formatting options (compact, verbose, detailed)

**Research Areas:**

- Color coding based on thresholds (high cost = red, normal = yellow)
- Configurable cost precision (2 vs 4 decimal places)
- Human-readable large numbers (1.5K lines instead of 1500)
- Session comparison (current vs previous session stats)
- Breakdown displays (lines added vs removed separately)

**Known Limitations:**

- No breakdown of lines added vs removed (just total)
- Duration loses seconds for long sessions (shows `"1h 23m"`, not `"1h 23m 45s"`)
- Cost always 4 decimals even for large values (`$100.0000` looks odd)
- No configuration options (colors, icons, formats all hardcoded)

---

## 📋 Modification Policy

**Safe to Modify (Extension Points):**

- ✅ Add new session statistic displays (commits, files, tokens, etc.)
- ✅ Adjust color or icon choices for different statistics
- ✅ Extend display structures with additional fields
- ✅ Add helper functions for complex formatting logic
- ✅ Create alternative formatting functions (compact, verbose, etc.)

**Modify with Extreme Care (Breaking Changes):**

- ⚠️ Display struct fields - breaks statusline orchestrator
- ⚠️ Function signatures - breaks all calling code
- ⚠️ HasInfo semantics - breaks zero-state handling
- ⚠️ Color/Icon field types - breaks display rendering
- ⚠️ GetFormattedCost() format - breaks cost display expectations

**NEVER Modify (Foundational):**

- ❌ Pure function guarantee (stateless, no side effects)
- ❌ Graceful degradation (always return valid display)
- ❌ Data vs Presentation separation (SessionContext provides data)
- ❌ Non-blocking guarantee (no I/O operations)

---

## 📖 Biblical Foundation

**Core Verse:**

- **"So teach us to number our days, that we may apply our hearts unto wisdom"** - Psalm 90:12
  - Awareness of time and stewardship - tracking work to understand patterns and use time wisely

**Supporting Verse:**

- **"Redeeming the time, because the days are evil"** - Ephesians 5:16
  - Wise stewardship of creative work

**Application:** Show session progress (lines changed, time spent, cost incurred) to foster awareness and wise stewardship of creative work. Not tracking for tracking's sake, but for wisdom in how time and resources are used.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
