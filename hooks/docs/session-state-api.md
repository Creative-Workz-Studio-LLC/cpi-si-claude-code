# Session State Library - API Documentation

**Package:** `hooks/lib/session`
**Version:** 3.0.0
**Last Updated:** 2025-11-12

---

## Overview

The Session State Library provides hooks-layer access to session state operations for the CPI-SI Hooks System. This is a thin orchestration wrapper that delegates all operations to `system/lib/sessiontime` (the authoritative implementation) while maintaining a backward-compatible interface for existing hooks.

**Key Characteristics:**
- Pure delegation wrapper - no business logic
- Zero duplication - all implementation in `system/lib/sessiontime`
- Backward-compatible type exports for hooks ecosystem
- Stateless - no cleanup required

---

## Installation

### Import

```go
import "hooks/lib/session"
```

### Requirements

- Go 1.21 or later
- `system/lib/sessiontime` must be accessible in module path
- Proper go.mod configuration with replace directives

### Build Verification

```bash
cd ~/.claude/hooks/lib/session
go build .
go vet .
```

---

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "hooks/lib/session"
)

func main() {
    // Get current session state
    state, err := session.GetSessionState()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Session ID: %s\n", state.SessionID)
    fmt.Printf("Started: %s\n", state.StartFormatted)
    fmt.Printf("Compactions: %d\n", state.CompactionCount)

    // Increment compaction count
    count, err := session.IncrementCompactionCount()
    if err != nil {
        log.Printf("Failed to increment: %v", err)
        return
    }

    fmt.Printf("New compaction count: %d\n", count)
}
```

---

## Public API Reference

### Types

#### SessionState

```go
type SessionState = sessiontime.SessionState
```

Re-exported from `system/lib/sessiontime` for backward compatibility.

**Purpose:** Represents complete session state with timing, configuration, and tracking information.

**Important:** This is a type alias, not a duplicate definition. Changes to `SessionState` happen in `system/lib/sessiontime` and automatically propagate to this wrapper.

**Fields:** See `system/lib/sessiontime` documentation for complete field reference.

**Common Fields:**
- `SessionID` (string) - Unique session identifier
- `StartFormatted` (string) - Human-readable start time
- `CompactionCount` (int) - Number of compactions this session
- `CircadianPhase` (string) - Time of day phase (morning/afternoon/evening/night)

---

### Functions

#### IncrementCompactionCount

```go
func IncrementCompactionCount() (int, error)
```

Increments the compaction count and returns the new value.

**What It Does:**
Delegates to `system/lib/sessiontime.IncrementCompactionCount()` which reads current session state, increments the `CompactionCount` field, writes updated state back to file, and returns the new count.

**Parameters:** None

**Returns:**
- `int` - New compaction count after increment
- `error` - Error from system library (file read/write/JSON errors)

**Health Impact:**
- Delegation success: +80 points
- Expected result: +20 points
- System library failure: -50 points

**Example:**

```go
count, err := session.IncrementCompactionCount()
if err != nil {
    log.Printf("Failed to increment compaction count: %v", err)
    return
}
fmt.Printf("Compaction count: %d\n", count)
```

**Typical Usage:** Called by hooks that trigger on compaction events (e.g., `session/cmd-pre-compact`).

---

#### GetCompactionCount

```go
func GetCompactionCount() (int, error)
```

Returns the current compaction count from session state.

**What It Does:**
Delegates to `system/lib/sessiontime.GetCompactionCount()` which reads current session state and extracts the `CompactionCount` field.

**Parameters:** None

**Returns:**
- `int` - Current compaction count
- `error` - Error from system library (file read/JSON parse errors)

**Health Impact:**
- Delegation success: +80 points
- Expected result: +20 points
- System library failure: -50 points

**Example:**

```go
count, err := session.GetCompactionCount()
if err != nil {
    log.Printf("Failed to get compaction count: %v", err)
    return
}
fmt.Printf("Current compaction count: %d\n", count)
```

**Typical Usage:** Query compaction count without modifying it.

---

#### GetSessionState

```go
func GetSessionState() (*SessionState, error)
```

Returns the complete current session state.

**What It Does:**
Delegates to `system/lib/sessiontime.ReadSession()` which reads the session state file, parses JSON, and returns the complete `SessionState` struct.

**Parameters:** None

**Returns:**
- `*SessionState` - Pointer to current session state with all fields
- `error` - Error from system library (file read/JSON parse errors)

**Health Impact:**
- Delegation success: +80 points
- Expected result: +20 points
- System library failure: -50 points

**Example:**

```go
state, err := session.GetSessionState()
if err != nil {
    log.Printf("Failed to get session state: %v", err)
    return
}

