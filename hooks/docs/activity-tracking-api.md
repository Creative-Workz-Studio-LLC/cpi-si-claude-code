# Activity Tracking Library API

**Type:** Library
**Location:** `hooks/lib/session/activity.go`
**Purpose:** Recent file modification tracking for session context awareness
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

The **Activity Tracking Library** provides workspace activity awareness by detecting recently modified files using configurable time windows and exclusion patterns. It enables:

- **Session context awareness** - Understanding what work has been done recently
- **Non-intrusive display** - Shows relevant activity without overwhelming output
- **Configurable behavior** - Time windows, exclusions, and thresholds customizable
- **Graceful fallback** - Works with hardcoded defaults if configuration missing

**Biblical Foundation:**
*"Whatever your hand finds to do, do it with all your might" - Ecclesiastes 9:10 (WEB)*

Faithful work leaves traces - recent activity reveals current focus and stewardship.

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Config-Driven** | All behavior controlled via activity-tracking.jsonc | Users customize without code changes |
| **Non-Blocking** | Failures never disrupt session operations | Context awareness serves workflow, doesn't block it |
| **Graceful Fallback** | Hardcoded defaults when config missing | Always works, even in degraded state |
| **Smart Thresholds** | Detailed/summary/silent modes based on file count | Balance information vs noise |
| **Display Library Integration** | Uses system/lib/display for formatted output | Consistent formatting with health tracking |

### Design Decisions

**Why find command instead of native Go filesystem walk?**

- Simple, widely available, battle-tested
- `-mmin` flag provides exact time-based filtering
- `-not -path` patterns for exclusions well-understood
- Future: Consider native Go walk for performance

**Why configurable exclusion patterns?**

- Different projects have different noise (node_modules, vendor, target, etc.)
- Build systems vary (npm, cargo, gradle, etc.)
- Users can add project-specific exclusions without code changes
- Reduces false positives in activity detection

**Why three display modes (detailed/summary/silent)?**

- <=10 files: Detailed count gives good context
- >10 but <=50: Summary awareness without overwhelming
- >50: Silent (likely build artifacts, not meaningful activity)
- Thresholds configurable for different workflows

---

## Public API

### CheckRecentActivity

**Purpose:** Display recently modified files in workspace for session context

**Signature:**

```go
func CheckRecentActivity(workspace string)
```

**Parameters:**

- `workspace` - Root directory to search for modified files

**Returns:** None (prints to stdout using display library)

**Example Usage:**

```go
import "hooks/lib/session"

// In session start hook
session.CheckRecentActivity("/home/user/project")
// Output: 📂 Recently modified (last hour): 7 file(s)

// In session stop hook
session.CheckRecentActivity(workspace)
// Shows what was worked on during session
```

**When to Use:**

- Session start: Show recent work context for awareness
- Session stop: Confirm what was modified during session
- Manual checks: See current workspace activity

**Health Scoring Impact:**

- Complete success: +30 pts (detection + display successful)
- Partial success: +20 pts (detection worked, display had issues)
- Command failure: -10 pts (find command failed)
- No activity: +5 pts (valid result, workspace quiet)

**Behavior:**

1. Builds find command with configured time window (-mmin)
2. Applies configured exclusion patterns (-not -path)
3. Executes find command, captures output
4. Formats using display library with configured thresholds
5. Prints to stdout (or stays silent if appropriate)

---

## Configuration System

### Configuration File Location

```
~/.claude/cpi-si/system/data/config/session/activity-tracking.jsonc
```

### Configuration Structure

```jsonc
{
  "metadata": {
    "name": "Activity Tracking Configuration",
    "version": "1.0.0",
    // ...
  },

  "time_window": {
    "minutes": 60,  // How far back to search
    "description": "How many minutes back to search for modified files"
  },

  "exclusion_patterns": [
    {
      "pattern": "*/.*",
      "description": "Hidden files and directories",
      "reasoning": "System/config files, not user work"
    },
    // ... more patterns
  ],

  "display": {
    "threshold_detailed": 10,   // Show count when <= this many
    "threshold_summary": 50,    // Show summary when > detailed but <= this
    "show_file_list": false,    // Display actual file paths
    "icon": "📂"                // Emoji for messages
  }
}
```

### Fallback Defaults

If configuration file missing or malformed:

- **Time window:** 60 minutes
- **Exclusions:** `*/.*`, `*/node_modules/*`
- **Detailed threshold:** 10 files
- **Summary threshold:** 50 files
- **File list display:** false
- **Icon:** 📂

### Extending Configuration

Add new patterns to `exclusion_patterns` array:

```jsonc
{
  "pattern": "*/build/*",
  "description": "Build output directories",
  "reasoning": "Compiled artifacts, not source work"
}
```

