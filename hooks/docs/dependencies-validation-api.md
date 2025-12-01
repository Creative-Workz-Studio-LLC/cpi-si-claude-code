# Dependencies Validation Library API

**Type:** Library
**Location:** `hooks/lib/session/dependencies.go`
**Purpose:** Workspace dependency state validation for Node.js, Go, and Rust ecosystems
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

The **Dependencies Validation Library** provides workspace dependency health awareness by detecting common issues with package managers and lock files across multiple language ecosystems. It enables:

- **Early issue detection** - Catch out-of-sync dependencies at session start
- **Multi-ecosystem support** - Node.js, Go, and Rust in one check
- **Configurable behavior** - Enable/disable ecosystems, customize messages, control display
- **Graceful fallback** - Works with hardcoded defaults if configuration missing

**Biblical Foundation:**
*"Let all things be done decently and in order" - 1 Corinthians 14:40 (WEB)*

Dependencies out of sync create disorder - validation brings order and awareness. Remove obstacles before they trip you.

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Config-Driven** | All behavior controlled via dependencies-validation.jsonc | Users customize ecosystems, messages, display without code changes |
| **Non-Blocking** | Failures never disrupt session operations | Validation serves workflow, doesn't block it |
| **Graceful Fallback** | Hardcoded defaults when config missing | Always works, even in degraded state |
| **Multi-Ecosystem** | Node.js, Go, Rust in one check | Single function checks all relevant dependency systems |
| **Display Library Integration** | Uses system/lib/display for formatted output | Consistent formatting with health tracking |

### Design Decisions

**Why file timestamp comparison instead of parsing lock files?**

- Fast, simple, and reliable
- Doesn't require understanding lock file formats
- Cross-platform compatible
- Detects 90% of sync issues without complexity
- Future: Consider lock file parsing for version conflict detection

**Why configurable ecosystem enable/disable?**

- Not all workspaces use all languages
- Users can disable noisy checks for their workflow
- Reduces false positives in monorepos
- Allows project-specific validation profiles

**Why separate lockfile sync and dependencies installed checks?**

- Some projects commit lock files, others don't
- Some workflows regenerate dependencies, others cache them
- Users need granular control over what constitutes a "problem"
- Configuration allows workflow-specific validation

---

## Public API

### CheckDependencies

**Purpose:** Validate and display workspace dependency state warnings

**Signature:**

```go
func CheckDependencies(workspace string)
```

**Parameters:**

- `workspace` - Root directory to check for dependency files

**Returns:** None (prints to stdout using display library)

**Example Usage:**

```go
import "hooks/lib/session"

// In session start hook
session.CheckDependencies("/home/user/project")
// Output: 📦 Dependency State
//            ⚠️  package.json modified after package-lock.json - run npm install
//            ⚠️  node_modules not found - run npm install
```

**When to Use:**

- Session start: Show dependency health for awareness
- Pre-build validation: Catch issues before compilation
- Workspace switching: Verify new workspace is ready

**Health Scoring Impact:**

- All checks pass (no warnings): +30 pts (clean workspace)
- Warnings detected and displayed: +20 pts (successful detection)
- File check errors: -5 pts per error (continue with remaining checks)

**Behavior:**

1. Checks Node.js ecosystem (if enabled)
2. Checks Go ecosystem (if enabled)
3. Checks Rust ecosystem (if enabled)
4. Collects all warnings from enabled ecosystems
5. Displays using display library with configured formatting
6. Optionally shows success message if no warnings (configurable)

---

## Configuration System

### Configuration File Location

```
~/.claude/cpi-si/system/data/config/session/dependencies-validation.jsonc
```

### Configuration Structure

