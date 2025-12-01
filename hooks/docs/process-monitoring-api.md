# Process Monitoring API Documentation

**Library**: `hooks/lib/session` (processes.go)
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

The Process Monitoring library provides configurable, non-blocking detection of active development servers on common ports at session start and end. It helps developers track what's running without manual checking, prevents forgotten background processes, and provides awareness of development environment state.

### Design Rationale

**Non-blocking by Design**: Port detection failures never interrupt session flow. If `lsof` fails or a port check times out, the session continues without disruption.

**Configuration-Driven**: Port list, display settings, and behavior preferences are externalized to JSONC configuration, allowing customization without code changes.

**Context-Aware**: Different messaging for session start ("Active dev servers") vs session end ("Dev servers still running") provides appropriate context.

**Timeout Protection**: All port checks execute with configurable timeout (default 2 seconds) to prevent hanging session start/end.

### Key Features

- Configurable port list (customize which ports to monitor)
- Non-blocking execution (detection failures don't interrupt sessions)
- Timeout protection (port checks can't hang session)
- Display customization (icons, messages, formatting)
- Graceful fallback (works with hardcoded defaults if config missing)
- Dual-context support (session start vs session end messaging)

---

## Public API

### CheckRunningProcesses()

Detects and reports active development servers at session start.

**Signature**:
```go
func CheckRunningProcesses()
```

**Parameters**: None

**Returns**: None (outputs directly to stdout if ports detected)

**Behavior**:
1. Checks if monitoring is enabled in configuration
2. Retrieves configured port list
3. Checks each port using `lsof` with timeout
4. Formats output with session-start context
5. Prints to stdout if processes found

**Example Usage**:
```go
import "hooks/lib/session"

func sessionStart() {
    session.CheckRunningProcesses()
    // Output: 🔌 Active dev servers on ports: 3000, 8080
}
```

**Health Impact**:
- All ports checked successfully: +40 points
- Partial checks succeed: +20 points
- All checks fail: -10 points

**When to Use**: Call at session start to notify developer of active processes.

---

### CheckRunningProcessesAsReminder()

Detects and reports active development servers at session end (reminder context).

**Signature**:
```go
func CheckRunningProcessesAsReminder()
```

**Parameters**: None

**Returns**: None (outputs directly to stdout if ports detected)

**Behavior**:
1. Checks if monitoring is enabled in configuration
2. Retrieves configured port list
3. Checks each port using `lsof` with timeout
4. Formats output with session-end context
5. Prints to stdout if processes found

**Example Usage**:
```go
import "hooks/lib/session"

func sessionEnd() {
    session.CheckRunningProcessesAsReminder()
    // Output: 🔌 Dev servers still running on ports: 3000, 8080
}
```

**Health Impact**:
- All ports checked successfully: +40 points
- Partial checks succeed: +20 points
- All checks fail: -10 points

**When to Use**: Call at session end to remind developer of processes still running.

---

## Configuration System

### Configuration File Location

**Path**: `~/.claude/cpi-si/system/data/config/session/processes.jsonc`

**Format**: JSONC (JSON with comments)

**Validation**: Falls back to hardcoded defaults if file missing or invalid

### Configuration Structure

```jsonc
{
  "metadata": {
    "name": "Process Monitoring Configuration",
    "version": "1.0.0",
    "author": "Your Name",
    "created": "YYYY-MM-DD",
    "last_updated": "YYYY-MM-DD"
  },

  "ports": {
    "enabled": true,                     // Master switch for monitoring
    "monitored_ports": [
      {
        "number": "3000",
        "description": "React/Next.js default dev server",
        "enabled": true
      },
      {
        "number": "8000",
        "description": "Django/Python HTTP server",
        "enabled": true
      }
      // ... more ports
    ],
    "custom_ports": ["9000", "4000"]     // Additional ports without descriptions
  },

  "display": {
    "show_at_start": true,               // Show at session start
    "show_at_end": true,                 // Show at session end
    "icon": "🔌",                        // Icon for notification
    "start_message": "Active dev servers on ports:",
    "end_message": "Dev servers still running on ports:",
    "separator": ", ",                   // Separator between port numbers
    "show_descriptions": false           // Include port descriptions (verbose)
  },

  "behavior": {
    "silent_failures": true,             // Continue silently on error
    "timeout_seconds": 2,                // Max wait time per port check
    "check_command": "lsof",             // Command to use for checking
    "require_lsof": false                // Fail if lsof not available
  }
}
```

### Configuration Options

#### Ports Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | true | Master switch for all port monitoring |
| `monitored_ports` | array | See below | List of port configs with descriptions |
| `monitored_ports[].number` | string | - | Port number as string |
| `monitored_ports[].description` | string | - | What typically runs on this port |
| `monitored_ports[].enabled` | bool | true | Whether to check this port |
| `custom_ports` | array | [] | Additional port numbers (strings only) |

**Default Monitored Ports**:
- `3000` - React/Next.js default dev server
- `8000` - Django/Python HTTP server
- `8080` - Generic HTTP server (Spring Boot, Tomcat)
- `5173` - Vite dev server
- `4200` - Angular CLI dev server

#### Display Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `show_at_start` | bool | true | Display at session start |
| `show_at_end` | bool | true | Display at session end |
| `icon` | string | "🔌" | Icon prefix for notification |
| `start_message` | string | "Active dev servers on ports:" | Message at start |
| `end_message` | string | "Dev servers still running on ports:" | Message at end |
| `separator` | string | ", " | Separator between port numbers |
| `show_descriptions` | bool | false | Include port descriptions (verbose) |

#### Behavior Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `silent_failures` | bool | true | Continue silently if checks fail |
| `timeout_seconds` | int | 2 | Max wait time per port check |
| `check_command` | string | "lsof" | Command for port checking |
| `require_lsof` | bool | false | Fail if lsof not available |

### Health Scoring

Configuration and execution tracked through health scoring system:

**Configuration Loading** (40 points):
- Config loaded successfully: +20
- Config missing (using defaults): +15
- Config invalid (using defaults): +10

**Port Detection** (40 points):
- All ports checked successfully: +40
- Partial port checks succeed: +20
- All port checks fail: -10

**Display Output** (20 points):
- Output formatted correctly: +20
- Output has minor issues: +10
- Output fails: +0

**Total**: 100 points possible

---

## Integration Patterns

### Basic Integration

```go
package main

import "hooks/lib/session"

func main() {
    // At session start
    session.CheckRunningProcesses()

    // ... session work ...

    // At session end
    session.CheckRunningProcessesAsReminder()
}
```

### Session Start Hook Integration

```go
// In session/cmd-start/start.go
package main

import (
    "hooks/lib/session"
)

func main() {
    // ... other session start logic ...

    // Check for running processes
    session.CheckRunningProcesses()

    // ... continue with session initialization ...
}
```

### Session End Hook Integration

```go
// In session/cmd-end/end.go
package main

import (
    "hooks/lib/session"
)

func main() {
    // ... other session end logic ...

    // Remind about running processes
    session.CheckRunningProcessesAsReminder()

    // ... finish session cleanup ...
}
```

### Conditional Monitoring

Monitoring can be disabled entirely through configuration without code changes:

```jsonc
{
  "ports": {
    "enabled": false  // Disables all monitoring
  }
}
```

Or disable specific contexts:

```jsonc
{
  "display": {
    "show_at_start": false,  // Don't show at start
    "show_at_end": true      // Only show at end
  }
}
```

---

## Troubleshooting

### Problem: No output even when processes are running

**Symptoms**: Port monitoring seems silent, no output displayed

**Possible Causes**:
1. Monitoring disabled in configuration
2. Display disabled for current context
3. Processes not listening on TCP
4. `lsof` not detecting processes

**Diagnostic Steps**:
```bash
# Check if ports.enabled is true
cat ~/.claude/cpi-si/system/data/config/session/processes.jsonc | grep -A5 '"ports"'

# Check display settings
cat ~/.claude/cpi-si/system/data/config/session/processes.jsonc | grep -A10 '"display"'

# Manually check if process listening
lsof -i :3000 -sTCP:LISTEN
```

**Solutions**:
- Set `ports.enabled: true` in configuration
- Set `display.show_at_start: true` or `display.show_at_end: true`
- Verify processes are actually listening: `netstat -tuln | grep :3000`
- Check process is bound to TCP, not just UDP

---

### Problem: "lsof: command not found"

**Symptoms**: Library silently fails to detect any processes

**Expected Behavior**: This is normal - non-blocking design means missing `lsof` doesn't interrupt session

**Diagnostic Steps**:
```bash
# Check if lsof installed
which lsof

# Install lsof if missing (Ubuntu/Debian)
sudo apt install lsof

# Install lsof if missing (macOS)
brew install lsof
```

**Solutions**:
- Install `lsof` package if port monitoring desired
- Library continues working without `lsof` (silently skips detection)
- Set `behavior.require_lsof: false` to explicitly allow missing `lsof`

---

### Problem: Processes not detected on custom ports

**Symptoms**: Custom ports never show in output even when processes running

**Possible Causes**:
1. Custom port not added to configuration
2. Process not actually listening
3. Typo in port number

**Diagnostic Steps**:
```bash
# Check custom_ports in configuration
cat ~/.claude/cpi-si/system/data/config/session/processes.jsonc | grep -A5 'custom_ports'

# Verify process listening
lsof -i :9000 -sTCP:LISTEN  # Replace 9000 with your port

# Check all listening ports
lsof -iTCP -sTCP:LISTEN
```

**Solutions**:
- Add port to `custom_ports` array: `"custom_ports": ["9000"]`
- Verify process is bound to port: `netstat -tuln`
- Ensure port number is string, not number: `"9000"` not `9000`

---

### Problem: Configuration changes not taking effect

**Symptoms**: Changed configuration but behavior unchanged

**Cause**: Configuration loaded once at package import (init function)

**Solution**:
- Restart session to reload configuration
- Package only reads configuration file once when imported
- Changes require session restart to take effect

---

### Problem: Timeout errors or slow detection

**Symptoms**: Port checking takes long time or times out

**Possible Causes**:
1. Timeout too short for system
2. `lsof` running slowly
3. Many ports configured

**Diagnostic Steps**:
```bash
# Test lsof performance
time lsof -i :3000 -sTCP:LISTEN -t

# Check timeout setting
cat ~/.claude/cpi-si/system/data/config/session/processes.jsonc | grep timeout_seconds
```

**Solutions**:
- Increase `behavior.timeout_seconds` to 5 or 10
- Reduce number of monitored ports
- Check system performance (high load may slow `lsof`)

---

## Future Roadmap

### Planned Features

✓ **Configuration-driven port list** - COMPLETED (v2.0.0)
✓ **Timeout protection** - COMPLETED (v2.0.0)
⏳ **Parallel port checking** - Performance improvement for many ports
⏳ **Alternative detection methods** - Use `ss` command as fallback
⏳ **Process name detection** - Show what's running, not just port number
⏳ **Health tracking integration** - Full health scoring implementation

### Research Areas

- **Cross-platform port detection**: Windows and macOS alternatives to `lsof`
- **Process ownership detection**: Show user running each process
- **Automatic cleanup offers**: Interactive prompt to kill forgotten processes
- **Integration with session patterns**: Learn typical port usage patterns
- **Notification system**: Alert for long-running processes across sessions

### Integration Targets

- **Health scoring system**: Track detection success rates over time
- **Session logging**: Record port activity for pattern analysis
- **Pattern learning**: Identify typical development environments
- **Display customization**: Color coding, severity levels for different ports

### Known Limitations

- **TCP only**: UDP processes not currently detected
- **Sequential checking**: Port checks run sequentially (parallel would be faster)
- **lsof dependency**: No fallback if `lsof` unavailable
- **No process identification**: Only shows port numbers, not process names
- **Configuration requires restart**: Changes don't apply mid-session
- **No interactive cleanup**: No option to kill processes from output

### Version History

**2.0.0 (2025-11-12)** - Configuration-driven monitoring
- Added `processes.jsonc` configuration file support
- Implemented timeout protection (2-second default)
- Display customization (icons, messages, separators)
- Graceful fallback to defaults if config missing
- Non-blocking execution guaranteed by timeout

**1.0.0 (2024-10-24)** - Initial implementation
- Hardcoded port list (3000, 8000, 8080, 5173, 4200)
- Basic `lsof` integration
- Dual-context support (start vs end messaging)
- Silent failure design established

---

## Quick Reference

### Default Ports Monitored

| Port | Description |
|------|-------------|
| 3000 | React/Next.js default dev server |
| 8000 | Django/Python HTTP server |
| 8080 | Generic HTTP server (Spring Boot, Tomcat) |
| 5173 | Vite dev server |
| 4200 | Angular CLI dev server |

### Common Commands

```bash
# Check configuration
cat ~/.claude/cpi-si/system/data/config/session/processes.jsonc

# Manually test port detection
lsof -i :3000 -sTCP:LISTEN -t

# See all listening TCP ports
lsof -iTCP -sTCP:LISTEN

# Test with timeout
timeout 2 lsof -i :3000 -sTCP:LISTEN -t

# Alternative port check (if lsof unavailable)
netstat -tuln | grep :3000
```

### Configuration Quick Edits

```bash
# Disable monitoring entirely
# Set: "enabled": false in ports section

# Disable start display only
# Set: "show_at_start": false in display section

# Add custom port
# Add to: "custom_ports": ["9000", "4000"]

# Increase timeout
# Set: "timeout_seconds": 5 in behavior section

# Change separator
# Set: "separator": " | " in display section
```

---

## Related Documentation

- **4-Block Structure**: See `standards/code/4-block/CWS-STD-001-DOC-4-block.md`
- **Configuration Standards**: See `standards/code/4-block/CWS-STD-003-DOC-configuration.md`
- **Health Scoring**: See `knowledge-base/algorithms/CPSI-ALG-001-DOC-health-scoring.md`
- **Session Library Overview**: See `hooks/lib/session/README.md`

---

**For questions, issues, or contributions**:
- Review the modification policy in processes.go CLOSING section
- Follow the 4-block structure pattern
- Test thoroughly before committing (`go build && go vet`)
- Document all changes comprehensively (What/Why/How pattern)
- Verify non-blocking behavior preserved

*"Watch and pray, that ye enter not into temptation" - Matthew 26:41 (KJV)*
