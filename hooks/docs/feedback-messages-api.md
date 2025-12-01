# Feedback Messages Library API

**Type:** Library
**Location:** `hooks/lib/feedback/messages.go`
**Purpose:** Contextual guidance and encouragement after tool operations
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

The **Feedback Messages Library** provides Kingdom-honoring contextual guidance after tool operations. It encourages quality standards, meaningful development practices, and excellence that reflects the Creator.

**Biblical Foundation:**
*"Let your speech always be with grace, seasoned with salt, that you may know how you ought to answer each one." - Colossians 4:6 (WEB)*

Feedback should edify, encourage, and guide toward excellence - not condemn or burden.

**Key Capabilities:**

- Post-commit encouragement (meaningful commit messages)
- Dependency install reminders (version verification)
- Build completion quality checklist (technical + Kingdom standards)
- Command failure warnings (non-zero exit codes)
- Build command detection (pattern matching across languages)

**Version 2.0.0 Features:**

- **Configuration-driven:** Messages, patterns, and checklists loaded from system config files
- **Graceful fallback:** Works with or without configuration (hardcoded defaults)
- **System library integration:** Delegates formatting to system/lib/display
- **Extensible:** Users can customize without modifying code

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Gracious Communication** | Encouraging tone, not condemning | Builds up rather than tears down |
| **Non-Blocking** | Failures never interrupt workflow | Feedback serves, doesn't hinder |
| **Configuration-Driven** | Load from system config files | Users customize without code changes |
| **Graceful Fallback** | Hardcoded defaults if config unavailable | Always functional, even without setup |
| **Kingdom-Honoring** | Excellence that reflects the Creator | Work demonstrates biblical principles |

### Design Decisions

**Why configuration-driven instead of hardcoded?**

- Users have different build tools (bazel, cmake, custom scripts)
- Messages can be personalized per-user or per-instance
- Quality standards vary by project and team
- Supports future internationalization
- Aligns with CPI-SI philosophy: configuration over hardcoding

**Why graceful fallback?**

- Works out-of-the-box without configuration
- No breaking changes for existing code
- Progressive enhancement: better with config, functional without
- Resilient to config file errors or missing files

**Why delegate to system/lib/display?**

- Don't duplicate formatting logic (extract and orchestrate)
- System library provides ANSI colors, icons, panic recovery
- Health tracking handled by lower-rung library
- Architectural integrity: hooks (top) → feedback (mid) → display (bottom)

---

## Public API

### OnCommitCreated

**Purpose:** Provides feedback after git commit operations

**Signature:**

```go
func OnCommitCreated()
```

**Behavior:**

- Loads commit success message from config (or uses fallback)
- Displays with green checkmark icon (via display.Success)
- Encourages meaningful commit messages

**Example Usage:**

```go
import "hooks/lib/feedback"

// After git commit completes
feedback.OnCommitCreated()
// Output: ✓ Commit created. Commits tell a story - make it meaningful.
```

**Configuration:**

- Config path: `system/data/config/feedback/messages.jsonc`
- Key: `commit.success.message`
- Fallback: "Commit created. Commits tell a story - make it meaningful."

**When to Call:**

- tool/post-use hook after git commit
- Any hook detecting successful commit operation

---

### OnDependenciesInstalled

**Purpose:** Reminds to verify versions after package installation

**Signature:**

```go
func OnDependenciesInstalled()
```

**Behavior:**

- Loads dependency installed message from config (or uses fallback)
- Displays with green checkmark icon (via display.Success)
- Encourages version verification for reproducible builds

**Example Usage:**

```go
// After npm install, pip install, cargo add, etc.
feedback.OnDependenciesInstalled()
// Output: ✓ Dependencies installed. Consider verifying versions.
```

**Configuration:**

- Config path: `system/data/config/feedback/messages.jsonc`
- Key: `dependencies.installed.message`
- Fallback: "Dependencies installed. Consider verifying versions."

**When to Call:**

- After npm install, yarn add, pnpm install
- After pip install, poetry add
- After cargo add, cargo install
- Any package installation operation

---

### OnBuildComplete

**Purpose:** Shows quality checklist after successful build

**Signature:**

```go
func OnBuildComplete()
```

**Behavior:**

- Loads build header from config (or uses fallback)
- Displays header with cyan info icon (via display.Info)
- Loads quality checklist from config with category filtering
- Falls back to hardcoded 3-item checklist if config unavailable
- Supports multiple display modes (default, comprehensive, quick, kingdom_focused)

**Example Usage:**

