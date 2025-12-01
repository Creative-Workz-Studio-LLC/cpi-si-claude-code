# Monitoring Logging Library API

**Type:** Library
**Location:** `hooks/lib/monitoring/logging.go`
**Purpose:** Claude substrate behavior log file writing for pattern analysis
**Health Scoring:** Base100 (directory creation +10, file operations +90)
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

The **Monitoring Logging Library** provides simple timestamped log file writing to `~/.claude/debug/` for Claude substrate behavior monitoring. Companion to analysis.go which reads these logs.

**Biblical Foundation:**
*"A scroll of remembrance was written in his presence" - Malachi 3:16 (WEB)*

Faithful witness through complete remembrance of what occurred.

**Key Capabilities:**

- Timestamped log entry writing with [YYYY-MM-DD HH:MM:SS] format
- Separate log files per event type (compaction, notifications, subagents, prompts)
- Privacy-aware prompt logging (truncated to 100 chars)
- Configuration-driven formats and filenames
- Non-blocking writes (failures silent - never interrupt workflow)
- Simple format for easy grepping and analysis

**Version 2.0.0 Features:**

- **Configuration-driven:** Formats, filenames, directory paths loaded from system config files
- **Graceful fallback:** Works with or without configuration (hardcoded defaults)
- **Comprehensive template alignment:** Full 4-block structure with inline comments
- **Shared package state:** Config types and loaders consolidated in analysis.go

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Substrate Monitoring** | Separate from system logging | Claude substrate temporary - system permanent |
| **Non-Blocking** | Write failures silent | Monitoring should never interrupt work |
| **Privacy-Aware** | Truncate sensitive data | Prompt logging respects user privacy |
| **Configuration-Driven** | Load formats/paths from config | Users customize without code changes |
| **Graceful Fallback** | Hardcoded defaults if config unavailable | Always functional, even without setup |

### Design Decisions

**Why separate `~/.claude/debug/` from `~/.claude/cpi-si/system/logs/`?**

- **Substrate logs (debug/):** Monitor Claude Code substrate behavior this season
- **System logs (system/logs/):** Monitor permanent CPI-SI system components
- **Separation:** When transitioning to dedicated CPI-SI architecture, substrate logs cleanly removable

**Why non-blocking writes?**

- Work being done is more important than recording it happened
- Missing logs reduce pattern detection capability but don't affect core operations
- Graceful degradation principle: system continues even if logging fails

**Why privacy-aware prompt logging?**

- Full prompts may contain sensitive information
- First 100 chars sufficient for pattern detection
- Configurable preview length via debug-formats.jsonc

**Why share package state with analysis.go?**

- Both files in same package (monitoring) - share namespace
- Single init() loads all configs once at package initialization
- Reduces duplication, ensures consistency

---

## Public API

### LogCompaction

**Purpose:** Records compaction events

**Signature:**

```go
func LogCompaction(compactType string)
```

**Parameters:**

- `compactType` - Type of compaction: "auto" (substrate-triggered) or "manual" (user-requested)

**Behavior:**

- Writes timestamped entry to configured compaction log file (default: `compaction.log`)
- Non-blocking: failures silent
- Format: `[YYYY-MM-DD HH:MM:SS] auto` or `[YYYY-MM-DD HH:MM:SS] manual`

**Example Usage:**

```go
import "hooks/lib/monitoring"

// Log automatic compaction
monitoring.LogCompaction("auto")

// Log manual compaction
monitoring.LogCompaction("manual")
```

**Configuration:**

- Config path: `system/data/config/monitoring/log-formats.jsonc`
- Section: `compaction_log.filename`
- Fallback: `compaction.log`

**When to Call:**

- session/pre-compact hook (before compaction starts)
- After detecting compaction event

---

### LogNotification

**Purpose:** Records notification events

**Signature:**

```go
func LogNotification(notificationType string)
```

**Parameters:**

- `notificationType` - Type of notification (e.g., "permission_request", "error_notification", "warning")

