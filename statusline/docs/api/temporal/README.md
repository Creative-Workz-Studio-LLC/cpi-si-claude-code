<div align="center">

# Temporal Awareness Display Library

**Pure Presentation Layer for Temporal Awareness**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Library](https://img.shields.io/badge/Type-Library-blue?style=flat)
![Functions](https://img.shields.io/badge/Functions-4-green?style=flat)
![Awareness](https://img.shields.io/badge/Awareness-Time%20%7C%20Session%20%7C%20Schedule%20%7C%20Calendar-orange?style=flat)

*Transforms temporal awareness data into formatted display structures with visual indicators*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [🎯 What It Provides](#-what-it-provides) • [🔑 Key Features](#-key-features) • [💻 Demo](#-demo) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - Temporal Awareness Display Library
**Purpose:** Transform temporal awareness data (time, session, schedule, calendar) into formatted displays
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Temporal Awareness Display Library](#temporal-awareness-display-library)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
  - [🎯 What It Provides](#-what-it-provides)
    - [Temporal Awareness Display Functions](#temporal-awareness-display-functions)
    - [Display Examples](#display-examples)
  - [🔑 Key Features](#-key-features)
    - [Time of Day Awareness](#time-of-day-awareness)
    - [Session Phase Classification](#session-phase-classification)
    - [Schedule Display](#schedule-display)
    - [Calendar Display](#calendar-display)
    - [Graceful Degradation](#graceful-degradation)
  - [💻 Demo](#-demo)
    - [Basic Usage](#basic-usage)
    - [Statusline Integration](#statusline-integration)
    - [Handling Zero-State](#handling-zero-state)
  - [🏗️ Architecture](#️-architecture)
    - [Component Type](#component-type)
    - [Dependencies](#dependencies)
    - [Design Principle](#design-principle)
  - [🏥 Health Scoring](#-health-scoring)
  - [⚡ Performance](#-performance)
    - [Function Performance](#function-performance)
  - [🔧 Extension Points](#-extension-points)
    - [Adding New Temporal Displays](#adding-new-temporal-displays)
    - [Custom Formatting Options](#custom-formatting-options)
  - [⚙️ Configuration](#️-configuration)
    - [Icon Choices](#icon-choices)
  - [🔍 Troubleshooting](#-troubleshooting)
    - [Expected Behaviors](#expected-behaviors)
    - [If Unexpected Results Occur](#if-unexpected-results-occur)
  - [🛠️ Modification Policy](#️-modification-policy)
    - [Safe to Modify](#safe-to-modify)
    - [Modify with Care](#modify-with-care)
    - [Never Modify](#never-modify)
  - [🚀 Future Expansions](#-future-expansions)
    - [Planned Features](#planned-features)
    - [Research Areas](#research-areas)
    - [Known Limitations](#known-limitations)
  - [📚 References](#-references)
    - [API Documentation](#api-documentation)
    - [Source Code](#source-code)
    - [Related Documentation](#related-documentation)
  - [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**Temporal Awareness Display Library** transforms temporal awareness data (time of day, session phase, schedule, calendar) into formatted display strings suitable for statusline space constraints.

> [!NOTE]
> **Pure presentation layer.** Receives temporal data from hooks/lib/temporal TemporalContext, outputs formatted displays (presentation). Never duplicates temporal tracking logic.

**Design:** Assessment vs Presentation separation - hooks/lib/temporal provides data, statusline/lib/temporal formats for display.

**Library:** [`statusline/lib/temporal`](../../../lib/temporal/)

**Source Code:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go)

**Version:** 1.0.0

---

## 🎯 What It Provides

### Temporal Awareness Display Functions

Four functions for comprehensive temporal awareness presentation:

| Function | Purpose | Display Example |
|----------|---------|-----------------|
| **GetTimeOfDayDisplay** | Time of day with appropriate icon | `🌅 morning` |
| **GetSessionPhaseDisplay** | Session duration and phase classification | `45m (fresh)` |
| **GetScheduleDisplay** | Current schedule state | `Work schedule (work)` |
| **GetCalendarDisplay** | Date and week information | `Tuesday, Nov 4, 2025 - Week 45` |

### Display Examples

**Formatted Output:**

```bash
🌅 morning          # Time of day with icon
45m (fresh)         # Session duration + phase
Work schedule       # Current schedule
Tue, Nov 4, W45     # Calendar with week
```

**Display Structures:**

```go
TimeOfDayDisplay{
    Icon:  "🌅",
    Label: "morning",
    Color: "yellow"
}

SessionPhaseDisplay{
    Duration: "45m",
    Phase:    "fresh",
    Color:    "green"
}

ScheduleDisplay{
    Label: "Work schedule",
    Type:  "work"
}

CalendarDisplay{
    DisplayString: "Tuesday, November 4, 2025 - Week 45"
}
```

---

## 🔑 Key Features

### Time of Day Awareness

Icon selection based on time:

| Time | Icon | Label | Color |
|:----:|:----:|-------|:-----:|
| Morning (6-12) | 🌅 | morning | yellow |
| Afternoon (12-18) | ☀️ | afternoon | orange |
| Evening (18-22) | 🌆 | evening | purple |
| Night (22-6) | 🌙 | night | blue |

### Session Phase Classification

Duration display with phase classification:

- **Fresh (0-30min)**: Green - optimal focus time
- **Working (30min-2h)**: Yellow - active session
- **Long (2h+)**: Orange - extended session

### Schedule Display

Current activity with type classification:

- Work schedule (work)
- Sleep schedule (sleep)
- Custom activities (custom)

### Calendar Display

Full date formatting with week number:

- Full weekday name
- Month, day, year
- ISO week number

### Graceful Degradation

All functions return valid structures even when data unavailable:

```go
// No session data
phaseDisplay := GetSessionPhaseDisplay(nil)
// → SessionPhaseDisplay{Duration: "", Phase: ""}

// Calling code checks empty fields
if phaseDisplay.Duration != "" {
    // Display session info
}
```

---

## 💻 Demo

### Basic Usage

```go
import "statusline/lib/temporal"
import "hooks/lib/temporal"

func main() {
    tempCtx, _ := temporal.GetTemporalContext()

    // Get temporal displays
    todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
    phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
    schedDisplay := temporal.GetScheduleDisplay(tempCtx)
    calDisplay := temporal.GetCalendarDisplay(tempCtx)

    // Display time of day
    fmt.Printf("%s %s\n", todDisplay.Icon, todDisplay.Label)
    // Output: 🌅 morning

    // Display session phase
    if phaseDisplay.Duration != "" {
        fmt.Printf("%s (%s)\n", phaseDisplay.Duration, phaseDisplay.Phase)
        // Output: 45m (fresh)
    }

    // Display schedule
    if schedDisplay.Label != "" {
        fmt.Printf("%s (%s)\n", schedDisplay.Label, schedDisplay.Type)
        // Output: Work schedule (work)
    }

    // Display calendar
    if calDisplay.DisplayString != "" {
        fmt.Println(calDisplay.DisplayString)
        // Output: Tuesday, November 4, 2025 - Week 45
    }
}
```

### Statusline Integration

```go
func buildStatusline(tempCtx *temporal.TemporalContext) string {
    var parts []string

    // Time of day
    todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
    parts = append(parts, fmt.Sprintf("%s %s", todDisplay.Icon, todDisplay.Label))

    // Session phase
    phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
    if phaseDisplay.Duration != "" {
        parts = append(parts, fmt.Sprintf("%s (%s)", phaseDisplay.Duration, phaseDisplay.Phase))
    }

    // Schedule
    schedDisplay := temporal.GetScheduleDisplay(tempCtx)
    if schedDisplay.Label != "" {
        parts = append(parts, schedDisplay.Label)
    }

    // Calendar
    calDisplay := temporal.GetCalendarDisplay(tempCtx)
    if calDisplay.DisplayString != "" {
        parts = append(parts, calDisplay.DisplayString)
    }

    return strings.Join(parts, " | ")
    // → "🌅 morning | 45m (fresh) | Work schedule | Tuesday, Nov 4, 2025 - Week 45"
}
```

### Handling Zero-State

```go
// Empty temporal context
emptyCtx := &temporal.TemporalContext{}

phaseDisplay := temporal.GetSessionPhaseDisplay(emptyCtx)
// → Duration: "", Phase: "" (no data to display)

schedDisplay := temporal.GetScheduleDisplay(emptyCtx)
// → Label: "", Type: "" (no data to display)

calDisplay := temporal.GetCalendarDisplay(emptyCtx)
// → DisplayString: "" (no data to display)

// Only time of day always has data (uses string labels)
todDisplay := temporal.GetTimeOfDayDisplay(emptyCtx)
// → Icon: "☀️", Label: "" (default icon, empty label)
```

---

## 🏗️ Architecture

### Component Type

**Ladder:** Library (Middle Rung)
**Role:** Presentation layer for temporal awareness

### Dependencies

**Lower Rungs (Data Sources):**

- `hooks/lib/temporal` - Temporal data collection (TemporalContext)
- `system/lib/display` - Color constants (optional for extensions)

**Higher Rungs (Consumers):**

- `statusline` - Main orchestrator

### Design Principle

> **Assessment vs Presentation Separation**
>
> Data collection happens in hooks/lib/temporal (assessment - tracking time, session, schedule). Presentation happens here (formatting with icons and visual indicators). Never duplicate temporal tracking logic.

**Data → Presentation Flow:**

```bash
TemporalContext                  Display Functions
(Data Source)                    (Presentation)
     ↓                                 ↓
ExternalTime.TimeOfDay          GetTimeOfDayDisplay()
     ↓                                 ↓
                                TimeOfDayDisplay{
                                  Icon: "🌅",
                                  Label: "morning",
                                  Color: "yellow"
                                }

InternalTime.ElapsedFormatted   GetSessionPhaseDisplay()
     ↓                                 ↓
                                SessionPhaseDisplay{
                                  Duration: "45m",
                                  Phase: "fresh",
                                  Color: "green"
                                }
```

---

## 🏥 Health Scoring

Pure presentation library - no health scoring infrastructure needed.

All operations guaranteed to succeed through graceful degradation:

- No session duration → Returns empty Duration/Phase fields
- No schedule → Returns empty Label/Type fields
- No calendar → Returns empty DisplayString

> [!NOTE]
> **This library cannot fail.** Always returns valid display structures. Health tracking would measure "successfully did nothing" which provides no value.

---

## ⚡ Performance

### Function Performance

| Function | Operation | Time | Memory |
|----------|-----------|:----:|:------:|
| **GetTimeOfDayDisplay** | Switch statement | <1 μs | ~50 bytes |
| **GetSessionPhaseDisplay** | Conditional check | <1 μs | ~50 bytes |
| **GetScheduleDisplay** | Field extraction | <1 μs | ~50 bytes |
| **GetCalendarDisplay** | sprintf formatting | <5 μs | ~100 bytes |

**Total Memory:** Four struct allocations per statusline render (~200 bytes total)

**Optimization:** This library needs no optimization. Field extraction and string formatting are already optimal for this use case.

---

## 🔧 Extension Points

### Adding New Temporal Displays

**Pattern:**

1. Check if hooks/lib/temporal provides the data (if not, extend temporal lib first)
2. Define new display structure in SETUP Types section
3. Create Get[Aspect]Display() function in BODY
4. Extract data from TemporalContext, format, return display struct
5. Update API documentation with new display
6. Add tests for new display function

**Example - Adding Timezone Display:**

```go
// In SETUP Types:
type TimezoneDisplay struct {
    Name   string
    Offset string
    Color  string
}

// In BODY:
func GetTimezoneDisplay(tempCtx *temporal.TemporalContext) TimezoneDisplay {
    if tempCtx.ExternalTime.Timezone == "" {
        return TimezoneDisplay{}
    }

    return TimezoneDisplay{
        Name:   tempCtx.ExternalTime.Timezone,
        Offset: tempCtx.ExternalTime.TimezoneOffset,
        Color:  display.Gray,
    }
}
```

### Custom Formatting Options

**Future enhancement:** Configurable display formats (compact, verbose, detailed)

```go
// Potential future API
func GetCalendarDisplayCompact(tempCtx) CalendarDisplay   // "Nov 4, W45"
func GetCalendarDisplayVerbose(tempCtx) CalendarDisplay   // "Tuesday, November 4, 2025, Week 45 of 2025"
```

---

## ⚙️ Configuration

### Icon Choices

**Time of Day Icons (Hardcoded):**

```go
const (
    IconMorning   = "🌅"  // Sunrise
    IconAfternoon = "☀️"   // Sun
    IconEvening   = "🌆"  // Cityscape at dusk
    IconNight     = "🌙"  // Crescent moon
)
```

**Modification:** Edit icon constants in `lib/temporal/temporal.go` function implementations

**Future:** External configuration file for runtime icon customization

---

## 🔍 Troubleshooting

**This library has no common failure modes - all operations guaranteed to succeed.**

### Expected Behaviors

**Empty fields when no data**

- No session duration → Duration: "" (correct, not an error)
- No schedule → Label: "" (correct, not an error)
- No calendar → DisplayString: "" (correct, not an error)

**Empty display data is zero-state**

- This is correct behavior
- Calling code checks field emptiness to determine display

### If Unexpected Results Occur

**Problem:** Time of day icon seems wrong

- **Cause:** TimeOfDay field in TemporalContext incorrect
- **Solution:** Check hooks/lib/temporal implementation
- **Expected values:** "morning", "afternoon", "evening", "night"

**Problem:** Session phase always empty

- **Cause:** ElapsedFormatted field not populated in TemporalContext
- **Solution:** Verify hooks/lib/temporal tracking session duration
- **Check:** Call temporal.GetTemporalContext() directly to verify

**Problem:** Calendar format looks wrong

- **Cause:** Calendar fields in TemporalContext incomplete
- **Solution:** Check ExternalCalendar population
- **Expected format:** "Weekday, Month Day, Year - Week N"

---

## 🛠️ Modification Policy

### Safe to Modify

- ✅ Add new temporal displays (timezone, circadian phase, etc.)
- ✅ Adjust icon choices for different times of day
- ✅ Extend display structures with additional fields
- ✅ Add helper functions for complex formatting logic
- ✅ Create alternative formatting functions (compact, verbose, etc.)

### Modify with Care

- ⚠️ Display struct fields - breaks statusline orchestrator
- ⚠️ Function signatures - breaks all calling code
- ⚠️ Icon selection logic - changes visual appearance
- ⚠️ Color field types - breaks display rendering

### Never Modify

- ❌ Pure function guarantee (stateless, no side effects)
- ❌ Graceful degradation (always return valid display)
- ❌ Data vs Presentation separation (hooks/lib/temporal provides data)
- ❌ Non-blocking guarantee (no I/O operations)

---

## 🚀 Future Expansions

### Planned Features

- ✅ Time of day display - **COMPLETED**
- ✅ Session phase display - **COMPLETED**
- ✅ Schedule display - **COMPLETED**
- ✅ Calendar display - **COMPLETED**
- ⏳ Timezone display
- ⏳ Circadian phase display
- ⏳ Next activity countdown
- ⏳ Custom formatting options (compact, verbose, detailed)

### Research Areas

- Color coding based on time of day (different morning/evening colors)
- Schedule conflict indicators
- Holiday-specific formatting
- Productivity pattern recognition

### Known Limitations

- No timezone display (assumes local time)
- Schedule display limited to single line
- Calendar display without holiday context
- No color variation based on time of day
- Icon choices hardcoded (not configurable)

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **GetTimeOfDayDisplay** | Time of day formatting with icon selection | [timeofday.md](./timeofday.md) |
| **GetSessionPhaseDisplay** | Session duration and phase classification | [sessionphase.md](./sessionphase.md) |
| **GetScheduleDisplay** | Current schedule state formatting | [schedule.md](./schedule.md) |
| **GetCalendarDisplay** | Date and week information formatting | [calendar.md](./calendar.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementations** | [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) | All four display functions with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7-library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |
| **Temporal Lib** | Temporal data collection library | [hooks/lib/temporal](~/.claude/hooks/lib/temporal/) |

---

## 📖 Biblical Foundation

**Core Verse:**

- **"To every thing there is a season, and a time to every purpose under the heaven"** - Ecclesiastes 3:1
  - Awareness of time and seasons - recognizing the appropriate moment for work, rest, and reflection

**Supporting Verse:**

- **"So teach us to number our days, that we may apply our hearts unto wisdom"** - Psalm 90:12
  - Counting time not for control, but for wisdom in living according to God's seasons

**Application:** Show temporal awareness (time of day, session duration, schedule, calendar) to foster wisdom in timing and rhythm. Not tracking time for productivity optimization, but for wisdom in honoring God's seasons and purposes.

**Why This Matters:** Many time-tracking systems focus on extraction ("maximize productivity"). This library focuses on wisdom ("honor the season"). Morning vs evening, fresh vs long session, work vs rest - each has its appropriate purpose and rhythm under heaven.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
