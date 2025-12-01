<div align="center">

# GetCostDisplay & GetFormattedCost Reference

**Session Cost Tracking and USD Formatting**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Functions](https://img.shields.io/badge/Functions-2-blue?style=flat)
![Precision](https://img.shields.io/badge/Precision-4%20Decimals-orange?style=flat)

*Formats session cost data and provides standardized USD string formatting for statusline*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signatures](#-signatures) • [⚙ Parameters](#-parameters) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetCostDisplay & GetFormattedCost Functions
**Purpose:** Session cost presentation and USD formatting
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

- [📖 Overview](#-overview)
- [📝 Signatures](#-signatures)
  - [GetCostDisplay](#getcostdisplay)
  - [GetFormattedCost](#getformattedcost)
- [⚙ Parameters](#-parameters)
  - [GetCostDisplay Parameters](#getcostdisplay-parameters)
    - [ctx (SessionContext)](#ctx-sessioncontext)
  - [GetFormattedCost Parameters](#getformattedcost-parameters)
    - [cost (float64)](#cost-float64)
- [🔄 Returns](#-returns)
  - [CostDisplay (GetCostDisplay)](#costdisplay-getcostdisplay)
  - [string (GetFormattedCost)](#string-getformattedcost)
- [🎯 Behavior](#-behavior)
  - [GetCostDisplay Behavior](#getcostdisplay-behavior)
  - [GetFormattedCost Behavior](#getformattedcost-behavior)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [With Color](#with-color)
  - [Zero-State Handling](#zero-state-handling)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Direct Formatting](#direct-formatting)
  - [Testing](#testing)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [🔧 Edge Cases](#-edge-cases)
  - [Zero Cost](#zero-cost)
  - [Negative Cost (Data Error)](#negative-cost-data-error)
  - [Very Small Cost (Micro-transactions)](#very-small-cost-micro-transactions)
  - [Large Cost](#large-cost)
  - [Precision Edge](#precision-edge)
  - [Exact Dollar Amount](#exact-dollar-amount)
- [🔧 Extension Points](#-extension-points)
  - [Threshold-Based Colors](#threshold-based-colors)
  - [Configurable Precision](#configurable-precision)
  - [Human-Readable Large Costs](#human-readable-large-costs)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Cost precision looks wrong](#cost-precision-looks-wrong)
  - [HasInfo false when expecting cost](#hasinfo-false-when-expecting-cost)
  - [Want different precision (not 4 decimals)](#want-different-precision-not-4-decimals)
  - [Large costs look odd with 4 decimals](#large-costs-look-odd-with-4-decimals)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetCostDisplay** formats session cost tracking data into display structure suitable for statusline. **GetFormattedCost** provides consistent USD string formatting (4 decimal precision) across statusline displays.

> [!NOTE]
> **Two complementary functions.** GetCostDisplay extracts cost data, GetFormattedCost formats it for display. Pure presentation functions, no I/O operations, cannot fail.

**Design:** Pure presentation layer - GetCostDisplay receives cost data from SessionContext.Cost, outputs formatted display structure. GetFormattedCost converts float64 to standardized USD strings with 4 decimal precision.

**Library:** [`statusline/lib/session`](../../../lib/session/)

**Source Code:** [`lib/session/session.go`](../../../lib/session/session.go)

**Library README:** [Session Library Documentation](README.md)

---

## 📝 Signatures

### GetCostDisplay

```go
func GetCostDisplay(ctx types.SessionContext) CostDisplay
```

Extracts cost data from SessionContext and packages for statusline display.

### GetFormattedCost

```go
func GetFormattedCost(cost float64) string
```

Helper function to format cost as USD string with 4 decimal precision.

---

## ⚙ Parameters

### GetCostDisplay Parameters

#### ctx (SessionContext)

Session context containing cost tracking data.

**Required Fields:**

| Field | Type | Purpose |
|-------|:----:|---------|
| **Cost.TotalCostUSD** | `float64` | Session cost in USD |

**Example:**

```go
ctx := types.SessionContext{
    Cost: types.CostInfo{
        TotalCostUSD: 0.0123,  // $0.0123 USD
    },
}
```

**Validation:** None required - zero or negative values treated as no cost

### GetFormattedCost Parameters

#### cost (float64)

Raw USD cost value.

**Range:** Any float64 value (typically 0.0001 to 100.0 for sessions)

**Example:**

```go
cost := 0.0123  // $0.0123 USD
```

---

## 🔄 Returns

### CostDisplay (GetCostDisplay)

Formatted display structure with all presentation elements.

```go
type CostDisplay struct {
    Cost    float64 // Cost in USD
    Color   string  // Terminal color code for display
    Icon    string  // Visual icon representing cost (e.g., "💰")
    HasInfo bool    // True if cost tracked, false if no data
}
```

**Field Details:**

| Field | Type | Purpose | Example |
|-------|------|---------|---------|
| **Cost** | `float64` | Raw USD value | `0.0123` |
| **Color** | `string` | Terminal color | `display.Yellow` |
| **Icon** | `string` | Visual icon | `"💰"` |
| **HasInfo** | `bool` | Data available? | `true` / `false` |

**Cost Field:** Raw USD value from SessionContext (e.g., 0.0123) - Use GetFormattedCost() to convert to display string

**Color:**

- Currently: `display.Yellow` (indicates cost awareness)
- Future: Threshold-based colors (high cost = red, normal = yellow, low = green)

**Icon:**

- Currently: `💰` (money icon)
- Future: State-based icons

**HasInfo:**

- `true`: Cost tracked, Cost value available
- `false`: No cost (zero or negative USD)

**Guarantee:** Always returns valid CostDisplay (never nil, never errors)

### string (GetFormattedCost)

Formatted cost string with 4 decimal precision.

**Format:** `$%.4f` (always 4 decimal places)

**Examples:**

| Input | Output | Notes |
|:-----:|--------|-------|
| `0.0123` | `"$0.0123"` | Typical token cost |
| `1.5` | `"$1.5000"` | Trailing zeros added |
| `0.0001` | `"$0.0001"` | Micro-cost |
| `100.0` | `"$100.0000"` | Large cost (looks odd) |
| `0.0` | `"$0.0000"` | Zero cost |

---

## 🎯 Behavior

### GetCostDisplay Behavior

USD value → display structure → formatted string:

```go
if ctx.Cost.TotalCostUSD <= 0 {
    return CostDisplay{HasInfo: false}  // Zero-state
}

return CostDisplay{
    Cost:    ctx.Cost.TotalCostUSD,
    Color:   display.Yellow,
    Icon:    "💰",
    HasInfo: true,
}
```

**Display Examples:**

| Scenario | TotalCostUSD | Cost Field | HasInfo | Formatted |
|----------|:------------:|:----------:|:-------:|-----------|
| Active session | 0.0123 | 0.0123 | `true` | `"$0.0123"` |
| No cost | 0.0 | 0.0 | `false` | N/A |
| High cost | 1.5678 | 1.5678 | `true` | `"$1.5678"` |
| Micro cost | 0.0001 | 0.0001 | `true` | `"$0.0001"` |

### GetFormattedCost Behavior

Converts float64 to standardized USD string format:

```go
return fmt.Sprintf("$%.4f", cost)
```

**Why 4 Decimals:**

- Token costs often very small (e.g., $0.0001 per request)
- 2 decimals would round to $0.00 (not useful)
- 4 decimals provides meaningful precision for micro-costs
- Standardized across all cost displays

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/session"
import "statusline/lib/types"

func buildStatusline(ctx types.SessionContext) string {
    // Get cost display
    costDisplay := session.GetCostDisplay(ctx)
    // → Cost: 0.0123

    var parts []string

    // Only add if cost tracked
    if costDisplay.HasInfo {
        costPart := fmt.Sprintf("%s %s", costDisplay.Icon, session.GetFormattedCost(costDisplay.Cost))
        parts = append(parts, costPart)
        // → "💰 $0.0123"
    }

    // ... add other parts ...

    return strings.Join(parts, " | ")
}
```

### With Color

```go
costDisplay := session.GetCostDisplay(ctx)

if costDisplay.HasInfo {
    // Apply color for terminal output
    colored := fmt.Sprintf("%s%s %s%s",
        costDisplay.Color,                        // Start color
        costDisplay.Icon,                         // Icon
        session.GetFormattedCost(costDisplay.Cost), // Formatted cost
        display.Reset)                            // Reset color
    fmt.Println(colored)
    // → Displays in yellow: 💰 $0.0123
}
```

### Zero-State Handling

```go
emptyCtx := types.SessionContext{}
costDisplay := session.GetCostDisplay(emptyCtx)

// Check HasInfo before using
if costDisplay.HasInfo {
    // Won't execute - no cost tracked
} else {
    // Correct behavior - no display needed
    fmt.Println("Session has no cost data")
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

### Direct Formatting

```go
// Using GetFormattedCost directly (without GetCostDisplay)
rawCost := 0.0123
formattedCost := session.GetFormattedCost(rawCost)
fmt.Println(formattedCost)  // → "$0.0123"

// Array of costs
costs := []float64{0.0001, 0.0050, 0.0123, 1.5000}
for _, cost := range costs {
    fmt.Println(session.GetFormattedCost(cost))
}
// → "$0.0001"
// → "$0.0050"
// → "$0.0123"
// → "$1.5000"
```

### Testing

```go
func TestCostDisplay(t *testing.T) {
    testCases := []struct {
        name        string
        costUSD     float64
        hasInfo     bool
        formatted   string
    }{
        {"Active session", 0.0123, true, "$0.0123"},
        {"No cost", 0.0, false, "$0.0000"},
        {"High cost", 1.5678, true, "$1.5678"},
        {"Micro cost", 0.0001, true, "$0.0001"},
        {"Negative (error)", -0.01, false, "-$0.0100"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ctx := types.SessionContext{
                Cost: types.CostInfo{
                    TotalCostUSD: tc.costUSD,
                },
            }

            result := session.GetCostDisplay(ctx)
            formatted := session.GetFormattedCost(tc.costUSD)

            if result.HasInfo != tc.hasInfo {
                t.Errorf("Expected HasInfo %v, got %v", tc.hasInfo, result.HasInfo)
            }
            if formatted != tc.formatted {
                t.Errorf("Expected %q, got %q", tc.formatted, formatted)
            }
        })
    }
}
```

---

## 🏥 Health Scoring

**No health scoring** - pure presentation library.

All operations guaranteed to succeed through graceful degradation:

- Zero cost → Returns `CostDisplay{HasInfo: false}`
- Negative cost (data error) → Treated as zero, returns `HasInfo: false}`
- Very small / very large costs → Formatted with consistent 4 decimal precision

**Philosophy:** These functions cannot fail. Always return valid CostDisplay structures or formatted strings.

---

## ⚡ Performance

**Time Complexity:** O(1) - Simple operations

**Typical Performance:**

| Function | Operation | Memory | Time |
|----------|-----------|:------:|:----:|
| **GetCostDisplay** | Field extraction + struct allocation | ~50 bytes | <1μs |
| **GetFormattedCost** | Sprintf with float formatting | ~15 bytes | <1μs |

**Memory:**

- **GetCostDisplay:** Single CostDisplay struct allocation (~50 bytes)
- **GetFormattedCost:** String allocation (~15 bytes for typical costs)
- Garbage collected automatically

**Optimization:** Not needed - operations already optimal

---

## 🔧 Edge Cases

### Zero Cost

```go
ctx := SessionContext{Cost: CostInfo{TotalCostUSD: 0.0}}
display := GetCostDisplay(ctx)
// → CostDisplay{HasInfo: false}
// Behavior: No cost = no display needed
```

### Negative Cost (Data Error)

```go
ctx := SessionContext{Cost: CostInfo{TotalCostUSD: -0.01}}
display := GetCostDisplay(ctx)
// → CostDisplay{HasInfo: false}
// Behavior: Treats negative as zero (graceful degradation)

formatted := GetFormattedCost(-0.01)
// → "-$0.0100"
// Note: GetFormattedCost doesn't validate - formats any float64
```

### Very Small Cost (Micro-transactions)

```go
ctx := SessionContext{Cost: CostInfo{TotalCostUSD: 0.0001}}
display := GetCostDisplay(ctx)
// → CostDisplay{Cost: 0.0001, HasInfo: true}
formatted := GetFormattedCost(0.0001)
// → "$0.0001"
// Behavior: 4 decimals captures micro-costs
```

### Large Cost

```go
ctx := SessionContext{Cost: CostInfo{TotalCostUSD: 100.0}}
display := GetCostDisplay(ctx)
// → CostDisplay{Cost: 100.0, HasInfo: true}
formatted := GetFormattedCost(100.0)
// → "$100.0000"
// Note: Looks odd with trailing zeros - consider configurable precision
```

### Precision Edge

```go
formatted := GetFormattedCost(0.00005)
// → "$0.0001" (rounded to 4 decimals)
// Note: Values beyond 4 decimals will be rounded

formatted = GetFormattedCost(0.123456789)
// → "$0.1235" (rounded to 4 decimals)
```

### Exact Dollar Amount

```go
formatted := GetFormattedCost(1.0)
// → "$1.0000"
// Behavior: Always shows 4 decimals, even for whole dollars
```

---

## 🔧 Extension Points

### Threshold-Based Colors

**Pattern:**

```go
func GetCostDisplay(ctx types.SessionContext) CostDisplay {
    if ctx.Cost.TotalCostUSD <= 0 {
        return CostDisplay{HasInfo: false}
    }

    // Threshold-based color selection
    color := display.Green   // Low cost (< $0.01)
    if ctx.Cost.TotalCostUSD > 1.0 {
        color = display.Red  // High cost
    } else if ctx.Cost.TotalCostUSD > 0.1 {
        color = display.Yellow  // Normal cost
    }

    return CostDisplay{
        Cost:    ctx.Cost.TotalCostUSD,
        Color:   color,  // State-based
        Icon:    "💰",
        HasInfo: true,
    }
}
```

### Configurable Precision

**Pattern:**

```go
// Different precision levels
func GetFormattedCostPrecise(cost float64) string {
    return fmt.Sprintf("$%.6f", cost)  // 6 decimals for very small costs
}

func GetFormattedCostStandard(cost float64) string {
    return fmt.Sprintf("$%.2f", cost)  // 2 decimals for typical costs
}

func GetFormattedCostAuto(cost float64) string {
    // Choose precision based on cost magnitude
    if cost < 0.01 {
        return fmt.Sprintf("$%.4f", cost)  // 4 decimals for micro-costs
    }
    return fmt.Sprintf("$%.2f", cost)  // 2 decimals otherwise
}
```

### Human-Readable Large Costs

**Pattern:**

```go
func GetFormattedCostReadable(cost float64) string {
    if cost >= 100 {
        return fmt.Sprintf("$%.2f", cost)  // Drop trailing zeros for large costs
    }
    return fmt.Sprintf("$%.4f", cost)  // Keep precision for small costs
}
// Examples:
// $0.0123 → "$0.0123" (small, keep precision)
// $100.00 → "$100.00" (large, use 2 decimals)
```

---

## 🔍 Troubleshooting

### Cost precision looks wrong

**Problem:** Display shows unexpected number of decimal places

**Check:**

1. Verify using GetFormattedCost() (always 4 decimals)
2. Confirm not using custom formatting elsewhere
3. Check if rounding occurred due to float64 precision limits

**Expected:** Always exactly 4 decimal places (`$0.0123`, `$1.5000`, etc.)

### HasInfo false when expecting cost

**Problem:** HasInfo returns false despite session having cost

**Check:**

1. Verify TotalCostUSD is > 0
2. Check if value is negative (gracefully treated as zero)
3. Confirm SessionContext.Cost properly populated

**Expected:** Zero or negative cost = `HasInfo: false` (correct behavior)

### Want different precision (not 4 decimals)

**Problem:** Need 2 or 6 decimals instead of 4

**Solution:** Create custom formatting function (see Extension Points above)

```go
// Custom 2-decimal formatting
func formatCost2Decimals(cost float64) string {
    return fmt.Sprintf("$%.2f", cost)
}
```

### Large costs look odd with 4 decimals

**Problem:** `$100.0000` looks awkward

**Solution:** Use adaptive precision based on cost magnitude (see Extension Points above)

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetLinesModifiedDisplay** | Lines modified formatting | [lines.md](lines.md) |
| **GetDurationDisplay** | Session duration formatting | [duration.md](duration.md) |

**Future Functions:**

| Function | Purpose | Status |
|----------|---------|:------:|
| `GetFormattedCostPrecise()` | 6 decimal precision | ⏳ Planned |
| `GetFormattedCostStandard()` | 2 decimal precision | ⏳ Planned |
| `GetFormattedCostAuto()` | Adaptive precision | 🔬 Research |
| `GetCostBreakdownDisplay()` | Input/output cost breakdown | 🔬 Research |

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Session Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetLinesModifiedDisplay** | Lines modified formatting | [lines.md](lines.md) |
| **GetDurationDisplay** | Duration formatting | [duration.md](duration.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/session/session.go`](../../../lib/session/session.go) | GetCostDisplay & GetFormattedCost with comprehensive inline documentation |
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
  - Awareness of resources to foster wisdom - tracking cost to understand resource usage

**Supporting Verses:**

- **"For which of you, intending to build a tower, sitteth not down first, and counteth the cost..."** - Luke 14:28
  - Count the cost before undertaking work
- **"Owe no man any thing, but to love one another"** - Romans 13:8
  - Responsible stewardship of resources

**Application:** Show session cost ($0.0123) to create awareness of resource usage without judgment - the cost itself doesn't measure value, but awareness of cost fosters wise reflection on resource stewardship.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
