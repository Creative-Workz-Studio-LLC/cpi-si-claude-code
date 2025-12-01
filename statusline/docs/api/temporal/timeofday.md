<div align="center">

# GetTimeOfDayDisplay Function Reference

**Time of Day Formatting with Icon Selection**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Icons](https://img.shields.io/badge/Icons-4%20Time%20Periods-orange?style=flat)
![Pure](https://img.shields.io/badge/Pure-No%20Side%20Effects-green?style=flat)

*Transforms time of day data into formatted display with appropriate icon selection*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [🔄 Returns](#-returns) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetTimeOfDayDisplay Function
**Purpose:** Time of day display formatting with icon selection
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📖 Overview](#-overview)
- [📝 Signature](#-signature)
- [⚙ Parameters](#-parameters)
  - [tempCtx (*temporal.TemporalContext)](#tempctx-temporaltemporalcontext)
- [🔄 Returns](#-returns)
  - [TimeOfDayDisplay](#timeofdaydisplay)
- [🎯 Behavior](#-behavior)
  - [Icon Selection Logic](#icon-selection-logic)
  - [Display Examples](#display-examples)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [Icon Only](#icon-only)
  - [With Color](#with-color)
  - [Conditional Display](#conditional-display)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Testing](#testing)
- [🔧 Edge Cases](#-edge-cases)
  - [Empty Time of Day](#empty-time-of-day)
  - [Unknown Time Value](#unknown-time-value)
  - [Nil Context](#nil-context)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [📊 Understanding Time of Day](#-understanding-time-of-day)
  - [What Time of Day Represents](#what-time-of-day-represents)
  - [Why Icon-Based](#why-icon-based)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Icon seems wrong for current time](#icon-seems-wrong-for-current-time)
  - [Label is empty](#label-is-empty)
  - [Want different icons](#want-different-icons)
  - [Want time-specific colors](#want-time-specific-colors)
- [🔧 Extension Points](#-extension-points)
  - [Adding More Time Periods](#adding-more-time-periods)
  - [Seasonal Icons](#seasonal-icons)
  - [Configurable Icons](#configurable-icons)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetTimeOfDayDisplay** formats time of day data with appropriate icon selection suitable for statusline display.

> [!NOTE]
> **Pure presentation function.** Receives time of day from hooks/lib/temporal TemporalContext, selects appropriate icon based on time period, outputs formatted display. Cannot fail, no side effects.

**Design:** Assessment vs Presentation separation - hooks/lib/temporal provides data (TimeOfDay), this function formats for display with icon selection.

**Library:** [`statusline/lib/temporal`](../../../lib/temporal/)

**Source Code:** [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go)

**Library README:** [Temporal Library Documentation](README.md)

---

## 📝 Signature

```go
func GetTimeOfDayDisplay(tempCtx *temporal.TemporalContext) TimeOfDayDisplay
```

**Parameters:** `tempCtx` - Pointer to temporal context

**Returns:** `TimeOfDayDisplay` structure with formatted data

---

## ⚙ Parameters

### tempCtx (*temporal.TemporalContext)

Pointer to temporal context containing time of day information.

**Source:** `hooks/lib/temporal.GetTemporalContext()`

**Required Fields:**

| Field | Type | Description |
|-------|------|-------------|
| **ExternalTime.TimeOfDay** | string | Time period label ("morning", "afternoon", "evening", "night") |

**Example:**

```go
tempCtx, _ := temporal.GetTemporalContext()
// tempCtx.ExternalTime.TimeOfDay = "morning"

todDisplay := GetTimeOfDayDisplay(tempCtx)
```

**Validation:** Accepts any TemporalContext (never fails)

---

## 🔄 Returns

### TimeOfDayDisplay

Formatted display structure with time of day presentation elements.

```go
type TimeOfDayDisplay struct {
    Icon  string // Visual icon representing time of day
    Label string // Time of day label (morning/afternoon/evening/night)
    Color string // Terminal color code
}
```

**Field Details:**

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| **Icon** | string | Visual icon for time period | `"🌅"` (morning) |
| **Label** | string | Time of day label | `"morning"` |
| **Color** | string | Terminal color code | `"yellow"` |

**Icon Mapping:**

| Time Period | Icon | Meaning |
|-------------|:----:|---------|
| **Morning** | 🌅 | Sunrise - new beginning, fresh start |
| **Afternoon** | ☀️ | Sun - peak activity, full light |
| **Evening** | 🌆 | Cityscape/Sunset - winding down, transition |
| **Night** | 🌙 | Moon - rest, darkness, completion |
| **Unknown/Default** | ☀️ | Default to afternoon icon |

**Guarantee:** Always returns valid TimeOfDayDisplay (never nil, never errors)

---

## 🎯 Behavior

### Icon Selection Logic

Time-based icon mapping via switch statement:

```go
timeOfDayIcon := "☀️"  // Default (afternoon/day)

switch tempCtx.ExternalTime.TimeOfDay {
case "morning":
    timeOfDayIcon = "🌅"  // Sunrise
case "afternoon":
    timeOfDayIcon = "☀️"  // Sun
case "evening":
    timeOfDayIcon = "🌆"  // Sunset
case "night":
    timeOfDayIcon = "🌙"  // Moon
}

return TimeOfDayDisplay{
    Icon:  timeOfDayIcon,
    Label: tempCtx.ExternalTime.TimeOfDay,
    Color: display.Yellow,
}
```

### Display Examples

**Morning:**

```go
// TemporalContext: TimeOfDay = "morning"
display := GetTimeOfDayDisplay(tempCtx)
// → TimeOfDayDisplay{Icon: "🌅", Label: "morning", Color: "yellow"}
```

**Afternoon:**

```go
// TemporalContext: TimeOfDay = "afternoon"
display := GetTimeOfDayDisplay(tempCtx)
// → TimeOfDayDisplay{Icon: "☀️", Label: "afternoon", Color: "yellow"}
```

**Evening:**

```go
// TemporalContext: TimeOfDay = "evening"
display := GetTimeOfDayDisplay(tempCtx)
// → TimeOfDayDisplay{Icon: "🌆", Label: "evening", Color: "yellow"}
```

**Night:**

```go
// TemporalContext: TimeOfDay = "night"
display := GetTimeOfDayDisplay(tempCtx)
// → TimeOfDayDisplay{Icon: "🌙", Label: "night", Color: "yellow"}
```

**Empty/Unknown:**

```go
// TemporalContext: TimeOfDay = ""
display := GetTimeOfDayDisplay(tempCtx)
// → TimeOfDayDisplay{Icon: "☀️", Label: "", Color: "yellow"}
// Defaults to afternoon icon, empty label
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

    // Get time of day display
    todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
    // → Icon: "🌅", Label: "morning"

    var parts []string

    // Add to statusline
    parts = append(parts, fmt.Sprintf("%s %s", todDisplay.Icon, todDisplay.Label))
    // → "🌅 morning"

    return strings.Join(parts, " | ")
}
```

### Icon Only

```go
todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)

// Use just the icon
parts = append(parts, todDisplay.Icon)
// → "🌅"
```

### With Color

```go
todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)

// Apply color for terminal output
colored := fmt.Sprintf("%s%s %s%s",
    todDisplay.Color,      // Start color
    todDisplay.Icon,       // Icon
    todDisplay.Label,      // Label
    display.Reset)         // Reset color

fmt.Println(colored)
// → Displays in yellow: 🌅 morning
```

### Conditional Display

```go
todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)

// Only display if label available
if todDisplay.Label != "" {
    parts = append(parts, fmt.Sprintf("%s %s", todDisplay.Icon, todDisplay.Label))
}
```

### Full Statusline Integration

```go
func buildTemporalSection(tempCtx *temporal.TemporalContext) string {
    var parts []string

    // Time of day
    todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
    parts = append(parts, fmt.Sprintf("%s %s", todDisplay.Icon, todDisplay.Label))

    // Session phase
    phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
    if phaseDisplay.Duration != "" {
        parts = append(parts, fmt.Sprintf("%s (%s)", phaseDisplay.Duration, phaseDisplay.Phase))
    }

    return strings.Join(parts, " | ")
    // → "🌅 morning | 45m (fresh)"
}
```

### Testing

```go
func TestTimeOfDayDisplay(t *testing.T) {
    testCases := []struct {
        timeOfDay    string
        expectedIcon string
    }{
        {"morning", "🌅"},
        {"afternoon", "☀️"},
        {"evening", "🌆"},
        {"night", "🌙"},
        {"", "☀️"},  // Default
    }

    for _, tc := range testCases {
        ctx := &temporal.TemporalContext{
            ExternalTime: temporal.ExternalTime{TimeOfDay: tc.timeOfDay},
        }

        display := GetTimeOfDayDisplay(ctx)

        if display.Icon != tc.expectedIcon {
            t.Errorf("TimeOfDay %q: expected icon %q, got %q",
                tc.timeOfDay, tc.expectedIcon, display.Icon)
        }
    }
}
```

---

## 🔧 Edge Cases

### Empty Time of Day

```go
// TemporalContext with empty TimeOfDay
emptyCtx := &temporal.TemporalContext{
    ExternalTime: temporal.ExternalTime{TimeOfDay: ""},
}

display := GetTimeOfDayDisplay(emptyCtx)
// Returns: TimeOfDayDisplay{Icon: "☀️", Label: "", Color: "yellow"}
// Behavior: Default icon (afternoon), empty label
```

### Unknown Time Value

```go
// TemporalContext with unrecognized TimeOfDay
unknownCtx := &temporal.TemporalContext{
    ExternalTime: temporal.ExternalTime{TimeOfDay: "twilight"},
}

display := GetTimeOfDayDisplay(unknownCtx)
// Returns: TimeOfDayDisplay{Icon: "☀️", Label: "twilight", Color: "yellow"}
// Behavior: Default icon, preserves original label
```

### Nil Context

```go
// Nil pointer (edge case - should not happen in practice)
display := GetTimeOfDayDisplay(nil)
// Returns: TimeOfDayDisplay{Icon: "☀️", Label: "", Color: "yellow"}
// Behavior: Graceful degradation, default icon
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Icon selection | +50 | Switch statement execution |
| Display structure creation | +50 | Build return value |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** Icon selection with default fallback guarantees valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(1) - single switch statement

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Switch statement** | O(1) | String comparison |
| **Struct allocation** | ~50 bytes | Single TimeOfDayDisplay |
| **Execution time** | <1 μs | Microsecond |

**Memory:** Single TimeOfDayDisplay struct allocation (~50 bytes)

**Optimization:** Not needed - switch statement already optimal

---

## 📊 Understanding Time of Day

### What Time of Day Represents

**Time of day** = Current period of the day based on external system time

| Period | Typical Hours | Characteristics |
|--------|:-------------:|-----------------|
| **Morning** | 6 AM - 12 PM | Fresh start, rising energy, new beginning |
| **Afternoon** | 12 PM - 6 PM | Peak activity, full light, prime work time |
| **Evening** | 6 PM - 10 PM | Winding down, transition, reflection |
| **Night** | 10 PM - 6 AM | Rest, darkness, restoration |

*Exact boundaries determined by hooks/lib/temporal implementation*

### Why Icon-Based

**Visual recognition advantages:**

- Icons provide instant recognition without reading text
- Universal understanding across languages
- Space-efficient in statusline
- Aesthetic appeal
- Quick cognitive processing

**Icon symbolism:**

- 🌅 Sunrise = New beginning, hope, fresh start
- ☀️ Sun = Peak activity, brightness, energy
- 🌆 Cityscape/Sunset = Transition, completion, beauty
- 🌙 Moon = Rest, peace, renewal

---

## 🔍 Troubleshooting

### Icon seems wrong for current time

**Problem:** Displaying afternoon icon (☀️) but it's morning

**Check:**

1. Verify TemporalContext.ExternalTime.TimeOfDay value
2. Check hooks/lib/temporal time classification logic
3. Debug: Print tempCtx.ExternalTime.TimeOfDay before calling function

**Expected:** Icon matches TimeOfDay field value, not necessarily actual clock time (depends on temporal lib implementation)

### Label is empty

**Problem:** Icon displays but label is blank

**Cause:** ExternalTime.TimeOfDay field empty in TemporalContext

**Solution:**

- Verify hooks/lib/temporal populating TimeOfDay
- Call `temporal.GetTemporalContext()` to check data

**Expected:** If hooks lib doesn't provide data, library defaults gracefully

### Want different icons

**Problem:** Prefer different emoji for time periods

**Solution:** Modify icon selection in temporal.go GetTimeOfDayDisplay():

```go
// Current:
case "morning":
    timeOfDayIcon = "🌅"

// Custom:
case "morning":
    timeOfDayIcon = "🌄"  // Sunrise over mountains
```

### Want time-specific colors

**Problem:** All times show same yellow color

**Solution:** Add color logic in temporal.go:

```go
timeColor := display.Yellow  // Default

switch tempCtx.ExternalTime.TimeOfDay {
case "morning":
    timeColor = display.Green  // Fresh/new
case "evening":
    timeColor = display.Blue   // Calm/winding down
case "night":
    timeColor = display.Gray   // Rest/low activity
}

return TimeOfDayDisplay{
    Icon:  timeOfDayIcon,
    Label: tempCtx.ExternalTime.TimeOfDay,
    Color: timeColor,
}
```

---

## 🔧 Extension Points

### Adding More Time Periods

**Finer granularity:**

```go
// Add more specific times
case "dawn":
    timeOfDayIcon = "🌄"  // Pre-sunrise
case "dusk":
    timeOfDayIcon = "🌇"  // Post-sunset
case "midnight":
    timeOfDayIcon = "🌃"  // Deep night
```

### Seasonal Icons

**Contextual icons based on season:**

```go
func GetTimeOfDayDisplay(tempCtx *temporal.TemporalContext) TimeOfDayDisplay {
    season := tempCtx.ExternalCalendar.Season  // Would need this field

    if tempCtx.ExternalTime.TimeOfDay == "morning" {
        switch season {
        case "winter":
            return TimeOfDayDisplay{Icon: "❄️"}
        case "spring":
            return TimeOfDayDisplay{Icon: "🌸"}
        case "summer":
            return TimeOfDayDisplay{Icon: "🌅"}
        case "fall":
            return TimeOfDayDisplay{Icon: "🍂"}
        }
    }
    // ... rest of logic
}
```

### Configurable Icons

```go
type TimeOfDayIcons struct {
    Morning   string
    Afternoon string
    Evening   string
    Night     string
}

func GetTimeOfDayDisplayWithIcons(tempCtx *temporal.TemporalContext, icons TimeOfDayIcons) TimeOfDayDisplay {
    // Use custom icons instead of hardcoded defaults
}
```

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetSessionPhaseDisplay** | Session duration and phase classification | [sessionphase.md](sessionphase.md) |
| **GetScheduleDisplay** | Current schedule state formatting | [schedule.md](schedule.md) |
| **GetCalendarDisplay** | Date and week information formatting | [calendar.md](calendar.md) |

**Dependencies:**

| Component | Purpose | Location |
|-----------|---------|----------|
| **hooks/lib/temporal.GetTemporalContext()** | Temporal data source | [hooks/lib/temporal](~/.claude/hooks/lib/temporal/) |

**Future Functions:**

- `GetCircadianPhaseDisplay()` - Biological rhythm indicator
- `GetCustomTimeDisplay()` - User-configurable time periods

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Temporal Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetSessionPhaseDisplay Function** | Session phase formatting | [sessionphase.md](sessionphase.md) |
| **GetScheduleDisplay Function** | Schedule formatting | [schedule.md](schedule.md) |
| **GetCalendarDisplay Function** | Calendar formatting | [calendar.md](calendar.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/temporal/temporal.go`](../../../lib/temporal/temporal.go) | GetTimeOfDayDisplay with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7-library ecosystem |

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
  - Awareness of time and seasons - recognizing that different times call for different approaches

**Supporting Verse:**

- **"Weeping may endure for a night, but joy cometh in the morning"** - Psalm 30:5
  - Each time period has its purpose and promise

**Application:** Visual representation of time of day reminds us that each period has its purpose. Morning brings fresh starts and new mercies (Lamentations 3:22-23), evening brings rest and reflection. God designed rhythms into creation - honoring them leads to wisdom.

**Why This Matters:** Modern culture often treats all hours as identical, pushing against natural rhythms. This display acknowledges God's design: mornings for beginning, afternoons for working, evenings for winding down, nights for resting. Not productivity optimization - wisdom in honoring the season.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