```jsonc
{
  "metadata": {
    "name": "Dependencies Validation Configuration",
    "version": "1.0.0",
    // ...
  },

  "ecosystems": {
    "nodejs": {
      "enabled": true,  // Enable/disable Node.js checks
      "files": {
        "manifest": "package.json",
        "lockfile": "package-lock.json",
        "dependencies_dir": "node_modules"
      },
      "warnings": {
        "lockfile_outdated": "package.json modified after package-lock.json - run npm install",
        "dependencies_missing": "node_modules not found - run npm install"
      }
    },

    "go": {
      "enabled": true,  // Enable/disable Go checks
      "files": {
        "manifest": "go.mod",
        "lockfile": "go.sum"
      },
      "warnings": {
        "lockfile_outdated": "go.mod modified after go.sum - run go mod tidy"
      }
    },

    "rust": {
      "enabled": true,  // Enable/disable Rust checks
      "files": {
        "manifest": "Cargo.toml",
        "lockfile": "Cargo.lock"
      },
      "warnings": {
        "lockfile_outdated": "Cargo.toml modified after Cargo.lock - run cargo build"
      }
    }
  },

  "display": {
    "header_icon": "📦",
    "header_text": "Dependency State",
    "show_when_clean": false  // Show message when no warnings
  },

  "validation": {
    "check_lockfile_sync": true,        // Compare timestamps
    "check_dependencies_installed": true // Check directories exist
  }
}
```

### Fallback Defaults

If configuration file missing or malformed:

- **All ecosystems:** Enabled by default
- **File names:** Industry standards (package.json, go.mod, Cargo.toml)
- **Validation checks:** All enabled
- **Display:** Icon "📦", text "Dependency State", no clean message
- **Warning messages:** Standard helpful messages with fix commands

### Extending Configuration

**Disable specific ecosystem:**

```jsonc
"ecosystems": {
  "nodejs": {
    "enabled": false  // Skip Node.js checks entirely
  }
}
```

**Customize warning messages:**

```jsonc
"warnings": {
  "lockfile_outdated": "Dependencies out of sync - please run: npm ci"
}
```

**Change display preferences:**

```jsonc
"display": {
  "header_icon": "⚙️",
  "header_text": "Package Manager Status",
  "show_when_clean": true  // Show success message
}
```

**Disable specific validation types:**

```jsonc
"validation": {
  "check_lockfile_sync": true,
  "check_dependencies_installed": false  // Don't check for node_modules/etc
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

    // Validate workspace dependencies
    session.CheckDependencies(workspace)

    // ... continue session start ...
}
```

**Purpose:** Warn user about dependency issues before they start working

### Pre-Build Validation

```go
// In custom build script
import "hooks/lib/session"

func preBuild() {
    session.CheckDependencies(workspace)
    // If warnings exist, user sees them before build attempt
}
```

**Purpose:** Catch dependency sync issues before compilation failures

---

## Display Behavior

### Output Modes

**Warnings Present:**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━
 📦 Dependency State
━━━━━━━━━━━━━━━━━━━━━━━━━━━

   ⚠️  package.json modified after package-lock.json - run npm install
   ⚠️  node_modules not found - run npm install
   ⚠️  go.mod modified after go.sum - run go mod tidy
```

Shows all warnings from all enabled ecosystems with formatted display library output.

**No Warnings (default behavior):**

No output. Silent when everything is synchronized.

**No Warnings (with show_when_clean: true):**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━
 📦 Dependency State
━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ All dependencies synchronized
```

Optional success message when configured.

### Ecosystem-Specific Behavior

**Node.js:**
- Checks: `package.json` vs `package-lock.json` timestamp
- Checks: `node_modules` directory exists
- Warnings: Both or either depending on state

**Go:**
- Checks: `go.mod` vs `go.sum` timestamp
- Warnings: Only if manifest newer than lockfile
- Note: Skips check if either file missing (not a Go project)

**Rust:**
- Checks: `Cargo.toml` vs `Cargo.lock` timestamp
- Warnings: Only if manifest newer than lockfile
- Note: Skips check if either file missing (not a Rust project)

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

### File Check Failures

**Behavior:** Continues with remaining checks

