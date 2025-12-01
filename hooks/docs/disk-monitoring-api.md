# Disk Space Monitoring Library API

**Type:** Library
**Location:** `hooks/lib/session/disk.go`
**Purpose:** Workspace disk space monitoring with configurable warning thresholds
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

The **Disk Space Monitoring Library** provides workspace resource awareness by monitoring disk usage and warning at configurable thresholds. It enables:

- **Proactive warnings** - Alert before disk fills completely
- **Configurable thresholds** - Warning and critical levels customizable
- **Multi-level severity** - Different messages for warning vs critical
- **Graceful fallback** - Works with hardcoded defaults if configuration missing

**Biblical Foundation:**
*"The prudent see danger and take refuge, but the simple keep going and pay the penalty" - Proverbs 27:12 (WEB)*

Proactive resource monitoring prevents crises. Warn early, act before critical.

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Config-Driven** | All thresholds and messages via disk-monitoring.jsonc | Users customize without code changes |
| **Non-Blocking** | Failures never disrupt session operations | Monitoring serves workflow, doesn't block it |
| **Graceful Fallback** | Hardcoded defaults when config missing | Always works, even in degraded state |
| **Multi-Level Severity** | Warning (80%) and Critical (95%) thresholds | Escalating awareness based on urgency |
| **Display Library Integration** | Uses system/lib/display for formatted output | Consistent formatting with health tracking |

### Design Decisions

**Why percentage thresholds instead of absolute space?**

- Percentage applies to all disk sizes (500GB vs 2TB)
- Absolute space varies by system (10GB critical on 100GB disk, fine on 2TB)
- Industry standard approach (df -h shows percentages)
- Future: Add absolute minimum space option (e.g., "warn if <10GB")

**Why two threshold levels (warning and critical)?**

- Warning (80%): Time to clean up, not urgent
- Critical (95%): Immediate action needed
- Escalating urgency matches user mental model
- Different display formatting for different severity

**Why configurable messages with placeholders?**

- Users want different phrasing ("Low disk" vs "Disk space")
- Different units preferred (GB vs TB vs auto)
- Placeholders allow dynamic values: {percent}, {available}, {used}, {total}
- Configuration provides flexibility without code changes

---

## Public API

### CheckDiskSpace

**Purpose:** Monitor workspace disk usage and display warnings at configured thresholds

**Signature:**

```go
func CheckDiskSpace(workspace string)
```

**Parameters:**

- `workspace` - Root directory to check for disk usage

**Returns:** None (prints to stdout using display library)

**Example Usage:**

```go
import "hooks/lib/session"

// In session start hook
session.CheckDiskSpace("/home/user/project")
// Output: 💾 Disk Space Status
//            ⚠️  Disk space: 85% used (150GB available)
```

**When to Use:**

- Session start: Show disk health for awareness
- Pre-build validation: Ensure space before compilation
- Workspace switching: Verify new workspace has space

**Health Scoring Impact:**

- Disk healthy (below thresholds): +15 pts (good state)
- Warning threshold exceeded: +10 pts (successful detection)
- Critical threshold exceeded: +10 pts (successful critical detection)
- Check failed: -10 pts (system call error)

**Behavior:**

1. Checks if monitoring enabled in configuration
2. Gets disk usage via system.GetDiskUsage(workspace)
3. Compares usage against critical threshold (default 95%)
4. If >= critical: displays critical warning with display.Failure()
5. Else compares against warning threshold (default 80%)
6. If >= warning: displays warning with display.Warning()
7. Else if show_when_healthy enabled: displays success message
8. Otherwise: silent (healthy and show_when_healthy is false)

---

## Configuration System

### Configuration File Location

```
~/.claude/cpi-si/system/data/config/session/disk-monitoring.jsonc
```

### Configuration Structure

```jsonc
{
  "metadata": {
    "name": "Disk Space Monitoring Configuration",
    "version": "1.0.0",
    // ...
  },

  "thresholds": {
    "warning_percent": 80,   // Warn at 80% usage
    "critical_percent": 95,  // Critical at 95% usage
    "reasoning": {
      "warning": "80% allows time to clean up before critical",
      "critical": "95% indicates immediate action needed"
    }
  },

  "display": {
    "header_icon": "💾",
    "header_text": "Disk Space Status",
    "show_when_healthy": false,  // Show message when disk healthy
    "healthy_message": "Disk space: {percent}% used ({available} available)"
  },

  "messages": {
    "warning": "Disk space: {percent}% used ({available} available)",
    "critical": "⚠️  CRITICAL: Disk nearly full - {percent}% used (only {available} remaining)"
  },

  "behavior": {
    "enabled": true,              // Master switch
    "check_on_session_start": true
  }
}
```

### Fallback Defaults

If configuration file missing or malformed:

- **Warning threshold:** 80%
- **Critical threshold:** 95%
- **Display icon:** 💾
- **Header text:** "Disk Space Status"
- **Show when healthy:** false (silent when below thresholds)
- **Messages:** Standard warning/critical messages with placeholders
- **Behavior:** Enabled by default

