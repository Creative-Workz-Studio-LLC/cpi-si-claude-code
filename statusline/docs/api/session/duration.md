<div align="center">

# GetDurationDisplay Function Reference

**Session Duration Tracking and Formatting**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Format](https://img.shields.io/badge/Format-Human%20Readable-green?style=flat)

*Formats session duration from milliseconds into human-readable display for statusline*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetDurationDisplay Function
**Purpose:** Session duration presentation for session statistics
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
  - [DurationDisplay](#durationdisplay)
- [🎯 Behavior](#-behavior)
  - [Conversion Process](#conversion-process)
  - [Duration Format Examples](#duration-format-examples)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [With Color](#with-color)
  - [Zero-State Handling](#zero-state-handling)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Testing](#testing)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [🔧 Edge Cases](#-edge-cases)
  - [Zero Duration](#zero-duration)
  - [Negative Duration (Data Error)](#negative-duration-data-error)
  - [Very Short Duration (< 1 second)](#very-short-duration--1-second)
  - [Very Long Duration (Multi-Day)](#very-long-duration-multi-day)
  - [Exact Hour](#exact-hour)
- [🔧 Extension Points](#-extension-points)
  - [Threshold-Based Colors](#threshold-based-colors)
  - [Alternative Format Functions](#alternative-format-functions)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Duration format unexpected](#duration-format-unexpected)
  - [HasInfo false when expecting duration](#hasinfo-false-when-expecting-duration)
  - [Duration string empty but HasInfo true](#duration-string-empty-but-hasinfo-true)
  - [Want millisecond precision](#want-millisecond-precision)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetDurationDisplay** formats session duration from milliseconds into human-readable display structure suitable for statusline.

> [!NOTE]
> **Pure presentation function.** Receives duration data from SessionContext.Cost, converts to human-readable format using system/lib/sessiontime. No I/O operations, cannot fail.

**Design:** Pure presentation layer - receives duration data from SessionContext.Cost, converts to human-readable format using system/lib/sessiontime, outputs formatted display with visual indicators.

**Library:** [`statusline/lib/session`](../../../lib/session/)

**Source Code:** [`lib/session/session.go`](../../../lib/session/session.go)

**Library README:** [Session Library Documentation](README.md)

---

## 📝 Signature

```go
func GetDurationDisplay(ctx types.SessionContext) DurationDisplay
```

---

## ⚙ Parameters

### ctx (SessionContext)

Session context containing cost tracking data with duration information.

**Required Fields:**

| Field | Type | Purpose |
|-------|:----:|---------|
| **Cost.TotalDurationMS** | `int64` | Session duration in milliseconds |

**Example:**

```go
ctx := types.SessionContext{
    Cost: types.CostInfo{
        TotalDurationMS: 5400000,  // 1 hour 30 minutes = 90 minutes = 5,400,000 ms
    },
}
```

**Validation:** None required - zero or negative values treated as no duration

---

## 🔄 Returns

### DurationDisplay

Formatted display structure with all presentation elements.

```go
type DurationDisplay struct {
    Duration string // Human-readable duration (e.g., "1h 23m", "45s")
    Color    string // Terminal color code for display
    Icon     string // Visual icon representing time (e.g., "⏱")
    HasInfo  bool   // True if duration tracked, false if no data
}
```

**Field Details:**

| Field | Type | Purpose | Example |
|-------|------|---------|---------|
| **Duration** | `string` | Human-readable format | `"1h 23m"` |
| **Color** | `string` | Terminal color | `display.Gray` |
| **Icon** | `string` | Visual icon | `"⏱"` |
| **HasInfo** | `bool` | Data available? | `true` / `false` |

**Duration Format:** Human-readable from `system/lib/sessiontime.FormatDuration()`

- Examples: `"45s"`, `"1m 23s"`, `"1h 23m"`, `"2h 15m 30s"`
- Empty string if no duration tracked

**Color:**

- Currently: `display.Gray` (passive time tracking indicator)
- Future: State-based colors (long sessions = yellow, very long = red)

**Icon:**

- Currently: `⏱` (timer/stopwatch icon)
- Future: State-based icons

**HasInfo:**

- `true`: Duration tracked, Duration string available
- `false`: No duration (zero or negative milliseconds)

**Guarantee:** Always returns valid DurationDisplay (never nil, never errors)

---

## 🎯 Behavior

### Conversion Process

Milliseconds → time.Duration → human-readable string:

```go
if ctx.Cost.TotalDurationMS <= 0 {
    return DurationDisplay{HasInfo: false}  // Zero-state
}

// Convert ms to time.Duration
duration := sessiontime.FormatDuration(time.Duration(ctx.Cost.TotalDurationMS) * time.Millisecond)

return DurationDisplay{
    Duration: duration,
    Color:    display.Gray,
    Icon:     "⏱",
    HasInfo:  true,
}
```

### Duration Format Examples

| Duration MS | Seconds | Minutes | Hours | Format | Notes |
|:----------:|:-------:|:-------:|:-----:|--------|-------|
| 45000 | 45 | - | - | `"45s"` | Short session |
| 83000 | 83 | 1.38 | - | `"1m 23s"` | Medium session |
| 4980000 | 4980 | 83 | 1.38 | `"1h 23m"` | Long session (seconds dropped) |
| 8130000 | 8130 | 135.5 | 2.26 | `"2h 15m 30s"` | Very long (seconds included) |
| 0 | 0 | 0 | 0 | `""` | No duration → `HasInfo: false` |

**Note:** Actual formatting rules depend on `system/lib/sessiontime` implementation. See system lib documentation for precise behavior.

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/session"
import "statusline/lib/types"

func buildStatusline(ctx types.SessionContext) string {
    // Get duration display
    durationDisplay := session.GetDurationDisplay(ctx)
    // → Duration: "1h 23m"

    var parts []string

    // Only add if duration tracked
    if durationDisplay.HasInfo {
        durationPart := fmt.Sprintf("%s %s", durationDisplay.Icon, durationDisplay.Duration)
        parts = append(parts, durationPart)
        // → "⏱ 1h 23m"
    }

    // ... add other parts ...

    return strings.Join(parts, " | ")
}
```

### With Color

```go
durationDisplay := session.GetDurationDisplay(ctx)

if durationDisplay.HasInfo {
    // Apply color for terminal output
    colored := fmt.Sprintf("%s%s %s%s",
        durationDisplay.Color,     // Start color
        durationDisplay.Icon,       // Icon
        durationDisplay.Duration,   // Formatted time
        display.Reset)              // Reset color
    fmt.Println(colored)
    // → Displays in gray: ⏱ 1h 23m
}
```

### Zero-State Handling

```go
emptyCtx := types.SessionContext{}
durationDisplay := session.GetDurationDisplay(emptyCtx)

// Check HasInfo before using
if durationDisplay.HasInfo {
    // Won't execute - no duration tracked
} else {
    // Correct behavior - no display needed
    fmt.Println("Session just started")
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
func TestDurationDisplay(t *testing.T) {
    testCases := []struct {
        name       string
        durationMS int64
        expected   string
        hasInfo    bool
    }{
        {"Short session", 45000, "45s", true},
        {"Medium session", 83000, "1m 23s", true},
        {"Long session", 4980000, "1h 23m", true},
        {"No duration", 0, "", false},
        {"Negative (error)", -1000, "", false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ctx := types.SessionContext{
                Cost: types.CostInfo{
                    TotalDurationMS: tc.durationMS,
                },
            }

            result := session.GetDurationDisplay(ctx)

            // Note: Duration format depends on system/lib/sessiontime
            // Test checks HasInfo and non-empty duration, not exact format
            if result.HasInfo != tc.hasInfo {
                t.Errorf("Expected HasInfo %v, got %v", tc.hasInfo, result.HasInfo)
            }
            if tc.hasInfo && result.Duration == "" {
                t.Error("Expected non-empty duration string")
            }
        })
    }
}
```

---

## 🏥 Health Scoring

**No health scoring** - pure presentation library.

All operations guaranteed to succeed through graceful degradation:

- Zero duration → Returns `DurationDisplay{HasInfo: false}`
- Negative duration (data error) → Treated as zero, returns `HasInfo: false}`
- Very short (<1s) / very long (multi-day) → Formatted according to system/lib/sessiontime rules

**Philosophy:** This function cannot fail. Always returns valid DurationDisplay structures.

---

## ⚡ Performance

**Time Complexity:** O(1) - Conversion and formatting operations

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Millisecond conversion** | O(1) | Arithmetic operation |
| **FormatDuration call** | O(1) | String building |
| **Struct allocation** | ~50 bytes | DurationDisplay |
| **Execution time** | <10 μs | Microseconds |

**Memory:**

- Single DurationDisplay struct allocation (~50 bytes)
- Duration string allocation (~10-20 bytes depending on length)
- Garbage collected automatically

**Optimization:** Not needed - conversion and string formatting already optimal

---

## 🔧 Edge Cases

### Zero Duration

```go
ctx := SessionContext{Cost: CostInfo{TotalDurationMS: 0}}
display := GetDurationDisplay(ctx)
// → DurationDisplay{HasInfo: false}
// Behavior: No duration = no display needed
```

### Negative Duration (Data Error)

```go
ctx := SessionContext{Cost: CostInfo{TotalDurationMS: -1000}}
display := GetDurationDisplay(ctx)
// → DurationDisplay{HasInfo: false}
// Behavior: Treats negative as zero (graceful degradation)
```

### Very Short Duration (< 1 second)

```go
ctx := SessionContext{Cost: CostInfo{TotalDurationMS: 500}}  // 500ms
display := GetDurationDisplay(ctx)
// → DurationDisplay{Duration: "0s" or "1s", HasInfo: true}
// Behavior: Depends on system/lib/sessiontime rounding behavior
```

### Very Long Duration (Multi-Day)

```go
ctx := SessionContext{Cost: CostInfo{TotalDurationMS: 86400000}}  // 24 hours
display := GetDurationDisplay(ctx)
// → DurationDisplay{Duration: "24h 0m" or "1d", HasInfo: true}
// Behavior: Depends on system/lib/sessiontime format (may show days)
```

### Exact Hour

```go
ctx := SessionContext{Cost: CostInfo{TotalDurationMS: 3600000}}  // Exactly 1 hour
display := GetDurationDisplay(ctx)
// → DurationDisplay{Duration: "1h" or "1h 0m", HasInfo: true}
// Behavior: Depends on system/lib/sessiontime formatting rules
```

---

## 🔧 Extension Points

### Threshold-Based Colors

**Pattern:**

```go
func GetDurationDisplay(ctx types.SessionContext) DurationDisplay {
    if ctx.Cost.TotalDurationMS <= 0 {
        return DurationDisplay{HasInfo: false}
    }

    duration := sessiontime.FormatDuration(time.Duration(ctx.Cost.TotalDurationMS) * time.Millisecond)

    // Threshold-based color selection
    color := display.Gray  // Default (< 1 hour)
    hours := float64(ctx.Cost.TotalDurationMS) / 3600000
    if hours > 4 {
        color = display.Red  // Very long session
    } else if hours > 2 {
        color = display.Yellow  // Long session
    }

    return DurationDisplay{
        Duration: duration,
        Color:    color,  // State-based
        Icon:     "⏱",
        HasInfo:  true,
    }
}
```

### Alternative Format Functions

Different verbosity levels:

**Pattern:**

```go
// Compact format
func GetDurationDisplayCompact(ctx) DurationDisplay {
    // Returns: "1.5h" instead of "1h 30m"
}

// Verbose format
func GetDurationDisplayVerbose(ctx) DurationDisplay {
    // Returns: "1 hour, 30 minutes" instead of "1h 30m"
}

// Precise format
func GetDurationDisplayPrecise(ctx) DurationDisplay {
    // Returns: "1h 30m 45s" (always include seconds)
}
```

---

## 🔍 Troubleshooting

### Duration format unexpected

**Problem:** Display shows unexpected duration format

**Check:**

1. Verify TotalDurationMS value (milliseconds, not seconds)
2. Confirm understanding of `system/lib/sessiontime.FormatDuration()` behavior
3. Check system lib documentation for formatting rules

**Diagnosis:**

```go
ms := ctx.Cost.TotalDurationMS
fmt.Printf("MS: %d, Seconds: %.1f, Minutes: %.1f\n",
    ms, float64(ms)/1000, float64(ms)/60000)
// Verify millisecond value is correct
```

### HasInfo false when expecting duration

**Problem:** HasInfo returns false despite session having duration

**Check:**

1. Verify TotalDurationMS is > 0
2. Check if value is negative (gracefully treated as zero)
3. Confirm SessionContext.Cost properly populated

**Expected:** Zero or negative duration = `HasInfo: false` (correct behavior)

### Duration string empty but HasInfo true

**Problem:** HasInfo true but Duration field is empty string

**Cause:** This should not happen - indicates `system/lib/sessiontime.FormatDuration()` issue

**Solution:**

1. Check `system/lib/sessiontime` implementation
2. Verify FormatDuration never returns empty string for valid duration
3. File bug if system lib behavior incorrect

### Want millisecond precision

**Problem:** Need to show milliseconds, not just seconds

**Solution:** This function shows human-readable time (45s, not 45,123ms). For millisecond precision:

```go
// Custom formatting
ms := ctx.Cost.TotalDurationMS
msDisplay := fmt.Sprintf("%dms", ms)  // → "45123ms"

// Or with conversion
seconds := float64(ms) / 1000
secDisplay := fmt.Sprintf("%.3fs", seconds)  // → "45.123s"
```

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetLinesModifiedDisplay** | Lines modified formatting | [lines.md](lines.md) |
| **GetCostDisplay** | Session cost formatting | [cost.md](cost.md) |

**Dependencies:**

| Function | Purpose | Location |
|----------|---------|----------|
| **sessiontime.FormatDuration** | Duration formatter (lower rung) | [`system/lib/sessiontime`](~/.claude/cpi-si/system/lib/sessiontime/) |

**Future Functions:**

| Function | Purpose | Status |
|----------|---------|:------:|
| `GetDurationDisplayCompact()` | Abbreviated format (1.5h) | ⏳ Planned |
| `GetDurationDisplayVerbose()` | Full words (1 hour, 30 minutes) | ⏳ Planned |
| `GetDurationBreakdownDisplay()` | Separate hours/minutes/seconds fields | 🔬 Research |

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Session Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetLinesModifiedDisplay** | Lines modified formatting | [lines.md](lines.md) |
| **GetCostDisplay** | Cost formatting | [cost.md](cost.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/session/session.go`](../../../lib/session/session.go) | GetDurationDisplay with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |
| **Duration Formatting** | [`system/lib/sessiontime`](~/.claude/cpi-si/system/lib/sessiontime/) | Duration formatting utility (lower rung) |

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
  - Awareness of time to foster wisdom - tracking session duration to understand how time is spent

**Supporting Verse:**

- **"Redeeming the time, because the days are evil"** - Ephesians 5:16
  - Wise stewardship of time

**Application:** Show how long the session has been active (1h 23m) to create awareness without judgment - the time itself doesn't measure productivity, but awareness of duration fosters wise reflection on time stewardship.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
