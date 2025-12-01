<div align="center">

# GetScheduleDisplay Function Reference

**Schedule Information Display Formatting for Statusline**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Pure](https://img.shields.io/badge/Pure-Function-success?style=flat)
![HasInfo](https://img.shields.io/badge/Pattern-HasInfo-orange?style=flat)

*Formats current schedule information with HasInfo pattern for zero-state handling*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [🔄 Returns](#-returns) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetScheduleDisplay Function
**Purpose:** Schedule information text formatting for statusline display
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
  - [ScheduleDisplay](#scheduledisplay)
- [🎯 Behavior](#-behavior)
  - [Simple Field Extraction with HasInfo](#simple-field-extraction-with-hasinfo)
  - [Display Examples](#display-examples)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [With Type Indicator](#with-type-indicator)
  - [Type-Specific Icons](#type-specific-icons)
  - [Conditional Display by Type](#conditional-display-by-type)
  - [Full Statusline Integration](#full-statusline-integration)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [⚙ Configuration](#-configuration)
  - [Schedule Type Constants](#schedule-type-constants)
- [🔧 Edge Cases](#-edge-cases)
  - [No Schedule Active](#no-schedule-active)
  - [Long Schedule Description](#long-schedule-description)
  - [Unknown Schedule Type](#unknown-schedule-type)
  - [Nil Context](#nil-context)
- [📊 Understanding Schedule Display](#-understanding-schedule-display)
  - [What Schedule Represents](#what-schedule-represents)
  - [Why Schedule Display](#why-schedule-display)
- [🔧 Extension Points](#-extension-points)
  - [Adding Schedule Colors](#adding-schedule-colors)
  - [Adding Schedule Icons](#adding-schedule-icons)
  - [Adding Schedule Metadata](#adding-schedule-metadata)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Schedule never displays](#schedule-never-displays)
  - [Schedule seems wrong](#schedule-seems-wrong)
  - [Want custom schedule types](#want-custom-schedule-types)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetScheduleDisplay** formats current schedule information suitable for statusline display.

> [!NOTE]
> **Pure presentation function.** Receives schedule data, applies HasInfo pattern for zero-state handling, outputs formatted display. No I/O, no side effects.

**Design:** Pure presentation layer - receives schedule data from hooks/lib/temporal TemporalContext, extracts description and type, outputs formatted display structure with HasInfo flag.

**Library:** [`statusline/lib/temporal`](../../../lib/temporal/)

**Source Code:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go)

**Library README:** [Temporal Library Documentation](README.md)

---

## 📝 Signature

```go
func GetScheduleDisplay(tempCtx *temporal.TemporalContext) ScheduleDisplay
```

---

## ⚙ Parameters

### tempCtx (*temporal.TemporalContext)

Pointer to temporal context containing schedule information.

**Source:** `hooks/lib/temporal.GetTemporalContext()`

**Required Fields:**

| Field | Type | Purpose |
|-------|------|---------|
| `InternalSchedule.Description` | string | Human-readable schedule description |
| `InternalSchedule.Type` | string | Schedule type classification |

**Example:**

```go
tempCtx, _ := temporal.GetTemporalContext()
// tempCtx.InternalSchedule.Description = "Work schedule"
// tempCtx.InternalSchedule.Type = "work"

schedDisplay := GetScheduleDisplay(tempCtx)
```

**Validation:** Accepts any TemporalContext (never fails)

---

## 🔄 Returns

### ScheduleDisplay

Formatted display structure with schedule presentation elements.

```go
type ScheduleDisplay struct {
    Label   string // Schedule description
    Type    string // Schedule type (work/sleep/personal/etc.)
    HasInfo bool   // True if schedule data available
}
```

**Field Details:**

| Field | Description | Example Values |
|-------|-------------|----------------|
| **Label** | Direct copy from `InternalSchedule.Description` | `"Work schedule"`, `"Sleep schedule"`, `"Personal time"` |
| **Type** | Direct copy from `InternalSchedule.Type` | `"work"`, `"sleep"`, `"personal"`, `"break"` |
| **HasInfo** | True when Label not empty | `true` (has data), `false` (zero-state) |

**Common Schedule Types:**

| Type | Description | Typical Hours |
|------|-------------|:-------------:|
| `"work"` | Active working hours | 9am-5pm |
| `"sleep"` | Rest/sleep period | 10pm-6am |
| `"personal"` | Personal time, leisure | Variable |
| `"break"` | Scheduled break | 15-30 min |
| `"meeting"` | Specific meetings | Variable |
| `"focus"` | Deep work/focus time | Variable |

**Guarantee:** Always returns valid ScheduleDisplay (never nil, never errors)

---

## 🎯 Behavior

### Simple Field Extraction with HasInfo

Straightforward data extraction with zero-state detection:

```go
// Extract schedule data
schedLabel := tempCtx.InternalSchedule.Description
schedType := tempCtx.InternalSchedule.Type

// HasInfo = true when we have a description
hasInfo := schedLabel != ""

return ScheduleDisplay{
    Label:   schedLabel,
    Type:    schedType,
    HasInfo: hasInfo,
}
```

### Display Examples

**Work Schedule:**

```go
// TemporalContext: Description = "Work schedule", Type = "work"
display := GetScheduleDisplay(tempCtx)
// → ScheduleDisplay{Label: "Work schedule", Type: "work", HasInfo: true}
```

**Sleep Schedule:**

```go
// TemporalContext: Description = "Sleep schedule (crosses midnight)", Type = "sleep"
display := GetScheduleDisplay(tempCtx)
// → ScheduleDisplay{Label: "Sleep schedule (crosses midnight)", Type: "sleep", HasInfo: true}
```

**Personal Time:**

```go
// TemporalContext: Description = "Personal time", Type = "personal"
display := GetScheduleDisplay(tempCtx)
// → ScheduleDisplay{Label: "Personal time", Type: "personal", HasInfo: true}
```

**No Schedule Data:**

```go
// TemporalContext: Description = "", Type = ""
display := GetScheduleDisplay(tempCtx)
// → ScheduleDisplay{Label: "", Type: "", HasInfo: false}
// Zero-state: no schedule information
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

    // Get schedule display
    schedDisplay := temporal.GetScheduleDisplay(tempCtx)
    // → Label: "Work schedule", Type: "work", HasInfo: true

    var parts []string

    // Add to statusline if available
    if schedDisplay.HasInfo {
        parts = append(parts, schedDisplay.Label)
        // → "Work schedule"
    }

    return strings.Join(parts, " | ")
}
```

### With Type Indicator

```go
schedDisplay := temporal.GetScheduleDisplay(tempCtx)

// Display label with type in parentheses
if schedDisplay.HasInfo {
    parts = append(parts, fmt.Sprintf("%s (%s)", schedDisplay.Label, schedDisplay.Type))
    // → "Work schedule (work)"
}
```

### Type-Specific Icons

```go
schedDisplay := temporal.GetScheduleDisplay(tempCtx)

if schedDisplay.HasInfo {
    // Add icon based on schedule type
    var icon string
    switch schedDisplay.Type {
    case "work":
        icon = "💼"
    case "sleep":
        icon = "😴"
    case "personal":
        icon = "🏠"
    case "break":
        icon = "☕"
    case "meeting":
        icon = "📅"
    case "focus":
        icon = "🎯"
    default:
        icon = "📋"
    }

    parts = append(parts, fmt.Sprintf("%s %s", icon, schedDisplay.Label))
    // → "💼 Work schedule"
}
```

### Conditional Display by Type

```go
schedDisplay := temporal.GetScheduleDisplay(tempCtx)

// Only show certain schedule types
if schedDisplay.HasInfo {
    if schedDisplay.Type == "work" || schedDisplay.Type == "sleep" {
        parts = append(parts, schedDisplay.Label)
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
| Description extraction | +25 | Extract schedule label |
| Type extraction | +25 | Extract schedule type |
| HasInfo flag determination | +25 | Detect zero-state correctly |
| Struct population | +25 | All fields correctly set |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** Simple field extraction with HasInfo pattern guarantees valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(1) - direct field access

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Field extraction** | O(1) | Direct struct field access |
| **String comparison** | O(1) | Empty string check for HasInfo |
| **Struct allocation** | ~50 bytes | ScheduleDisplay |
| **Execution time** | <1 μs | Microseconds |

**Memory:**

- Single ScheduleDisplay struct allocation (~50 bytes)
- No heap allocations beyond result struct

**Optimization:** Not needed - already optimal for simple field extraction

---

## ⚙ Configuration

### Schedule Type Constants

```go
const (
    ScheduleTypeWork     = "work"     // Active working hours
    ScheduleTypeSleep    = "sleep"    // Rest/sleep period
    ScheduleTypePersonal = "personal" // Personal time, leisure
    ScheduleTypeBreak    = "break"    // Scheduled break
    ScheduleTypeMeeting  = "meeting"  // Specific meetings
    ScheduleTypeFocus    = "focus"    // Deep work/focus time
)
```

**Location:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) (SETUP section)

**Why constants:**

- Single source of truth for schedule types
- Easy to add new schedule types
- Consistent across all schedule-related code

---

## 🔧 Edge Cases

### No Schedule Active

```go
// TemporalContext with empty schedule
emptyCtx := &temporal.TemporalContext{
    InternalSchedule: temporal.InternalSchedule{
        Description: "",
        Type:        "",
    },
}

display := GetScheduleDisplay(emptyCtx)
// Returns: ScheduleDisplay{Label: "", Type: "", HasInfo: false}
// Behavior: Zero-state, no schedule data
```

### Long Schedule Description

```go
// TemporalContext: Description = "Work schedule (important meeting at 2pm)", Type = "work"
display := GetScheduleDisplay(tempCtx)
// Returns: ScheduleDisplay{Label: "Work schedule (important meeting at 2pm)", Type: "work", HasInfo: true}
// Behavior: Full description preserved (statusline orchestrator may truncate if needed)
```

### Unknown Schedule Type

```go
// TemporalContext: Description = "Custom activity", Type = "custom"
display := GetScheduleDisplay(tempCtx)
// Returns: ScheduleDisplay{Label: "Custom activity", Type: "custom", HasInfo: true}
// Behavior: Accepts any type value, no validation
```

### Nil Context

```go
// Nil pointer (edge case - should not happen in practice)
display := GetScheduleDisplay(nil)
// Returns: ScheduleDisplay{Label: "", Type: "", HasInfo: false}
// Behavior: Graceful degradation, zero-state
```

---

## 📊 Understanding Schedule Display

### What Schedule Represents

**Schedule** = Current planned activity or time block based on internal schedule tracking

**Purpose:**

| Benefit | Description |
|---------|-------------|
| **Time block awareness** | Know what time block you're in |
| **Activity alignment** | Maintain awareness of planned activities |
| **Reality check** | Recognize when actual work doesn't match schedule |
| **Rhythm tracking** | Track adherence to planned rhythm |

### Why Schedule Display

**Use cases:**

```go
// Alignment check
"Am I working during work hours?"

// Rest reminder
"Should I be sleeping?"

// Context awareness
"What's the current planned activity?"

// Pattern recognition
"Do I consistently work outside work schedule?"
```

**Benefits:**

- Recognize when to shift activities
- Maintain healthy work-life boundaries
- Notice patterns of schedule misalignment
- Honor seasons of work and rest

---

## 🔧 Extension Points

### Adding Schedule Colors

**Type-based color coding:**

```go
type ScheduleDisplay struct {
    Label   string
    Type    string
    Color   string // Add this field
    HasInfo bool
}

func GetScheduleDisplay(tempCtx *temporal.TemporalContext) ScheduleDisplay {
    schedLabel := tempCtx.InternalSchedule.Description
    schedType := tempCtx.InternalSchedule.Type
    hasInfo := schedLabel != ""

    // Color based on schedule type
    schedColor := display.Gray // Default
    if hasInfo {
        switch schedType {
        case ScheduleTypeWork:
            schedColor = display.Blue
        case ScheduleTypeSleep:
            schedColor = display.Cyan
        case ScheduleTypePersonal:
            schedColor = display.Green
        case ScheduleTypeBreak:
            schedColor = display.Yellow
        }
    }

    return ScheduleDisplay{
        Label:   schedLabel,
        Type:    schedType,
        Color:   schedColor,
        HasInfo: hasInfo,
    }
}
```

### Adding Schedule Icons

**Built-in icon selection:**

```go
type ScheduleDisplay struct {
    Label   string
    Type    string
    Icon    string // Add this field
    HasInfo bool
}

func GetScheduleDisplay(tempCtx *temporal.TemporalContext) ScheduleDisplay {
    schedLabel := tempCtx.InternalSchedule.Description
    schedType := tempCtx.InternalSchedule.Type
    hasInfo := schedLabel != ""

    // Icon based on schedule type
    schedIcon := "📋" // Default
    if hasInfo {
        switch schedType {
        case ScheduleTypeWork:
            schedIcon = "💼"
        case ScheduleTypeSleep:
            schedIcon = "😴"
        case ScheduleTypePersonal:
            schedIcon = "🏠"
        case ScheduleTypeBreak:
            schedIcon = "☕"
        case ScheduleTypeMeeting:
            schedIcon = "📅"
        case ScheduleTypeFocus:
            schedIcon = "🎯"
        }
    }

    return ScheduleDisplay{
        Label:   schedLabel,
        Type:    schedType,
        Icon:    schedIcon,
        HasInfo: hasInfo,
    }
}
```

### Adding Schedule Metadata

**Enhanced schedule information:**

```go
type ScheduleDisplay struct {
    Label     string
    Type      string
    StartTime string  // When schedule block started
    EndTime   string  // When it ends
    Progress  float64 // Percentage through schedule (0-100)
    HasInfo   bool
}

func GetScheduleDisplay(tempCtx *temporal.TemporalContext) ScheduleDisplay {
    // Extract schedule data
    schedLabel := tempCtx.InternalSchedule.Description
    schedType := tempCtx.InternalSchedule.Type
    hasInfo := schedLabel != ""

    // Calculate progress through schedule block if data available
    var progress float64
    if hasInfo && tempCtx.InternalSchedule.StartTime != "" {
        progress = calculateScheduleProgress(
            tempCtx.InternalSchedule.StartTime,
            tempCtx.InternalSchedule.EndTime,
            tempCtx.ExternalTime.Current,
        )
    }

    return ScheduleDisplay{
        Label:     schedLabel,
        Type:      schedType,
        StartTime: tempCtx.InternalSchedule.StartTime,
        EndTime:   tempCtx.InternalSchedule.EndTime,
        Progress:  progress,
        HasInfo:   hasInfo,
    }
}
```

---

## 🔍 Troubleshooting

### Schedule never displays

**Problem:** schedDisplay.HasInfo always false

**Check:**

1. Verify TemporalContext.InternalSchedule.Description populated
2. Check hooks/lib/temporal schedule tracking configured
3. Debug: Print tempCtx.InternalSchedule before calling function

**Expected:** hooks/lib/temporal must provide schedule data for this to work

### Schedule seems wrong

**Problem:** Shows "Work schedule" during sleep hours

**Cause:** InternalSchedule data incorrect in TemporalContext

**Solution:**

- Verify hooks/lib/temporal schedule logic
- Check schedule configuration/patterns
- Ensure time-based schedule transitions working

**Expected:** Schedule matches actual time of day and planned activities

### Want custom schedule types

**Problem:** Need additional schedule types beyond work/sleep/personal

**Solution:** Extend schedule types in hooks/lib/temporal (data layer):

```go
// In hooks/lib/temporal (not statusline lib):
type Schedule struct {
    Description string
    Type        string  // Add new types: "focus", "meeting", "break", etc.
    StartTime   string
    EndTime     string
}
```

**Note:** This is data collection logic, belongs in hooks/lib/temporal

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetTimeOfDayDisplay** | Time of day with icon | [timeofday.md](timeofday.md) |
| **GetSessionPhaseDisplay** | Session duration and phase | [sessionphase.md](sessionphase.md) |
| **GetCalendarDisplay** | Date and week information | [calendar.md](calendar.md) |

**Dependencies:**

| Component | Purpose | Link |
|-----------|---------|------|
| **hooks/lib/temporal** | Temporal data source | [hooks/lib/temporal/](~/.claude/hooks/lib/temporal/) |

**Future Functions:**

- `GetScheduleProgressDisplay()` - Progress through current schedule block
- `GetNextScheduleDisplay()` - Upcoming schedule information
- `GetScheduleConflictDisplay()` - Schedule vs reality mismatch indicator

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Temporal Library README** | Library overview and integration guide | [README.md](README.md) |
| **Time of Day Function** | GetTimeOfDayDisplay reference | [timeofday.md](timeofday.md) |
| **Session Phase Function** | GetSessionPhaseDisplay reference | [sessionphase.md](sessionphase.md) |
| **Calendar Function** | GetCalendarDisplay reference | [calendar.md](calendar.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) | GetScheduleDisplay with comprehensive inline documentation |
| **Data Source** | [`hooks/lib/temporal/`](~/.claude/hooks/lib/temporal/) | Schedule tracking and classification |
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

- **"To every thing there is a season, and a time to every purpose under the heaven"** - Ecclesiastes 3:1
  - Recognizing appointed times - there is a time for work, a time for rest, a time for each activity

**Supporting Verse:**

- **"Be very careful, then, how you live—not as unwise but as wise, making the most of every opportunity"** - Ephesians 5:15-16
  - Living intentionally within time - not rigidity, but rhythmic awareness

**Application:** Schedule display helps us live intentionally within seasons and times. Knowing "this is work time" or "this is rest time" helps us honor the season we're in. Not rigid adherence to schedules, but awareness that enables wise stewardship of our hours - working faithfully during work time, resting genuinely during rest time, honoring the rhythms God established.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"To every thing there is a season, and a time to every purpose under the heaven." - Ecclesiastes 3:1*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
