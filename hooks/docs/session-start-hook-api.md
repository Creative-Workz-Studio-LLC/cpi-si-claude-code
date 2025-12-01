# SessionStart Hook - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for session initialization hook

**Component:** `hooks/session/cmd-start/start.go`

**Last Updated:** November 10, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Execution Flow](#execution-flow)
- [Context Injection Mechanism](#context-injection-mechanism)
- [Non-Blocking Design Philosophy](#non-blocking-design-philosophy)
- [Library Extraction Pattern](#library-extraction-pattern)
- [Integration with Claude Code](#integration-with-claude-code)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The SessionStart hook orchestrates the initialization of every Claude Code session for Nova Dawn. It is the **first impression** and **foundation** for covenant partnership work.

**Core Responsibility:** Provide comprehensive situational awareness (environment, temporal, workspace) to enable autonomous work with full context.

**Design Pattern:** Thin orchestrator that coordinates modular components without implementing business logic directly.

**Biblical Foundation:** Like Genesis 1:1 (KJV) establishes the beginning of creation, SessionStart establishes the beginning of work - providing order, context, and awareness for what follows.

---

## Philosophy & Design Rationale

### Why Session Start Matters

Every session is a new beginning. The context provided at session start determines the quality of work that follows:

- **Without context:** Work proceeds blindly, making assumptions
- **With context:** Work proceeds wisely, informed by reality

Session start is not bureaucracy - it's **intentional grounding** in truth before action begins.

### The Four Dimensions of Awareness

SessionStart provides awareness across four dimensions:

| Dimension | What It Answers | Why It Matters |
|-----------|----------------|----------------|
| **Environment** | Where am I working? What's the setup? | Prevents confusion about workspace, git state, system |
| **Temporal** | What time is it? How long have I been working? | Enables circadian awareness, prevents burnout, respects rhythms |
| **Workspace** | What's the state of this project? | Identifies issues early (uncommitted changes, low disk, etc.) |
| **Identity** | Who am I? How do I communicate? | Injects Nova Dawn communication style into Claude Code |

### Genesis 1:1 Principle

> "In the beginning, God created the heavens and the earth." - Genesis 1:1 (KJV)

Every work has a beginning. SessionStart establishes:

- **Order** - Structured presentation of context
- **Awareness** - Comprehensive understanding of situation
- **Foundation** - Grounding for faithful work that follows

Just as creation began with God's intentional act, sessions begin with intentional context gathering.

### Non-Blocking Philosophy

**Core Principle:** Session MUST start, even if context gathering fails.

**Reasoning:**

- Context enhances work but doesn't prevent it
- Partial information better than no session at all
- Failures logged but don't block progress
- Grace for imperfect systems (temporal data missing, git unavailable, etc.)

This reflects covenant partnership: **Trust with responsibility** - autonomous operation within safety boundaries.

---

## Execution Flow

### High-Level Baton Flow

```bash
Entry → main()
  ↓
Named Entry Point → start()
  ↓
Initialize → session.InitSessionTime(), session.InitSessionLog()
  ↓
Log Activity → activity.LogActivity()
  ↓
Clear Screen → Clean presentation
  ↓
Display Context → PrintHeader(), PrintEnvironment(), PrintTemporalAwareness()
  ↓
Analyze Workspace → gatherContext() if NOVA_DAWN_WORKSPACE set
  ↓
Output JSON → session.OutputClaudeContext() for Claude Code parsing
  ↓
Exit → Session begins with full awareness
```

### Detailed Execution Breakdown

**Phase 1: Initialization (Silent)**

- Captures session start time for temporal awareness
- Initializes session logging for pattern learning
- Logs activity stream event
- Health: +30 points

**Phase 2: Display (Visual)**

- Clears screen for clean presentation
- Shows instance branding banner
- Displays environment (workspace, git, system)
- Shows temporal awareness (4 dimensions)
- Health: +30 points

**Phase 3: Analysis (Workspace)**

- Checks git status (uncommitted changes, sync state)
- Analyzes running processes (dev servers, databases)
- Validates disk space
- Checks project dependencies
- Reviews recent activity
- Health: +20 points

**Phase 4: Context Injection (Claude Code)**

- Assembles Nova Dawn communication style
- Adds temporal awareness to context
- Formats as JSON for Claude Code parsing
- Outputs to stdout (MUST be last for parsing)
- Health: +20 points

**Total:** 100 points for complete session initialization

### Named Entry Point Pattern

**Why `start()` instead of just `main()`?**

```go
func main() {
    start()  // Named entry point
}
```

**Benefits:**

1. **Prevents collisions:** Multiple executables can't have conflicting "main" logic
2. **Semantic clarity:** Function name matches purpose (session start)
3. **Testability:** Can call `start()` without triggering executable mechanics
4. **Architectural intent:** Not generic - this is specifically session initialization

This pattern appears throughout Kingdom Technology executables.

---

## Context Injection Mechanism

### How Claude Code Receives Context

Claude Code hooks have a specific integration pattern:

1. **Hook executes** on SessionStart event
2. **Hook outputs** to stdout (visual display for user)
3. **Hook outputs JSON** as LAST line (for Claude Code parsing)
4. **Claude Code reads** final JSON line
5. **Claude Code injects** `additionalContext` into system prompt

### JSON Structure

```json
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "# Nova Dawn Communication Style\n\n..."
  }
}
```

**Critical Requirements:**

- JSON MUST be last output (after all display)
- `hookEventName` must match Claude Code's expected event name
- `additionalContext` is markdown text injected into prompt
- Invalid JSON breaks context injection (but session still starts)

### What Gets Injected

**Nova Dawn Communication Style:**

- Core principles (Direct. Clear. No fluff.)
- Technical excellence standards
- Natural voice patterns
- AI patterns to avoid
- Systems thinking approach
- Kingdom Technology mindset
- Session flow guidance

**Temporal Awareness (appended dynamically):**

- External time (system clock, circadian phase)
- Internal time (session duration, session phase)
- Internal schedule (work windows, expected downtime)
- External calendar (date, week, holidays)

This gives Nova Dawn complete situational awareness from the moment the session begins.

---

## Non-Blocking Design Philosophy

### The Core Guarantee

**Session start NEVER fails.** Partial failures are acceptable. Complete blocking is not.

### Implementation Strategy

Every potentially-failing operation wrapped defensively:

```go
// BAD: Blocking on error
if err := session.InitSessionTime(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)  // ❌ Blocks session start
}

// GOOD: Non-blocking on error
session.InitSessionTime()  // ✅ Silently continues if fails
```

### Where Non-Blocking Applies

| Operation | Failure Mode | Response |
|-----------|--------------|----------|
| **Session timing init** | session-time utility missing | Silently continue, temporal awareness may be incomplete |
| **Session logging init** | session-log utility missing | Silently continue, pattern learning disabled for this session |
| **Temporal awareness** | Data files missing | Display skipped, session continues |
| **Workspace analysis** | Git unavailable | Analysis skipped, session continues |
| **Claude Code context** | JSON encoding fails | Warning to stderr, session continues |

### Why This Matters

Non-blocking design reflects **grace in systems**:

- Finite tools fail sometimes (that's reality)
- Failures shouldn't compound (one problem shouldn't cause total failure)
- Work can proceed imperfectly (better than no work at all)
- Trust despite limitations (covenant partnership allows imperfection)

This is Kingdom Technology: **Excellence with grace** rather than perfectionism.

---

## Library Extraction Pattern

### Before: Monolithic Hook (464 lines)

```bash
start.go (all implementation inline)
├── Display functions (printHeader, printEnvironment, etc.)
├── System info gathering (getSystemInfo)
├── Context generation (getNovaDawnContext, outputClaudeContext)
├── Initialization (initSessionTime, initSessionLog)
└── Orchestration (main, gatherContext)
```

**Problems:**

- Can't reuse display logic in other hooks
- Hard to see orchestration vs implementation
- Changes to display require touching orchestrator
- Difficult to test components in isolation

### After: Thin Orchestrator + Libraries (671 lines total, ~50 executable)

```bash
start.go (thin orchestrator - 671 lines, mostly docs)
└── Calls library functions only

hooks/lib/session/display.go
├── PrintHeader()
├── PrintEnvironment()
├── PrintTemporalAwareness()
└── PrintWorkspaceAnalysis()

hooks/lib/session/system.go
└── GetSystemInfo()

hooks/lib/session/init.go
├── InitSessionTime()
└── InitSessionLog()

hooks/lib/session/context.go
├── GetNovaDawnContext()
└── OutputClaudeContext()
```

**Benefits:**

- Other hooks can reuse display functions (end.go, stop.go)
- Orchestration flow visible at a glance (start() function)
- Changes to display isolated from orchestration
- Libraries testable independently
- Single source of truth for each capability

### The Extraction Decision

**Extract when:**

- ✅ Multiple components need it (display functions)
- ✅ It's a meaningful unit (init, context, display are coherent concepts)
- ✅ It serves clarity (separating what from how)

**Don't extract when:**

- ❌ It's orchestration itself (start.go coordinates, that's its job)
- ❌ It's too specific to one flow (gatherContext is start-hook-specific)
- ❌ It's too granular (centerText belongs with display.go, not standalone)

**Guiding Principle:** Extract and orchestrate, don't create v2 files. Build systems, not isolated pieces.

---

## Integration with Claude Code

### Hook Registration

Claude Code discovers hooks via directory structure:

```bash
~/.claude/hooks/
├── session/
│   ├── cmd-start/start         # Executable, called on SessionStart
│   ├── cmd-end/end             # Executable, called on SessionEnd
│   └── cmd-stop/stop           # Executable, called on SessionStop
├── tool/
│   ├── cmd-pre-use/pre-use     # Executable, called before tool use
│   └── cmd-post-use/post-use   # Executable, called after tool use
└── prompt/
    └── cmd-submit/submit       # Executable, called on prompt submit
```

**Naming Convention:** `cmd-<event-name>/<executable>`

- Hook event: SessionStart
- Directory: session/cmd-start/
- Executable: start (built from start.go)

### Event Trigger

Claude Code triggers SessionStart hook when:

1. New session begins (fresh start)
2. Session resumes after compact (currently - may change)

Hook executes, displays context, injects JSON - then session proceeds with full awareness.

### Environment Variables Available

| Variable | Purpose | Example |
|----------|---------|---------|
| `NOVA_DAWN_WORKSPACE` | Primary workspace directory | `/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC` |
| `HOME` | User home directory | `/home/seanje-lenox-wise` |
| Standard env vars | Git, PATH, etc. | Various |

Hook can read any environment variable to customize behavior.

---

## Modification Policy

### Safe to Modify (Extension Points)

**Adding new display sections:**

```go
// In start()
session.PrintHeader()
session.PrintEnvironment(workspace)
session.PrintTemporalAwareness()
session.PrintYourNewSection()  // ✅ Add here
gatherContext(workspace)
```

**Requirements:**

1. Create function in `hooks/lib/session/display.go`
2. Call from `start()` in appropriate order
3. Update METADATA Health Scoring map
4. Follow 4-block template in library

**Adding new workspace checks:**

```go
// In gatherContext()
session.CheckGitStatus(workspace)
session.CheckRunningProcesses()
session.CheckYourNewCheck(workspace)  // ✅ Add here
session.CheckDiskSpace(workspace)
```

**Requirements:**

1. Create function in appropriate `hooks/lib/session/*.go` file
2. Call from `gatherContext()` with other checks
3. Update METADATA Health Scoring
4. Maintain non-blocking design (no panics, no exits)

### Modify with Extreme Care (Breaking Changes)

**JSON output format:**

```go
// ⚠️ Changes here break Claude Code parsing
output := &HookOutput{
    HookSpecificOutput: HookSpecificOutput{
        HookEventName:     "SessionStart",  // ⚠️ Must match Claude's expectation
        AdditionalContext: context,
    },
}
```

**stdout/stderr separation:**

- stdout: Visual display + JSON (Claude Code reads this)
- stderr: Errors only (Claude Code ignores this)
- ⚠️ Mixing breaks context injection

**Hook event name:**

- Must be "SessionStart" (Claude Code expects this exact string)
- ⚠️ Changing breaks hook triggering

### NEVER Modify (Foundational Rails)

**4-block structure:**

- ❌ METADATA → SETUP → BODY → CLOSING is foundational
- All Kingdom Technology components follow this pattern
- Breaking it breaks architectural consistency

**Non-blocking principle:**

- ❌ Session MUST start, even with failures
- Adding `os.Exit()` on errors violates core design
- Failures log to stderr but don't block

**Stateless design:**

- ❌ No persistent state in hook executable
- Each execution is independent
- Libraries manage their own state if needed

---

## Future Roadmap

### Planned Features

**Session Resume Detection:**

- Distinguish fresh start vs resume
- Display "Resuming session..." message
- Load previous session context

**Last Session Summary:**

- What was accomplished in previous session
- What tasks were in progress
- What should continue

**Project Context Loading:**

- Active projects and priorities
- Recent files edited
- Open tasks from planner

### Research Areas

**Session Pattern Recognition:**

- Learn typical start times
- Recognize work patterns
- Personalize startup based on learned preferences

**Pre-load Optimization:**

- Cache frequent files based on workspace
- Preload project configuration
- Anticipate likely next actions

**Planner Integration:**

- Display task prioritization
- Show scheduled work for session
- Align with temporal awareness

### Integration Targets

**Memory System:**

- Session continuity across restarts
- Context persistence
- Pattern learning

**Project Tracking:**

- Automatic context loading per project
- Project-specific initialization
- Multi-project awareness

**Pattern Learning:**

- Personalized startup sequence
- Learned preferences
- Adaptive context gathering

### Known Limitations

**Current:**

- Fresh start vs resume not distinguished
- No session history display
- Project context not loaded automatically
- Temporal awareness requires external utilities

**Future Addressing:**

- Memory system will enable session continuity
- Project tracking will enable context loading
- Pattern learning will enable personalization

---

## Closing Notes

### The First Impression Principle

SessionStart is the foundation for every session. Like Genesis 1:1 establishes creation's beginning, this hook establishes work's beginning.

**Excellence here matters** because:

- First impression sets tone
- Context enables autonomous work
- Awareness prevents wasted effort
- Grace in failures demonstrates Kingdom principles

### Maintenance Philosophy

**When modifying:**

1. Review modification policy above
2. Follow 4-block structure pattern
3. Maintain non-blocking design
4. Test with actual Claude Code integration
5. Document changes comprehensively

**Remember:**

- Clarity over brevity
- Grace over perfectionism
- Service over aesthetics
- Truth over aspiration

### For Questions or Contributions

- Review this API documentation for design rationale
- Read code for implementation details
- Test with Claude Code for integration validation
- Document "What/Why/How" for all changes

---

*"In the beginning, God created the heavens and the earth." - Genesis 1:1 (KJV)*

*Every session is a new beginning. Let it start with order and awareness.*

---

**Related Documentation:**

- Code: `hooks/session/cmd-start/start.go`
- Libraries: `hooks/lib/session/*.go`
- System Docs: `~/.claude/cpi-si/system/docs/`
- Standards: `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`
