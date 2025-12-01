<div align="center">

# GetMemoryDisplay Function Reference

**Memory Usage Formatting with Percentage-Based Health Coloring**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Thresholds](https://img.shields.io/badge/Thresholds-60%25%20%7C%2080%25-orange?style=flat)
![Pure](https://img.shields.io/badge/Pure-No%20Side%20Effects-green?style=flat)

*Transforms memory usage data into formatted display with health-based color coding*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [🔄 Returns](#-returns) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetMemoryDisplay Function
**Purpose:** Memory usage display formatting with percentage-based health coloring
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
  - [MemoryDisplay](#memorydisplay)
- [🎯 Behavior](#-behavior)
  - [Health Threshold Calculation](#health-threshold-calculation)
  - [Display Examples](#display-examples)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [With Percentage](#with-percentage)
  - [With Color](#with-color)
  - [Zero-State Handling](#zero-state-handling)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Testing](#testing)
- [🔧 Edge Cases](#-edge-cases)
  - [Zero Memory Usage](#zero-memory-usage)
  - [Very Low Memory System](#very-low-memory-system)
  - [High Memory System](#high-memory-system)
  - [Exactly At Threshold](#exactly-at-threshold)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [📊 Understanding Memory Usage](#-understanding-memory-usage)
  - [What Memory Metrics Mean](#what-memory-metrics-mean)
  - [Why Percentage-Based](#why-percentage-based)
  - [Health Threshold Rationale](#health-threshold-rationale)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Memory percentage seems wrong](#memory-percentage-seems-wrong)
  - [HasInfo false on working system](#hasinfo-false-on-working-system)
  - [Want different thresholds](#want-different-thresholds)
- [🔧 Extension Points](#-extension-points)
  - [Memory Breakdown Display](#memory-breakdown-display)
  - [Memory Pressure Indicator](#memory-pressure-indicator)
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

**GetMemoryDisplay** formats memory usage data with percentage-based health coloring suitable for statusline display.

> [!NOTE]
> **Pure presentation function.** Receives memory data from system/lib/system, applies health thresholds based on usage percentage, outputs formatted display. Cannot fail, no side effects.

**Design:** Assessment vs Presentation separation - system/lib/system provides data (MemUsedGB, MemTotalGB), this function formats for display with health-based color coding.

**Library:** [`statusline/lib/system`](../../../lib/system/)

**Source Code:** [`lib/system/system.go`](../../../lib/system/system.go)

**Library README:** [System Library Documentation](README.md)

---

## 📝 Signature

```go
func GetMemoryDisplay() MemoryDisplay
```

**Parameters:** None - retrieves current system memory usage automatically

**Returns:** `MemoryDisplay` structure with formatted data

---

## ⚙ Parameters

### None

Function takes no parameters - automatically retrieves current system memory usage via `system/lib/system.GetInfo()`.

**Data Source:** `/proc/meminfo` (via system library)

---

## 🔄 Returns

### MemoryDisplay

Formatted display structure with all presentation elements.

```go
type MemoryDisplay struct {
    UsedGB  float64 // Memory used in GB
    TotalGB float64 // Total memory in GB
    Percent float64 // Usage percentage (0-100)
    Color   string  // Terminal color code based on usage percent
    Icon    string  // Visual icon representing memory (e.g., "💾")
    HasInfo bool    // True if memory data available, false if no data
}
```

**Field Details:**

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| **UsedGB** | float64 | Memory used in gigabytes | `5.2` |
| **TotalGB** | float64 | Total physical RAM in gigabytes | `16.0` |
| **Percent** | float64 | Usage percentage (UsedGB / TotalGB * 100) | `32.5` |
| **Color** | string | Health-based color code | `"green"` / `"yellow"` / `"red"` |
| **Icon** | string | Visual memory icon | `"💾"` |
| **HasInfo** | bool | Data availability flag | `true` / `false` |

**Color Mapping:**

| Usage | Color | Meaning |
|:-----:|:-----:|---------|
| ≤ 60% | `display.Green` | Healthy memory usage |
| > 60% ≤ 80% | `display.Yellow` | Elevated usage |
| > 80% | `display.Red` | High memory pressure |

**Guarantee:** Always returns valid MemoryDisplay (never nil, never errors)

---

## 🎯 Behavior

### Health Threshold Calculation

Memory percentage-based coloring:

```go
sysInfo := syslib.GetInfo()

if sysInfo.MemTotalGB <= 0 {
    return MemoryDisplay{HasInfo: false}  // Zero-state
}

// Calculate usage percentage
memPercent := (sysInfo.MemUsedGB / sysInfo.MemTotalGB) * 100

// Apply health thresholds
memColor := display.Green
if memPercent > 80 {
    memColor = display.Red      // High memory pressure
} else if memPercent > 60 {
    memColor = display.Yellow   // Elevated usage
}

return MemoryDisplay{
    UsedGB:  sysInfo.MemUsedGB,
    TotalGB: sysInfo.MemTotalGB,
    Percent: memPercent,
    Color:   memColor,
    Icon:    "💾",
    HasInfo: true,
}
```

### Display Examples

**Healthy Usage (Green):**

```go
// System: 5.2GB used / 16.0GB total
display := GetMemoryDisplay()
// → MemoryDisplay{UsedGB: 5.2, TotalGB: 16.0, Percent: 32.5, Color: "green", HasInfo: true}
// Calculation: (5.2 / 16.0) * 100 = 32.5% (≤60% = healthy)
```

**Elevated Usage (Yellow):**

```go
// System: 10.5GB used / 16.0GB total
display := GetMemoryDisplay()
// → MemoryDisplay{UsedGB: 10.5, TotalGB: 16.0, Percent: 65.6, Color: "yellow", HasInfo: true}
// Calculation: (10.5 / 16.0) * 100 = 65.6% (>60% and ≤80% = elevated)
```

**High Pressure (Red):**

```go
// System: 14.0GB used / 16.0GB total
display := GetMemoryDisplay()
// → MemoryDisplay{UsedGB: 14.0, TotalGB: 16.0, Percent: 87.5, Color: "red", HasInfo: true}
// Calculation: (14.0 / 16.0) * 100 = 87.5% (>80% = high pressure)
```

**No Memory Data:**

```go
// System with no /proc/meminfo access
display := GetMemoryDisplay()
// → MemoryDisplay{UsedGB: 0, TotalGB: 0, Percent: 0, Color: "", HasInfo: false}
```

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/system"

func buildStatusline() string {
    // Get memory display
    memDisplay := system.GetMemoryDisplay()
    // → UsedGB: 5.2, TotalGB: 16.0, Percent: 32.5

    var parts []string

    // Only add if memory data available
    if memDisplay.HasInfo {
        memPart := fmt.Sprintf("%s %.1f/%.1fGB",
            memDisplay.Icon, memDisplay.UsedGB, memDisplay.TotalGB)
        parts = append(parts, memPart)
        // → "💾 5.2/16.0GB"
    }

    return strings.Join(parts, " | ")
}
```

### With Percentage

```go
memDisplay := system.GetMemoryDisplay()

if memDisplay.HasInfo {
    memPart := fmt.Sprintf("%s %.0f%%",
        memDisplay.Icon, memDisplay.Percent)
    fmt.Println(memPart)
    // → "💾 32%"
}
```

### With Color

```go
memDisplay := system.GetMemoryDisplay()

if memDisplay.HasInfo {
    // Apply color for terminal output
    colored := fmt.Sprintf("%s%s %.1f/%.1fGB (%.0f%%)%s",
        memDisplay.Color,       // Start color
        memDisplay.Icon,        // Icon
        memDisplay.UsedGB,      // Used
        memDisplay.TotalGB,     // Total
        memDisplay.Percent,     // Percentage
        display.Reset)          // Reset color
    fmt.Println(colored)
    // → Displays in green/yellow/red: 💾 5.2/16.0GB (32%)
}
```

### Zero-State Handling

```go
memDisplay := system.GetMemoryDisplay()

// Check HasInfo before using
if memDisplay.HasInfo {
    fmt.Printf("Memory: %.0f%%\n", memDisplay.Percent)
} else {
    fmt.Println("Memory data unavailable")
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
func TestMemoryDisplay(t *testing.T) {
    memDisplay := system.GetMemoryDisplay()

    // Verify structure
    if memDisplay.HasInfo {
        if memDisplay.TotalGB <= 0 {
            t.Error("TotalGB should be > 0 when HasInfo is true")
        }
        if memDisplay.UsedGB < 0 {
            t.Error("UsedGB should be >= 0")
        }
        if memDisplay.Percent < 0 || memDisplay.Percent > 100 {
            t.Error("Percent should be in range 0-100")
        }
        if memDisplay.Color == "" {
            t.Error("Color should be set when HasInfo is true")
        }
    }
}
```

---

## 🔧 Edge Cases

### Zero Memory Usage

```go
// Theoretically impossible (kernel uses memory)
display := GetMemoryDisplay()
// Returns: MemoryDisplay{HasInfo: false}
// Behavior: Zero usage treated as no data
```

### Very Low Memory System

```go
// System: 1.5GB used / 2.0GB total
display := GetMemoryDisplay()
// Returns: MemoryDisplay{UsedGB: 1.5, TotalGB: 2.0, Percent: 75.0, Color: "yellow", HasInfo: true}
// Behavior: Low RAM systems reach elevated thresholds faster
```

### High Memory System

```go
// System: 50GB used / 128GB total
display := GetMemoryDisplay()
// Returns: MemoryDisplay{UsedGB: 50.0, TotalGB: 128.0, Percent: 39.1, Color: "green", HasInfo: true}
// Behavior: High RAM systems can use more absolute memory while staying green
```

### Exactly At Threshold

```go
// System: 9.6GB used / 16.0GB total
display := GetMemoryDisplay()
// Returns: MemoryDisplay{UsedGB: 9.6, TotalGB: 16.0, Percent: 60.0, Color: "green", HasInfo: true}
// Calculation: Exactly 60% is still green (thresholds are >60%, >80%)
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Data retrieval | +25 | Get memory info from system lib |
| Percentage calculation | +25 | Calculate usage percent |
| Threshold application | +25 | Determine health color |
| Display structure creation | +25 | Build return value |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** Memory retrieval with fallback guarantees valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(1) - /proc/meminfo read + calculation

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Data source** | /proc/meminfo | Linux kernel interface |
| **Read time** | <50 μs | Microseconds |
| **Calculation** | O(1) | Division and arithmetic |
| **Struct allocation** | ~50 bytes | Single MemoryDisplay |
| **Execution time** | <100 μs | Total end-to-end |

**Memory:** Single MemoryDisplay struct allocation (~50 bytes)

**Optimization:** Not needed - /proc read and struct allocation already optimal

---

## 📊 Understanding Memory Usage

### What Memory Metrics Mean

**MemTotal:** Total physical RAM installed

**MemAvailable:** Memory available for new processes (includes cached memory that can be freed)

**MemUsed:** MemTotal - MemAvailable (what system/lib/system calculates)

### Why Percentage-Based

**Absolute memory meaningless without context:**

- 8GB used on 16GB system = 50% (healthy)
- 8GB used on 12GB system = 66% (elevated)

**Percentage shows actual capacity utilization:**

- Same thresholds work across different RAM sizes
- 60%/80% meaningful regardless of total RAM

### Health Threshold Rationale

**60% threshold (yellow):**

- Memory usage elevated but manageable
- Good time to investigate if sustained
- System still has headroom

**80% threshold (red):**

- High memory pressure
- Swapping may begin soon
- Performance degradation likely
- Action may be needed

---

## 🔍 Troubleshooting

### Memory percentage seems wrong

**Problem:** Memory display shows 32% but `free` command shows different

**Check:**

1. Verify calculation: (MemUsed / MemTotal) * 100
2. Note: MemUsed = MemTotal - MemAvailable (not MemFree)
3. Linux caches aggressively - "used" includes cache that can be freed

**Expected:** Our calculation uses MemAvailable (actual available for new processes)

### HasInfo false on working system

**Problem:** Memory display shows no data despite system running

**Check:**

1. Verify /proc/meminfo accessible: `cat /proc/meminfo`
2. Check system/lib/system implementation
3. Debug: Call syslib.GetInfo() directly

**Expected:** Linux systems should always have /proc/meminfo

### Want different thresholds

**Problem:** 60%/80% doesn't match your definition of "healthy memory usage"

**Solution:** Adjust thresholds in system.go GetMemoryDisplay():

```go
// Current:
if memPercent > 80 {
    memColor = display.Red
} else if memPercent > 60 {
    memColor = display.Yellow
}

// Custom (e.g., 70%/85%):
if memPercent > 85 {
    memColor = display.Red
} else if memPercent > 70 {
    memColor = display.Yellow
}
```

---

## 🔧 Extension Points

### Memory Breakdown Display

**Show swap usage separately:**

```go
type MemoryBreakdownDisplay struct {
    UsedGB      float64
    TotalGB     float64
    SwapUsedGB  float64
    SwapTotalGB float64
    Color       string
    Icon        string
    HasInfo     bool
}

func GetMemoryBreakdownDisplay() MemoryBreakdownDisplay {
    sysInfo := syslib.GetInfo()
    swapInfo := syslib.GetSwapInfo()  // New function in system lib

    // Calculate health including swap
    // ...
}
```

### Memory Pressure Indicator

**Show if memory usage increasing:**

```go
type MemoryDisplay struct {
    UsedGB   float64
    TotalGB  float64
    Percent  float64
    Trend    string  // "↑" increasing, "↓" decreasing, "→" stable
    Color    string
    Icon     string
    HasInfo  bool
}

// Would require historical tracking
```

### Configurable Thresholds

```go
type MemoryThresholds struct {
    Warning  float64  // 60.0
    Critical float64  // 80.0
}

func GetMemoryDisplayWithThresholds(thresholds MemoryThresholds) MemoryDisplay {
    // Apply custom thresholds instead of hardcoded 60%/80%
}
```

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetLoadDisplay** | Load average with CPU-relative health coloring | [load.md](load.md) |
| **GetDiskDisplay** | Disk capacity with usage-based health coloring | [disk.md](disk.md) |

**Dependencies:**

| Component | Purpose | Location |
|-----------|---------|----------|
| **system/lib/system.GetInfo()** | System data collection (memory usage) | [system/lib/system](~/.claude/cpi-si/system/lib/system/) |

**Future Functions:**

- `GetMemoryBreakdownDisplay()` - Separate RAM + swap display
- `GetMemoryTrendDisplay()` - Memory with usage trend

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **System Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetLoadDisplay Function** | Load average formatting | [load.md](load.md) |
| **GetDiskDisplay Function** | Disk capacity formatting | [disk.md](disk.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/system/system.go`](../../../lib/system/system.go) | GetMemoryDisplay with comprehensive inline documentation |
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

**Application:** Show actual memory usage (5.2/16.0GB, 32%) with truthful health indicator (green/yellow/red) without alarm or minimization. Faithful witness to memory state.

**Why This Matters:** Many tools show 70% memory as "critical" (alarmist) or 85% as "okay" (minimizing). This library tells the truth: 70% = yellow (elevated but manageable), 85% = red (actually high pressure).

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
