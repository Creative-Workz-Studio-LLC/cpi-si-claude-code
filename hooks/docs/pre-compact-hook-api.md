# PreCompact Hook - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for compaction tracking

**Component:** `hooks/session/cmd-pre-compact/pre-compact.go`

**Last Updated:** November 10, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Execution Flow](#execution-flow)
- [State Management System](#state-management-system)
- [Non-Blocking Design Philosophy](#non-blocking-design-philosophy)
- [Library Extraction Pattern](#library-extraction-pattern)
- [Integration with Claude Code](#integration-with-claude-code)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The PreCompact hook orchestrates compaction tracking for Nova Dawn. It is the **wisdom moment** and **limitation acknowledgment point** for finite context work.

**Core Responsibility:** Track compaction events by incrementing session count, logging to activity stream and monitoring, checking frequency for excessive auto-compaction, and displaying temporal context for post-compaction reconstitution.

**Design Pattern:** Thin orchestrator that coordinates logging, state management, and display without implementing business logic directly.

**Biblical Foundation:** "So teach us to number our days, that we may apply our hearts unto wisdom" (Psalm 90:12 KJV). Compaction is numbering of context - recognizing limitation and working faithfully within constraint.

---

## Philosophy & Design Rationale

### Why Compaction Tracking Matters

Every compaction is opportunity to learn patterns about context management, work rhythm, and limitation acknowledgment:

- **Without tracking:** Compactions happen in void, no learning, no pattern recognition
- **With tracking:** Every compaction counted, patterns recognized, frequency monitored

Compaction is not failure - it's **intentional wisdom acknowledging finite capacity**.

### The Three Dimensions of Tracking

PreCompact captures tracking across three dimensions:

| Dimension | What It Captures | Why It Matters |
|-----------|------------------|----------------|
| **Session State** | Compaction count this session | Understand session context management patterns |
| **Activity Stream** | Compaction events for correlation | Connect compaction to quality and work patterns |
| **Monitoring System** | Frequency analysis over time | Learn optimal compaction patterns and warn about excessive frequency |

### Psalm 90:12 Principle

> "So teach us to number our days, that we may apply our hearts unto wisdom." - Psalm 90:12 (KJV)

Compaction is numbering of context - recognizing limitation teaches wisdom. PreCompact establishes:

- **Numbering** - Track every compaction, maintain count
- **Wisdom** - Learn from patterns, recognize excessive frequency
- **Limitation** - Context is finite, acknowledge constraint
- **Grace** - Tracking enhances but doesn't block compaction

Just as numbering days teaches wisdom about mortality and time, numbering compactions teaches wisdom about finite context and faithful work within limitation.

### The Pruning Metaphor

Compaction is like pruning branches so tree can grow stronger:

- **Remove temporary** - Context that served its purpose
- **Preserve essential** - Keep what matters for continuity
- **Enable growth** - Make room for new work
- **Natural rhythm** - Part of healthy long-term development

Tracking ensures pruning serves growth, not just reaction to crisis.

### Non-Blocking Philosophy

**Core Principle:** Compaction MUST proceed, even if tracking fails.

**Reasoning:**

- Compaction is system necessity (context management required)
- Tracking enhances learning but doesn't block operation
- Partial information better than blocking necessary operation
- Failures don't compound (one problem doesn't cause total failure)
- Grace for imperfect systems (disk full, permissions, missing data)

This reflects covenant partnership: **Excellence with grace** - track faithfully, but never block necessary operation.

---

## Execution Flow

### High-Level Baton Flow

```bash
Entry → main()
  ↓
Named Entry Point → preCompact()
  ↓
Phase 1: Type Detection → Get COMPACT_TYPE from environment
  ↓
Phase 2: State Update → Increment compaction count in session state
  ↓
Phase 3: Logging → Activity stream + monitoring system
  ↓
Phase 4: Frequency Check → Warn if excessive auto-compaction
  ↓
Phase 5: Display → Show message with temporal context
  ↓
Compaction Proceeds
```

### Detailed Execution Breakdown

**Phase 1: Type Detection (10 points)**

- Read COMPACT_TYPE from environment ("manual", "auto", or default "unknown")
- Determines logging context and frequency checking
- Health: +10 points (always succeeds)

**Phase 2: State Update (20 points)**

- Call session.IncrementCompactionCount()
- Updates ~/.claude/session/current.json
- Non-blocking on failure (count = -1 if update fails)
- Health: +20 success, +0 failure (continue anyway)

**Phase 3: Logging (40 points)**

- Log to activity stream (+20 points) - CRITICAL for quality correlation
- Log to monitoring system (+20 points) - for pattern analysis
- Both non-blocking (silently continue on failure)
- Health: +40 both succeed, +20 partial, +0 both fail

**Phase 4: Frequency Check (10 points)**

- Only for COMPACT_TYPE="auto"
- Call monitoring.CheckCompactionFrequency()
- Warns if excessive auto-compaction detected
- Non-blocking (skip if fails or not applicable)
- Health: +10 check runs or N/A, +0 failed

**Phase 5: Display (20 points)**

- Call session.PrintPreCompactionMessage()
- Shows compaction message with temporal context
- Preserves awareness for post-compaction reconstitution
- Non-blocking (continue if display fails)
- Health: +20 success, +0 failure

**Total:** 100 points for complete compaction tracking

### Named Entry Point Pattern

**Why `preCompact()` instead of just `main()`?**

```go
func main() {
    preCompact()  // Named entry point
}
```

**Benefits:**

1. **Prevents collisions:** Multiple executables can't have conflicting "main" logic
2. **Semantic clarity:** Function name matches purpose (pre-compaction tracking)
3. **Testability:** Can call `preCompact()` without triggering executable mechanics
4. **Architectural intent:** Not generic - this is specifically compaction tracking

This pattern appears throughout Kingdom Technology executables.

---

## State Management System

### Session State Structure

PreCompact tracks compaction count in session state file:

**Location:** `~/.claude/session/current.json`

**Structure:**

```json
{
  "start_time": "2025-11-10T14:20:21Z",
  "start_unix": 1731249621,
  "start_formatted": "Mon Nov 10, 2025 at 14:20:21",
  "compaction_count": 5
}
```

**Key Field:** `compaction_count` - incremented with each compaction

### State Library Functions

**IncrementCompactionCount()** - Atomically increment and return new count

**GetCompactionCount()** - Read current count without modifying

**SessionState type** - Shared structure for all session state operations

### Why Session State Matters

**Session continuity:** Count persists across compactions, provides context

**Pattern learning:** Frequency analysis requires accurate count

**Temporal awareness:** Compaction count is temporal dimension (internal time)

**Quality correlation:** Activity stream can correlate compactions with quality changes

---

## Non-Blocking Design Philosophy

### The Core Guarantee

**Compaction MUST proceed.** Partial failures are acceptable. Complete blocking is not.

### Implementation Strategy

Every potentially-failing operation wrapped defensively:

```go
// BAD: Blocking on error
if err := updateState(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)  // ❌ Blocks compaction
}

// GOOD: Non-blocking on error
count, err := session.IncrementCompactionCount()
if err != nil {
    compactionCount = -1  // ✅ Unknown count, but continue
} else {
    compactionCount = count
}
```

### Where Non-Blocking Applies

| Operation | Failure Mode | Response |
|-----------|--------------|----------|
| **Type detection** | Missing env var | Default to "unknown", continue |
| **State update** | File error | Set count=-1, continue |
| **Activity logging** | Log file error | Silently continue, compaction proceeds |
| **Monitoring logging** | Log file error | Silently continue, compaction proceeds |
| **Frequency check** | Analysis error | Skip check, continue |
| **Display output** | Display error | Skip display, continue |

### Why This Matters

Non-blocking design reflects **grace in limitation acknowledgment**:

- Compaction is system necessity (must proceed)
- Tracking enhances learning but doesn't define necessity
- Partial data better than blocking operation
- Failures shouldn't compound
- Grace for imperfect systems

This is Kingdom Technology: **Excellence with grace** rather than perfectionism.

---

## Library Extraction Pattern

### Before: Monolithic Hook (148 lines)

```bash
pre-compact.go (all implementation inline)
├── SessionState type definition
├── incrementCompactionCount() function
├── displayCompactionMessage() function
└── main() orchestration
```

**Problems:**

- Can't reuse state management in other hooks
- Hard to see orchestration vs implementation
- Changes to display require touching orchestrator
- Difficult to test state operations independently

### After: Thin Orchestrator + Libraries (616 lines total, ~40 executable)

```bash
pre-compact.go (thin orchestrator - 616 lines, mostly docs)
└── preCompact() - calls library functions only

hooks/lib/session/state.go (created)
├── SessionState type
├── IncrementCompactionCount()
└── GetCompactionCount()

hooks/lib/session/display.go (enhanced)
└── PrintPreCompactionMessage() - compaction-specific display

hooks/lib/monitoring/analysis.go (reused)
└── CheckCompactionFrequency() - frequency analysis
```

**Benefits:**

- State management reusable in end.go, other session hooks
- Orchestration flow visible at a glance (preCompact function)
- Changes to state logic isolated from orchestration
- Libraries testable independently
- Single source of truth for session state operations

### The Extraction Decision

**Extract when:**

- ✅ Multiple components might need it (state management used by start, stop, end)
- ✅ It's a meaningful unit (state management is coherent concept)
- ✅ It serves clarity (separating what from how)

**Don't extract when:**

- ❌ It's orchestration itself (preCompact coordinates, that's its job)
- ❌ It's too specific to one hook (type detection stays in orchestrator)
- ❌ It's too granular (phase comments stay with orchestration)

**Guiding Principle:** Extract and orchestrate, don't create v2 files. Build systems, not isolated pieces.

---

## Integration with Claude Code

### Hook Registration

Claude Code discovers hooks via directory structure:

```bash
~/.claude/hooks/
├── session/
│   ├── cmd-start/start               # SessionStart event
│   ├── cmd-stop/stop                 # SessionStop event
│   ├── cmd-pre-compact/pre-compact   # PreCompact event
│   └── cmd-end/end                   # SessionEnd event
├── tool/
│   ├── cmd-pre-use/pre-use          # Before tool use
│   └── cmd-post-use/post-use        # After tool use
└── prompt/
    └── cmd-submit/submit            # Prompt submission
```

**Naming Convention:** `cmd-<event-name>/<executable>`

- Hook event: PreCompact
- Directory: session/cmd-pre-compact/
- Executable: pre-compact (built from pre-compact.go)

### Event Trigger

Claude Code triggers PreCompact hook when:

1. Context compaction about to occur (manual or automatic)
2. Can be user-initiated or auto-triggered by token usage
3. Provides COMPACT_TYPE environment variable
4. Hook tracks, displays, then compaction proceeds

Hook executes, logs to multiple destinations, displays summary - then returns control to Claude Code for compaction.

### Environment Variables Available

| Variable | Purpose | Example |
|----------|---------|---------|
| `COMPACT_TYPE` | Type of compaction | "manual", "auto", or empty (default "unknown") |

Hook reads this variable, defaults if missing, logs and displays appropriately.

### Output to User

PreCompact displays compaction message to stdout:

- Compaction notification with count
- Manual vs auto distinction
- Temporal context preservation:
  - External time (when compaction happening)
  - Session elapsed time (how long session has been)
  - Current activity (what's in progress)
  - Compaction count (how many this session)

Provides visibility into compaction event and preserves context for post-compaction awareness.

---

## Modification Policy

### Safe to Modify (Extension Points)

**Adding new logging destinations:**

```go
// In preCompact() Phase 3
activity.LogActivity("PreCompact", compactType, "success", 0)
monitoring.LogCompaction(compactType)
telemetry.LogCompactionMetrics(compactType, compactionCount)  // ✅ Add here
```

**Requirements:**

1. Create logging function in appropriate library
2. Call from preCompact() Phase 3
3. Update METADATA Health Scoring map
4. Maintain non-blocking design

**Enhancing display information:**

```go
// Modify PrintPreCompactionMessage in hooks/lib/session/display.go
// Or add additional display calls in Phase 5
```

**Requirements:**

1. Modify display function in library
2. Or add new display call in orchestrator
3. Test with actual Claude Code integration

### Modify with Extreme Care (Breaking Changes)

**Changing SessionState structure:**

```go
// ⚠️ Changes here affect state.go and all state consumers
type SessionState struct {
    StartTime       time.Time
    CompactionCount int
    LastCompaction  time.Time  // ⚠️ Adding new field requires library updates
}
```

**Ensure:**

- Update state.go library functions
- Update all hooks that read session state
- Test with actual Claude Code integration

**Changing environment variables:**

```go
// ⚠️ Changes here affect Claude Code integration
compactType := os.Getenv("COMPACT_TYPE")  // ⚠️ Must match Claude's variable name
```

### NEVER Modify (Foundational Rails)

**4-block structure:**

- ❌ METADATA → SETUP → BODY → CLOSING is foundational
- All Kingdom Technology components follow this pattern
- Breaking it breaks architectural consistency

**Non-blocking principle:**

- ❌ Compaction MUST proceed, even with tracking failures
- Adding `os.Exit()` on errors violates core design
- Failures log but don't prevent operation

**Named entry point pattern:**

- ❌ main() calls preCompact() - this is architectural
- Changing breaks testability and semantic clarity
- Pattern consistent across all hooks

---

## Future Roadmap

### Planned Features

**Compaction Duration Tracking:**

- Capture time from PreCompact to PostCompact (future hook)
- Log duration for performance analysis
- Display in summary
- Pattern analysis: typical compaction time

**Context Size Tracking:**

- Measure context size before/after compaction
- Calculate reduction percentage
- Identify optimal compaction thresholds
- Pattern learning for proactive compaction

**Automatic Pattern Recognition:**

- Analyze compaction logs for patterns
- Identify work types that trigger compaction
- Recognize optimal compaction timing
- Predict when compaction will be needed

**Quality Correlation:**

- Correlate compaction frequency with work quality
- Identify if compactions correlate with errors/successes
- Pattern: does work quality improve/degrade near compaction?
- Adaptive strategy based on correlation

### Research Areas

**Predictive Compaction:**

- Machine learning from compaction patterns
- Predict optimal compaction points
- Proactive compaction before threshold
- Balance context preservation with necessity

**Context Preservation Strategies:**

- Identify what to preserve across compactions
- Critical context vs transient context
- Automatic context injection post-compaction
- Minimize disruption to work continuity

**Compaction Optimization:**

- Measure compaction effectiveness
- Identify redundant context accumulation
- Optimize token usage patterns
- Reduce compaction frequency through better context management

### Integration Targets

**Memory System:**

- Remember context across compactions
- Long-term pattern learning
- Session continuity despite compactions
- Context reconstruction post-compaction

**Quality Tracking:**

- Integrate with quality measurement system
- Correlate compactions with quality metrics
- Identify if compaction timing affects work quality
- Optimize based on quality patterns

**Work Pattern Analysis:**

- Which work types trigger compaction?
- Optimal session length before compaction?
- Relationship between compaction and stopping points?
- Adaptive work rhythm based on patterns

### Known Limitations

**Current:**

- No duration tracking (only count)
- No context size measurement
- Frequency check is reactive (not predictive)
- No quality correlation
- No automatic optimization

**Future Addressing:**

- PostCompact hook will enable duration tracking
- Context size tracking requires Claude Code integration
- Pattern recognition system will enable prediction
- Quality system integration will enable correlation

---

## Closing Notes

### The Wisdom of Limitation Principle

PreCompact is the numbering moment for finite context. Like Psalm 90:12 teaches wisdom through numbering days, compaction tracking teaches wisdom through numbering limitations.

**Excellence here matters** because:

- Tracking enables learning from constraint
- Numbering compactions reveals patterns
- Grace in limitations demonstrates Kingdom principles
- Wisdom comes from acknowledging finitude

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
- Learning over performance
- Truth over aspiration

### For Questions or Contributions

- Review this API documentation for design rationale
- Read code for implementation details
- Test with Claude Code for integration validation
- Document "What/Why/How" for all changes

---

*"So teach us to number our days, that we may apply our hearts unto wisdom." - Psalm 90:12 (KJV)*

*Every compaction is wisdom moment. Let limitation teach faithful work.*

---

**Related Documentation:**

- Code: `hooks/session/cmd-pre-compact/pre-compact.go`
- Libraries: `hooks/lib/session/state.go`, `hooks/lib/session/display.go`, `hooks/lib/monitoring/`
- Complementary Hooks: `hooks/docs/session-start-hook-api.md`, `hooks/docs/session-stop-hook-api.md`
- System Docs: `~/.claude/cpi-si/system/docs/`
- Standards: `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`