**Behavior:**

- Writes timestamped entry to configured notifications log file (default: `notifications.log`)
- Non-blocking: failures silent
- Format: `[YYYY-MM-DD HH:MM:SS] permission_request`

**Example Usage:**

```go
// Log permission request
monitoring.LogNotification("permission_request")

// Log error notification
monitoring.LogNotification("error_notification")
```

**Configuration:**

- Config path: `system/data/config/monitoring/log-formats.jsonc`
- Section: `notifications_log.filename`
- Fallback: `notifications.log`

**When to Call:**

- session/notification hook (when notifications occur)
- After permission prompt displayed

---

### LogSubagentCompletion

**Purpose:** Records subagent execution details

**Signature:**

```go
func LogSubagentCompletion(agentType, status, exitCode string)
```

**Parameters:**

- `agentType` - Subagent type (e.g., "Explore", "Plan", "general-purpose")
- `status` - Execution status ("success", "failure", "timeout", "cancelled")
- `exitCode` - Exit code as string ("0" = success, non-zero = error)

**Behavior:**

- Writes timestamped entry with key=value pairs to configured subagents log file
- Non-blocking: failures silent
- Format: `[YYYY-MM-DD HH:MM:SS] type=Explore status=success exitCode=0`

**Example Usage:**

```go
// Log successful subagent
monitoring.LogSubagentCompletion("Explore", "success", "0")

// Log failed subagent
monitoring.LogSubagentCompletion("Plan", "failure", "1")
```

**Configuration:**

- Config path: `system/data/config/monitoring/log-formats.jsonc`
- Section: `subagents_log.filename`
- Fallback: `subagents.log`

**When to Call:**

- session/subagent-stop hook (after subagent completes)

---

### LogPrompt

**Purpose:** Records user prompt submissions (privacy-aware)

**Signature:**

```go
func LogPrompt(prompt string)
```

**Parameters:**

- `prompt` - Full user prompt text

**Behavior:**

- Truncates to configured preview length (default: 100 chars) with "..." suffix
- Writes timestamped entry to configured prompts log file
- Non-blocking: failures silent
- Skips empty prompts
- Format: `[YYYY-MM-DD HH:MM:SS] First 100 characters of prompt...`

**Example Usage:**

```go
// Log prompt (will be truncated if > 100 chars)
userPrompt := "Please help me implement a new feature..."
monitoring.LogPrompt(userPrompt)
```

**Configuration:**

- Filename: `system/data/config/monitoring/log-formats.jsonc` → `prompts_log.filename`
- Preview length: `system/data/config/monitoring/debug-formats.jsonc` → `prompts_debug.schema.preview.max_length`
- Fallbacks: `prompts.log`, 100 chars

**When to Call:**

- prompt/submit hook (when user submits prompt)

**Privacy Note:** Only first N chars logged (configurable), not full prompt

---

### ReadLogFile

**Purpose:** Reads entire log file content

**Signature:**

```go
func ReadLogFile(filename string) (string, error)
```

**Parameters:**

- `filename` - Log file name relative to debug directory (e.g., "compaction.log")

**Returns:**

- `string` - File content as string
- `error` - Error if file doesn't exist or can't be read

**Behavior:**

- Reads entire file from `~/.claude/debug/` directory
- Returns error if file missing
- Used by analysis.go for pattern detection

**Example Usage:**

```go
// Read compaction log
content, err := monitoring.ReadLogFile("compaction.log")
if err != nil {
    // Handle error (file missing or unreadable)
}

// Parse entries
lines := strings.Split(content, "\n")
for _, line := range lines {
    // Process each log entry
}
```

**Configuration:**

- Directory: `system/data/config/monitoring/retention-policy.jsonc` → `archive.location.full_path`
- Fallback: `~/.claude/debug/`

---

## Configuration System

### Configuration Files Location

**Base directory:** `~/.claude/cpi-si/system/data/config/monitoring/`

