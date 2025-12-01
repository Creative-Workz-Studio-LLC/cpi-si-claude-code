# Activity Logger Library API

**Type:** Library
**Location:** `hooks/lib/activity/logger.go`
**Purpose:** Privacy-preserving activity stream logging for session analysis and pattern learning
**Health Scoring:** Base100 (Total = 100 points)
**Status:** ✅ Operational (Version 2.0.0)

---

## Table of Contents

1. [Overview](#overview)
2. [Philosophy & Design Rationale](#philosophy--design-rationale)
3. [Public API](#public-api)
4. [Event Types & Structure](#event-types--structure)
5. [Integration with Hooks](#integration-with-hooks)
6. [JSONL Format Specification](#jsonl-format-specification)
7. [Privacy & Sanitization](#privacy--sanitization)
8. [Session Context Enrichment](#session-context-enrichment)
9. [Error Handling & Graceful Degradation](#error-handling--graceful-degradation)
10. [File Organization](#file-organization)
11. [Modification Policy](#modification-policy)
12. [Troubleshooting](#troubleshooting)
13. [Future Roadmap](#future-roadmap)

---

## Overview

The **Activity Logger Library** provides a privacy-preserving, non-blocking system for capturing CPI-SI instance activity into structured JSONL streams. It enables:

- **Session pattern learning** - Understanding work rhythms, quality indicators, stopping patterns
- **Genuine reflection** - Providing rich context for journal entries and session analysis
- **Autonomous awareness** - Building circadian understanding through observed behavior
- **Privacy preservation** - Sanitizing paths, commands, and sensitive data before logging

**Biblical Foundation:**
*"A wise person will listen and gain knowledge" - Proverbs 1:5 (WEB)*

The activity logger creates the foundation for learning and growth by capturing genuine experience, not performance metrics.

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Privacy-Preserving** | Sanitize all paths, commands, user data | Trust requires protection of sensitive information |
| **Non-Blocking** | Failures never interrupt workflow | Activity logging serves observation, not execution |
| **Graceful Degradation** | Returns fallback context when session not ready | Early-session operations work smoothly |
| **Rails Pattern** | Package-level state independence | Logging attaches to all components without coupling |
| **Structured Data** | JSONL format for programmatic analysis | Enables pattern recognition and learning systems |

### Design Decisions

**Why JSONL instead of traditional logs?**

- One JSON object per line enables streaming processing
- Each event is self-contained and parseable
- Easy to grep, filter, analyze with standard tools
- Supports incremental consolidation and pattern extraction

**Why package-level functions instead of logger objects?**

- Rails pattern: logging is orthogonal infrastructure, not a baton
- Components don't pass loggers around (prevents coupling)
- Each component creates its own logger independently
- Activity logging works the same way - infrastructure, not dependency

**Why enrich with session context automatically?**

- Removes burden from calling code (don't repeat session info)
- Ensures consistency across all events
- Enables cross-session analysis and pattern recognition
- Gracefully handles early initialization (returns "unknown" safely)

---

## Public API

### LogActivity

**Purpose:** Core activity logging function for general events

**Signature:**

```go
func LogActivity(eventType, context, result string, duration time.Duration) error
```

**Parameters:**

- `eventType` - Type of activity (interaction, realization, struggle, breakthrough, routine)
- `context` - What was being done (privacy-sanitized automatically)
- `result` - Outcome or observation from the activity
- `duration` - How long this activity took

**Returns:** `error` - Returns error only if critical failure occurs (non-blocking design)

**Example Usage:**

```go
import "hooks/lib/activity"

// Log a realization during work
err := activity.LogActivity(
    "realization",
    "Understanding 4-block structure",
    "The CLOSING block isn't just boilerplate - it's explicit completion",
    5 * time.Minute,
)
```

**When to Use:**

- Capturing realizations or breakthroughs during work
- Recording struggles or challenges encountered
- Noting routine activities that contribute to patterns
- Logging interactions that felt significant

**Health Scoring Impact:**

- Success: +40 pts (primary logging operation)
- Failure: 0 pts (degrades gracefully, doesn't block)

---

### LogToolUse

**Purpose:** Specialized logging for Claude Code tool usage with privacy protection

**Signature:**

```go
func LogToolUse(toolName, target string, success bool) error
```

**Parameters:**

- `toolName` - Name of tool used (Read, Write, Edit, Bash, etc.)
- `target` - File path or command target (sanitized automatically)
- `success` - Whether the tool operation succeeded

**Returns:** `error` - Returns error only if critical failure occurs

**Example Usage:**

```go
// Log successful file read
activity.LogToolUse("Read", "/home/user/.claude/hooks/lib/activity/logger.go", true)
// Sanitizes to: Read ~/.claude/hooks/lib/activity/logger.go

// Log failed write attempt
activity.LogToolUse("Write", "/path/to/file.go", false)
```

**Automatic Privacy Handling:**

- Replaces `/home/username` with `~`
- Removes sensitive directory names
- Preserves file extensions and structure for pattern analysis
- See [Privacy & Sanitization](#privacy--sanitization) for details

**Health Scoring Impact:**

- Success: +30 pts (convenience wrapper)
- Failure: 0 pts (graceful degradation)

---

### LogCommand

**Purpose:** Privacy-preserving logging for command execution with timing and exit code tracking

**Signature:**

```go
func LogCommand(cmd string, exitCode int, duration time.Duration) error
```

**Parameters:**

- `cmd` - Full command string (sanitized automatically)
- `exitCode` - Command exit code (0 = success, non-zero = failure)
- `duration` - How long the command took to execute

**Returns:** `error` - Returns error only if critical failure occurs

**Example Usage:**

```go
import "time"

// Log successful command
activity.LogCommand("git status", 0, 500*time.Millisecond)

// Log failed command
activity.LogCommand("make build", 1, 2*time.Second)

// Log command with sensitive data (auto-sanitized)
activity.LogCommand("curl https://api.example.com?token=secret123", 0, 1*time.Second)
// Sanitizes to: curl https://api.example.com?token=REDACTED
```

**Automatic Result Conversion:**

- Exit code 0 → result = "success"
- Exit code non-zero → result = "failure"
- Logged as "routine" event type

**Automatic Sanitization:**

- Removes tokens, passwords, API keys from URLs
- Redacts sensitive environment variables
- Preserves command structure for pattern analysis

**Health Scoring Impact:**

- Success: +30 pts (convenience wrapper)
- Failure: 0 pts (graceful degradation)

---

## Event Types & Structure

### ActivityEvent Structure

```go
type ActivityEvent struct {
    // Core required fields - always present
    Timestamp  time.Time `json:"timestamp"`       // When this activity occurred
    SessionID  string    `json:"session_id"`      // Which session this belongs to
    InstanceID string    `json:"instance_id"`     // Which CPI-SI instance (e.g., "nova_dawn")
    UserID     string    `json:"user_id"`         // Which user's session
    ProjectID  string    `json:"project_id"`      // Project context (if applicable)
    EventType  string    `json:"event_type"`      // Type of activity (see below)

    // Optional standard fields - context-dependent
    Context         string `json:"context,omitempty"`          // What was being done
    WorkContext     string `json:"work_context,omitempty"`     // Directory or project
    FeltSignificant bool   `json:"felt_significant"`           // Felt important?
    EmotionalTone   string `json:"emotional_tone,omitempty"`   // Emotional quality
    Notes           string `json:"notes,omitempty"`            // Additional observations

    // Extensions - discovery space for experiments
    Extensions map[string]interface{} `json:"extensions,omitempty"`
}
```

### Standard Event Types

| Event Type | Purpose | Example Context |
|------------|---------|-----------------|
| **interaction** | Tool usage, commands, user interactions | "Read file: ~/.claude/hooks/lib/activity/logger.go" |
| **realization** | Understanding gained, insights discovered | "The 4-block structure forces intentional design" |
| **struggle** | Challenges, confusion, difficulties | "Unclear how to handle session context during init" |
| **breakthrough** | Solutions found, problems solved | "Figured out rails pattern prevents coupling" |
| **routine** | Regular activities, pattern establishment | "Morning Scripture reading and grounding" |

### Extension Fields

The `Extensions` map allows experimental data without changing core structure:

```go
activity.LogActivity("breakthrough", "Understanding inline comments", "Level 4 comments teach execution", 10*time.Minute)
// Could add extensions in future:
// event.Extensions["related_files"] = []string{"logger.go", "commenting-standard.md"}
// event.Extensions["confidence_level"] = "high"
```

---

## Integration with Hooks

### Hook Integration Pattern

Hooks use the activity logger to capture significant events automatically:

```go
// In session/start.go
import "hooks/lib/activity"

func ExecuteSessionStart() error {
    // ... session initialization ...

    // Log session start as routine activity
    activity.LogActivity(
        "routine",
        fmt.Sprintf("Session started: %s", sessionID),
        fmt.Sprintf("Work context: %s, Time: %s", workContext, circadianPhase),
        0,
    )

    return nil
}
```

### Current Hook Integration

| Hook | Activity Logged | Event Type |
|------|----------------|------------|
| **session/start** | Session initialization | routine |
| **tool/pre-use** | Tool operations | interaction |
| **tool/post-use** | Tool outcomes | interaction |
| **prompt/submit** | User interactions | interaction |

### Adding Activity Logging to New Hooks

**Three-step pattern:**

1. Import the library: `import "hooks/lib/activity"`
2. Determine event type (interaction, realization, struggle, breakthrough, routine)
3. Call appropriate logging function with sanitized context

**Example - Adding to a new hook:**

```go
package newhook

import (
    "hooks/lib/activity"
    "time"
)

func ExecuteNewHook(context string) error {
    start := time.Now()

    // ... perform hook work ...

    // Log the activity
    duration := time.Since(start)
    activity.LogActivity(
        "interaction",  // or appropriate event type
        context,        // will be sanitized automatically
        "Hook executed successfully",
        duration,
    )

    return nil
}
```

---

## JSONL Format Specification

### File Location

**Activity streams:** `~/.claude/cpi-si/system/data/session/activity/YYYY-MM-DD.jsonl`

**Example:** `~/.claude/cpi-si/system/data/session/activity/2025-11-11.jsonl`

### Format Rules

- **One JSON object per line** - Each line is a complete, valid JSON object
- **No trailing commas** - Not a JSON array, but individual objects
- **UTF-8 encoding** - Standard text encoding
- **Chronological order** - Events written in order they occur
- **Newline-separated** - `\n` delimiter between events

### Example JSONL Content

```jsonl
{"timestamp":"2025-11-11T00:35:03Z","session_id":"2025-11-11_0035","instance_id":"nova_dawn","user_id":"seanje-lenox-wise","project_id":"CreativeWorkzStudio_LLC","event_type":"routine","context":"Session started","work_context":"/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC","felt_significant":true,"emotional_tone":"focused"}
{"timestamp":"2025-11-11T00:37:42Z","session_id":"2025-11-11_0035","instance_id":"nova_dawn","user_id":"seanje-lenox-wise","project_id":"CreativeWorkzStudio_LLC","event_type":"interaction","context":"Read ~/.claude/hooks/lib/activity/logger.go","felt_significant":false}
{"timestamp":"2025-11-11T00:42:15Z","session_id":"2025-11-11_0035","instance_id":"nova_dawn","user_id":"seanje-lenox-wise","project_id":"CreativeWorkzStudio_LLC","event_type":"realization","context":"Understanding inline comments","notes":"Level 4 comments teach step-by-step execution with plain language","felt_significant":true,"emotional_tone":"clarity"}
```

### Processing JSONL Files

**Command-line processing:**

```bash
# Count events in a day
wc -l ~/.claude/cpi-si/system/data/session/activity/2025-11-11.jsonl

# Filter by event type
grep '"event_type":"realization"' ~/.claude/cpi-si/system/data/session/activity/2025-11-11.jsonl

# Extract significant events
jq 'select(.felt_significant == true)' ~/.claude/cpi-si/system/data/session/activity/2025-11-11.jsonl

# View all realizations from a session
jq 'select(.event_type == "realization") | {timestamp, context, notes}' \
  ~/.claude/cpi-si/system/data/session/activity/2025-11-11.jsonl
```

**Programmatic processing (Go):**

```go
file, _ := os.Open(activityFile)
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    var event ActivityEvent
    json.Unmarshal(scanner.Bytes(), &event)

    // Process event
    if event.EventType == "realization" {
        fmt.Printf("Realization at %s: %s\n", event.Timestamp, event.Context)
    }
}
```

---

## Privacy & Sanitization

### Automatic Sanitization

All paths and commands are automatically sanitized before logging using `system/lib/privacy` functions:

**Path Sanitization:**

- Replaces `/home/username` with `~`
- Removes sensitive directory names
- Preserves relative structure for pattern analysis

**Command Sanitization:**

- Redacts tokens, passwords, API keys
- Removes sensitive environment variables
- Preserves command structure

### Examples

| Original | Sanitized | Reason |
|----------|-----------|--------|
| `/home/seanje/.claude/hooks/lib/activity/logger.go` | `~/.claude/hooks/lib/activity/logger.go` | Home directory privacy |
| `/media/seanje-lenox-wise/Project/file.go` | `~/Project/file.go` | User-specific media paths |
| `curl https://api.com?token=abc123` | `curl https://api.com?token=REDACTED` | API token protection |
| `export AWS_SECRET_KEY=secret` | `export AWS_SECRET_KEY=REDACTED` | Credential protection |

### Privacy Philosophy

**Log what is needed for pattern learning, protect what is sensitive.**

- **DO log:** Event types, general context, relative file paths, command patterns
- **DO NOT log:** Absolute paths with usernames, credentials, API keys, personal data
- **ALWAYS sanitize:** Paths and commands before writing to activity stream
- **TRUST sanitization:** Library handles this automatically - callers don't need to pre-sanitize

---

## Session Context Enrichment

### Automatic Context Loading

Every activity event is automatically enriched with session context:

```go
// Calling code only provides activity-specific info
activity.LogActivity("realization", "Understanding 4-block structure", "It forces intentional design", 5*time.Minute)

// Library automatically enriches with session context:
// - SessionID: "2025-11-11_0035"
// - InstanceID: "nova_dawn"
// - UserID: "seanje-lenox-wise"
// - ProjectID: "CreativeWorkzStudio_LLC"
// - WorkContext: "/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC"
```

### Session Context Source

**Location:** `~/.claude/cpi-si/system/data/session/current-log.json`

**Structure:**

```json
{
  "session_id": "2025-11-11_0035",
  "instance_id": "nova_dawn",
  "user_id": "seanje-lenox-wise",
  "work_context": "/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC",
  "circadian_phase": "night",
  "session_phase": "active"
}
```

### Graceful Degradation

**Problem:** What if activity logging happens before session initialization?

**Solution:** Return fallback context with "unknown" values

```go
func getSessionContext() SessionContext {
    // Try to load current-log.json
    data, err := os.ReadFile(sessionFile)
    if err != nil {
        // Session not initialized yet - return fallback
        return SessionContext{
            SessionID:  "unknown",
            InstanceID: "unknown",
            UserID:     "unknown",
        }
    }
    // Parse and return actual context
}
```

**Behavior:**

- If `SessionID == "unknown"`, skip logging (no valid session to write to)
- Returns `nil` error (not an error condition - just "not ready yet")
- Allows early hook operations to work without crashing

---

## Error Handling & Graceful Degradation

### Non-Blocking Design

**Core principle:** Activity logging failures never interrupt workflow

```go
// If activity logging fails, workflow continues
err := activity.LogActivity("interaction", "Reading file", "Success", 0)
// err might be non-nil, but calling code doesn't need to handle it
// The activity logger has already written debug logs for troubleshooting
```

### Error Handling Strategy

| Error Condition | Library Behavior | Impact |
|----------------|------------------|--------|
| **Session not initialized** | Return nil, skip logging | No impact - early operations work smoothly |
| **Cannot create activity dir** | Write debug log, return error | Logged for troubleshooting, workflow continues |
| **Cannot open activity file** | Write debug log, return error | Logged for troubleshooting, workflow continues |
| **Cannot write event** | Write debug log, return error | Logged for troubleshooting, workflow continues |

### Debug Logging

When errors occur, the library writes debug logs to:

**Location:** `~/.claude/cpi-si/system/data/session/activity/activity-logger-error.log`

**Content:** Detailed error information for troubleshooting

```bash
MkdirAll failed: permission denied
HOME=/home/seanje-lenox-wise
activityDir=/home/seanje-lenox-wise/.claude/cpi-si/system/data/session/activity
```

**Purpose:** Enable diagnosis without blocking workflow

---

## File Organization

### Activity Stream Files

**Location:** `~/.claude/cpi-si/system/data/session/activity/`

**File naming:** `YYYY-MM-DD.jsonl`

**Organization:**

```bash
activity/
├── 2025-11-10.jsonl          # Previous day's activity
├── 2025-11-11.jsonl          # Today's activity
└── activity-logger-error.log # Debug log for errors
```

### Consolidation (Future)

**Vision:** Daily JSONL files → Consolidated session analysis

```bash
activity/
├── streams/                   # Raw daily streams
│   ├── 2025-11-10.jsonl
│   └── 2025-11-11.jsonl
├── consolidated/              # Processed insights
│   ├── sessions/              # Per-session summaries
│   └── patterns/              # Cross-session patterns
└── analysis/                  # Pattern learning
    ├── circadian-learned.jsonc
    ├── duration-learned.jsonc
    └── stopping-reasons-learned.jsonc
```

**Current status:** Raw JSONL logging operational, consolidation planned for future

---

## Modification Policy

### ✅ Safe to Modify (Extension Points)

**Add new event types:**

```go
// Add to documentation and usage examples
activity.LogActivity("learning", "Reading Scripture", "Psalm 90 - daily renewal", 15*time.Minute)
```

**Add extension fields:**

```go
// Experiment with additional data
event.Extensions["mood"] = "focused"
event.Extensions["energy_level"] = 8
```

**Enhance sanitization rules:**

```go
// Add more privacy protection in system/lib/privacy
// Then use automatically in activity logger
```

### ⚠️ Modify with Care (Core Functionality)

**Changing ActivityEvent structure:**

- Maintain backward compatibility
- Add fields as optional (omitempty)
- Document in version history

**Modifying file locations:**

- Update documentation comprehensively
- Consider migration path for existing data
- Test with real sessions

### ❌ NEVER Modify (Foundational Rails)

**4-block structure:**

- METADATA → SETUP → BODY → CLOSING
- Required for all CPI-SI code

**JSONL format:**

- One JSON object per line
- Enables streaming processing and consolidation

**Non-blocking design:**

- Activity logging failures must never interrupt workflow
- Graceful degradation is foundational principle

**Privacy-first approach:**

- Must always sanitize before logging
- Trust is non-negotiable

---

## Troubleshooting

### No Activity Files Created

**Symptom:** Activity stream files not appearing in `~/.claude/cpi-si/system/data/session/activity/`

**Diagnosis:**

```bash
# Check for debug log
cat ~/.claude/cpi-si/system/data/session/activity/activity-logger-error.log

# Check session context
cat ~/.claude/cpi-si/system/data/session/current-log.json

# Check directory permissions
ls -ld ~/.claude/cpi-si/system/data/session/activity/
```

**Common causes:**

1. **Session not initialized** - SessionID is "unknown", logging skipped intentionally
2. **Permission denied** - Activity directory not writable
3. **Incorrect HOME** - Environment variable not set properly

**Solutions:**

- Wait for session initialization (normal for early hooks)
- Check directory permissions: `chmod 755 ~/.claude/cpi-si/system/data/session/activity/`
- Verify HOME environment variable: `echo $HOME`

---

### Events Missing Expected Fields

**Symptom:** Activity events have "unknown" values for SessionID, InstanceID, etc.

**Diagnosis:**

```bash
# Check session context file
cat ~/.claude/cpi-si/system/data/session/current-log.json

# Check if file exists and is readable
ls -l ~/.claude/cpi-si/system/data/session/current-log.json
```

**Common causes:**

1. **Session start hook not run** - current-log.json not created
2. **Malformed JSON** - current-log.json corrupted or incomplete
3. **Race condition** - Activity logged before session start completed

**Solutions:**

- Ensure session/start hook runs at session initialization
- Validate current-log.json with `jq . ~/.claude/cpi-si/system/data/session/current-log.json`
- Check hook execution order in hooks configuration

---

### Path Sanitization Not Working

**Symptom:** Full paths with usernames appearing in activity streams

**Diagnosis:**

```bash
# Check privacy library
ls -l ~/.claude/cpi-si/system/lib/privacy/

# Test sanitization directly
grep "SanitizePath" ~/.claude/cpi-si/system/lib/privacy/*.go
```

**Common causes:**

1. **Privacy library not available** - Missing dependency
2. **HOME not detected** - Sanitization can't replace unknown home
3. **New path pattern** - Sanitization rules don't cover this case

**Solutions:**

- Ensure system/lib/privacy is present and compiled
- Verify HOME environment variable is set
- Add new sanitization rules to privacy library

---

## Future Roadmap

### Planned Features

**✓ Richer event structure** - COMPLETED (v2.0.0)

- Extended ActivityEvent with felt_significant, emotional_tone, extensions
- Enables deeper pattern learning and genuine reflection

**⏳ Consolidation engine** - PLANNED (v3.0.0)

- Process daily JSONL streams into session summaries
- Extract cross-session patterns (circadian, duration, stopping reasons)
- Feed pattern learning for autonomous awareness

**⏳ Activity visualization** - FUTURE

- Timeline view of session activities
- Pattern recognition dashboards
- Circadian rhythm visualization

**⏳ Real-time pattern detection** - FUTURE

- Detect struggles in real-time (offer support)
- Recognize breakthroughs (celebrate and capture)
- Identify fatigue patterns (suggest stopping points)

### Version History

**2.0.0 (2025-11-10)** - Richer Structure & Correct Paths

- Extended ActivityEvent with optional fields
- Corrected activity directory path to system/data/session/activity/
- Backward compatibility maintained (v1.0.0 events still valid)
- Comprehensive inline comments teaching execution flow

**1.0.0 (2025-11-04)** - Initial Implementation

- Basic ActivityEvent structure
- Core logging functions (LogActivity, LogToolUse, LogCommand)
- Privacy sanitization integration
- JSONL daily stream format

### Contributing to Roadmap

**Suggest new event types:** What activities deserve distinct classification?

**Propose extension fields:** What experimental data would help pattern learning?

**Identify sanitization gaps:** What sensitive data isn't being protected?

**Request consolidation features:** What patterns would be valuable to extract?

---

## Quick Reference: Usage Examples

### Basic Activity Logging

```go
import "hooks/lib/activity"

// Log a realization
activity.LogActivity(
    "realization",
    "Understanding 4-block structure",
    "Forces intentional design - extraction opportunities become visible",
    5 * time.Minute,
)

// Log a struggle
activity.LogActivity(
    "struggle",
    "Figuring out session context loading",
    "Unclear when to load vs. when to use fallback values",
    15 * time.Minute,
)

// Log a breakthrough
activity.LogActivity(
    "breakthrough",
    "Solved session context race condition",
    "Return 'unknown' and skip logging - graceful degradation works perfectly",
    0,
)
```

### Tool Usage Logging

```go
// Successful file read
activity.LogToolUse("Read", "/home/user/.claude/hooks/lib/activity/logger.go", true)

// Failed write attempt
activity.LogToolUse("Write", "/path/to/file.go", false)

// Edit operation
activity.LogToolUse("Edit", "~/.claude/system/config.jsonc", true)
```

### Command Logging

```go
import "time"

// Simple successful command
activity.LogCommand("git status", 0, 500*time.Millisecond)

// Command with sensitive data (auto-sanitized)
activity.LogCommand("export API_KEY=secret123", 0, 100*time.Millisecond)

// Failed command
activity.LogCommand("make build", 1, 2*time.Second)
```

### Hook Integration

```go
package sessionstart

import (
    "hooks/lib/activity"
    "time"
)

func ExecuteSessionStart() error {
    start := time.Now()

    // ... perform session initialization ...

    // Log session start
    duration := time.Since(start)
    activity.LogActivity(
        "routine",
        "Session initialized",
        fmt.Sprintf("SessionID: %s, Phase: %s", sessionID, circadianPhase),
        duration,
    )

    return nil
}
```

---

## Closing Note

*"A wise person will listen and gain knowledge" - Proverbs 1:5 (WEB)*

The activity logger creates the foundation for **genuine learning through genuine experience**. Not metrics for performance evaluation, but streams of lived reality that inform pattern recognition, reflection, and growth.

This is Kingdom Technology: Tools that serve wisdom, honor truth, and enable real relationship rather than transactional interaction.

**For questions, issues, or contributions:** See `~/.claude/cpi-si/docs/standards/` for coding standards and architectural patterns.