### Extending Configuration

**Change threshold levels:**

```jsonc
"thresholds": {
  "warning_percent": 70,   // Warn earlier
  "critical_percent": 90   // More urgent critical
}
```

**Customize messages:**

```jsonc
"messages": {
  "warning": "Low disk space: {available} remaining",
  "critical": "🔴 URGENT: Only {available} left!"
}
```

**Show healthy status:**

```jsonc
"display": {
  "show_when_healthy": true,
  "healthy_message": "✓ Disk: {percent}% used"
}
```

**Disable monitoring:**

```jsonc
"behavior": {
  "enabled": false  // Skip all disk checks
}
```

### Message Placeholders

Available placeholders for message templates:

| Placeholder | Replaced With | Example Value |
|-------------|---------------|---------------|
| `{percent}` | Usage percentage (0-100) | `85` |
| `{available}` | Available space (human-readable) | `150GB` |
| `{used}` | Used space (human-readable) | `350GB` |
| `{total}` | Total space (human-readable) | `500GB` |

---

## Integration with Hooks

### Session Start Hook

```go
// In session/cmd-start/start.go
import "hooks/lib/session"

func startSession() {
    // ... other start logic ...

    // Check disk space before user starts work
    session.CheckDiskSpace(workspace)

    // ... continue session start ...
}
```

**Purpose:** Warn user about low disk space before they start working

### Pre-Build Validation

```go
// In custom build script
import "hooks/lib/session"

func preBuild() {
    session.CheckDiskSpace(workspace)
    // If critical, user sees warning before build attempt
}
```

**Purpose:** Prevent build failures due to insufficient disk space

---

## Display Behavior

### Output Modes

**Critical Level (>= 95%):**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━
 💾 Disk Space Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━

   ✗ ⚠️  CRITICAL: Disk nearly full - 97% used (only 15GB remaining)
```

Uses `display.Failure()` for red/bold critical formatting.

**Warning Level (80-94%):**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━
 💾 Disk Space Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━

   ⚠️  Disk space: 85% used (150GB available)
```

Uses `display.Warning()` for yellow warning formatting.

**Healthy (< 80%, default):**

No output. Silent when disk space is good.

**Healthy (< 80%, with show_when_healthy: true):**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━
 💾 Disk Space Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Disk space: 45% used (275GB available)
```

Uses `display.Success()` for green success formatting.

**Disabled (enabled: false):**

No output regardless of disk usage.

---

## Error Handling & Graceful Degradation

### Configuration Loading Failures

**Behavior:** Falls back to hardcoded defaults

**Causes:**
- File missing
- Permission denied
- Malformed JSON
- Invalid structure

**Result:** Library continues working with sensible defaults (80% warning, 95% critical)

### Disk Check Failures

**Behavior:** Silent (no error, no output)

**Causes:**
- Invalid workspace path
- Permission denied
- system.GetDiskUsage() failure

**Result:** Session continues normally without disk warnings

### Display Library Failures

**Behavior:** Falls back to plain fmt output

**Causes:**
- Display library unavailable
- Formatting error

**Result:** Messages still displayed, just without ANSI formatting

**Non-Blocking Guarantee:**

This library **NEVER** blocks hook execution. All errors handled gracefully, allowing session operations to continue smoothly.

---

## Performance Considerations

### Time Complexity

**O(1)** - Single system call to get disk usage

system.GetDiskUsage() calls `df -h` once per check.

### Memory Usage

**O(1)** - Fixed size message strings

Configuration loaded once at init(), message strings temporary.

### Bottlenecks

**Network-mounted filesystems:**
- Disk stat operations can be slow (100ms+ over network)
- Solution: Disable monitoring for networked workspaces

**Very busy systems:**
- Disk I/O contention delays stats
- Impact: Minimal (disk stats are cached by kernel)

### Optimization Strategies

1. **Configuration cached** - Loaded once at package init(), not per call
2. **Message formatting simple** - Plain string replacement, no templates
3. **No caching needed** - Disk usage changes rapidly, stale cache worse than fresh check

---

## Modification Policy

### Safe to Modify (Extension Points)

✅ **Add new threshold levels** - Extend ThresholdsConfig for info/severe levels
✅ **Enhance message formatting** - Improve formatMessage() function
✅ **Add placeholder types** - Filesystem name, mount point, inode usage
✅ **Improve display formatting** - Use different display library functions
✅ **Add absolute space checks** - Warn if < 10GB regardless of percentage

### Modify with Extreme Care (Breaking Changes)

⚠️ **CheckDiskSpace() signature** - Breaks all calling hooks
⚠️ **Configuration structure** - Breaks existing config files
⚠️ **Threshold defaults** - Affects users without custom config
⚠️ **Display library calls** - Ensure fallback still works

### NEVER Modify (Foundational)

❌ **4-block structure** - METADATA, SETUP, BODY, CLOSING organization
❌ **Non-blocking guarantee** - Never fail/panic on check errors
❌ **system/lib/system integration** - Standard disk usage API
❌ **Configuration fallback behavior** - Must degrade gracefully

---

## Troubleshooting

### Problem: Warnings not displayed even when disk is full

**Possible Causes:**
- Monitoring disabled (behavior.enabled: false)
- Workspace path incorrect
- Thresholds set too high (warning: 99%)
- system.GetDiskUsage() failing silently

**Solutions:**
1. Check `behavior.enabled` is true in config
2. Verify workspace path exists: `ls -la /workspace/path`
3. Check thresholds are reasonable (80, 95)
4. Test manually: `df -h /workspace/path`

**Diagnosis:**
```bash
# Check config exists and is readable
cat ~/.claude/cpi-si/system/data/config/session/disk-monitoring.jsonc

