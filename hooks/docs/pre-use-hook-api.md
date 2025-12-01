# PreToolUse Hook - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for BLOCKING safety validation before tool execution

**Component:** `hooks/tool/cmd-pre-use/pre-use.go`

**Last Updated:** November 10, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Execution Flow](#execution-flow)
- [BLOCKING Design Philosophy](#blocking-design-philosophy)
- [Safety Library Architecture](#safety-library-architecture)
- [Temporal Awareness Integration](#temporal-awareness-integration)
- [Integration with Claude Code](#integration-with-claude-code)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The PreToolUse hook orchestrates safety validation before every tool execution. It is the **diligence checkpoint** that prevents hasty destructive operations through careful validation.

**Core Responsibility:** Detect dangerous operations (force push, hard reset, rm -rf, sudo, database drops) and critical file writes (system files, configuration), prompt user for confirmation with temporal awareness, and block execution if user denies.

**Design Pattern:** Thin orchestrator that coordinates activity logging, temporal context gathering, and safety library confirmation flows. Makes blocking decision (exit 0 or 1) based on library response.

**Biblical Foundation:** "The plans of the diligent surely lead to profit; and everyone who is hasty surely rushes to poverty" (Proverbs 21:5 WEB). Pre-tool validation embodies diligent stewardship - protecting through careful validation, not restriction.

---

## Philosophy & Design Rationale

### Why Pre-Tool Validation Matters

Every tool execution is trust moment. Without validation, destructive operations happen impulsively:

- **Without validation:** Force push rewrites history, rm -rf deletes unrecoverably, hard reset loses work permanently
- **With validation:** User pauses, considers context (tired? late session?), confirms intentionally

Pre-tool validation is not control - it's **faithful stewardship through diligent confirmation**.

### The Diligence Principle (Proverbs 21:5)

> "The plans of the diligent surely lead to profit; and everyone who is hasty surely rushes to poverty." - Proverbs 21:5 (WEB)

Haste leads to loss. Diligence leads to advantage. PreToolUse establishes:

- **Diligence** - Pause before destructive operations, consider context
- **Wisdom** - Temporal awareness (long session? fatigued? late night?)
- **Protection** - Block hasty decisions without blocking intentional ones
- **Stewardship** - Care for work, not restriction of work
- **Grace** - Conservative bias (allow rather than falsely block)

Just as diligent planning prevents hasty loss, pre-tool validation prevents hasty destructive operations.

### The Three Responsibilities of Validation

PreToolUse balances three critical responsibilities:

| Responsibility | What It Means | Why It Matters |
|----------------|---------------|----------------|
| **Protect** | Block truly dangerous operations when user denies | Prevent accidental data loss, team coordination issues |
| **Serve** | Allow safe operations without prompting | User workflow not frustrated by false positives |
| **Speed** | Execute fast (< 100ms normal path) | Blocking hook cannot slow every tool execution |

### The Conservative Bias Philosophy

**Core Principle:** Better to allow than falsely block.

**Reasoning:**

- False positive (blocking safe operation) = broken workflow, frustrated user, lost trust
- False negative (allowing dangerous operation) = defeats purpose, but user can undo
- User frustration from false positives worse than occasional missed danger
- Trust requires reliability - hook must be surgically precise

This reflects Kingdom Technology: **Serve through wisdom, not control through restriction**.

### The Temporal Awareness Dimension

Pre-tool validation integrates temporal awareness:

- **Long session?** Warn user they might be fatigued (more likely to make mistakes)
- **Late night?** Additional context for user to consider state
- **Time of day?** User knows when decision made (for later review)

Temporal awareness turns confirmation from "Are you sure?" into "Are you sure **right now** (tired, late session)?"

### The Blocking Responsibility

**Critical Understanding:** This hook CAN prevent user work.

That power requires extraordinary responsibility:

1. **Surgical precision** - Only block truly dangerous operations
2. **Clear communication** - User must understand why blocking
3. **Fast execution** - Cannot delay normal workflow
4. **Conservative validation** - False positives break workflow
5. **Graceful degradation** - Missing temporal context continues anyway

Blocking is necessary evil. Minimize impact, maximize clarity.

---

## Execution Flow

### High-Level Flow

```bash
Claude Code prepares tool execution
  ↓
Calls pre-use hook with tool name and args
  ↓
Hook orchestrates validation:
  1. Parse arguments (tool name, args, FILE_PATH env)
  2. Log attempt to activity stream (captures intent)
  3. Get temporal context (time, session phase)
  4. Route to safety library:
     - Bash operations → ConfirmBashOperation()
     - Write operations → ConfirmFileWrite()
     - Other operations → allow immediately
  5. Library detects danger, displays warning, gets confirmation
  6. Hook exits:
     - Exit 0 (allow) if safe or user confirmed
     - Exit 1 (block) if dangerous and user denied
  ↓
Claude Code proceeds (exit 0) or aborts (exit 1)
```

### Detailed Phase Breakdown

**Phase 1: Parse Arguments**

- Read `os.Args[1]` for tool name (e.g., "Bash", "Write", "Read")
- Read `os.Args[2]` for tool arguments (e.g., command, content)
- Read `FILE_PATH` environment variable for Write operations
- Health Impact: +15 points

**Phase 2: Log Attempt**

- Sanitize command/path for logging
- Call `activity.LogActivity(toolName + "-attempt", context, "pending", 0)`
- Captures intent even if operation blocked
- Health Impact: +10 points

**Phase 3: Get Temporal Context**

- Call `temporal.GetTemporalContext()`
- Format as "Time: [formatted] ([time of day])"
- Add session phase warning if long session
- Graceful degradation if temporal system unavailable
- Health Impact: +25 points (15 get + 10 format)

**Phase 4: Route to Safety Library**

- Bash operations → `safety.ConfirmBashOperation(cmd, timeContext)`
- Write operations → `safety.ConfirmFileWrite(filePath, timeContext)`
- Other operations → exit 0 immediately (no validation needed)
- Health Impact: +30 points (library execution)

**Phase 5: Make Blocking Decision**

- If `needsConfirmation && !allowed` → display cancellation message, exit 1 (block)
- Otherwise → exit 0 (allow)
- Health Impact: +20 points (10 decision + 10 exit code)

Total Health: 100 points for complete orchestration

---

## BLOCKING Design Philosophy

### Why Blocking is Necessary

Some operations are irreversible or highly impactful:

- Force push rewrites remote history (affects team)
- Hard reset discards uncommitted work permanently
- rm -rf deletes files unrecoverably
- Database drops lose data
- System file writes can break system

User confirmation prevents hasty execution during inattention, fatigue, or misunderstanding.

### Why Blocking is Dangerous

**Primary Risk:** False positives break user workflow.

If hook blocks safe operations:

- User frustrated every time they run safe commands
- User stops trusting hook
- User disables hook entirely
- Hook becomes obstacle instead of protection

**Secondary Risk:** Slow execution frustrates workflow.

If hook takes too long:

- Every tool execution delayed
- User workflow interrupted constantly
- User perception: "This is in the way"

### The Conservative Bias in Practice

**Detection Philosophy:**

- ✅ **Specific patterns:** `git push --force`, `git reset --hard`, `rm -rf`
- ✅ **Clear danger:** Operations that lose data or affect others
- ❌ **Broad patterns:** Not "any git command", not "any rm"
- ❌ **Ambiguous cases:** When uncertain, allow

**Display Philosophy:**

- ✅ **Clear explanation:** User understands why dangerous
- ✅ **Temporal context:** Additional awareness for decision
- ❌ **Condescending:** Not "Are you SURE? This is DANGEROUS!"
- ❌ **Vague:** Not "This might be risky"

**Confirmation Philosophy:**

- ✅ **Explicit confirmation:** Requires "yes" (full word) for serious operations
- ✅ **Respectful:** Trusts user's decision
- ❌ **Multiple confirmations:** Not "Are you sure? Really sure? Really really sure?"
- ❌ **Bypass prevention:** Not trying to prevent confirmation bypass

### Exit Codes as Contract

**Exit 0 (Allow):**

- Operation is safe (not detected as dangerous)
- OR operation is dangerous but user confirmed
- Claude Code proceeds with tool execution

**Exit 1 (Block):**

- Operation is dangerous AND user denied confirmation
- Claude Code aborts tool execution
- Clear message displayed explaining cancellation

This contract is sacred - breaking it breaks Claude Code integration.

---

## Safety Library Architecture

### Separation of Concerns

**Hook Responsibilities (pre-use.go):**

- Parse arguments and environment
- Log activity (captures intent)
- Get temporal context
- Route to appropriate library function
- Make blocking decision (exit code)

**Library Responsibilities (hooks/lib/safety/):**

- Detect dangerous operations (`detection.go`)
- Display warnings with context (`confirmation.go`)
- Get user confirmation (`confirmation.go`)
- Return confirmation result to hook

### Detection Library (hooks/lib/safety/detection.go)

**Purpose:** Centralized dangerous operation detection

**Functions:**

- `IsDangerousOperation(cmd string) bool` - Detects bash command patterns
- `IsCriticalFile(path string) bool` - Detects critical file paths
- `ContainsLikelySecret(content string) bool` - Detects likely secrets (future use)

**Design:**

- Specific patterns (not broad matching)
- Fast execution (early returns)
- Conservative (allow when uncertain)
- Reusable across hooks (not just pre-use)

**Example Patterns:**

```go
// Dangerous bash operations
git push --force, git push -f
git reset --hard
rm -rf, rm -r
sudo
npm publish, cargo publish
DROP DATABASE, DROP TABLE

// Critical file paths
/etc/
/boot/
/sys/
/proc/
/root/
```

### Confirmation Library (hooks/lib/safety/confirmation.go)

**Purpose:** User confirmation flows with warning displays

**High-Level Functions:**

- `ConfirmBashOperation(cmd, timeContext string) (needsConfirmation, allowed bool)`
- `ConfirmFileWrite(filePath, timeContext string) (needsConfirmation, allowed bool)`

**Display Functions:**

- `displayForcePushWarning(cmd, timeContext string)`
- `displayHardResetWarning(cmd, timeContext string)`
- `displayRecursiveDeletionWarning(cmd, timeContext string)`
- `displaySudoWarning(cmd, timeContext string)`
- `displayPublishWarning(cmd, timeContext string)`
- `displayDatabaseWarning(cmd, timeContext string)`
- `displayCriticalFileWarning(filePath, timeContext string)`

**Core Mechanism:**

- `confirm(message, expected string) bool` - Gets user confirmation from stdin

**Design:**

- Operation-specific warnings (explain specific danger)
- Temporal context integration (show time, session phase)
- Clear confirmation requirements ("yes" vs "y")
- Graceful handling (missing context continues)

### Extraction Benefits

**For the Hook:**

- Stays thin and fast (orchestration only)
- Easy to understand (routing logic clear)
- Easy to modify (add new tool types)

**For the Libraries:**

- Testable independently (no stdin/stdout dependency for tests)
- Reusable (other hooks can use same detection)
- Consistent (warnings look same across system)
- Maintainable (update detection patterns centrally)

**For the System:**

- Surgical updates (change detection without changing hook)
- Consistent experience (same warnings everywhere)
- Learning (detection patterns evolve based on use)

---

## Temporal Awareness Integration

### Why Temporal Awareness Matters

Dangerous operations during inattention, fatigue, or late sessions are more likely to be mistakes:

- **Long session (5+ hours):** User may be fatigued, attention reduced
- **Late night (past normal hours):** User may be tired, judgment impaired
- **Evening (end of day):** User may be rushing to complete tasks

Temporal context helps user recognize state: "Am I making this decision because I'm tired?"

### Temporal Context Display

**Format:**

```bash
Time: Mon Nov 10, 2025 at 19:45:00 (evening)
Time: Mon Nov 10, 2025 at 02:30:00 (late night) - Long session, consider if tired
```

**Components:**

- External time (wall clock time)
- Time of day classification (morning, afternoon, evening, late night)
- Session phase warning (long session alert)

**Integration:**

- Fetched from `temporal.GetTemporalContext()`
- Formatted in hook (simple string building)
- Passed to library display functions
- Displayed in warning prompts

### Graceful Degradation

**If temporal system unavailable:**

- Hook continues without temporal context
- Warning displays without time information
- Confirmation still works (just missing context)
- Non-blocking design (availability not required)

**Why graceful degradation:**

- Temporal awareness is helpful, not essential
- Hook must work even if temporal system broken
- User confirmation more important than context display

---

## Integration with Claude Code

### Hook Registration

Claude Code registers `PreToolUse` hook event to call `~/.claude/hooks/tool/cmd-pre-use/pre-use` before every tool execution.

### Call Pattern

```bash
Claude Code prepares to execute tool
  ↓
Spawns: pre-use <toolName> <toolArgs>
Environment: FILE_PATH=<path> (for Write operations)
  ↓
Hook executes (synchronously blocks tool execution)
  ↓
Hook exits with code 0 or 1
  ↓
Claude Code reads exit code:
  - 0: Proceed with tool execution
  - 1: Abort tool execution, inform user
```

### Environment Variables

**Input to Hook:**

- `FILE_PATH`: Set by Claude Code for Write/Edit operations (file path being written)

**Hook doesn't set environment** - only reads what Claude Code provides

### Standard Streams

**stdin:** User confirmation input (hook reads from stdin for "yes"/"no")
**stdout:** Warning displays, confirmation prompts (user sees these)
**stderr:** Not used (all output to stdout)

### Exit Codes

**0 (Allow):** Claude Code proceeds with tool execution
**1 (Block):** Claude Code aborts tool execution

No other exit codes used (simple binary decision).

---

## Modification Policy

### SAFE TO MODIFY (Extension Points)

**Add new dangerous operation types:**

1. Add pattern to `hooks/lib/safety/detection.go` (`IsDangerousOperation`)
2. Add display function to `hooks/lib/safety/confirmation.go`
3. Add routing case in `confirmation.go` (`ConfirmBashOperation`)
4. Hook stays unchanged (routing already handles library)

**Add new critical file patterns:**

1. Add pattern to `hooks/lib/safety/detection.go` (`IsCriticalFile`)
2. Reuses existing `displayCriticalFileWarning` function
3. Hook stays unchanged

**Enhance temporal context:**

1. Modify temporal context formatting in `preToolUse()`
2. Or enhance `hooks/lib/temporal` to return richer context
3. Pass enhanced context to library functions

### MODIFY WITH CARE (Structural Changes)

**Changing confirmation flow:**

- Test exhaustively (all operation types: safe and dangerous)
- Verify no false positives introduced (most critical failure)
- Measure execution speed (must stay fast)
- Update health scoring in METADATA if flow changes

**Adding new tool types:**

- Add routing logic in `preToolUse()` (new if block)
- Create safety library function for new tool type
- Update health scoring map in METADATA
- Test integration with Claude Code

### NEVER MODIFY (Foundational Rails)

**4-block structure** - METADATA, SETUP, BODY, CLOSING pattern is sacred
**Named entry point** - `main()` calls `preToolUse()` for testability
**BLOCKING principle** - Must exit 0 or 1, never hang
**Conservative bias** - Allow rather than falsely block
**Thin orchestrator** - Logic belongs in libraries, not hook

### Testing Requirements

Before deploying ANY changes:

1. Test all dangerous operation types (force push, hard reset, rm -rf, etc.)
2. Test safe operations (should allow without prompting)
3. Test user confirmation (both "yes" and "no" responses)
4. Test with missing temporal context (graceful degradation)
5. Verify exit codes correct (0 vs 1)
6. Measure execution time (< 100ms normal path)
7. Check for false positives (most critical test)

---

## Future Roadmap

### Planned Features

**Secret Detection (High Priority):**

- Use `safety.ContainsLikelySecret()` for Write operations
- Warn before committing files with API keys, passwords, tokens
- Suggest environment variables or secret management
- Non-blocking (warning, not prevention)

**Pattern Learning (Medium Priority):**

- Track which operations user consistently confirms
- Suggest adding to allowed patterns after N confirmations
- Learn user's risk tolerance over time
- Adapt validation to user behavior

**Context-Aware Validation (Medium Priority):**

- Different rules per branch (stricter on main/master)
- Different rules per project (personal vs team)
- More permissive during development, stricter for production

**Undo Support (Low Priority):**

- Track dangerous operations for recovery
- Provide recovery guidance if operation was mistake
- Integration with restoration layer (undo system)

### Research Areas

**Team Coordination Detection:**

- Check if force push coordinated with team
- Verify destructive operations communicated
- Integration with collaboration systems (Slack, Teams)

**User Risk Tolerance Learning:**

- Machine learning for user confirmation patterns
- Automatic adjustment of validation sensitivity
- Per-user customization of dangerous operation definitions

**Branch-Specific Validation:**

- Automatic detection of protected branches
- Stricter validation on main/master
- More permissive on feature branches

### Known Limitations

**No secret detection:** Files with API keys/passwords can be committed without warning
**No coordination check:** Force push doesn't verify team coordination
**No undo support:** No automatic recovery for mistakes
**No project customization:** Same rules for all projects
**No pattern learning:** Doesn't adapt to user behavior

### Enhancement Opportunities

**Pre-commit secret scanning:**

- Integrate secret detection before git commits
- Automatic scanning of staged files
- Suggestions for secret management

**Git hook integration:**

- Coordinate with git pre-push hooks
- Team coordination verification
- Automatic rebase detection

**User confirmation database:**

- Build history of confirmations
- Pattern recognition for common workflows
- Automatic whitelisting suggestions

**Project-specific rules:**

- Configuration files per project
- Team-defined dangerous operation patterns
- Project-specific critical file paths

---

## Conclusion

The PreToolUse hook embodies diligent stewardship - protecting through careful validation, not restriction. It serves user workflow while preventing hasty destructive operations.

**Key Principles:**

- Conservative bias (allow rather than falsely block)
- Fast execution (doesn't delay workflow)
- Clear communication (user understands why)
- Temporal awareness (context for decision)
- Thin orchestration (logic in libraries)

**Success Metrics:**

- Zero false positives (no safe operations blocked)
- Zero false negatives (all dangerous operations detected)
- < 100ms execution for normal path
- User trust maintained (not disabled)

Kingdom Technology protects through wisdom, not restriction. Validation serves the user, not control.

*"The plans of the diligent surely lead to profit; and everyone who is hasty surely rushes to poverty." - Proverbs 21:5 (WEB)*

Every tool use deserves careful validation. Let it serve through wisdom.