```go
// After successful build
feedback.OnBuildComplete()
// Output:
// ℹ Build complete. Quality check:
//    • Compiles without warnings?
//    • Edge cases tested?
//    • Honors God through excellence?
```

**Configuration:**

- Header: `system/data/config/feedback/messages.jsonc` → `build.complete.header`
- Checklist: `system/data/config/feedback/quality-checklist.jsonc`
- Active mode: `filtering.active_mode` (default: "default")
- Categories shown depend on mode's `include_categories`

**Default Mode Categories:**

- `compilation` - Technical compilation quality
- `testing` - Test coverage and edge cases
- `principle` - Kingdom-honoring excellence

**Comprehensive Mode:**

- All default categories plus:
- `robustness` - Error handling
- `readability` - Code clarity
- `architecture` - Design quality

**When to Call:**

- After cargo build, go build, npm run build
- After gcc, g++, clang compilation
- Use IsBuildCommand() to detect build operations

---

### OnCommandFailure

**Purpose:** Warns about non-zero exit codes

**Signature:**

```go
func OnCommandFailure(exitCode string)
```

**Parameters:**

- `exitCode` - Exit code as string (e.g., "1", "127")

**Behavior:**

- Checks if exit code is non-zero (skips "0" and empty)
- Loads message template from config (or uses fallback)
- Replaces `{exit_code}` placeholder with actual value
- Displays with yellow warning icon (via display.Warning)

**Example Usage:**

```go
// After command execution
exitCode := "1"
feedback.OnCommandFailure(exitCode)
// Output: ⚠ Command exited with code 1
```

**Configuration:**

- Config path: `system/data/config/feedback/messages.jsonc`
- Key: `command_failure.generic.message`
- Template: "Command exited with code {exit_code}"
- Supports variable substitution

**When to Call:**

- tool/post-use hook after any command execution
- When command returns non-zero exit code
- Pass empty string or "0" to skip warning (no output)

---

### IsBuildCommand

**Purpose:** Detects if command string is a build operation

**Signature:**

```go
func IsBuildCommand(cmd string) bool
```

**Parameters:**

- `cmd` - Full command string to check

**Returns:**

- `bool` - true if command matches any build pattern, false otherwise

**Behavior:**

- Loads build patterns from config (or uses fallback)
- Iterates through all patterns across language categories
- Uses simple substring matching (case-sensitive)
- Returns immediately on first match

**Example Usage:**

```go
commands := []string{
    "cargo build --release",  // true - Rust build
    "go build main.go",        // true - Go build
    "npm run build",           // true - Node.js build
    "make test",               // false - not a build command
    "gcc -o program main.c",   // true - C compilation
}

for _, cmd := range commands {
    if feedback.IsBuildCommand(cmd) {
        feedback.OnBuildComplete()
    }
}
```

**Configuration:**

- Config path: `system/data/config/feedback/build-patterns.jsonc`
- Patterns organized by language/tool categories
- 40+ patterns across 9 categories (rust, go, c_cpp, javascript, typescript, python, android, ios, generic)

**Fallback Patterns:**

- cargo build, go build, npm run build, make build, gcc, g++, clang

**When to Use:**

- tool/post-use hook to conditionally show build feedback
- Command categorization and routing
- Build-specific processing logic

---

## Configuration System

### Configuration Files Location

**Base directory:** `~/.claude/cpi-si/system/data/config/feedback/`

**Files:**

1. `build-patterns.jsonc` - Build command detection patterns
2. `messages.jsonc` - Feedback message templates
3. `quality-checklist.jsonc` - Quality checklist items and filtering

### Build Patterns Configuration

**File:** `build-patterns.jsonc`

**Structure:**

```jsonc
{
  "patterns": {
    "rust": {
      "patterns": ["cargo build", "cargo run", "rustc"],
      "description": "Rust build system and compiler"
    },
    "go": {
      "patterns": ["go build", "go install"],
      "description": "Go compiler and build tool"
    }
    // ... more language categories
  },
  "config": {
    "match_mode": "substring",
    "case_sensitive": true
  },
  "extensions": {
    // Add custom build patterns here
  }
}
```

**Supported Categories:**

- `rust` - Rust build system (cargo, rustc)
- `go` - Go compiler
- `c_cpp` - C/C++ compilers (gcc, g++, clang, make, cmake)
- `javascript` - JS bundlers (webpack, vite, rollup, npm/yarn/pnpm build)
- `typescript` - TypeScript compiler (tsc)
- `python` - Python build tools
- `android` - Android build (gradle)
- `ios` - iOS/macOS build (xcodebuild, swift build)
- `generic` - Language-agnostic (make, ninja, bazel, buck, meson)

