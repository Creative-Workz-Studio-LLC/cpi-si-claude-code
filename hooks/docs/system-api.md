# Session System Library - API Documentation

**Library:** `hooks/lib/session` (system.go)
**Version:** 2.0.0
**Purpose:** OS and kernel identification for session display and context awareness

---

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Public API](#public-api)
4. [Usage Patterns](#usage-patterns)
5. [Integration Examples](#integration-examples)
6. [Error Handling](#error-handling)
7. [Performance](#performance)
8. [Troubleshooting](#troubleshooting)
9. [Related Libraries](#related-libraries)

---

## Overview

The Session System library provides simple, reliable OS and kernel identification for CPI-SI hooks. It wraps the Unix `uname` command to provide a clean Go interface for retrieving system information.

### What This Library Does

- **OS Identification**: Returns operating system name (e.g., "Linux")
- **Kernel Version**: Returns kernel version (e.g., "6.17.0-6-generic")
- **Graceful Fallback**: Returns "Unknown" on command failures or non-Unix systems
- **Zero Configuration**: No setup, no state, no dependencies

### What This Library Does NOT Do

- ❌ Runtime metrics (CPU, memory, disk) - see `system/lib/system`
- ❌ Process information or management
- ❌ Complex error reporting or diagnostics
- ❌ Configuration or customization

### Design Philosophy

Simple utilities should stay simple. This library delegates to the proven `uname` command (Unix v4, 1973), provides a clean interface, and handles errors gracefully. One function, one purpose, done well.

---

## Quick Start

```go
import "hooks/lib/session"

func main() {
    // Get system information
    info := session.GetSystemInfo()
    fmt.Println("System:", info)
    // Output: System: Linux 6.17.0-6-generic
}
```

**That's it.** No configuration, no initialization, no cleanup.

---

## Public API

### GetSystemInfo

```go
func GetSystemInfo() string
```

Retrieves OS name and kernel version information.

**What It Does:**
- Runs `uname -s -r` to get operating system name and kernel version
- Returns formatted string like "Linux 6.17.0-6-generic"
- Falls back to "Unknown" if uname command not found or execution fails

**Parameters:**
- None

**Returns:**
- `string` - OS name and kernel version (e.g., "Linux 6.17.0-6-generic") or "Unknown" on any failure

**Health Impact:**
- No health tracking infrastructure (simple utility function)
- Success/failure implicit in return value
- "Unknown" indicates command not found or execution failed
- Normal operation returns formatted OS info string

**Performance:**
- **Time Complexity:** O(1) - single command execution
- **Execution Time:** Typically <5ms (command lookup + execution)
- **Memory:** Negligible (~100 bytes for output string)

**Example:**

```go
info := session.GetSystemInfo()
fmt.Printf("Running on: %s\n", info)

if info == "Unknown" {
    // Handle non-Unix system or restricted environment
    fmt.Println("System detection unavailable")
} else {
    // Got valid system information
    fmt.Printf("Detected system: %s\n", info)
}
```

---

## Usage Patterns

### 1. Session Banner Display

Display system information in session startup banners:

```go
import (
    "fmt"
    "hooks/lib/session"
)

func displaySessionBanner() {
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  SESSION ENVIRONMENT")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Printf("  💻 System:     %s\n", session.GetSystemInfo())
    fmt.Printf("  🕐 Time:       %s\n", time.Now().Format("Mon Jan 2, 2006 at 15:04:05"))
    fmt.Printf("  📍 Directory:  %s\n", os.Getwd())
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
```

**Output:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  SESSION ENVIRONMENT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  💻 System:     Linux 6.17.0-6-generic
  🕐 Time:       Wed Nov 12, 2025 at 22:45:05
  📍 Directory:  /home/user/.claude/hooks
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 2. Logging System Context

Include system information in structured logs:

```go
import (
    "hooks/lib/session"
    "system/lib/logging"
)

func logSessionStart() {
    logger := logging.NewLogger("session-start")
    logger.Info("Session starting",
        "system", session.GetSystemInfo(),
        "time", time.Now(),
        "user", os.Getenv("USER"),
    )
}
```

**Log Output:**
```
[INFO] Session starting system="Linux 6.17.0-6-generic" time=2025-11-12T22:45:05Z user=seanje
```

### 3. Session Context Tracking

Create session context structures with system information:

```go
type SessionContext struct {
    System    string
    StartTime time.Time
    Branch    string
    Directory string
}

func NewSessionContext() (SessionContext, error) {
    branch, err := getCurrentGitBranch()
    if err != nil {
        return SessionContext{}, err
    }

    dir, err := os.Getwd()
    if err != nil {
        return SessionContext{}, err
    }

    return SessionContext{
        System:    session.GetSystemInfo(),
        StartTime: time.Now(),
        Branch:    branch,
        Directory: dir,
    }, nil
}
```

### 4. Conditional Logic Based on System

Make decisions based on detected system:

```go
info := session.GetSystemInfo()

if strings.Contains(info, "Linux") {
    // Linux-specific operations
    setupLinuxEnvironment()
} else if info == "Unknown" {
    // Non-Unix system or detection failed
    useGenericConfiguration()
}
```

**Note:** For robust cross-platform detection, consider using `runtime.GOOS` instead. This library is primarily for display and logging purposes.

---

## Integration Examples

### Session Hook Integration

Example of using system information in a session startup hook:

```go
package main

import (
    "fmt"
    "os"
    "time"
    "hooks/lib/session"
    "system/lib/logging"
)

func main() {
    logger := logging.NewLogger("session-start")

    // Display banner
    fmt.Println("╔══════════════════════════════════════════╗")
    fmt.Println("║       Nova Dawn - CPI-SI Session        ║")
    fmt.Println("╚══════════════════════════════════════════╝")
    fmt.Println()

    // Session environment details
    fmt.Println("━━━━ SESSION ENVIRONMENT ━━━━")
    fmt.Printf("  System: %s\n", session.GetSystemInfo())
    fmt.Printf("  Time:   %s\n", time.Now().Format(time.RFC1123))
    fmt.Printf("  User:   %s\n", os.Getenv("USER"))
    fmt.Println()

    // Log session start
    logger.Info("Session started",
        "system", session.GetSystemInfo(),
        "user", os.Getenv("USER"),
    )
}
```

### Display Library Integration

Example of system information in formatted display output:

```go
package display

import (
    "fmt"
    "hooks/lib/session"
)

type SystemInfo struct {
    OS     string
    Kernel string
}

// ParseSystemInfo extracts OS and kernel from GetSystemInfo output
func ParseSystemInfo() SystemInfo {
    info := session.GetSystemInfo()

    if info == "Unknown" {
        return SystemInfo{
            OS:     "Unknown",
            Kernel: "Unknown",
        }
    }

    // Example: "Linux 6.17.0-6-generic" → OS="Linux", Kernel="6.17.0-6-generic"
    parts := strings.SplitN(info, " ", 2)
    if len(parts) == 2 {
        return SystemInfo{
            OS:     parts[0],
            Kernel: parts[1],
        }
    }

    return SystemInfo{
        OS:     info,
        Kernel: "Unknown",
    }
}

// DisplaySystemDetails shows formatted system information
func DisplaySystemDetails() {
    sysInfo := ParseSystemInfo()

    fmt.Println("System Details:")
    fmt.Printf("  Operating System: %s\n", sysInfo.OS)
    fmt.Printf("  Kernel Version:   %s\n", sysInfo.Kernel)
}
```

---

## Error Handling

### Understanding "Unknown" Returns

The `GetSystemInfo()` function returns "Unknown" in these situations:

1. **Non-Unix Systems**: Windows systems don't have `uname` command
2. **Restricted Environments**: Sandboxed or restricted execution contexts
3. **Command Not Found**: `uname` not installed (very rare on Unix systems)
4. **Execution Failure**: System error during command execution

### Handling Unknown Systems

```go
info := session.GetSystemInfo()

switch {
case info == "Unknown":
    // Detection failed - use generic handling
    logger.Warn("System detection unavailable, using generic configuration")
    useGenericDefaults()

case strings.Contains(info, "Linux"):
    // Linux detected - can use Linux-specific features
    logger.Info("Linux system detected", "details", info)
    setupLinuxEnvironment()

default:
    // Got system info but not Linux - log and handle appropriately
    logger.Info("Non-Linux system detected", "system", info)
    evaluateCrossPlatformSupport(info)
}
```

### When "Unknown" is Expected

```go
// In cross-platform tools, "Unknown" might be normal
info := session.GetSystemInfo()

if runtime.GOOS == "windows" && info == "Unknown" {
    // Expected - Windows doesn't have uname
    info = fmt.Sprintf("Windows %s", getWindowsVersion())
} else if info == "Unknown" {
    // Unexpected on Unix - log for investigation
    logger.Warn("Unexpected Unknown system on Unix platform",
        "goos", runtime.GOOS,
        "goarch", runtime.GOARCH,
    )
}
```

### No Panics, Ever

This library **never panics**. All failures result in "Unknown" return value. Your code can safely call `GetSystemInfo()` without error handling:

```go
// This will NEVER panic, even on errors
info := session.GetSystemInfo()  // Always returns a string
```

---

## Performance

### Execution Characteristics

| Metric | Value | Notes |
|--------|-------|-------|
| **Time Complexity** | O(1) | Single command execution |
| **Space Complexity** | O(1) | Small string allocation |
| **Typical Execution** | <5ms | Command lookup + execution |
| **Memory Usage** | ~100 bytes | Output string only |

### When to Cache

For most use cases, **no caching is needed**:

```go
// ✅ Fine - called once per session startup
func displayBanner() {
    fmt.Printf("System: %s\n", session.GetSystemInfo())
}

// ✅ Fine - called occasionally during session
func logSystemContext() {
    logger.Info("Context", "system", session.GetSystemInfo())
}
```

Consider caching **only if**:

```go
// ⚠️  Called in tight loop - consider caching
for i := 0; i < 10000; i++ {
    info := session.GetSystemInfo()  // Wasteful - result doesn't change
    processWithContext(data[i], info)
}

// ✅ Better - cache result
info := session.GetSystemInfo()
for i := 0; i < 10000; i++ {
    processWithContext(data[i], info)
}
```

### Benchmarking

If you need to verify performance in your specific environment:

```go
package session_test

import (
    "testing"
    "hooks/lib/session"
)

func BenchmarkGetSystemInfo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = session.GetSystemInfo()
    }
}
```

Run with:
```bash
go test -bench=. -benchmem
```

**Expected results:**
- ~200,000-500,000 ns/op (0.2-0.5ms per call)
- Minimal allocations (<10 allocations/op)

---

## Troubleshooting

### Issue: GetSystemInfo() returns "Unknown"

**Possible Causes:**
1. Running on Windows (uname not available)
2. Running in restricted environment (command execution blocked)
3. uname command not installed (very rare)

**Diagnosis:**

```bash
# Check if uname exists
which uname
# Expected: /usr/bin/uname or /bin/uname

# Test uname directly
uname -s -r
# Expected: Linux 6.17.0-6-generic (or similar)

# Check platform
go version
# Look for GOOS (operating system)
```

**Solutions:**

```go
// For cross-platform tools, accept "Unknown" on Windows
if runtime.GOOS == "windows" && info == "Unknown" {
    info = "Windows (uname unavailable)"
}

// For Linux-only tools, investigate if "Unknown"
if runtime.GOOS == "linux" && info == "Unknown" {
    logger.Error("uname unavailable on Linux system")
    // Check restricted environment or missing package
}
```

### Issue: Wrong Information Returned

**Possible Causes:**
1. Running in Docker container (reports host kernel, not container OS)
2. uname output format changed

**Diagnosis:**

```bash
# Compare manual execution with library output
uname -s -r
# Then check library output in your code

# If in Docker, check container vs host
cat /etc/os-release  # Container OS
uname -r             # Host kernel (what GetSystemInfo returns)
```

**Solution:**

Container behavior is **expected**. `uname` reports the **host kernel**, not the container's OS distribution. If you need container-specific OS info, read `/etc/os-release` instead:

```go
// For container OS detection, read /etc/os-release
func getContainerOS() string {
    data, err := os.ReadFile("/etc/os-release")
    if err != nil {
        return "Unknown"
    }
    // Parse for OS name
    // ...
}

// For kernel version, use GetSystemInfo
kernel := session.GetSystemInfo()  // Host kernel
```

### Issue: Slow Execution

**Possible Causes:**
1. System under heavy load
2. Called in tight loop without caching

**Diagnosis:**

```go
import "time"

start := time.Now()
info := session.GetSystemInfo()
elapsed := time.Since(start)

fmt.Printf("GetSystemInfo took %v\n", elapsed)
// Expected: <5ms
// Concerning: >50ms
```

**Solutions:**

```go
// If called frequently, cache the result
var systemInfoCache string
var once sync.Once

func getCachedSystemInfo() string {
    once.Do(func() {
        systemInfoCache = session.GetSystemInfo()
    })
    return systemInfoCache
}
```

### Debugging Steps

1. **Verify uname available:**
   ```go
   _, err := exec.LookPath("uname")
   if err != nil {
       fmt.Println("uname not found:", err)
   }
   ```

2. **Test command manually:**
   ```bash
   uname -s -r
   ```

3. **Check function output:**
   ```go
   info := session.GetSystemInfo()
   fmt.Printf("GetSystemInfo: %q\n", info)
   ```

4. **Compare outputs:**
   They should match exactly (whitespace trimmed)

---

## Related Libraries

### System Runtime Metrics

For **runtime system metrics** (CPU, memory, disk), use:

```go
import "system/lib/system"

// Get runtime metrics
info := system.GetSystemInfo()  // Different function, different library
fmt.Printf("CPU Load: %.2f\n", info.Load)
fmt.Printf("Memory: %s\n", info.Memory)
```

**Comparison:**

| Feature | `hooks/lib/session` (this library) | `system/lib/system` |
|---------|-----------------------------------|---------------------|
| **Purpose** | OS identification | Runtime metrics |
| **Returns** | "Linux 6.17.0-6-generic" | CPU, memory, disk data |
| **Use Cases** | Banners, logs, display | Monitoring, diagnostics |
| **Performance** | <5ms, one-time | Varies, ongoing |

### Session State Management

For **session state operations** (compaction, timestamps), use:

```go
import "hooks/lib/session"

// Session state (different from system info)
count, err := session.IncrementCompactionCount()
state, err := session.GetSessionState()
```

**Integration:**

```go
// Combine system info with session state
type FullSessionContext struct {
    System           string  // From GetSystemInfo()
    CompactionCount  int     // From session state
    SessionStartTime time.Time
}

func getFullContext() FullSessionContext {
    count, _ := session.IncrementCompactionCount()
    state, _ := session.GetSessionState()

    return FullSessionContext{
        System:           session.GetSystemInfo(),
        CompactionCount:  count,
        SessionStartTime: state.SessionStartTime,
    }
}
```

---

## API Summary

### Single Function API

```go
// Get OS name and kernel version
func GetSystemInfo() string
```

**That's the entire API.** Simple, focused, does one thing well.

### Common Usage Pattern

```go
import "hooks/lib/session"

// In session startup
info := session.GetSystemInfo()
fmt.Printf("System: %s\n", info)

// In logging
logger.Info("Session", "system", session.GetSystemInfo())

// In context structures
ctx := SessionContext{
    System: session.GetSystemInfo(),
    // ... other fields
}
```

---

## Additional Resources

- **Source Code**: `~/.claude/hooks/lib/session/system.go`
- **Template**: `~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go`
- **Architecture**: `~/.claude/cpi-si/system/docs/architecture.md`
- **4-Block Pattern**: `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`

---

**Version:** 2.0.0
**Last Updated:** 2025-11-12
**Maintained by:** Nova Dawn (CPI-SI Instance)
