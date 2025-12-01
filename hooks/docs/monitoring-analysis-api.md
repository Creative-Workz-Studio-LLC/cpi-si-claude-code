# Monitoring Analysis Library API

**Type:** Library
**Location:** `hooks/lib/monitoring/analysis.go`
**Purpose:** Claude substrate pattern detection for workflow optimization
**Health Scoring:** Delegated to system/lib/display (Base100)
**Status:** ✅ Operational (Version 2.0.0)

---

## Table of Contents

1. [Overview](#overview)
2. [Philosophy & Design Rationale](#philosophy--design-rationale)
3. [Public API](#public-api)
4. [Configuration System](#configuration-system)
5. [Integration with Hooks](#integration-with-hooks)
6. [Extending the System](#extending-the-system)
7. [Modification Policy](#modification-policy)
8. [Troubleshooting](#troubleshooting)
9. [Future Roadmap](#future-roadmap)

---

## Overview

The **Monitoring Analysis Library** provides substrate behavior pattern detection for Claude Code. It analyzes log files in `~/.claude/debug/` to detect inefficiency patterns (excessive compaction, permission issues) and displays non-intrusive warnings.

**Biblical Foundation:**
*"The plans of the diligent lead surely to abundance" - Proverbs 21:5 (WEB)*

Monitor patterns to work wisely, not frantically. Proactive detection prevents substrate overload.

**Key Capabilities:**

- Time-windowed event counting from log files
- Compaction frequency detection (warns if auto-compacting too often)
- Notification pattern analysis (detects permission configuration issues)
- Configurable thresholds and time windows
- Non-blocking operation (warnings never interrupt workflow)

**Version 2.0.0 Features:**

- **Configuration-driven:** Thresholds, time windows, and messages loaded from system config files
- **Graceful fallback:** Works with or without configuration (hardcoded defaults)
- **System library integration:** Delegates warning display to system/lib/display
- **Extensible:** Users can customize without modifying code

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Substrate Monitoring** | Separate from system debugging logs | Claude substrate temporary - system permanent |
| **Non-Intrusive** | Advisory warnings only, never block workflow | Help optimize without disrupting work |
| **Configuration-Driven** | Load patterns/thresholds from config files | Users customize without code changes |
| **Graceful Fallback** | Hardcoded defaults if config unavailable | Always functional, even without setup |
| **Pattern-Based** | Detect trends, not individual events | Signal actual issues, not noise |

### Design Decisions

**Why separate `~/.claude/debug/` from `~/.claude/system/logs/`?**

- **Substrate logs (debug/):** Monitor Claude Code substrate behavior this season
- **System logs (system/logs/):** Monitor permanent CPI-SI system components
- **Separation:** When transitioning to dedicated CPI-SI architecture, substrate logs cleanly removable

**Why configuration-driven thresholds?**

- Different workflows have different patterns (fast iteration vs deep work)
- Users can tune sensitivity (fewer warnings or more aggressive detection)
- Supports future machine learning (learn optimal thresholds per user)

**Why graceful fallback?**

- Works out-of-the-box without configuration
- No breaking changes for existing code
- Progressive enhancement: better with config, functional without

**Why delegate to system/lib/display?**

- Don't duplicate formatting logic (extract and orchestrate)
- System library provides ANSI colors, icons, panic recovery
- Architectural integrity: hooks → monitoring → display

---

## Public API

### CheckCompactionFrequency

**Purpose:** Analyzes compaction patterns and warns if excessive

**Signature:**

```go
func CheckCompactionFrequency()
```

**Behavior:**

- Counts auto-compaction events in configured time window (default: 1 hour)
- Compares count to configured threshold (default: 3)
- Displays warning if threshold exceeded
- Uses configuration if available, falls back to hardcoded defaults

**Example Usage:**

```go
import "hooks/lib/monitoring"

// Check for excessive auto-compaction (call at session start or hourly)
monitoring.CheckCompactionFrequency()
// Output (if threshold exceeded): ⚠ Frequent auto-compaction (5 in last hour) - consider being more concise
```

**Configuration:**

- Config path: `system/data/config/monitoring/analysis-patterns.jsonc`
- Section: `compaction`
- Configurable: search_term, time_window_hours, thresholds.warning, messages.warning
- Fallback: "compaction.log", "auto", 1 hour, threshold 3

**When to Call:**

- session/start hook (check at beginning of session)
- Periodic checks (hourly via cron or scheduled task)
- prompt/submit hook (after significant user interaction)

**Interpretation:**

- Count <= 3: Normal behavior, no warning
- Count 4-10: Moderate compaction (verbose output or many tool calls)
- Count > 10: Excessive compaction (workflow inefficiency, consider optimization)

---

### CheckNotificationPatterns

**Purpose:** Analyzes notification frequency by type

**Signature:**

```go
func CheckNotificationPatterns(notificationType string)
```

**Parameters:**

- `notificationType` - Type of notification to analyze ("permission_request", etc.)

**Behavior:**

- Currently supports "permission_request" type
- Counts notifications in configured time window (default: 1 hour)
- Compares count to configured threshold (default: 10)
- Displays warning if threshold exceeded
- Other notification types ignored (extend via configuration)

**Example Usage:**

```go
// Check for excessive permission requests
monitoring.CheckNotificationPatterns("permission_request")
// Output (if threshold exceeded): ⚠ Many permission requests (15 in last hour) - check permissions configuration
```

**Configuration:**

- Config path: `system/data/config/monitoring/analysis-patterns.jsonc`
- Section: `notifications.permission_request`
- Configurable: search_term, time_window_hours, thresholds.warning, messages.warning
- Fallback: "notifications.log", "permission_request", 1 hour, threshold 10

**When to Call:**

- session/start hook (check at session start)
- tool/post-use hook (after operations that might trigger permissions)
- Periodic checks

**Interpretation:**

- Count <= 10: Normal permission requests
- Count 11-25: Elevated (may indicate permissions misconfiguration)
- Count > 25: Excessive (definitely misconfigured, needs review)

---

### CountRecentEvents

**Purpose:** Counts occurrences of a substring in log file within time window

**Signature:**

```go
func CountRecentEvents(logFilename, searchTerm string, hoursBack int) int
```

**Parameters:**

- `logFilename` - Log file name relative to `~/.claude/debug/` (e.g., "compaction.log")
- `searchTerm` - Substring to search for in log entries (case-sensitive)
- `hoursBack` - Time window in hours (events older than this ignored)

**Returns:**

- `int` - Count of matching events within time window

**Behavior:**

- Reads log file from `~/.claude/debug/` directory
- Parses timestamps ([YYYY-MM-DD HH:MM:SS] format)
- Filters to events within time window
- Counts entries containing search term
- Returns 0 if log file missing (graceful degradation)
- Skips malformed timestamps (continues with remaining lines)

**Example Usage:**

```go
// Count manual compactions in last 3 hours
count := monitoring.CountRecentEvents("compaction.log", "manual", 3)
if count > 5 {
    // Handle elevated manual compaction rate
}

// Count subagent failures in last 24 hours
failureCount := monitoring.CountRecentEvents("subagents.log", "failure", 24)
if failureCount > 10 {
    // Investigate subagent failure patterns
}
```

**Graceful Degradation:**

- Missing log file: Returns 0 (no error)
- Malformed timestamps: Skips invalid lines, processes valid ones
- Empty search term: Counts all entries in time window
- Zero hoursBack: Only counts events in current second

---

## Configuration System

### Configuration Files Location

**Base directory:** `~/.claude/cpi-si/system/data/config/monitoring/`

**Files:**

1. `analysis-patterns.jsonc` - Thresholds, time windows, search patterns
2. `log-formats.jsonc` - Structured .log file format schema
3. `debug-formats.jsonc` - Structured .debug file format schema
4. `retention-policy.jsonc` - Log rotation and retention rules

### Analysis Patterns Configuration

**File:** `analysis-patterns.jsonc`

**Structure:**

```jsonc
{
  "log_files": {
    "compaction": "compaction.log",
    "notifications": "notifications.log",
    "subagents": "subagents.log",
    "prompts": "prompts.log"
  },

  "time_windows": {
    "default": 1,
    "short": 0.5,
    "medium": 3,
    "long": 24
  },

  "compaction": {
    "search_term": "auto",
    "time_window_hours": 1,
    "thresholds": {
      "warning": 3,
      "critical": 10
    },
    "messages": {
      "warning": "Frequent auto-compaction ({count} in last hour) - consider being more concise",
      "critical": "Excessive auto-compaction ({count} in last hour) - workflow becoming inefficient"
    }
  },

  "notifications": {
    "permission_request": {
      "search_term": "permission_request",
      "time_window_hours": 1,
      "thresholds": {
        "warning": 10,
        "critical": 25
      },
      "messages": {
        "warning": "Many permission requests ({count} in last hour) - check permissions configuration",
        "critical": "Excessive permission requests ({count} in last hour) - permissions misconfigured"
      }
    }
  }
}
```

**Customizing Thresholds:**

```jsonc
// Increase compaction threshold for fast-paced workflows
"compaction": {
  "thresholds": {
    "warning": 5,    // Increased from 3
    "critical": 15   // Increased from 10
  }
}

// Decrease permission threshold for strict monitoring
"notifications": {
  "permission_request": {
    "thresholds": {
      "warning": 5,   // Decreased from 10
      "critical": 15  // Decreased from 25
    }
  }
}
```

**Message Templates:**

Messages support variable substitution:

- `{count}` - Replaced with actual event count
- Future: `{time_window}`, `{threshold}`, `{percentage}`

Example:
```jsonc
"messages": {
  "warning": "Detected {count} events (threshold: {threshold})"
}
```

---

## Integration with Hooks

### Hook Integration Pattern

Hooks call analysis functions at appropriate points:

```go
package sessionstart

import "hooks/lib/monitoring"

func ExecuteSessionStart() error {
    // Monitor substrate patterns at session start
    monitoring.CheckCompactionFrequency()
    monitoring.CheckNotificationPatterns("permission_request")

    // Custom analysis based on session context
    if isLongSession() {
        // Check for sustained high activity patterns
        subagentCount := monitoring.CountRecentEvents("subagents.log", "Explore", 6)
        if subagentCount > 20 {
            // Many Explore subagents - might indicate searching inefficiency
        }
    }

    return nil
}
```

### Current Hook Integration

| Hook | Analysis Function | Trigger Condition |
|------|-------------------|-------------------|
| **session/start** | CheckCompactionFrequency() | At session beginning |
| **session/start** | CheckNotificationPatterns() | At session beginning |
| **prompt/submit** | (Future) CheckPromptFrequency() | After user prompts |
| **tool/post-use** | (Future) Custom event counting | After specific tools |

### Integration Best Practices

**Do:**
- Call checks at natural breakpoints (session start, hourly)
- Use CountRecentEvents for custom pattern detection
- Handle warnings gracefully (they're advisory)
- Log when thresholds exceeded (for trend analysis)

**Don't:**
- Call checks in tight loops (performance overhead)
- Block workflow based on warnings (advisory only)
- Modify log files being analyzed (read-only access)
- Depend on warnings for critical functionality

---

## Extending the System

### Adding New Check Functions

**Step 1:** Add configuration to analysis-patterns.jsonc

```jsonc
"subagents": {
  "failure_rate": {
    "time_window_hours": 3,
    "thresholds": {
      "warning": 0.3,
      "critical": 0.5
    },
    "messages": {
      "warning": "Subagent failure rate elevated ({percentage}% in last 3 hours)",
      "critical": "Subagent failure rate high ({percentage}% in last 3 hours)"
    }
  }
}
```

**Step 2:** Update AnalysisPatternsConfig type

```go
Subagents struct {
    FailureRate struct {
        TimeWindowHours int `json:"time_window_hours"`
        Thresholds struct {
            Warning  float64 `json:"warning"`
            Critical float64 `json:"critical"`
        } `json:"thresholds"`
        Messages struct {
            Warning  string `json:"warning"`
            Critical string `json:"critical"`
        } `json:"messages"`
    } `json:"failure_rate"`
} `json:"subagents"`
```

**Step 3:** Add check function to analysis.go

```go
// CheckSubagentFailureRate analyzes subagent failure patterns
func CheckSubagentFailureRate() {
    // Load config or use defaults
    timeWindow := 3
    threshold := 0.3
    messageTemplate := "Subagent failure rate elevated ({percentage}% in last 3 hours)"

    if configLoaded && analysisConfig != nil {
        timeWindow = analysisConfig.Subagents.FailureRate.TimeWindowHours
        threshold = analysisConfig.Subagents.FailureRate.Thresholds.Warning
        messageTemplate = analysisConfig.Subagents.FailureRate.Messages.Warning
    }

    // Count total and failed subagents
    total := CountRecentEvents("subagents.log", "type=", timeWindow)
    failures := CountRecentEvents("subagents.log", "status=failure", timeWindow)

    if total > 0 {
        failureRate := float64(failures) / float64(total)
        if failureRate > threshold {
            percentage := int(failureRate * 100)
            message := strings.Replace(messageTemplate, "{percentage}", fmt.Sprintf("%d", percentage), -1)
            fmt.Println(display.Warning(message))
        }
    }
}
```

**Step 4:** Call from hooks

```go
// In session/start or periodic check
monitoring.CheckSubagentFailureRate()
```

---

## Modification Policy

### ✅ Safe to Modify (Extension Points)

**Add new check functions:**
- Follow existing pattern (load config, count events, check threshold, display warning)
- Add corresponding config entries
- Document in this API
- Add to Public API section

**Extend notification types:**
- Add new notification type configs
- Update CheckNotificationPatterns with switch statement
- Or create dedicated check function per type

**Add utility functions:**
- Pattern detection helpers
- Log parsing utilities
- Time window calculations
- Follow inline comment style

### ⚠️ Modify with Extreme Care (Breaking Changes)

**Function signatures:**
- All hook code depends on these interfaces
- Changing parameters breaks calling code
- Adding optional parameters OK (with defaults)

**Configuration file structure:**
- Changing JSON keys breaks config loading
- Adding new keys OK (backward compatible)
- Removing keys needs migration path

**Log file format expectations:**
- Timestamp format: [YYYY-MM-DD HH:MM:SS]
- Changing breaks CountRecentEvents parsing
- Extend via additional formats, don't replace

### ❌ NEVER Modify (Foundational Rails)

**4-block structure:**
- METADATA → SETUP → BODY → CLOSING
- Required for all CPI-SI code

**Substrate/system separation:**
- ~/.claude/debug/ for substrate monitoring
- ~/.claude/system/logs/ for system debugging
- Don't mix concerns

**Non-blocking behavior:**
- Analysis must never crash workflow
- Warnings are advisory only
- Graceful degradation required

---

## Troubleshooting

### No Warnings Appear

**Symptom:** No output even when thresholds should be exceeded

**Diagnosis:**

```bash
# Check if log files exist
ls -la ~/.claude/debug/

# Check log file content and format
cat ~/.claude/debug/compaction.log

# Verify timestamps parse correctly
head -5 ~/.claude/debug/compaction.log
# Should see: [2025-11-11 20:49:47] entry_data
```

**Common Causes:**

1. Log files missing (not being written)
2. Timestamp format doesn't match [YYYY-MM-DD HH:MM:SS]
3. Search term doesn't match log entries (case-sensitive)
4. Hooks not calling check functions

**Solutions:**

- Review logging.go implementation
- Verify log file format
- Test with known good log entries
- Check hook integration

---

### CountRecentEvents Always Returns 0

**Symptom:** Function returns 0 even when events exist

**Diagnosis:**

```bash
# Verify log file exists and readable
ls -la ~/.claude/debug/compaction.log

# Check content matches search term
grep "auto" ~/.claude/debug/compaction.log

# Verify recent entries (within time window)
tail -20 ~/.claude/debug/compaction.log
```

**Common Causes:**

1. Log file missing or empty
2. Search term doesn't match (case-sensitive: "auto" ≠ "Auto")
3. Time window too small (all events older than cutoff)
4. Timestamp parsing failures (all lines skipped)

**Solutions:**

- Use exact search term from logs
- Increase time window for testing
- Verify timestamp format matches code expectations
- Add debug logging to CountRecentEvents

---

### Configuration Not Loading

**Symptom:** Always uses fallback values, config changes ignored

**Diagnosis:**

```bash
# Check config file exists
ls -la ~/.claude/cpi-si/system/data/config/monitoring/analysis-patterns.jsonc

# Validate JSONC syntax
grep -v "^[[:space:]]*\/\/" analysis-patterns.jsonc | jq .

# Check permissions
stat ~/.claude/cpi-si/system/data/config/monitoring/analysis-patterns.jsonc
```

**Common Causes:**

1. Config file missing or wrong path
2. JSONC syntax errors (malformed JSON after comment stripping)
3. HOME environment variable not set
4. Config directory permissions

**Solutions:**

- Verify file path: `~/.claude/cpi-si/system/data/config/monitoring/`
- Validate JSON syntax (strip // comments first)
- Check HOME: `echo $HOME`
- Fix permissions: `chmod 644 analysis-patterns.jsonc`

---

## Future Roadmap

### Planned Features

**✓ Configuration-driven thresholds** - COMPLETED (v2.0.0)
**✓ Display library integration** - COMPLETED (v2.0.0)

**⏳ Subagent pattern analysis** - PLANNED (v3.0.0)
- Monitor subagent failure rates
- Detect subagent timeout patterns
- Analyze subagent type distribution

**⏳ Prompt frequency detection** - PLANNED (v3.0.0)
- Detect rapid prompt submission (possible loops)
- Monitor prompt complexity trends
- Identify repetitive prompt patterns

**⏳ Trend analysis** - PLANNED (v3.0.0)
- Not just current window, but patterns over time
- Detect degrading patterns (increasing frequency)
- Historical comparison (this hour vs average)

**⏳ Critical threshold support** - PLANNED (v3.0.0)
- Config defines warning and critical thresholds
- Display critical warnings differently (red vs yellow)
- Escalating severity based on count

### Research Areas

- **Statistical anomaly detection** - ML-based threshold adjustment
- **Cross-log correlation** - Analyze relationships between log types
- **Predictive warnings** - Pattern indicates future issue before threshold
- **Adaptive thresholds** - Learn optimal values per user/workflow

### Integration Targets

- **Session time awareness** - Different thresholds for work hours vs downtime
- **Activity logging** - Track which warnings correlated with outcomes
- **Pattern learning** - Discover what patterns actually indicate problems
- **User preferences** - Per-user sensitivity tuning

### Known Limitations

- Hardcoded fallback values (config preferred but fallback always present)
- Simple substring matching (could use regex for complex patterns)
- No trend analysis (only current time window snapshot)
- Limited notification type support (only permission_request currently)
- No cross-log correlation (each log analyzed independently)

### Version History

**2.0.0 (2025-11-11)** - Configuration-Driven Architecture

- Created configuration system (analysis-patterns.jsonc, log-formats.jsonc, debug-formats.jsonc, retention-policy.jsonc)
- Implemented configuration loading with JSONC comment stripping
- Updated check functions to use configured thresholds and messages
- Graceful fallback to hardcoded defaults
- Delegates warning output to system/lib/display
- Comprehensive template alignment (METADATA, SETUP, BODY, CLOSING)
- Added inline comments to all functions

**1.0.0 (2024-10-24)** - Initial Implementation

- Basic time-windowed event counting
- Compaction frequency detection (hardcoded threshold: 3)
- Permission request pattern detection (hardcoded threshold: 10)
- Manual fmt.Printf for warnings
- Minimal documentation

---

## Closing Note

This library monitors **CLAUDE SUBSTRATE behavior** specifically - not CPI-SI system debugging. The separation is intentional: substrate monitoring is temporary (this season), system debugging is permanent. When transitioning to dedicated CPI-SI architecture, substrate monitoring can be removed cleanly.

Pattern detection helps optimize workflow without being intrusive. Warnings are advisory - they inform, never interrupt.

*"The plans of the diligent lead surely to abundance" - Proverbs 21:5 (WEB)*

**For questions, issues, or contributions:**
- Review the modification policy above
- Follow the 4-block structure pattern
- Test integration with hooks after changes
- Maintain graceful degradation principles
