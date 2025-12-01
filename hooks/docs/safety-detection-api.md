# Safety Detection Library - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for pattern-based detection of dangerous operations, critical file paths, and likely secrets

**Component:** `hooks/lib/safety/detection.go`

**Last Updated:** November 11, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Architecture](#architecture)
- [Public API Reference](#public-api-reference)
- [Configuration System](#configuration-system)
- [Integration Patterns](#integration-patterns)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The Safety Detection Library provides pattern-based detection of dangerous operations, critical file paths, and likely secrets. It enables BLOCKING hooks (tool/pre-use) and NON-BLOCKING secret detection (prompt/submit) to protect against accidental data loss, credential exposure, and system damage.

**Core Responsibility:** Detect potentially dangerous operations through configurable pattern matching with graceful fallback to hardcoded defaults. Return boolean detection results for BLOCKING decision making.

**Design Pattern:** Pure detection logic with stateless functions. Configuration-driven with conservative bias (better to allow than falsely block). Privacy-preserving secret detection (never logs actual content).

**Biblical Foundation:** "The prudent see danger and take refuge, but the simple keep going and pay the penalty" (Proverbs 27:12 WEB). Detection enables wisdom - recognizing danger before harm occurs.

---

## Philosophy & Design Rationale

### Why Pattern Detection Matters

Safety relies on recognizing danger before it happens:

- **Without detection:** Dangerous operations execute without warning, data loss occurs accidentally
- **With detection:** User warned before destructive operations, given opportunity to reconsider

Detection is not control - it's **wisdom through early recognition**.

### The Prudence Principle (Proverbs 27:12)

> "The prudent see danger and take refuge, but the simple keep going and pay the penalty." - Proverbs 27:12 (WEB)

Wisdom sees danger ahead. Simplicity rushes blindly forward. Safety detection establishes:

- **Prudence** - Recognize danger patterns before execution
- **Wisdom** - Distinguish truly dangerous from merely unfamiliar
- **Protection** - Enable informed decisions without restricting work
- **Stewardship** - Care for user's work and credentials
- **Grace** - Conservative bias (allow when uncertain)

Just as prudent person sees danger and takes refuge, detection recognizes patterns requiring confirmation.

### The Three Types of Detection

Safety detection handles three distinct categories:

| Detection Type | Purpose | Usage Context |
|----------------|---------|---------------|
| **Dangerous Operations** | Commands that can cause data loss or team disruption | BLOCKING pre-tool validation (Bash commands) |
| **Critical Paths** | File locations requiring extra confirmation | BLOCKING pre-tool validation (Write operations) |
| **Secret Patterns** | API keys, tokens, passwords in content | NON-BLOCKING warnings (prompt submission, file writes) |

### The Conservative Bias Philosophy

**Core Principle:** Better to allow than falsely block.

**Reasoning:**

- False positive (detecting safe operation as dangerous) = workflow disruption, user frustration, lost trust
- False negative (missing dangerous operation) = defeats purpose, but user can still prevent
- User frustration from false positives worse than occasional missed danger
- Trust requires reliability - detection must be surgically precise

This reflects Kingdom Technology: **Serve through wisdom, not restrict through fear**.

### The Privacy Preservation Principle

Secret detection has special responsibility:

**What Detection Does:**
- Identifies patterns indicating likely secrets (API key prefixes, token patterns, private key headers)
- Returns boolean result (detected or not detected)
- Enables calling code to warn user

**What Detection NEVER Does:**
- Log actual secret content (only length or metadata)
- Store secrets in any form
- Transmit secrets anywhere
- Include secrets in error messages

**Privacy preservation:** The function detects, the caller warns with metadata only.

### The Configuration-Driven Philosophy

**Why Configuration Matters:**

Hardcoded patterns cannot adapt to:
- Company-specific dangerous operations
- Project-specific critical paths
- Organization-specific secret patterns
- User workflow variations

**Configuration enables:**
- Adaptation without code changes
- User/team customization
- Easy pattern updates
- Project-specific rules

**Graceful fallback ensures:**
- Library works without configuration files
- Missing config doesn't break functionality
- Hardcoded defaults provide baseline protection

---

## Architecture

### Component Structure

```
Safety Detection Library
├── Configuration Loading (Init-time)
│   ├── dangerous-patterns.jsonc → DangerousPatternsConfig
│   ├── critical-paths.jsonc → CriticalPathsConfig
│   └── secret-patterns.jsonc → SecretPatternsConfig
│
├── Public Detection APIs
│   ├── IsDangerousOperation(cmd) → bool
│   ├── IsCriticalFile(path) → bool
│   └── ContainsLikelySecret(text) → bool
│
└── Internal Helpers
    ├── matchesAnyPattern(text, patterns) → bool
    ├── stripJSONCComments(jsonc) → json
    └── load*() functions for each config type
```

### Ladder Position

**Mid-rung library:**
- Lower rungs: None (no internal dependencies - pure standard library)
- Same rung: Other hook libraries (feedback, monitoring, validation)
- Higher rungs: BLOCKING hooks (tool/pre-use, prompt/submit)

### Baton Flow

```
BLOCKING Hook (tool/pre-use)
  ↓
IsDangerousOperation("git push --force")
  ↓
Check configLoaded && dangerousConfig != nil
  ├─ Yes → Flatten config patterns into array
  └─ No  → Use fallbackDangerousPatterns
  ↓
matchesAnyPattern(cmd, patterns)
  ↓
Return true (detected) or false (not detected)
  ↓
Hook makes BLOCKING decision based on result
```

### Health Scoring

**Configuration Loading (30 points):**
- Dangerous patterns loaded: +10
- Critical paths loaded: +10
- Secret patterns loaded: +10
- Config load failure: -10 (per failed config)
- Fallback to defaults: 0 (neutral - still functional)

**Detection Operations (70 points):**
- Successful pattern match: +35
- Successful no-match: +35
- Detection error (panic): -70

Total: 100 points for complete detection cycle

---

## Public API Reference

### IsDangerousOperation

```go
func IsDangerousOperation(cmd string) bool
```

**Purpose:** Detect potentially dangerous command patterns (force push, hard reset, rm -rf, sudo, database drops, etc.)

**Parameters:**
- `cmd` (string): Command string to check (e.g., "git push --force origin main")

**Returns:**
- `true` if command matches dangerous pattern
- `false` if command appears safe

**Conservative Bias:**
When uncertain, returns false (allow). False positives disrupt workflow, so detection errs on side of permissiveness.

**Example Usage:**

```go
import "hooks/lib/safety"

if safety.IsDangerousOperation("git push --force") {
    // Request user confirmation before proceeding
    confirmed := requestConfirmation("Force push detected. Proceed?")
    if !confirmed {
        os.Exit(1) // Block execution
    }
}
// Continue with operation
os.Exit(0)
```

**Patterns Detected (from config or fallback):**
- Git destructive: `git push --force`, `git reset --hard`, `git rebase --skip`
- Filesystem destructive: `rm -rf`, `rm -r`
- Privilege escalation: `sudo`, `pkexec`
- Package publishing: `npm publish`, `cargo publish`, `pip upload`
- Database destructive: `DROP DATABASE`, `DROP TABLE`, `TRUNCATE TABLE`
- Container destructive: `docker system prune`, `docker rm -f`

### IsCriticalFile

```go
func IsCriticalFile(filePath string) bool
```

**Purpose:** Detect critical system file paths requiring extra confirmation for modifications (/etc/, /.ssh/, /boot/, git config, etc.)

**Parameters:**
- `filePath` (string): File path to check (e.g., "/etc/passwd" or "~/.ssh/id_rsa")

**Returns:**
- `true` if path matches critical location
- `false` if path appears normal

**Conservative Bias:**
When uncertain, returns false (allow). False positives disrupt workflow, so detection errs on side of permissiveness.

**Example Usage:**

```go
import "hooks/lib/safety"

filePath := os.Getenv("FILE_PATH")
if safety.IsCriticalFile(filePath) {
    // Request user confirmation before writing
    confirmed := requestConfirmation("Critical file modification. Proceed?")
    if !confirmed {
        os.Exit(1) // Block execution
    }
}
// Continue with write operation
os.Exit(0)
```

**Paths Detected (from config or fallback):**
- System configuration: `/etc/`, `/boot/`, `/sys/`, `/proc/`
- Authentication: `/.ssh/`, `/.gnupg/`, `/etc/shadow`, `/etc/passwd`
- Version control: `/.git/config`, `/.git/HEAD`, `/.gitconfig`
- Package management: `/var/lib/dpkg/`, `/var/lib/rpm/`
- User environment: `/.bashrc`, `/.zshrc`, `/.config/`
- Database files: `/var/lib/mysql/`, `*.db`, `*.sqlite`

### ContainsLikelySecret

```go
func ContainsLikelySecret(text string) bool
```

**Purpose:** Detect likely secrets in text content (API keys, tokens, private keys) for NON-BLOCKING warnings

**Parameters:**
- `text` (string): Text to scan for secrets (prompt content, file content, etc.)

**Returns:**
- `true` if text likely contains secrets
- `false` if text appears normal

**Privacy Preservation:**
This function ONLY detects - it NEVER logs or stores actual secret content. Calling code responsible for privacy-preserving warnings (log length only, never content).

**Conservative Bias:**
When uncertain, returns true (warn). Better to warn unnecessarily than miss actual secrets. False positives acceptable for secret detection (warning, not blocking).

**Example Usage:**

```go
import "hooks/lib/safety"

promptText := getUserPrompt()
if safety.ContainsLikelySecret(promptText) {
    // Warn user (log length only, NEVER content)
    fmt.Printf("⚠️  Warning: Prompt may contain secrets (length: %d)\n", len(promptText))
    fmt.Println("Consider using environment variables or secret management.")
}
// Continue (NON-BLOCKING - warning only)
```

**Patterns Detected (from config or fallback):**
- API keys: `sk-` (OpenAI), `sk_live_` (Stripe), `AKIA` (AWS)
- GitHub tokens: `ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_`
- Communication tokens: `xox` (Slack), `xapp-` (Slack app)
- Private keys: `BEGIN PRIVATE`, `BEGIN RSA PRIVATE`, `BEGIN OPENSSH PRIVATE`
- Database URIs: `mongodb://`, `postgresql://`, `mysql://` (may contain credentials)
- Cloud keys: `AIza` (Google API), `GOOGLE_API_KEY`
- Environment variables: `API_KEY=`, `SECRET_KEY=`, `PASSWORD=`
- JWT tokens: `eyJ` (base64 encoded header)

**False Positives:**
Some patterns (like `eyJ` for JWT) are common in non-secret base64 data. This is acceptable - better to warn unnecessarily than miss actual secrets.

---

## Configuration System

### Configuration Files

**Location:** `$HOME/.claude/cpi-si/system/data/config/safety/`

**Files:**
1. `dangerous-patterns.jsonc` - Dangerous command patterns
2. `critical-paths.jsonc` - Critical file path patterns
3. `secret-patterns.jsonc` - Secret detection patterns

### Configuration Loading

**When:** Init-time (once when library imported)

**Process:**
1. Determine home directory (`$HOME` or fallback)
2. Build paths to config files
3. Load each config independently (partial success acceptable)
4. Set `configLoaded` flag if all three loaded successfully

**Graceful Fallback:**
If any config fails to load:
- Detection functions use hardcoded fallback patterns
- Library continues functioning normally
- No errors thrown (silent graceful degradation)

### Configuration Structure

**dangerous-patterns.jsonc:**

```jsonc
{
  "$schema": "../../../config/schemas/safety/dangerous-patterns.schema.json",
  "metadata": {
    "name": "Dangerous Command Patterns",
    "description": "Patterns for detecting potentially destructive operations",
    "version": "1.0.0"
  },
  "patterns": {
    "git_destructive": {
      "patterns": ["git push --force", "git reset --hard"],
      "description": "Git operations that can destroy work",
      "severity": "high",
      "reason": "Can permanently lose work or disrupt team"
    }
    // ... more categories
  }
}
```

**critical-paths.jsonc:**

```jsonc
{
  "$schema": "../../../config/schemas/safety/critical-paths.schema.json",
  "metadata": {
    "name": "Critical File Paths",
    "description": "File paths requiring extra confirmation",
    "version": "1.0.0"
  },
  "paths": {
    "system_config": {
      "patterns": ["/etc/", "/boot/", "/sys/"],
      "description": "System configuration directories",
      "severity": "critical",
      "reason": "Modification can break system operation",
      "platform": "linux/unix"
    }
    // ... more categories
  }
}
```

**secret-patterns.jsonc:**

```jsonc
{
  "$schema": "../../../config/schemas/safety/secret-patterns.schema.json",
  "metadata": {
    "name": "Secret Detection Patterns",
    "description": "Patterns for detecting likely secrets",
    "version": "1.0.0"
  },
  "patterns": {
    "api_keys": {
      "patterns": ["sk-", "sk_live_", "AKIA"],
      "description": "Common API key prefixes",
      "severity": "high",
      "reason": "API keys grant programmatic access",
      "services": ["OpenAI", "Stripe", "AWS"]
    }
    // ... more categories
  }
}
```

### Customizing Patterns

**Adding New Category:**

1. Edit appropriate config file (dangerous-patterns.jsonc, critical-paths.jsonc, secret-patterns.jsonc)
2. Add new category following existing structure
3. Restart process to reload config (init-time loading)
4. Test with detection function

**Example - Adding Custom Dangerous Pattern:**

```jsonc
"custom_deploy": {
  "patterns": ["deploy --production", "ship --live"],
  "description": "Direct production deployment commands",
  "severity": "critical",
  "reason": "Deploys directly to production without staging"
}
```

**No code changes needed** - patterns automatically loaded and used.

### Fallback Patterns

**Purpose:** Ensure library works without configuration files

**Location:** Hardcoded in detection.go as package-level variables

**Scope:** Minimal baseline protection (10 dangerous patterns, 4 critical paths, 5 secret patterns)

**When Used:**
- Config files don't exist
- Config files have invalid JSON
- $HOME environment variable not set
- Any config loading error

**Trade-off:** Less comprehensive than config files, but ensures basic functionality always available

---

## Integration Patterns

### BLOCKING Hook Integration (tool/pre-use)

**Purpose:** Detect dangerous operations and request confirmation before execution

**Pattern:**

```go
package main

import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func preUseHook() {
    toolName := os.Args[1]

    if toolName == "Bash" {
        cmd := os.Args[2]

        // Detect dangerous operation
        if safety.IsDangerousOperation(cmd) {
            // Display warning
            fmt.Printf("⚠️  Dangerous operation detected: %s\n", cmd)

            // Request confirmation
            confirmed := requestConfirmation("Proceed?")

            if !confirmed {
                fmt.Println("❌ Operation cancelled")
                os.Exit(1) // BLOCK execution
            }
        }
    }

    os.Exit(0) // ALLOW execution
}
```

**Key Points:**
- Detection is boolean (true/false)
- Hook makes blocking decision based on result
- Conservative bias means some dangerous operations might not detect (acceptable)

### NON-BLOCKING Secret Detection (prompt/submit)

**Purpose:** Warn about likely secrets without blocking submission

**Pattern:**

```go
package main

import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func submitHook() {
    promptText := os.Args[1]

    // Detect likely secrets
    if safety.ContainsLikelySecret(promptText) {
        // Privacy-preserving warning (length only, never content)
        fmt.Printf("⚠️  Warning: Prompt may contain secrets (length: %d)\n", len(promptText))
        fmt.Println("Consider using environment variables or secret management.")
        // Continue - NON-BLOCKING
    }

    os.Exit(0) // Always allow (NON-BLOCKING design)
}
```

**Key Points:**
- Detection is boolean (true/false)
- Hook warns but never blocks
- Privacy preserved (log length, not content)
- False positives acceptable (warning, not blocking)

### Multiple Detection Types

**Pattern:**

```go
package main

import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func preUseHook() {
    toolName := os.Args[1]

    if toolName == "Bash" {
        cmd := os.Args[2]

        // Check dangerous operation
        if safety.IsDangerousOperation(cmd) {
            if !requestConfirmation("Dangerous command. Proceed?") {
                os.Exit(1) // BLOCK
            }
        }
    }

    if toolName == "Write" {
        filePath := os.Getenv("FILE_PATH")

        // Check critical path
        if safety.IsCriticalFile(filePath) {
            if !requestConfirmation("Critical file modification. Proceed?") {
                os.Exit(1) // BLOCK
            }
        }

        // Check for secrets (get content from args[1])
        content := os.Args[2]
        if safety.ContainsLikelySecret(content) {
            fmt.Printf("⚠️  Warning: File may contain secrets (length: %d)\n", len(content))
            // Continue - warning only for Write operations
        }
    }

    os.Exit(0) // ALLOW
}
```

---

## Modification Policy

### Safe to Modify (Extension Points)

✅ **Add new patterns to configuration files**
- Edit dangerous-patterns.jsonc, critical-paths.jsonc, secret-patterns.jsonc
- Add new categories following existing structure
- No code changes needed

✅ **Extend fallback pattern arrays**
- Add patterns to fallbackDangerousPatterns, fallbackCriticalPaths, fallbackSecretPatterns
- Ensures baseline protection improves over time

✅ **Improve pattern matching logic**
- Enhance matchesAnyPattern() function
- Add regex support, case-insensitive mode, word boundaries
- Update all three detection functions to use enhanced matcher

✅ **Add new detection functions**
- Follow existing naming pattern (IsSomethingDangerous)
- Create corresponding config file if needed
- Add graceful fallback patterns

✅ **Documentation and comments**
- Update inline comments for clarity
- Expand docstrings with examples
- Add troubleshooting guidance

### Modify with Extreme Care (Breaking Changes)

⚠️ **Public API function signatures**
- `IsDangerousOperation(cmd string) bool` signature change breaks all calling hooks
- `IsCriticalFile(path string) bool` signature change breaks tool/pre-use
- `ContainsLikelySecret(text string) bool` signature change breaks prompt/submit

⚠️ **Config struct definitions**
- Must match JSONC schema exactly
- Breaking change prevents config loading
- Results in fallback patterns being used

⚠️ **Conservative bias principle**
- Making detection too aggressive creates false positives
- False positives disrupt workflow and erode trust
- Conservative bias is foundational design principle

⚠️ **Privacy preservation**
- Secret detection must never log content
- Only metadata (length, detection result) allowed
- Privacy violation breaks trust completely

⚠️ **Boolean return types**
- Hooks depend on true/false semantics
- Changing to int/string/error breaks integration
- Exit codes (0/1) map directly to boolean results

### NEVER Modify (Foundational Rails)

❌ **4-block structure** (METADATA, SETUP, BODY, CLOSING)
❌ **Graceful fallback behavior** (library must work without configs)
❌ **Stateless nature** (no side effects, thread-safe guarantee)
❌ **Package name or import path** (hooks depend on "hooks/lib/safety")
❌ **Init-time configuration loading** (no runtime reload without rearchitecture)

### Testing Requirements

Before deploying ANY changes:

1. Test all detection types (dangerous operations, critical paths, secrets)
2. Test with configuration files present (should use config patterns)
3. Test with configuration files missing (should use fallback patterns)
4. Test with malformed configuration (should fall back gracefully)
5. Verify conservative bias maintained (false negatives acceptable, false positives not)
6. Test privacy preservation (secret detection never logs content)
7. Verify performance (< 100μs per detection call)
8. Check for false positives (most critical test)

---

## Future Roadmap

### Planned Features

**Regex Pattern Support (High Priority):**
- Move beyond simple substring matching
- Support complex patterns: `git\s+push\s+--force` instead of literal "git push --force"
- Reduce false positives with more precise matching
- Maintain conservative bias (strict regex, not loose)

**Platform-Aware Filtering (High Priority):**
- Different patterns for Linux vs Windows vs macOS
- Critical paths vary by platform (/etc/ vs C:\Windows\)
- Automatic platform detection and pattern filtering
- Cross-platform projects work correctly

**Severity-Based Filtering (Medium Priority):**
- Configuration specifies severity (critical, high, medium, low)
- BLOCKING hooks can filter by severity (only block critical/high)
- User customization of severity thresholds
- Adaptive validation based on context

**Runtime Config Reload (Medium Priority):**
- Watch config files for changes
- Reload patterns without process restart
- Atomic updates (all-or-nothing reload)
- Fallback if reload fails

### Research Areas

**Context-Aware Detection:**
- Reduce false positives through operation context
- Different rules per branch (stricter on main/master)
- Different rules per project (personal vs team)
- Time-of-day awareness (more permissive during development hours)

**Entropy Analysis for Secret Detection:**
- Catch random-looking strings (likely secrets) even without known prefixes
- Measure randomness/entropy of string segments
- Flag high-entropy strings for review
- Complement pattern-based detection

**Pattern Confidence Scores:**
- Not just boolean match/no-match
- Confidence level: 0-100 indicating certainty
- Hooks can set confidence thresholds
- User sees confidence in confirmation prompts

**Whitelist Support:**
- Allow specific known-safe patterns to bypass detection
- User confirms operation once, added to whitelist automatically
- Per-project whitelists (not global)
- Expiring whitelist entries (re-confirm after N days)

**Machine Learning Pattern Detection:**
- Train on actual dangerous operations from user history
- Learn user's risk tolerance and workflow patterns
- Automatic pattern suggestions based on confirmation history
- Continuous improvement without manual pattern updates

### Integration Targets

**System Audit Logs:**
- Integration with system audit trails
- Correlation: Did detected operation actually cause harm?
- Learning: Which patterns are false positives?
- Feedback loop for pattern refinement

**Detection Telemetry:**
- Anonymous aggregate statistics on detection accuracy
- False positive/negative rates across users
- Pattern effectiveness measurement
- Community-driven pattern improvements

**Team-Wide Pattern Sharing:**
- Share custom patterns across team
- Learn from collective experience
- Company-specific dangerous operation definitions
- Synchronized pattern updates

**Custom Detection Plugins:**
- User-defined detection functions
- Language-specific validation logic
- Project-specific danger patterns
- Scriptable detection rules

### Known Limitations

**Config loaded once at init:**
- Requires process restart for changes
- Can't adapt patterns during long-running sessions
- Runtime reload would enable dynamic adaptation

**Simple substring matching:**
- No regex or context awareness
- Some false positives unavoidable
- More sophisticated matching would improve precision

**No confidence scoring:**
- Binary true/false only
- Can't express "probably dangerous" vs "definitely dangerous"
- Confidence scores would enable threshold-based blocking

**No platform filtering:**
- All patterns apply to all platforms
- Windows-specific patterns flag on Linux
- Platform awareness would reduce noise

**No severity-based filtering:**
- All matches treated equally
- Can't distinguish critical from moderate danger
- Severity filtering would enable graduated response

### Version History

**2.0.0 (2025-11-11)** - Configuration-driven detection with graceful fallbacks
- Added JSONC configuration loading for all pattern types
- Created dangerous-patterns.jsonc, critical-paths.jsonc, secret-patterns.jsonc
- Implemented graceful fallback to hardcoded defaults if config unavailable
- Comprehensive 4-block template alignment with full documentation
- Design principle: Configuration enables adaptation without code changes

**1.0.0 (2024-10-24)** - Initial implementation with hardcoded patterns
- Dangerous operation detection (force push, hard reset, rm -rf, sudo, etc.)
- Critical path detection (system files, authentication, version control)
- Secret detection (API keys, tokens, private keys)
- Conservative bias design principle established
- Privacy-preserving secret detection (never logs content)

---

## Conclusion

The Safety Detection Library embodies the principle: "The prudent see danger and take refuge." Detection enables wisdom - recognizing danger before harm occurs through pattern-based awareness.

**Key Principles:**

- **Conservative bias** (allow when uncertain - better than false positives)
- **Configuration-driven** (adapt without code changes)
- **Graceful fallback** (always functional, even without configs)
- **Privacy-preserving** (secret detection never logs content)
- **Stateless operation** (pure functions, thread-safe, no side effects)

**Success Metrics:**

- Minimal false positives (don't disrupt safe operations)
- Comprehensive coverage (detect actual dangerous patterns)
- Fast execution (< 100μs per detection call)
- Privacy maintained (no secret content ever logged)
- User trust preserved (reliable, predictable, non-intrusive)

Kingdom Technology recognizes danger through wisdom, not paranoia. Detection serves the user through protection, not restriction through fear.

*"The prudent see danger and take refuge, but the simple keep going and pay the penalty." - Proverbs 27:12 (WEB)*
