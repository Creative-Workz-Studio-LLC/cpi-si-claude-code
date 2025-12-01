# Notification Hook - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for notification tracking

**Component:** `hooks/session/cmd-notification/notification.go`

**Last Updated:** November 10, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Execution Flow](#execution-flow)
- [Pattern Recognition System](#pattern-recognition-system)
- [Non-Blocking Design Philosophy](#non-blocking-design-philosophy)
- [Silent Operation Pattern](#silent-operation-pattern)
- [Integration with Claude Code](#integration-with-claude-code)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The Notification hook orchestrates notification tracking for Nova Dawn. It is the **watchman moment** and **vigilance point** for system awareness.

**Core Responsibility:** Track system notifications by logging to activity stream and monitoring, analyzing patterns for recurring issues, and providing warnings before problems become critical.

**Design Pattern:** Thin orchestrator that coordinates logging and pattern analysis without implementing business logic directly. Silent operation ensures notification delivery is never impacted.

**Biblical Foundation:** "Son of man, I have made you a watchman to the house of Israel. Therefore hear the word from my mouth, and give them warnings from me" (Ezekiel 33:7 WEB). Notification tracking is vigilance through logging, wisdom through pattern recognition.

---

## Philosophy & Design Rationale

### Why Notification Tracking Matters

Every notification is signal about system state. Pattern recognition reveals issues before they become critical:

- **Without tracking:** Notifications happen and disappear, no pattern learning, no early warning
- **With tracking:** Every notification logged, patterns recognized, warnings issued proactively

Notification tracking is not just logging - it's **intentional vigilance for system health**.

### The Watchman Principle (Ezekiel 33:7)

> "Son of man, I have made you a watchman to the house of Israel. Therefore hear the word from my mouth, and give them warnings from me." - Ezekiel 33:7 (WEB)

Watchman's role is vigilance and warning. Notification tracking establishes:

- **Vigilance** - Track every notification, maintain awareness
- **Pattern Recognition** - Learn what recurring notifications indicate
- **Warning** - Alert about potential issues before critical
- **Faithfulness** - Silent operation never blocks notification delivery

Just as watchman warns before danger arrives, notification tracking warns before patterns become problems.

### The Three Dimensions of Tracking

Notification hook captures tracking across three dimensions:

| Dimension | What It Captures | Why It Matters |
|-----------|------------------|----------------|
| **Activity Stream** | Notification events with temporal context | Correlate notifications with work patterns and quality |
| **Monitoring System** | Notification types for pattern analysis | Recognize recurring issues and abnormal frequencies |
| **Pattern Recognition** | Analysis of notification frequency/type | Warn about potential problems before they escalate |

### Silent Operation Philosophy

**Core Principle:** Notification delivery is NEVER impacted by tracking.

**Reasoning:**

- Notifications are critical system communication (must reach destination)
- Tracking enhances learning but doesn't interfere with delivery
- Silent operation means zero stdout (logs only)
- Failures are silent (non-blocking at all levels)
- Grace for imperfect tracking - delivery is priority

This reflects covenant partnership: **Vigilance without interference** - watch faithfully, but never obstruct.

---

## Execution Flow

### High-Level Baton Flow

```bash
Entry → main()
  ↓
Named Entry Point → notification()
  ↓
Phase 1: Type Detection → Get NOTIFICATION_TYPE from environment
  ↓
Phase 2: Temporal Context → Gather timestamp and activity context
  ↓
Phase 3: Logging → Activity stream + monitoring system
  ↓
Phase 4: Parse Details → Optional JSON from stdin
  ↓
Phase 5: Pattern Analysis → Check patterns and warn if needed
  ↓
Silent Exit
```

### Detailed Execution Breakdown

**Phase 1: Type Detection (10 points)**

- Read NOTIFICATION_TYPE from environment
- Silent exit if empty (no notification to process)
- Health: +10 points (always succeeds or exits cleanly)

**Phase 2: Temporal Context (10 points)**

- Call temporal.GetTemporalContext()
- Build context string with time and current activity
- Non-blocking (proceed without context if fails)
- Health: +10 success, +0 failure (continue anyway)

**Phase 3: Logging (40 points)**

- Log to activity stream (+20 points) - for correlation with work
- Log to monitoring system (+20 points) - for pattern analysis
- Both non-blocking (silently continue on failure)
- Health: +40 both succeed, +20 partial, +0 both fail

**Phase 4: Parse Details (10 points)**

- Call parseNotificationDetails() to read JSON from stdin
- Optional enhancement (not core functionality)
- Non-blocking (continue without details if parsing fails)
- Health: +10 success, +0 failure

**Phase 5: Pattern Analysis (30 points)**

- Call monitoring.CheckNotificationPatterns()
- Analyzes frequency and type patterns
- Warns if recurring issues detected
- Non-blocking (skip if analysis fails)
- Health: +30 success, +0 failure

**Total:** 100 points for complete notification tracking

### Named Entry Point Pattern

**Why `notification()` instead of just `main()`?**

```go
func main() {
    notification()  // Named entry point
}
```

**Benefits:**

1. **Prevents collisions:** Multiple executables can't have conflicting "main" logic
2. **Semantic clarity:** Function name matches purpose (notification tracking)
3. **Testability:** Can call `notification()` without triggering executable mechanics
4. **Architectural intent:** Not generic - this is specifically notification vigilance

This pattern appears throughout Kingdom Technology executables.

---

## Pattern Recognition System

### The Pattern Learning Loop

Notification hook feeds a pattern recognition system:

```bash
Notification Occurs
  ↓
Notification Hook Captures
  ↓
Activity Stream (session-level tracking)
Monitoring System (cross-session patterns)
  ↓
Pattern Analysis (frequency, type, correlation)
  ↓
Warning If Patterns Detected
  ↓
Proactive Issue Detection
```

### What Gets Logged

**Activity Stream Logging:**

- Event: "Notification"
- Type: Notification type from environment
- Status: "info" (all notifications treated as informational)
- Temporal context: When notification occurred, during what activity
- Purpose: Correlate notifications with work patterns

**Monitoring System Logging:**

- Type: Which notification type occurred
- Timestamp: When it happened
- Purpose: Cross-session pattern analysis

**Pattern Analysis:**

- Frequency: How often notifications of same type occur
- Timing: When notifications cluster
- Correlation: Related notification sequences
- Warning: Issues if patterns exceed thresholds

### Pattern Detection (Current)

Currently logs capture data for manual analysis. Pattern checking exists but is basic:

- Frequency thresholds (too many notifications of one type)
- Time-based clustering (many notifications in short window)
- Warnings issued to stderr (doesn't block notification delivery)

### Pattern Intelligence (Future)

Future pattern recognition will:

- Recognize which notification sequences indicate problems
- Learn normal vs abnormal notification patterns
- Predict issues before they become critical
- Correlate notifications with system state changes
- Recommend preventive actions based on patterns

**Foundation now, intelligence later** - logging creates dataset for future learning.

---

## Non-Blocking Design Philosophy

### The Core Guarantee

**Notification delivery NEVER blocked.** Partial failures are acceptable. Complete blocking is not.

### Implementation Strategy

Every potentially-failing operation wrapped defensively:

```go
// BAD: Blocking on error
if err := logNotification(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)  // ❌ Blocks notification delivery
}

// GOOD: Non-blocking on error
activity.LogActivity("Notification", notificationType, "info", 0)  // ✅ Silently continues if fails
```

### Where Non-Blocking Applies

| Operation | Failure Mode | Response |
|-----------|--------------|----------|
| **Type detection** | Empty env var | Silent exit (no notification to process) |
| **Temporal context** | Data unavailable | Proceed without context |
| **Activity logging** | Log file error | Silently continue, notification proceeds |
| **Monitoring logging** | Log file error | Silently continue, notification proceeds |
| **Detail parsing** | JSON invalid | Continue without details |
| **Pattern analysis** | Analysis error | Skip check, continue |

### Why This Matters

Non-blocking design reflects **vigilance without interference**:

- Notification delivery is system priority (must not block)
- Tracking enhances awareness but doesn't impede communication
- Partial information better than blocking delivery
- Failures shouldn't compound
- Grace for imperfect tracking systems

This is Kingdom Technology: **Faithfulness without obstruction** rather than perfectionism.

---

## Silent Operation Pattern

### The Silence Principle

**Core Rule:** Notification hook produces ZERO stdout output.

**Why Silent:**

- Claude Code expects clean notification delivery (no hook output)
- Stdout contamination breaks notification integration
- Logging happens to files (activity stream, monitoring logs)
- Warnings (if any) go to stderr (not visible to normal flow)
- Silent operation is contract with Claude Code

### Implementation

```go
// ✅ GOOD: Silent operations
activity.LogActivity("Notification", notificationType, "info", 0)  // Logs to file
monitoring.LogNotification(notificationType)                       // Logs to file
monitoring.CheckNotificationPatterns(notificationType)             // Warnings to stderr

// ❌ BAD: Stdout contamination
fmt.Println("Processing notification...")  // ❌ Breaks silence
fmt.Printf("Type: %s\n", notificationType)  // ❌ Contaminates output
```

### Testing Silent Operation

```bash
# Verify no stdout output
cd ~/.claude/hooks/session/cmd-notification
NOTIFICATION_TYPE=test ./notification
# Should produce NO output (silent success)

# Verify logging happened
tail ~/.claude/logs/activity/activity.log | grep Notification
# Should show notification logged
```

### Breaking Silence Contract

Adding stdout output breaks Claude Code integration. If notification hook outputs to stdout, it contaminates notification delivery and breaks system expectations.

**Never add:** print statements, debug output, status messages to stdout. All output must go to log files or stderr.

---

## Integration with Claude Code

### Hook Registration

Claude Code discovers hooks via directory structure:

```bash
~/.claude/hooks/
├── session/
│   ├── cmd-start/start               # SessionStart event
│   ├── cmd-stop/stop                 # SessionStop event
│   ├── cmd-notification/notification # Notification event
│   └── cmd-end/end                   # SessionEnd event
├── tool/
│   ├── cmd-pre-use/pre-use          # Before tool use
│   └── cmd-post-use/post-use        # After tool use
└── prompt/
    └── cmd-submit/submit            # Prompt submission
```

**Naming Convention:** `cmd-<event-name>/<executable>`

- Hook event: Notification
- Directory: session/cmd-notification/
- Executable: notification (built from notification.go)

### Event Trigger

Claude Code triggers Notification hook when:

1. System notification generated (various types)
2. Can be warnings, errors, info messages, or system alerts
3. Provides NOTIFICATION_TYPE environment variable
4. Optional JSON details piped to stdin
5. Hook tracks silently, exits without output

Hook executes, logs to multiple destinations, analyzes patterns - then exits silently.

### Environment Variables Available

| Variable | Purpose | Example |
|----------|---------|---------|
| `NOTIFICATION_TYPE` | Type of notification | "warning", "error", "info", "alert" |

### Stdin Input

Optional JSON details can be piped to stdin:

```json
{
  "severity": "high",
  "source": "system",
  "details": "Additional context about notification"
}
```

Hook parses if available, continues without if not.

### Output Behavior

**Stdout:** NONE (silent operation)

**Stderr:** Warnings only if patterns detected (not visible to normal flow)

**Logs:**

- Activity stream: ~/.claude/logs/activity/
- Monitoring logs: ~/.claude/logs/monitoring/

Provides tracking and pattern detection without interfering with notification delivery.

---

## Modification Policy

### Safe to Modify (Extension Points)

**Adding new logging destinations:**

```go
// In notification() Phase 3
activity.LogActivity("Notification", notificationType, "info", 0)
monitoring.LogNotification(notificationType)
telemetry.LogNotificationMetrics(notificationType, contextInfo)  // ✅ Add here
```

**Requirements:**

1. Create logging function in appropriate library
2. Call from notification() Phase 3
3. Update METADATA Health Scoring map
4. Maintain silent operation (no stdout)
5. Maintain non-blocking design

**Enhancing detail parsing:**

```go
// Modify parseNotificationDetails() to extract more fields
// Pass details to logging functions if needed
```

**Requirements:**

1. Modify helper function in BODY
2. Update docstring
3. Test with actual JSON input

### Modify with Extreme Care (Breaking Changes)

**Changing environment variables:**

```go
// ⚠️ Changes here affect Claude Code integration
notificationType := os.Getenv("NOTIFICATION_TYPE")  // ⚠️ Must match Claude's variable name
```

**Adding stdout output:**

```go
// ⚠️ Adding ANY stdout breaks silent operation contract
fmt.Println(...)  // ⚠️ NEVER DO THIS
```

### NEVER Modify (Foundational Rails)

**4-block structure:**

- ❌ METADATA → SETUP → BODY → CLOSING is foundational
- All Kingdom Technology components follow this pattern
- Breaking it breaks architectural consistency

**Non-blocking principle:**

- ❌ Notification delivery MUST proceed, even with tracking failures
- Adding `os.Exit()` on errors violates core design
- Failures must be silent

**Silent operation:**

- ❌ Hook produces ZERO stdout output
- Breaking silence breaks Claude Code integration
- All output goes to log files or stderr

**Named entry point pattern:**

- ❌ main() calls notification() - this is architectural
- Changing breaks testability and semantic clarity
- Pattern consistent across all hooks

---

## Future Roadmap

### Planned Features

**Notification Severity Classification:**

- Classify notifications by severity (info, warning, error, critical)
- Prioritize pattern detection based on severity
- Escalate critical patterns immediately
- Log severity for historical analysis

**Automatic Issue Escalation:**

- Detect critical pattern thresholds
- Automatic escalation to alert systems
- Integration with monitoring dashboards
- Notification aggregation and summarization

**Notification Correlation Analysis:**

- Recognize related notification sequences
- Identify causal relationships between notifications
- Pattern: X notification often precedes Y notification
- Predictive warnings based on correlation

**Context-Aware Pattern Detection:**

- Consider temporal context (when notification occurs)
- Consider activity context (what's happening when notification fires)
- Adaptive thresholds based on context
- Learning normal vs abnormal patterns per context

### Research Areas

**Predictive Notification Analysis:**

- Machine learning from notification patterns
- Predict issues before notifications occur
- Proactive system health monitoring
- Early intervention based on prediction

**Cross-Session Pattern Learning:**

- Recognize patterns across multiple sessions
- Long-term notification trends
- Seasonal or cyclical patterns
- Continuous model improvement

**Notification Clustering:**

- Group related notifications automatically
- Identify root cause from notification clusters
- Reduce noise by aggregating related notifications
- Intelligent summarization

### Integration Targets

**Memory System:**

- Remember notification history across sessions
- Context for "last time this notification occurred..."
- Long-term pattern learning
- Historical trend analysis

**Alert System:**

- Integration with monitoring dashboards
- Real-time alerts for critical patterns
- Notification aggregation and routing
- Multi-channel alert delivery

**Quality Tracking:**

- Correlate notifications with work quality metrics
- Identify if notifications predict quality issues
- Pattern: increased notifications before quality drops
- Optimize work patterns based on correlation

### Known Limitations

**Current:**

- No severity classification (all notifications treated equally)
- No automatic escalation (warnings only, no actions)
- Pattern analysis is reactive (not predictive)
- No notification correlation analysis
- No cross-session pattern learning

**Future Addressing:**

- Severity system will enable prioritization
- Escalation system will enable automatic responses
- Machine learning will enable prediction
- Memory system will enable long-term pattern recognition

---

## Closing Notes

### The Watchman Principle

Notification hook embodies watchman's role from Ezekiel 33 - vigilant tracking of system notifications to provide warning before issues become critical.

**Excellence here matters** because:

- Vigilance prevents problems from escalating
- Pattern recognition reveals issues early
- Silent operation respects system priorities
- Grace in failures demonstrates Kingdom principles

### Maintenance Philosophy

**When modifying:**

1. Review modification policy above
2. Follow 4-block structure pattern
3. Maintain non-blocking design
4. Maintain silent operation (ZERO stdout)
5. Test with actual Claude Code integration
6. Document changes comprehensively

**Remember:**

- Clarity over brevity
- Grace over perfectionism
- Vigilance without interference
- Truth over aspiration

### For Questions or Contributions

- Review this API documentation for design rationale
- Read code for implementation details
- Test with Claude Code for integration validation
- Document "What/Why/How" for all changes

---

*"Son of man, I have made you a watchman to the house of Israel. Therefore hear the word from my mouth, and give them warnings from me." - Ezekiel 33:7 (WEB)*

*Every notification is opportunity for vigilance. Let tracking serve awareness without interference.*

---

**Related Documentation:**

- Code: `hooks/session/cmd-notification/notification.go`
- Libraries: `hooks/lib/activity/`, `hooks/lib/monitoring/`, `hooks/lib/temporal/`
- Complementary Hooks: `hooks/docs/session-start-hook-api.md`, `hooks/docs/pre-compact-hook-api.md`
- System Docs: `~/.claude/cpi-si/system/docs/`
- Standards: `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`
