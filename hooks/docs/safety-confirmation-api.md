# Safety Confirmation Library - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for user confirmation orchestration in BLOCKING safety validation

**Component:** `hooks/lib/safety/confirmation.go`

**Last Updated:** November 11, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Architecture](#architecture)
- [Public API Reference](#public-api-reference)
- [Integration Patterns](#integration-patterns)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The Safety Confirmation Library orchestrates user confirmation flows for dangerous operations detected by the Safety Detection Library. It displays contextual warnings with temporal awareness, prompts users for explicit confirmation, validates responses for exact match, and returns boolean results for BLOCKING decision making.

**Core Responsibility:** Route dangerous operations to appropriate warning displays, get user confirmation with exact match validation, integrate temporal context for enhanced awareness, return BLOCKING decision results.

**Design Pattern:** Confirmation orchestration layer connecting detection (what's dangerous) to decision (allowed/blocked). Stateless user interaction with conservative response validation. Privacy-preserving secret warnings (NON-BLOCKING awareness).

**Biblical Foundation:** "The plans of the diligent surely lead to profit; and everyone who is hasty surely rushes to poverty" (Proverbs 21:5 WEB). Confirmation requires pause before destructive operations - diligence prevents hasty loss.

---

## Philosophy & Design Rationale

### Why Confirmation Matters

Safety requires both detection AND informed decision-making:

- **Detection alone:** Recognizes danger but doesn't prevent execution
- **Confirmation:** Pauses execution, communicates risk, enables informed choice
- **Together:** Detection identifies, confirmation communicates and decides

Confirmation is not gatekeeping - it's **protection through informed pause**.

### The Diligence Principle (Proverbs 21:5)

> "The plans of the diligent surely lead to profit; and everyone who is hasty surely rushes to poverty." - Proverbs 21:5 (WEB)

Diligence requires pause. Haste rushes to loss. Safety confirmation establishes:

- **Pause** - Interrupt between detection and execution
- **Communication** - Clear warning explaining specific danger
- **Context** - Temporal awareness (session length, time of day)
- **Decision** - User makes informed choice after understanding risk
- **Protection** - Prevent hasty destructive operations

Just as diligent plans lead to profit, diligent confirmation prevents hasty loss.

### The Companion Scripture (Proverbs 22:3)

> "A prudent person foresees danger and takes precautions. The simpleton goes blindly on and suffers the consequences." - Proverbs 22:3 (NLT)

Detection foresees danger (Proverbs 27:12). Confirmation takes precautions (Proverbs 22:3). Together they enable prudence:

- **Detection** - "See danger ahead"
- **Confirmation** - "Take precautions before proceeding"

The simpleton rushes through confirmations without reading. The prudent reads warnings, understands risks, decides wisely.

### The Three Confirmation Types

Safety confirmation handles three distinct operation categories:

| Confirmation Type | Purpose | Response Validation | Blocking |
|-------------------|---------|---------------------|----------|
| **Bash Operations** | Dangerous commands (force push, rm -rf, sudo) | Varies by severity ("yes" or "y") | BLOCKING |
| **File Writes** | Critical file paths (/etc/, /.ssh/, git config) | Full "yes" required | BLOCKING |
| **Secret Detection** | API keys, tokens in prompts | No response needed | NON-BLOCKING (warning only) |

### The Contextual Warning Philosophy

**Principle:** Warnings should communicate *why* dangerous and provide *temporal context* for enhanced awareness.

**Why Temporal Context Matters:**

- **Long session warnings:** After 4+ hours, fatigue increases error risk → "⏰ Note: Session has been running 4+ hours - extra caution advised"
- **Time of day awareness:** Late night (after 22:00) or very early (before 06:00) → "🌙 Note: It's quite late - consider if this can wait until morning"
- **Enhanced decision quality:** Context helps user recognize when they might be rushing or tired

**Warning Design Principles:**

1. **Specific** - Explain exact danger for this operation type
2. **Educational** - Help user understand why dangerous
3. **Contextual** - Include temporal awareness when relevant
4. **Respectful** - Clear without condescension
5. **Actionable** - Show exactly what's being confirmed

### The Conservative Validation Philosophy

**Principle:** Ambiguous responses default to denial for safety.

**Response Validation Strategy:**

| Response | Interpretation | Reasoning |
|----------|----------------|-----------|
| "yes" (full word) | Confirmed | Explicit, intentional, high-severity operations |
| "y" (single char) | Confirmed | Quick confirmation, medium-severity operations |
| Anything else | Denied | Typos, uncertainty, accidental enter → default to safe |

**Why This Matters:**

- User accidentally hits enter → Operation blocked (safe default)
- User types "yse" (typo) → Operation blocked (better than executing mistake)
- User types "yeah" → Operation blocked (require exact match)
- Forces intentional confirmation, prevents hasty accidents

This reflects Kingdom Technology: **Protection through intentional pause, not restriction through fear**.

### The Privacy Preservation Principle

Secret detection warnings are NON-BLOCKING and privacy-preserving:

**What Confirmation Does:**
- Displays warning that likely secret detected
- Suggests user review prompt before submitting
- Returns immediately without prompting (NON-BLOCKING)

**What Confirmation NEVER Does:**
- Log which secret pattern matched
- Display any portion of secret content
- Block execution (awareness only, user decides)
- Store or transmit detected secrets

**Design:** Detection identifies (boolean), confirmation warns (metadata only), user decides.

---

## Architecture

### Architectural Position

```
┌──────────────────────────────────────────────────────────┐
│ HOOK LAYER (tool/pre-use, prompt/submit)                │ Top Rung
│ - Orchestrates safety validation                        │ (Commands)
│ - Makes BLOCKING decisions                              │
└────────────────────┬─────────────────────────────────────┘
                     │ calls
                     ↓
┌──────────────────────────────────────────────────────────┐
│ CONFIRMATION LAYER (hooks/lib/safety/confirmation.go)   │ Mid Rung
│ - Routes to appropriate warning displays                │ (Libraries)
│ - Prompts for user confirmation                         │
│ - Validates responses with exact match                  │
│ - Returns BLOCKING decision (needs, allowed)            │
└────────────────────┬─────────────────────────────────────┘
                     │ uses
                     ↓
┌──────────────────────────────────────────────────────────┐
│ DETECTION LAYER (hooks/lib/safety/detection.go)         │ Mid Rung
│ - Pattern matching (IsDangerousOperation, etc.)         │ (Libraries)
│ - Configuration loading                                 │
│ - Boolean detection results                             │
└────────────────────┬─────────────────────────────────────┘
                     │ uses
                     ↓
┌──────────────────────────────────────────────────────────┐
│ DISPLAY LAYER (system/lib/display)                      │ Lower Rung
│ - ANSI color formatting                                 │ (System Libs)
│ - Warning(), Info() formatters                          │
└──────────────────────────────────────────────────────────┘
```

### Execution Flow (BLOCKING Operations)

**Bash Command Validation:**

```
tool/pre-use hook receives Bash command
    ↓
Calls: ConfirmBashOperation(cmd, timeContext)
    ↓
confirmation.go: IsDangerousOperation(cmd)?
    ├─ false → return (false, true)  [no confirmation, allowed]
    └─ true  → continue to warning routing
         ↓
    Route to appropriate warning display:
    - Force push → displayForcePushWarning()
    - Hard reset → displayHardResetWarning()
    - rm -rf → displayRmRfWarning()
    - sudo/pkexec → displaySudoWarning()
    - npm publish → displayPublishWarning()
    - Database commands → displayDatabaseWarning()
         ↓
    Prompt user: confirm("message", "yes" or "y")
         ↓
    Read stdin, validate exact match
         ↓
    Return (true, allowed_bool) to hook
         ↓
Hook makes BLOCKING decision:
    - allowed = true  → Continue execution (exit 0)
    - allowed = false → Block execution (exit 1)
```

**File Write Validation:**

```
tool/pre-use hook receives file write operation
    ↓
Calls: ConfirmFileWrite(filePath, timeContext)
    ↓
confirmation.go: IsCriticalFile(filePath)?
    ├─ false → return (false, true)  [no confirmation, allowed]
    └─ true  → continue to warning
         ↓
    displayCriticalFileWarning(filePath, timeContext)
         ↓
    Prompt user: confirm("Confirm write operation?", "yes")
         ↓
    Read stdin, validate exact "yes" match
         ↓
    Return (true, allowed_bool) to hook
         ↓
Hook makes BLOCKING decision
```

**Secret Detection (NON-BLOCKING):**

```
prompt/submit hook receives user prompt
    ↓
detection.go: ContainsLikelySecret(prompt)?
    ├─ false → Continue without warning
    └─ true  → Call DisplaySecretWarning()
         ↓
    Display NON-BLOCKING warning
    (Suggest review, don't block)
         ↓
    Return immediately
         ↓
Hook continues (no blocking)
```

### Stateless Design

Confirmation is **stateless** - each confirmation flow is independent:

- No package-level state tracking confirmation history
- No session-based confirmation counting
- Each operation confirmed individually on its own merits
- Temporal context passed in (not tracked internally)

**Why Stateless:**
- Simpler reasoning about behavior
- No state corruption risks
- Temporal context comes from hook (session-aware layer)
- Each dangerous operation deserves its own confirmation

---

## Public API Reference

### ConfirmBashOperation

```go
func ConfirmBashOperation(cmd string, timeContext string) (needsConfirmation bool, allowed bool)
```

**Purpose:** Orchestrate confirmation for dangerous Bash commands (force push, hard reset, rm -rf, sudo, publish, database operations)

**Parameters:**
- `cmd` (string): The Bash command to validate (e.g., "git push --force origin main")
- `timeContext` (string): Temporal awareness context from hook ("long_session", "late_night", "very_early", "normal", "")

**Returns:**
- `needsConfirmation` (bool): True if confirmation was required (operation was dangerous)
- `allowed` (bool): True if user confirmed operation, false if denied

**Confirmation Behavior:**

| Operation Type | Warning Function | Required Response | Reasoning |
|----------------|------------------|-------------------|-----------|
| Force push | `displayForcePushWarning` | "yes" (full word) | High severity - team coordination critical |
| Hard reset | `displayHardResetWarning` | "yes" (full word) | High severity - data loss critical |
| rm -rf | `displayRmRfWarning` | "yes" (full word) | High severity - permanent deletion |
| Database DROP/TRUNCATE | `displayDatabaseWarning` | "yes" (full word) | High severity - data loss critical |
| sudo/pkexec | `displaySudoWarning` | "y" (single char) | Medium severity - system modification |
| npm/cargo publish | `displayPublishWarning` | "y" (single char) | Medium severity - public release |

**Temporal Context Integration:**

Warnings include context-specific notes:
- `"long_session"` (4+ hours) → "⏰ Note: Session has been running 4+ hours - extra caution advised"
- `"late_night"` (after 22:00) → "🌙 Note: It's quite late - consider if this can wait until morning"
- `"very_early"` (before 06:00) → "🌅 Note: It's very early - make sure you're fully alert"
- `"normal"` or `""` → No additional temporal note

**Example Usage:**

```go
import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func main() {
    cmd := os.Getenv("BASH_COMMAND")
    timeContext := os.Getenv("TIME_CONTEXT")  // From session-aware hook

    needs, allowed := safety.ConfirmBashOperation(cmd, timeContext)

    if needs && !allowed {
        fmt.Println("❌ Operation cancelled by user")
        os.Exit(1)  // BLOCK execution
    }

    // Continue with operation
    fmt.Println("✓ Operation allowed")
    os.Exit(0)
}
```

**Conservative Behavior:**
- Non-dangerous commands: Return (false, true) immediately - no confirmation needed
- Ambiguous responses: Interpreted as denial - safe default
- Unknown operation types: Allowed to pass (conservative bias from detection layer)

---

### ConfirmFileWrite

```go
func ConfirmFileWrite(filePath string, timeContext string) (needsConfirmation bool, allowed bool)
```

**Purpose:** Orchestrate confirmation for critical file path writes (system config, authentication, version control, databases)

**Parameters:**
- `filePath` (string): File path being written (e.g., "/etc/passwd", "~/.ssh/id_rsa")
- `timeContext` (string): Temporal awareness context from hook ("long_session", "late_night", etc.)

**Returns:**
- `needsConfirmation` (bool): True if confirmation was required (file path was critical)
- `allowed` (bool): True if user confirmed write, false if denied

**Confirmation Behavior:**
- **All critical file writes:** Require full "yes" response (high severity - system modification)
- **Warning display:** `displayCriticalFileWarning(filePath, timeContext)` shows specific file path and temporal context
- **Response validation:** Exact "yes" match required (no shortcuts for critical system files)

**Example Usage:**

```go
import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func main() {
    filePath := os.Getenv("FILE_PATH")
    timeContext := os.Getenv("TIME_CONTEXT")  // From session-aware hook

    needs, allowed := safety.ConfirmFileWrite(filePath, timeContext)

    if needs && !allowed {
        fmt.Println("❌ Write operation cancelled by user")
        os.Exit(1)  // BLOCK execution
    }

    // Continue with write
    fmt.Println("✓ Write operation allowed")
    os.Exit(0)
}
```

**Critical Paths Detected:**
- System configuration: `/etc/`, `/boot/`, `/sys/`, `/proc/`
- Authentication: `/.ssh/`, `/.gnupg/`, `/etc/shadow`, `/etc/passwd`, `/etc/sudoers`
- Version control: `/.git/config`, `/.git/HEAD`, `/.gitconfig`
- Package management: `/var/lib/dpkg/`, `/var/lib/rpm/`, `/usr/local/lib/`
- User environment: `/.bashrc`, `/.bash_profile`, `/.zshrc`, `/.profile`, `/.config/`
- Database files: `/var/lib/mysql/`, `/var/lib/postgresql/`, `*.db`, `*.sqlite`
- Windows system: `C:\Windows\`, `C:\Program Files\`, `C:\System32\`

**Conservative Behavior:**
- Non-critical paths: Return (false, true) immediately - no confirmation needed
- Conservative bias inherited from IsCriticalFile() - better to allow than falsely block

---

### DisplaySecretWarning

```go
func DisplaySecretWarning()
```

**Purpose:** Display NON-BLOCKING warning for likely secret detection in user prompts (API keys, tokens, passwords)

**Parameters:** None - warning is generic (privacy-preserving, no details logged)

**Returns:** None - displays warning and returns immediately (NON-BLOCKING)

**Warning Behavior:**
- Displays prominent warning that potential secret detected
- Suggests user review prompt before submitting
- Does NOT block execution (awareness only)
- Does NOT log which pattern matched (privacy-preserving)
- Does NOT show any portion of secret content

**Example Usage:**

```go
import (
    "os"
    "hooks/lib/safety"
)

func main() {
    prompt := os.Getenv("USER_PROMPT")

    if safety.ContainsLikelySecret(prompt) {
        safety.DisplaySecretWarning()  // NON-BLOCKING awareness
        // Note: Detection uses detection.go, confirmation only displays warning
    }

    // Always continue - NON-BLOCKING design
    os.Exit(0)
}
```

**Warning Display:**
```
════════════════════════════════════════════════════════════
⚠️  POTENTIAL SECRET DETECTED IN PROMPT
════════════════════════════════════════════════════════════

Your prompt may contain API keys, tokens, or other sensitive credentials.

Please review before submitting:
  • API keys often start with patterns like "sk-", "AKIA", etc.
  • Consider if this information should be shared
  • Secrets in prompts may be logged or processed

This is a warning only - you may proceed if intentional.
```

**Privacy Preservation:**
- Function displays generic warning (no pattern details)
- Never logs what was detected or which pattern matched
- User must review prompt themselves to identify secret
- Detection happens, confirmation warns safely

**Design Philosophy:**
NON-BLOCKING because user may have legitimate reasons to include credentials in prompts (e.g., asking how to securely store a key they just generated). Warning provides awareness, user makes final decision.

---

## Integration Patterns

### Pattern 1: BLOCKING Pre-Tool Validation (Bash Commands)

**Hook:** `tool/pre-use.go`

**Context:** User attempts to execute Bash command via Claude Code

**Implementation:**

```go
package main

import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func main() {
    // Get command from environment
    cmd := os.Getenv("BASH_COMMAND")
    if cmd == "" {
        os.Exit(0)  // Not a bash command, allow
    }

    // Get temporal context from session-aware code
    timeContext := determineTimeContext()  // "long_session", "late_night", etc.

    // Request confirmation for dangerous operations
    needs, allowed := safety.ConfirmBashOperation(cmd, timeContext)

    // Make BLOCKING decision
    if needs && !allowed {
        fmt.Println("\n❌ Operation cancelled by user")
        fmt.Println("   Command not executed")
        os.Exit(1)  // BLOCK execution
    }

    // Allow operation to proceed
    os.Exit(0)
}

func determineTimeContext() string {
    // Implement session duration and time-of-day logic
    // Return: "long_session", "late_night", "very_early", "normal", or ""
    return ""  // Simplified for example
}
```

**Flow:**
1. Hook intercepts Bash command before execution
2. Calls ConfirmBashOperation with command and temporal context
3. If dangerous and not confirmed: Exit 1 (BLOCK)
4. If safe or confirmed: Exit 0 (ALLOW)

---

### Pattern 2: BLOCKING Pre-Tool Validation (File Writes)

**Hook:** `tool/pre-use.go`

**Context:** User attempts to write file via Claude Code

**Implementation:**

```go
package main

import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func main() {
    // Get file path from environment
    filePath := os.Getenv("FILE_PATH")
    if filePath == "" {
        os.Exit(0)  // Not a file write, allow
    }

    // Get temporal context
    timeContext := determineTimeContext()

    // Request confirmation for critical file writes
    needs, allowed := safety.ConfirmFileWrite(filePath, timeContext)

    // Make BLOCKING decision
    if needs && !allowed {
        fmt.Println("\n❌ Write operation cancelled by user")
        fmt.Println("   File not modified:", filePath)
        os.Exit(1)  // BLOCK execution
    }

    // Allow write to proceed
    os.Exit(0)
}
```

**Flow:**
1. Hook intercepts file write before execution
2. Calls ConfirmFileWrite with path and temporal context
3. If critical and not confirmed: Exit 1 (BLOCK)
4. If normal or confirmed: Exit 0 (ALLOW)

---

### Pattern 3: NON-BLOCKING Secret Detection (Prompt Submit)

**Hook:** `prompt/submit.go`

**Context:** User submits prompt to Claude Code

**Implementation:**

```go
package main

import (
    "os"
    "hooks/lib/safety"
)

func main() {
    // Get user prompt from environment
    prompt := os.Getenv("USER_PROMPT")
    if prompt == "" {
        os.Exit(0)
    }

    // Check for likely secrets (using detection.go)
    if safety.ContainsLikelySecret(prompt) {
        // Display NON-BLOCKING warning
        safety.DisplaySecretWarning()
        // Warning displayed, continue anyway (NON-BLOCKING)
    }

    // Always allow prompt submission (NON-BLOCKING design)
    os.Exit(0)
}
```

**Flow:**
1. Hook receives user prompt before submission
2. Checks for likely secrets via detection.ContainsLikelySecret()
3. If detected: Display warning (awareness)
4. Always exit 0 (NON-BLOCKING)

**Why NON-BLOCKING:**
User may have legitimate reasons to include credentials in prompts. Blocking would prevent valid use cases. Warning provides awareness, user decides.

---

### Pattern 4: Combined Safety Validation

**Hook:** `tool/pre-use.go` (comprehensive safety)

**Context:** Validate both Bash commands AND file writes with temporal awareness

**Implementation:**

```go
package main

import (
    "fmt"
    "os"
    "hooks/lib/safety"
)

func main() {
    timeContext := determineTimeContext()

    // Check Bash commands
    if cmd := os.Getenv("BASH_COMMAND"); cmd != "" {
        needs, allowed := safety.ConfirmBashOperation(cmd, timeContext)
        if needs && !allowed {
            fmt.Println("❌ Bash operation cancelled")
            os.Exit(1)  // BLOCK
        }
    }

    // Check file writes
    if filePath := os.Getenv("FILE_PATH"); filePath != "" {
        needs, allowed := safety.ConfirmFileWrite(filePath, timeContext)
        if needs && !allowed {
            fmt.Println("❌ File write cancelled")
            os.Exit(1)  // BLOCK
        }
    }

    // All checks passed
    os.Exit(0)
}

func determineTimeContext() string {
    // Session duration logic
    sessionHours := getSessionDuration()  // Hypothetical function
    if sessionHours >= 4 {
        return "long_session"
    }

    // Time of day logic
    hour := time.Now().Hour()
    if hour >= 22 || hour < 6 {
        if hour >= 22 {
            return "late_night"
        }
        return "very_early"
    }

    return "normal"
}
```

**Flow:**
1. Determine temporal context once
2. Validate Bash commands (if present)
3. Validate file writes (if present)
4. Block if either check fails confirmation
5. Allow if all checks pass

---

## Modification Policy

### When to Modify This Library

**Appropriate Modifications:**

1. **Add new operation type routing** - New dangerous operation categories discovered
   - Example: Container operations (docker rm -f), cloud destructive commands (aws s3 rm --recursive)
   - Add new case in ConfirmBashOperation switch, create new display function

2. **Enhance warning displays** - Improve clarity without changing structure
   - Better explanations of specific dangers
   - Additional contextual information
   - Improved formatting using system/lib/display

3. **Refine temporal context integration** - Better session awareness
   - More nuanced time-of-day warnings
   - Session fatigue indicators
   - Pattern recognition (repeated dangerous operations)

4. **Fix bugs in response validation** - Ensure exact match logic works correctly
   - Edge cases in input trimming
   - Consistent "yes" vs "y" requirements

**Inappropriate Modifications:**

1. **Changing conservative validation** - Don't make response matching more permissive
   - ❌ Accepting "yeah", "yep", "sure" → Breaks intentional confirmation requirement
   - ✅ Keep exact "yes" or "y" validation

2. **Adding confirmation history/tracking** - Preserve stateless design
   - ❌ "User confirmed 3 rm -rf in a row, auto-deny next" → Adds state complexity
   - ✅ Each operation confirmed independently

3. **Blocking secret detection** - Keep NON-BLOCKING design
   - ❌ Making DisplaySecretWarning block execution → Prevents legitimate use
   - ✅ Warning only, user decides

4. **Logging secret content** - Never compromise privacy preservation
   - ❌ Adding logs of which pattern matched → Privacy violation
   - ✅ Generic warning, no pattern details

### Coordination with Detection Library

Confirmation depends on detection (IsDangerousOperation, IsCriticalFile, ContainsLikelySecret). Changes to detection patterns may require confirmation updates:

**Detection adds new pattern → Confirmation routing:**
- If new dangerous operation type: Add routing case + display function
- If expanding existing category: Confirmation works automatically

**Example:**
```
Detection adds: "git push --force-with-lease" to dangerous operations
Confirmation response: Add case to force push routing OR create new display function
```

**Coordination Principle:** Detection identifies WHAT is dangerous. Confirmation explains WHY and HOW to user. Keep this separation clear.

---

## Future Roadmap

### Planned Enhancements

**1. Adaptive Confirmation Severity**

**Current:** Fixed "yes" or "y" requirements per operation type

**Future:** Context-aware severity adjustment
- Late night + long session + force push → Require typing full operation ("git push --force origin main")
- Normal hours + fresh session + sudo → Quick "y" confirmation sufficient
- Repeated same operation in short time → Escalate to full command typing

**Reasoning:** Confirmation severity should match risk AND user state. Tired users need more friction.

**2. Operation Context Display**

**Current:** Show command/path + temporal context

**Future:** Show git branch, affected files, recent history
- Force push: Show commits about to be overwritten
- rm -rf: Show directory contents preview (file count, size)
- Critical file write: Show current file diff (before/after)

**Reasoning:** Better decisions require understanding full impact.

**3. Confirmation Pattern Learning**

**Current:** Stateless - each confirmation independent

**Future:** Learn user patterns for enhanced warnings (NOT auto-allow)
- User never force pushes to main → Extra warning if attempting
- User frequently modifies ~/.bashrc → Less severe warning (familiar operation)
- User making first sudo operation → Educational note about implications

**Reasoning:** Personalized warnings more effective than generic (still require confirmation).

**4. Multi-Step Confirmations**

**Current:** Single yes/no decision

**Future:** Progressive disclosure for complex operations
1. "This will force push to main - understand implications?" [yes/no]
2. "You will overwrite these commits: [list]" [continue/cancel]
3. "Team members may be affected - confirmed?" [type branch name to confirm]

**Reasoning:** Most dangerous operations deserve graduated confirmation, not single gate.

**5. Confirmation Templates**

**Current:** Hardcoded warning displays

**Future:** Configurable warning templates (JSONC)
- Customize warnings per team/project
- Add company-specific safety notes
- Localization support (multiple languages)

**Reasoning:** Different contexts need different communication styles.

### Non-Goals

**Will NOT implement:**

1. **Auto-allow based on history** - Always require explicit confirmation
2. **Confirmation bypass flags** - No "skip confirmation" environment variables
3. **Silent confirmations** - Always display warning, never assume
4. **Machine-readable-only responses** - Always human-interactive

**Reasoning:** Confirmation exists to create pause for reflection. Shortcuts defeat purpose.

---

## Appendix: Warning Display Reference

### Force Push Warning

**Trigger:** `git push --force` or `git push -f`

**Severity:** HIGH - Requires "yes"

**Warning Content:**
- Explains: Rewrites remote history, affects team members
- Shows: Full command, branch affected
- Context: Temporal awareness note if applicable
- Guidance: Coordinate with team before force pushing

---

### Hard Reset Warning

**Trigger:** `git reset --hard`

**Severity:** HIGH - Requires "yes"

**Warning Content:**
- Explains: Discards ALL uncommitted changes permanently
- Shows: Full command
- Context: Temporal awareness note
- Guidance: Consider stashing changes first

---

### rm -rf Warning

**Trigger:** `rm -rf` or `rm -r`

**Severity:** HIGH - Requires "yes"

**Warning Content:**
- Explains: Permanently deletes files/directories recursively
- Shows: Full command
- Context: Temporal awareness note
- Guidance: Verify path before confirming

---

### Sudo/Pkexec Warning

**Trigger:** `sudo` or `pkexec` commands

**Severity:** MEDIUM - Requires "y"

**Warning Content:**
- Explains: Elevated privileges can modify system
- Shows: Full command
- Context: Temporal awareness note
- Icon: 🔐 (security implications)

---

### Publish Warning

**Trigger:** `npm publish`, `cargo publish`, `pip upload`

**Severity:** MEDIUM - Requires "y"

**Warning Content:**
- Explains: Publishes code to public registry
- Shows: Full command
- Context: Temporal awareness note
- Icon: 📦 (package publication)

---

### Database Warning

**Trigger:** `DROP DATABASE`, `DROP TABLE`, `TRUNCATE TABLE`

**Severity:** HIGH - Requires "yes"

**Warning Content:**
- Explains: Permanently deletes database/table data
- Shows: Full SQL command
- Context: Temporal awareness note
- Guidance: Verify backups exist

---

### Critical File Warning

**Trigger:** System config, authentication, git config paths

**Severity:** HIGH - Requires "yes"

**Warning Content:**
- Explains: Modifying critical system/config file
- Shows: Exact file path
- Context: Temporal awareness note
- Warning: Errors can break system/authentication

---

### Secret Detection Warning

**Trigger:** API keys, tokens, private keys in text

**Severity:** AWARENESS - NON-BLOCKING

**Warning Content:**
- Explains: Potential secret detected (generic)
- Guidance: Review prompt before submitting
- Privacy: NO pattern details, NO content logged
- Decision: User decides whether to proceed

---

**End of Safety Confirmation Library API Documentation**
