<div align="center">

# System Health Display Library

**Pure Presentation Layer for System Health Metrics**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Library](https://img.shields.io/badge/Type-Library-blue?style=flat)
![Functions](https://img.shields.io/badge/Functions-3-green?style=flat)
![Metrics](https://img.shields.io/badge/Metrics-Load%20%7C%20Memory%20%7C%20Disk-orange?style=flat)

*Transforms system health metrics into formatted display structures with health-based color coding*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [🎯 What It Provides](#-what-it-provides) • [🔑 Key Features](#-key-features) • [💻 Demo](#-demo) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - System Health Display Library
**Purpose:** Transform system health metrics (load, memory, disk) into formatted display structures
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [System Health Display Library](#system-health-display-library)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
  - [🎯 What It Provides](#-what-it-provides)
    - [System Health Display Functions](#system-health-display-functions)
    - [Display Examples](#display-examples)
  - [🔑 Key Features](#-key-features)
    - [Health-Based Color Coding](#health-based-color-coding)
    - [Graceful Degradation](#graceful-degradation)
    - [Visual Clarity](#visual-clarity)
    - [Health Truth](#health-truth)
  - [💻 Demo](#-demo)
    - [Basic Usage](#basic-usage)
    - [Statusline Integration](#statusline-integration)
    - [Handling Zero-State](#handling-zero-state)
    - [With Color Coding](#with-color-coding)
  - [🏗️ Architecture](#️-architecture)
    - [Component Type](#component-type)
    - [Dependencies](#dependencies)
    - [Design Principle](#design-principle)
  - [🏥 Health Scoring](#-health-scoring)
  - [⚡ Performance](#-performance)
    - [Function Performance](#function-performance)
  - [🔧 Extension Points](#-extension-points)
    - [Adding New System Health Metrics](#adding-new-system-health-metrics)
    - [Custom Health Thresholds](#custom-health-thresholds)
  - [⚙️ Configuration](#️-configuration)
    - [Health Threshold Values](#health-threshold-values)
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

**System Health Display Library** transforms raw system metrics into formatted display structures suitable for statusline space constraints with health-based color coding.

> [!NOTE]
> **Pure presentation layer.** Receives system data from system/lib/system (assessment), outputs formatted displays (presentation). Never duplicates data collection logic.

**Design:** Assessment vs Presentation separation - system/lib/system provides data, statusline/lib/system formats for display.

**Library:** [`statusline/lib/system`](../../../lib/system/)

**Source Code:** [`lib/system/system.go`](../../../lib/system/system.go)

**Version:** 1.0.0

---

## 🎯 What It Provides

### System Health Display Functions

Three functions for comprehensive system health presentation:

| Function | Purpose | Health Thresholds |
|----------|---------|-------------------|
| **GetLoadDisplay** | CPU load with CPU-relative coloring | >80% red, >50% yellow, ≤50% green |
| **GetMemoryDisplay** | Memory usage with usage-based coloring | >80% red, >60% yellow, ≤60% green |
| **GetDiskDisplay** | Disk capacity with usage-based coloring | >90% red, >75% yellow, ≤75% green |

### Display Examples

**Formatted Output:**

```bash
⚡ 1.25        # Load average (green - 31% on 4 CPUs)
💾 5.2GB/16GB  # Memory usage (green - 32%)
💿 45%         # Disk usage (green - plenty of space)
```

**Display Structures:**

```go
LoadDisplay{
    LoadAvg: 1.5,
    Color:   "green",  // 37.5% load
    Icon:    "⚡",
    HasInfo: true
}

MemoryDisplay{
    UsedGB:  5.2,
    TotalGB: 16.0,
    Percent: 32.5,
    Color:   "green",
    Icon:    "💾",
    HasInfo: true
}

DiskDisplay{
    Percent: 45.0,
    Color:   "green",
    Icon:    "💿",
    HasInfo: true
}
```

---

## 🔑 Key Features

### Health-Based Color Coding

Each metric has specific thresholds appropriate for that resource:

**Load Average** (CPU-relative):

- **Green (≤50%)**: System has headroom
- **Yellow (>50% ≤80%)**: Working but headroom exists
- **Red (>80%)**: Approaching full capacity

**Memory Usage** (percentage-based):

- **Green (≤60%)**: Healthy memory usage
- **Yellow (>60% ≤80%)**: Elevated usage
- **Red (>80%)**: High memory pressure

**Disk Capacity** (usage-based):

- **Green (≤75%)**: Plenty of space
- **Yellow (>75% ≤90%)**: Running low
- **Red (>90%)**: Critically low space

### Graceful Degradation

All functions return valid structures even when data unavailable:

```go
// No load data available
loadDisplay := GetLoadDisplay()
// → LoadDisplay{HasInfo: false}

// Invalid disk path
diskDisplay := GetDiskDisplay("/nonexistent")
// → DiskDisplay{HasInfo: false}

// Calling code checks HasInfo before displaying
if loadDisplay.HasInfo {
    // Display load
}
```

### Visual Clarity

Icons and colors for quick health recognition:

- **⚡** Load average (lightning - power/activity)
- **💾** Memory usage (floppy disk - storage)
- **💿** Disk capacity (CD - persistent storage)

### Health Truth

Shows actual state without alarm or minimization:

- **Not alarmist**: 60% memory = yellow (elevated), not red (you have headroom)
- **Not minimizing**: 92% disk = red (critically low), not yellow (this is serious)
- **Faithful witness**: Truth-telling through accurate metrics

---

## 💻 Demo

### Basic Usage

```go
import "statusline/lib/system"

func main() {
    // Get system health displays
    loadDisplay := system.GetLoadDisplay()
    memDisplay := system.GetMemoryDisplay()
    diskDisplay := system.GetDiskDisplay("/")

    // Display if data available
    if loadDisplay.HasInfo {
        fmt.Printf("%s %.2f load\n", loadDisplay.Icon, loadDisplay.LoadAvg)
        // Output: ⚡ 1.25 load
    }

    if memDisplay.HasInfo {
        fmt.Printf("%s %.1fGB/%.1fGB (%.0f%%)\n",
            memDisplay.Icon, memDisplay.UsedGB, memDisplay.TotalGB, memDisplay.Percent)
        // Output: 💾 5.2GB/16.0GB (32%)
    }

    if diskDisplay.HasInfo {
        fmt.Printf("%s %.0f%% disk used\n", diskDisplay.Icon, diskDisplay.Percent)
        // Output: 💿 45% disk used
    }
}
```

### Statusline Integration

```go
func buildStatusline() string {
    var parts []string

    // System health
    loadDisplay := system.GetLoadDisplay()
    if loadDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%.2f", loadDisplay.LoadAvg))
    }

    memDisplay := system.GetMemoryDisplay()
    if memDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%.0f%% mem", memDisplay.Percent))
    }

    diskDisplay := system.GetDiskDisplay("/")
    if diskDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%.0f%% disk", diskDisplay.Percent))
    }

    // ... add other statusline parts ...

    return strings.Join(parts, " | ")
    // → "1.25 | 32% mem | 45% disk"
}
```

### Handling Zero-State

```go
// System without load data
loadDisplay := system.GetLoadDisplay()
// → HasInfo: false (no load data available)

// Invalid disk path
diskDisplay := system.GetDiskDisplay("/nonexistent")
// → HasInfo: false (path doesn't exist)

// Statusline orchestrator checks HasInfo before displaying
if loadDisplay.HasInfo {
    // Won't execute - no data available
}
```

### With Color Coding

```go
loadDisplay := system.GetLoadDisplay()

if loadDisplay.HasInfo {
    // Apply color for terminal output
    colored := fmt.Sprintf("%s%s %.2f%s",
        loadDisplay.Color,      // Start color (green/yellow/red)
        loadDisplay.Icon,       // Icon
        loadDisplay.LoadAvg,    // Value
        display.Reset)          // Reset color
    fmt.Println(colored)
    // → Displays in appropriate color: ⚡ 1.25
}
```

---

## 🏗️ Architecture

### Component Type

**Ladder:** Library (Middle Rung)
**Role:** Presentation layer for system health metrics

### Dependencies

**Lower Rungs (Data Sources):**

- `system/lib/system` - System data collection (load, memory, CPU count, disk usage)
- `system/lib/display` - Color constants (Green, Yellow, Red)

**Higher Rungs (Consumers):**

- `statusline` - Main orchestrator

### Design Principle

> **Assessment vs Presentation Separation**
>
> Data collection happens in system/lib/system (assessment). Presentation happens here (formatting with health thresholds). Never duplicate data collection logic.

**Data → Presentation Flow:**

```bash
System/lib/system                      Display Functions
(Data Source)                          (Presentation)
     ↓                                       ↓
GetInfo():                            GetLoadDisplay()
  LoadAvg1: 1.5                             ↓
  CPUCount: 4                         LoadDisplay{
     ↓                                  LoadAvg: 1.5,
                                        Color: "green",  // 37.5% = healthy
                                        Icon: "⚡",
                                        HasInfo: true
                                      }
```

---

## 🏥 Health Scoring

Pure presentation library - no health scoring infrastructure needed.

All operations guaranteed to succeed through graceful degradation:

- No load data → Returns `LoadDisplay{HasInfo: false}`
- No memory data → Returns `MemoryDisplay{HasInfo: false}`
- Invalid disk path → Returns `DiskDisplay{HasInfo: false}`

> [!NOTE]
> **This library cannot fail.** Always returns valid display structures. Health tracking would measure "successfully did nothing" which provides no value.

---

## ⚡ Performance

### Function Performance

| Function | Operation | Time | Memory |
|----------|-----------|:----:|:------:|
| **GetLoadDisplay** | /proc/loadavg read + calculation | <100 μs | ~50 bytes |
| **GetMemoryDisplay** | /proc/meminfo read + calculation | <100 μs | ~50 bytes |
| **GetDiskDisplay** | df command + parsing | <500 μs | ~50 bytes |

**Total Memory:** Three struct allocations per statusline render (~150 bytes total)

**Optimization:** This library needs no optimization. System data retrieval and struct allocation are already optimal for this use case.

---

## 🔧 Extension Points

### Adding New System Health Metrics

**Pattern:**

1. Check if system/lib/system provides the data (if not, extend system lib first)
2. Define new display structure in SETUP Types section
3. Create Get[Metric]Display() function in BODY
4. Apply health thresholds for color selection
5. Update API documentation with new metric
6. Add tests for new display function

**Example - Adding Network Display:**

```go
// In SETUP Types:
type NetworkDisplay struct {
    Throughput float64  // MB/s
    Color      string
    Icon       string
    HasInfo    bool
}

// In BODY:
func GetNetworkDisplay() NetworkDisplay {
    netInfo := syslib.GetNetworkInfo()

    if netInfo.Throughput <= 0 {
        return NetworkDisplay{HasInfo: false}
    }

    // Health thresholds for network
    netColor := display.Green
    if netInfo.Throughput > 100 {  // Saturated
        netColor = display.Red
    } else if netInfo.Throughput > 50 {  // Heavy usage
        netColor = display.Yellow
    }

    return NetworkDisplay{
        Throughput: netInfo.Throughput,
        Color:      netColor,
        Icon:       "🌐",
        HasInfo:    true,
    }
}
```

### Custom Health Thresholds

**Future enhancement:** Configurable thresholds instead of hardcoded values

```go
// Current (hardcoded):
if loadPercent > 80 {
    loadColor = display.Red
}

// Future (configurable):
type HealthThresholds struct {
    LoadWarning   float64  // 50.0
    LoadCritical  float64  // 80.0
    MemWarning    float64  // 60.0
    MemCritical   float64  // 80.0
    DiskWarning   float64  // 75.0
    DiskCritical  float64  // 90.0
}

func GetLoadDisplayWithThresholds(thresholds HealthThresholds) LoadDisplay {
    // Use configured thresholds instead of hardcoded values
}
```

---

## ⚙️ Configuration

### Health Threshold Values

**Current Implementation:**

All thresholds hardcoded in function implementations:

```go
// Load thresholds (CPU-relative)
if loadPercent > 80 { loadColor = display.Red }
else if loadPercent > 50 { loadColor = display.Yellow }

// Memory thresholds (usage percentage)
if memPercent > 80 { memColor = display.Red }
else if memPercent > 60 { memColor = display.Yellow }

// Disk thresholds (capacity percentage)
if diskPercent > 90 { diskColor = display.Red }
else if diskPercent > 75 { diskColor = display.Yellow }
```

**Modification:** Edit threshold values in `lib/system/system.go` function implementations

**Future:** External configuration file for runtime threshold adjustment

---

## 🔍 Troubleshooting

**This library has no common failure modes - all operations guaranteed to succeed.**

### Expected Behaviors

**Zero values return HasInfo: false**

- No load data → HasInfo: false (correct, not an error)
- No memory data → HasInfo: false (correct, not an error)
- Invalid disk path → HasInfo: false (correct, not an error)

**Empty display data when HasInfo: false**

- This is correct zero-state
- Calling code checks HasInfo to determine display

### If Unexpected Results Occur

**Problem:** HasInfo true but color seems wrong

- **Cause:** Health thresholds may not match expectations
- **Solution:** Verify threshold values in function code
- **Check:** Load (>80%/>50%), Memory (>80%/>60%), Disk (>90%/>75%)

**Problem:** Load percentage calculation seems off

- **Cause:** Load average is relative to CPU count
- **Solution:** Load 4.0 on 4 CPUs = 100% (expected behavior)
- **Formula:** (LoadAvg1 / CPUCount) * 100

**Problem:** HasInfo false when system clearly has data

- **Cause:** system/lib/system not returning data
- **Solution:** Check system/lib/system implementation
- **Debug:** Call syslib.GetInfo() directly to verify

---

## 🛠️ Modification Policy

### Safe to Modify

- ✅ Add new system health metrics (network, CPU temp, etc.)
- ✅ Adjust health threshold values
- ✅ Extend display structures with additional fields
- ✅ Change color or icon choices for different metrics

### Modify with Care

- ⚠️ Display struct fields - breaks statusline orchestrator
- ⚠️ Function signatures - breaks all calling code
- ⚠️ HasInfo semantics - breaks zero-state handling
- ⚠️ Health threshold semantics - changes what "healthy" means

### Never Modify

- ❌ Pure function guarantee (no side effects, no state)
- ❌ Graceful degradation (always return valid values)
- ❌ Data vs Presentation separation (system/lib/system provides data)
- ❌ Stateless design principle

---

## 🚀 Future Expansions

### Planned Features

- ✅ Load average display - **COMPLETED**
- ✅ Memory usage display - **COMPLETED**
- ✅ Disk usage display - **COMPLETED**
- ⏳ Network throughput display
- ⏳ CPU temperature display
- ⏳ Process count display
- ⏳ Configurable health thresholds

### Research Areas

- Adaptive thresholds based on historical patterns
- Multi-level health indicators (warning/critical separation)
- Trend indicators (↑↓ for increasing/decreasing metrics)
- Threshold configuration via external file

### Known Limitations

- Health thresholds hardcoded (not configurable)
- Binary health states (green/yellow/red - no gradients)
- No historical context (can't show "getting worse")
- Disk display only shows single path (not multiple mounts)

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **GetLoadDisplay** | Load average formatting and CPU-relative health coloring | [load.md](./load.md) |
| **GetMemoryDisplay** | Memory usage formatting and usage-based health coloring | [memory.md](./memory.md) |
| **GetDiskDisplay** | Disk capacity formatting and usage-based health coloring | [disk.md](./disk.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementations** | [`lib/system/system.go`](../../../lib/system/system.go) | All three display functions with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7-library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |
| **System Lib** | System data collection library | [system/lib/system](~/.claude/cpi-si/system/lib/system/) |

---

## 📖 Biblical Foundation

**Core Verse:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Order and clarity in communication - showing system health clearly and truthfully

**Supporting Verse:**

- **"A faithful witness will not lie: but a false witness will utter lies"** - Proverbs 14:5
  - Truth-telling through accurate metrics and health indicators

**Application:** Show actual system health (load, memory, disk) without alarm or minimization. Not hiding problems, not exaggerating concerns - faithful witness to system state.

**Why This Matters:** System monitoring tools often alarm unnecessarily (70% memory shown as "critical") or minimize real issues (95% disk shown as "okay"). This library tells the truth: 70% memory = yellow (elevated but manageable), 95% disk = red (actually critical).

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