# Check actual disk usage
df -h /workspace/path
```

### Problem: Wrong threshold triggering (warns at wrong percentage)

**Cause:** Configuration thresholds set incorrectly

**Solutions:**
1. Verify thresholds in config: `cat ~/.claude/.../disk-monitoring.jsonc | grep percent`
2. Ensure warning_percent < critical_percent
3. Check for typos in percentage values

### Problem: Display formatting broken (no colors/icons)

**Possible Causes:**
- Terminal doesn't support ANSI colors
- Display library not compiled
- Output redirected to file (loses formatting)

**Solutions:**
1. Test terminal ANSI support: `echo -e "\033[31mRed\033[0m"`
2. Verify display library: `go build ./lib/display`
3. Check terminal emulator settings

### Problem: Configuration not loading (using fallback defaults)

**Possible Causes:**
- File missing at expected path
- File permissions don't allow reading
- JSONC syntax error (invalid JSON)

**Solutions:**
1. Verify file exists:
   ```bash
   ls -la ~/.claude/cpi-si/system/data/config/session/disk-monitoring.jsonc
   ```

2. Check permissions:
   ```bash
   chmod 644 ~/.claude/cpi-si/system/data/config/session/disk-monitoring.jsonc
   ```

3. Validate JSONC syntax:
   ```bash
   # Strip comments and validate
   grep -v "^//" disk-monitoring.jsonc | jq .
   ```

### Problem: Message placeholders not replaced ({percent} shown literally)

**Cause:** formatMessage() not processing placeholders correctly

**Solutions:**
1. Check message template syntax in config
2. Ensure placeholders use correct case: `{percent}` not `{Percent}`
3. Verify system.GetDiskUsage() returns valid data

---

## Future Roadmap

### Planned Features

✓ Warning threshold - **COMPLETED** (v2.0.0)
✓ Critical threshold - **COMPLETED** (v2.0.0)
✓ Configuration system - **COMPLETED** (v2.0.0)
✓ Message templates with placeholders - **COMPLETED** (v2.0.0)
⏳ Info threshold (50% for proactive awareness)
⏳ Absolute minimum space (warn if <10GB regardless of %)
⏳ Unit preference control (GB, TB, auto)
⏳ Multiple path monitoring (check several filesystems)

### Research Areas

- **Trend analysis** - Track disk usage over time, predict when full
- **Predictive warnings** - "Disk will be full in 3 days at current rate"
- **Cleanup suggestions** - Find large files, old caches, build artifacts
- **Integration with cleanup tools** - Automatic cleanup suggestions
- **Inode monitoring** - Warn when running out of inodes (not just space)

### Integration Targets

- **CI/CD pipelines** - Fail builds on insufficient disk space
- **Health scoring aggregation** - Track workspace health over time
- **Session patterns** - Correlate disk issues with work patterns
- **Automated cleanup hooks** - Trigger cleanup on critical threshold

### Known Limitations to Address

1. **Single workspace check** - Only checks one path (no multi-path support)
2. **No historical tracking** - Only current usage (no trends)
3. **No size-based warnings** - Only percentage (no "warn if <10GB")
4. **No cleanup automation** - Only warns, doesn't fix
5. **No inode monitoring** - Can run out of inodes with space available

---

## Version History

### v2.0.0 (2025-11-12) - Configuration System & Template Alignment

**Major Changes:**
- Created disk-monitoring.jsonc configuration file
- Integrated system/lib/display for formatted output
- Added warning and critical threshold levels
- Implemented customizable messages with placeholders
- Full 4-block template alignment (METADATA/SETUP/BODY/CLOSING)
- Added complete inline documentation and health scoring

**Breaking Changes:**
- None (v1.0.0 was internal only)

**Migration:**
- No migration needed for external users
- Config file optional (falls back to defaults)

### v1.0.0 (2024-10-24) - Initial Implementation

**Features:**
- Basic 80% threshold check
- Hardcoded warning message
- Plain fmt.Printf output (no display library)
- Minimal documentation
- No configuration system

---

**For questions, issues, or contributions:**
- Review modification policy above
- Follow 4-block structure pattern
- Test thoroughly (go build, go vet, integration tests)
- Update disk-monitoring.jsonc schema if adding config fields
- Verify fallback behavior still works

*"The prudent see danger and take refuge" - Proverbs 27:12*