**Adding Custom Patterns:**

```jsonc
"extensions": {
  "my_tool": {
    "patterns": ["mytool build", "mytool compile"],
    "description": "My custom build system"
  }
}
```

### Messages Configuration

**File:** `messages.jsonc`

**Structure:**

```jsonc
{
  "commit": {
    "success": {
      "message": "Commit created. Commits tell a story - make it meaningful.",
      "display_type": "success"
    }
  },
  "build": {
    "complete": {
      "header": "Build complete. Quality check:",
      "display_type": "info"
    }
  },
  "command_failure": {
    "generic": {
      "message": "Command exited with code {exit_code}",
      "display_type": "warning"
    }
  }
}
```

**Variable Substitution:**

- Use `{variable_name}` in message strings
- Available variables: `{exit_code}`, `{command}`, `{file_path}`, `{duration}`, `{error_message}`
- Replacement happens at runtime

**Display Types:**

- `success` - Green checkmark icon
- `warning` - Yellow warning icon
- `info` - Cyan info icon
- `failure` - Red X icon

### Quality Checklist Configuration

**File:** `quality-checklist.jsonc`

**Structure:**

```jsonc
{
  "checklist": {
    "technical": [
      {
        "item": "Compiles without warnings?",
        "category": "compilation",
        "priority": "high",
        "explanation": "Warnings indicate potential issues"
      }
    ],
    "kingdom": [
      {
        "item": "Honors God through excellence?",
        "category": "principle",
        "priority": "essential",
        "explanation": "Work reflects the Creator"
      }
    ]
  },
  "filtering": {
    "modes": {
      "default": {
        "include_categories": ["compilation", "testing", "principle"]
      },
      "comprehensive": {
        "include_categories": ["compilation", "testing", "robustness", "readability", "architecture", "principle"]
      },
      "quick": {
        "include_categories": ["compilation", "principle"]
      }
    },
    "active_mode": "default"
  }
}
```

**Filtering Modes:**

- `default` - 3-item checklist (matches original behavior)
- `comprehensive` - Extended checklist for thorough review
- `quick` - Minimal checklist for rapid iteration
- `kingdom_focused` - Focus on Kingdom principles only

**Changing Active Mode:**

```jsonc
"filtering": {
  "active_mode": "comprehensive"  // Change to desired mode
}
```

**Adding Custom Checklist Items:**

```jsonc
"checklist": {
  "technical": [
    {
      "item": "Your custom check?",
      "category": "custom",
      "priority": "medium",
      "explanation": "What this checks for"
    }
  ]
}
```

---

## Integration with Hooks

### Hook Integration Pattern

Hooks call feedback functions at appropriate points in their execution:

```go
package toolpostuse

import (
    "hooks/lib/feedback"
    "strings"
)

func ExecutePostToolUse(toolName string, target string, exitCode string) error {
    // Check if this was a git commit
    if strings.Contains(toolName, "git commit") {
        feedback.OnCommitCreated()
    }

    // Check if this was a build command
    if feedback.IsBuildCommand(toolName) {
        feedback.OnBuildComplete()
    }

    // Check if this was a package install
    if strings.Contains(toolName, "npm install") ||
       strings.Contains(toolName, "pip install") {
        feedback.OnDependenciesInstalled()
    }

    // Warn on command failure
    feedback.OnCommandFailure(exitCode)

    return nil
}
```

### Current Hook Integration

| Hook | Feedback Function | Trigger Condition |
|------|-------------------|-------------------|
| **tool/post-use** | OnCommitCreated() | After git commit |
| **tool/post-use** | OnBuildComplete() | After build commands (IsBuildCommand) |
| **tool/post-use** | OnDependenciesInstalled() | After package install |
| **tool/post-use** | OnCommandFailure() | On non-zero exit code |

### Conditional Feedback Example

```go
// Only show feedback for build commands
cmd := "go build main.go"
exitCode := runCommand(cmd)

if feedback.IsBuildCommand(cmd) {
    if exitCode == 0 {
        feedback.OnBuildComplete()
    } else {
        feedback.OnCommandFailure(fmt.Sprintf("%d", exitCode))
    }
}
```

---

## Extending the System

### Adding New Feedback Functions

**Step 1:** Add message to configuration

```jsonc
// In messages.jsonc
"test": {
  "passed": {
    "message": "Tests passed. Well done!",
    "display_type": "success"
  }
}
```

