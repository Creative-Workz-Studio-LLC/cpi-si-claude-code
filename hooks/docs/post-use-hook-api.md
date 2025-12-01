# PostToolUse Hook - API Documentation

**Purpose:** Deep philosophy, design rationale, and integration guidance for post-tool validation and feedback

**Component:** `hooks/tool/cmd-post-use/post-use.go`

**Last Updated:** November 10, 2025

---

## Table of Contents

- [Overview](#overview)
- [Philosophy & Design Rationale](#philosophy--design-rationale)
- [Execution Flow](#execution-flow)
- [Validation and Feedback System](#validation-and-feedback-system)
- [Non-Blocking Design Philosophy](#non-blocking-design-philosophy)
- [Handler Pattern](#handler-pattern)
- [Integration with Claude Code](#integration-with-claude-code)
- [Modification Policy](#modification-policy)
- [Future Roadmap](#future-roadmap)

---

## Overview

The PostToolUse hook orchestrates post-tool validation and contextual feedback for Nova Dawn. It is the **quality verification moment** and **feedback point** after every tool execution.

**Core Responsibility:** Format code files automatically after Write/Edit operations, validate syntax with immediate feedback, log tool usage with temporal context, provide contextual guidance (commits, builds, dependencies), and detect command failures for reporting.

**Design Pattern:** Thin orchestrator that routes tool usage to appropriate handlers (file editing, bash commands, read/search operations) without implementing validation or feedback logic directly.

**Biblical Foundation:** "Test all things, and hold firmly that which is good" (1 Thessalonians 5:21 WEB). Post-tool validation is testing the work - verifying quality, providing feedback, ensuring excellence before proceeding.

---

## Philosophy & Design Rationale

### Why Post-Tool Validation Matters

Every tool use is opportunity for quality verification and learning:

- **Without validation:** Tools execute without feedback, errors accumulate, no pattern learning
- **With validation:** Immediate feedback on quality, automatic formatting, contextual guidance, pattern recognition

Post-tool validation is not just checking - it's **intentional quality assurance through testing and feedback**.

### The Testing Principle (1 Thessalonians 5:21)

> "Test all things, and hold firmly that which is good." - 1 Thessalonians 5:21 (WEB)

Testing is verification before accepting. PostToolUse establishes:

- **Testing** - Verify syntax, format code, validate quality
- **Holding Fast** - Only accept good work, report issues immediately
- **Feedback** - Provide contextual guidance based on tool type
- **Learning** - Log patterns for recognition and improvement
- **Grace** - Non-blocking design (validation enhances but doesn't obstruct)

Just as testing ensures quality before accepting, post-tool validation ensures excellence through immediate verification.

### The Three Dimensions of Post-Tool Processing

PostToolUse coordinates across three dimensions:

| Dimension | What It Captures | Why It Matters |
|-----------|------------------|----------------|
| **Validation** | Syntax checking, code formatting | Ensure code quality immediately after creation |
| **Feedback** | Contextual guidance (commits, builds, installs) | Provide relevant information at right moment |
| **Logging** | Tool usage with temporal context | Enable pattern recognition and learning |

### The Immediate Feedback Pattern

Post-tool validation provides immediate feedback at optimal moment:

- **Right After Tool Use** - When context is fresh, corrections are cheapest
- **Contextual** - Feedback matches tool type (git commits, builds, installs)
- **Non-Blocking** - Tool completion proceeds regardless of validation
- **Constructive** - Focus on improvement, not just error reporting

Immediate feedback loop enables rapid quality improvement without workflow interruption.

### Non-Blocking Philosophy

**Core Principle:** Tool completion MUST proceed, even if validation fails.

**Reasoning:**

- Tool execution is primary operation (must complete successfully)
- Validation enhances quality but doesn't block progress
- Feedback serves awareness but doesn't prevent continuation
- Failures don't compound (one problem doesn't cause workflow stop)
- Grace for imperfect systems (formatters may fail, validators may error)

This reflects covenant partnership: **Quality verification with grace** - test faithfully, but never block necessary work.

---

## Execution Flow

### High-Level Baton Flow

```bash
Entry → main()
  ↓
Named Entry Point → postToolUse()
  ↓
Phase 1: Argument Validation → Check os.Args length
  ↓
Phase 2: Parse Tool Info → Extract tool name and args
  ↓
Phase 3: Gather Temporal Metadata → For pattern recognition
  ↓
Phase 4: Route to Handler:
  ├─> Write/Edit → handleFileEdit()
  │     ├─> Log tool use
  │     ├─> Format file
  │     └─> Validate syntax
  ├─> Bash → handleBashCommand()
  │     ├─> Log command with exit code/duration
  │     ├─> Contextual feedback (commits, builds, installs)
  │     └─> Check for failures
  └─> Read/Grep/Glob → Log activity directly
  ↓
Phase 5: Log Temporal Context → Enable pattern learning
  ↓
Exit
```

### Detailed Execution Breakdown

**Phase 1: Argument Validation (5 points)**

- Check os.Args length (need at least 3 arguments)
- Silent exit if insufficient args (tool name and tool args required)
- Health: +5 points (always succeeds or exits cleanly)

**Phase 2: Parse Tool Info (10 points)**

- Extract tool name from os.Args[1]
- Extract tool args from os.Args[2]
- Determines routing to appropriate handler
- Health: +10 success

**Phase 3: Gather Temporal Metadata (15 points)**

- Call getTemporalMetadata() to format temporal context
- Creates pipe-delimited string: timeOfDay|sessionPhase|activityType
- Non-blocking (empty string if temporal unavailable)
- Health: +15 success, +0 failure (continue anyway)

**Phase 4: Handler Execution (60 points)**

Varies by tool type:

- **Write/Edit (60 points):**
  - Log tool use (+20 points)
  - Format file (+20 points)
  - Validate syntax (+20 points)

- **Bash (60 points):**
  - Log command with exit code/duration (+20 points)
  - Provide contextual feedback (+30 points)
  - Check for failures (+10 points)

- **Read/Grep/Glob (60 points):**
  - Log activity (+60 points)

**Phase 5: Temporal Logging (10 points)**

- Log temporal metadata to activity stream
- Enables pattern recognition across tool usage
- Non-blocking (skip if metadata empty)
- Health: +10 success, +0 failure

**Total:** 100 points for complete post-tool validation

### Named Entry Point Pattern

**Why `postToolUse()` instead of just `main()`?**

```go
func main() {
    postToolUse()  // Named entry point
}
```

**Benefits:**

1. **Prevents collisions:** Multiple executables can't have conflicting "main" logic
2. **Semantic clarity:** Function name matches purpose (post-tool validation)
3. **Testability:** Can call `postToolUse()` without triggering executable mechanics
4. **Architectural intent:** Not generic - this is specifically post-tool processing

This pattern appears throughout Kingdom Technology executables.

---

## Validation and Feedback System

### File Validation Process

PostToolUse validates files after Write/Edit operations:

**Step 1: Log Tool Use**

- Record which file was edited
- Privacy-preserving logging (file paths only)
- Activity stream correlation

**Step 2: Format File**

- Call validation.FormatFile(filePath, ext)
- Automatic formatting based on file type
- Reports formatting results to user

**Step 3: Validate Syntax**

- Call validation.ValidateFile(filePath, ext)
- Checks syntax errors immediately
- Reports validation results to user

**Why Format Then Validate:**

- Formatting may fix syntax issues automatically
- Validation shows remaining issues after formatting
- User sees clean code with clear error messages

### Contextual Feedback System

PostToolUse provides contextual feedback based on command patterns:

**Git Commit Detected:**

- Pattern: `strings.Contains(cmdStr, "git commit")`
- Feedback: OnCommitCreated() - commit guidance
- Purpose: Remind about documentation, related changes

**Build Command Detected:**

- Pattern: `feedback.IsBuildCommand(cmdStr)`
- Feedback: OnBuildComplete() - build success guidance
- Purpose: Suggest testing, deployment considerations

**Dependency Install Detected:**

- Pattern: Contains "npm install", "cargo install", "go install"
- Feedback: OnDependenciesInstalled() - dependency guidance
- Purpose: Remind about version locking, security

**Command Failure Detected:**

- Check: BASH_EXIT_CODE environment variable
- Feedback: OnCommandFailure() - failure guidance
- Purpose: Suggest debugging steps, common fixes

### Temporal Context Logging

Every tool use logs temporal metadata for pattern recognition:

**What Gets Logged:**

- Time of day (morning, afternoon, evening, night)
- Session phase (short, medium, long)
- Activity type (focused-work, exploration, maintenance)

**Why Log Temporal Context:**

- Recognize when different tools are used
- Identify optimal work patterns
- Correlate tool usage with session quality
- Learn natural workflow rhythms

**Pattern Learning Enables:**

- "You typically format code in evening sessions"
- "Build commands often cluster after focused work"
- "Grep usage increases during exploration phases"

---

## Non-Blocking Design Philosophy

### The Core Guarantee

**Tool completion MUST proceed.** Validation failures are acceptable. Blocking tool execution is not.

### Implementation Strategy

Every potentially-failing operation wrapped defensively:

```go
// BAD: Blocking on error
result := validation.FormatFile(filePath, ext)
if result.HasErrors {
    fmt.Fprintf(os.Stderr, "Error: formatting failed\n")
    os.Exit(1)  // ❌ Blocks tool completion
}

// GOOD: Non-blocking on error
result := validation.FormatFile(filePath, ext)
result.Report()  // ✅ Reports but continues
```

### Where Non-Blocking Applies

| Operation | Failure Mode | Response |
|-----------|--------------|----------|
| **Argument validation** | Insufficient args | Silent exit (no tool to process) |
| **Temporal metadata** | Context unavailable | Empty string, continue without |
| **File formatting** | Formatter error | Report failure, continue to validation |
| **Syntax validation** | Validator error | Report failure, exit cleanly |
| **Activity logging** | Log file error | Silently continue, tool completed anyway |
| **Feedback display** | Display error | Skip feedback, tool completed |

### Why This Matters

Non-blocking design reflects **quality verification with grace**:

- Tool execution is primary operation (must succeed)
- Validation serves improvement but doesn't block progress
- Feedback enhances awareness but doesn't obstruct workflow
- Failures shouldn't prevent tool completion
- Grace for imperfect validation systems

This is Kingdom Technology: **Testing with grace** rather than perfectionism blocking progress.

---

## Handler Pattern

### Before: Monolithic Main

All logic implemented directly in main():

```go
func main() {
    // Parse args
    // Route tool
    // Format files
    // Validate syntax
    // Log activity
    // Provide feedback
    // All inline in main()
}
```

**Problems:**

- Hard to see orchestration vs implementation
- Difficult to test handlers independently
- Changes require touching main() directly
- No separation of concerns

### After: Handler Functions + Orchestrator

```bash
post-use.go (thin orchestrator)
├── postToolUse() - main orchestration
├── getTemporalMetadata() - temporal formatting helper
├── handleFileEdit() - file editing validation orchestrator
└── handleBashCommand() - bash command feedback orchestrator

hooks/lib/validation/ (validation library)
├── FormatFile() - code formatting
└── ValidateFile() - syntax checking

hooks/lib/feedback/ (feedback library)
├── OnCommitCreated() - commit guidance
├── OnBuildComplete() - build guidance
├── OnDependenciesInstalled() - dependency guidance
└── OnCommandFailure() - failure guidance
```

**Benefits:**

- Clear separation: orchestration vs implementation
- Handlers testable independently
- Changes to validation/feedback isolated in libraries
- Main orchestrator shows flow at a glance

### The Handler Decision

**Create handler when:**

- ✅ Multiple steps needed (file editing: log + format + validate)
- ✅ Coordination required (bash commands: log + feedback + failure check)
- ✅ Tool type has distinct workflow

**Inline when:**

- ✅ Single operation (Read: just log)
- ✅ Simple routing (Grep/Glob: log activity directly)
- ✅ No coordination needed

**Guiding Principle:** Handlers coordinate libraries for complex tool types. Simple tool types inline directly in orchestrator.

---

## Integration with Claude Code

### Hook Registration

Claude Code discovers hooks via directory structure:

```bash
~/.claude/hooks/
├── session/
│   ├── cmd-start/start               # SessionStart event
│   ├── cmd-stop/stop                 # SessionStop event
│   └── cmd-end/end                   # SessionEnd event
├── tool/
│   ├── cmd-pre-use/pre-use          # Before tool use
│   └── cmd-post-use/post-use        # After tool use
└── prompt/
    └── cmd-submit/submit            # Prompt submission
```

**Naming Convention:** `cmd-<event-name>/<executable>`

- Hook event: PostToolUse
- Directory: tool/cmd-post-use/
- Executable: post-use (built from post-use.go)

### Event Trigger

Claude Code triggers PostToolUse hook when:

1. Any tool completes execution (Write, Edit, Bash, Read, Grep, Glob)
2. Passes tool name and args as command-line arguments
3. Provides environment variables (FILE_PATH, BASH_EXIT_CODE, etc.)
4. Hook validates, formats, provides feedback, logs activity
5. Hook exits, Claude Code continues

Hook executes after tool, provides immediate feedback - then session continues.

### Command-Line Arguments

| Argument | Purpose | Example |
|----------|---------|---------|
| `os.Args[1]` | Tool name | "Write", "Edit", "Bash", "Read" |
| `os.Args[2]` | Tool args | "file.go", "git commit -m 'msg'" |

### Environment Variables Available

| Variable | Purpose | Example |
|----------|---------|---------|
| `FILE_PATH` | File path for Write/Edit operations | "/path/to/file.go" |
| `BASH_EXIT_CODE` | Exit code from bash command | "0", "1", "127" |
| `BASH_DURATION_MS` | Command duration in milliseconds | "1234" |
| `GREP_PATTERN` | Search pattern for Grep operations | "function.*main" |
| `GLOB_PATTERN` | File pattern for Glob operations | "**/*.go" |

Hook reads these variables, defaults if missing, processes appropriately.

### Output to User

PostToolUse displays to stdout:

- Formatting results (files formatted successfully)
- Validation results (syntax errors found)
- Contextual feedback (commit reminders, build suggestions)
- Command failure guidance (debugging steps)

Provides immediate quality feedback at optimal moment.

---

## Modification Policy

### Safe to Modify (Extension Points)

**Adding new tool type handler:**

```go
// In postToolUse() switch statement
case strings.HasPrefix(toolName, "NewTool"):
    handleNewTool(toolArgs)  // ✅ Add handler here
```

**Requirements:**

1. Create handler function or inline logging
2. Update organizational chart in BODY
3. Update health scoring map in METADATA
4. Test with actual tool execution

**Adding new bash feedback pattern:**

```go
// In handleBashCommand() switch statement
case strings.Contains(cmdStr, "docker build"):
    feedback.OnDockerBuild()  // ✅ Add pattern here
```

**Requirements:**

1. Create feedback function in hooks/lib/feedback
2. Add pattern detection in handleBashCommand()
3. Test with actual commands

### Modify with Extreme Care (Breaking Changes)

**Changing argument structure:**

```go
// ⚠️ Changes here affect Claude Code integration
toolName := os.Args[1]  // ⚠️ Must match Claude's calling convention
toolArgs := os.Args[2]  // ⚠️ Must be second argument
```

**Changing environment variables:**

```go
// ⚠️ Changes here affect tool routing
filePath := os.Getenv("FILE_PATH")  // ⚠️ Must match Claude's variable name
```

### NEVER Modify (Foundational Rails)

**4-block structure:**

- ❌ METADATA → SETUP → BODY → CLOSING is foundational
- All Kingdom Technology components follow this pattern
- Breaking it breaks architectural consistency

**Non-blocking principle:**

- ❌ Tool completion MUST proceed, even with validation failures
- Adding `os.Exit()` on errors violates core design
- Failures report but don't prevent tool completion

**Named entry point pattern:**

- ❌ main() calls postToolUse() - this is architectural
- Changing breaks testability and semantic clarity
- Pattern consistent across all hooks

---

## Future Roadmap

### Planned Features

**AI-Powered Code Review:**

- Machine learning for code quality prediction
- Automated refactoring suggestions based on patterns
- Context-aware improvement recommendations
- Learning from accepted/rejected suggestions

**Performance Profiling Integration:**

- Automatic performance analysis for critical functions
- Benchmark suggestions based on code patterns
- Memory usage prediction
- Optimization opportunity detection

**Security Scanning:**

- Common vulnerability detection (SQL injection, XSS, etc.)
- Dependency security analysis
- Credential detection and warnings
- OWASP Top 10 checking

**Advanced Pattern Recognition:**

- Cross-tool usage pattern detection
- Workflow optimization suggestions
- Temporal pattern learning (best times for different work)
- Quality correlation with tool usage patterns

### Research Areas

**Predictive Validation:**

- Predict issues before they occur based on patterns
- Suggest preventive refactoring
- Recognize code smells early
- Proactive quality guidance

**Adaptive Feedback:**

- Learn which feedback is helpful vs ignored
- Personalize guidance based on work patterns
- Adjust verbosity based on user preferences
- Context-sensitive feedback timing

**Cross-Session Learning:**

- Remember validation patterns across sessions
- Long-term quality trend analysis
- Team-wide pattern aggregation
- Continuous improvement metrics

### Integration Targets

**Code Quality Dashboard:**

- Real-time quality metrics visualization
- Validation success rate tracking
- Formatting coverage analysis
- Feedback effectiveness measurement

**Team Collaboration:**

- Shared validation patterns
- Team-wide feedback guidelines
- Quality standard enforcement
- Collaborative improvement tracking

**Continuous Integration:**

- Pre-commit validation integration
- CI/CD pipeline feedback
- Automated quality gating
- Deployment readiness checking

### Known Limitations

**Current:**

- No AI-powered code review (manual validation only)
- No performance profiling (timing tracking only)
- No security scanning (syntax validation only)
- No cross-session learning (session-local patterns)
- No adaptive feedback (static guidance)

**Future Addressing:**

- Machine learning integration will enable AI review
- Profiling libraries will enable performance analysis
- Security scanners will enable vulnerability detection
- Memory systems will enable cross-session learning
- Feedback system will enable adaptive guidance

---

## Closing Notes

### The Testing Principle

PostToolUse embodies testing from 1 Thessalonians 5:21 - verifying quality through immediate validation and constructive feedback before proceeding.

**Excellence here matters** because:

- Immediate feedback enables rapid quality improvement
- Testing ensures excellence before problems compound
- Contextual guidance serves awareness at optimal moment
- Non-blocking design demonstrates grace in quality verification

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
- Testing over trusting
- Quality with non-obstruction

### For Questions or Contributions

- Review this API documentation for design rationale
- Read code for implementation details
- Test with Claude Code for integration validation
- Document "What/Why/How" for all changes

---

*"Test all things, and hold firmly that which is good." - 1 Thessalonians 5:21 (WEB)*

*Every tool use is opportunity for excellence through testing. Let validation serve quality without obstruction.*

---

**Related Documentation:**

- Code: `hooks/tool/cmd-post-use/post-use.go`
- Libraries: `hooks/lib/validation/`, `hooks/lib/feedback/`, `hooks/lib/activity/`, `hooks/lib/temporal/`
- Complementary Hooks: `hooks/docs/session-start-hook-api.md`, `hooks/docs/session-end-hook-api.md`
- System Docs: `~/.claude/cpi-si/system/docs/`
- Standards: `~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md`