fmt.Printf("Session ID: %s\n", state.SessionID)
fmt.Printf("Started: %s\n", state.StartFormatted)
fmt.Printf("Compactions: %d\n", state.CompactionCount)
fmt.Printf("Circadian Phase: %s\n", state.CircadianPhase)
```

**Typical Usage:** Access complete session information for hooks that need full context.

---

## Usage Patterns

### Pattern 1: Compaction Tracking

Track compactions in pre-compaction hooks:

```go
package main

import (
    "log"
    "hooks/lib/session"
)

func main() {
    // Increment compaction count
    count, err := session.IncrementCompactionCount()
    if err != nil {
        // Log error but don't fail hook
        log.Printf("Warning: Failed to track compaction: %v", err)
        return
    }

    log.Printf("Compaction #%d\n", count)
    // Continue with compaction operations...
}
```

### Pattern 2: Session Context Access

Access session information in any hook:

```go
package main

import (
    "fmt"
    "log"
    "hooks/lib/session"
)

func main() {
    state, err := session.GetSessionState()
    if err != nil {
        log.Printf("Warning: Could not get session state: %v", err)
        return
    }

    // Use session context for decisions
    if state.CircadianPhase == "night" {
        fmt.Println("Late night session detected")
        // Adjust behavior for late-night work
    }

    fmt.Printf("Session %s started at %s\n",
        state.SessionID,
        state.StartFormatted)
}
```

### Pattern 3: Read-Only Queries

Query session information without modification:

```go
package main

import (
    "fmt"
    "hooks/lib/session"
)

func main() {
    // Just read compaction count
    count, err := session.GetCompactionCount()
    if err != nil {
        // Handle error appropriately
        return
    }

    if count > 5 {
        fmt.Println("Warning: Multiple compactions this session")
        // Notify user or adjust behavior
    }
}
```

---

## Error Handling

All functions return errors that originate from `system/lib/sessiontime`. The wrapper propagates these errors unchanged.

### Common Errors

**File Not Found:**
```go
count, err := session.GetCompactionCount()
if err != nil {
    // Session state file doesn't exist yet
    // This can happen early in session lifecycle
    log.Printf("Session state not available: %v", err)
}
```

**JSON Parse Errors:**
```go
state, err := session.GetSessionState()
if err != nil {
    // Session state file corrupted or invalid JSON
    log.Printf("Failed to parse session state: %v", err)
}
```

**File Write Errors:**
```go
count, err := session.IncrementCompactionCount()
if err != nil {
    // Permission issues or disk full
    log.Printf("Failed to update session state: %v", err)
}
```

### Error Handling Best Practices

**In Hooks:**
```go
// Non-critical operations - log but don't fail hook
count, err := session.IncrementCompactionCount()
if err != nil {
    log.Printf("Warning: %v", err)
    return  // Exit gracefully, don't crash hook
}
```

**In Critical Operations:**
```go
// Critical operations - propagate error
state, err := session.GetSessionState()
if err != nil {
    return fmt.Errorf("failed to get session state: %w", err)
}
```

---

## Performance Considerations

### Wrapper Overhead

- **Function call overhead:** Negligible (single delegation)
- **Memory overhead:** None (stateless wrapper)
- **Type alias overhead:** Zero (compile-time only)

### Actual Performance

All actual performance characteristics are in `system/lib/sessiontime`:
- File I/O happens in system library
- JSON parsing happens in system library
- See `system/lib/sessiontime` documentation for operation costs

### Caching

This wrapper does **not** cache results. Every call hits the system library, which reads from disk. If you need to access session state multiple times in rapid succession, store the result:

```go
// Inefficient - reads file 3 times
state1, _ := session.GetSessionState()
fmt.Println(state1.SessionID)
state2, _ := session.GetSessionState()
fmt.Println(state2.StartFormatted)
state3, _ := session.GetSessionState()
fmt.Println(state3.CompactionCount)

// Efficient - reads file once
state, err := session.GetSessionState()
if err != nil {
    log.Fatal(err)
}
fmt.Println(state.SessionID)
fmt.Println(state.StartFormatted)
fmt.Println(state.CompactionCount)
```

---

## Troubleshooting

### Problem: `undefined: sessiontime`

**Cause:** `system/lib/sessiontime` not in module path

**Solution:** Verify go.mod has correct replace directive

```bash
# Check system module configuration
cat ~/.claude/cpi-si/system/go.mod
```

### Problem: `SessionState` type not found

**Cause:** Import path incorrect or system library not built

**Solution:** Verify system library is accessible

```bash
cd ~/.claude/cpi-si/system
go build system/lib/sessiontime
```

### Problem: Wrapper functions return errors

**Cause:** System library operations failing

**Solution:** See `system/lib/sessiontime` troubleshooting

**Note:** Wrapper propagates errors unchanged - diagnose in system library

---

## Integration Examples

### Session Start Hook

```go
// hooks/session/cmd-start
package main

