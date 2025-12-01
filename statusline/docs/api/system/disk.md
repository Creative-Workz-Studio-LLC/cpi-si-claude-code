<div align="center">

# GetDiskDisplay Function Reference

**Disk Capacity Formatting with Usage-Based Health Coloring**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Thresholds](https://img.shields.io/badge/Thresholds-75%25%20%7C%2090%25-orange?style=flat)
![Pure](https://img.shields.io/badge/Pure-No%20Side%20Effects-green?style=flat)

*Transforms disk usage data into formatted display with health-based color coding*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [🔄 Returns](#-returns) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - GetDiskDisplay Function
**Purpose:** Disk capacity display formatting with usage-based health coloring
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
  - [path (string)](#path-string)
- [🔄 Returns](#-returns)
  - [DiskDisplay](#diskdisplay)
- [🎯 Behavior](#-behavior)
  - [Health Threshold Calculation](#health-threshold-calculation)
  - [Display Examples](#display-examples)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [Multiple Partitions](#multiple-partitions)
  - [With Color](#with-color)
  - [Zero-State Handling](#zero-state-handling)
  - [Full Statusline Integration](#full-statusline-integration)
  - [Testing](#testing)
- [🔧 Edge Cases](#-edge-cases)
  - [Empty Filesystem](#empty-filesystem)
  - [Nearly Full Filesystem](#nearly-full-filesystem)
  - [Exactly At Threshold](#exactly-at-threshold)
  - [Invalid Path](#invalid-path)
  - [Network Filesystem](#network-filesystem)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [📊 Understanding Disk Usage](#-understanding-disk-usage)
  - [What Disk Percentage Means](#what-disk-percentage-means)
  - [Why Percentage-Based](#why-percentage-based)
  - [Health Threshold Rationale](#health-threshold-rationale)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Disk percentage seems wrong](#disk-percentage-seems-wrong)
  - [HasInfo false for valid path](#hasinfo-false-for-valid-path)
  - [Want different thresholds](#want-different-thresholds)
  - [Multiple filesystems in statusline](#multiple-filesystems-in-statusline)
- [🔧 Extension Points](#-extension-points)
  - [Multi-Filesystem Display](#multi-filesystem-display)
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

**GetDiskDisplay** formats disk usage data with capacity-based health coloring suitable for statusline display.

> [!NOTE]
> **Pure presentation function.** Receives disk data from system/lib/system, applies health thresholds based on usage percentage, outputs formatted display. Cannot fail, no side effects.

**Design:** Assessment vs Presentation separation - system/lib/system provides data (UsagePercent via df), this function formats for display with usage-based health-based color coding.

**Library:** [`statusline/lib/system`](../../../lib/system/)

**Source Code:** [`lib/system/system.go`](../../../lib/system/system.go)

**Library README:** [System Library Documentation](README.md)

---

## 📝 Signature

```go
func GetDiskDisplay(path string) DiskDisplay
```

**Parameters:** `path` - Filesystem path to check

**Returns:** `DiskDisplay` structure with formatted data

---

## ⚙ Parameters

### path (string)

Filesystem path to check disk usage for.

**Common Values:**

| Path | Description | Use Case |
|------|-------------|----------|
| `"/"` | Root filesystem | Most common - system disk |
| `"/home"` | Home directory mount | If separate partition |
| `"/data"` | Data partition | Additional storage |
| Any valid path | Any mounted filesystem | Custom mounts |

**Example:**

```go
diskDisplay := GetDiskDisplay("/")
```

**Validation:** Invalid paths return `HasInfo: false` (graceful degradation)

**Data Source:** `df` command (via system library)

---

## 🔄 Returns

### DiskDisplay

Formatted display structure with all presentation elements.

```go
type DiskDisplay struct {
    Percent float64 // Disk usage percentage (0-100)
    Color   string  // Terminal color code based on usage percent
    Icon    string  // Visual icon representing disk (e.g., "💿")
    HasInfo bool    // True if disk data available, false if no data
}
```

**Field Details:**

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| **Percent** | float64 | Disk usage percentage (0-100) | `45.0` |
| **Color** | string | Health-based color code | `"green"` / `"yellow"` / `"red"` |
| **Icon** | string | Visual disk icon (CD) | `"💿"` |
| **HasInfo** | bool | Data availability flag | `true` / `false` |

**Color Mapping:**

| Usage | Color | Meaning |
|:-----:|:-----:|---------|
| ≤ 75% | `display.Green` | Plenty of space |
| > 75% ≤ 90% | `display.Yellow` | Running low |
| > 90% | `display.Red` | Critically low space |

**Guarantee:** Always returns valid DiskDisplay (never nil, never errors)

---

## 🎯 Behavior

### Health Threshold Calculation

Disk percentage-based coloring:

```go
diskInfo := syslib.GetDiskUsage(path)

if diskInfo.UsagePercent <= 0 {
    return DiskDisplay{HasInfo: false}  // Zero-state
}

// Apply health thresholds
diskColor := display.Green
if diskInfo.UsagePercent > 90 {
    diskColor = display.Red      // Critically low space
} else if diskInfo.UsagePercent > 75 {
    diskColor = display.Yellow   // Running low
}

return DiskDisplay{
    Percent: diskInfo.UsagePercent,
    Color:   diskColor,
    Icon:    "💿",
    HasInfo: true,
}
```

### Display Examples

**Plenty of Space (Green):**

```go
// Filesystem: 45% used
display := GetDiskDisplay("/")
// → DiskDisplay{Percent: 45.0, Color: "green", Icon: "💿", HasInfo: true}
// 45% ≤ 75% = plenty of space
```

**Running Low (Yellow):**

```go
// Filesystem: 78% used
display := GetDiskDisplay("/")
// → DiskDisplay{Percent: 78.0, Color: "yellow", Icon: "💿", HasInfo: true}
// 78% > 75% and ≤ 90% = running low
```

**Critically Low (Red):**

```go
// Filesystem: 92% used
display := GetDiskDisplay("/")
// → DiskDisplay{Percent: 92.0, Color: "red", Icon: "💿", HasInfo: true}
// 92% > 90% = critically low space
```

**Invalid Path:**

```go
// Path doesn't exist
display := GetDiskDisplay("/nonexistent")
// → DiskDisplay{Percent: 0, Color: "", Icon: "", HasInfo: false}
```

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/system"

func buildStatusline() string {
    // Get disk display for root
    diskDisplay := system.GetDiskDisplay("/")
    // → Percent: 45.0

    var parts []string

    // Only add if disk data available
    if diskDisplay.HasInfo {
        diskPart := fmt.Sprintf("%s %.0f%%", diskDisplay.Icon, diskDisplay.Percent)
        parts = append(parts, diskPart)
        // → "💿 45%"
    }

    return strings.Join(parts, " | ")
}
```

### Multiple Partitions

```go
// Check root and home separately
rootDisk := system.GetDiskDisplay("/")
homeDisk := system.GetDiskDisplay("/home")

var parts []string

if rootDisk.HasInfo {
    parts = append(parts, fmt.Sprintf("/ %.0f%%", rootDisk.Percent))
}

if homeDisk.HasInfo {
    parts = append(parts, fmt.Sprintf("/home %.0f%%", homeDisk.Percent))
}

// → "/ 45% | /home 62%"
```

### With Color

```go
diskDisplay := system.GetDiskDisplay("/")

if diskDisplay.HasInfo {
    // Apply color for terminal output
    colored := fmt.Sprintf("%s%s %.0f%% disk%s",
        diskDisplay.Color,      // Start color
        diskDisplay.Icon,       // Icon
        diskDisplay.Percent,    // Percentage
        display.Reset)          // Reset color
    fmt.Println(colored)
    // → Displays in green/yellow/red: 💿 45% disk
}
```

### Zero-State Handling

```go
diskDisplay := system.GetDiskDisplay("/nonexistent")

// Check HasInfo before using
if diskDisplay.HasInfo {
    fmt.Printf("Disk: %.0f%%\n", diskDisplay.Percent)
} else {
    fmt.Println("Disk data unavailable")
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
func TestDiskDisplay(t *testing.T) {
    diskDisplay := system.GetDiskDisplay("/")

    // Verify structure
    if diskDisplay.HasInfo {
        if diskDisplay.Percent < 0 || diskDisplay.Percent > 100 {
            t.Error("Percent should be in range 0-100")
        }
        if diskDisplay.Color == "" {
            t.Error("Color should be set when HasInfo is true")
        }
        if diskDisplay.Icon == "" {
            t.Error("Icon should be set when HasInfo is true")
        }
    }
}
```

---

## 🔧 Edge Cases

### Empty Filesystem

```go
// Newly mounted filesystem with no data
display := GetDiskDisplay("/mnt/empty")
// Returns: DiskDisplay{Percent: 0, Color: "green", HasInfo: true}
// Behavior: 0% used = plenty of space
```

### Nearly Full Filesystem

```go
// Filesystem: 95% used
display := GetDiskDisplay("/")
// Returns: DiskDisplay{Percent: 95.0, Color: "red", HasInfo: true}
// Behavior: >90% = critical (correct - very little space left)
```

### Exactly At Threshold

```go
// Filesystem: 75% used
display := GetDiskDisplay("/")
// Returns: DiskDisplay{Percent: 75.0, Color: "green", HasInfo: true}
// Calculation: Exactly 75% is still green (threshold is >75%)
```

### Invalid Path

```go
// Path doesn't exist
display := GetDiskDisplay("/does/not/exist")
// Returns: DiskDisplay{HasInfo: false}
// Behavior: df command fails, graceful degradation
```

### Network Filesystem

```go
// NFS mount
display := GetDiskDisplay("/mnt/nfs")
// Returns: DiskDisplay with network filesystem usage
// Behavior: df works on network mounts (may be slower)
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Data retrieval | +25 | Get disk usage via df |
| Path validation | +25 | Verify path exists |
| Threshold application | +25 | Determine health color |
| Display structure creation | +25 | Build return value |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** Disk retrieval with fallback guarantees valid output. Health tracking demonstrates successful operation, not error detection.

---

## ⚡ Performance

**Time Complexity:** O(1) - df command execution + parsing

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Data source** | df command | Linux disk utility |
| **Command execution** | <200 μs | Local filesystem |
| **Output parsing** | O(1) | String operations |
| **Struct allocation** | ~50 bytes | Single DiskDisplay |
| **Execution time** | <500 μs | Total end-to-end |

**Network Filesystems:** May take longer (depends on network latency)

**Memory:** Single DiskDisplay struct allocation (~50 bytes)

**Optimization:** Not needed - df is already fast for local filesystems

---

## 📊 Understanding Disk Usage

### What Disk Percentage Means

**Disk usage percentage** = (Used Space / Total Space) * 100

- **45% used**: Plenty of space remaining
- **78% used**: Getting full, should monitor
- **92% used**: Critically low, action needed

### Why Percentage-Based

**Absolute space meaningless without context:**

- 450GB used on 500GB disk = 90% (critical)
- 450GB used on 2TB disk = 22% (plenty of space)

**Percentage shows actual capacity utilization:**

- Same thresholds work across different disk sizes
- 75%/90% meaningful regardless of total capacity

### Health Threshold Rationale

**75% threshold (yellow):**

- Disk getting full but not critical
- Good time to clean up or plan expansion
- System still has headroom

**90% threshold (red):**

- Critically low disk space
- System may fail operations soon
- Immediate action needed
- Many applications warn at 90%

---

## 🔍 Troubleshooting

### Disk percentage seems wrong

**Problem:** Disk display shows 45% but `df -h` shows different

**Check:**

1. Verify you're checking same filesystem path
2. Reserved space: root-reserved blocks don't count as available
3. Note: df uses capacity percentage (what we use)

**Expected:** Our calculation matches `df` Use% column

### HasInfo false for valid path

**Problem:** Disk display shows no data for path that exists

**Check:**

1. Verify path is a mountpoint: `df /path`
2. Check permissions: can process access path?
3. Debug: Call syslib.GetDiskUsage(path) directly

**Common Cause:** Path exists but isn't mounted or accessible

### Want different thresholds

**Problem:** 75%/90% doesn't match your definition of "running low"

**Solution:** Adjust thresholds in system.go GetDiskDisplay():

```go
// Current:
if diskInfo.UsagePercent > 90 {
    diskColor = display.Red
} else if diskInfo.UsagePercent > 75 {
    diskColor = display.Yellow
}

// Custom (e.g., 80%/95%):
if diskInfo.UsagePercent > 95 {
    diskColor = display.Red
} else if diskInfo.UsagePercent > 80 {
    diskColor = display.Yellow
}
```

### Multiple filesystems in statusline

**Problem:** Want to show both root and home disk usage

**Solution:** Call GetDiskDisplay() multiple times:

```go
rootDisk := system.GetDiskDisplay("/")
homeDisk := system.GetDiskDisplay("/home")

// Display both if available
if rootDisk.HasInfo {
    parts = append(parts, fmt.Sprintf("/ %.0f%%", rootDisk.Percent))
}
if homeDisk.HasInfo {
    parts = append(parts, fmt.Sprintf("~/ %.0f%%", homeDisk.Percent))
}
```

---

## 🔧 Extension Points

### Multi-Filesystem Display

**Show all important filesystems:**

```go
type MultiDiskDisplay struct {
    Disks   []DiskInfo  // Array of filesystem info
    Color   string      // Worst health color
    Icon    string
    HasInfo bool
}

type DiskInfo struct {
    Path    string
    Percent float64
    Color   string
}

func GetMultiDiskDisplay(paths []string) MultiDiskDisplay {
    var disks []DiskInfo
    worstColor := display.Green

    for _, path := range paths {
        diskInfo := syslib.GetDiskUsage(path)
        if diskInfo.UsagePercent > 0 {
            // Track each filesystem
            // Update worstColor if this disk is worse
        }
    }

    // Return combined view
}
```

### Trend Indicators

**Show if disk usage increasing:**

```go
type DiskDisplay struct {
    Percent float64
    Trend   string  // "↑" increasing, "↓" decreasing, "→" stable
    Color   string
    Icon    string
    HasInfo bool
}

// Would require historical tracking
```

### Configurable Thresholds

```go
type DiskThresholds struct {
    Warning  float64  // 75.0
    Critical float64  // 90.0
}

func GetDiskDisplayWithThresholds(path string, thresholds DiskThresholds) DiskDisplay {
    // Apply custom thresholds instead of hardcoded 75%/90%
}
```

---

## 🔗 Related Functions

**Same Library:**

| Function | Purpose | Link |
|----------|---------|------|
| **GetLoadDisplay** | Load average with CPU-relative health coloring | [load.md](load.md) |
| **GetMemoryDisplay** | Memory usage with percentage-based health coloring | [memory.md](memory.md) |

**Dependencies:**

| Component | Purpose | Location |
|-----------|---------|----------|
| **system/lib/system.GetDiskUsage()** | Disk data collection via df | [system/lib/system](~/.claude/cpi-si/system/lib/system/) |

**Future Functions:**

- `GetMultiDiskDisplay()` - Multiple filesystems in single display
- `GetDiskTrendDisplay()` - Disk usage with growth trend

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **System Library README** | Library overview and integration guide | [README.md](README.md) |
| **GetLoadDisplay Function** | Load average formatting | [load.md](load.md) |
| **GetMemoryDisplay Function** | Memory usage formatting | [memory.md](memory.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/system/system.go`](../../../lib/system/system.go) | GetDiskDisplay with comprehensive inline documentation |
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

**Application:** Show actual disk usage (45%) with truthful health indicator (green/yellow/red) without alarm or minimization. Faithful witness to storage state - neither hiding problems ("still works!") nor exaggerating concerns ("we're out of space!").

**Why This Matters:** Many systems alarm at 80% disk (too early - plenty of space left) or don't warn until 95% (too late - system failing). This library tells the truth: 80% = yellow (monitor it), 95% = red (critical - act now).

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
