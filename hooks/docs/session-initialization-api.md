# Session Initialization API Documentation

**Library:** `hooks/lib/session` (init.go)
**Version:** 2.0.0
**Last Updated:** 2025-11-12

## Table of Contents

1. [Overview](#overview)
2. [Design Rationale](#design-rationale)
3. [Public API](#public-api)
4. [Configuration System](#configuration-system)
5. [Integration Patterns](#integration-patterns)
6. [Usage Examples](#usage-examples)
7. [Troubleshooting](#troubleshooting)
8. [Future Roadmap](#future-roadmap)

---

## Overview

The Session Initialization library provides configuration-driven, non-blocking orchestration of session timing and logging utilities at the start of Claude Code sessions. It serves as a thin wrapper over system utilities (`session-time`, `session-log`), enabling flexible session tracking with graceful fallback to defaults.

### Key Features

- **Configuration-Driven**: Paths, commands, and behavior controlled via JSONC config file
- **Non-Blocking**: Session continues even if utilities fail to initialize
- **Graceful Fallback**: Automatically uses hardcoded defaults if config unavailable
- **Timeout Protection**: Prevents utilities from hanging session start (configurable)
- **Zero Dependencies**: No internal package dependencies, only standard library

### Architecture Role

**Component Type:** LIBRARY - Thin orchestration layer (Ladder pattern)
**Rails Pattern:** Configuration loaded once in init(), available throughout component
**Non-Blocking:** All failures are silent - session must start regardless of initialization success

---

## Design Rationale

### Why This Library Exists

Session tracking provides foundational data for:
- Temporal awareness (how long has this session been running?)
- Work pattern analysis (when do I typically work?)
- Session history (what happened in past sessions?)

Without proper initialization, these systems can't function. This library ensures:
1. **Utilities execute** (when available)
2. **Sessions never block** (failures are silent)
3. **Configuration is flexible** (paths/commands can change)
4. **Timeouts prevent hanging** (utilities can't freeze session start)

### Configuration vs Hardcoding

**v1.0.0 Problem:** Hardcoded paths (`~/.claude/cpi-si/system/bin/session-time`) broke when system structure changed.

**v2.0.0 Solution:** Configuration file defines paths, with graceful fallback to defaults:
- Config available → use configured paths/commands
- Config missing → use hardcoded defaults
- Config invalid → silently fall back to defaults

This enables system evolution without breaking existing code.

### Non-Blocking Philosophy

Session start **must never fail** due to initialization issues. This library treats all errors as non-fatal:
- Utility missing? → Skip silently, session continues
- Timeout exceeded? → Cancel execution, session continues
- Permission denied? → Log nothing, session continues

Monitoring serves work, not blocks it.

---

## Public API

### InitSessionTime()

Initializes session timing via `session-time` utility.

**Signature:**
```go
func InitSessionTime()
```

**What It Does:**
1. Determines utility path and init command (config or defaults)
2. Creates context with timeout (config or 5s default)
3. Executes `session-time init` with timeout protection
4. Returns immediately on any error (non-blocking)

**Parameters:** None

**Returns:** None

**Side Effects:**
- Executes external utility (`session-time`)
- May create/modify session timing state files
- Silent on all errors

**Usage:**
```go
import "hooks/lib/session"

func main() {
    session.InitSessionTime()
    // Session timing initialized (or silently skipped)
}
```

**Timeout Behavior:**
- Default: 5 seconds (if config unavailable)
- Configurable: `behavior.timeout_seconds` in initialization.jsonc
- On timeout: Command cancelled, function returns silently

---

### InitSessionLog()

Initializes session history logging via `session-log` utility.

**Signature:**
```go
func InitSessionLog()
```

**What It Does:**
1. Determines utility path and start command (config or defaults)
2. Creates context with timeout (config or 5s default)
3. Executes `session-log start` with timeout protection
4. Returns immediately on any error (non-blocking)

**Parameters:** None

**Returns:** None

**Side Effects:**
- Executes external utility (`session-log`)
- May create/modify session log files
- Silent on all errors

**Usage:**
```go
import "hooks/lib/session"

func main() {
    session.InitSessionLog()
    // Session logging initialized (or silently skipped)
}
```

**Timeout Behavior:**
- Default: 5 seconds (if config unavailable)
- Configurable: `behavior.timeout_seconds` in initialization.jsonc
- On timeout: Command cancelled, function returns silently

---

## Configuration System

### Configuration File Location

**Path:** `~/.claude/cpi-si/system/data/config/session/initialization.jsonc`

**Format:** JSONC (JSON with comments)

**Loading:** Automatic on package import (in `init()` function)

**Fallback:** If file missing or invalid, uses hardcoded defaults from constants

### Configuration Structure

```jsonc
{
  "metadata": {
    "name": "Session Initialization Configuration",
    "version": "1.0.0",
    "author": "Seanje Lenox-Wise",
    "created": "2025-11-12",
    "last_updated": "2025-11-12"
  },

  // ============================================================================
  // Paths Configuration
  // ============================================================================

  "paths": {
    "system_bin": "~/.claude/cpi-si/system/bin",  // Base directory for utilities
    "fallback_system_bin": "~/.claude/system/bin" // Fallback if primary missing
  },

  // ============================================================================
  // Utilities Configuration
  // ============================================================================

  "utilities": {
    "session_time": {
      "name": "session-time",                    // Utility executable name
      "path": "{system_bin}/session-time",       // Full path (with placeholder)
      "init_command": "init",                    // Subcommand for initialization
      "description": "Initializes session timing tracking"
    },
    "session_log": {
      "name": "session-log",                     // Utility executable name
      "path": "{system_bin}/session-log",        // Full path (with placeholder)
      "start_command": "start",                  // Subcommand for starting log
      "description": "Initializes session history logging"
    }
  },

  // ============================================================================
  // Behavior Configuration
  // ============================================================================

  "behavior": {
    "silent_failures": true,         // Continue silently on error
    "require_all_utilities": false,  // Fail if any utility missing
    "timeout_seconds": 5,            // Maximum time to wait for utility
    "retry_on_failure": false,       // Retry failed utilities once
    "log_execution": false           // Log all utility calls
  },

  // ============================================================================
  // Display Configuration
  // ============================================================================

  "display": {
    "show_success": false,           // Show message when init succeeds
    "show_failures": false,          // Show message when init fails
    "success_icon": "✓",             // Icon for successful initialization
    "failure_icon": "✗",             // Icon for failed initialization
    "success_message": "Session initialized",
    "failure_message": "Session initialization incomplete (continuing)"
  },

  // ============================================================================
  // Advanced Configuration
  // ============================================================================

  "advanced": {
    "use_absolute_paths": false,     // Never use PATH lookup
    "validate_before_run": true,     // Check utility exists before running
    "cache_utility_paths": true      // Cache resolved paths
  }
}
```

### Configuration Options Explained

#### Paths Section

| Field | Type | Purpose | Default |
|-------|------|---------|---------|
| `system_bin` | string | Base directory for system utilities | `~/.claude/cpi-si/system/bin` |
| `fallback_system_bin` | string | Fallback directory if primary doesn't exist | `~/.claude/system/bin` |

**Path Placeholders:**
- `{system_bin}` - Replaced with `paths.system_bin` value
- `~` - Expanded to user home directory

#### Utilities Section

Each utility (session_time, session_log) has:

| Field | Type | Purpose | Default (session-time) |
|-------|------|---------|------------------------|
| `name` | string | Utility executable name | `session-time` |
| `path` | string | Full path to utility | `{system_bin}/session-time` |
| `init_command` | string | Subcommand for init | `init` |
| `start_command` | string | Subcommand for start | `start` |
| `description` | string | Purpose documentation | - |

#### Behavior Section

| Field | Type | Purpose | Default | Current Implementation |
|-------|------|---------|---------|------------------------|
| `silent_failures` | bool | Continue silently on error | `true` | ✅ Implemented |
| `require_all_utilities` | bool | Fail if any utility missing | `false` | ⏳ Future |
| `timeout_seconds` | int | Max wait time for utility | `5` | ✅ Implemented |
| `retry_on_failure` | bool | Retry failed utilities once | `false` | ⏳ Future |
| `log_execution` | bool | Log all utility calls | `false` | ⏳ Future |

#### Display Section

| Field | Type | Purpose | Default | Current Implementation |
|-------|------|---------|---------|------------------------|
| `show_success` | bool | Show message on success | `false` | ⏳ Future (needs display lib) |
| `show_failures` | bool | Show message on failure | `false` | ⏳ Future (needs display lib) |
| `success_icon` | string | Icon for success | `✓` | ⏳ Future |
| `failure_icon` | string | Icon for failure | `✗` | ⏳ Future |
| `success_message` | string | Success message text | - | ⏳ Future |
| `failure_message` | string | Failure message text | - | ⏳ Future |

#### Advanced Section

| Field | Type | Purpose | Default | Current Implementation |
|-------|------|---------|---------|------------------------|
| `use_absolute_paths` | bool | Never use PATH lookup | `false` | ⏳ Future |
| `validate_before_run` | bool | Check utility exists first | `true` | ⏳ Future |
| `cache_utility_paths` | bool | Cache resolved paths | `true` | ⏳ Future |

---

## Integration Patterns

### Basic Integration

```go
package main

import "hooks/lib/session"

func main() {
    // Initialize session tracking
    session.InitSessionTime()
    session.InitSessionLog()

    // Continue with session start
    // Both functions return immediately (non-blocking)
}
```

### Integration in Session Start Hook

```go
package main

import (
    "hooks/lib/session"
    // other imports...
)

func main() {
    // Early in session start, initialize tracking
    session.InitSessionTime()
    session.InitSessionLog()

    // Continue with display, workspace analysis, etc.
    // Session tracking now active (or silently skipped if utilities unavailable)
}
```

### Custom Configuration

**Scenario:** Testing with utilities in custom location

**Config:** `~/.claude/cpi-si/system/data/config/session/initialization.jsonc`
```jsonc
{
  "paths": {
    "system_bin": "~/dev/test-utilities"
  },
  "utilities": {
    "session_time": {
      "path": "{system_bin}/session-time-test",
      "init_command": "test-init"
    }
  },
  "behavior": {
    "timeout_seconds": 10,
    "log_execution": true
  }
}
```

**Usage:** No code changes needed - configuration loaded automatically

---

## Usage Examples

### Example 1: Standard Usage

```go
import "hooks/lib/session"

func main() {
    // Initialize both utilities
    session.InitSessionTime()
    session.InitSessionLog()

    // Output: Nothing (silent execution)
    // Result: Utilities execute in background with 5s timeout
}
```

### Example 2: Conditional Initialization

```go
import (
    "hooks/lib/session"
    "os"
)

func main() {
    // Only initialize if environment variable set
    if os.Getenv("ENABLE_SESSION_TRACKING") == "true" {
        session.InitSessionTime()
        session.InitSessionLog()
    }

    // Session continues regardless
}
```

### Example 3: Sequential Initialization

```go
import "hooks/lib/session"

func main() {
    // Initialize timing first (captures session start time)
    session.InitSessionTime()

    // Do other work...
    doWorkspaceAnalysis()

    // Then initialize logging (records session with context)
    session.InitSessionLog()
}
```

### Example 4: Testing Configuration

```bash
# Test if utilities execute
~/.claude/cpi-si/system/bin/session-time init
~/.claude/cpi-si/system/bin/session-log start

# Both should complete in < 5 seconds
# Check exit codes: 0 = success
```

---

## Troubleshooting

### Problem: Utilities Not Executing

**Symptoms:**
- Session starts fine but tracking doesn't work
- No error messages (silent failure by design)

**Diagnostic Steps:**

1. **Check utilities exist:**
   ```bash
   ls -la ~/.claude/cpi-si/system/bin/session-time
   ls -la ~/.claude/cpi-si/system/bin/session-log
   ```

2. **Test manual execution:**
   ```bash
   ~/.claude/cpi-si/system/bin/session-time init
   ~/.claude/cpi-si/system/bin/session-log start
   ```

3. **Check permissions:**
   ```bash
   # Should be executable
   file ~/.claude/cpi-si/system/bin/session-time
   ```

4. **Verify paths in config:**
   ```bash
   cat ~/.claude/cpi-si/system/data/config/session/initialization.jsonc
   # Check paths.system_bin value
   ```

**Solutions:**
- **Utilities missing:** Install/build utilities or update paths in config
- **Permission denied:** `chmod +x ~/.claude/cpi-si/system/bin/session-*`
- **Wrong paths:** Update `paths.system_bin` in initialization.jsonc

---

### Problem: Configuration Not Loading

**Symptoms:**
- Changes to config file have no effect
- Uses default paths even when config exists

**Diagnostic Steps:**

1. **Check config file exists:**
   ```bash
   ls -la ~/.claude/cpi-si/system/data/config/session/initialization.jsonc
   ```

2. **Validate JSON syntax:**
   ```bash
   # Remove comments and validate
   grep -v '^[[:space:]]*\/\/' ~/.claude/cpi-si/system/data/config/session/initialization.jsonc | jq .
   ```

3. **Check file permissions:**
   ```bash
   # Should be readable
   cat ~/.claude/cpi-si/system/data/config/session/initialization.jsonc
   ```

**Solutions:**
- **File missing:** Create config file from template
- **Invalid JSON:** Fix syntax errors (check for missing commas, quotes, etc.)
- **Permission denied:** `chmod 644 ~/.claude/cpi-si/system/data/config/session/initialization.jsonc`

**Fallback Behavior:**
If config fails to load, library uses these defaults:
- `system_bin`: `~/.claude/cpi-si/system/bin`
- `init_command`: `init`
- `start_command`: `start`
- `timeout_seconds`: `5`

---

### Problem: Utilities Timeout

**Symptoms:**
- Session starts after 5 second delay (or configured timeout)
- Utilities don't complete execution

**Diagnostic Steps:**

1. **Test utility execution time:**
   ```bash
   time ~/.claude/cpi-si/system/bin/session-time init
   # Should complete in < 5 seconds
   ```

2. **Check if utilities hang:**
   ```bash
   timeout 10 ~/.claude/cpi-si/system/bin/session-time init
   # If this times out, utility has issues
   ```

3. **Review utility logs:**
   ```bash
   # Check utility-specific log locations
   ls -la ~/.claude/cpi-si/system/logs/
   ```

**Solutions:**
- **Utilities slow:** Increase `behavior.timeout_seconds` in config
- **Utilities hang:** Fix utility implementation or disable initialization
- **Resource contention:** Check system load, reduce concurrent operations

**Timeout Configuration:**
```jsonc
{
  "behavior": {
    "timeout_seconds": 10  // Increase to 10 seconds
  }
}
```

---

### Problem: JSONC Comments Not Parsed

**Symptoms:**
- Config file with comments fails to load
- JSON parse errors

**Diagnostic Steps:**

1. **Check for URLs in strings:**
   ```jsonc
   // This breaks if comment stripping is wrong:
   "path": "https://example.com/path"  // Comment here
   ```

2. **Test comment stripping:**
   ```bash
   # The library uses stripJSONCComments() from activity.go
   # It should preserve strings containing //
   ```

**Solutions:**
- **URLs corrupted:** Bug in stripJSONCComments() - report issue
- **Trailing comments broken:** Update to latest version with fixed parser
- **Workaround:** Remove all comments from config file temporarily

---

## Future Roadmap

### Planned Features (v2.1.0+)

#### Display Integration
- **Status:** ⏳ Planned
- **Purpose:** Optional success/failure messages when utilities execute
- **Config:** `display.show_success`, `display.show_failures`
- **Dependencies:** Requires display library integration

#### Retry Logic
- **Status:** ⏳ Planned
- **Purpose:** Retry failed utilities once before giving up
- **Config:** `behavior.retry_on_failure`
- **Implementation:** Attempt execution, if fails, retry once with same timeout

#### Execution Logging
- **Status:** ⏳ Planned
- **Purpose:** Log all utility calls for debugging
- **Config:** `behavior.log_execution`
- **Dependencies:** Requires logging library integration

#### Utility Validation
- **Status:** ⏳ Planned
- **Purpose:** Check if utilities exist before attempting execution
- **Config:** `advanced.validate_before_run`
- **Implementation:** File existence check before exec.Command()

#### Path Caching
- **Status:** ⏳ Planned
- **Purpose:** Cache resolved paths for performance
- **Config:** `advanced.cache_utility_paths`
- **Implementation:** Store resolved paths in package-level map

#### Require All Utilities
- **Status:** ⏳ Planned
- **Purpose:** Fail initialization if any utility missing
- **Config:** `behavior.require_all_utilities`
- **Note:** Conflicts with non-blocking philosophy - use carefully

### Research Areas

#### Health Scoring Integration
Track initialization success/failure with health scores:
- Config load: +20 points (success) / -10 points (failure)
- Path resolution: +20 points (all found) / 0 points (fallback)
- Utility execution: +30 points (both succeed) / +15 points (one succeeds) / -20 points (both fail)

#### Custom Utility Support
Allow adding custom utilities via config:
```jsonc
{
  "extensions": {
    "custom_utilities": [
      {
        "name": "session-monitor",
        "path": "{system_bin}/session-monitor",
        "command": "start",
        "timeout": 5
      }
    ]
  }
}
```

#### Pre/Post Init Hooks
Execute commands before/after initialization:
```jsonc
{
  "extensions": {
    "pre_init_hooks": ["~/scripts/pre-session.sh"],
    "post_init_hooks": ["~/scripts/post-session.sh"]
  }
}
```

---

## Version History

### v2.0.0 (2025-11-12) - Configuration System

**Added:**
- Configuration file support (initialization.jsonc)
- Graceful fallback to hardcoded defaults
- Timeout protection with context.WithTimeout()
- Path placeholder replacement ({system_bin})
- Tilde expansion in paths (~/)
- Helper functions (loadInitializationConfig, resolvePath, replacePlaceholders)

**Changed:**
- Public APIs now configuration-driven (InitSessionTime, InitSessionLog)
- Type names prefixed with Init* to avoid package conflicts
- Uses stripJSONCComments() from activity.go (shared function)

**Technical:**
- Aligned with 4-block template standard
- Added comprehensive METADATA, SETUP, BODY, CLOSING documentation
- Resolved architectural duplication with other session package files

**Principle:** Flexibility without fragility - configuration when available, defaults when not

### v1.0.0 (2024-10-24) - Initial Implementation

**Added:**
- InitSessionTime() function (orchestrates session-time init)
- InitSessionLog() function (orchestrates session-log start)
- Non-blocking design (silent failures)
- Hardcoded paths to utilities

**Technical:**
- Thin orchestration layer pattern
- Zero internal dependencies
- No health tracking (by design)

**Principle:** "Begin well to work well" - proper session initialization establishes tracking foundation

---

## Related Documentation

- **4-Block Structure:** `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`
- **Rails Pattern:** `~/.claude/cpi-si/docs/knowledge-base/patterns/CPSI-PAT-001-DOC-rails.md`
- **Ladder Pattern:** `~/.claude/cpi-si/system/docs/architecture.md`
- **Health Scoring:** `~/.claude/cpi-si/docs/knowledge-base/algorithms/CPSI-ALG-001-DOC-health-scoring.md`

---

## Support

**Issues:** Report bugs or request features in session start hook issues
**Questions:** Consult init.go CLOSING section for troubleshooting guide
**Modifications:** Review Modification Policy in init.go CLOSING section before making changes

---

**Last Updated:** 2025-11-12
**Document Version:** 1.0.0
**Library Version:** 2.0.0