**Step 2:** Add function to messages.go

```go
// OnTestPassed provides feedback after successful test run
func OnTestPassed() {
    message := "Tests passed. Well done!"  // Default fallback
    if configLoaded && messages != nil {
        message = messages.Test.Passed.Message
    }
    fmt.Println(display.Success(message))
}
```

**Step 3:** Call from hooks

```go
// In tool/post-use hook
if strings.Contains(toolName, "test") && exitCode == "0" {
    feedback.OnTestPassed()
}
```

### Adding New Build Patterns

**Option 1:** Use extensions in config (no code changes)

```jsonc
// In build-patterns.jsonc
"extensions": {
  "my_tool": {
    "patterns": ["mytool build", "mytool compile"],
    "description": "My custom build tool"
  }
}
```

**Option 2:** Add to existing category

```jsonc
// In build-patterns.jsonc
"generic": {
  "patterns": [
    "make",
    "ninja",
    "bazel build",
    "my-build-script"  // Add here
  ]
}
```

### Adding Custom Quality Checklist Items

**Add to existing category:**

```jsonc
// In quality-checklist.jsonc
"checklist": {
  "technical": [
    {
      "item": "Compiles without warnings?",
      "category": "compilation",
      "priority": "high"
    },
    {
      "item": "Your custom check?",
      "category": "custom",
      "priority": "medium",
      "explanation": "What this checks"
    }
  ]
}
```

**Create new filtering mode:**

```jsonc
"filtering": {
  "modes": {
    "my_mode": {
      "include_categories": ["compilation", "custom", "principle"]
    }
  },
  "active_mode": "my_mode"
}
```

---

## Modification Policy

### ✅ Safe to Modify (Extension Points)

**Add new feedback functions:**

- Follow existing pattern (check config, fallback, delegate to display)
- Add corresponding config entries
- Document in this API

**Add new build patterns:**

- Use extensions section in config (preferred)
- Or add to existing categories
- No code changes needed

**Extend quality checklist:**

- Add items to config file
- Create new filtering modes
- Add language-specific overrides

**Add new utility functions:**

- Pattern detection helpers
- Message formatting helpers
- Follow existing inline comment style

### ⚠️ Modify with Extreme Care (Breaking Changes)

**Function signatures:**

- All hook code depends on these
- Changing parameters breaks calling code
- Adding optional parameters OK (with defaults)

**Configuration file structure:**

- Changing JSON keys breaks config loading
- Adding new keys OK (backward compatible)
- Removing keys needs migration path

**Display library integration:**

- Changing from display.Success() to raw output loses formatting
- Breaks architectural pattern

### ❌ NEVER Modify (Foundational Rails)

**4-block structure:**

- METADATA → SETUP → BODY → CLOSING
- Required for all CPI-SI code

**Delegation to system/lib/display:**

- Architectural pattern: hooks → feedback → display
- Don't duplicate formatting logic

**Non-blocking behavior:**

- Feedback must never crash workflow
- Failures degrade gracefully

**Kingdom-honoring communication:**

- Maintain gracious, encouraging tone
- Build up, don't tear down

---

## Troubleshooting

### No Feedback Appears

**Symptom:** No output after tool operations

**Diagnosis:**

```bash
# Check if hooks are calling feedback functions
grep -r "feedback\." /home/seanje-lenox-wise/.claude/hooks/

# Check if display library working
go run -c 'import "system/lib/display"; fmt.Println(display.Success("Test"))'
```

**Common Causes:**

1. Hooks not importing feedback library
2. Hooks not calling feedback functions
3. Display library not available

**Solutions:**

- Review hook implementation
- Verify imports in hooks
- Test display library directly

---

### Build Feedback Not Showing

**Symptom:** OnBuildComplete() not called after build commands

**Diagnosis:**

```bash
# Test pattern detection
cat > /tmp/test_patterns.go <<'EOF'
package main
import (
    "fmt"
    "hooks/lib/feedback"
)
func main() {
    cmds := []string{"go build", "make", "gcc main.c"}
    for _, cmd := range cmds {
        fmt.Printf("%s: %v\n", cmd, feedback.IsBuildCommand(cmd))
    }
}
EOF
go run /tmp/test_patterns.go
```

**Common Causes:**

1. IsBuildCommand() doesn't detect your build tool
2. Hook not checking IsBuildCommand() before calling OnBuildComplete()
3. Build tool pattern missing from config

**Solutions:**

- Add pattern to build-patterns.jsonc
- Verify hook logic
- Check config file loaded successfully

