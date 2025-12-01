# Session Reminders API Documentation

**Library**: `hooks/lib/session` (reminders.go)
**Version**: 2.0.0
**Last Updated**: 2025-11-12

---

## Table of Contents

1. [Overview](#overview)
2. [Public API](#public-api)
3. [Configuration System](#configuration-system)
4. [Integration Patterns](#integration-patterns)
5. [Troubleshooting](#troubleshooting)
6. [Future Roadmap](#future-roadmap)

---

## Overview

### Purpose

The Session Reminders library provides gentle, configurable awareness of uncommitted work at session end. It prevents developers from forgetting changes, maintains workspace discipline, and provides contextual reminders without interrupting workflow.

### Design Rationale

**Non-Blocking by Design**: Reminder failures never interrupt session flow. If git checks fail or configuration is missing, the session continues without disruption.

**Configuration-Driven**: Reminder behavior, display settings, and threshold preferences are externalized to JSONC configuration, allowing customization without code changes.

**Threshold-Based**: Only shows reminders when uncommitted changes meet or exceed configured threshold (default: any changes).

**Graceful Fallback**: Works with hardcoded defaults if configuration missing, ensuring reminders always available.

**Biblical Foundation**: "Remember therefore from whence thou art fallen, and repent" - Revelation 2:5 (KJV). Gentle reminders prevent forgetfulness.

### Key Features

- Configurable threshold (customize when reminders appear)
- Non-blocking execution (reminder failures don't interrupt sessions)
- Display customization (icons, messages, formatting)
- Git repository awareness (optional git-only checking)
- Message templating ({count} placeholder for dynamic content)
- Silent failure mode (continue without warnings)
- Graceful fallback (works with defaults if config missing)

---

## Public API

### RemindUncommittedWork()

Detects and displays gentle reminder about uncommitted work at session end.

**Signature**:
```go
func RemindUncommittedWork(workspace string)
```

**Parameters**:
- `workspace` (string): Path to workspace directory to check for uncommitted changes

**Returns**: None (outputs directly to stdout if reminder threshold met)

**Behavior**:
1. Checks if reminders enabled in configuration
2. Retrieves behavior settings (displayEnabled, checkGitOnly, silentFailures)
3. Verifies workspace is git repository (if checkGitOnly enabled)
4. Gets uncommitted change count from git
5. Formats message using template with threshold check
6. Prints to stdout if threshold met

**Example Usage**:
```go
import "hooks/lib/session"

func sessionEnd() {
    workspace := "/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC"
    session.RemindUncommittedWork(workspace)
    // Output: ⚠️  Reminder: 15 uncommitted change(s) in workspace
}
```

**Health Impact**:
- Reminder displayed successfully: +40 points
- Reminder skipped (below threshold): +30 points
- Configuration loaded correctly: +20 points
- Git check succeeded: +10 points
- Display disabled (intentional): 0 points
- Git check failed: -10 points

**When to Use**: Call at session end to remind developer of uncommitted work before closing session.

---

## Configuration System

### Configuration File Location

**Path**: `~/.claude/cpi-si/system/data/config/session/reminders.jsonc`

**Format**: JSONC (JSON with comments)

**Validation**: Falls back to hardcoded defaults if file missing or invalid

### Configuration Structure

```jsonc
{
  "metadata": {
    "name": "Session Reminders Configuration",
    "description": "Controls reminder display and behavior at session end",
    "version": "1.0.0",
    "author": "Your Name",
    "created": "YYYY-MM-DD",
    "last_updated": "YYYY-MM-DD"
  },

  "reminders": {
    "uncommitted_work": {
      "enabled": true,                     // Show uncommitted work reminder
      "icon": "⚠️",                        // Icon for reminder
      "message": "Reminder: {count} uncommitted change(s) in workspace",
      "threshold": 0,                      // Minimum changes to trigger (0 = any)
      "show_details": false                // Future: show file list
    }
  },

  "display": {
    "enabled": true,                       // Master switch for all reminders
    "prefix_newline": true,                // Add newline before reminders
    "group_reminders": false,              // Future: group multiple reminders together
    "color_coding": false                  // Future: color code by severity
  },

  "behavior": {
    "silent_failures": true,               // Continue silently on error
    "check_git_only": true,                // Only check if workspace is git repo
    "cache_results": false                 // Future: cache results within session
  },

  "extensions": {
    "custom_reminders": [],                // Future: custom reminder hooks
    "severity_levels": {},                 // Future: categorize by importance
    "notification_hooks": []               // Future: external notifications
  }
}
```

### Configuration Options

#### Reminders Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `uncommitted_work.enabled` | bool | true | Whether to show uncommitted work reminder |
| `uncommitted_work.icon` | string | "⚠️" | Icon prefix for reminder display |
| `uncommitted_work.message` | string | "Reminder: {count} uncommitted change(s) in workspace" | Message template ({count} placeholder) |
| `uncommitted_work.threshold` | int | 0 | Minimum changes to trigger reminder (0 = any changes) |
| `uncommitted_work.show_details` | bool | false | Future: show file list in reminder |

**Message Template**:
- Use `{count}` placeholder for uncommitted change count
- Example: `"You have {count} unsaved file(s)"` → `"You have 15 unsaved file(s)"`
- Placeholder replaced at runtime with actual count

**Threshold Behavior**:
- `threshold: 0` - Show reminder for any changes (1 or more)
- `threshold: 5` - Only show if 5 or more uncommitted changes
- `threshold: 10` - Only show if 10 or more uncommitted changes
- Below threshold: No reminder displayed

#### Display Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | true | Master switch for all reminders (overrides individual settings) |
| `prefix_newline` | bool | true | Add newline before reminder for visual separation |
| `group_reminders` | bool | false | Future: group multiple reminders together |
| `color_coding` | bool | false | Future: color code reminders by severity |

**Display Formatting**:
- `prefix_newline: true` - Adds `\n` before reminder for spacing
- `prefix_newline: false` - Reminder appears immediately after previous output
- Icon + message format: `{icon}  {message}\n`
- Example: `⚠️  Reminder: 15 uncommitted change(s) in workspace\n`

#### Behavior Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `silent_failures` | bool | true | Continue silently if checks fail (no warnings) |
| `check_git_only` | bool | true | Only check workspaces that are git repositories |
| `cache_results` | bool | false | Future: cache git status within session |

**Behavior Details**:

**silent_failures**:
- `true` - Errors/warnings not displayed (non-git repo, git check failure, etc.)
- `false` - Shows warning when workspace not git repo: `⚠️  Workspace is not a git repository: /path`

**check_git_only**:
- `true` - Skip reminder if workspace not a git repository (silent return)
- `false` - Attempt reminder even for non-git workspaces (will fail gracefully)

**cache_results** (Future):
- `true` - Cache git status results for session duration
- `false` - Check git status fresh each time (current behavior)

### Health Scoring

Configuration and execution tracked through health scoring system:

**Configuration Loading** (20 points):
- Config loaded successfully: +20
- Config missing (using defaults): +15
- Config invalid (using defaults): +10

**Git Repository Check** (10 points):
- Workspace is git repo: +10
- Workspace not git repo (check_git_only enabled): 0
- Git check failed: -10

**Reminder Display** (40 points):
- Reminder displayed (threshold met): +40
- Reminder skipped (below threshold): +30
- Display disabled (intentional): 0
- Display failed: -20

**Message Formatting** (30 points):
- Message formatted correctly: +30
- Message has minor issues: +15
- Message formatting failed: 0

**Total**: 100 points possible

---

## Integration Patterns

### Basic Integration

```go
package main

import "hooks/lib/session"

func main() {
    workspace := "/path/to/project"

    // ... session work ...

    // At session end
    session.RemindUncommittedWork(workspace)
}
```

### Session End Hook Integration

```go
// In session/cmd-end/end.go
package main

import (
    "hooks/lib/session"
    "os"
)

func main() {
    // Get workspace from environment or argument
    workspace := os.Getenv("CLAUDE_WORKSPACE")
    if workspace == "" {
        workspace, _ = os.Getwd()
    }

    // ... other session end logic ...

    // Remind about uncommitted work
    session.RemindUncommittedWork(workspace)

    // ... finish session cleanup ...
}
```

### Multi-Workspace Integration

```go
// Check multiple workspaces
workspaces := []string{
    "/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC",
    "/home/seanje-lenox-wise/.claude/hooks",
    "/tmp",
}

for _, workspace := range workspaces {
    session.RemindUncommittedWork(workspace)
}
```

### Conditional Reminders

Reminders can be disabled entirely through configuration without code changes:

```jsonc
{
  "display": {
    "enabled": false  // Disables all reminders
  }
}
```

Or disable specific reminder types:

```jsonc
{
  "reminders": {
    "uncommitted_work": {
      "enabled": false  // Disable only uncommitted work reminder
    }
  }
}
```

Or adjust threshold to only show for significant changes:

```jsonc
{
  "reminders": {
    "uncommitted_work": {
      "threshold": 10  // Only remind if 10+ uncommitted changes
    }
  }
}
```

---

## Troubleshooting

### Problem: No reminder even when uncommitted changes exist

**Symptoms**: Uncommitted work present but no reminder displayed

**Possible Causes**:
1. Reminders disabled in configuration
2. Threshold set too high (changes below threshold)
3. Workspace not a git repository (check_git_only enabled)
4. Display disabled globally

**Diagnostic Steps**:
```bash
# Check if display enabled
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc | grep -A5 '"display"'

# Check threshold setting
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc | grep -A10 '"uncommitted_work"'

# Verify uncommitted changes exist
cd /path/to/workspace
git status --porcelain | wc -l

# Check if workspace is git repo
git rev-parse --git-dir
```

**Solutions**:
- Set `display.enabled: true` in configuration
- Set `reminders.uncommitted_work.enabled: true` in configuration
- Lower threshold: `threshold: 0` for any changes
- Verify workspace is actually a git repository
- Check git status manually to confirm changes exist

---

### Problem: Reminder appears when no changes exist

**Symptoms**: Reminder displays but `git status` shows clean working tree

**Possible Causes**:
1. Cached git status (future feature)
2. Configuration loading incorrectly
3. Threshold set to negative value

**Diagnostic Steps**:
```bash
# Check git status manually
cd /path/to/workspace
git status

# Check threshold value
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc | grep threshold

# Verify configuration loaded correctly
# (Check logs or add debug output)
```

**Solutions**:
- Ensure `threshold: 0` or positive value
- Restart session to reload configuration
- Verify git repository is in expected state

---

### Problem: Warning appears: "Workspace is not a git repository"

**Symptoms**: Warning displayed even though workspace is valid

**Expected Behavior**: This warning only appears when `silent_failures: false` and workspace is not a git repository

**Diagnostic Steps**:
```bash
# Check silent_failures setting
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc | grep silent_failures

# Verify workspace is git repo
cd /path/to/workspace
git rev-parse --git-dir
```

**Solutions**:
- Set `behavior.silent_failures: true` to suppress warnings
- Ensure workspace path points to actual git repository
- Set `behavior.check_git_only: true` to skip non-git workspaces silently

---

### Problem: Custom message template not working

**Symptoms**: Changed message in configuration but old message still appears

**Possible Causes**:
1. Configuration not reloaded (session still running)
2. Syntax error in configuration file
3. Message missing {count} placeholder

**Diagnostic Steps**:
```bash
# Validate JSONC syntax
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc

# Check message value
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc | grep message
```

**Solutions**:
- Restart session to reload configuration (config loaded once at init)
- Ensure message includes `{count}` placeholder for dynamic count
- Verify JSONC syntax is valid (no trailing commas, quotes correct)
- Example valid message: `"You have {count} unsaved change(s)"`

---

### Problem: Threshold not working as expected

**Symptoms**: Reminder appears even when below threshold, or doesn't appear when above threshold

**Possible Causes**:
1. Threshold value incorrect type (string instead of int)
2. Configuration not loaded
3. Fallback to default threshold (0)

**Diagnostic Steps**:
```bash
# Check threshold type and value
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc | grep -A1 threshold

# Test with specific threshold
# (Temporarily set to known value like 5)

# Count uncommitted changes manually
cd /path/to/workspace
git status --porcelain | wc -l
```

**Solutions**:
- Ensure threshold is integer, not string: `"threshold": 5` (not `"threshold": "5"`)
- Set explicit threshold value in configuration
- Restart session to reload configuration
- Verify uncommitted count matches expectations

---

### Problem: Configuration changes not taking effect

**Symptoms**: Changed configuration but behavior unchanged

**Cause**: Configuration loaded once at package import (init function)

**Solution**:
- Restart session to reload configuration
- Package only reads configuration file once when imported
- Changes require session restart to take effect
- No dynamic reloading currently implemented

---

## Future Roadmap

### Planned Features

✓ **Configuration-driven reminders** - COMPLETED (v2.0.0)
✓ **Threshold-based display** - COMPLETED (v2.0.0)
✓ **Message templating** - COMPLETED (v2.0.0)
⏳ **Multiple reminder types** - Long-running processes, unsaved buffers, etc.
⏳ **File detail display** - Show which files are uncommitted
⏳ **Severity levels** - Categorize reminders by importance
⏳ **Custom reminder hooks** - User-defined reminder scripts
⏳ **Cache support** - Cache git status within session

### Research Areas

- **Multiple workspace awareness**: Track multiple workspaces in single session
- **Reminder grouping**: Display multiple reminders together with visual grouping
- **Color coding**: Use terminal colors to indicate severity
- **Interactive remediation**: Offer to commit/stash changes directly from reminder
- **Cross-session persistence**: Track uncommitted work across multiple sessions
- **Integration with session patterns**: Learn typical commit patterns, suggest when unusual

### Integration Targets

- **Health scoring system**: Track reminder display success rates over time
- **Session logging**: Record reminder activity for pattern analysis
- **Pattern learning**: Identify typical commit behavior and workflow
- **Display customization**: Color coding, formatting options, severity indicators

### Extensibility Design

Configuration designed for future expansion:

**Adding New Reminder Types**:
```jsonc
{
  "reminders": {
    "uncommitted_work": { /* existing */ },
    "long_running_processes": {
      "enabled": true,
      "icon": "🔌",
      "message": "{count} process(es) still running",
      "threshold": 1
    },
    "unsaved_buffers": {
      "enabled": true,
      "icon": "📝",
      "message": "{count} unsaved buffer(s)",
      "threshold": 0
    }
  }
}
```

**Custom Reminder Hooks** (Future):
```jsonc
{
  "extensions": {
    "custom_reminders": [
      {
        "name": "check_todos",
        "script": "/path/to/check_todos.sh",
        "icon": "✓",
        "threshold": 1
      }
    ]
  }
}
```

### Known Limitations

- **Single workspace per call**: RemindUncommittedWork checks one workspace at a time
- **Git-only**: No support for other version control systems (svn, hg, etc.)
- **Configuration requires restart**: Changes don't apply mid-session
- **No interactive remediation**: Can't commit/stash directly from reminder
- **No cross-session tracking**: Each session checks git status fresh
- **No file details**: Only shows count, not which files uncommitted

### Version History

**2.0.0 (2025-11-12)** - Configuration-driven reminders
- Added `reminders.jsonc` configuration file support
- Implemented threshold-based display (configurable minimum changes)
- Message templating with `{count}` placeholder
- Display customization (icons, messages, prefix newline)
- Behavior customization (silent failures, git-only checking)
- Graceful fallback to defaults if config missing
- All constants properly used (displayEnabled, checkGitOnly, silentFailures)

**1.0.0 (2024-10-24)** - Initial implementation
- Hardcoded message format and icon
- Fixed threshold (any changes)
- Basic git integration via system/lib/git
- Simple output to stdout
- Silent failure design established

---

## Quick Reference

### Default Configuration Values

| Setting | Default Value | Description |
|---------|---------------|-------------|
| Icon | ⚠️ | Warning triangle emoji |
| Message | "Reminder: {count} uncommitted change(s) in workspace" | Template with placeholder |
| Threshold | 0 | Show for any changes (1 or more) |
| Display Enabled | true | Reminders shown by default |
| Prefix Newline | true | Visual separation before reminder |
| Silent Failures | true | No warnings on errors |
| Check Git Only | true | Skip non-git workspaces |

### Common Commands

```bash
# Check configuration
cat ~/.claude/cpi-si/system/data/config/session/reminders.jsonc

# Verify uncommitted changes
cd /path/to/workspace
git status --porcelain | wc -l

# See what would be reminded
git status --short

# Test git repository detection
git rev-parse --git-dir

# Check if workspace is git repo (exit code)
git rev-parse --is-inside-work-tree
```

### Configuration Quick Edits

```bash
# Disable all reminders
# Set: "enabled": false in display section

# Disable uncommitted work reminder only
# Set: "enabled": false in uncommitted_work section

# Change threshold to 10 changes
# Set: "threshold": 10 in uncommitted_work section

# Show warnings on errors
# Set: "silent_failures": false in behavior section

# Customize message
# Set: "message": "You forgot {count} file(s)!" in uncommitted_work section

# Change icon
# Set: "icon": "💾" in uncommitted_work section

# Remove newline spacing
# Set: "prefix_newline": false in display section
```

---

## Related Documentation

- **4-Block Structure**: See `standards/code/4-block/CWS-STD-001-DOC-4-block.md`
- **Configuration Standards**: See `standards/code/4-block/CWS-STD-003-DOC-configuration.md`
- **Health Scoring**: See `knowledge-base/algorithms/CPSI-ALG-001-DOC-health-scoring.md`
- **Session Library Overview**: See `hooks/lib/session/README.md`
- **Git Library API**: See `system/lib/git/README.md` (GetInfo, IsGitRepository functions)

---

**For questions, issues, or contributions**:
- Review the modification policy in reminders.go CLOSING section
- Follow the 4-block structure pattern
- Test thoroughly before committing (`go build && go vet`)
- Document all changes comprehensively (What/Why/How pattern)
- Verify non-blocking behavior preserved

*"Remember now thy Creator in the days of thy youth" - Ecclesiastes 12:1 (KJV)*