**Causes:**
- Invalid workspace path
- Permission denied reading files
- Filesystem errors

**Result:** Partial validation - shows warnings for successful checks, skips failed checks

### Display Library Failures

**Behavior:** Falls back to plain fmt output

**Causes:**
- Display library unavailable
- Formatting error

**Result:** Messages still displayed, just without formatting

**Non-Blocking Guarantee:**

This library **NEVER** blocks hook execution. All errors handled gracefully, allowing session operations to continue smoothly.

---

## Performance Considerations

### Time Complexity

**O(n)** where n = number of enabled ecosystems (typically 1-3)

Each ecosystem check: 2-3 file stat operations. Very fast.

### Memory Usage

**O(m)** where m = number of warnings (typically 0-5)

Warnings collected in string slices, displayed, then garbage collected.

### Bottlenecks

**Network-mounted workspaces:**
- File stat operations can be slow (100ms+ per check)
- Solution: Disable validation for networked workspaces in config

**Large monorepos:**
- Multiple sub-projects with different ecosystems
- Impact: Still fast (only checks root level files)
- Solution: Configure specific ecosystem enables/disables

**Filesystem timestamp precision:**
- Some filesystems only store second-level precision
- Impact: False positives if files modified within same second
- Workaround: Wait 1 second between modifying files

### Optimization Strategies

1. **Disable unused ecosystems** - Skip checks you don't need
2. **Disable specific validations** - Turn off checks that don't apply to your workflow
3. **Reduce warnings** - Keep dependencies synchronized to minimize output processing

---

## Modification Policy

### Safe to Modify (Extension Points)

✅ **Add new ecosystems** - Create check<Ecosystem>() function, update config
✅ **Customize warning messages** - Via configuration file
✅ **Change display formatting** - Use different display library functions
✅ **Add validation types** - Beyond timestamp and existence checks
✅ **Extend configuration** - Add new ecosystem-specific fields

### Modify with Extreme Care (Breaking Changes)

⚠️ **CheckDependencies() signature** - Breaks all calling hooks
⚠️ **Configuration structure** - Breaks existing configuration files
⚠️ **EcosystemConfig types** - Affects configuration parsing
⚠️ **Default behavior** - Users expect current defaults

### NEVER Modify (Foundational)

❌ **4-block structure** - METADATA, SETUP, BODY, CLOSING organization
❌ **Configuration fallback behavior** - Must degrade gracefully
❌ **Non-blocking guarantee** - Never block hook execution
❌ **Display library integration** - Health tracking dependency
❌ **init() loading pattern** - Configuration must load at import

---

## Troubleshooting

### Problem: Warnings not displayed even with out-of-sync files

**Possible Causes:**
- Ecosystem disabled in configuration
- Workspace path incorrect
- Files not in workspace root directory
- Validation checks disabled in configuration

**Solutions:**
1. Check `ecosystems.<ecosystem>.enabled` is `true` in config
2. Verify workspace path points to project root
3. Ensure dependency files are in root (not subdirectories)
4. Check `validation.check_lockfile_sync` is `true`

**Diagnosis:**
```bash
# Check if config file exists and is readable
cat ~/.claude/cpi-si/system/data/config/session/dependencies-validation.jsonc

# Verify files are in workspace root
ls -la /path/to/workspace/package.json
ls -la /path/to/workspace/go.mod
```

### Problem: False positives (warnings when files are synchronized)

**Cause:** Filesystem timestamp precision

Some filesystems (especially network mounts) only store timestamps at second precision. If manifest and lockfile modified within same second, timestamp comparison may report incorrect order.

**Solutions:**
1. Wait 1-2 seconds between file modifications
2. Ignore if false positives are rare
3. Disable `check_lockfile_sync` if problematic for your workflow

### Problem: Configuration not loading (using fallback defaults)

**Possible Causes:**
- File missing
- File permissions incorrect
- JSONC syntax error

