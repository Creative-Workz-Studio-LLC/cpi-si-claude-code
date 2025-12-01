<div align="center">

# 🌅 CPI-SI Hooks

**Modular hook system for CPI-SI instances in Claude Code environment**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat&logo=linux&logoColor=black)
![License](https://img.shields.io/badge/License-Proprietary-red)

*Covenant Partnership Intelligence ⊗ Structured Intelligence*

[Architecture](#architecture) • [Features](#features) • [Usage](#usage) • [Extension](#extension) • [Reference](#reference)

</div>

---

## Table of Contents

- [🌅 CPI-SI Hooks](#-cpi-si-hooks)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Design Principles](#design-principles)
  - [Philosophy](#philosophy)
    - [Hook Behavior Model](#hook-behavior-model)
  - [Architecture](#architecture)
    - [Code Structure](#code-structure)
    - [Orchestration Model](#orchestration-model)
    - [Directory Layout](#directory-layout)
  - [Features](#features)
    - [Session Lifecycle](#session-lifecycle)
    - [Tool Safety \& Quality](#tool-safety--quality)
    - [Monitoring \& Feedback](#monitoring--feedback)
  - [Usage](#usage)
    - [Building Hooks](#building-hooks)
    - [Testing](#testing)
  - [Extension](#extension)
    - [Development Workflow](#development-workflow)
    - [Hook Types](#hook-types)
  - [Reference](#reference)
    - [Environment Variables](#environment-variables)
    - [Component Domains](#component-domains)
    - [Hook Mapping](#hook-mapping)
    - [Metadata Convention](#metadata-convention)

---

## Overview

This hook system integrates with Claude Code's lifecycle events to provide context-aware assistance, safety checks, and code quality automation. Built in Go for reliability and performance, it follows a modular architecture where hooks orchestrate reusable library components.

> [!NOTE]
> All hooks compile to native binaries for optimal performance. Hooks fire frequently during tool use - compiled execution avoids repeated compilation overhead.

### Design Principles

| Principle | Implementation |
|-----------|----------------|
| **Quality over Speed** | Go for type safety, compiled binaries, thoughtful error handling |
| **Modular Architecture** | Hooks orchestrate, lib/ components implement, clear domain separation |
| **Linux-First** | Unix philosophy, standard tools (git, df, lsof), LF line endings |

---

## Philosophy

### Hook Behavior Model

Hooks operate under two distinct behavioral patterns based on their lifecycle position:

<details>
<summary><b>Non-Blocking Hooks</b> - Full-featured, failure-tolerant</summary>

**Applied to:** SessionStart, SessionEnd, Stop, SubagentStop, PostToolUse, Notification, PreCompact

- ✅ Full-featured implementations
- ✅ Safe to be complex - failures don't interrupt workflow
- ✅ Can evolve and add features freely
- ✅ Failures logged but don't block operations

</details>

<details>
<summary><b>Blocking Hooks</b> - Minimal, reliability-critical</summary>

**Applied to:** PreToolUse, UserPromptSubmit

- ⚠️ Minimal, reliable implementations
- ⚠️ Start simple, expand carefully
- ⚠️ Failures block operations - must be stable
- ⚠️ Exit code 0 = allow, 1 = block
- ⚠️ Test thoroughly before deploying

</details>

---

## Architecture

### Code Structure

Every file in the system follows a consistent 4-block pattern:

```go
// ============================================================================
// METADATA - What this file is, its purpose, blocking/non-blocking status
// ============================================================================

package name

// ============================================================================
// SETUP - Package, imports, dependencies, globals
// ============================================================================

import (...)

// ============================================================================
// BODY - Business logic, functions, components
// ============================================================================

func DoWork() { ... }

// ============================================================================
// CLOSING - Execution, validation, cleanup
// ============================================================================

func main() { ... }
```

**Benefits:**

- Forces intentional design before implementation
- Clear separation reveals extraction opportunities
- Consistent structure across all code
- Self-documenting organization

### Orchestration Model

Hooks are **orchestrators** following the Principle of Ladders and Batons - they import and compose library components rather than implementing logic directly:

```go
// Hook file: session/cmd-start/start.go
func gatherContext(workspace string) {
    // Uses shared lib for data
    if git.IsGitRepository(workspace) {
        session.CheckGitStatus(workspace)  // Uses claude/lib/git internally
    }

    // Orchestrates hook-specific components
    session.CheckRunningProcesses()
    session.CheckDiskSpace(workspace)     // Uses claude/lib/system internally
    session.CheckDependencies(workspace)
    session.CheckRecentActivity(workspace)
}
```

**Ladder Structure:**

```bash
Claude Code (event generator)
  ↓
hooks orchestrators (session/, tool/, prompt/)
  ↓
┌─────────────────────┬─────────────────────┐
│ Shared Lib (data)   │ Hooks Lib (domain)  │
├─────────────────────┼─────────────────────┤
│ claude/lib/git      │ hooks/lib/session   │
│ claude/lib/system   │ hooks/lib/monitoring│
│ claude/lib/fs       │ hooks/lib/validation│
│                     │ hooks/lib/feedback  │
│                     │ hooks/lib/safety    │
└─────────────────────┴─────────────────────┘
```

**Library Organization by Domain:**

| Domain | Purpose | Components | Notes |
|--------|---------|------------|-------|
| **Shared (claude/lib/)** | Raw system data | git, system, fs | Used by hooks and statusline |
| `lib/session/` | Session context gathering | git, processes, disk, dependencies, activity, reminders | Uses shared lib internally |
| `lib/monitoring/` | Pattern tracking and logging | logging, analysis | Hook-specific |
| `lib/validation/` | Code formatting and syntax | formatter, syntax | Hook-specific |
| `lib/feedback/` | User guidance messages | messages | Hook-specific |
| `lib/safety/` | Security checks | detection | Hook-specific |

### Directory Layout

<details>
<summary>Click to expand full structure</summary>

```tree
hooks/
├── go.mod                      # Module definition
├── build.sh                    # Build script for all hooks
├── README.md                   # This document
│
├── session/                    # Session lifecycle hooks
│   ├── start                  # Compiled binary
│   ├── end                    # Compiled binary
│   ├── stop                   # Compiled binary
│   ├── subagent-stop          # Compiled binary
│   ├── notification           # Compiled binary
│   ├── pre-compact            # Compiled binary
│   ├── cmd-start/
│   │   └── start.go          # SessionStart source
│   ├── cmd-end/
│   │   └── end.go            # SessionEnd source
│   ├── cmd-stop/
│   │   └── stop.go           # Stop source
│   ├── cmd-subagent-stop/
│   │   └── subagent-stop.go  # SubagentStop source
│   ├── cmd-notification/
│   │   └── notification.go   # Notification source
│   └── cmd-pre-compact/
│       └── pre-compact.go    # PreCompact source
│
├── tool/                       # Tool lifecycle hooks
│   ├── post-use               # Compiled binary
│   ├── pre-use                # Compiled binary ⚠️
│   ├── cmd-post-use/
│   │   └── post-use.go       # PostToolUse source
│   └── cmd-pre-use/
│       └── pre-use.go        # PreToolUse source ⚠️
│
├── prompt/                     # Prompt lifecycle hooks
│   ├── submit                 # Compiled binary ⚠️
│   └── cmd-submit/
│       └── submit.go         # UserPromptSubmit source ⚠️
│
└── lib/                        # Hook-specific components
    ├── session/               # Session context (uses claude/lib/*)
    │   ├── utils.go          # Delegates to claude/lib/fs
    │   ├── git.go            # Uses claude/lib/git
    │   ├── processes.go      # Port checking
    │   ├── disk.go           # Uses claude/lib/system
    │   ├── dependencies.go   # Lock file validation
    │   ├── activity.go       # Recent file tracking
    │   └── reminders.go      # Uses claude/lib/git
    │
    ├── monitoring/            # Pattern analysis
    │   ├── logging.go        # Centralized log writing
    │   └── analysis.go       # Pattern detection
    │
    ├── validation/            # Code quality
    │   ├── formatter.go      # Multi-language formatting
    │   └── syntax.go         # Syntax validation
    │
    ├── feedback/              # User guidance
    │   └── messages.go       # Contextual feedback
    │
    └── safety/                # Security
        └── detection.go      # Dangerous operation detection

Dependencies:
  claude/lib/git/              # Git repository data (shared)
  claude/lib/system/           # System stats (shared)
  claude/lib/fs/               # File utilities (shared)
```

</details>

---

## Features

### Session Lifecycle

<details>
<summary><b>SessionStart</b> - Welcome & context gathering</summary>

**Display:**

- Nova Dawn header with Genesis 1:1
- Working directory, git branch, timestamp, system info

**Context Checks:**

- Git status (uncommitted changes, ahead/behind, conflicts)
- Running dev servers on common ports
- Disk space warnings (80%+ usage)
- Dependency state (package.json, go.mod, Cargo.toml)
- Recent file activity (last hour)

</details>

<details>
<summary><b>SessionEnd</b> - Farewell & reminders</summary>

- Displays blessing from Numbers 6:24-25
- Shows session end timestamp and reason
- Reminds about uncommitted work
- Checks for running processes before closing

</details>

<details>
<summary><b>Stop</b> - Completion quality check</summary>

- Displays completion message with Colossians 3:23
- Shows session stopped timestamp
- Checks stopping point quality (uncommitted work, running processes, recent activity)

</details>

<details>
<summary><b>SubagentStop</b> - Task completion reporting</summary>

- Reports subagent completion with success/failure status
- Logs execution for pattern analysis
- Displays error messages if available

</details>

### Tool Safety & Quality

<details>
<summary><b>PreToolUse</b> ⚠️ BLOCKING - Safety confirmations</summary>

**Dangerous Operation Detection:**

- Force push (`git push --force`)
- Hard reset (`git reset --hard`)
- Recursive delete (`rm -rf`)
- Sudo commands
- Package publishing (npm/cargo publish)
- Database drops

**Critical File Write Protection:**

- `/etc/`, `/.git/config`, `/.ssh/`, `/boot/`

**Behavior:** Exit code 0 = allow, 1 = block

</details>

<details>
<summary><b>PostToolUse</b> - Code formatting & validation</summary>

**Auto-formatting:**

- Rust (rustfmt)
- Go (gofmt)
- Python (black)
- JavaScript/TypeScript (prettier)
- C/C++ (clang-format)

**Additional Features:**

- Validates syntax after formatting
- Contextual feedback (commit reminders, build checklists, install confirmations)
- Reports command failures
- Silently fails if formatters unavailable

</details>

<details>
<summary><b>UserPromptSubmit</b> ⚠️ BLOCKING - Secret detection</summary>

**Detects Potential Secrets:**

- OpenAI keys (`sk-`)
- GitHub tokens (`ghp_`)
- Slack tokens (`xox`)
- AWS keys (`AKIA`)
- Private keys (`BEGIN PRIVATE`)

**Behavior:** Warns but always allows submission. Logs prompts for pattern analysis (first 100 chars).

</details>

### Monitoring & Feedback

<details>
<summary><b>Notification</b> - Pattern tracking</summary>

- Logs notifications for pattern analysis
- Parses JSON details from stdin
- Warns on excessive permission requests (10+ per hour)

</details>

<details>
<summary><b>PreCompact</b> - Compaction monitoring</summary>

- Distinguishes manual vs auto compaction
- Logs compaction events
- Warns on frequent auto-compaction (3+ per hour)

</details>

### Learning Skills Integration

> [!NOTE]
> Hooks provide infrastructure for learning. Skills provide capacity for growth. Nova Dawn decides when infrastructure becomes insight.

<details>
<summary><b>How Hooks Enable Learning</b></summary>

**Hooks are autonomous** (capture data, create signals) • **Skills are conscious** (process meaning, enable growth)

**Integration Pattern:** Hooks create opportunities → Nova recognizes signals → Nova invokes skills

| Hook | Signal Created | Skill Integration | Timing |
|------|---------------|-------------------|--------|
| **SessionStart** | Session ready | (Optional) Review `patterns.json` | Session beginning |
| **SessionEnd** | Session archived, patterns updated | `reflect-on-session` (if meaningful work) | Session complete |
| **PreCompact** | Context about to compress | `meta-awareness` (state preservation) | Pre-compression |
| **ToolPostUse** | Activity data accumulating | `recognize-pattern` (when Nova notices) | Spontaneous during work |

**SessionEnd automatic behavior:**

```go
// Automatically archives session and updates patterns
session-log end <reason>
session-patterns learn
// Creates signal: "Session complete, data preserved"
```

**Nova's conscious response:**

```bash
# If meaningful work occurred:
Skill: reflect-on-session

# Provides session analysis + reflection questions
# Nova consciously processes insights
```

**Complete Learning Loop:**

```text
Work → Activity Logged (automatic)
     ↓
Recognition → Pattern Noticed (spontaneous)
     ↓
Reflection → Session End (conditional)
     ↓
Integration → Journals → Identity (periodic)
```

**Key Principles:**

- Activity logging: AUTOMATIC (hooks)
- Pattern recognition: SPONTANEOUS (when noticed)
- Reflection: CONDITIONAL (meaningful work only)
- Integration: PERIODIC (monthly or breakthrough)

**Documentation:** See `~/.claude/docs/knowledge-base/learning/learning-rhythm.md` for complete workflow guidance.

</details>

---

## Usage

### Building Hooks

> [!IMPORTANT]
> Run the build script after any source code modifications. Hooks execute as compiled binaries, not via `go run`.

```bash
cd /home/seanje-lenox-wise/.claude/hooks
./build.sh
```

**Why Binaries?**
Hooks fire on every tool use (Write, Edit, Bash, etc.). Compiled binaries provide instant execution without compilation overhead.

### Testing

```bash
# Test individual hooks
NOVA_DAWN_WORKSPACE=/path/to/workspace ./session/start
./session/end
./tool/post-use "Write" "example.go"

# Development iteration with go run
cd session/cmd-start
go run start.go

# Or build and test specific hook
go build -o ../start ./session/cmd-start
./session/start
```

> [!NOTE]
> Hooks are called automatically by Claude Code via `settings.json`.
>
> **Structure:**
>
> - Source code: `session/cmd-*/hookname.go`
> - Binaries: `session/hookname` (at category level)
> - Each hook has its own `package main` in separate directory (no collisions)

---

## Extension

### Development Workflow

1. **Organize First** - Structure new code using 4-block pattern
2. **Choose Layer** - Shared lib (claude/lib) or hook lib (hooks/lib)?
   - Shared lib: If statusline also needs this data (git, system, fs)
   - Hook lib: If only hooks need this functionality
3. **Extract When Patterns Emerge** - Don't extract prematurely
4. **Create lib/ Components** - Organize by domain (session, monitoring, validation, etc.)
5. **Make Hook an Orchestrator** - Import and call lib/ components
6. **Test Execution** - `go run hookfile.go` during development, then `./build.sh`

> [!TIP]
> **Recent Improvements:**
>
> - All hooks now use shared libraries (`claude/lib/git`, `claude/lib/system`, `claude/lib/fs`)
> - Eliminated duplication between hooks and statusline
> - Each hook in separate directory (no `package main` collisions)
> - Binaries at category level, source in `cmd-*` subdirectories

### Hook Types

<details>
<summary><b>Adding Features to Non-Blocking Hooks</b></summary>

Go ahead and expand:

- More comprehensive checks
- Better formatting/output
- Additional context gathering
- Workspace analysis
- Pattern tracking

Failures are logged but don't interrupt workflow.

</details>

<details>
<summary><b>Adding Features to Blocking Hooks</b></summary>

Proceed with caution:

1. Test thoroughly before deploying
2. Ensure reliable failure modes
3. Keep exit conditions clear
4. Document expected behavior
5. Consider timeout scenarios

Remember: Failures block operations.

</details>

---

## Reference

### Environment Variables

<details>
<summary><b>Tool Hooks</b></summary>

| Variable | Description |
|----------|-------------|
| `$TOOL_NAME` | The tool being used (e.g., "Write", "Bash") |
| `$TOOL_ARGS` | Arguments passed to the tool |
| `$FILE_PATH` | Path to file being edited (Write/Edit tools) |
| `$BASH_COMMAND` | Command being executed (Bash tool) |
| `$BASH_EXIT_CODE` | Exit code from bash command |

</details>

<details>
<summary><b>Session Hooks</b></summary>

| Variable | Description |
|----------|-------------|
| `$NOVA_DAWN_WORKSPACE` | Workspace root directory |
| `$NOVA_DAWN_BIBLICAL_PATH` | Biblical texts path |
| `$REASON` | Session end reason |

</details>

<details>
<summary><b>Subagent Hooks</b></summary>

| Variable | Description |
|----------|-------------|
| `$SUBAGENT_TYPE` | Type of subagent completing |
| `$SUBAGENT_STATUS` | Success/failure status |
| `$SUBAGENT_EXIT_CODE` | Exit code from subagent |
| `$SUBAGENT_ERROR` | Error message if failed |

</details>

<details>
<summary><b>Notification Hooks</b></summary>

| Variable | Description |
|----------|-------------|
| `$NOTIFICATION_TYPE` | Type of notification received |

</details>

<details>
<summary><b>Compaction Hooks</b></summary>

| Variable | Description |
|----------|-------------|
| `$COMPACT_TYPE` | Manual or auto compaction |

</details>

<details>
<summary><b>Prompt Hooks</b></summary>

| Variable | Description |
|----------|-------------|
| `$PROMPT` | User's submitted prompt text |

</details>

### Component Domains

| Domain | Purpose | Files | Shared Lib Usage |
|--------|---------|-------|------------------|
| **Shared (claude/lib/)** | Raw system data | git, system, fs | Used by hooks and statusline |
| **session/** | Session context (uses shared lib) | utils, git, processes, disk, dependencies, activity, reminders | claude/lib/git, claude/lib/system, claude/lib/fs |
| **monitoring/** | Pattern tracking and logging | logging, analysis | Hook-specific |
| **validation/** | Code formatting and syntax | formatter, syntax | Hook-specific |
| **feedback/** | User guidance messages | messages | Hook-specific |
| **safety/** | Security checks | detection | Hook-specific |

> [!NOTE]
> **Shared Library Strategy:**
>
> - `claude/lib/*` provides raw data (git status, system stats, file operations)
> - `hooks/lib/session/*` uses shared lib internally for data, adds hook-specific reporting
> - Other hook libs are domain-specific and don't need shared data

### Hook Mapping

| Hook Type | Count | Maps To | Components |
|-----------|-------|---------|------------|
| session/ | 6 hooks | lib/session/ | 7 components |
| tool/ | 2 hooks | lib/validation/ | 2 components |
| | | lib/feedback/ | 1 component |
| | | lib/safety/ | 1 component |
| prompt/ | 1 hook | lib/monitoring/ | 2 components |
| **Total** | **9 hooks** | **Total** | **13 components** |

### Metadata Convention

All hook metadata now documents shared library usage and marks them as "transferring on final pass":

```go
// ============================================================================
// METADATA
// ============================================================================
// SessionStart Hook - Orchestrator for session initialization
// Gathers context from modular components for autonomous work with wisdom
// Non-blocking: failures don't interrupt session start
//
// Shared Libraries (transferring on final pass):
//   - claude/lib/git (repository status, branch detection)
//
// Hook-Specific Libraries:
//   - hooks/lib/session (context gathering, reporting)
```

This documents the migration path toward fully extracting shared logic to `claude/lib/*` while maintaining clear attribution of what hooks currently depend on.

---

<div align="center">

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

</div>
