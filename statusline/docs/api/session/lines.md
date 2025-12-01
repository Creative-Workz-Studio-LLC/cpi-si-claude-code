<div align="center">

# GetLinesModifiedDisplay Function Reference

**Lines Modified Tracking for Session Statistics**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Calculation](https://img.shields.io/badge/Calculation-Added%20%2B%20Removed-orange?style=flat)

*Formats lines modified information (added + removed) for statusline display*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetLinesModifiedDisplay Function
**Purpose:** Lines modified presentation for session statistics
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
  - [ctx (SessionContext)](#ctx-sessioncontext)
- [🔄 Returns](#-returns)
  - [LinesModifiedDisplay](#linesmodifieddisplay)
- [🎯 Behavior](#-behavior)
  - [Calculation Logic](#calculation-logic)
  - [Display Examples](#display-examples)
  - [Semantic Meaning](#semantic-meaning)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [With Color](#with-color)
  - [Zero-State Handling](#zero-state-handling)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Testing](#testing)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [🔧 Edge Cases](#-edge-cases)
  - [Zero Total (No Changes)](#zero-total-no-changes)
  - [Negative Values (Data Error)](#negative-values-data-error)
  - [Very Large Values](#very-large-values)
  - [Only Additions](#only-additions)
  - [Only Removals](#only-removals)
- [🔧 Extension Points](#-extension-points)
  - [Adding Breakdown Display](#adding-breakdown-display)
  - [Threshold-Based Colors](#threshold-based-colors)
- [🔍 Troubleshooting](#-troubleshooting)
  - [TotalLines seems wrong](#totallines-seems-wrong)
  - [HasInfo false when expecting data](#hasinfo-false-when-expecting-data)
  - [Want net change, not total](#want-net-change-not-total)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetLinesModifiedDisplay** formats lines modified information (added + removed) from session data into display structure suitable for statusline.

> [!NOTE]
> **Pure presentation function.** Receives lines data from SessionContext.Cost, calculates total churn, outputs formatted display. No I/O operations, cannot fail.

**Design:** Pure presentation layer - receives lines data from SessionContext.Cost, calculates total editing activity (churn), outputs formatted display with visual indicators.

**Library:** [`statusline/lib/session`](../../../lib/session/)

**Source Code:** [`lib/session/session.go`](../../../lib/session/session.go)

**Library README:** [Session Library Documentation](README.md)

---

## 📝 Signature

```go
func GetLinesModifiedDisplay(ctx types.SessionContext) LinesModifiedDisplay
```

---

## ⚙ Parameters

### ctx (SessionContext)

Session context containing cost tracking data with lines modified information.

**Required Fields:**

| Field | Type | Purpose |
|-------|:----:|---------|
| **Cost.TotalLinesAdded** | `int` | Number of lines added during session |
| **Cost.TotalLinesRemoved** | `int` | Number of lines removed during session |

**Example:**

```go
ctx := types.SessionContext{
    Cost: types.CostInfo{
        TotalLinesAdded: 100,
        TotalLinesRemoved: 50,
    },
}
```

**Validation:** None required - negative values treated as zero (no changes)

---

## 🔄 Returns

### LinesModifiedDisplay

Formatted display structure with all presentation elements.

```go
type LinesModifiedDisplay struct {
    TotalLines int    // Total lines added + removed
    Color      string // Terminal color code for display
    Icon       string // Visual icon representing editing (e.g., "📝")
    HasInfo    bool   // True if changes tracked, false if no data
}
```

**Field Details:**

| Field | Type | Purpose | Example |
|-------|------|---------|---------|
| **TotalLines** | `int` | Total editing activity (churn) | `150` |
| **Color** | `string` | Terminal color | `display.Yellow` |
| **Icon** | `string` | Visual icon | `"📝"` |
| **HasInfo** | `bool` | Data available? | `true` / `false` |

**TotalLines Calculation:** `TotalLinesAdded + TotalLinesRemoved`

- Represents total editing activity (code churn), not net change
- Zero or negative → `HasInfo: false`

**Color:**

- Currently: `display.Yellow` (indicates active editing)
- Future: Threshold-based colors (high churn = red, normal = yellow, low = green)

**Icon:**

- Currently: `📝` (edit/pencil icon)
- Future: State-based icons

**HasInfo:**

- `true`: Changes tracked, TotalLines available
- `false`: No changes (zero or negative total)

**Guarantee:** Always returns valid LinesModifiedDisplay (never nil, never errors)

---

## 🎯 Behavior

### Calculation Logic

Total lines = added + removed:

```go
totalLinesModified := ctx.Cost.TotalLinesAdded + ctx.Cost.TotalLinesRemoved

if totalLinesModified <= 0 {
    return LinesModifiedDisplay{HasInfo: false}  // Zero-state
}

return LinesModifiedDisplay{
    TotalLines: totalLinesModified,
    Color:      display.Yellow,
    Icon:       "📝",
    HasInfo:    true,
}
```

### Display Examples

| Scenario | TotalLinesAdded | TotalLinesRemoved | TotalLines | HasInfo |
|----------|:---------------:|:-----------------:|:----------:|:-------:|
| Active Editing | 100 | 50 | 150 | `true` |
| No Changes | 0 | 0 | 0 | `false` |
| Only Additions | 75 | 0 | 75 | `true` |
| Only Removals | 0 | 30 | 30 | `true` |

### Semantic Meaning

**Total Churn:** Represents total editing activity, not net change

- **Added 100, Removed 50:** TotalLines = 150 (not 50 net)
- **Why:** Shows work volume, not just end result
- **Example:** Refactoring may have zero net lines but high churn

**When to Use:**

| ✅ Good For | ❌ Not Good For |
|------------|----------------|
| Showing session activity level | Net codebase growth (use added - removed) |
| Tracking editing intensity | Code quality metrics (just quantity) |
| Monitoring code churn | Productivity measurement (lines ≠ productivity) |

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/session"
import "statusline/lib/types"

func buildStatusline(ctx types.SessionContext) string {
    // Get lines modified display
    linesDisplay := session.GetLinesModifiedDisplay(ctx)
    // → TotalLines: 150 (100 added + 50 removed)

    var parts []string

    // Only add if changes tracked
    if linesDisplay.HasInfo {
        linesPart := fmt.Sprintf("%s %d lines", linesDisplay.Icon, linesDisplay.TotalLines)
        parts = append(parts, linesPart)
        // → "📝 150 lines"
    }

    // ... add other parts ...

    return strings.Join(parts, " | ")
}
```

### With Color

```go
linesDisplay := session.GetLinesModifiedDisplay(ctx)

if linesDisplay.HasInfo {
    // Apply color for terminal output
    colored := fmt.Sprintf("%s%s %d lines%s",
        linesDisplay.Color,      // Start color
        linesDisplay.Icon,        // Icon
        linesDisplay.TotalLines,  // Value
        display.Reset)            // Reset color
    fmt.Println(colored)
    // → Displays in yellow: 📝 150 lines
}
```

### Zero-State Handling

```go
emptyCtx := types.SessionContext{}
linesDisplay := session.GetLinesModifiedDisplay(emptyCtx)

// Check HasInfo before using
if linesDisplay.HasInfo {
    // Won't execute - no changes tracked
} else {
    // Correct behavior - no display needed
    fmt.Println("No changes yet")
}
```

### Full Statusline Integration

```go
func buildSessionStatusline(ctx types.SessionContext) string {
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

    return strings.Join(parts, " | ")
    // → "150 lines | 1h 23m | $0.0123"
}
```

### Testing

```go
func TestLinesModifiedDisplay(t *testing.T) {
    testCases := []struct {
        name     string
        added    int
        removed  int
        expected int
        hasInfo  bool
    }{
        {"Active editing", 100, 50, 150, true},
        {"No changes", 0, 0, 0, false},
        {"Only additions", 75, 0, 75, true},
        {"Only removals", 0, 30, 30, true},
        {"Negative (error)", -10, 0, 0, false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ctx := types.SessionContext{
                Cost: types.CostInfo{
                    TotalLinesAdded:   tc.added,
                    TotalLinesRemoved: tc.removed,
                },
            }

            result := session.GetLinesModifiedDisplay(ctx)

            if result.TotalLines != tc.expected {
                t.Errorf("Expected %d, got %d", tc.expected, result.TotalLines)
            }
            if result.HasInfo != tc.hasInfo {
                t.Errorf("Expected HasInfo %v, got %v", tc.hasInfo, result.HasInfo)
            }
        })
    }
}
```

---

## 🏥 Health Scoring

**No health scoring** - pure presentation library.

All operations guaranteed to succeed through graceful degradation:

- Zero lines modified → Returns `LinesModifiedDisplay{HasInfo: false}`
- Negative values (data error) → Treated as zero, returns `HasInfo: false`
- Very large values → No truncation, displays full value

**Philosophy:** This function cannot fail. Always returns valid LinesModifiedDisplay structures.

---

## ⚡ Performance

**Time Complexity:** O(1) - Simple addition operation

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Arithmetic** | 1 addition | TotalLinesAdded + TotalLinesRemoved |
| **Struct allocation** | ~50 bytes | LinesModifiedDisplay |
| **Execution time** | <1 μs | Microseconds |

**Memory:**

- Single LinesModifiedDisplay struct allocation (~50 bytes)
- No additional heap allocations
- Garbage collected automatically

**Optimization:** Not needed - arithmetic and struct allocation already optimal

---

## 🔧 Edge Cases

### Zero Total (No Changes)

```go
ctx := SessionContext{Cost: CostInfo{TotalLinesAdded: 0, TotalLinesRemoved: 0}}
display := GetLinesModifiedDisplay(ctx)
// → LinesModifiedDisplay{HasInfo: false}
// Behavior: Zero changes = no display needed
```

### Negative Values (Data Error)

```go
ctx := SessionContext{Cost: CostInfo{TotalLinesAdded: -10, TotalLinesRemoved: 0}}
display := GetLinesModifiedDisplay(ctx)
// → LinesModifiedDisplay{HasInfo: false}
// Behavior: Treats negative as zero (graceful degradation)
```

### Very Large Values

```go
ctx := SessionContext{Cost: CostInfo{TotalLinesAdded: 5000, TotalLinesRemoved: 3000}}
display := GetLinesModifiedDisplay(ctx)
// → LinesModifiedDisplay{TotalLines: 8000, HasInfo: true}
// Behavior: No truncation - displays full value
// Note: Future enhancement may format as "8K lines"
```

### Only Additions

```go
ctx := SessionContext{Cost: CostInfo{TotalLinesAdded: 100, TotalLinesRemoved: 0}}
display := GetLinesModifiedDisplay(ctx)
// → LinesModifiedDisplay{TotalLines: 100, HasInfo: true}
// Behavior: Total churn includes additions only
```

### Only Removals

```go
ctx := SessionContext{Cost: CostInfo{TotalLinesAdded: 0, TotalLinesRemoved: 50}}
display := GetLinesModifiedDisplay(ctx)
// → LinesModifiedDisplay{TotalLines: 50, HasInfo: true}
// Behavior: Total churn includes removals only
```

---

## 🔧 Extension Points

### Adding Breakdown Display

Show added vs removed separately:

**Pattern:**

```go
// New display type
type LinesBreakdownDisplay struct {
    Added    int
    Removed  int
    Color    string
    Icon     string
    HasInfo  bool
}

func GetLinesBreakdownDisplay(ctx types.SessionContext) LinesBreakdownDisplay {
    if ctx.Cost.TotalLinesAdded <= 0 && ctx.Cost.TotalLinesRemoved <= 0 {
        return LinesBreakdownDisplay{HasInfo: false}
    }

    return LinesBreakdownDisplay{
        Added:   ctx.Cost.TotalLinesAdded,
        Removed: ctx.Cost.TotalLinesRemoved,
        Color:   display.Yellow,
        Icon:    "📝",
        HasInfo: true,
    }
}
// Display: "+100 -50"
```

### Threshold-Based Colors

**Pattern:**

```go
func GetLinesModifiedDisplay(ctx types.SessionContext) LinesModifiedDisplay {
    totalLinesModified := ctx.Cost.TotalLinesAdded + ctx.Cost.TotalLinesRemoved

    if totalLinesModified <= 0 {
        return LinesModifiedDisplay{HasInfo: false}
    }

    // Threshold-based color selection
    color := display.Green   // Low churn (< 100)
    if totalLinesModified > 500 {
        color = display.Red  // High churn
    } else if totalLinesModified > 100 {
        color = display.Yellow  // Normal churn
    }

    return LinesModifiedDisplay{
        TotalLines: totalLinesModified,
        Color:      color,  // State-based
        Icon:       "📝",
        HasInfo:    true,
    }
}
```

---

## 🔍 Troubleshooting

### TotalLines seems wrong

**Problem:** Display shows unexpected line count

**Check:**

1. Verify `TotalLinesAdded + TotalLinesRemoved = expected total`
2. Confirm understanding: Total = added + removed (not net change)
3. Check SessionContext population - where do these values come from?

**Example:**

```go
// Added 100, Removed 50
// TotalLines = 150 (not 50)
// This is CORRECT - shows total churn, not net change
```

### HasInfo false when expecting data

**Problem:** HasInfo returns false despite having session activity

**Check:**

1. Verify TotalLinesAdded and TotalLinesRemoved are > 0
2. Check if values are negative (gracefully treated as zero)
3. Confirm SessionContext.Cost properly populated

**Diagnosis:**

```go
fmt.Printf("Added: %d, Removed: %d, Total: %d\n",
    ctx.Cost.TotalLinesAdded,
    ctx.Cost.TotalLinesRemoved,
    ctx.Cost.TotalLinesAdded + ctx.Cost.TotalLinesRemoved)
// If total <= 0, HasInfo will be false (correct behavior)
```

### Want net change, not total

**Problem:** Need net lines (added - removed), not total churn

**Solution:** This function shows churn. For net change, calculate separately:

```go
netLines := ctx.Cost.TotalLinesAdded - ctx.Cost.TotalLinesRemoved
// Net = 100 - 50 = 50 (growth)

totalChurn := linesDisplay.TotalLines
// Total = 100 + 50 = 150 (activity)
```

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetDurationDisplay** | Session duration formatting | [duration.md](duration.md) |
| **GetCostDisplay** | Session cost formatting | [cost.md](cost.md) |

**Future Functions:**

| Function | Purpose | Status |
|----------|---------|:------:|
| `GetLinesBreakdownDisplay()` | Separate added/removed display | ⏳ Planned |
| `GetNetLinesDisplay()` | Net change (added - removed) | ⏳ Planned |
| `GetFilesChangedDisplay()` | File count instead of line count | 🔬 Research |

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Session Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetDurationDisplay** | Duration formatting | [duration.md](duration.md) |
| **GetCostDisplay** | Cost formatting | [cost.md](cost.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/session/session.go`](../../../lib/session/session.go) | GetLinesModifiedDisplay with comprehensive inline documentation |
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

- **"So teach us to number our days, that we may apply our hearts unto wisdom"** - Psalm 90:12
  - Awareness of work to foster wisdom - tracking lines changed to understand editing patterns and session activity

**Application:** Show what changed (lines modified) without judgment - the number itself doesn't measure quality, but awareness of activity fosters wise reflection on how time was used.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