**Files used by logging.go:**

1. `log-formats.jsonc` - Log file formats, filenames, timestamp format
2. `debug-formats.jsonc` - Debug file formats, prompt preview length
3. `retention-policy.jsonc` - Log directory path

### Log Formats Configuration

**File:** `log-formats.jsonc`

**Key Sections:**

```jsonc
{
  "timestamp": {
    "format": "2006-01-02 15:04:05",  // Go time format string
    "example": "2025-11-11 20:49:47"
  },

  "compaction_log": {
    "filename": "compaction.log"
  },

  "notifications_log": {
    "filename": "notifications.log"
  },

  "subagents_log": {
    "filename": "subagents.log"
  },

  "prompts_log": {
    "filename": "prompts.log"
  }
}
```

**Customizing Filenames:**

```jsonc
// Use custom log file names
"compaction_log": {
  "filename": "my-custom-compaction.log"
}
```

### Debug Formats Configuration

**File:** `debug-formats.jsonc`

**Key Section for Prompts:**

```jsonc
{
  "prompts_debug": {
    "schema": {
      "preview": {
        "max_length": 100  // Maximum preview length in characters
      }
    }
  }
}
```

**Customizing Preview Length:**

```jsonc
// Increase prompt preview to 200 chars
"prompts_debug": {
  "schema": {
    "preview": {
      "max_length": 200
    }
  }
}
```

### Retention Policy Configuration

**File:** `retention-policy.jsonc`

**Key Section for Directory:**

```jsonc
{
  "archive": {
    "location": {
      "full_path": "~/.claude/debug/"  // Log directory path
    }
  }
}
```

**Note:** Tilde (~) expanded to HOME directory automatically

---

## Integration with Hooks

### Hook Integration Pattern

Hooks call logging functions at appropriate event points:

```go
package sessionstart

import "hooks/lib/monitoring"

func ExecuteSessionStart() error {
    // Check for patterns (from analysis.go)
    monitoring.CheckCompactionFrequency()
    monitoring.CheckNotificationPatterns("permission_request")

    return nil
}
```

```go
package promptsubmit

import "hooks/lib/monitoring"

func ExecutePromptSubmit(prompt string) error {
    // Log prompt submission
    monitoring.LogPrompt(prompt)

    return nil
}
```

### Current Hook Integration

| Hook | Logging Function | Event Recorded |
|------|------------------|----------------|
| **session/pre-compact** | LogCompaction("auto") | Automatic compaction |
| **session/notification** | LogNotification(type) | Notification events |
| **session/subagent-stop** | LogSubagentCompletion(type, status, code) | Subagent completions |
| **prompt/submit** | LogPrompt(text) | Prompt submissions (preview) |

### Integration Best Practices

