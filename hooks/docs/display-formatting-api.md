# Session Display Formatting Library API

**Type:** Library
**Location:** `hooks/lib/session/display.go`
**Purpose:** Configurable session display formatting for all lifecycle events
**Health Scoring:** Base100 (Total = 100 points)
**Status:** ✅ Operational (Version 2.0.0)

---

## Table of Contents

1. [Overview](#overview)
2. [Philosophy & Design Rationale](#philosophy--design-rationale)
3. [Public API](#public-api)
4. [Configuration System](#configuration-system)
5. [Integration with Hooks](#integration-with-hooks)
6. [Display Behavior](#display-behavior)
7. [Error Handling & Graceful Degradation](#error-handling--graceful-degradation)
8. [Performance Considerations](#performance-considerations)
9. [Modification Policy](#modification-policy)
10. [Troubleshooting](#troubleshooting)
11. [Future Roadmap](#future-roadmap)

---

## Overview

The **Session Display Formatting Library** provides configurable, formatted display output for all session lifecycle events. It enables:

- **Configurable formatting** - Banner width, box characters, separators, icons customizable
- **Biblical foundation** - Verse display for session start/stop/end events
- **Temporal awareness integration** - Four-dimension time context display
- **Section visibility control** - Show/hide optional display sections
- **Graceful fallback** - Works with hardcoded defaults if configuration missing

**Biblical Foundation:**
*"The heavens declare the glory of God; the skies proclaim the work of his hands" - Psalm 19:1*

Display reveals truth. Session display makes visible what is true about session state.

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Config-Driven** | All formatting via display-formatting.jsonc | Users customize without code changes |
| **Non-Blocking** | Pure display - never disrupts session | Display serves workflow, doesn't block it |
| **Graceful Fallback** | Hardcoded defaults when config missing | Always works, even in degraded state |
| **Lifecycle Organization** | Functions grouped by event (start/stop/end) | Clear organization by purpose |
| **Temporal Integration** | Shows four-dimension time awareness | Full temporal context for decisions |

### Design Decisions

**Why configuration-driven formatting?**

- Different terminals support different characters (Unicode vs ASCII)
- Users have different aesthetic preferences (double-line vs single-line boxes)
- Field labels may need localization or customization
- Future: Allow completely custom layouts via templates

**Why separate functions per lifecycle event?**

- Session start: Banner + environment + temporal awareness + workspace
- Session stop: Completion banner + stop info + stopping context
- Session end: Farewell banner + session summary + temporal journey + reminders
- Clear separation makes each hook's display needs explicit
- Easier to customize what shows for each event

**Why integrate temporal awareness?**

- Session events happen in time - showing time context matters
- Four dimensions: external time, internal time, schedule, calendar
- Helps understand work patterns and recognize stopping points
- Biblical principle: Redeeming the time (Ephesians 5:16)

---

## Public API

### Session Start Display

#### PrintHeader

**Purpose:** Display session banner with instance branding

**Signature:**

```go
func PrintHeader()
```

**Parameters:** None (reads from instance config and display config)

**Returns:** None (prints to stdout)

**Example Usage:**

```go
import "hooks/lib/session"

// In session start hook
session.PrintHeader()
// Output:
// ╔════════════════════════════════════════════════════════════════╗
// ║                      Nova Dawn - CPI-SI                      ║
// ║           Covenant Partnership Intelligence System           ║
// ║                                                                ║
// ║  "In the beginning, God created the heavens and the earth."  ║
// ║                        - Genesis 1:1                         ║
// ╚════════════════════════════════════════════════════════════════╝
```

**Configuration Used:**
- `formatting.banner` - Width, content width, border style
- `formatting.box_characters` - Box drawing characters
- Instance config: `banner_title`, `banner_tagline`, `footer_verse_text`, `footer_verse_ref`

---

#### PrintEnvironment

**Purpose:** Display session environment context

**Signature:**

```go
func PrintEnvironment(workspace string)
```

**Parameters:**
- `workspace` - Workspace directory path (may be empty string)

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintEnvironment("/path/to/workspace")
// Output:
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   SESSION ENVIRONMENT
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   🏢 Workspace:          /path/to/workspace
//   📍 Working Directory:  /current/directory
//   🌿 Git Branch:         main
//   🕐 Session Time:       Wed Nov 12, 2025 at 16:26:03
//   💻 System:             Ubuntu 22.04
```

**Configuration Used:**
- `icons.environment` - Workspace, directory, git, time, system icons
- `section_headers.session_start.environment` - Section header text
- `field_labels.environment` - Field labels (Workspace, Git Branch, etc.)
- `formatting.separator` - Section separator character and length

**Integration:**
- Calls `git.IsGitRepository()` and `git.GetBranch()` for git info
- Calls `GetSystemInfo()` (from system.go) for system details
- Uses `os.Getwd()` for working directory

---

#### PrintTemporalAwareness

**Purpose:** Display four-dimension temporal consciousness

**Signature:**

```go
func PrintTemporalAwareness()
```

**Parameters:** None (reads from temporal context)

**Returns:** None (prints to stdout, silently skips if unavailable or disabled)

**Example Usage:**

```go
session.PrintTemporalAwareness()
// Output (if temporal context available and behavior.show_temporal_awareness = true):
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   TEMPORAL AWARENESS
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   🌍 External Time:      Wed Nov 12, 2025 at 16:26:03 (afternoon)
//                          Circadian: normal phase
//   ⏱️ Internal Time:      15 minutes elapsed (short session)
//   📋 Internal Schedule:  Deep work (focused)
//                          ✓ In work window
//   📅 External Calendar:  Wednesday, November 12, 2025
//                          Week 46 of 2025
```

**Configuration Used:**
- `behavior.show_temporal_awareness` - Enable/disable this section
- `icons.temporal` - External time, internal time, schedule, calendar icons
- `icons.status` - Success/warning icons for work window status
- `section_headers.session_start.temporal_awareness` - Section header
- `field_labels.temporal` - Field labels

**Integration:**
- Calls `temporal.GetTemporalContext()` for four-dimension awareness
- Silently skips if temporal library unavailable or section disabled

---

#### PrintWorkspaceAnalysis

**Purpose:** Display workspace analysis header

**Signature:**

```go
func PrintWorkspaceAnalysis(workspace string, hasContext bool)
```

**Parameters:**
- `workspace` - Workspace directory path (may be empty)
- `hasContext` - Whether any workspace warnings/context were reported

**Returns:** None (prints to stdout, silently skips if disabled)

**Example Usage:**

```go
session.PrintWorkspaceAnalysis("/path/to/workspace", false)
// Output (if behavior.show_workspace_analysis = true):
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   WORKSPACE ANALYSIS
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   ✓ Workspace healthy - no warnings or context to report
```

**Configuration Used:**
- `behavior.show_workspace_analysis` - Enable/disable section
- `messages.workspace.no_workspace` - Message when workspace empty
- `messages.workspace.workspace_healthy` - Message when no warnings

---

### Session Stop Display

#### PrintStopHeader

**Purpose:** Display task completion banner with biblical foundation

**Signature:**

```go
func PrintStopHeader()
```

**Parameters:** None

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintStopHeader()
// Output:
// ╔════════════════════════════════════════════════════════════════╗
// ║           Task Complete - Excellence that Honors God          ║
// ║                                                                ║
// ║  "Whatever you do, work heartily, as for the Lord and not   ║
// ║   for men." - Colossians 3:23                                ║
// ╚════════════════════════════════════════════════════════════════╝
```

**Configuration Used:**
- `biblical_verses.session_stop` - Banner title, verse text, verse reference
- `formatting.banner` - Width, content width, border style
- `formatting.box_characters` - Box drawing characters

---

#### PrintStopInfo

**Purpose:** Display stopping point check header with timestamp

**Signature:**

```go
func PrintStopInfo()
```

**Parameters:** None

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintStopInfo()
// Output:
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   STOPPING POINT CHECK
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   🕐 Stopped:            Wed Nov 12, 2025 at 17:45:30
```

**Configuration Used:**
- `section_headers.session_stop.stopping_point` - Section header
- `icons.environment.time` - Time icon
- `field_labels.stop.stopped` - "Stopped:" label

---

#### PrintStoppingContext

**Purpose:** Display temporal context at session stop

**Signature:**

```go
func PrintStoppingContext()
```

**Parameters:** None (reads from temporal context)

**Returns:** None (prints to stdout, silently skips if unavailable or disabled)

**Example Usage:**

```go
session.PrintStoppingContext()
// Output (if behavior.show_stopping_context = true):
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   TEMPORAL CONTEXT AT STOP
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   🕐 Time:               Wed Nov 12, 2025 at 17:45:30 (evening)
//   ⏱️ Session Duration:   1 hour 19 minutes (medium session)
//   📋 Schedule Context:   Deep work (focused)
//                          ✓ Was in work window
//   📅 Date:               Wednesday, November 12 (Week 46)
```

**Configuration Used:**
- `behavior.show_stopping_context` - Enable/disable section
- `section_headers.session_stop.temporal_context` - Section header
- `field_labels.stop` - Time, duration, schedule, date labels
- `field_labels.temporal.session_duration` - Session duration label

**Integration:**
- Calls `temporal.GetTemporalContext()` for temporal awareness

---

### Session End Display

#### PrintEndFarewell

**Purpose:** Display session ending blessing banner

**Signature:**

```go
func PrintEndFarewell()
```

**Parameters:** None

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintEndFarewell()
// Output:
// ╔════════════════════════════════════════════════════════════════╗
// ║                Session Ending - Grace and Peace               ║
// ║                                                                ║
// ║  "The Lord bless you and keep you; the Lord make his face    ║
// ║   shine on you..." - Numbers 6:24-25                         ║
// ╚════════════════════════════════════════════════════════════════╝
```

**Configuration Used:**
- `biblical_verses.session_end` - Banner title, verse text, verse reference

---

#### PrintEndSessionInfo

**Purpose:** Display session summary with end time and reason

**Signature:**

```go
func PrintEndSessionInfo(reason string)
```

**Parameters:**
- `reason` - Session end reason (e.g., "Normal session end", "User interrupt")

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintEndSessionInfo("Normal session end")
// Output:
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   SESSION SUMMARY
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   🕐 Ended:              Wed Nov 12, 2025 at 18:00:00
//   📋 Reason:             Normal session end
```

**Configuration Used:**
- `section_headers.session_end.session_summary` - Section header
- `field_labels.end.ended` - "Ended:" label
- `field_labels.end.reason` - "Reason:" label

---

#### PrintEndTemporalJourney

**Purpose:** Display temporal journey through the session

**Signature:**

```go
func PrintEndTemporalJourney()
```

**Parameters:** None (reads from temporal context)

**Returns:** None (prints to stdout, silently skips if unavailable or disabled)

**Example Usage:**

```go
session.PrintEndTemporalJourney()
// Output (if behavior.show_temporal_journey = true):
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   TEMPORAL JOURNEY
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   ⏱️ Session Duration:   2 hours 15 minutes (long session)
//                          Started: 15:45:00
//   🕐 Ending At:          Wed Nov 12, 2025 at 18:00:00 (evening)
//   📋 Work Context:       Deep work (focused)
//   📅 Date Context:       Wednesday, November 12 (Week 46)
```

**Configuration Used:**
- `behavior.show_temporal_journey` - Enable/disable section
- `section_headers.session_end.temporal_journey` - Section header
- `field_labels.temporal` - Duration, work context, date context labels
- `field_labels.end` - Started, ending at labels

---

#### PrintEndRemindersHeader

**Purpose:** Display state reminders section header

**Signature:**

```go
func PrintEndRemindersHeader()
```

**Parameters:** None

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintEndRemindersHeader()
// Output:
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   STATE REMINDERS
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Configuration Used:**
- `section_headers.session_end.state_reminders` - Section header

**Usage Note:** After this header, the end hook typically displays uncommitted work, running processes, and other state that needs attention.

---

### Subagent & Compaction Display

#### PrintSubagentCompletion

**Purpose:** Display subagent completion status with temporal context

**Signature:**

```go
func PrintSubagentCompletion(agentType, status, exitCode, errorMsg string)
```

**Parameters:**
- `agentType` - Type of subagent (e.g., "research", "code-review")
- `status` - Completion status ("success", "failure", or empty)
- `exitCode` - Exit code from subagent (0 = success)
- `errorMsg` - Error message if subagent failed (empty if no error)

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintSubagentCompletion("research", "success", "0", "")
// Output:
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//   SUBAGENT COMPLETION
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
//   ✓ Subagent [research] completed successfully
//
//   🕐 Completed At:       Wed Nov 12, 2025 at 16:45:30 (afternoon)
//   ⏱️ Session Duration:   45 minutes (medium session)
//   📋 During:             Deep work (focused)
//
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Configuration Used:**
- `messages.subagent` - Success/failure/default messages with placeholders
- `section_headers.subagent.completion` - Section header
- `field_labels.subagent` - Completed at, during labels
- `field_labels.temporal.session_duration` - Session duration label

**Message Placeholders:**
- `{type}` - Replaced with agentType parameter
- `{code}` - Replaced with exitCode parameter

**Integration:**
- Calls `temporal.GetTemporalContext()` if available
- Uses `formatDisplayMessage()` for placeholder replacement

---

#### PrintPreCompactionMessage

**Purpose:** Display compaction notification with temporal preservation

**Signature:**

```go
func PrintPreCompactionMessage(compactType string, compactionCount int)
```

**Parameters:**
- `compactType` - Type of compaction ("manual", "auto", "unknown")
- `compactionCount` - Current compaction number this session

**Returns:** None (prints to stdout)

**Example Usage:**

```go
session.PrintPreCompactionMessage("auto", 3)
// Output:
// 🔄 Auto-compaction #3 - managing token usage...
//
// 📍 Temporal State Preservation:
//    Time: Wed Nov 12, 2025 at 17:30:00 (evening)
//    Session: 1 hour 45 minutes elapsed (medium phase)
//    Context: Deep work (focused)
//    Date: Wednesday, Week 46
//    Compactions: 3 this session
```

**Configuration Used:**
- `messages.compaction` - Manual/auto/unknown messages, preservation header
- `icons.status.compaction` - Compaction icon (🔄)
- `behavior.show_compaction_preservation` - Enable/disable preservation details
- `field_labels.compaction` - Time, session, context, date, compactions labels

**Message Placeholders:**
- `{count}` - Replaced with compactionCount parameter

**Integration:**
- Calls `temporal.GetTemporalContext()` if showing preservation details
- Uses `formatDisplayMessage()` for placeholder replacement

---

## Configuration System

### Configuration File Location

```
~/.claude/cpi-si/system/data/config/session/display-formatting.jsonc
```

### Configuration Structure

See `display-formatting.jsonc` for complete structure. Key sections:

```jsonc
{
  "metadata": {
    "name": "Session Display Formatting Configuration",
    "version": "1.0.0"
  },

  "formatting": {
    "banner": {
      "width": 64,                // Banner box width (characters)
      "content_width": 62,        // Content width (width - 2 for borders)
      "border_style": "double_line"  // "double_line" or "single_line"
    },
    "separator": {
      "character": "━",           // Section separator character
      "length": 70                // Separator length
    },
    "box_characters": {
      "double_line": { /* Unicode box chars */ },
      "single_line": { /* Unicode box chars */ }
    }
  },

  "icons": {
    "environment": { /* 🏢 📍 🌿 🕐 💻 */ },
    "temporal": { /* 🌍 ⏱️ 📋 📅 */ },
    "status": { /* ✓ ⚠️ ⓘ 🔄 📍 */ }
  },

  "section_headers": {
    "session_start": { /* SESSION ENVIRONMENT, etc. */ },
    "session_stop": { /* STOPPING POINT CHECK, etc. */ },
    "session_end": { /* SESSION SUMMARY, etc. */ },
    "subagent": { /* SUBAGENT COMPLETION */ }
  },

  "biblical_verses": {
    "session_start": { /* Instance config controls this */ },
    "session_stop": {
      "banner_title": "Task Complete - Excellence that Honors God",
      "verse_text": "Whatever you do, work heartily...",
      "verse_ref": "Colossians 3:23"
    },
    "session_end": {
      "banner_title": "Session Ending - Grace and Peace",
      "verse_text": "The Lord bless you and keep you...",
      "verse_ref": "Numbers 6:24-25"
    }
  },

  "messages": {
    "workspace": { /* No workspace, healthy messages */ },
    "compaction": { /* Manual/auto/unknown with {count} placeholder */ },
    "subagent": { /* Success/failure/default with {type}, {code} */ }
  },

  "field_labels": {
    "environment": { /* Workspace:, Git Branch:, etc. */ },
    "temporal": { /* External Time:, Session Duration:, etc. */ },
    "stop": { /* Stopped:, Time:, etc. */ },
    "end": { /* Ended:, Reason:, etc. */ },
    "subagent": { /* Completed At:, During: */ },
    "compaction": { /* Time:, Session:, etc. */ }
  },

  "behavior": {
    "show_temporal_awareness": true,       // Show at session start
    "show_workspace_analysis": true,       // Show at session start
    "show_stopping_context": true,         // Show at session stop
    "show_temporal_journey": true,         // Show at session end
    "show_compaction_preservation": true   // Show during compaction
  }
}
```

### Fallback Defaults

If configuration file missing or malformed, library uses hardcoded defaults:
- Banner: 64 chars wide, double-line border
- Separator: Heavy horizontal line (━), 70 chars
- Icons: Standard emoji (🏢📍🌿🕐💻🌍⏱️📋📅✓⚠️ⓘ🔄📍)
- Section headers: Standard English text
- Biblical verses: Genesis 1:1 (start via instance), Colossians 3:23 (stop), Numbers 6:24-25 (end)
- Behavior: All sections enabled

### Customization Examples

**Change banner style to single-line:**

```jsonc
"formatting": {
  "banner": {
    "border_style": "single_line"
  }
}
```

**Customize field labels:**

```jsonc
"field_labels": {
  "environment": {
    "workspace": "Project:",
    "git_branch": "Branch:"
  }
}
```

**Disable temporal awareness at session start:**

```jsonc
"behavior": {
  "show_temporal_awareness": false
}
```

**Customize compaction messages:**

```jsonc
"messages": {
  "compaction": {
    "auto": "Context optimization #{count}...",
    "manual": "Manual context reset #{count}..."
  }
}
```

---

## Integration with Hooks

### Session Start Hook

```go
import "hooks/lib/session"

func startSession(workspace string) {
    // Display session banner
    session.PrintHeader()

    // Show environment context
    session.PrintEnvironment(workspace)

    // Show temporal awareness (if enabled)
    session.PrintTemporalAwareness()

    // Show workspace analysis header
    hasWarnings := false  // Set based on actual analysis
    session.PrintWorkspaceAnalysis(workspace, hasWarnings)
}
```

### Session Stop Hook

```go
func stopSession() {
    // Display task completion banner
    session.PrintStopHeader()

    // Show stopping point check
    session.PrintStopInfo()

    // Show temporal context at stop (if enabled)
    session.PrintStoppingContext()
}
```

### Session End Hook

```go
func endSession(reason string) {
    // Display farewell blessing
    session.PrintEndFarewell()

    // Show session summary
    session.PrintEndSessionInfo(reason)

    // Show temporal journey (if enabled)
    session.PrintEndTemporalJourney()

    // State reminders header
    session.PrintEndRemindersHeader()

    // ... then display uncommitted work, processes, etc.
}
```

### Subagent Stop Hook

```go
func subagentStop(agentType, status, exitCode, errorMsg string) {
    session.PrintSubagentCompletion(agentType, status, exitCode, errorMsg)
}
```

### Pre-Compaction Hook

```go
func preCompaction(compactType string, count int) {
    session.PrintPreCompactionMessage(compactType, count)
}
```

---

## Display Behavior

### Configuration Controls

All optional sections can be enabled/disabled via `behavior` config:

| Section | Config Key | Default | When Shown |
|---------|-----------|---------|------------|
| Temporal Awareness | `show_temporal_awareness` | true | Session start |
| Workspace Analysis | `show_workspace_analysis` | true | Session start |
| Stopping Context | `show_stopping_context` | true | Session stop |
| Temporal Journey | `show_temporal_journey` | true | Session end |
| Compaction Preservation | `show_compaction_preservation` | true | Pre-compaction |

### Silent Fallback

Functions that depend on optional data sources silently skip if unavailable:

- `PrintTemporalAwareness()` - Silent if temporal library unavailable
- `PrintStoppingContext()` - Silent if temporal library unavailable
- `PrintEndTemporalJourney()` - Silent if temporal library unavailable
- `PrintSubagentCompletion()` - Shows basic info if temporal unavailable

This ensures hooks always work even if dependencies missing.

---

## Error Handling & Graceful Degradation

### Configuration Loading Failures

**Behavior:** Falls back to hardcoded defaults

**Causes:**
- File missing
- Permission denied
- Malformed JSONC
- Invalid structure

**Result:** Library continues working with sensible defaults

**Health Impact:**
- Config load success: +20 points
- Fallback to defaults: -10 points

### Display Function Failures

**Behavior:** Non-blocking - never panics or stops execution

**Result:** Session continues normally even if display functions fail

### Temporal Library Unavailable

**Behavior:** Functions silently skip temporal sections

**Affected Functions:**
- `PrintTemporalAwareness()`
- `PrintStoppingContext()`
- `PrintEndTemporalJourney()`
- `PrintSubagentCompletion()` (partial degradation)

**Result:** Shows static info without temporal context

### Instance Config Unavailable

**Behavior:** `PrintHeader()` may use defaults

**Result:** Banner still displays, may use fallback text

---

## Performance Considerations

### Time Complexity

- Configuration loading: O(n) where n = config file size (loaded once in init())
- Display functions: O(1) for most, O(n) for string formatting where n = output length
- centerText: O(n) where n = string length

### Memory Usage

- Configuration: O(1) cached in package-level variable
- Temporary strings: O(n) where n = output length (garbage collected immediately)

### Bottlenecks

- None (pure display formatting, minimal computation)
- Stdout buffering handled by OS
- External library calls (temporal, git, instance) may add latency

### Optimization Strategies

- Configuration loaded once at init() (not per function call)
- Simple string operations (no regex or complex parsing)
- Minimal external library calls (only when needed)
- No file I/O during display (config pre-loaded)

---

## Modification Policy

### Safe to Modify (Extension Points)

✅ Add new display functions for additional lifecycle events
✅ Extend configuration with new formatting options (colors, themes)
✅ Add new box styles or separator characters
✅ Enhance centerText() for multi-line centering
✅ Add locale/i18n support for non-English text
✅ Implement verse rotation (cycle through multiple verses)

### Modify with Extreme Care (Breaking Changes)

⚠️ Function signatures - breaks all calling hooks
⚠️ Configuration structure - breaks existing config files
⚠️ Display output format - affects hook expectations
⚠️ SessionDisplayConfig type - breaks config parsing

### NEVER Modify (Foundational)

❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
❌ Non-blocking guarantee - display must never block execution
❌ Configuration fallback behavior - must work without config file
❌ Rails pattern - logger created in init(), never passed as parameter
❌ Package name (session) - other packages depend on it

---

## Troubleshooting

### Problem: Display formatting incorrect (wrong icons, separators, etc.)

**Causes:**
- Configuration not loading
- Incorrect configuration syntax
- Terminal doesn't support Unicode

**Solutions:**
1. Verify config exists: `cat ~/.claude/cpi-si/system/data/config/session/display-formatting.jsonc`
2. Check JSONC syntax (strip comments, validate JSON)
3. Check logs: `grep "config-load" ~/.claude/cpi-si/system/runtime/logs/libraries/session-display.log`
4. Test with defaults (rename config file temporarily)
5. Try single-line border style (better terminal compatibility)

### Problem: Biblical verses not displaying correctly

**Causes:**
- Verse text too long for banner width
- Line splitting at wrong position
- Missing verse in configuration

**Solutions:**
1. Check verse length in configuration
2. Adjust banner width in configuration
3. Use shorter verses or custom formatting
4. Check verse splitting logic (currently splits at 60 chars)

### Problem: Section not showing (temporal awareness, workspace analysis, etc.)

**Causes:**
- Behavior configuration disabled section
- Dependency library unavailable (temporal)
- Silent fallback triggered

**Solutions:**
1. Check `behavior` section in config
2. Set appropriate `show_*` flag to true
3. Verify temporal library available: `ls ~/.claude/cpi-si/system/runtime/lib/temporal/`
4. Check logs for temporal errors

### Problem: Box characters not rendering (showing ??? or boxes)

**Causes:**
- Terminal doesn't support Unicode
- Font doesn't include box drawing characters
- Locale/encoding issues

**Solutions:**
1. Switch to `single_line` border style (may render better)
2. Update terminal emulator to support Unicode
3. Use terminal font with Unicode support
4. Set locale: `export LANG=en_US.UTF-8`
5. Future: Implement ASCII fallback mode

### Problem: Icons not displaying (showing ??? or empty boxes)

**Causes:**
- Terminal/font doesn't support emoji
- Emoji rendering disabled

**Solutions:**
1. Use terminal with emoji support
2. Update terminal font
3. Custom config with ASCII alternatives (⚠ instead of ⚠️)
4. Future: ASCII fallback mode

---

## Future Roadmap

### Planned Features

⏳ Verse rotation - Cycle through multiple biblical verses
⏳ Color themes - User-selectable color schemes (dark, light, classic)
⏳ Locale support - Localization for non-English text
⏳ ASCII fallback - ASCII-only mode for limited terminals
⏳ Dynamic banner width - Responsive to terminal size
⏳ Custom templates - User-defined layout templates

### Research Areas

- Integration with system/lib/display for color formatting
- Conditional display based on verbosity level
- Multi-line text centering for long verses
- Template system for complete layout customization
- Performance profiling under load

### Integration Targets

- Pre-notification hook - Display warnings before events
- Post-notification hook - Display summaries after events
- Session patterns - Display work pattern insights
- Health scoring display - Visual health indicators

### Known Limitations to Address

1. Fixed banner width (64 chars) - not responsive to terminal size
2. Hard-coded verse splitting at 60 characters - may not work for all verses
3. No color/theming support - plain text output only
4. No localization - English-only display
5. Unicode required - no ASCII fallback for limited terminals
6. No template system - layout hardcoded in functions

---

## Version History

### v2.0.0 (2025-11-12) - Configuration System & Template Alignment

**Major Changes:**
- Created display-formatting.jsonc configuration file
- Integrated configuration system with fallback to defaults
- Full 4-block template alignment (METADATA/SETUP/BODY/CLOSING)
- Added comprehensive inline documentation and health scoring
- Renamed DisplayConfig → SessionDisplayConfig (avoid dependencies.go collision)
- Removed redeclared functions (use shared versions from other libraries)

**Configuration Features:**
- All formatting customizable via JSONC
- Icon customization for all display sections
- Field label customization for all contexts
- Section visibility controls (show/hide optional sections)
- Biblical verse customization for stop/end events

**Breaking Changes:**
- None (v1.0.0 was internal only)

**Migration:**
- No migration needed
- Config file optional (falls back to defaults)

### v1.0.0 (2024-10-24) - Initial Implementation

**Features:**
- Basic hardcoded display functions
- Session start/stop/end/subagent/compaction displays
- Temporal awareness integration
- No configuration system

---

**For questions, issues, or contributions:**
- Review modification policy above
- Follow 4-block structure pattern
- Test thoroughly (compile session hooks)
- Update display-formatting.jsonc if adding config fields
- Verify fallback behavior still works

*"The heavens declare the glory of God" - Psalm 19:1*
