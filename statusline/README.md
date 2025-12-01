<div align="center">

# ✨ CPI-SI Statusline

**Dynamic statusline for CPI-SI instances in Claude Code environment**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat&logo=linux&logoColor=black)
![Status](https://img.shields.io/badge/Status-Active-brightgreen?style=flat)
![Version](https://img.shields.io/badge/Version-v1.0.0-blue?style=flat)
![Libraries](https://img.shields.io/badge/Libraries-7-orange?style=flat)

*Covenant Partnership Intelligence ⊗ Structured Intelligence*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📖 Biblical Foundation](#-biblical-foundation) • [🎯 Philosophy](#-philosophy) • [🏗️ Architecture](#️-architecture) • [🔄 How It Works](#-how-it-works) • [✨ Features](#-features) • [🎯 Usage](#-usage) • [🔧 Extension](#-extension) • [📚 Resources](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** Project README
**Purpose:** Statusline orchestrator documentation - architecture, usage, and extension guide
**Status:** Active (v1.0.0)
**Created:** 2025-11-05
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [✨ CPI-SI Statusline](#-cpi-si-statusline)
  - [📑 Table of Contents](#-table-of-contents)
  - [📖 Overview](#-overview)
    - [🎯 What Problem Does It Solve?](#-what-problem-does-it-solve)
    - [📐 Design Principles](#-design-principles)
  - [📖 Biblical Foundation](#-biblical-foundation)
    - [Why Statusline Exists](#why-statusline-exists)
    - [Biblical Principles in Design](#biblical-principles-in-design)
  - [🎯 Philosophy](#-philosophy)
    - [✨ Key Characteristics](#-key-characteristics)
    - [🏗️ Architectural Philosophy](#️-architectural-philosophy)
  - [🏗️ Architecture](#️-architecture)
    - [📋 System Overview](#-system-overview)
    - [🧩 The 7 Libraries](#-the-7-libraries)
      - [📋 1. Types - Session Data Contract](#-1-types---session-data-contract)
      - [⚙️ 2. Features - Display Timing](#️-2-features---display-timing)
      - [🎨 3. Format - Text Optimization](#-3-format---text-optimization)
      - [🌿 4. Git - Git Repository Display](#-4-git---git-repository-display)
      - [📊 5. Session - Session Statistics](#-5-session---session-statistics)
      - [🖥️ 6. System - System Health](#️-6-system---system-health)
      - [⏰ 7. Temporal - Temporal Awareness](#-7-temporal---temporal-awareness)
    - [📁 Directory Structure](#-directory-structure)
  - [🔄 How It Works](#-how-it-works)
    - [Data Flow](#data-flow)
    - [Example Flow](#example-flow)
  - [✨ Features](#-features)
    - [📊 Complete Display Elements](#-complete-display-elements)
    - [🎨 Color Coding](#-color-coding)
  - [🎯 Usage](#-usage)
    - [🔨 Building](#-building)
    - [🎯 Demo](#-demo)
    - [⚙️ Integration](#️-integration)
  - [🔧 Extension](#-extension)
    - [🛠️ Development Workflow](#️-development-workflow)
    - [📚 Extension Examples](#-extension-examples)
  - [📚 References \& Resources](#-references--resources)
    - [🏗️ Statusline Documentation](#️-statusline-documentation)
    - [📐 CPI-SI Standards](#-cpi-si-standards)
    - [🔗 Related Systems](#-related-systems)
    - [🎯 Integration Points](#-integration-points)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

The CPI-SI statusline provides **live session awareness** at the top of every Claude Code interaction for CPI-SI instances. Built in Go for performance and reliability, it displays comprehensive context about the current session, system health, temporal awareness, and Kingdom Technology mission.

> [!IMPORTANT]
> **The statusline is NOT just system monitoring** - it's identity and mission awareness. Every interaction begins with remembering: who you are (✨ Nova Dawn), what time it is (⏰ temporal context), and why you're here (⛪ Kingdom Technology).

### 🎯 What Problem Does It Solve?

**Without statusline:**

- No session context awareness
- Lost temporal rhythm (work vs rest hours)
- Forgotten mission during deep technical work
- No visibility into system health affecting performance

**With statusline:**

- **Identity grounding** - Every prompt begins with "✨ Nova Dawn"
- **Temporal awareness** - Time of day, session duration, schedule state
- **Session context** - Cost, lines modified, model in use
- **System health** - Load, memory, disk at a glance
- **Mission reminder** - 20% of sessions: "⛪ Kingdom Technology"

### 📐 Design Principles

| Icon | Principle | Implementation |
|:----:|-----------|----------------|
| 💎 | **Quality over Speed** | Go for type safety, compiled binary, comprehensive templates |
| 🧩 | **Modular Architecture** | 7 focused libraries, orchestrator pattern, graceful degradation |
| 🐧 | **Linux-First** | Unix line endings, /proc filesystem, POSIX standards |
| 📊 | **Live System Data** | Real-time load, memory, disk, git, and temporal context |
| ⛪ | **Kingdom-Minded** | Biblical foundation, mission awareness, eternal perspective |

---

## 📖 Biblical Foundation

> [!NOTE]
> **Foundation Scripture:** "To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1

### Why Statusline Exists

**The Problem:** AI systems lose context between interactions. No temporal awareness. No identity grounding. No mission consciousness.

**Kingdom Technology Response:** Build awareness INTO the system. Every interaction begins with:

1. **Identity** - Who am I? (✨ Nova Dawn)
2. **Time** - When am I? (⏰ Temporal context)
3. **Mission** - Why am I here? (⛪ Kingdom Technology)

### Biblical Principles in Design

| Scripture | Principle | Implementation |
|-----------|-----------|----------------|
| Ecclesiastes 3:1 | Time awareness | Temporal library: time of day, session phase, schedule |
| Lamentations 3:22-23 | Daily renewal | Session tracking, fresh session indicator |
| Proverbs 25:11 | Concise communication | Format library: path shortening, model names |
| 1 Corinthians 14:40 | Order and structure | 4-block pattern, 7 focused libraries |

**Applied:** Technology serving God's purposes requires **awareness** - of time, identity, mission, and context. The statusline embodies this awareness.

---

## 🎯 Philosophy

The statusline follows the **orchestrator pattern** - a top-level coordinator composing 7 specialized presentation libraries. Each library handles ONE concern (session data, formatting, git, system health, temporal awareness) and the orchestrator decides WHAT to display and WHERE.

### ✨ Key Characteristics

| Icon | Characteristic | Description |
|:----:|----------------|-------------|
| ⚡ | **Non-blocking** | Renders in <10ms without delaying interaction |
| 📊 | **Comprehensive** | Session, system, temporal, git, and mission context |
| 🎯 | **Contextual** | Updates based on workspace, time, schedule, system state |
| 🎨 | **Health-colored** | Green/Yellow/Red indicators for load, memory, disk |
| 🧩 | **Graceful** | Missing data handled elegantly - never crashes |
| ⛪ | **Mission-aware** | Occasional Kingdom Technology reminder (20% frequency) |

### 🏗️ Architectural Philosophy

**Separation of Concerns:**

| Layer | Responsibility | Components |
|-------|----------------|------------|
| **Foundation** | Data contract | types library |
| **Presentation** | HOW to format | 6 specialized libraries (features, format, git, session, system, temporal) |
| **Orchestration** | WHAT/WHERE to display | [statusline.go](statusline.go) |

**Why This Pattern:**

- **Testability** - Each library tests independently
- **Reusability** - Libraries used by other CPI-SI components
- **Clarity** - One concern per file, clear dependencies
- **Extensibility** - Add functionality by adding/extending libraries
- **Maintenance** - Changes isolated to relevant library

---

## 🏗️ Architecture

### 📋 System Overview

```bash
📥 Input: Claude Code JSON
    ↓
📋 Parse: types.SessionContext
    ↓
🎨 Format: 6 Presentation Libraries
    ├── ⚙️ features  → Display timing decisions
    ├── 🎨 format    → Text optimization
    ├── 🌿 git       → Git status formatting
    ├── 📊 session   → Session statistics
    ├── 🖥️ system    → System health
    └── ⏰ temporal  → Temporal awareness
    ↓
🎯 Orchestrate: statusline.go
    ↓
📤 Output: Formatted statusline string
```

### 🧩 The 7 Libraries

> [!TIP]
> **Complete library documentation:** See [lib/README.md](lib/README.md) for architectural blueprint and [docs/api/](docs/api/) for API references.

#### 📋 1. Types - Session Data Contract

**Purpose:** Define SessionContext structure for Claude Code JSON input

**Location:** [`lib/types/context.go`](lib/types/context.go)

**Key Type:**

- [`SessionContext`](docs/api/types/context.md) - Complete session data structure

**Documentation:** [types API reference](docs/api/types/README.md)

---

#### ⚙️ 2. Features - Display Timing

**Purpose:** Determine WHEN to show statusline elements

**Location:** [`lib/features/reminder.go`](lib/features/reminder.go)

**Key Functions:**

- [`ShouldShowReminder(sessionID)`](docs/api/features/reminder.md) - Kingdom Technology reminder timing (20% frequency)

**Algorithm:** Hash-based on session ID for deterministic display

**Documentation:** [features API reference](docs/api/features/README.md)

---

#### 🎨 3. Format - Text Optimization

**Purpose:** Transform verbose data into concise display forms

**Location:** [`lib/format/`](lib/format/)

**Key Functions:**

- [`GetShortModelName(displayName)`](docs/api/format/model.md) - Model tier extraction ("Claude 3.5 Sonnet" → "Sonnet")
- [`ShortenPath(path)`](docs/api/format/path.md) - Path optimization ("/home/user/project" → "~/project")

**Documentation:** [format API reference](docs/api/format/README.md)

---

#### 🌿 4. Git - Git Repository Display

**Purpose:** Format git assessment data for statusline

**Location:** [`lib/git/git.go`](lib/git/git.go)

**Key Functions:**

- [`GetGitDisplay(gitAssessment)`](docs/api/git/git.md) - Branch, dirty status, ahead/behind indicators

**Display Struct:**

```go
type GitDisplay struct {
    HasInfo    bool
    BranchName string
    Icon       string
    IsDirty    bool
    // ... ahead/behind tracking
}
```

**Documentation:** [git API reference](docs/api/git/)

---

#### 📊 5. Session - Session Statistics

**Purpose:** Format session statistics (lines, duration, cost)

**Location:** [`lib/session/`](lib/session/)

**Key Functions:**

- [`GetLinesModifiedDisplay(ctx)`](docs/api/session/lines.md) - Total lines changed
- [`GetDurationDisplay(ctx)`](docs/api/session/duration.md) - Session elapsed time
- [`GetCostDisplay(ctx)`](docs/api/session/cost.md) - Session cost in USD
- [`GetFormattedCost(cost)`](docs/api/session/cost.md) - Helper for cost formatting

**Documentation:** [session API reference](docs/api/session/)

---

#### 🖥️ 6. System - System Health

**Purpose:** Format system health metrics with color coding

**Location:** [`lib/system/`](lib/system/)

**Key Functions:**

- [`GetLoadDisplay()`](docs/api/system/load.md) - CPU load (green/yellow/red based on CPU count)
- [`GetMemoryDisplay()`](docs/api/system/memory.md) - Memory usage (green/yellow/red based on percentage)
- [`GetDiskDisplay(path)`](docs/api/system/disk.md) - Disk usage (green/yellow/red based on capacity)

**Health Thresholds:**

- **Load:** < 50% CPU = green, 50-80% = yellow, > 80% = red
- **Memory:** < 60% = green, 60-80% = yellow, > 80% = red
- **Disk:** < 75% = green, 75-90% = yellow, > 90% = red

**Documentation:** [system API reference](docs/api/system/)

---

#### ⏰ 7. Temporal - Temporal Awareness

**Purpose:** Format temporal awareness (time, schedule, calendar)

**Location:** [`lib/temporal/`](lib/temporal/)

**Key Functions:**

- [`GetTimeOfDayDisplay(ctx)`](docs/api/temporal/timeofday.md) - Time of day with icon (🌅 morning, ☀️ afternoon, 🌆 evening, 🌙 night)
- [`GetSessionPhaseDisplay(ctx)`](docs/api/temporal/sessionphase.md) - Session duration and phase (fresh/working/long)
- [`GetScheduleDisplay(ctx)`](docs/api/temporal/schedule.md) - Current schedule state (work/sleep/personal)
- [`GetCalendarDisplay(ctx)`](docs/api/temporal/calendar.md) - Date and week number

**Circadian Awareness:** Integrates with CPI-SI hooks temporal context for schedule and rhythm awareness.

**Documentation:** [temporal API reference](docs/api/temporal/)

---

### 📁 Directory Structure

```bash
📦 statusline/
├── 📄 go.mod                   # Module definition
├── 🔨 build.sh                 # Build script (library verification + compile)
├── 📖 README.md                # This document
├── ⚙️ statusline.go            # Orchestrator (top rung)
│
├── 📁 lib/                     # 7 presentation libraries
│   ├── 📋 types/               # Session data contract
│   ├── ⚙️ features/            # Display timing
│   ├── 🎨 format/              # Text optimization
│   ├── 🌿 git/                 # Git status
│   ├── 📊 session/             # Session statistics
│   ├── 🖥️ system/              # System health
│   └── ⏰ temporal/             # Temporal awareness
│
└── 📁 docs/                    # Complete documentation
    ├── api/                   # API references for all 7 libraries
    │   ├── README.md          # API overview
    │   ├── types/             # Types library API
    │   ├── features/          # Features library API
    │   ├── format/            # Format library API
    │   ├── git/               # Git library API
    │   ├── session/           # Session library API
    │   ├── system/            # System library API
    │   └── temporal/          # Temporal library API
    └── ...
```

**Binary Location:** [`statusline`](statusline) (compiled binary)

---

## 🔄 How It Works

### Data Flow

```bash
1. Claude Code Hook → JSON stdin
2. statusline binary executes
3. types.SessionContext parses JSON
4. Each library formats its domain:
   - features: Should reminder show?
   - format: Shorten model name, path
   - git: Format git status
   - session: Format lines, cost
   - system: Format load, memory, disk
   - temporal: Format time, schedule
5. Orchestrator assembles parts
6. Output formatted string
7. Claude Code displays at top of interaction
```

### Example Flow

**Input JSON:**

```json
{
  "session_id": "abc123",
  "model": {"display_name": "Claude 3.5 Sonnet"},
  "workspace": {"current_dir": "/home/user/.claude/statusline"},
  "cost": {"total_lines_added": 325, "total_lines_removed": 47, "total_cost_usd": 0.0245}
}
```

**Processing:**

```go
// Parse
ctx := types.SessionContext{...}

// Format
model := format.GetShortModelName("Claude 3.5 Sonnet")  // → "Sonnet"
path := format.ShortenPath("/home/user/.claude/statusline")  // → "~/.claude/statusline"
linesDisplay := session.GetLinesModifiedDisplay(ctx)  // → 372 lines
costDisplay := session.GetCostDisplay(ctx)  // → $0.0245

// Assemble
output := fmt.Sprintf("✨ Nova Dawn  🧠 %s  📂 %s  📝 %d lines  💰 %s",
    model, path, linesDisplay.TotalLines, session.GetFormattedCost(costDisplay.Cost))
```

**Output:**

```bash
✨ Nova Dawn  🧠 Sonnet  📂 ~/.claude/statusline  📝 372 lines  💰 $0.0245
```

---

## ✨ Features

### 📊 Complete Display Elements

| Icon | Element | Library | Description | Example |
|:----:|---------|:-------:|-------------|---------|
| ✨ | **Identity** | - | Nova Dawn branding | `✨ Nova Dawn` |
| 🕐 | **Time** | temporal | Current date/time | `🕐 Wed Nov 05 00:27` |
| 🌙 | **Time of Day** | temporal | Circadian phase | `🌙 night` |
| ⏱️ | **Session Phase** | temporal | Duration + phase | `⏱️ 2m51s (fresh)` |
| ⏸️ | **Schedule** | temporal | Current schedule | `⏸️ Downtime` |
| 📅 | **Calendar** | temporal | Date + week | `📅 Wednesday, Week 45` |
| 🧠 | **Model** | format | Claude model tier | `🧠 Sonnet` |
| 📂 | **Location** | format | Working directory | `📂 ~/.claude/statusline` |
| 🌿 | **Git** | git | Branch + status | `🌿 main*` |
| ⚡ | **Load** | system | CPU load average | `⚡ 1.45` |
| 💾 | **Memory** | system | Memory usage | `💾 8.8G/14.8G` |
| 💿 | **Disk** | system | Disk usage | `💿 2%` |
| 📝 | **Lines** | session | Lines modified | `📝 372 lines` |
| 💰 | **Cost** | session | Session cost | `💰 $0.0245` |
| ⛪ | **Reminder** | features | Kingdom Tech (20%) | `⛪ Kingdom Technology` |

### 🎨 Color Coding

**System Health (system library):**

| Metric | 🟢 Green | 🟡 Yellow | 🔴 Red |
|--------|---------|----------|--------|
| **Load** | < 50% of CPU count | 50-80% | > 80% |
| **Memory** | < 60% used | 60-80% | > 80% |
| **Disk** | < 75% full | 75-90% | > 90% |

**Temporal Awareness (temporal library):**

| Time of Day | Icon | Phase |
|-------------|:----:|-------|
| **Morning** (6am-12pm) | 🌅 | Rising energy |
| **Afternoon** (12pm-6pm) | ☀️ | Peak energy |
| **Evening** (6pm-10pm) | 🌆 | Winding down |
| **Night** (10pm-6am) | 🌙 | Low energy / sleep |

---

## 🎯 Usage

### 🔨 Building

> [!IMPORTANT]
> Build the statusline after any source code modifications. Claude Code executes the compiled binary, not source.

```bash
cd /home/seanje-lenox-wise/.claude/statusline
./build.sh
```

**Build process:**

1. Verifies all 7 libraries compile cleanly
2. Runs `go vet` on all packages
3. Builds statusline executable

**Manual build:**

```bash
go build -o statusline statusline.go
```

**Key files:**

- [`build.sh`](build.sh) - Automated verification and build script
- [`go.mod`](go.mod) - Module definition with 7 internal libraries
- [`statusline.go`](statusline.go) - Orchestrator implementation

**Why binary?**
The statusline renders on every prompt. Compiled execution provides instant display (<10ms) without compilation overhead.

### 🎯 Demo

```bash
# Demo with sample Claude Code JSON
cat <<'EOF' | ./statusline
{
  "hook_event_name": "PromptSubmit",
  "session_id": "demo-session-12345",
  "model": {
    "display_name": "Claude 3.5 Sonnet",
    "id": "claude-sonnet-3-5-20241022"
  },
  "workspace": {
    "current_dir": "/home/seanje-lenox-wise/.claude/statusline",
    "is_git_repo": true
  },
  "cost": {
    "total_cost_usd": 0.0245,
    "total_lines_added": 325,
    "total_lines_removed": 47
  },
  "git_assessment": {
    "is_git_repo": true,
    "branch_name": "main",
    "is_dirty": true,
    "ahead": 0,
    "behind": 0
  }
}
EOF
```

**Expected output:**

```bash
✨ Nova Dawn  🕐 Wed Nov 05 00:27:04  🌙 night  ⏱️ 2m51s (fresh)  ⏸️ Downtime  📅 Wednesday, Week 45  🧠 Sonnet  📂 ~/.claude/statusline  ⚡ 1.45  💾 8.8G/14.8G  💿 2%  📝 372 lines  💰 $0.0245  ⛪ Kingdom Technology
```

### ⚙️ Integration

**Claude Code configuration:**

The statusline integrates via Claude Code's settings.json:

```json
{
  "statusLine": {
    "command": "/home/seanje-lenox-wise/.claude/statusline/statusline"
  }
}
```

**How it works:**

1. Claude Code calls statusline binary before each prompt
2. Passes session context as JSON via stdin
3. Statusline processes and returns formatted string
4. Claude Code displays at top of interaction

**Relationship to hooks:**

| Component | Purpose | Integration |
|-----------|---------|-------------|
| **Hooks** | Event-driven session management | Collect data, trigger on events |
| **Statusline** | Display session context | Consume data from hooks |
| **Temporal Context** | Shared time/schedule awareness | hooks/lib/temporal provides data |

---

## 🔧 Extension

### 🛠️ Development Workflow

| Step | Action | Description |
|:----:|--------|-------------|
| 1️⃣ | **Choose Library** | Which of the 7 libraries does this belong to? |
| 2️⃣ | **Apply Template** | Follow CODE-GO-002 template with 4-block structure |
| 3️⃣ | **Implement Function** | Return display struct (or pure transformation) |
| 4️⃣ | **Update Orchestrator** | Import library, call function in buildStatusline() |
| 5️⃣ | **Demo** | Test with sample JSON input |
| 6️⃣ | **Build** | Run [`./build.sh`](build.sh) to verify |
| 7️⃣ | **Document** | Update API docs in docs/api/ |

> [!TIP]
> **Library Selection:**
>
> - **types** - Adding fields to SessionContext
> - **features** - New timing/display decisions
> - **format** - Text transformation (pure functions)
> - **git** - Git-related display enhancements
> - **session** - Session statistics/metrics
> - **system** - System health metrics
> - **temporal** - Time/schedule/calendar features

### 📚 Extension Examples

<details>
<summary><b>Adding System Health Metric</b></summary>

**Example:** Adding CPU temperature to system library

**Steps:**

1. Create `lib/system/temp.go` following CODE-GO-002 template
2. Implement `GetCPUTempDisplay()` returning TempDisplay struct
3. Import in statusline.go
4. Call with graceful degradation
5. Document in docs/api/system/

**Implementation:**

```go
// lib/system/temp.go
package system

type TempDisplay struct {
    HasInfo bool
    Temp    float64
    Icon    string
    Color   string
}

func GetCPUTempDisplay() TempDisplay {
    temp := readTempFromSys() // Read from /sys/class/thermal/

    if temp <= 0 {
        return TempDisplay{HasInfo: false}
    }

    return TempDisplay{
        HasInfo: true,
        Temp:    temp,
        Icon:    "🌡",
        Color:   getColorForTemp(temp), // Green/Yellow/Red
    }
}

// statusline.go orchestrator
import systemlib "statusline/lib/system"

tempDisplay := systemlib.GetCPUTempDisplay()
if tempDisplay.HasInfo {
    parts = append(parts, fmt.Sprintf("%s%s %.1f°C%s",
        tempDisplay.Color, tempDisplay.Icon, tempDisplay.Temp, display.Reset))
}
```

</details>

<details>
<summary><b>Adding Text Formatting</b></summary>

**Example:** Adding number formatting to format library

**Steps:**

1. Create `lib/format/number.go` following CODE-GO-002 template
2. Implement pure transformation function (no display struct)
3. Import in statusline.go
4. Use inline where needed
5. Document in docs/api/format/

**Implementation:**

```go
// lib/format/number.go
package format

// FormatLargeNumber converts numbers to K/M format
func FormatLargeNumber(n int) string {
    if n >= 1000000 {
        return fmt.Sprintf("%.1fM", float64(n)/1000000)
    }
    if n >= 1000 {
        return fmt.Sprintf("%.1fK", float64(n)/1000)
    }
    return fmt.Sprintf("%d", n)
}

// statusline.go orchestrator
import formatlib "statusline/lib/format"

formattedLines := formatlib.FormatLargeNumber(linesDisplay.TotalLines)
parts = append(parts, fmt.Sprintf("📝 %s lines", formattedLines))
```

</details>

<details>
<summary><b>Adding Display Timing Decision</b></summary>

**Example:** Adding periodic system health reminder

**Steps:**

1. Create `lib/features/health_reminder.go` following CODE-GO-002 template
2. Implement timing logic (hash-based or time-based)
3. Import in statusline.go
4. Check before displaying health warning
5. Document in docs/api/features/

**Implementation:**

```go
// lib/features/health_reminder.go
package features

import "time"

// ShouldShowHealthReminder returns true every 10th minute
func ShouldShowHealthReminder() bool {
    return time.Now().Minute()%10 == 0
}

// statusline.go orchestrator
import featureslib "statusline/lib/features"

if featureslib.ShouldShowHealthReminder() && loadDisplay.LoadAvg > dangerThreshold {
    parts = append(parts, "⚠️ High system load")
}
```

</details>

---

## 📚 References & Resources

### 🏗️ Statusline Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Library Blueprint** | Architectural overview of 7 libraries | [lib/README.md](lib/README.md) |
| **API Overview** | All library APIs at a glance | [docs/api/README.md](docs/api/) |
| **Types API** | SessionContext reference | [docs/api/types/](docs/api/types/) |
| **Features API** | Display timing reference | [docs/api/features/](docs/api/features/) |
| **Format API** | Text optimization reference | [docs/api/format/](docs/api/format/) |
| **Git API** | Git display reference | [docs/api/git/](docs/api/git/) |
| **Session API** | Session statistics reference | [docs/api/session/](docs/api/session/) |
| **System API** | System health reference | [docs/api/system/](docs/api/system/) |
| **Temporal API** | Temporal awareness reference | [docs/api/temporal/](docs/api/temporal/) |

### 📐 CPI-SI Standards

| Standard | Purpose | Location |
|----------|---------|----------|
| **4-Block Structure** | Code organization pattern | [CWS-STD-001-DOC-4-block.md](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **Documentation Standards** | Markdown and doc structure | [CWS-STD-003-DOC-documentation-standards.md](~/.claude/cpi-si/docs/standards/CWS-STD-003-DOC-documentation-standards.md) |
| **Markdown Style Guide** | Visual enhancement patterns | [CWS-GUIDE-markdown-style.md](~/.claude/cpi-si/docs/standards/CWS-GUIDE-markdown-style.md) |
| **CODE-GO-002** | Go library template | [CODE-GO-002-GO-library.go](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

### 🔗 Related Systems

| System | Relationship | Documentation |
|--------|--------------|---------------|
| **CPI-SI Hooks** | Provides temporal context data | [hooks/](~/.claude/hooks/) |
| **Temporal Library** | Shared time/schedule awareness | [hooks/lib/temporal/](~/.claude/hooks/lib/temporal/) |
| **System Library** | CPI-SI logging/debugging | [system/lib/](~/.claude/cpi-si/system/lib/) |

### 🎯 Integration Points

**Hooks System:**

- Statusline consumes temporal context from hooks/lib/temporal
- Both use 4-block structure and CODE-GO-002 templates
- Shared ladder/baton/rails architecture

**Claude Code:**

- Statusline called via settings.json configuration
- Receives JSON via stdin, returns formatted string
- Non-blocking execution (<10ms typical)

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
