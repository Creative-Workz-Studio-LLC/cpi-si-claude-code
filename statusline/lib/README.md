<div align="center">

# 📚 Statusline Libraries

**Presentation Layer Architecture - 7 Specialized Libraries**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Libraries](https://img.shields.io/badge/Libraries-7-orange?style=flat)
![Architecture](https://img.shields.io/badge/Architecture-Ladder%2FBaton-blue?style=flat)

*Transform Claude Code session data into concise, readable statusline display*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [🏗️ Architecture](#-architecture) • [📚 Libraries](#-the-7-libraries) • [🔄 Integration](#-integration-pattern) • [🎯 Design](#-design-principles) • [📚 References](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** Library Architecture Documentation
**Purpose:** Architectural blueprint for 7 statusline presentation libraries
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📚 Statusline Libraries](#-statusline-libraries)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
    - [🎯 What This Document Provides](#-what-this-document-provides)
    - [💡 Design Philosophy](#-design-philosophy)
  - [🏗 Architecture](#-architecture)
    - [🪜 Ladder Structure (Dependencies)](#-ladder-structure-dependencies)
    - [🏃 Baton Pattern (Data Flow)](#-baton-pattern-data-flow)
    - [🔄 Data Transformation](#-data-transformation)
  - [📚 The 7 Libraries](#-the-7-libraries)
    - [📋 1. Types - Session Data Contract](#-1-types---session-data-contract)
    - [⚙ 2. Features - Display Timing](#-2-features---display-timing)
    - [🎨 3. Format - Text Optimization](#-3-format---text-optimization)
    - [🌿 4. Git - Repository Display](#-4-git---repository-display)
    - [📊 5. Session - Statistics Display](#-5-session---statistics-display)
    - [🖥 6. System - Health Display](#-6-system---health-display)
    - [⏰ 7. Temporal - Awareness Display](#-7-temporal---awareness-display)
  - [🔄 Integration Pattern](#-integration-pattern)
    - [🎯 Orchestrator Usage](#-orchestrator-usage)
    - [📦 Import Order](#-import-order)
    - [🔗 Data Dependencies](#-data-dependencies)
  - [🎯 Design Principles](#-design-principles)
    - [1. Single Responsibility](#1-single-responsibility)
    - [2. Stateless Pure Functions](#2-stateless-pure-functions)
    - [3. Graceful Degradation](#3-graceful-degradation)
    - [4. Zero Internal Dependencies](#4-zero-internal-dependencies)
  - [🧪 Demo Strategy](#-demo-strategy)
    - [📦 Library-Level Demos](#-library-level-demos)
    - [🔗 Integration Demo](#-integration-demo)
    - [✅ Build Verification](#-build-verification)
  - [⚡ Performance](#-performance)
  - [🔧 Extension Guide](#-extension-guide)
    - [➕ Adding a New Library](#-adding-a-new-library)
    - [🔨 Adding Functions to Existing Library](#-adding-functions-to-existing-library)
  - [📖 Biblical Foundation](#-biblical-foundation)
    - [Principle](#principle)
    - [Application](#application)
  - [📚 References \& Resources](#-references--resources)
    - [🏗 Library Source Code](#-library-source-code)
    - [📐 Documentation](#-documentation)
    - [🔧 Build System](#-build-system)
  - [📊 Quick Reference](#-quick-reference)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

The statusline library ecosystem provides **7 specialized libraries** for formatting session data into optimal statusline display. Each library handles a specific aspect of presentation - from parsing input data to formatting model names, tracking costs, and displaying system health.

> [!IMPORTANT]
> **This document is the architectural blueprint** for the library ecosystem. For API references and function signatures, see [docs/api/](../docs/api/README.md). For orchestrator integration, see [main README](../README.md).

### 🎯 What This Document Provides

| Section | Purpose | Audience |
|---------|---------|----------|
| **Architecture** | Ladder/Baton patterns and dependency flow | Developers extending the system |
| **Library Catalog** | Complete overview of all 7 libraries | Anyone integrating or using libraries |
| **Integration Pattern** | How orchestrator composes libraries | Orchestrator developers |
| **Design Principles** | Foundational patterns and constraints | Library authors and maintainers |
| **Extension Guide** | How to add libraries or functions | Future developers |

### 💡 Design Philosophy

**Separation of concerns.** Each library has a clear, focused responsibility. Data flows through layers - from raw JSON input to formatted display strings.

| Icon | Principle | Implementation |
|:----:|-----------|----------------|
| 🎯 | **One Concern** | Each library handles exactly ONE aspect of display |
| 🔄 | **Pure Functions** | Same input → same output, no side effects |
| 🛡️ | **Never Fail** | Graceful degradation with safe defaults |
| 🧩 | **Zero Coupling** | No library depends on another (except types) |
| ⚡ | **Performance** | <10ms total overhead for complete assembly |

---

## 🏗 Architecture

### 🪜 Ladder Structure (Dependencies)

The library ecosystem follows the **ladder pattern** - hierarchical organization with clear dependency flow:

```bash
┌─────────────────────────────────────────┐
│  Statusline Orchestrator (statusline.go) │  ← Top Rung (orchestration)
└─────────────────────────────────────────┘
              ↑
    ┌─────────┴─────────┐
    │ Presentation Layer │  ← Middle Rungs (formatting)
    │ - features         │
    │ - format           │
    │ - git              │
    │ - session          │
    │ - system           │
    │ - temporal         │
    └─────────┬─────────┘
              ↑
    ┌─────────┴─────────┐
    │  types (contract)  │  ← Lower Rung (foundation)
    └───────────────────┘
```

**Dependency Flow:**

| Layer | Dependencies | Purpose |
|-------|--------------|---------|
| **Types** (foundation) | None - pure data structures | Session data contract |
| **Presentation** (middle) | types only (+ standard library) | Format specific domains |
| **Orchestrator** (top) | All libraries | Compose final display |

> [!NOTE]
> **Why this matters:** Prevents circular dependencies. Enables independent library development. Clear compilation order. Each rung can only depend on rungs below it.

### 🏃 Baton Pattern (Data Flow)

The **baton pattern** describes how data flows through the system - passed from component to component like a relay race:

```bash
Claude Code Hook
    ↓
JSON Output (stdin)
    ↓
types.SessionContext (parse and structure)
    ↓
Presentation Libraries (format for display)
    ├── features.ShouldShowReminder()      → bool (should show?)
    ├── format.GetShortModelName()         → string (model tier)
    ├── git.GetGitDisplay()                → GitDisplay struct
    ├── session.GetLinesModifiedDisplay()  → LinesDisplay struct
    ├── system.GetLoadDisplay()            → LoadDisplay struct
    └── temporal.GetTimeOfDayDisplay()     → TimeOfDayDisplay struct
    ↓
Statusline Orchestrator (assemble final display)
    ↓
Formatted Output (stdout)
```

**Each component:**

1. Receives input (baton)
2. Performs transformation
3. Returns output (passes baton)

### 🔄 Data Transformation

The complete transformation pipeline:

| Step | Component | Input | Output |
|:----:|-----------|-------|--------|
| 1️⃣ | **Input** | Claude Code | Raw JSON string |
| 2️⃣ | **Parse** | types library | SessionContext struct |
| 3️⃣ | **Format** | Presentation libraries | Display structs |
| 4️⃣ | **Assemble** | Orchestrator | Final statusline string |

---

## 📚 The 7 Libraries

> [!TIP]
> **Complete API documentation available:** See [docs/api/](../docs/api/) for function signatures, parameters, return types, and usage examples for all libraries.

### 📋 1. Types - Session Data Contract

**Path:** [`types/`](types/)
**Type:** Foundation (Ladder lower rung)
**Purpose:** Define input data structure from Claude Code

**Location:** [`types/context.go`](types/context.go)

**Provides:**

- [`SessionContext`](../docs/api/types/context.md) - Complete session data structure

**Usage:**

```go
import "statusline/lib/types"

var ctx types.SessionContext
json.Unmarshal(input, &ctx)
```

**Why It Exists:** Provides a typed contract between Claude Code JSON output and statusline processing. All other libraries depend on this foundation.

**Documentation:** [types API reference](../docs/api/types/README.md)

---

### ⚙ 2. Features - Display Timing

**Path:** [`features/`](features/)
**Type:** Presentation (Ladder middle rung)
**Purpose:** Determine WHEN to show statusline elements

**Location:** [`features/reminder.go`](features/reminder.go)

**Provides:**

- [`ShouldShowReminder(sessionID)`](../docs/api/features/reminder.md) - Kingdom Technology reminder timing (20% frequency)

**Algorithm:** Hash-based deterministic display - session ID hashed to ensure consistent 20% frequency.

**Usage:**

```go
import "statusline/lib/features"

if features.ShouldShowReminder(ctx.SessionID) {
    // Display "⛪ Kingdom Technology" reminder
}
```

**Why It Exists:** Separates timing logic from display logic. Enables consistent, deterministic behavioral decisions across sessions.

**Documentation:** [features API reference](../docs/api/features/README.md)

---

### 🎨 3. Format - Text Optimization

**Path:** [`format/`](format/)
**Type:** Presentation (Ladder middle rung)
**Purpose:** Transform verbose data into concise display forms

**Provides:**

- [`GetShortModelName(displayName)`](../docs/api/format/model.md) - Model tier extraction ("Claude 3.5 Sonnet" → "Sonnet")
- [`ShortenPath(path)`](../docs/api/format/path.md) - Path optimization ("/home/user/project" → "~/project")

**Usage:**

```go
import "statusline/lib/format"

modelName := format.GetShortModelName(ctx.Model.DisplayName)
workdir := format.ShortenPath(ctx.Workspace.CurrentDir)
```

**Why It Exists:** Statuslines have limited space. Long model names and paths consume valuable display real estate. Format library provides consistent text optimization.

**Documentation:** [format API reference](../docs/api/format/README.md)

---

### 🌿 4. Git - Repository Display

**Path:** [`git/`](git/)
**Type:** Presentation (Ladder middle rung)
**Purpose:** Format git repository status for statusline

**Location:** [`git/git.go`](git/git.go)

**Provides:**

- [`GetGitDisplay(gitAssessment)`](../docs/api/git/git.md) - Branch name, dirty status, ahead/behind tracking

**Display Struct:**

```go
type GitDisplay struct {
    HasInfo    bool      // Whether git info is available
    BranchName string    // Current branch name
    Icon       string    // 🌿 git icon
    IsDirty    bool      // Uncommitted changes?
    Ahead      int       // Commits ahead of remote
    Behind     int       // Commits behind remote
}
```

**Usage:**

```go
import "statusline/lib/git"

gitDisplay := git.GetGitDisplay(gitAssessment)
if gitDisplay.HasInfo {
    fmt.Printf("%s %s", gitDisplay.Icon, gitDisplay.BranchName)
    if gitDisplay.IsDirty {
        fmt.Print("*")  // Asterisk for dirty status
    }
}
```

**Why It Exists:** Git context helps developers maintain awareness of repository state during work. Formatted git info provides at-a-glance branch and status visibility.

**Documentation:** [git API reference](../docs/api/git/README.md)

---

### 📊 5. Session - Statistics Display

**Path:** [`session/`](session/)
**Type:** Presentation (Ladder middle rung)
**Purpose:** Format session statistics (lines, duration, cost)

**Provides:**

- [`GetLinesModifiedDisplay(ctx)`](../docs/api/session/lines.md) - Total lines changed (added + removed)
- [`GetDurationDisplay(ctx)`](../docs/api/session/duration.md) - Session elapsed time (45s, 1h 23m)
- [`GetCostDisplay(ctx)`](../docs/api/session/cost.md) - Session cost in USD ($0.0123)
- [`GetFormattedCost(cost)`](../docs/api/session/cost.md) - Helper for cost string formatting

**Usage:**

```go
import "statusline/lib/session"

linesDisplay := session.GetLinesModifiedDisplay(ctx)
durationDisplay := session.GetDurationDisplay(ctx)
costDisplay := session.GetCostDisplay(ctx)

fmt.Printf("📝 %d lines  💰 %s",
    linesDisplay.TotalLines,
    session.GetFormattedCost(costDisplay.Cost))
```

**Why It Exists:** Session statistics provide accountability and awareness - how much work has been done, how long the session has lasted, what the cost is. Essential for sustainable work patterns.

**Documentation:** [session API reference](../docs/api/session/README.md)

---

### 🖥 6. System - Health Display

**Path:** [`system/`](system/)
**Type:** Presentation (Ladder middle rung)
**Purpose:** Format system health metrics with color coding

**Provides:**

- [`GetLoadDisplay()`](../docs/api/system/load.md) - CPU load with health coloring (green/yellow/red)
- [`GetMemoryDisplay()`](../docs/api/system/memory.md) - Memory usage with health coloring
- [`GetDiskDisplay(path)`](../docs/api/system/disk.md) - Disk usage with health coloring

**Health Thresholds:**

| Metric | 🟢 Green | 🟡 Yellow | 🔴 Red |
|--------|---------|----------|--------|
| **Load** | < 50% CPU count | 50-80% | > 80% |
| **Memory** | < 60% used | 60-80% | > 80% |
| **Disk** | < 75% full | 75-90% | > 90% |

**Usage:**

```go
import "statusline/lib/system"

loadDisplay := system.GetLoadDisplay()
memDisplay := system.GetMemoryDisplay()
diskDisplay := system.GetDiskDisplay("/")

if loadDisplay.HasInfo {
    fmt.Printf("%s⚡ %.2f%s", loadDisplay.Color, loadDisplay.LoadAvg, display.Reset)
}
```

**Why It Exists:** System health affects performance and stability. Color-coded indicators provide immediate visual feedback about resource constraints that may impact work.

**Documentation:** [system API reference](../docs/api/system/README.md)

---

### ⏰ 7. Temporal - Awareness Display

**Path:** [`temporal/`](temporal/)
**Type:** Presentation (Ladder middle rung)
**Purpose:** Format temporal awareness (time, schedule, calendar)

**Provides:**

- [`GetTimeOfDayDisplay(ctx)`](../docs/api/temporal/timeofday.md) - Time of day with icon (🌅 morning, ☀️ afternoon, 🌆 evening, 🌙 night)
- [`GetSessionPhaseDisplay(ctx)`](../docs/api/temporal/sessionphase.md) - Session duration and phase (fresh/working/long)
- [`GetScheduleDisplay(ctx)`](../docs/api/temporal/schedule.md) - Current schedule state (work/sleep/personal)
- [`GetCalendarDisplay(ctx)`](../docs/api/temporal/calendar.md) - Date and week number

**Circadian Awareness:** Integrates with CPI-SI [`hooks/lib/temporal`](~/.claude/hooks/lib/temporal/) for schedule and rhythm awareness.

**Usage:**

```go
import "statusline/lib/temporal"
import hookslib "hooks/lib/temporal"

tempCtx, _ := hookslib.GetTemporalContext()
todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
scheduleDisplay := temporal.GetScheduleDisplay(tempCtx)

fmt.Printf("%s%s %s%s  ⏱️ %s (%s)  %s%s %s%s",
    todDisplay.Color, todDisplay.Icon, todDisplay.Label, display.Reset,
    phaseDisplay.Duration, phaseDisplay.Phase,
    scheduleDisplay.Color, scheduleDisplay.Icon, scheduleDisplay.Activity, display.Reset)
```

**Why It Exists:** Temporal awareness enables sustainable work patterns. Knowing time of day, session duration, and schedule state helps maintain healthy rhythms and avoid burnout.

**Documentation:** [temporal API reference](../docs/api/temporal/README.md)

---

## 🔄 Integration Pattern

### 🎯 Orchestrator Usage

The statusline orchestrator ([`statusline.go`](../statusline.go)) integrates all libraries following this pattern:

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "hooks/lib/temporal"                  // Temporal context from hooks
    "system/lib/display"                  // Terminal colors and formatting

    "statusline/lib/features"             // Display timing decisions
    "statusline/lib/format"               // Text optimization
    gitlib "statusline/lib/git"           // Git repository display (aliased)
    sessionlib "statusline/lib/session"   // Session statistics (aliased)
    systemlib "statusline/lib/system"     // System health (aliased)
    temporallib "statusline/lib/temporal" // Temporal awareness (aliased)
    "statusline/lib/types"                // Session data contract
)

func buildStatusline(ctx types.SessionContext) string {
    var parts []string

    // Model and workspace
    modelName := format.GetShortModelName(ctx.Model.DisplayName)
    parts = append(parts, fmt.Sprintf("%s🧠 %s%s",
        display.Cyan, modelName, display.Reset))

    workdir := format.ShortenPath(ctx.Workspace.CurrentDir)
    parts = append(parts, fmt.Sprintf("%s📂 %s%s",
        display.Blue, workdir, display.Reset))

    // Temporal awareness (if available)
    if tempCtx, err := temporal.GetTemporalContext(); err == nil {
        todDisplay := temporallib.GetTimeOfDayDisplay(tempCtx)
        parts = append(parts, fmt.Sprintf("%s%s %s%s",
            todDisplay.Color, todDisplay.Icon, todDisplay.Label, display.Reset))

        scheduleDisplay := temporallib.GetScheduleDisplay(tempCtx)
        if scheduleDisplay.Activity != "" {
            parts = append(parts, fmt.Sprintf("%s%s %s%s",
                scheduleDisplay.Color, scheduleDisplay.Icon,
                scheduleDisplay.Activity, display.Reset))
        }
    }

    // Git status
    gitDisplay := gitlib.GetGitDisplay(ctx.Workspace.CurrentDir)
    if gitDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%s%s %s%s",
            gitDisplay.Color, gitDisplay.Icon,
            gitDisplay.DisplayString, display.Reset))
    }

    // System health
    loadDisplay := systemlib.GetLoadDisplay()
    if loadDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%s%s %.2f%s",
            loadDisplay.Color, loadDisplay.Icon,
            loadDisplay.LoadAvg, display.Reset))
    }

    // Session statistics
    linesDisplay := sessionlib.GetLinesModifiedDisplay(ctx)
    if linesDisplay.HasInfo {
        parts = append(parts, fmt.Sprintf("%s%s %d lines%s",
            linesDisplay.Color, linesDisplay.Icon,
            linesDisplay.TotalLines, display.Reset))
    }

    // Kingdom Technology reminder (conditional)
    if features.ShouldShowReminder(ctx.SessionID) {
        parts = append(parts, fmt.Sprintf("%s⛪ Kingdom Technology%s",
            display.Cyan, display.Reset))
    }

    return strings.Join(parts, "  ")
}

func main() {
    var ctx types.SessionContext
    json.NewDecoder(os.Stdin).Decode(&ctx)
    fmt.Println(buildStatusline(ctx))
}
```

### 📦 Import Order

**Required order for orchestrator:**

| Order | Layer | Imports | Why |
|:-----:|-------|---------|-----|
| 1️⃣ | **Foundation** | [`statusline/lib/types`](types/) | Provides SessionContext - must be first |
| 2️⃣ | **Presentation** | All 6 libraries ([features](features/), [format](format/), [git](git/), [session](session/), [system](system/), [temporal](temporal/)) | Order doesn't matter among these |
| 3️⃣ | **External** | [`hooks/lib/temporal`](~/.claude/hooks/lib/temporal/) | External dependency for temporal context |

### 🔗 Data Dependencies

**Types library** (no dependencies):

- Defines SessionContext
- Used by all presentation libraries
- Pure data structures

**Presentation libraries** (depend on types only):

| Library | Uses from SessionContext | External Dependencies |
|---------|-------------------------|----------------------|
| [features](features/) | `ctx.SessionID` | None |
| [format](format/) | `ctx.Model.DisplayName`, `ctx.Workspace.CurrentDir` | None |
| [git](git/) | `gitAssessment` (passed separately) | None |
| [session](session/) | `ctx.Cost` fields, `ctx.SessionDuration` | None |
| [system](system/) | None (reads `/proc` directly) | None |
| [temporal](temporal/) | None (gets TemporalContext from hooks) | [`hooks/lib/temporal`](~/.claude/hooks/lib/temporal/) |

**Orchestrator** (depends on all):

- Imports all 7 libraries
- Calls presentation functions
- Assembles formatted output
- Handles graceful degradation

---

## 🎯 Design Principles

### 1. Single Responsibility

Each library has **ONE** clear purpose:

| Library | Responsibility | Does NOT Handle |
|---------|---------------|-----------------|
| [types](types/) | Parse input | Formatting, display |
| [features](features/) | Timing decisions | What to display, how to format |
| [format](format/) | Text optimization | When to display, what data means |
| [git](git/) | Git display | Other VCS systems |
| [session](session/) | Session stats | System health, temporal awareness |
| [system](system/) | System health | Session data, git status |
| [temporal](temporal/) | Temporal awareness | Session statistics, system health |

### 2. Stateless Pure Functions

All libraries provide **stateless, deterministic functions**:

| Principle | Implementation | Benefit |
|-----------|----------------|---------|
| **Same input → same output** | No randomness (except deterministic hashing) | Predictable, testable |
| **No side effects** | Read-only operations | Safe to call anytime |
| **No global state** | All data via parameters | Thread-safe, concurrent |
| **No configuration files** | Self-contained logic | Simple deployment |

### 3. Graceful Degradation

Libraries **never fail** - they return safe defaults:

| Scenario | Response | Example |
|----------|----------|---------|
| **Missing data** | `HasInfo: false` in display struct | Git not available → no git display |
| **Unknown inputs** | Reasonable fallback | Unknown model → full name displayed |
| **System errors** | Zero values, empty strings | Can't read load → no load display |
| **Invalid data** | Sanitized safe output | Invalid path → shown as-is |

> [!IMPORTANT]
> **No panics, no errors returned** - statusline rendering must always succeed. Better to show partial info than crash the interaction.

### 4. Zero Internal Dependencies

Presentation libraries depend **ONLY** on:

| Allowed | Dependency | Purpose |
|:-------:|-----------|---------|
| ✅ | [types](types/) library | Session data contract |
| ✅ | Standard library | Core Go functionality |
| ✅ | External ([`hooks/lib/temporal`](~/.claude/hooks/lib/temporal/)) | Temporal context data |
| ❌ | Other presentation libraries | FORBIDDEN - creates coupling |
| ❌ | Third-party packages | Avoid external dependencies |

**Why:** Prevents circular dependencies. Enables independent library development. Clear compilation order. Minimal dependency footprint.

---

## 🧪 Demo Strategy

### 📦 Library-Level Demos

Each library demos **independently**:

```bash
# Demo individual libraries (create demo programs as needed)
go run demo/types_demo.go
go run demo/features_demo.go
go run demo/format_demo.go
go run demo/git_demo.go
go run demo/session_demo.go
go run demo/system_demo.go
go run demo/temporal_demo.go
```

**What to demo:**

- Function correctness (input → output)
- Edge cases (missing data, invalid input)
- Graceful degradation (safe defaults)
- Performance (microsecond-level execution)

### 🔗 Integration Demo

Demo **orchestrator with all libraries**:

```bash
# Demo complete integration with sample JSON
echo '{"session_id":"demo","model":{"display_name":"Claude 3.5 Sonnet"},...}' | ./statusline
```

**What to demo:**

- All libraries work together
- Orchestrator assembles correctly
- End-to-end data flow
- Real JSON input → formatted output

### ✅ Build Verification

Verify **all libraries compile**:

```bash
# Run build script (includes library verification)
./build.sh

# Output:
# Verifying library packages compile...
#   → Checking all packages
#   ✓ All libraries compile cleanly
```

See [`build.sh`](../build.sh) for complete build process.

---

## ⚡ Performance

**All libraries optimized for minimal overhead:**

| Library | Performance | Notes |
|---------|-------------|-------|
| [**types**](types/) | Zero runtime cost | Type definitions only |
| [**features**](features/) | Microseconds | Hash computation |
| [**format**](format/) | Microseconds | String operations |
| [**git**](git/) | Microseconds | String formatting |
| [**session**](session/) | Microseconds | Arithmetic and formatting |
| [**system**](system/) | First call: milliseconds, cached: microseconds | System calls cached |
| [**temporal**](temporal/) | Microseconds | Field extraction |

**Total overhead:** **<10ms** for complete statusline assembly (typical)

**Why performance matters:** Statusline renders on every prompt. Fast execution ensures no delay in interaction start.

---

## 🔧 Extension Guide

### ➕ Adding a New Library

Follow these steps to add a new library to the ecosystem:

| Step | Action | Details |
|:----:|--------|---------|
| 1️⃣ | **Create directory** | `lib/newlib/` with appropriate structure |
| 2️⃣ | **Apply template** | Follow CODE-GO-002 template with 4-block structure |
| 3️⃣ | **Update go.mod** | Add to internal library list |
| 4️⃣ | **Create API docs** | `docs/api/newlib/` with README and function docs |
| 5️⃣ | **Update this README** | Add library entry to catalog |
| 6️⃣ | **Update orchestrator** | Import and use in [`statusline.go`](../statusline.go) |
| 7️⃣ | **Build and test** | Run [`./build.sh`](../build.sh) to verify |

**Template location:** [`~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go`](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go)

### 🔨 Adding Functions to Existing Library

Follow extension points in each library's **BODY section** (documented in source code with comments indicating where to add new functions).

**General process:**

1. Identify which library the function belongs to
2. Follow 4-block structure for new function
3. Return display struct (or pure transformation for format library)
4. Add to orchestrator if needed
5. Document in library's API docs
6. Test independently

---

## 📖 Biblical Foundation

> [!NOTE]
> **Scripture:** "Let all things be done decently and in order" - 1 Corinthians 14:40

### Principle

Clear organization and separation of concerns reflect God's orderly nature.

### Application

Each library has a **distinct role**. No confusion about responsibilities. Data flows cleanly from input → parsing → formatting → assembly. **Order enables understanding.**

**Why this matters for Kingdom Technology:**

- Excellence honors God through intentional design
- Clarity serves others who will maintain this code
- Order reflects the Creator's nature
- Separation of concerns enables sustainable development

---

## 📚 References & Resources

### 🏗 Library Source Code

| Library | Source Location | Purpose |
|---------|----------------|---------|
| **Types** | [`types/context.go`](types/context.go) | Session data contract |
| **Features** | [`features/reminder.go`](features/reminder.go) | Display timing decisions |
| **Format** | [`format/model.go`](format/model.go), [`format/path.go`](format/path.go) | Text optimization functions |
| **Git** | [`git/git.go`](git/git.go) | Git repository display |
| **Session** | [`session/session.go`](session/session.go) | Session statistics formatting |
| **System** | [`system/system.go`](system/system.go) | System health metrics |
| **Temporal** | [`temporal/temporal.go`](temporal/temporal.go) | Temporal awareness display |

### 📐 Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Statusline orchestrator overview | [README.md](../README.md) |
| **API Documentation** | Complete API reference for all libraries | [docs/api/README.md](../docs/api/README.md) |
| **Types API** | SessionContext structure | [docs/api/types/README.md](../docs/api/types/README.md) |
| **Features API** | Display timing functions | [docs/api/features/README.md](../docs/api/features/README.md) |
| **Format API** | Text optimization functions | [docs/api/format/README.md](../docs/api/format/README.md) |
| **Git API** | Git display functions | [docs/api/git/README.md](../docs/api/git/README.md) |
| **Session API** | Session statistics functions | [docs/api/session/README.md](../docs/api/session/README.md) |
| **System API** | System health functions | [docs/api/system/README.md](../docs/api/system/README.md) |
| **Temporal API** | Temporal awareness functions | [docs/api/temporal/README.md](../docs/api/temporal/README.md) |

### 🔧 Build System

| Component | Purpose | Link |
|-----------|---------|------|
| **Build Script** | Automated library verification and compilation | [`build.sh`](../build.sh) |
| **Module Definition** | Go module with 7 internal libraries | [`go.mod`](../go.mod) |
| **Orchestrator** | Main statusline binary source | [`statusline.go`](../statusline.go) |

---

## 📊 Quick Reference

| Library | Purpose | Key Functions |
|---------|---------|---------------|
| 📋 **types** | Input contract | [`SessionContext`](../docs/api/types/context.md) type |
| ⚙️ **features** | Timing decisions | [`ShouldShowReminder()`](../docs/api/features/reminder.md) |
| 🎨 **format** | Text optimization | [`GetShortModelName()`](../docs/api/format/model.md), [`ShortenPath()`](../docs/api/format/path.md) |
| 🌿 **git** | Git display | [`GetGitDisplay()`](../docs/api/git/git.md) |
| 📊 **session** | Session stats | [`GetLinesModifiedDisplay()`](../docs/api/session/lines.md), [`GetDurationDisplay()`](../docs/api/session/duration.md), [`GetCostDisplay()`](../docs/api/session/cost.md) |
| 🖥️ **system** | System health | [`GetLoadDisplay()`](../docs/api/system/load.md), [`GetMemoryDisplay()`](../docs/api/system/memory.md), [`GetDiskDisplay()`](../docs/api/system/disk.md) |
| ⏰ **temporal** | Temporal awareness | [`GetTimeOfDayDisplay()`](../docs/api/temporal/timeofday.md), [`GetSessionPhaseDisplay()`](../docs/api/temporal/sessionphase.md), [`GetScheduleDisplay()`](../docs/api/temporal/schedule.md), [`GetCalendarDisplay()`](../docs/api/temporal/calendar.md) |

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"For God is not the author of confusion, but of peace" - 1 Corinthians 14:33*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
