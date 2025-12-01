<div align="center">

# 📚 Statusline API Documentation

**Complete API Reference for 7 Presentation Libraries**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Libraries](https://img.shields.io/badge/Libraries-7-orange?style=flat)
![Functions](https://img.shields.io/badge/Functions-14+-blue?style=flat)

*Transform Claude Code session data into concise, readable statusline display*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📚 Libraries](#-the-7-libraries) • [🚀 Quick Start](#-quick-start) • [🎯 Design](#-design-principles) • [📚 References](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation Index
**Purpose:** Complete API reference for all 7 statusline presentation libraries
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📖 Overview](#-overview)
  - [🎯 What This Documentation Provides](#-what-this-documentation-provides)
  - [💡 Design Philosophy](#-design-philosophy)
- [📚 The 7 Libraries](#-the-7-libraries)
  - [📋 1. Types - Session Data Contract](#-1-types---session-data-contract)
  - [⚙ 2. Features - Display Timing](#-2-features---display-timing)
  - [🎨 3. Format - Text Optimization](#-3-format---text-optimization)
  - [🌿 4. Git - Repository Display](#-4-git---repository-display)
  - [📊 5. Session - Statistics Display](#-5-session---statistics-display)
  - [🖥 6. System - Health Display](#-6-system---health-display)
  - [⏰ 7. Temporal - Awareness Display](#-7-temporal---awareness-display)
- [🚀 Quick Start](#-quick-start)
  - [Types Library](#types-library)
  - [Features Library](#features-library)
  - [Format Library](#format-library)
  - [Session Library](#session-library)
  - [System Library](#system-library)
  - [Temporal Library](#temporal-library)
- [🎯 Design Principles](#-design-principles)
  - [1. Presentation Logic Only](#1-presentation-logic-only)
  - [2. Space Optimization](#2-space-optimization)
  - [3. Graceful Degradation](#3-graceful-degradation)
- [📖 Biblical Foundation](#-biblical-foundation)
- [⚡ Performance](#-performance)
- [🔧 Extension Guide](#-extension-guide)
  - [Adding New Timing Strategies](#adding-new-timing-strategies)
  - [Adding New Format Functions](#adding-new-format-functions)
  - [Adding New Display Functions](#adding-new-display-functions)
- [📚 References & Resources](#-references--resources)
  - [🏗 Library Documentation](#-library-documentation)
  - [📐 Source Code](#-source-code)
  - [🔗 Related Documentation](#-related-documentation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

The Statusline API provides **7 specialized libraries** for formatting and timing decisions in statusline display. These libraries handle presentation logic - transforming verbose system data into concise, readable forms optimized for limited horizontal space.

> [!IMPORTANT]
> **This is the API reference index.** For architectural overview and library integration patterns, see [lib/README.md](../../lib/README.md). For project overview and usage, see [main README.md](../../README.md).

### 🎯 What This Documentation Provides

| Section | Purpose | Audience |
|---------|---------|----------|
| **Library Catalog** | Overview of all 7 libraries | Developers choosing which library to use |
| **Quick Start** | Working code examples | Developers integrating statusline libraries |
| **Design Principles** | Core patterns and constraints | Library authors and maintainers |
| **API References** | Detailed function signatures | Developers using specific functions |

### 💡 Design Philosophy

**Clarity through brevity** - show what matters without ceremony.

| Icon | Principle | Implementation |
|:----:|-----------|----------------|
| 🎯 | **One Concern** | Each library handles exactly ONE aspect of display |
| 🔄 | **Pure Functions** | Same input → same output, no side effects |
| 🛡 | **Never Fail** | Graceful degradation with safe defaults |
| ⚡ | **Performance** | <1μs per function call (typical) |

---

## 📚 The 7 Libraries

> [!TIP]
> **Complete architectural overview:** See [lib/README.md](../../lib/README.md) for ladder/baton patterns, dependency flow, and integration guidance.

### 📋 1. Types - Session Data Contract

**Purpose:** Define SessionContext structure for Claude Code JSON input

**Path:** [`lib/types/`](../../lib/types/)

**Key Type:**

- [`SessionContext`](types/context.md) - Complete session data structure

**Use Cases:**

- Parse Claude Code JSON input with type safety
- Access session metadata (ID, model, workspace, cost)
- Provide structured data to all statusline components
- Establish compile-time verified contract

**Documentation:** [Types Library README](types/README.md)

**Source Code:** [`lib/types/context.go`](../../lib/types/context.go)

---

### ⚙ 2. Features - Display Timing

**Purpose:** Determine WHEN to show statusline elements

**Path:** [`lib/features/`](../../lib/features/)

**Key Functions:**

- [`ShouldShowReminder(sessionID)`](features/reminder.md) - Kingdom Technology reminder timing (20% frequency)

**Use Cases:**

- Session-based display decisions
- Deterministic frequency control
- Reminder timing logic

**Documentation:** [Features Library README](features/README.md)

**Source Code:** [`lib/features/reminder.go`](../../lib/features/reminder.go)

---

### 🎨 3. Format - Text Optimization

**Purpose:** Transform verbose data into concise display forms

**Path:** [`lib/format/`](../../lib/format/)

**Key Functions:**

- [`GetShortModelName(displayName)`](format/model.md) - Model tier extraction ("Claude 3.5 Sonnet" → "Sonnet")
- [`ShortenPath(path)`](format/path.md) - Path optimization ("/home/user/project" → "~/project")

**Use Cases:**

- Model name shortening for limited space
- Path display optimization
- Space-efficient text presentation

**Documentation:** [Format Library README](format/README.md)

**Source Code:** [`lib/format/model.go`](../../lib/format/model.go), [`lib/format/path.go`](../../lib/format/path.go)

---

### 🌿 4. Git - Repository Display

**Purpose:** Format git repository status for statusline

**Path:** [`lib/git/`](../../lib/git/)

**Key Functions:**

- [`GetGitDisplay(workdir)`](git/git.md) - Git status formatting with branch, dirty status, ahead/behind

**Use Cases:**

- Branch name display with visual indicators
- Uncommitted changes indicator (*)
- Ahead/behind tracking (↑↓)
- Color coding and icon selection

**Documentation:** [Git Library README](git/README.md)

**Source Code:** [`lib/git/git.go`](../../lib/git/git.go)

---

### 📊 5. Session - Statistics Display

**Purpose:** Format session statistics (lines, duration, cost)

**Path:** [`lib/session/`](../../lib/session/)

**Key Functions:**

- [`GetLinesModifiedDisplay(ctx)`](session/lines.md) - Lines changed (added + removed)
- [`GetDurationDisplay(ctx)`](session/duration.md) - Session duration (45s, 1h 23m)
- [`GetCostDisplay(ctx)`](session/cost.md) - Session cost ($0.0123)
- [`GetFormattedCost(cost)`](session/cost.md) - Cost string helper

**Use Cases:**

- Lines modified tracking (total churn)
- Session duration display (human-readable)
- Cost tracking (token usage, API costs)
- Zero-state handling for missing data

**Documentation:** [Session Library README](session/README.md)

**Source Code:** [`lib/session/session.go`](../../lib/session/session.go)

---

### 🖥 6. System - Health Display

**Purpose:** Format system health metrics with color coding

**Path:** [`lib/system/`](../../lib/system/)

**Key Functions:**

- [`GetLoadDisplay()`](system/load.md) - CPU load with health coloring (green/yellow/red)
- [`GetMemoryDisplay()`](system/memory.md) - Memory usage with health coloring
- [`GetDiskDisplay(path)`](system/disk.md) - Disk usage with health coloring

**Health Thresholds:**

| Metric | 🟢 Green | 🟡 Yellow | 🔴 Red |
|--------|---------|----------|--------|
| **Load** | < 50% CPU count | 50-80% | > 80% |
| **Memory** | < 60% used | 60-80% | > 80% |
| **Disk** | < 75% full | 75-90% | > 90% |

**Documentation:** [System Library README](system/README.md)

**Source Code:** [`lib/system/system.go`](../../lib/system/system.go)

---

### ⏰ 7. Temporal - Awareness Display

**Purpose:** Format temporal awareness (time, schedule, calendar)

**Path:** [`lib/temporal/`](../../lib/temporal/)

**Key Functions:**

- [`GetTimeOfDayDisplay(ctx)`](temporal/timeofday.md) - Time of day with icon (🌅 morning, ☀ afternoon, 🌆 evening, 🌙 night)
- [`GetSessionPhaseDisplay(ctx)`](temporal/sessionphase.md) - Session duration and phase (fresh/working/long)
- [`GetScheduleDisplay(ctx)`](temporal/schedule.md) - Current schedule state (work/sleep/personal)
- [`GetCalendarDisplay(ctx)`](temporal/calendar.md) - Date and week number

**Circadian Awareness:** Integrates with CPI-SI [`hooks/lib/temporal`](~/.claude/hooks/lib/temporal/) for schedule and rhythm awareness.

**Documentation:** [Temporal Library README](temporal/README.md)

**Source Code:** [`lib/temporal/temporal.go`](../../lib/temporal/temporal.go)

---

## 🚀 Quick Start

### Types Library

```go
import (
    "encoding/json"
    "statusline/lib/types"
)

// Parse Claude Code JSON input
var ctx types.SessionContext
err := json.Unmarshal(claudeCodeJSON, &ctx)

// Access session data
sessionID := ctx.SessionID
modelName := ctx.Model.DisplayName
costUSD := ctx.Cost.TotalCostUSD
```

### Features Library

```go
import "statusline/lib/features"

// Determine if Kingdom Technology reminder should display
if features.ShouldShowReminder(ctx.SessionID) {
    parts = append(parts, "⛪ Kingdom Technology")
}
```

### Format Library

```go
import "statusline/lib/format"

// Shorten model name for display
modelName := format.GetShortModelName(ctx.Model.DisplayName)
// "Claude 3.5 Sonnet" → "Sonnet"

// Shorten working directory path
workdir := format.ShortenPath(ctx.Workspace.CurrentDir)
// "/home/user/project" → "~/project"
```

### Session Library

```go
import "statusline/lib/session"

// Get session statistics displays
linesDisplay := session.GetLinesModifiedDisplay(ctx)
durationDisplay := session.GetDurationDisplay(ctx)
costDisplay := session.GetCostDisplay(ctx)

// Check if data available before displaying
if linesDisplay.HasInfo {
    fmt.Printf("%s%s %d lines%s",
        linesDisplay.Color, linesDisplay.Icon,
        linesDisplay.TotalLines, display.Reset)
}
```

### System Library

```go
import "statusline/lib/system"

// Get system health displays
loadDisplay := system.GetLoadDisplay()
memDisplay := system.GetMemoryDisplay()
diskDisplay := system.GetDiskDisplay("/")

// Check if data available before displaying
if loadDisplay.HasInfo {
    fmt.Printf("%s%s %.2f%s",
        loadDisplay.Color, loadDisplay.Icon,
        loadDisplay.LoadAvg, display.Reset)
}
```

### Temporal Library

```go
import "statusline/lib/temporal"
import hookslib "hooks/lib/temporal"

// Get temporal context
tempCtx, _ := hookslib.GetTemporalContext()

// Get temporal awareness displays
todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
schedDisplay := temporal.GetScheduleDisplay(tempCtx)

// Display with color coding
fmt.Printf("%s%s %s%s",
    todDisplay.Color, todDisplay.Icon,
    todDisplay.Label, display.Reset)
```

---

## 🎯 Design Principles

### 1. Presentation Logic Only

These libraries handle **HOW** to display information, not **WHAT** to display:

| Component | Responsibility |
|-----------|----------------|
| **Features library** | WHEN to show elements (timing decisions) |
| **Presentation libraries** | HOW to format data (formatting transformations) |
| **Statusline orchestrator** | WHAT to show and WHERE (assembly and layout) |

### 2. Space Optimization

Statusline has limited horizontal space (typically 80-120 characters total):

| Challenge | Solution | Example |
|-----------|----------|---------|
| Long model names | Show tier only | "Claude 3.5 Sonnet" → "Sonnet" |
| Deep paths | Use ~ for home, basename for long paths | "/home/user/project" → "~/project" |
| Consistency | Deterministic - same input → same output | Always produces identical results |

### 3. Graceful Degradation

All functions designed to **never fail**:

| Scenario | Response | Benefit |
|----------|----------|---------|
| Unknown inputs | Safe defaults | System keeps running |
| Missing data | `HasInfo: false` in display struct | Optional element skipped |
| System errors | Zero values, empty strings | Partial display shown |
| Invalid data | Sanitized output | No crashes |

> [!IMPORTANT]
> **No panics, no errors returned** - statusline rendering must always succeed. Better to show partial info than crash the interaction.

---

## 📖 Biblical Foundation

**Core Verses:**

- **"To every thing there is a season, and a time to every purpose under the heaven"** - Ecclesiastes 3:1
  - Wisdom in timing (knowing when to speak and when to be silent)
- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Order and clarity in communication
- **"A word fitly spoken is like apples of gold in pictures of silver"** - Proverbs 25:11
  - Concise, beautiful communication
- **"Make the path straight"** - Matthew 3:3
  - Clarity in direction and presentation

---

## ⚡ Performance

All library functions are lightweight:

| Library | Performance | Notes |
|---------|-------------|-------|
| [**types**](../../lib/types/) | Zero runtime cost | Type definitions only |
| [**features**](../../lib/features/) | < 1μs | Hash computation |
| [**format**](../../lib/format/) | < 1μs | String operations |
| [**git**](../../lib/git/) | < 1μs | String formatting |
| [**session**](../../lib/session/) | < 1μs | Arithmetic and formatting |
| [**system**](../../lib/system/) | First call: ~1ms, cached: < 1μs | System calls cached |
| [**temporal**](../../lib/temporal/) | < 1μs | Field extraction |

**Total overhead:** <10ms for complete statusline assembly (typical)

---

## 🔧 Extension Guide

### Adding New Timing Strategies

Follow `Should[Strategy]Reminder` pattern in features library.

See [features/README.md](features/README.md) for extension points.

### Adding New Format Functions

Follow stateless pure function pattern in format library.

See [format/README.md](format/README.md) for extension points.

### Adding New Display Functions

Each library has extension points documented in its BODY section.

See individual library README files for specific guidance.

---

## 📚 References & Resources

### 🏗 Library Documentation

| Library | API Reference | Source Code |
|---------|---------------|-------------|
| **Types** | [types/README.md](types/README.md) | [lib/types/](../../lib/types/) |
| **Features** | [features/README.md](features/README.md) | [lib/features/](../../lib/features/) |
| **Format** | [format/README.md](format/README.md) | [lib/format/](../../lib/format/) |
| **Git** | [git/README.md](git/README.md) | [lib/git/](../../lib/git/) |
| **Session** | [session/README.md](session/README.md) | [lib/session/](../../lib/session/) |
| **System** | [system/README.md](system/README.md) | [lib/system/](../../lib/system/) |
| **Temporal** | [temporal/README.md](temporal/README.md) | [lib/temporal/](../../lib/temporal/) |

### 📐 Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Library Implementation** | [lib/](../../lib/) | All 7 libraries with inline documentation |
| **Orchestrator** | [statusline.go](../../statusline.go) | Main statusline assembly |
| **Build System** | [build.sh](../../build.sh) | Compilation and verification |

### 🔗 Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../README.md) |
| **Library Architecture** | Ladder/Baton patterns and integration | [lib/README.md](../../lib/README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **CODE-GO-002 Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