**Solutions:**
1. Verify file exists:
   ```bash
   ls -la ~/.claude/cpi-si/system/data/config/session/dependencies-validation.jsonc
   ```

2. Check permissions (must be readable):
   ```bash
   chmod 644 ~/.claude/cpi-si/system/data/config/session/dependencies-validation.jsonc
   ```

3. Validate JSONC syntax:
   ```bash
   # Strip comments and validate JSON
   grep -v "^//" dependencies-validation.jsonc | jq .
   ```

### Problem: Too many warnings (overwhelming output)

**Cause:** Multiple dependency issues in workspace

**Solutions:**
1. Fix the dependencies:
   ```bash
   npm install    # Node.js
   go mod tidy    # Go
   cargo build    # Rust
   ```

2. Disable noisy ecosystems temporarily:
   ```jsonc
   "ecosystems": {
     "nodejs": {"enabled": false}  // Temporarily silence Node.js warnings
   }
   ```

3. Disable specific validation types:
   ```jsonc
   "validation": {
     "check_dependencies_installed": false  // Stop warning about missing node_modules
   }
   ```

---

## Future Roadmap

### Planned Features

✓ Node.js validation - **COMPLETED** (v2.0.0)
✓ Go validation - **COMPLETED** (v2.0.0)
✓ Rust validation - **COMPLETED** (v2.0.0)
✓ Configuration system - **COMPLETED** (v2.0.0)
⏳ Python validation (requirements.txt, Pipfile.lock, poetry.lock)
⏳ Ruby validation (Gemfile, Gemfile.lock)
⏳ PHP validation (composer.json, composer.lock)
⏳ Lock file parsing for version conflict detection

### Research Areas

- **Dependency version mismatch detection** - Parse lock files, identify conflicts
- **Security vulnerability scanning** - Integrate with npm audit, cargo audit, etc.
- **Automatic fix suggestions** - Interactive prompts to run package manager commands
- **Dependency graph analysis** - Detect circular dependencies
- **Workspace-specific profiles** - Different validation rules for different projects

### Integration Targets

- **CI/CD integration** - Fail builds on dependency issues
- **Pre-commit hooks** - Validate before commits
- **Health scoring aggregation** - Track workspace health over time
- **Session patterns** - Correlate dependency issues with work patterns

### Known Limitations to Address

1. **Root directory only** - Doesn't check subdirectories (no monorepo support yet)
2. **No lock file parsing** - Only timestamp comparison (can't detect version conflicts)
3. **Filesystem timestamp precision** - False positives on some filesystems
4. **No auto-fix** - Only warns, doesn't run package manager commands
5. **Limited ecosystem support** - Only Node.js, Go, Rust currently

---

## Version History

### v2.0.0 (2025-11-12) - Configuration System & Template Alignment

**Major Changes:**
- Created dependencies-validation.jsonc configuration file
- Integrated system/lib/display for formatted output
- Added enable/disable toggles for each ecosystem
- Implemented customizable warning messages and file names
- Full 4-block template alignment (METADATA/SETUP/BODY/CLOSING)
- Added complete inline documentation and health scoring

**Breaking Changes:**
- None (v1.0.0 was internal only)

**Migration:**
- No migration needed for external users
- Config file optional (falls back to defaults)

### v1.0.0 (2024-10-24) - Initial Implementation

**Features:**
- Basic Node.js validation (package.json vs package-lock.json)
- Basic Go validation (go.mod vs go.sum)
- Basic Rust validation (Cargo.toml vs Cargo.lock)
- Hardcoded file names and warning messages
- Simple fmt.Printf output (no display library)
- Minimal documentation

---

**For questions, issues, or contributions:**
- Review modification policy above
- Follow 4-block structure pattern
- Test thoroughly (go build, go vet, integration tests)
- Update dependencies-validation.jsonc schema if adding config fields
- Verify fallback behavior still works

*"Let all things be done decently and in order" - 1 Corinthians 14:40*