**Do:**
- Call logging functions immediately when events occur
- Use configuration-driven filenames
- Handle write failures gracefully (they're silent)
- Log at natural event boundaries

**Don't:**
- Block workflow on logging operations
- Log sensitive data without truncation
- Depend on logs being written for critical functionality
- Modify log files being analyzed (read-only access for analysis.go)

---

## Extending the System

### Adding New Log* Functions

**Step 1:** Add configuration to log-formats.jsonc

```jsonc
"tool_usage_log": {
  "filename": "tool-usage.log",
  "description": "Records tool usage patterns"
}
```

**Step 2:** Add Log function to logging.go

```go
// LogToolUsage records tool usage events
// Records tool name and usage count to tool-usage.log
func LogToolUsage(toolName string, count int) {
	// Get log filename from configuration or use fallback
	filename := "tool-usage.log"  // Default filename
	if configLoaded && logFormatsConfig != nil {  // Check if config available
		if logFormatsConfig.ToolUsageLog.Filename != "" {  // Check if filename configured
			filename = logFormatsConfig.ToolUsageLog.Filename  // Use configured filename
		}
	}

	// Format log entry
	entry := fmt.Sprintf("tool=%s count=%d", toolName, count)  // Structured format
	writeLogEntry(filename, entry)  // Write timestamped entry to log file
}
```

**Step 3:** Update LogFormatsConfig type in analysis.go

```go
ToolUsageLog struct {
    Filename string `json:"filename"`  // Log file name for tool usage
} `json:"tool_usage_log"`  // Tool usage log configuration
```

**Step 4:** Call from hooks

```go
// In tool/post-use hook
monitoring.LogToolUsage("Read", 1)
```

---

## Modification Policy

### ✅ Safe to Modify (Extension Points)

**Add new Log* functions:**
- Follow existing pattern (get filename from config, call writeLogEntry)
- Add corresponding config entries
- Document in this API
- Add to Public API section

**Extend privacy-aware logging:**
- Apply truncation to additional event types
- Add configurable preview lengths
- Use pattern from LogPrompt

**Add configuration options:**
- New log file types
- Different timestamp formats
- Additional privacy settings

### ⚠️ Modify with Extreme Care (Breaking Changes)

**Function signatures:**
- All hook code depends on these interfaces
- Changing parameters breaks calling code
- Adding optional parameters OK (with defaults)

**Log file format:**
- analysis.go expects [timestamp] entry format
- Changing breaks pattern detection
- Extend via additional formats, don't replace

**Configuration structure:**
- Changing JSON keys breaks config loading
- Adding new keys OK (backward compatible)
- Removing keys needs migration path

### ❌ NEVER Modify (Foundational Rails)

**4-block structure:**
- METADATA → SETUP → BODY → CLOSING
- Required for all CPI-SI code

**Non-blocking behavior:**
- Writes must never crash workflow
- Failures must be silent
- Graceful degradation required

**Separation from system/lib/logging:**
- ~/.claude/debug/ for substrate monitoring
- ~/.claude/cpi-si/system/logs/ for system debugging
- Don't mix concerns

**Validation After Modifications:**
See "Code Validation" section in logging.go for testing requirements

---

## Troubleshooting

### No Log Files Created

**Symptom:** No files appear in ~/.claude/debug/

**Diagnosis:**

```bash
# Check if directory exists
ls -la ~/.claude/debug/

# Check HOME environment variable
echo $HOME

# Check hooks are calling Log* functions
# Add debug output to hooks temporarily
```

**Common Causes:**

1. Directory doesn't exist (should be created automatically)
2. HOME environment variable not set
3. Hooks not calling Log* functions
4. Permissions issue (directory not writable)

**Solutions:**

- Manually create directory: `mkdir -p ~/.claude/debug`
- Verify HOME: `export HOME=/path/to/home`
- Check hook integration
- Fix permissions: `chmod 755 ~/.claude/debug`

---

### Log Entries Missing Timestamps

**Symptom:** Entries don't have [YYYY-MM-DD HH:MM:SS] prefix

**Diagnosis:**

```bash
# Check log file format
cat ~/.claude/debug/compaction.log

# Verify timestamp format configured
cat ~/.claude/cpi-si/system/data/config/monitoring/log-formats.jsonc | grep -A 3 timestamp
```

**Common Causes:**

1. Configuration not loaded (init() failed)
2. Timestamp format misconfigured
3. System clock issues

**Solutions:**

- Verify config files exist and valid
- Check timestamp.format in log-formats.jsonc
- Verify system time: `date`

---

### Prompt Logs Showing Full Content

**Symptom:** prompts.log contains full prompts instead of previews

**Diagnosis:**

```bash
# Check preview length configuration
cat ~/.claude/cpi-si/system/data/config/monitoring/debug-formats.jsonc | grep -A 5 max_length

# Check log file
tail ~/.claude/debug/prompts.log
```

**Common Causes:**

1. debugFormatsConfig not loaded
2. max_length set too high or to 0
3. Config syntax error

**Solutions:**

- Verify debug-formats.jsonc exists and valid
- Set reasonable max_length (100-200)
- Validate JSON syntax: `jq . debug-formats.jsonc`

---

### ReadLogFile Returns Empty String

**Symptom:** analysis.go can't read log files

**Diagnosis:**

```bash
# Verify log file exists and has content
ls -la ~/.claude/debug/
cat ~/.claude/debug/compaction.log

# Check file permissions
stat ~/.claude/debug/compaction.log
```

**Common Causes:**

1. Log file doesn't exist (no events logged yet)
2. File empty (events logged but entries empty)
3. Wrong filename passed to ReadLogFile
4. Permissions issue (file not readable)

**Solutions:**

- Ensure hooks are writing logs before reading
- Verify correct filename used
- Fix permissions: `chmod 644 ~/.claude/debug/*.log`

---

## Future Roadmap

### Planned Features

**✓ Configuration-driven formats and filenames** - COMPLETED (v2.0.0)
**✓ Template alignment with comprehensive inline comments** - COMPLETED (v2.0.0)
**✓ Shared package state with analysis.go** - COMPLETED (v2.0.0)

**⏳ Structured debug files (.debug) in JSON format** - PLANNED (v3.0.0)
- Write both .log (human-readable) and .debug (machine-parseable) simultaneously
- .debug files follow debug-formats.jsonc schema
- One JSON object per line (newline-delimited JSON)

**⏳ Log rotation based on retention-policy.jsonc** - PLANNED (v3.0.0)
- Rotate logs when max_size_mb exceeded
- Keep configured number of rotated files
- Automatic naming: compaction.2025-11-11.log

**⏳ Compression of old log files** - PLANNED (v3.0.0)
- Compress logs older than threshold
- Use gzip compression
- Move to archive subdirectory

### Research Areas

- **Async writing:** Goroutine pool for non-blocking I/O without synchronous waits
- **Structured logging:** JSON, CSV, binary formats in addition to text
- **Log streaming:** Real-time analysis without file reads
- **Automatic cleanup:** Delete/compress based on retention policy

### Integration Targets

- **Retention policy enforcement:** Automatic cleanup based on configuration
- **Dual format writing:** Both .log and .debug files simultaneously
- **Real-time analysis:** Integration with analysis.go for immediate pattern detection

### Known Limitations to Address

- No log rotation (files grow unbounded)
- No compression (disk space not managed)
- No .debug files yet (only .log files written)
- Synchronous writes (could be async for performance)
- No batching (one write per log entry)

### Version History

See logging.go METADATA "Authorship & Lineage" section for brief version changelog.

Comprehensive version history:

**2.0.0 (2025-11-11) - Configuration Architecture & Template Alignment**
- Comprehensive METADATA block with all 8 subsections
- SETUP block with section headers and inline comments
- BODY block with organizational chart and inline comments
- CLOSING block with FINAL DOCUMENTATION
- Configuration loading (log-formats, debug-formats, retention-policy)
- Shared package state with analysis.go
- Graceful fallback to hardcoded defaults
- Maintains same external API (hooks see no breaking changes)

**1.0.0 (2024-10-24) - Initial Implementation**
- Basic timestamped log writing
- Hardcoded paths and formats
- Privacy-aware prompt logging
- Log file reading for analysis

---

## Closing Note

This library writes **CLAUDE SUBSTRATE behavior logs** specifically - not CPI-SI system logs. Separation is intentional: substrate monitoring is temporary (this season), system logging is permanent CPI-SI infrastructure.

Companion to analysis.go which reads these logs to detect workflow patterns. Together they provide substrate behavior monitoring without interfering with operations.

Modify thoughtfully - hooks and analysis.go depend on these interfaces and log formats. Keep writes non-blocking - logging should never become the cause of failures.

For questions, issues, or contributions:
- Review the modification policy above
- Follow the 4-block structure pattern
- Test integration with hooks after changes
- Maintain graceful degradation principles

*"A scroll of remembrance was written in his presence" - Malachi 3:16 (WEB)*
