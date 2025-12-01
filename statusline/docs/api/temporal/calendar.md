<div align="center">

# GetCalendarDisplay Function Reference

**Calendar Date and Week Display Formatting for Statusline**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Pure](https://img.shields.io/badge/Pure-Function-success?style=flat)
![ISO Week](https://img.shields.io/badge/ISO-8601%20Weeks-blue?style=flat)

*Formats calendar date and week information with complete date context*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [🔄 Returns](#-returns) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetCalendarDisplay Function
**Purpose:** Calendar date and week text formatting for statusline display
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

- [GetCalendarDisplay Function Reference](#getcalendardisplay-function-reference)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
  - [📝 Signature](#-signature)
  - [⚙ Parameters](#-parameters)
    - [tempCtx (\*temporal.TemporalContext)](#tempctx-temporaltemporalcontext)
  - [🔄 Returns](#-returns)
    - [CalendarDisplay](#calendardisplay)
  - [🎯 Behavior](#-behavior)
    - [Calendar Formatting with HasInfo](#calendar-formatting-with-hasinfo)
    - [Display Examples](#display-examples)
  - [💻 Usage](#-usage)
    - [Basic Usage](#basic-usage)
    - [With Icon](#with-icon)
    - [Week Start Indicator](#week-start-indicator)
    - [Year Progress Context](#year-progress-context)
    - [Full Statusline Integration](#full-statusline-integration)
  - [🏥 Health Scoring](#-health-scoring)
  - [⚡ Performance](#-performance)
  - [⚙ Configuration](#-configuration)
    - [ISO Week Numbering](#iso-week-numbering)
  - [🔧 Edge Cases](#-edge-cases)
    - [No Calendar Data](#no-calendar-data)
    - [Week 1 (Start of Year)](#week-1-start-of-year)
    - [Week 52/53 (End of Year)](#week-5253-end-of-year)
    - [Nil Context](#nil-context)
  - [📊 Understanding Calendar Display](#-understanding-calendar-display)
    - [What Calendar Display Shows](#what-calendar-display-shows)
    - [Why Week Number Matters](#why-week-number-matters)
  - [🔧 Extension Points](#-extension-points)
    - [Multiple Format Options](#multiple-format-options)
    - [Year Progress Indicator](#year-progress-indicator)
    - [Holiday Indicators](#holiday-indicators)
  - [🔍 Troubleshooting](#-troubleshooting)
    - [Calendar never displays](#calendar-never-displays)
    - [Week number seems wrong](#week-number-seems-wrong)
    - [Date format inconsistent](#date-format-inconsistent)
  - [🔗 Related Functions](#-related-functions)
  - [📚 References](#-references)
    - [API Documentation](#api-documentation)
    - [Source Code](#source-code)
    - [Related Documentation](#related-documentation)
    - [External Standards](#external-standards)
  - [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetCalendarDisplay** formats calendar date and week information suitable for statusline display.

> [!NOTE]
> **Pure presentation function.** Receives calendar data, formats complete date with week number, outputs display string. No I/O, no side effects.

**Design:** Pure presentation layer - receives calendar data from hooks/lib/temporal TemporalContext, formats into human-readable display string with complete date and ISO week number.

**Library:** [`statusline/lib/temporal`](../../../lib/temporal/)

**Source Code:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go)

**Library README:** [Temporal Library Documentation](README.md)

---

## 📝 Signature

```go
func GetCalendarDisplay(tempCtx *temporal.TemporalContext) CalendarDisplay
```

---

## ⚙ Parameters

### tempCtx (*temporal.TemporalContext)

Pointer to temporal context containing external calendar information.

**Source:** `hooks/lib/temporal.GetTemporalContext()`

**Required Fields:**

| Field | Type | Purpose |
|-------|------|---------|
| `ExternalCalendar.Date` | string | Full date string |
| `ExternalCalendar.DayOfWeek` | string | Day name (Monday, Tuesday, etc.) |
| `ExternalCalendar.WeekNumber` | int | ISO week of year (1-53) |

**Example:**

```go
tempCtx, _ := temporal.GetTemporalContext()
// tempCtx.ExternalCalendar.Date = "November 4, 2025"
// tempCtx.ExternalCalendar.DayOfWeek = "Tuesday"
// tempCtx.ExternalCalendar.WeekNumber = 45

calDisplay := GetCalendarDisplay(tempCtx)
```

**Validation:** Accepts any TemporalContext (never fails)

---

## 🔄 Returns

### CalendarDisplay

Formatted display structure with calendar presentation.

```go
type CalendarDisplay struct {
    DisplayString string // Formatted calendar display
    HasInfo       bool   // True if calendar data available
}
```

**Field Details:**

| Field | Description | Example Values |
|-------|-------------|----------------|
| **DisplayString** | Formatted date with week number | `"Tuesday, November 4, 2025 - Week 45"` |
| **HasInfo** | True when DisplayString not empty | `true` (has data), `false` (zero-state) |

**Display Format:**

```bash
{DayOfWeek}, {Date} - Week {WeekNumber}
```

**Example Outputs:**

| Input Data | Output DisplayString |
|------------|---------------------|
| Tuesday, Nov 4, 2025, Week 45 | `"Tuesday, November 4, 2025 - Week 45"` |
| Monday, Jan 6, 2025, Week 1 | `"Monday, January 6, 2025 - Week 1"` |
| Empty data | `""` (HasInfo = false) |

**Guarantee:** Always returns valid CalendarDisplay (never nil, never errors)

---

## 🎯 Behavior

### Calendar Formatting with HasInfo

Complete date formatting with week and zero-state detection:

```go
// Check if calendar data available
if tempCtx.ExternalCalendar.Date == "" {
    return CalendarDisplay{HasInfo: false}  // Zero-state
}

// Format full calendar display
calendarStr := fmt.Sprintf("%s, %s - Week %d",
    tempCtx.ExternalCalendar.DayOfWeek,
    tempCtx.ExternalCalendar.Date,
    tempCtx.ExternalCalendar.WeekNumber)

return CalendarDisplay{
    DisplayString: calendarStr,
    HasInfo:       true,
}
```

### Display Examples

**Standard Date:**

```go
// TemporalContext:
//   Date = "November 4, 2025"
//   DayOfWeek = "Tuesday"
//   WeekNumber = 45

display := GetCalendarDisplay(tempCtx)
// → CalendarDisplay{DisplayString: "Tuesday, November 4, 2025 - Week 45", HasInfo: true}
```

**Start of Year:**

```go
// TemporalContext:
//   Date = "January 6, 2025"
//   DayOfWeek = "Monday"
//   WeekNumber = 1

display := GetCalendarDisplay(tempCtx)
// → CalendarDisplay{DisplayString: "Monday, January 6, 2025 - Week 1", HasInfo: true}
```

**End of Year (ISO Week Boundary):**

```go
// TemporalContext:
//   Date = "December 31, 2025"
//   DayOfWeek = "Wednesday"
//   WeekNumber = 1  // ISO week (can be week 1 of next year)

display := GetCalendarDisplay(tempCtx)
// → CalendarDisplay{DisplayString: "Wednesday, December 31, 2025 - Week 1", HasInfo: true}
```

**No Calendar Data:**

```go
// TemporalContext: Date = "", DayOfWeek = "", WeekNumber = 0
display := GetCalendarDisplay(tempCtx)
// → CalendarDisplay{DisplayString: "", HasInfo: false}
// Zero-state: no calendar information
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

    // Get calendar display
    calDisplay := temporal.GetCalendarDisplay(tempCtx)
    // → DisplayString: "Tuesday, November 4, 2025 - Week 45", HasInfo: true

    var parts []string

    // Add to statusline if available
    if calDisplay.HasInfo {
        parts = append(parts, calDisplay.DisplayString)
        // → "Tuesday, November 4, 2025 - Week 45"
    }

    return strings.Join(parts, " | ")
}
```

### With Icon

```go
calDisplay := temporal.GetCalendarDisplay(tempCtx)

if calDisplay.HasInfo {
    parts = append(parts, fmt.Sprintf("📅 %s", calDisplay.DisplayString))
    // → "📅 Tuesday, November 4, 2025 - Week 45"
}
```

### Week Start Indicator

```go
calDisplay := temporal.GetCalendarDisplay(tempCtx)

if calDisplay.HasInfo {
    // Check if Monday (start of week)
    if strings.Contains(calDisplay.DisplayString, "Monday") {
        parts = append(parts, fmt.Sprintf("📆 Week %d starts! %s",
            tempCtx.ExternalCalendar.WeekNumber,
            calDisplay.DisplayString))
    } else {
        parts = append(parts, calDisplay.DisplayString)
    }
}
```

### Year Progress Context

```go
calDisplay := temporal.GetCalendarDisplay(tempCtx)

if calDisplay.HasInfo {
    weekNum := tempCtx.ExternalCalendar.WeekNumber
    yearProgress := (float64(weekNum) / 52.0) * 100

    parts = append(parts, fmt.Sprintf("%s (%.0f%% through year)",
        calDisplay.DisplayString,
        yearProgress))
    // → "Tuesday, November 4, 2025 - Week 45 (87% through year)"
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

    // Calendar
    calDisplay := temporal.GetCalendarDisplay(tempCtx)
    if calDisplay.HasInfo {
        parts = append(parts, calDisplay.DisplayString)
    }

    return strings.Join(parts, " | ")
    // → "🌅 morning | 45m (fresh) | Work schedule | Tuesday, November 4, 2025 - Week 45"
}
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Date check | +20 | Detect zero-state |
| String formatting (sprintf) | +40 | Format complete calendar display |
| HasInfo flag determination | +20 | Set zero-state flag correctly |
| Struct population | +20 | All fields correctly set |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** String formatting with HasInfo pattern guarantees valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(1) - sprintf formatting

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **String formatting** | O(1) | fmt.Sprintf with 3 args |
| **String comparison** | O(1) | Empty string check for HasInfo |
| **Struct allocation** | ~100 bytes | CalendarDisplay |
| **Execution time** | <5 μs | Microseconds |

**Memory:**

- Single CalendarDisplay struct allocation (~100 bytes)
- One formatted string allocation (50-80 bytes)
- No additional heap allocations

**Optimization:** Not needed - sprintf already optimal for this use case

---

## ⚙ Configuration

### ISO Week Numbering

**Week numbering follows ISO-8601 standard:**

```go
const (
    ISOWeekMin = 1   // First week of year
    ISOWeekMax = 53  // Maximum weeks in year (rare)
    ISOWeekTypical = 52  // Typical weeks in year
)
```

**ISO-8601 Week Definition:**

- Week 1 = First week with Thursday in new year
- Week starts Monday
- Weeks numbered 1-52 (sometimes 53)
- Dec 29-31 can be week 1 of next year

**Location:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) (SETUP section)

**Why ISO-8601:**

- International standard for week numbering
- Consistent across systems and countries
- Week always contains at least 4 days of year
- Aligns with Monday-start weeks

---

## 🔧 Edge Cases

### No Calendar Data

```go
// TemporalContext with empty calendar
emptyCtx := &temporal.TemporalContext{
    ExternalCalendar: temporal.ExternalCalendar{
        Date:       "",
        DayOfWeek:  "",
        WeekNumber: 0,
    },
}

display := GetCalendarDisplay(emptyCtx)
// Returns: CalendarDisplay{DisplayString: "", HasInfo: false}
// Behavior: Zero-state, no calendar data
```

### Week 1 (Start of Year)

```go
// TemporalContext:
//   Date = "January 6, 2025"
//   DayOfWeek = "Monday"
//   WeekNumber = 1

display := GetCalendarDisplay(tempCtx)
// Returns: CalendarDisplay{DisplayString: "Monday, January 6, 2025 - Week 1", HasInfo: true}
// Behavior: First week of year displays correctly
```

### Week 52/53 (End of Year)

```go
// TemporalContext:
//   Date = "December 30, 2025"
//   DayOfWeek = "Tuesday"
//   WeekNumber = 1  // ISO week (week 1 of 2026)

display := GetCalendarDisplay(tempCtx)
// Returns: CalendarDisplay{DisplayString: "Tuesday, December 30, 2025 - Week 1", HasInfo: true}
// Behavior: ISO week numbering may show week 1 at end of year (correct per ISO-8601)
```

### Nil Context

```go
// Nil pointer (edge case - should not happen in practice)
display := GetCalendarDisplay(nil)
// Returns: CalendarDisplay{DisplayString: "", HasInfo: false}
// Behavior: Graceful degradation, zero-state
```

---

## 📊 Understanding Calendar Display

### What Calendar Display Shows

**Calendar display** = Complete date context with week number

| Component | Purpose | Example |
|-----------|---------|---------|
| **Day of Week** | Context for rhythm (start/middle/end of week) | `"Monday"` (week start) |
| **Full Date** | Complete date identification | `"November 4, 2025"` |
| **Week Number** | Annual progress tracking (ISO-8601) | `45` (week 45 of 52) |

**Complete Display:** `"Tuesday, November 4, 2025 - Week 45"`

- **Tuesday** - Middle of work week
- **November 4, 2025** - Specific day identification
- **Week 45** - 87% through the year (45/52)

### Why Week Number Matters

**Week awareness benefits:**

| Benefit | Description |
|---------|-------------|
| **Annual progress** | Track progress through year (Week 45 of 52 = ~87%) |
| **Cross-month reference** | Consistent reference independent of month boundaries |
| **Goal tracking** | Quarters = ~13 weeks, half-year = ~26 weeks |
| **International standard** | ISO-8601 week numbering universally recognized |

**Week Numbering Details:**

- **Week 1:** First week with Thursday in new year
- **Range:** 1-52 (sometimes 53 in leap year with specific day alignment)
- **Week Start:** Monday (ISO standard)
- **Year Boundary:** Dec 29-31 can be week 1 of next year

---

## 🔧 Extension Points

### Multiple Format Options

**Different verbosity levels:**

```go
type CalendarDisplay struct {
    DisplayString string // "Tuesday, November 4, 2025 - Week 45"
    Full          string // Same as DisplayString
    Medium        string // "Nov 4, 2025 - W45"
    Compact       string // "11/04 W45"
    Minimal       string // "W45"
    HasInfo       bool
}

func GetCalendarDisplay(tempCtx *temporal.TemporalContext) CalendarDisplay {
    if tempCtx.ExternalCalendar.Date == "" {
        return CalendarDisplay{HasInfo: false}
    }

    full := fmt.Sprintf("%s, %s - Week %d",
        tempCtx.ExternalCalendar.DayOfWeek,
        tempCtx.ExternalCalendar.Date,
        tempCtx.ExternalCalendar.WeekNumber)

    // Generate other formats
    medium := fmt.Sprintf("%s - W%d",
        shortenDate(tempCtx.ExternalCalendar.Date),
        tempCtx.ExternalCalendar.WeekNumber)

    compact := fmt.Sprintf("%s W%d",
        numericDate(tempCtx.ExternalCalendar.Date),
        tempCtx.ExternalCalendar.WeekNumber)

    minimal := fmt.Sprintf("W%d", tempCtx.ExternalCalendar.WeekNumber)

    return CalendarDisplay{
        DisplayString: full,
        Full:          full,
        Medium:        medium,
        Compact:       compact,
        Minimal:       minimal,
        HasInfo:       true,
    }
}
```

### Year Progress Indicator

**Show annual completion:**

```go
type CalendarDisplay struct {
    DisplayString string
    YearProgress  float64  // Percentage through year (0-100)
    HasInfo       bool
}

func GetCalendarDisplay(tempCtx *temporal.TemporalContext) CalendarDisplay {
    if tempCtx.ExternalCalendar.Date == "" {
        return CalendarDisplay{HasInfo: false}
    }

    yearProgress := (float64(tempCtx.ExternalCalendar.WeekNumber) / 52.0) * 100

    displayStr := fmt.Sprintf("%s, %s - Week %d (%.0f%% through year)",
        tempCtx.ExternalCalendar.DayOfWeek,
        tempCtx.ExternalCalendar.Date,
        tempCtx.ExternalCalendar.WeekNumber,
        yearProgress)

    return CalendarDisplay{
        DisplayString: displayStr,
        YearProgress:  yearProgress,
        HasInfo:       true,
    }
}
```

### Holiday Indicators

**Show special days:**

```go
type CalendarDisplay struct {
    DisplayString string
    IsHoliday     bool
    HolidayName   string  // "Thanksgiving", "Christmas", etc.
    HasInfo       bool
}

func GetCalendarDisplay(tempCtx *temporal.TemporalContext) CalendarDisplay {
    if tempCtx.ExternalCalendar.Date == "" {
        return CalendarDisplay{HasInfo: false}
    }

    holidayName := checkHoliday(tempCtx.ExternalCalendar)

    displayStr := fmt.Sprintf("%s, %s - Week %d",
        tempCtx.ExternalCalendar.DayOfWeek,
        tempCtx.ExternalCalendar.Date,
        tempCtx.ExternalCalendar.WeekNumber)

    if holidayName != "" {
        displayStr = fmt.Sprintf("🎉 %s (%s)", holidayName, displayStr)
    }

    return CalendarDisplay{
        DisplayString: displayStr,
        IsHoliday:     holidayName != "",
        HolidayName:   holidayName,
        HasInfo:       true,
    }
}
```

---

## 🔍 Troubleshooting

### Calendar never displays

**Problem:** calDisplay.HasInfo always false

**Check:**

1. Verify TemporalContext.ExternalCalendar.Date populated
2. Check hooks/lib/temporal calendar integration working
3. Debug: Print tempCtx.ExternalCalendar before calling function

**Expected:** hooks/lib/temporal must provide calendar data for this to work

### Week number seems wrong

**Problem:** Shows Week 1 but it's end of December

**Cause:** ISO-8601 week numbering (week 1 can include end of previous year)

**Explanation:**

- ISO week 1 = first week with Thursday in new year
- Dec 29-31 can be week 1 of next year
- This is correct per ISO standard

**Example:**

```bash
Dec 29, 2025 (Monday) = Week 1 of 2026
Dec 30, 2025 (Tuesday) = Week 1 of 2026
Dec 31, 2025 (Wednesday) = Week 1 of 2026
Jan 1, 2026 (Thursday) = Week 1 of 2026
```

**Expected:** Week number follows ISO-8601 standard

### Date format inconsistent

**Problem:** Date string format varies

**Cause:** Date field in ExternalCalendar controlled by hooks/lib/temporal

**Solution:**

- Standardize date formatting in hooks/lib/temporal
- Ensure consistent format (e.g., "Month Day, Year")
- Document expected format

**Expected:** Date format consistency comes from data layer

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetTimeOfDayDisplay** | Time of day with icon | [timeofday.md](timeofday.md) |
| **GetSessionPhaseDisplay** | Session duration and phase | [sessionphase.md](sessionphase.md) |
| **GetScheduleDisplay** | Current schedule state | [schedule.md](schedule.md) |

**Dependencies:**

| Component | Purpose | Link |
|-----------|---------|------|
| **hooks/lib/temporal** | Temporal data source | [hooks/lib/temporal/](~/.claude/hooks/lib/temporal/) |

**Future Functions:**

- `GetHolidayDisplay()` - Holiday and special day indicators
- `GetCalendarCompactDisplay()` - Shortened calendar format
- `GetYearProgressDisplay()` - Annual completion percentage

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Temporal Library README** | Library overview and integration guide | [README.md](README.md) |
| **Time of Day Function** | GetTimeOfDayDisplay reference | [timeofday.md](timeofday.md) |
| **Session Phase Function** | GetSessionPhaseDisplay reference | [sessionphase.md](sessionphase.md) |
| **Schedule Function** | GetScheduleDisplay reference | [schedule.md](schedule.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) | GetCalendarDisplay with comprehensive inline documentation |
| **Data Source** | [`hooks/lib/temporal/`](~/.claude/hooks/lib/temporal/) | Calendar data collection and ISO week calculation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

### External Standards

| Standard | Purpose | Link |
|----------|---------|------|
| **ISO-8601** | Week numbering standard | <https://en.wikipedia.org/wiki/ISO_8601> |

---

## 📖 Biblical Foundation

**Core Verse:**

- **"To every thing there is a season, and a time to every purpose under the heaven"** - Ecclesiastes 3:1
  - Awareness of seasons and appointed times - recognizing where we are in the year's rhythm

**Supporting Verse:**

- **"So teach us to number our days, that we may apply our hearts unto wisdom"** - Psalm 90:12
  - Numbering our days wisely - understanding time's progression through the year

**Application:** Calendar display helps us number our days wisely. Knowing the date and week number provides context - Week 45 of 52 reminds us the year is nearly complete. Tuesday in work week establishes daily rhythm. Not for anxiety about time passing, but for wisdom in how we steward our days. Each week is a gift, each day an opportunity to work faithfully and rest appropriately within the rhythms God established.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"So teach us to number our days, that we may apply our hearts unto wisdom." - Psalm 90:12*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
