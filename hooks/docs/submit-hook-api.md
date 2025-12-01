# UserPromptSubmit Hook - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for NON-BLOCKING prompt secret detection

**Component:** `hooks/prompt/cmd-submit/submit.go`

**Last Updated:** November 10, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Execution Flow](#execution-flow)
- [Secret Detection Pattern](#secret-detection-pattern)
- [Non-Blocking Design Philosophy](#non-blocking-design-philosophy)
- [Safety Library Integration](#safety-library-integration)
- [Privacy-Preserving Design](#privacy-preserving-design)
- [Integration with Claude Code](#integration-with-claude-code)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The UserPromptSubmit hook provides NON-BLOCKING secret detection before prompts are submitted. It is the **vigilance checkpoint** that guards words without disrupting workflow.

**Core Responsibility:** Detect potential secrets (API keys, passwords, tokens) in prompts before submission, display warning to user for awareness, log activity for monitoring, NEVER block submission (always allows prompt through).

**Design Pattern:** Thin orchestrator that coordinates privacy-preserving logging, secret detection, warning display, and async monitoring - without implementing detection logic directly.

**Biblical Foundation:** "Set a watch, O LORD, before my mouth; keep the door of my lips" (Psalm 141:3 KJV). Submit hook establishes vigilance without control - awareness without restriction.

---

## Philosophy & Design Rationale

### Why Prompt Secret Detection Matters

Every prompt submission is trust moment. Without detection, secrets leak unintentionally:

- **Without detection:** API keys submitted in code snippets, passwords shared for "debugging", tokens exposed in examples
- **With detection:** User warned before submission, can review and decide, awareness raised about sensitive data

Prompt secret detection is not censorship - it's **faithful vigilance through awareness**.

### The Vigilance Principle (Psalm 141:3)

> "Set a watch, O LORD, before my mouth; keep the door of my lips." - Psalm 141:3 (KJV)

Guard words before they go forth. Submit hook establishes:

- **Vigilance** - Watch for secrets before submission, raise awareness
- **Awareness** - User knows potential secret detected, can review
- **Freedom** - Never blocks submission (user decides after warning)
- **Stewardship** - Care for sensitive data without controlling flow
- **Grace** - False positives acceptable (better to warn than miss)

Just as guarding mouth prevents harmful words, submit hook guards against accidental secret exposure.

### The Three Responsibilities of Detection

Submit hook balances three critical responsibilities:

| Responsibility | What It Means | Why It Matters |
|----------------|---------------|----------------|
| **Detect** | Recognize likely secrets (API keys, passwords, tokens) | Prevent accidental exposure |
| **Warn** | Display clear warning without blocking | User awareness without disruption |
| **Allow** | Never prevent prompt submission | Workflow continues, user decides |

### The Non-Blocking Philosophy

**Core Principle:** Warnings, not prevention. Awareness, not control.

**Reasoning:**

- Prompt submission is core workflow (can't disrupt)
- False positives inevitable (can't block on uncertainty)
- User judgment trumps automated detection
- Workflow disruption worse than occasional missed secret
- Trust requires non-interference

This reflects Kingdom Technology: **Serve through awareness, not control through restriction**.

### The Privacy-Preserving Dimension

Submit hook respects prompt privacy:

- **Activity logging** - Records prompt LENGTH only, never content
- **No storage** - Prompt not saved, not logged to files
- **Async monitoring** - Non-blocking background logging
- **User control** - User decides what to submit

Privacy awareness turns detection from "surveillance" into "faithful stewardship without invasion".

### The Warning Responsibility

**Critical Understanding:** This hook WARNS but NEVER blocks.

That design requires extraordinary restraint:

1. **Clear warnings** - User understands what detected and why
2. **Fast execution** - Detection can't slow submission flow
3. **Privacy preservation** - Never log prompt content
4. **Non-blocking** - Always exits 0 (allow)
5. **Conservative detection** - Better to warn than miss

Warning is servant posture. Warn faithfully, trust user decision.

---

## Execution Flow

### High-Level Flow

```bash
Claude Code prepares prompt submission
  ↓
Calls submit hook with PROMPT environment variable
  ↓
Hook orchestrates detection:
  1. Read prompt from environment
  2. Log activity (prompt length only, preserves privacy)
  3. Check for likely secrets (safety library detection)
  4. Display warning if secrets detected (non-blocking)
  5. Async monitoring log (background, doesn't block)
  6. Exit 0 (always allow submission)
  ↓
Claude Code proceeds with prompt submission
```

### Detailed Phase Breakdown

**Phase 1: Read Prompt (10 points)**

- Read `PROMPT` environment variable
- No validation (empty prompt acceptable)
- Health Impact: +10 points (always succeeds)

**Phase 2: Privacy-Preserving Activity Logging (20 points)**

- Calculate prompt length only
- Log "PromptSubmit" activity with length (not content)
- Preserves privacy (no prompt content in logs)
- Non-blocking (proceed if logging fails)
- Health Impact: +20 success, +0 failure

**Phase 3: Secret Detection (30 points)**

- Call `safety.ContainsLikelySecret(prompt)` from library
- Fast pattern matching (doesn't slow submission)
- Conservative detection (better to warn than miss)
- Health Impact: +30 points (detection always runs)

**Phase 4: Warning Display (20 points)**

- If secrets detected, call `safety.DisplaySecretWarning()`
- Non-blocking display (failure doesn't prevent submission)
- Clear warning encourages user review
- Health Impact: +20 success, +0 failure

**Phase 5: Async Monitoring Log (10 points)**

- Spawn goroutine for `monitoring.LogPrompt(prompt)`
- Background logging doesn't block submission
- Monitoring system tracks patterns asynchronously
- Health Impact: +10 points (goroutine spawn)

**Phase 6: Allow Submission (10 points)**

- Always exit 0 (allow prompt submission)
- User informed, user decides
- Hook serves through awareness, not control
- Health Impact: +10 points (always succeeds)

Total Health: 100 points for complete vigilance without disruption

---

## Secret Detection Pattern

### What Gets Detected

**Pattern Categories:**

| Pattern Type | Examples | Why Detected |
|--------------|----------|--------------|
| **API Keys** | `sk-...`, `api_key=...`, `token=...` | Likely authentication credentials |
| **Passwords** | `password=...`, `pwd=...`, `pass:...` | Likely authentication secrets |
| **Tokens** | Long alphanumeric strings, JWT-like patterns | Likely access tokens |
| **Private Keys** | `-----BEGIN PRIVATE KEY-----` | Cryptographic secrets |
| **Connection Strings** | Database URLs with credentials | Infrastructure secrets |

### Detection Library (safety/detection.go)

**Function:** `ContainsLikelySecret(content string) bool`

**Design:**

- Fast pattern matching (regex + heuristics)
- Conservative bias (better false positive than false negative)
- No prompt storage (detection in-memory only)
- Stateless (no learning, no tracking)

**Example Patterns:**

```go
// API Key patterns
sk-[a-zA-Z0-9]{32,}
api_key=[a-zA-Z0-9]{16,}

// Password patterns
password=.{6,}
pwd=.{6,}

// Token patterns
Bearer [a-zA-Z0-9\-._~+/]+=*

// Private key patterns
-----BEGIN (RSA|DSA|EC|OPENSSH|PGP) PRIVATE KEY-----
```

### Warning Display (safety/confirmation.go)

**Function:** `DisplaySecretWarning()`

**Display:**

```bash
⚠️  Potential secret detected in prompt
   Review prompt before submitting
```

**Design:**

- Non-blocking (display, then continue)
- Non-condescending (trusts user judgment)
- Clear enough to raise awareness
- Brief enough to not disrupt flow

---

## Non-Blocking Design Philosophy

### The Core Guarantee

**Prompt submission MUST proceed.** Detection warns, never prevents.

### Why Non-Blocking is Sacred

**Primary Reason:** Prompt submission is core workflow.

If hook blocks submission:

- User workflow broken on every false positive
- Trust in system destroyed
- Hook becomes obstacle instead of servant
- User disables hook entirely

**Secondary Reason:** False positives inevitable.

Pattern matching can't achieve certainty:

- Code snippets contain "password" as variable names
- Example data looks like real credentials
- Learning materials demonstrate secret patterns
- Legitimate content triggers detection

**Tertiary Reason:** User judgment is superior.

User knows context:

- Is this a real secret or example?
- Is exposure acceptable for this use case?
- Should I rephrase or proceed?

Detection serves awareness. Decision remains with user.

### Implementation Strategy

Every operation designed for zero blocking:

```go
// ✅ Always allow
func userPromptSubmit() {
    prompt := os.Getenv("PROMPT")

    // Privacy-preserving logging (length only)
    activity.LogActivity("PromptSubmit", "length:"+strconv.Itoa(len(prompt)), "success", 0)

    // Warn if secrets detected (but continue regardless)
    if prompt != "" && safety.ContainsLikelySecret(prompt) {
        safety.DisplaySecretWarning()
        // Note: No exit here, continues to allow
    }

    // Async monitoring (doesn't block)
    go monitoring.LogPrompt(prompt)

    // Always allow submission
    os.Exit(ExitAllow)  // ✅ Exit 0
}
```

### Exit Code Contract

**Exit 0 (Allow) - Always:**

- Secrets detected → warn → allow
- No secrets detected → allow
- Detection failed → allow
- Display failed → allow
- Logging failed → allow

**Exit 1 (Block) - Never:**

- Hook NEVER exits 1
- Non-blocking is foundational design
- Blocking would violate core principle

This contract is sacred - breaking it breaks user trust.

---

## Safety Library Integration

### Separation of Concerns

**Hook Responsibilities (submit.go):**

- Read prompt from environment
- Log activity (privacy-preserving)
- Route to detection library
- Display warning if detected
- Async monitoring
- Allow submission (exit 0)

**Library Responsibilities (hooks/lib/safety/):**

- Detect likely secrets (`detection.go`)
- Display warning (`confirmation.go`)
- Pattern matching logic
- No prompt storage or logging

### Detection Library (hooks/lib/safety/detection.go)

**Purpose:** Centralized secret pattern detection

**Function:**

- `ContainsLikelySecret(content string) bool` - Pattern matching for secrets

**Design:**

- Fast execution (can't slow submission)
- Conservative detection (warn rather than miss)
- Stateless (no tracking)
- Reusable across hooks

### Warning Library (hooks/lib/safety/confirmation.go)

**Purpose:** Display warning messages

**Function:**

- `DisplaySecretWarning()` - Display generic secret warning

**Design:**

- Non-blocking display
- Clear without condescension
- Trusts user judgment
- Consistent with other warnings

### Extraction Benefits

**For the Hook:**

- Stays thin and fast (orchestration only)
- Easy to understand (routing logic clear)
- Easy to modify (add new checks)

**For the Libraries:**

- Testable independently
- Reusable (other hooks can detect secrets)
- Consistent (warnings look same across system)
- Maintainable (update patterns centrally)

**For the System:**

- Surgical updates (change detection without changing hook)
- Consistent experience (same detection everywhere)
- Pattern evolution (detection improves based on use)

---

## Privacy-Preserving Design

### Why Privacy Matters

Prompts contain sensitive information:

- User thoughts and intentions
- Project-specific context
- Potentially confidential data
- Personal communication style

Submit hook respects privacy while enabling monitoring.

### Privacy Preservation Strategy

| Component | What's Logged | What's NOT Logged |
|-----------|---------------|-------------------|
| **Activity Log** | Prompt length, timestamp | Prompt content |
| **Detection** | In-memory only | Never stored |
| **Warning** | Generic message | What was detected |
| **Monitoring** | Async background logging | Blocking user flow |

### Implementation

**Activity Logging:**

```go
// ✅ Privacy-preserving - length only
promptLength := strconv.Itoa(len(prompt))
activity.LogActivity("PromptSubmit", "length:"+promptLength, "success", 0)

// ❌ NOT done - would violate privacy
// activity.LogActivity("PromptSubmit", prompt, "success", 0)
```

**Detection:**

```go
// ✅ In-memory detection only
if safety.ContainsLikelySecret(prompt) {
    // Detection happens, nothing stored
}
```

**Monitoring:**

```go
// ✅ Async monitoring (non-blocking)
go monitoring.LogPrompt(prompt)
// Monitoring system handles storage policies
```

### The Privacy Principle

**User data is sacred.** Detection serves user without invading privacy.

- Prompt content not in activity logs
- Detection ephemeral (in-memory only)
- Monitoring async (user not blocked)
- User controls what submits

Privacy preservation reflects covenant partnership - serve without surveillance.

---

## Integration with Claude Code

### Hook Registration

Claude Code registers `UserPromptSubmit` hook event to call `~/.claude/hooks/prompt/cmd-submit/submit` before every prompt submission.

### Call Pattern

```bash
Claude Code prepares prompt submission
  ↓
Spawns: submit
Environment: PROMPT=<user_prompt_content>
  ↓
Hook executes (non-blocking, always allows)
  ↓
Hook exits with code 0
  ↓
Claude Code proceeds with prompt submission
```

### Environment Variables

**Input to Hook:**

- `PROMPT`: Full prompt content from user (multi-line, can be very long)

**Hook doesn't set environment** - only reads what Claude Code provides

### Standard Streams

**stdin:** Not used (no user confirmation needed, non-blocking)
**stdout:** Warning displays (user sees these if secrets detected)
**stderr:** Not used (all output to stdout)

### Exit Codes

**0 (Allow):** Claude Code proceeds with prompt submission (ALWAYS)
**1 (Block):** Never used (non-blocking design)

No other exit codes used (always allows).

---

## Modification Policy

### SAFE TO MODIFY (Extension Points)

**Add new secret patterns:**

1. Add pattern to `hooks/lib/safety/detection.go` (`ContainsLikelySecret`)
2. Test pattern with sample data
3. Verify false positive rate acceptable
4. Hook stays unchanged (routing already handles library)

**Enhance warning display:**

1. Modify `DisplaySecretWarning()` in `hooks/lib/safety/confirmation.go`
2. Keep warning clear and brief
3. Avoid condescension
4. Hook stays unchanged

**Add monitoring capabilities:**

1. Enhance `monitoring.LogPrompt()` functionality
2. Maintain async execution (goroutine)
3. Preserve non-blocking design
4. Hook stays unchanged

### MODIFY WITH CARE (Structural Changes)

**Changing detection logic:**

- Test exhaustively (various prompt types: code, text, examples)
- Verify performance (can't slow submission)
- Measure false positive rate
- Update health scoring in METADATA if flow changes

**Adding new checks:**

- Verify non-blocking design maintained
- Test impact on submission speed
- Ensure privacy preservation
- Update health scoring map in METADATA

### NEVER MODIFY (Foundational Rails)

**4-block structure** - METADATA, SETUP, BODY, CLOSING pattern is sacred
**Non-blocking principle** - MUST always exit 0, never blocks submission
**Privacy preservation** - Activity logs never contain prompt content
**Thin orchestrator** - Logic belongs in libraries, not hook
**Exit 0 guarantee** - Breaking this breaks user trust

### Testing Requirements

Before deploying ANY changes:

1. Test with various prompt types (code snippets, text, examples)
2. Test with known secrets (API keys, passwords, tokens)
3. Test with false positives (legitimate content that looks like secrets)
4. Verify non-blocking (submission always proceeds)
5. Verify privacy (activity logs don't contain prompt content)
6. Measure execution time (< 50ms normal path)
7. Check async monitoring works (goroutine spawns correctly)

---

## Future Roadmap

### Planned Features

**Enhanced Secret Detection (High Priority):**

- Machine learning for pattern recognition
- Context-aware detection (code vs text)
- Confidence scoring for warnings
- User feedback integration

**Secret Type Classification (Medium Priority):**

- Identify specific secret types detected
- Display targeted warnings (e.g., "Potential API key detected")
- Suggest remediation (e.g., "Use environment variables")
- Link to best practices documentation

**User Feedback Loop (Medium Priority):**

- Allow user to mark false positives
- Learn from user corrections
- Improve detection patterns over time
- Build user-specific detection profiles

**Context-Aware Detection (Low Priority):**

- Distinguish real secrets from examples
- Recognize quoted/escaped patterns differently
- Understand code context vs narrative
- Reduce false positives intelligently

### Research Areas

**Machine Learning Detection:**

- Train model on known secret patterns
- Continuous learning from user feedback
- Confidence scoring for warnings
- Adaptive pattern recognition

**Remediation Suggestions:**

- Detect secret type, suggest solution
- "Use environment variable: OPENAI_API_KEY"
- "Store in .env file, add to .gitignore"
- Link to security best practices

**Integration with Secret Management:**

- Detect secrets, suggest secret manager
- Integration with 1Password, Bitwarden, etc.
- Automatic substitution suggestions
- Secret rotation reminders

### Known Limitations

**No secret type classification:** Generic warning, doesn't identify what detected
**No user feedback:** Can't learn from false positives
**No remediation suggestions:** Warns but doesn't guide
**No context awareness:** Code examples trigger same as real secrets
**No confidence scoring:** All detections treated equally

### Enhancement Opportunities

**Pre-submission secret scanning:**

- Integration with git pre-commit hooks
- Scan staged files for secrets
- Prevent commits with secrets

**Secret rotation reminders:**

- Track detected secret patterns
- Remind user to rotate regularly
- Integration with secret management

**Team coordination:**

- Shared secret detection patterns
- Team-specific rules
- Organizational policy enforcement

**Documentation integration:**

- Link warnings to security docs
- Explain why pattern dangerous
- Best practices for secret management

---

## Conclusion

The UserPromptSubmit hook embodies vigilance without control - warning faithfully while trusting user decision. It serves workflow without disrupting it.

**Key Principles:**

- Non-blocking design (warns, never blocks)
- Privacy preservation (logs length, not content)
- Fast execution (doesn't slow submission)
- Conservative detection (warn rather than miss)
- Thin orchestration (logic in libraries)

**Success Metrics:**

- Zero blocked submissions (non-blocking is absolute)
- Meaningful warnings (user understands what detected)
- < 50ms execution for normal path
- Privacy maintained (no prompt content in activity logs)
- User trust preserved (not disabled)

Kingdom Technology protects through awareness, not restriction. Vigilance serves the user, not controls them.

*"Set a watch, O LORD, before my mouth; keep the door of my lips." - Psalm 141:3 (KJV)*

Every prompt submission deserves faithful vigilance. Let warnings serve without restricting.

---

**Related Documentation:**

- Code: `hooks/prompt/cmd-submit/submit.go`
- Libraries: `hooks/lib/safety/detection.go`, `hooks/lib/safety/confirmation.go`, `hooks/lib/activity/`, `hooks/lib/monitoring/`
- Complementary Hooks: `hooks/docs/pre-use-hook-api.md`
- System Docs: `~/.claude/cpi-si/system/docs/`
- Standards: `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`