import (
    "log"
    "hooks/lib/session"
)

func main() {
    state, err := session.GetSessionState()
    if err != nil {
        log.Printf("Warning: Could not get session state: %v", err)
        return
    }

    log.Printf("Session %s started\n", state.SessionID)
    log.Printf("Circadian phase: %s\n", state.CircadianPhase)
}
```

### Pre-Compaction Hook

```go
// hooks/session/cmd-pre-compact
package main

import (
    "log"
    "hooks/lib/session"
)

func main() {
    count, err := session.IncrementCompactionCount()
    if err != nil {
        log.Printf("Warning: Failed to track compaction: %v", err)
        return
    }

    log.Printf("Compaction #%d beginning\n", count)
}
```

### Session End Hook

```go
// hooks/session/cmd-end
package main

import (
    "fmt"
    "hooks/lib/session"
)

func main() {
    state, err := session.GetSessionState()
    if err != nil {
        return  // Session may already be cleaned up
    }

    fmt.Printf("Session %s ended\n", state.SessionID)
    fmt.Printf("Total compactions: %d\n", state.CompactionCount)
}
```

---

## Testing

### Basic Functionality Test

```go
package main

import (
    "testing"
    "hooks/lib/session"
)

func TestBasicFunctionality(t *testing.T) {
    // Test GetSessionState
    state, err := session.GetSessionState()
    if err != nil {
        t.Fatalf("GetSessionState failed: %v", err)
    }

    if state.SessionID == "" {
        t.Error("SessionID should not be empty")
    }

    // Test GetCompactionCount
    count1, err := session.GetCompactionCount()
    if err != nil {
        t.Fatalf("GetCompactionCount failed: %v", err)
    }

    // Test IncrementCompactionCount
    count2, err := session.IncrementCompactionCount()
    if err != nil {
        t.Fatalf("IncrementCompactionCount failed: %v", err)
    }

    if count2 != count1+1 {
        t.Errorf("Expected count %d, got %d", count1+1, count2)
    }
}
```

### Integration Test with Hooks

```bash
# Test with actual hooks
cd ~/.claude/hooks/session
./cmd-start   # Should access session state
./cmd-pre-compact   # Should increment compaction count
./cmd-end     # Should read final state
```

---

## Related Documentation

- **Authoritative Implementation:** `system/lib/sessiontime` - Where actual operations happen
- **Command Wrapper:** `system/runtime/cmd/session-time` - CLI access to session state
- **Architecture:** `~/.claude/cpi-si/system/docs/architecture.md` - Ladder/Baton/Rails model
- **4-Block Structure:** `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`

---

## Version History

### 3.0.0 (2025-11-12)
- Architectural consolidation - removed duplication
- Type re-export for backward compatibility
- Pure delegation pattern established
- Zero business logic in wrapper

### 2.0.0 (2025-11-10)
- Config inheritance from user/instance/project
- Richer SessionState structure
- Correct file path resolution

### 1.0.0 (2024-10-24)
- Initial implementation
- Basic session state operations
- Simple compaction tracking

---

## Contributing

### Adding New Wrapper Functions

1. Add function in source file's BODY "Public APIs" section
2. Follow delegation pattern: `return sessiontime.FunctionName(...)`
3. Document with full docstring (What/Parameters/Returns/Health/Example)
4. Update BODY Organizational Chart
5. Update METADATA Public API list
6. **Update this API documentation**

### Modification Guidelines

**Safe to modify:**
- ✅ Add new wrapper functions (follow delegation pattern)
- ✅ Add convenience methods that delegate to system library

**Modify with care:**
- ⚠️ Function signatures (breaks all calling hooks)
- ⚠️ SessionState type alias (breaks hooks accessing state)

**Never modify:**
- ❌ Delegation pattern (must always delegate to system library)
- ❌ SessionState definition (lives in `system/lib/sessiontime` only)

---

## Support

For questions, issues, or contributions:
- Review modification policy in source code
- Follow 4-block structure pattern
- Test with actual hooks before committing
- Document all changes comprehensively

---

*"Be watchful, stand firm in the faith" - 1 Corinthians 16:13 (WEB)*
