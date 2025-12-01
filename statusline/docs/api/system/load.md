<div align="center">

# GetLoadDisplay Function Reference

**Load Average Formatting with CPU-Relative Health Coloring**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Thresholds](https://img.shields.io/badge/Thresholds-50%25%20%7C%2080%25-orange?style=flat)
![Pure](https://img.shields.io/badge/Pure-No%20Side%20Effects-green?style=flat)

*Transforms load average data into formatted display with CPU-relative health coloring*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [🔄 Returns](#-returns) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetLoadDisplay Function
**Purpose:** Load average display formatting with CPU-relative health coloring
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
  - [None](#none)
- [🔄 Returns](#-returns)
  - [LoadDisplay](#loaddisplay)
- [🎯 Behavior](#-behavior)
  - [Health Threshold Calculation](#health-threshold-calculation)
  - [Display Examples](#display-examples)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [With Color](#with-color)
  - [Zero-State Handling](#zero-state-handling)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Testing](#testing)
- [🔧 Edge Cases](#-edge-cases)
  - [Zero Load](#zero-load)
  - [Load Exceeding CPU Count](#load-exceeding-cpu-count)
  - [Single CPU System](#single-cpu-system)
  - [High CPU Count System](#high-cpu-count-system)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [📊 Understanding Load Average](#-understanding-load-average)
  - [What Load Average Means](#what-load-average-means)
  - [Why Relative to CPU Count](#why-relative-to-cpu-count)
  - [Health Threshold Rationale](#health-threshold-rationale)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Load color seems wrong](#load-color-seems-wrong)
  - [HasInfo false on working system](#hasinfo-false-on-working-system)
  - [Want different thresholds](#want-different-thresholds)
- [🔧 Extension Points](#-extension-points)
  - [Multi-Level Health Indicators](#multi-level-health-indicators)
  - [Trend Indicators](#trend-indicators)
  - [Configurable Thresholds](#configurable-thresholds)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**GetLoadDisplay** formats load average data with CPU-relative health coloring suitable for statusline display.

> [!NOTE]
> **Pure presentation function.** Receives load data from system/lib/system, applies health thresholds relative to CPU count, outputs formatted display. Cannot fail, no side effects.

**Design:** Assessment vs Presentation separation - system/lib/system provides data (LoadAvg1, CPUCount), this function formats for display with CPU-relative health-based color coding.

**Library:** [`statusline/lib/system`](../../../lib/system/)

**Source Code:** [`lib/system/system.go`](../../../lib/system/system.go)

**Library README:** [System Library Documentation](README.md)

---

## 📝 Signature

```go
func GetLoadDisplay() LoadDisplay
```

**Parameters:** None - retrieves current system load average automatically

**Returns:** `LoadDisplay` structure with formatted data

---

## ⚙ Parameters

### None

Function takes no parameters - automatically retrieves current system load average via `system/lib/system.GetInfo()`.

**Data Source:** `/proc/loadavg` (via system library)

---

## 🔄 Returns

### LoadDisplay

Formatted display structure with all presentation elements.

```go
type LoadDisplay struct {
    LoadAvg float64 // 1-minute load average
    Color   string  // Terminal color code based on load/CPU ratio
    Icon    string  // Visual icon representing load (e.g., "⚡")
    HasInfo bool    // True if load data available, false if no data
}
```

**Field Details:**

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| **LoadAvg** | float64 | 1-minute load average from system | `1.25` |
| **Color** | string | Health-based color code (CPU-relative) | `"green"` / `"yellow"` / `"red"` |
| **Icon** | string | Visual load icon (lightning bolt) | `"⚡"` |
| **HasInfo** | bool | Data availability flag | `true` / `false` |

**Color Mapping (CPU-Relative):**

| Load/CPU Ratio | Color | Meaning |
|:-----:|:-----:|---------|
| ≤ 50% | `display.Green` | System has headroom |
| > 50% ≤ 80% | `display.Yellow` | Working but capacity exists |
| > 80% | `display.Red` | Approaching full capacity |

**Guarantee:** Always returns valid LoadDisplay (never nil, never errors)

---

## 🎯 Behavior

### Health Threshold Calculation

Load percentage relative to CPU count:

```go
sysInfo := syslib.GetInfo()

if sysInfo.LoadAvg1 <= 0 {
    return LoadDisplay{HasInfo: false}  // Zero-state
}

// Calculate load relative to CPU count
loadPercent := (sysInfo.LoadAvg1 / float64(sysInfo.CPUCount)) * 100

// Apply health thresholds
loadColor := display.Green
if loadPercent > 80 {
    loadColor = display.Red     // Stressed
} else if loadPercent > 50 {
    loadColor = display.Yellow  // Elevated
}

return LoadDisplay{
    LoadAvg: sysInfo.LoadAvg1,
    Color:   loadColor,
    Icon:    "⚡",
    HasInfo: true,
}
```

### Display Examples

**Healthy Load (Green):**

```go
// System: 4 CPUs, Load 1.5
display := GetLoadDisplay()
// → LoadDisplay{LoadAvg: 1.5, Color: "green", Icon: "⚡", HasInfo: true}
// Calculation: (1.5 / 4) * 100 = 37.5% (≤50% = healthy)
```

**Elevated Load (Yellow):**

```go
// System: 4 CPUs, Load 3.0
display := GetLoadDisplay()
// → LoadDisplay{LoadAvg: 3.0, Color: "yellow", Icon: "⚡", HasInfo: true}
// Calculation: (3.0 / 4) * 100 = 75% (>50% and ≤80% = elevated)
```

**Stressed Load (Red):**

```go
// System: 4 CPUs, Load 4.5
display := GetLoadDisplay()
// → LoadDisplay{LoadAvg: 4.5, Color: "red", Icon: "⚡", HasInfo: true}
// Calculation: (4.5 / 4) * 100 = 112.5% (>80% = stressed)
```

**No Load Data:**

```go
// System with no /proc/loadavg access
display := GetLoadDisplay()
// → LoadDisplay{LoadAvg: 0, Color: "", Icon: "", HasInfo: false}
```

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/system"

func buildStatusline() string {
    // Get load display
    loadDisplay := system.GetLoadDisplay()
    // → LoadAvg: 1.25

    var parts []string

    // Only add if load data available
    if loadDisplay.HasInfo {
        loadPart := fmt.Sprintf("%s %.2f", loadDisplay.Icon, loadDisplay.LoadAvg)
        parts = append(parts, loadPart)
        // → "⚡ 1.25"
    }

    return strings.Join(parts, " | ")
}
```

### With Color

```go
loadDisplay := system.GetLoadDisplay()

if loadDisplay.HasInfo {
    // Apply color for terminal output
    colored := fmt.Sprintf("%s%s %.2f%s",
        loadDisplay.Color,      // Start color
        loadDisplay.Icon,       // Icon
        loadDisplay.LoadAvg,    // Value
        display.Reset)          // Reset color
    fmt.Println(colored)
    // → Displays in green/yellow/red: ⚡ 1.25
}
```

### Zero-State Handling

```go
loadDisplay := system.GetLoadDisplay()

// Check HasInfo before using
if loadDisplay.HasInfo {
    fmt.Printf("Load: %.2f\n", loadDisplay.LoadAvg)
} else {
    fmt.Println("Load data unavailable")
}
```

### Full Statusline Integration

```go
func buildSystemHealth() string {
    var parts []string

    // Load
    loadDisplay := system.GetLoadDisplay()
    if loadDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%.2f", loadDisplay.LoadAvg))
    }

    // Memory
    memDisplay := system.GetMemoryDisplay()
    if memDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%.0f%% mem", memDisplay.Percent))
    }

    // Disk
    diskDisplay := system.GetDiskDisplay("/")
    if diskDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%.0f%% disk", diskDisplay.Percent))
    }

    return strings.Join(parts, " | ")
    // → "1.25 | 32% mem | 45% disk"
}
```

### Testing

```go
func TestLoadDisplay(t *testing.T) {
    loadDisplay := system.GetLoadDisplay()

    // Verify structure
    if loadDisplay.HasInfo {
        if loadDisplay.LoadAvg <= 0 {
            t.Error("LoadAvg should be > 0 when HasInfo is true")
        }
        if loadDisplay.Color == "" {
            t.Error("Color should be set when HasInfo is true")
        }
        if loadDisplay.Icon == "" {
            t.Error("Icon should be set when HasInfo is true")
        }
    }
}
```

---

## 🔧 Edge Cases

### Zero Load

```go
// System idle
display := GetLoadDisplay()
// Returns: LoadDisplay{HasInfo: false}
// Behavior: Zero load = no display (system truly idle or data unavailable)
```

### Load Exceeding CPU Count

```go
// System: 4 CPUs, Load 8.0
display := GetLoadDisplay()
// Returns: LoadDisplay{LoadAvg: 8.0, Color: "red", HasInfo: true}
// Calculation: (8.0 / 4) * 100 = 200% (>80% = red)
// Behavior: Load > CPU count indicates processes waiting for CPU
```

### Single CPU System

```go
// System: 1 CPU, Load 0.6
display := GetLoadDisplay()
// Returns: LoadDisplay{LoadAvg: 0.6, Color: "yellow", HasInfo: true}
// Calculation: (0.6 / 1) * 100 = 60% (>50% and ≤80% = yellow)
// Behavior: Single CPU systems reach elevated load faster
```

### High CPU Count System

```go
// System: 32 CPUs, Load 20.0
display := GetLoadDisplay()
// Returns: LoadDisplay{LoadAvg: 20.0, Color: "yellow", HasInfo: true}
// Calculation: (20.0 / 32) * 100 = 62.5% (>50% and ≤80% = yellow)
// Behavior: High CPU systems can sustain higher absolute load
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Data retrieval | +25 | Get load average from system lib |
| CPU-relative calculation | +25 | Calculate load percentage |
| Threshold application | +25 | Determine health color |
| Display structure creation | +25 | Build return value |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** Load retrieval with fallback guarantees valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(1) - /proc/loadavg read + calculation

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Data source** | /proc/loadavg | Linux kernel interface |
| **Read time** | <50 μs | Microseconds |
| **Calculation** | O(1) | Division and comparison |
| **Struct allocation** | ~50 bytes | Single LoadDisplay |
| **Execution time** | <100 μs | Total end-to-end |

**Memory:** Single LoadDisplay struct allocation (~50 bytes)

**Optimization:** Not needed - /proc read and struct allocation already optimal

---

## 📊 Understanding Load Average

### What Load Average Means

**Load average** = average number of processes in the run queue (running or waiting for CPU)

- **Load 1.0 on 1 CPU**: Fully utilized (one process always running)
- **Load 1.0 on 4 CPUs**: 25% utilized (one process running, 3 CPUs idle)
- **Load 4.0 on 4 CPUs**: Fully utilized (4 processes running)
- **Load 8.0 on 4 CPUs**: Overloaded (8 processes, 4 waiting)

### Why Relative to CPU Count

**Absolute load meaningless without context:**

- Load 4.0 on 4 CPUs = 100% (busy but not stressed)
- Load 4.0 on 32 CPUs = 12.5% (mostly idle)

**Health thresholds use load/CPU ratio:**

- Measures actual system capacity utilization
- Same thresholds work across different CPU counts

### Health Threshold Rationale

**50% threshold (yellow):**

- System has headroom but working
- Good time to investigate if sustained

**80% threshold (red):**

- System approaching full capacity
- Performance degradation likely
- Action may be needed

---

## 🔍 Troubleshooting

### Load color seems wrong

**Problem:** Load 2.0 shows yellow on 4 CPUs but expected green

**Check:**

1. Verify calculation: (2.0 / 4) * 100 = 50% (exactly at yellow threshold)
2. Threshold is > 50%, so 50% exactly is still green
3. 50.1% would be yellow

**Expected Behavior:** Thresholds are exclusive (>50%, >80%)

### HasInfo false on working system

**Problem:** Load display shows no data despite system running

**Check:**

1. Verify /proc/loadavg accessible: `cat /proc/loadavg`
2. Check system/lib/system implementation
3. Debug: Call syslib.GetInfo() directly

**Expected:** Linux systems should always have /proc/loadavg

### Want different thresholds

**Problem:** 50%/80% doesn't match your definition of "healthy"

**Solution:** Adjust thresholds in system.go GetLoadDisplay():

```go
// Current:
if loadPercent > 80 {
    loadColor = display.Red
} else if loadPercent > 50 {
    loadColor = display.Yellow
}

// Custom (e.g., 60%/90%):
if loadPercent > 90 {
    loadColor = display.Red
} else if loadPercent > 60 {
    loadColor = display.Yellow
}
```

---

## 🔧 Extension Points

### Multi-Level Health Indicators

**Future enhancement:**

```go
// Add warning level between yellow and red
if loadPercent > 95 {
    loadColor = display.Red       // Critical
} else if loadPercent > 80 {
    loadColor = display.Orange    // Warning
} else if loadPercent > 50 {
    loadColor = display.Yellow    // Elevated
}
```

### Trend Indicators

**Show if load increasing/decreasing:**

```go
type LoadDisplay struct {
    LoadAvg    float64
    Load5      float64  // 5-minute average
    Trend      string   // "↑" increasing, "↓" decreasing, "→" stable
    Color      string
    Icon       string
    HasInfo    bool
}

func GetLoadDisplay() LoadDisplay {
    sysInfo := syslib.GetInfo()

    trend := "→"
    if sysInfo.LoadAvg1 > sysInfo.LoadAvg5 * 1.1 {
        trend = "↑"  // Load increasing
    } else if sysInfo.LoadAvg1 < sysInfo.LoadAvg5 * 0.9 {
        trend = "↓"  // Load decreasing
    }

    // ... rest of function
}
```

### Configurable Thresholds

```go
type LoadThresholds struct {
    Warning  float64  // 50.0
    Critical float64  // 80.0
}

func GetLoadDisplayWithThresholds(thresholds LoadThresholds) LoadDisplay {
    // Apply custom thresholds instead of hardcoded 50%/80%
}
```

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetMemoryDisplay** | Memory usage with percentage-based health coloring | [memory.md](memory.md) |
| **GetDiskDisplay** | Disk capacity with usage-based health coloring | [disk.md](disk.md) |

**Dependencies:**

| Component | Purpose | Location |
|-----------|---------|----------|
| **system/lib/system.GetInfo()** | System data collection (load, CPU count) | [system/lib/system](~/.claude/cpi-si/system/lib/system/) |

**Future Functions:**

- `GetLoadTrendDisplay()` - Load with 5/15-minute trend
- `GetLoadBreakdownDisplay()` - Separate 1/5/15-minute averages

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **System Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetMemoryDisplay Function** | Memory usage formatting | [memory.md](memory.md) |
| **GetDiskDisplay Function** | Disk capacity formatting | [disk.md](disk.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/system/system.go`](../../../lib/system/system.go) | GetLoadDisplay with comprehensive inline documentation |
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

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Order and clarity in communication - showing system health clearly and truthfully

**Supporting Verse:**

- **"A faithful witness will not lie: but a false witness will utter lies"** - Proverbs 14:5
  - Truth-telling through accurate metrics and health indicators

**Application:** Show actual load average (1.25) with truthful health indicator (green/yellow/red) without alarm or minimization. Faithful witness to system state.

**Why This Matters:** Load average is meaningless without CPU context. Load 4.0 on 4 CPUs = 100% (working hard but not overloaded). Load 4.0 on 32 CPUs = 12.5% (mostly idle). This library tells the truth by showing CPU-relative load percentage.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