Adjust time windows:

```jsonc
"time_window": {
  "minutes": 120  // Look back 2 hours instead
}
```

Modify display thresholds:

```jsonc
"display": {
  "threshold_detailed": 20,  // Show count up to 20 files
  "threshold_summary": 100   // Summary up to 100 files
}
```

---

## Integration with Hooks

### Session Start Hook

```go
// In session/cmd-start/start.go
import "hooks/lib/session"

func startSession() {
    // ... other start logic ...

    // Show recent workspace activity for context
    session.CheckRecentActivity(workspace)

    // ... continue session start ...
}
```

**Purpose:** Ground user in recent work context at session start

### Session Stop Hook

```go
// In session/cmd-stop/stop.go
import "hooks/lib/session"

func stopSession() {
    // ... stop logic ...

    // Show what was modified during session
    session.CheckRecentActivity(workspace)

    // ... complete session stop ...
}
```

**Purpose:** Confirm what work was done before ending session

---

## Display Behavior

### Output Modes

**Detailed Mode (count <= threshold_detailed):**

```
📂 Recently modified (last hour): 7 file(s)
```

Shows specific count. Default threshold: 10 files.

**Summary Mode (threshold_detailed < count <= threshold_summary):**

```
📂 Recently modified (last hour): 35 files
```

General awareness without specifics. Default threshold: 50 files.

**Silent Mode (count > threshold_summary):**

No output. Likely build artifacts or mass changes, not meaningful context.

### File List Display (Optional)

When `show_file_list: true` in config (and count <= threshold_detailed):

```
📂 Recently modified (last hour): 3 file(s)
  - src/main.go
  - config/settings.json
  - README.md
```

Default: disabled (prevents output spam).

---

## Error Handling & Graceful Degradation

### Configuration Loading Failures

**Behavior:** Falls back to hardcoded defaults

**Causes:**
- File missing
- Permission denied
- Malformed JSON
- Invalid structure

**Result:** Library continues working with sensible defaults

### Find Command Failures

**Behavior:** Returns empty slice, stays silent

**Causes:**
- Invalid workspace path
- Permission denied
- find command not available

**Result:** No error thrown, hook continues normally

### Display Library Failures

**Behavior:** Falls back to plain fmt output

**Causes:**
- Display library unavailable
- Formatting error

**Result:** Message still displayed, just without formatting

**Non-Blocking Guarantee:**

This library **NEVER** blocks hook execution. All errors handled gracefully, allowing session operations to continue smoothly.

---

## Performance Considerations

### Time Complexity

**O(n)** where n = number of files in workspace

Limited by find command filesystem traversal speed.

### Memory Usage

**O(m)** where m = number of recently modified files

Find output captured in strings, scales with result count.

### Bottlenecks

**Large Workspaces (>100K files):**
- find command becomes slow
- Solution: Reduce time_window.minutes (fewer files to check)
- Solution: Add more exclusion_patterns (skip large directories)

**Many Exclusions (>20 patterns):**
- Each pattern adds overhead to find traversal
- Impact usually minor, but consider consolidation

**File List Display:**
- Disabled by default
- Large counts can spam output
- Only enable for debugging

### Optimization Strategies

1. **Reduce time window** - Check last 30 minutes instead of 60
2. **Add exclusions** - Skip known large directories (build/, dist/, vendor/)
3. **Increase thresholds** - Show summary for larger counts
4. **Disable file list** - Keep output concise

---

## Modification Policy

### Safe to Modify (Extension Points)

✅ **Add exclusion patterns** - Via configuration file (no code changes)
✅ **Adjust time windows** - Via configuration file
✅ **Modify thresholds** - Via configuration file
✅ **Add display modes** - Extend formatOutput() function
✅ **Add helper functions** - For additional detection modes

### Modify with Extreme Care (Breaking Changes)

⚠️ **CheckRecentActivity() signature** - Breaks all calling hooks (start.go, stop.go)
⚠️ **ActivityConfig struct fields** - Breaks configuration parsing
⚠️ **Display threshold defaults** - Affects all users without custom config
⚠️ **Exclusion pattern format** - Breaks existing configuration files
⚠️ **Output format** - Affects hook display expectations

### NEVER Modify (Foundational)

❌ **4-block structure** - METADATA, SETUP, BODY, CLOSING organization
❌ **Configuration fallback behavior** - Must degrade gracefully
❌ **Non-blocking guarantee** - Never block hook execution
❌ **Display library integration** - Health tracking dependency
❌ **init() loading pattern** - Configuration must load at import

---

## Troubleshooting

### Problem: No activity displayed even with recent work

**Possible Causes:**
- Time window too short (default 60 minutes)
- Exclusion patterns catching work files
- Display threshold exceeded (silent mode)