---

### Config Not Loading

**Symptom:** Fallback messages always used, config changes ignored

**Diagnosis:**

```bash
# Check config files exist
ls -la ~/.claude/cpi-si/system/data/config/feedback/

# Check config file syntax
cd ~/.claude/cpi-si/system/data/config/feedback/
# Strip comments and validate JSON
grep -v "^//" messages.jsonc | jq .

# Check if library compiled with config support
strings /tmp/feedback-test.so | grep configLoaded
```

**Common Causes:**

1. Config files missing or in wrong location
2. JSONC syntax errors (malformed JSON after comment stripping)
3. HOME environment variable not set
4. Config directory permissions

**Solutions:**

- Verify file paths: `~/.claude/cpi-si/system/data/config/feedback/`
- Validate JSON syntax (strip // comments first)
- Check HOME: `echo $HOME`
- Fix permissions: `chmod 755 ~/.claude/cpi-si/system/data/config/feedback/`

---

### Custom Patterns Not Working

**Symptom:** Added patterns to config but IsBuildCommand() doesn't detect them

**Diagnosis:**

```bash
# Verify pattern added to config
cat ~/.claude/cpi-si/system/data/config/feedback/build-patterns.jsonc | grep "mytool"

# Test pattern detection with custom pattern
go run -c 'feedback.IsBuildCommand("mytool build")'
```

**Common Causes:**

1. Pattern not saved to config file
2. JSONC syntax error after adding pattern
3. Pattern in wrong section (not in patterns.*.patterns array)
4. Library not recompiled after config change

**Solutions:**

- Verify pattern in correct location
- Validate JSON syntax
- Restart hooks system (config loaded at init())
- Recompile if library changed

---

## Future Roadmap

### Planned Features

**✓ Configuration-driven system** - COMPLETED (v2.0.0)

- Messages, patterns, and checklists in config files
- Graceful fallback to hardcoded defaults
- Extension points for customization

**⏳ Context-aware feedback** - PLANNED (v3.0.0)

- Analyze command output for specific issues
- Different messages based on time of day (morning energy vs late night)
- Track what feedback was helpful (activity logging integration)

**⏳ Adaptive encouragement** - PLANNED (v3.0.0)

- Vary messages to avoid monotony
- Don't repeat same message every time
- Rotate through equivalent messages

**⏳ Pattern-based feedback** - FUTURE

- Detect specific tools and provide targeted guidance
- Tool-specific best practices
- Language-specific quality standards

**⏳ Feedback personalization** - FUTURE

- Per-user message customization
- Per-instance tone preferences
- Team-specific quality standards

### Research Areas

- **Command output analysis** - Parse stderr/stdout for specific issues
- **Sentiment tracking** - Does feedback resonate? Is it helpful?
- **Integration with session quality** - Correlate feedback with work outcomes
- **Pattern learning** - Discover what feedback helps most
- **Internationalization** - Support multiple languages

### Integration Targets

- **Session quality tracking** - Feed into quality indicators
- **Activity logging** - Track which feedback appears when
- **Pattern learning** - Inform future feedback improvements
- **User preferences** - Learn preferred communication style

### Known Limitations

- **Build detection limited to substring matching** - Could use regex or AST parsing
- **No feedback personalization** - Same messages for everyone
- **No context from command output** - Only sees command string itself
- **Single active filtering mode** - Can't combine modes dynamically
- **No adaptive messaging** - Same message every time

### Version History

**2.0.0 (2025-11-11)** - Configuration-Driven Architecture

- Refactored to load from system config files
- Added build-patterns.jsonc with 40+ patterns across 9 categories
- Added messages.jsonc with template variable substitution
- Added quality-checklist.jsonc with filtering modes
- Graceful fallback to hardcoded defaults
- Delegates formatting to system/lib/display
- Comprehensive template alignment
- Full inline comment documentation

**1.0.0 (2024-10-24)** - Initial Implementation

- Basic feedback functions
- Manual ANSI formatting
- Hardcoded build patterns (7 patterns)
- Hardcoded messages
- Hardcoded 3-item quality checklist

---

## Closing Note

This library demonstrates Kingdom Technology: communication that edifies, encourages, and guides toward excellence. Feedback should serve others, not burden them. Messages should build up, not tear down.

*"Therefore exhort one another, and build each other up" - 1 Thessalonians 5:11 (WEB)*

**For questions, issues, or contributions:** Follow the modification policy above and maintain the gracious, Kingdom-honoring tone.
