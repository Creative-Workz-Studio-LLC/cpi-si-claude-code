<div align="center">

# GetSessionPhaseDisplay Function Reference

**Session Phase Display Formatting for Statusline**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Pure](https://img.shields.io/badge/Pure-Function-success?style=flat)
![Phase Colors](https://img.shields.io/badge/Phase%20Colors-3-purple?style=flat)

*Formats session duration and phase classification with duration-based color coding*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [🔄 Returns](#-returns) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetSessionPhaseDisplay Function
**Purpose:** Session phase text formatting for statusline display
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
  - [tempCtx (*temporal.TemporalContext)](#tempctx-temporaltemporalcontext)
- [🔄 Returns](#-returns)
  - [SessionPhaseDisplay](#sessionphasedisplay)
- [🎯 Behavior](#-behavior)
  - [Phase Color Assignment](#phase-color-assignment)
  - [Display Examples](#display-examples)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [Duration Only Display](#duration-only-display)
  - [With Color Formatting](#with-color-formatting)
  - [Conditional Phase Display](#conditional-phase-display)
  - [Full Statusline Integration](#full-statusline-integration)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [⚙ Configuration](#-configuration)
  - [Phase Color Constants](#phase-color-constants)
- [🔧 Edge Cases](#-edge-cases)
  - [New Session (No Duration)](#new-session-no-duration)
  - [Very Short Session](#very-short-session)
  - [Very Long Session](#very-long-session)
  - [Nil Context](#nil-context)
- [📊 Understanding Session Phase](#-understanding-session-phase)
  - [What Session Phase Represents](#what-session-phase-represents)
  - [Why Phase Classification](#why-phase-classification)
- [🔧 Extension Points](#-extension-points)
  - [Adding Phase-Specific Icons](#adding-phase-specific-icons)
  - [Adding Break Suggestions](#adding-break-suggestions)
  - [More Granular Phase Categories](#more-granular-phase-categories)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Duration never displays](#duration-never-displays)
  - [Phase seems wrong](#phase-seems-wrong)
  - [Want different phase thresholds](#want-different-phase-thresholds)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetSessionPhaseDisplay** formats session duration and phase classification suitable for statusline display.

> [!NOTE]
> **Pure presentation function.** Receives session duration data, applies phase-based color, outputs formatted display. No I/O, no side effects.

**Design:** Pure presentation layer - receives session duration from hooks/lib/temporal TemporalContext, applies phase color based on duration classification, outputs formatted display with visual indicators.

**Library:** [`statusline/lib/temporal`](../../../lib/temporal/)

**Source Code:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go)

**Library README:** [Temporal Library Documentation](README.md)

---

## 📝 Signature

```go
func GetSessionPhaseDisplay(tempCtx *temporal.TemporalContext) SessionPhaseDisplay
```

---

## ⚙ Parameters

### tempCtx (*temporal.TemporalContext)

Pointer to temporal context containing session duration information.

**Source:** `hooks/lib/temporal.GetTemporalContext()`

**Required Fields:**

| Field | Type | Purpose |
|-------|------|---------|
| `InternalTime.ElapsedFormatted` | string | Human-readable session duration |
| `InternalTime.SessionPhase` | string | Phase classification (fresh/working/long) |

**Example:**

```go
tempCtx, _ := temporal.GetTemporalContext()
// tempCtx.InternalTime.ElapsedFormatted = "45m"
// tempCtx.InternalTime.SessionPhase = "fresh"

phaseDisplay := GetSessionPhaseDisplay(tempCtx)
```

**Validation:** Accepts any TemporalContext (never fails)

---

## 🔄 Returns

### SessionPhaseDisplay

Formatted display structure with session phase presentation elements.

```go
type SessionPhaseDisplay struct {
    Duration string // Human-readable duration (e.g., "45m", "1h 23m")
    Phase    string // Phase classification (fresh/working/long)
    Color    string // Terminal color code based on phase
}
```

**Field Details:**

| Field | Description | Example Values |
|-------|-------------|----------------|
| **Duration** | Direct copy from `InternalTime.ElapsedFormatted` | `"45s"`, `"1m 23s"`, `"45m"`, `"1h 23m"` |
| **Phase** | Direct copy from `InternalTime.SessionPhase` | `"fresh"`, `"working"`, `"long"` |
| **Color** | Phase-based color code | Green (fresh), Yellow (long) |

**Color Mapping:**

| Phase | Color | Meaning |
|-------|:-----:|---------|
| `"fresh"` | 🟢 Green | Recently started session (< 1 hour typically) |
| `"working"` | 🟢 Green | Active session in progress (1-3 hours typically) |
| `"long"` | 🟡 Yellow | Extended session (> 3 hours typically) |
| Empty | None | Zero-state (no session data) |

**Guarantee:** Always returns valid SessionPhaseDisplay (never nil, never errors)

---

## 🎯 Behavior

### Phase Color Assignment

Duration-based color selection:

```go
// Check if duration available
if tempCtx.InternalTime.ElapsedFormatted == "" {
    return SessionPhaseDisplay{}  // Zero-state
}

// Determine color based on phase
phaseColor := display.Green  // Default (fresh/working)

if tempCtx.InternalTime.SessionPhase == "long" {
    phaseColor = display.Yellow  // Caution - consider break
}

return SessionPhaseDisplay{
    Duration: tempCtx.InternalTime.ElapsedFormatted,
    Phase:    tempCtx.InternalTime.SessionPhase,
    Color:    phaseColor,
}
```

### Display Examples

**Fresh Session (< 1 hour):**

```go
// TemporalContext: ElapsedFormatted = "45m", SessionPhase = "fresh"
display := GetSessionPhaseDisplay(tempCtx)
// → SessionPhaseDisplay{Duration: "45m", Phase: "fresh", Color: "green"}
```

**Working Session (1-3 hours):**

```go
// TemporalContext: ElapsedFormatted = "1h 23m", SessionPhase = "working"
display := GetSessionPhaseDisplay(tempCtx)
// → SessionPhaseDisplay{Duration: "1h 23m", Phase: "working", Color: "green"}
```

**Long Session (> 3 hours):**

```go
// TemporalContext: ElapsedFormatted = "3h 45m", SessionPhase = "long"
display := GetSessionPhaseDisplay(tempCtx)
// → SessionPhaseDisplay{Duration: "3h 45m", Phase: "long", Color: "yellow"}
```

**No Session Data:**

```go
// TemporalContext: ElapsedFormatted = "", SessionPhase = ""
display := GetSessionPhaseDisplay(tempCtx)
// → SessionPhaseDisplay{Duration: "", Phase: "", Color: ""}
// Zero-state: no data to display
```

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/temporal"
import "hooks/lib/temporal"

func buildStatusline() string {
    // Get temporal context
    tempCtx, _ := temporal.GetTemporalContext()

    // Get session phase display
    phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
    // → Duration: "45m", Phase: "fresh"

    var parts []string

    // Add to statusline if available
    if phaseDisplay.Duration != "" {
        parts = append(parts, fmt.Sprintf("%s (%s)", phaseDisplay.Duration, phaseDisplay.Phase))
        // → "45m (fresh)"
    }

    return strings.Join(parts, " | ")
}
```

### Duration Only Display

```go
phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)

// Display just the duration
if phaseDisplay.Duration != "" {
    parts = append(parts, phaseDisplay.Duration)
    // → "45m"
}
```

### With Color Formatting

```go
phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)

// Apply color for terminal output
if phaseDisplay.Duration != "" {
    colored := fmt.Sprintf("%s%s (%s)%s",
        phaseDisplay.Color,       // Start color
        phaseDisplay.Duration,    // Duration
        phaseDisplay.Phase,       // Phase
        display.Reset)            // Reset color

    fmt.Println(colored)
    // → Displays in green/yellow: 45m (fresh)
}
```

### Conditional Phase Display

```go
phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)

if phaseDisplay.Duration != "" {
    // Change display based on phase
    if phaseDisplay.Phase == "long" {
        parts = append(parts, fmt.Sprintf("⚠️ %s (consider break)", phaseDisplay.Duration))
    } else {
        parts = append(parts, phaseDisplay.Duration)
    }
}
```

### Full Statusline Integration

```go
func buildTemporalSection(tempCtx *temporal.TemporalContext) string {
    var parts []string

    // Time of day
    todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
    if todDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%s %s", todDisplay.Icon, todDisplay.Label))
    }

    // Session phase
    phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
    if phaseDisplay.Duration != "" {
        parts = append(parts, fmt.Sprintf("%s (%s)", phaseDisplay.Duration, phaseDisplay.Phase))
    }

    // Schedule
    schedDisplay := temporal.GetScheduleDisplay(tempCtx)
    if schedDisplay.HasInfo {
        parts = append(parts, schedDisplay.Label)
    }

    return strings.Join(parts, " | ")
    // → "🌅 morning | 45m (fresh) | Work schedule"
}
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Empty duration check | +25 | Detect zero-state |
| Phase color assignment | +25 | Apply correct color for phase |
| Struct population | +25 | All fields correctly set |
| Return valid struct | +25 | Always returns valid structure |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** Zero-state handling and phase-based color selection guarantee valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(1) - conditional check and assignment

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **String comparison** | O(1) | Empty string check |
| **Struct allocation** | ~50 bytes | SessionPhaseDisplay |
| **Execution time** | <1 μs | Microseconds |

**Memory:**

- Single SessionPhaseDisplay struct allocation (~50 bytes)
- No heap allocations beyond result struct

**Optimization:** Not needed - already optimal for simple field extraction

---

## ⚙ Configuration

### Phase Color Constants

```go
const (
    PhaseColorFresh   = display.Green   // < 1 hour typically
    PhaseColorWorking = display.Green   // 1-3 hours typically
    PhaseColorLong    = display.Yellow  // > 3 hours typically
)
```

**Location:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) (BODY section)

**Why constants:**

- Single source of truth for phase colors
- Easy to update if color scheme changes
- Consistent across all phase display code

---

## 🔧 Edge Cases

### New Session (No Duration)

```go
// TemporalContext with empty duration
emptyCtx := &temporal.TemporalContext{
    InternalTime: temporal.InternalTime{
        ElapsedFormatted: "",
        SessionPhase:     "",
    },
}

display := GetSessionPhaseDisplay(emptyCtx)
// Returns: SessionPhaseDisplay{Duration: "", Phase: "", Color: ""}
// Behavior: Zero-state, no display data
```

### Very Short Session

```go
// TemporalContext: ElapsedFormatted = "15s", SessionPhase = "fresh"
display := GetSessionPhaseDisplay(tempCtx)
// Returns: SessionPhaseDisplay{Duration: "15s", Phase: "fresh", Color: "green"}
// Behavior: Even very short durations display correctly
```

### Very Long Session

```go
// TemporalContext: ElapsedFormatted = "8h 45m", SessionPhase = "long"
display := GetSessionPhaseDisplay(tempCtx)
// Returns: SessionPhaseDisplay{Duration: "8h 45m", Phase: "long", Color: "yellow"}
// Behavior: Extremely long sessions still format correctly
```

### Nil Context

```go
// Nil pointer (edge case - should not happen in practice)
display := GetSessionPhaseDisplay(nil)
// Returns: SessionPhaseDisplay{Duration: "", Phase: "", Color: ""}
// Behavior: Graceful degradation, zero-state
```

---

## 📊 Understanding Session Phase

### What Session Phase Represents

**Session phase** = Classification of session duration relative to typical work periods

| Phase | Duration Range | Characteristics |
|-------|:--------------:|-----------------|
| **Fresh** | < 1 hour typically | High focus potential · Good for starting new work · Energy likely still high |
| **Working** | 1-3 hours typically | Deep work state · Momentum established · Productive flow |
| **Long** | > 3 hours typically | May need break soon · Fatigue possible · Consider pausing/reflecting |

*Exact thresholds determined by hooks/lib/temporal implementation*

### Why Phase Classification

**Awareness benefits:**

| Benefit | How It Helps |
|---------|--------------|
| **Break recognition** | Recognize when taking break might be wise |
| **Rhythm understanding** | Understand current work rhythm |
| **Pattern tracking** | Track productivity patterns over time |
| **Burnout prevention** | Avoid burnout through time awareness |

**Color coding purpose:**

- 🟢 **Green (fresh/working)** - Positive, productive state
- 🟡 **Yellow (long)** - Caution, consider break

---

## 🔧 Extension Points

### Adding Phase-Specific Icons

**Visual indicators for each phase:**

```go
// In temporal.go:
type SessionPhaseDisplay struct {
    Duration string
    Phase    string
    Icon     string  // Add this field
    Color    string
}

func GetSessionPhaseDisplay(tempCtx *temporal.TemporalContext) SessionPhaseDisplay {
    // ... existing logic ...

    phaseIcon := "⏱"  // Default timer
    switch tempCtx.InternalTime.SessionPhase {
    case "fresh":
        phaseIcon = "🌱"  // Growing/new
    case "working":
        phaseIcon = "⚙️"  // Active work
    case "long":
        phaseIcon = "⏰"  // Time warning
    }

    return SessionPhaseDisplay{
        Duration: tempCtx.InternalTime.ElapsedFormatted,
        Phase:    tempCtx.InternalTime.SessionPhase,
        Icon:     phaseIcon,
        Color:    phaseColor,
    }
}
```

### Adding Break Suggestions

**Display break recommendations:**

```go
type SessionPhaseDisplay struct {
    Duration       string
    Phase          string
    Color          string
    BreakSuggested bool  // New field
}

func GetSessionPhaseDisplay(tempCtx *temporal.TemporalContext) SessionPhaseDisplay {
    breakSuggested := false
    if tempCtx.InternalTime.SessionPhase == "long" {
        breakSuggested = true
    }

    return SessionPhaseDisplay{
        Duration:       tempCtx.InternalTime.ElapsedFormatted,
        Phase:          tempCtx.InternalTime.SessionPhase,
        Color:          phaseColor,
        BreakSuggested: breakSuggested,
    }
}
```

### More Granular Phase Categories

**Add intermediate phases:**

```go
// Add phases between existing ones
switch tempCtx.InternalTime.SessionPhase {
case "fresh":
    phaseColor = display.Green
case "warming":  // New phase (30m-1h)
    phaseColor = display.Cyan
case "working":
    phaseColor = display.Green
case "extended":  // New phase (2h-3h)
    phaseColor = display.Yellow
case "long":
    phaseColor = display.Red
}
```

---

## 🔍 Troubleshooting

### Duration never displays

**Problem:** phaseDisplay.Duration always empty

**Check:**

1. Verify TemporalContext.InternalTime.ElapsedFormatted populated
2. Check hooks/lib/temporal session tracking working
3. Debug: Print tempCtx.InternalTime before calling function

**Expected:** hooks/lib/temporal must track session duration for this to work

### Phase seems wrong

**Problem:** Shows "long" but session just started

**Cause:** InternalTime.SessionPhase value incorrect in TemporalContext

**Solution:**

- Verify hooks/lib/temporal phase classification logic
- Check session start time tracking
- Call `temporal.GetTemporalContext()` to verify phase value

**Expected:** Phase matches actual session duration

### Want different phase thresholds

**Problem:** 3 hours threshold for "long" doesn't match your work style

**Solution:** Adjust thresholds in hooks/lib/temporal (data layer, not this lib):

```go
// In hooks/lib/temporal (not statusline lib):
if elapsed < 1*time.Hour {
    phase = "fresh"
} else if elapsed < 2*time.Hour {  // Custom: 2 hours instead of 3
    phase = "working"
} else {
    phase = "long"
}
```

**Note:** This is data collection logic, not presentation - belongs in hooks/lib/temporal

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetTimeOfDayDisplay** | Time of day with icon | [timeofday.md](timeofday.md) |
| **GetScheduleDisplay** | Current schedule state | [schedule.md](schedule.md) |
| **GetCalendarDisplay** | Date and week information | [calendar.md](calendar.md) |

**Dependencies:**

| Component | Purpose | Link |
|-----------|---------|------|
| **hooks/lib/temporal** | Temporal data source | [hooks/lib/temporal/](~/.claude/hooks/lib/temporal/) |

**Future Functions:**

- `GetProductivityPhaseDisplay()` - Productivity pattern indicator
- `GetFocusDisplay()` - Focus/flow state indicator

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Temporal Library README** | Library overview and integration guide | [README.md](README.md) |
| **Time of Day Function** | GetTimeOfDayDisplay reference | [timeofday.md](timeofday.md) |
| **Schedule Function** | GetScheduleDisplay reference | [schedule.md](schedule.md) |
| **Calendar Function** | GetCalendarDisplay reference | [calendar.md](calendar.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) | GetSessionPhaseDisplay with comprehensive inline documentation |
| **Data Source** | [`hooks/lib/temporal/`](~/.claude/hooks/lib/temporal/) | Session duration tracking and phase classification |
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
  - Awareness of time passing - not to create anxiety, but to foster wisdom in how we spend our hours

**Supporting Verse:**

- **"Redeeming the time, because the days are evil"** - Ephesians 5:16
  - Wise stewardship of time - long sessions may need breaks, fresh sessions are opportunities

**Application:** Session phase display reminds us that time is passing. Not to create pressure or guilt, but to encourage wise stewardship. Awareness leads to wisdom - recognizing when to rest, when to continue, how to work faithfully within our finite capacity.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"So teach us to number our days, that we may apply our hearts unto wisdom." - Psalm 90:12*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