**Solutions:**
1. Check time_window.minutes in config (increase if needed)
2. Review exclusion_patterns (may be too broad)
3. Check display.threshold_summary (may need to increase)
4. Verify workspace path is correct

**Diagnosis:**
```bash
# Manually run find command to see what's detected
find /workspace/path -type f -mmin -60 -not -path '*/.*'
```

### Problem: Configuration not loading (using fallback defaults)

**Possible Causes:**
- File missing
- File permissions incorrect
- JSONC syntax error

**Solutions:**
1. Verify file exists:
   ```bash
   ls -la ~/.claude/cpi-si/system/data/config/session/activity-tracking.jsonc
   ```

2. Check permissions (must be readable):
   ```bash
   chmod 644 ~/.claude/cpi-si/system/data/config/session/activity-tracking.jsonc
   ```

3. Validate JSONC syntax:
   ```bash
   # Strip comments and validate JSON
   grep -v "^//" activity-tracking.jsonc | jq .
   ```

### Problem: Too much activity displayed (overwhelming output)

**Possible Causes:**
- Too few exclusion patterns for your workflow
- Display thresholds set too high

**Solutions:**
1. Add more exclusion patterns to config:
   ```jsonc
   {
     "pattern": "*/your-build-dir/*",
     "description": "Your build artifacts",
     "reasoning": "Not meaningful activity"
   }
   ```

2. Lower threshold_summary:
   ```jsonc
   "display": {
     "threshold_summary": 20  // Quieter behavior
   }
   ```

### Problem: Activity detection slow in large workspaces

**Cause:** find command traversing many files

**Solutions:**
1. Reduce time window:
   ```jsonc
   "time_window": {"minutes": 30}  // Half the search time
   ```

2. Add more exclusions:
   ```jsonc
   "exclusion_patterns": [
     {"pattern": "*/node_modules/*", ...},
     {"pattern": "*/vendor/*", ...},
     {"pattern": "*/target/*", ...}
   ]
   ```

**Note:** This is find command limitation, not library issue. Native Go walk could improve performance in future.

---

## Future Roadmap

### Planned Features

✓ Configuration system - **COMPLETED** (v2.0.0)
✓ Display library integration - **COMPLETED** (v2.0.0)
✓ Comprehensive exclusion patterns - **COMPLETED** (v2.0.0)
⏳ File type filters (show only .go, .rs, etc.)
⏳ Size-based filtering (ignore large binaries)
⏳ User-based filtering (show only my modifications)
⏳ Git integration (show only uncommitted changes)

### Research Areas

- **Alternative to find command** - Native Go filesystem walk for speed
- **Incremental scanning** - Cache results, only check new modifications
- **Language-specific activity** - Detect which projects are active
- **Intelligent exclusions** - Learn what user typically ignores
- **Activity correlation** - Relate to session duration, time of day

### Integration Targets

- **Git integration** - Uncommitted changes, branch activity
- **Session patterns** - Correlate activity with time of day, duration
- **Project detection** - Identify active projects from file paths
- **Health scoring aggregation** - Track detection quality over time
- **Temporal awareness** - Celestial correlation (work during daylight?)

### Known Limitations to Address

1. **find command performance** - Large workspaces (>100K files) slow
2. **No caching** - Every call rescans filesystem from scratch
3. **No git-awareness** - Shows all files, not just version-controlled
4. **Comment stripping** - Line-based only (no inline //, no /* */)
5. **No content analysis** - Just modification time, not what changed
6. **No categorization** - All files treated equally

---

## Version History

### v2.0.0 (2025-11-12) - Configuration System & Template Alignment

**Major Changes:**
- Created activity-tracking.jsonc configuration file
- Integrated system/lib/display for formatted output
- Added comprehensive exclusion patterns (8 default patterns)
- Implemented graceful fallback to hardcoded defaults
- Full 4-block template alignment (METADATA/SETUP/BODY/CLOSING)
- Added complete inline documentation and health scoring

**Breaking Changes:**
- None (v1.0.0 was internal only)

**Migration:**
- No migration needed for external users
- Config file optional (falls back to defaults)

### v1.0.0 (2024-10-24) - Initial Implementation

**Features:**
- Basic find command execution with hardcoded 60-minute window
- Hardcoded exclusions (dotfiles, node_modules only)
- Simple threshold display logic (10 files)
- Direct fmt.Printf output (no display library)
- Minimal documentation

---

**For questions, issues, or contributions:**
- Review modification policy above
- Follow 4-block structure pattern
- Test thoroughly (go build, go vet, integration tests)
- Update activity-tracking.jsonc schema if adding config fields
- Verify fallback behavior still works

*"Whatever your hand finds to do, do it with all your might" - Ecclesiastes 9:10*
