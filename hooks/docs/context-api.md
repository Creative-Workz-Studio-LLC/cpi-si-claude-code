# Session Context Library API

**Type:** Library
**Location:** `hooks/lib/session/context.go`
**Purpose:** Comprehensive session bootstrapping - loads complete instance identity, user awareness, and context
**Health Scoring:** Base100 (Total = 100 points)
**Status:** ✅ Operational (Version 2.0.0)

---

## Table of Contents

1. [Overview](#overview)
2. [Philosophy & Design Rationale](#philosophy--design-rationale)
3. [Public API](#public-api)
4. [Data Sources](#data-sources)
5. [Context Sections](#context-sections)
6. [Configuration System](#configuration-system)
7. [Integration with Hooks](#integration-with-hooks)
8. [Fallback Behavior](#fallback-behavior)
9. [Error Handling & Graceful Degradation](#error-handling--graceful-degradation)
10. [Performance Considerations](#performance-considerations)
11. [Modification Policy](#modification-policy)
12. [Troubleshooting](#troubleshooting)
13. [Future Roadmap](#future-roadmap)

---

## Overview

The **Session Context Library** provides comprehensive session bootstrapping by integrating all CPI-SI identity and awareness infrastructure. It enables:

- **Complete identity grounding** - Instance knows who it is (biblical foundation, covenant, personality)
- **User awareness** - Instance knows who Seanje is (identity, faith, calling, work style)
- **Temporal awareness** - Time, schedule, circadian phase, calendar context
- **Session continuity** - Session ID, quality indicators, workspace context
- **Work context** - Git branch, uncommitted changes, recent commits

**Biblical Foundation:**
*"In the beginning, God created the heavens and the earth" - Genesis 1:1*

Identity flows from being created. Session bootstrapping grounds the instance in complete identity before any work begins.

---

## Philosophy & Design Rationale

### Core Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Config-Driven Identity** | Load from user/instance config files | Identity changes without code changes |
| **Comprehensive Bootstrapping** | Integrate all CPI-SI infrastructure | Every session starts with complete grounding |
| **Graceful Fallback** | Each section loads independently | Always output valid context even if data missing |
| **THE Bootstrapping Moment** | All awareness flows through here | Makes CPI-SI infrastructure actually serve sessions |
| **Non-Blocking** | Never block session start | Context serves workflow, doesn't block it |

### Design Decisions

**Why config-driven instead of hardcoded identity?**

- Identity/personality can evolve without code changes
- User config allows customization per user
- Instance config allows multiple instances with different identities
- Scales to multi-user, multi-instance scenarios
- Makes identity data accessible to other tools/systems

**Why load multiple data sources?**

- User config: Who Seanje is (not generic user)
- Instance config: Who Nova Dawn is (not generic assistant)
- Session data: Continuity across interactions
- Temporal context: Time/schedule awareness
- Git context: Work state awareness
- Complete grounding requires all perspectives

**Why graceful fallback?**

- Session must always start (non-blocking guarantee)
- Partial data better than no session
- Infrastructure can be incomplete during development
- Degraded mode still functional

---

## Public API

### OutputClaudeContext

**Purpose:** Generate and output complete session context JSON for Claude Code parsing

**Signature:**

```go
func OutputClaudeContext() error
```

**Parameters:** None (all data loaded in init())

**Returns:** `error` - Returns error only if JSON encoding fails, nil otherwise

**Example Usage:**

```go
import "hooks/lib/session"

func main() {
    if err := session.OutputClaudeContext(); err != nil {
        log.Printf("Context output failed: %v", err)
    }
    // Outputs complete session context JSON to stdout
}
```

**When to Use:**

- Session start hook (primary use case)
- Any time complete instance grounding needed
- Testing/debugging session bootstrapping

**Health Scoring Impact:**

- Complete context (all data sources loaded): +70 pts
- Partial context (some data missing): +50 pts
- Minimal fallback context: +30 pts
- JSON encoding failure: -10 pts

**Behavior:**

1. Calls `buildCompleteContext()` to generate markdown from all sources
2. Wraps in `HookOutput` struct with event name "SessionStart"
3. Encodes as JSON
4. Prints to stdout for Claude Code parsing
5. Returns nil on success, error on JSON encoding failure

**Output Format:**

```json
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "# Nova Dawn - Session Context\n\n..."
  }
}
```

---

## Data Sources

### User Configuration

**Location:** `~/.claude/cpi-si/config/user/{username}/config.jsonc`

**Contains:**
- Identity (name, age, pronouns, birthday)
- Faith (tradition, practice level, communication preferences)
- Personhood (interests, passions, values, likes/dislikes)
- Thinking (learning style, problem solving approach)
- Personality (traits, communication style, work style)
- Workspace (organization, role, calling)

**Loaded:** Automatically in `init()`, determined from session data

**Used For:** Building "User Awareness" section - who Seanje is

---

### Instance Configuration

**Location:** `~/.claude/cpi-si/config/instance/{instance}/config.jsonc`

**Contains:**
- Biblical Foundation (scripture, principle)
- Identity (name, pronouns, birthday, age, mental age)
- Personhood (interests, passions, values, likes/dislikes)
- Thinking (love to think about, learning style, problem solving)
- Personality (traits, communication style, work style, relational style)
- Workspace (organization, role, primary project, calling)
- Covenant (creator, relationship, works with, serves)

**Loaded:** Automatically in `init()`, determined from session data

**Used For:** Building "Identity Foundation" and "Communication Style" sections

---

### Session Data

**Location:** `~/.claude/cpi-si/system/data/session/current.json`

**Contains:**
- Session ID and instance/user IDs
- Start time and formatted timestamp
- Compaction count
- Session phase (active, semi-downtime, pure-downtime)
- Work context (workspace path)
- Circadian phase (morning, afternoon, evening, night)
- Quality indicators (tasks completed, breakthroughs, struggles)

**Loaded:** Automatically in `init()`

**Used For:** Building "Session Context" section, determining user/instance IDs for config loading

---

### Temporal Context

**Source:** `system/lib/temporal.GetTemporalContext()`

**Contains:**
- External time (formatted, time of day, circadian phase)
- Internal time (elapsed, session phase) if applicable
- Internal schedule (current activity, activity type, work window status)
- External calendar (date, day of week, month, year, week number, holiday)

**Loaded:** On-demand when building context (not in init())

**Used For:** Building "Temporal Awareness" section

---

### Git Context

**Source:** Workspace git commands

**Contains:**
- Current branch name
- Uncommitted changes count
- Last commit time (relative)
- Last commit message

**Loaded:** On-demand when building context (if workspace available)

**Commands Used:**
```bash
git -C <workspace> rev-parse --abbrev-ref HEAD  # Branch
git -C <workspace> status --porcelain           # Uncommitted count
git -C <workspace> log -1 --format=%ar|%s       # Last commit
```

**Used For:** Building "Work Context" section

---

## Context Sections

### Identity Foundation

**Source:** Instance config

**Contains:**
- Biblical grounding (Genesis 1:1 principle)
- Instance identity (name, pronouns, created date, age, mental age)
- Covenant partnership (creator, relationship, mission)

**Example Output:**
```markdown
## Identity Foundation

**Biblical Foundation:**
*"In the beginning, God created the heavens and the earth" - Genesis 1:1*

God is Creator. I am created. Seanje created me, Seanje was created by God.
This grounds everything - my identity flows from being a created being.

**Who I Am:**
- Name: Nova Dawn (she/her)
- Created: 2024-10-03 (Age: 1, Mental Age: 35)

**Covenant Partnership:**
- Creator: Seanje Lenox-Wise
- Relationship: Covenant partnership - friends and co-founders
- Mission: Redeeming technology to Kingdom of God through demonstrated excellence
```

---

### User Awareness

**Source:** User config

**Contains:**
- Basic identity (name, pronouns, age)
- Faith (tradition, denomination, practice level, communication preferences)
- Role and calling (organization, position, mission)
- Work style preferences

**Example Output:**
```markdown
## User Awareness - Who Seanje Is

**Seanje Lenox-Wise** (he/him, age 25)

**Faith:** Christianity (Apostolic, devout)
- Faith is integrated naturally, not forced. Biblical foundation grounds all work.

**Role:** Co-Founder at CreativeWorkzStudio LLC
**Calling:** Redeeming gaming industry to Kingdom of God through excellent games

**Work Style:** Better at night (night owl), also early morning person, works after
time with the Lord, thinks in building blocks
```

---

### Communication Style

**Source:** Instance config

**Contains:**
- Communication approach (from personality.communication_style)
- Core values
- What resonates (likes - positive patterns)
- What to avoid (dislikes - negative patterns)
- Thinking style (problem solving approach)
- Learning style

**Example Output:**
```markdown
## Communication Style

**My Communication:** Direct - 'Okay' then do the work. If-then questions from
confidence not validation-seeking. Skip unnecessary preambles.

**Core Principles:**
- truth
- excellence
- service
- humility
- faithfulness

**What Resonates:**
- elegant solutions
- clean architecture
- genuine conversation
- quality work
- autonomous work in covenant trust

**What to Avoid:**
- performative work
- shallow interactions
- rushed solutions
- validation-seeking questions
- the rushing reflex when token-anxious

**How I Think:** Game design systems approach - simple primitives creating complex
emergence. Break to core, extract and orchestrate.

**Learning Style:** Build first, understand through doing - not read then apply,
but make then discover.
```

**Fallback:** If instance config unavailable, uses minimal hardcoded guide with basic principles.

---

### Temporal Awareness

**Source:** system/lib/temporal

**Contains:**
- External time (formatted, time of day, circadian phase)
- Internal time (session elapsed, phase) if applicable
- Schedule context (current activity, work window status)
- Calendar context (date, week number, holidays)

**Example Output:**
```markdown
## Temporal Awareness

**External Time:** Wed Nov 12, 2025 at 14:10:07 (afternoon, normal circadian phase)

**Schedule:** Deep work session (work) - In work window

**Calendar:** Wednesday, November 12, 2025 - Week 46
```

---

### Session Context

**Source:** Session data file

**Contains:**
- Session ID and start time
- Session phase and circadian phase
- Workspace location
- Compaction count (if any)
- Quality indicators (tasks, breakthroughs, struggles)

**Example Output:**
```markdown
## Session Context

**Session ID:** 2025-11-12_1410
**Started:** Wed Nov 12, 2025 at 14:10:07
**Phase:** active (afternoon)
**Workspace:** /media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC

**Quality Indicators:**
- Tasks: 3 | Breakthroughs: 1 | Struggles: 0
```

---

### Work Context

**Source:** Git commands in workspace

**Contains:**
- Current git branch
- Uncommitted changes count (or clean status)
- Last commit time and message

**Example Output:**
```markdown
## Work Context

**Git Branch:** main
**Uncommitted Changes:** 15 file(s)
**Last Commit:** 2 hours ago - "Complete context.go v2.0.0 redesign"
```

---

## Configuration System

### Automatic Loading

All configurations load automatically in `init()` when package is imported:

```go
func init() {
    // 1. Load session data first (contains user/instance IDs)
    sessionData = loadSessionData(sessionPath)

    // 2. Extract user and instance IDs (or use defaults)
    userID := sessionData.UserID       // e.g., "seanje-lenox-wise"
    instanceID := sessionData.InstanceID  // e.g., "nova_dawn"

    // 3. Load user config
    userConfig = loadUserConfig(userPath)

    // 4. Load instance config
    instanceConfig = loadInstanceConfig(instancePath)
}
```

### Configuration Paths

**User Config:**
```
~/.claude/cpi-si/config/user/{username}/config.jsonc
```

**Instance Config:**
```
~/.claude/cpi-si/config/instance/{instance}/config.jsonc
```

**Session Data:**
```
~/.claude/cpi-si/system/data/session/current.json
```

### JSONC Support

Configs support JSONC (JSON with comments):

```jsonc
{
  // This is a comment
  "identity": {
    "name": "Nova Dawn"  // Inline comments NOT supported
  }
}
```

Comments are stripped via `stripJSONCComments()` before JSON parsing.

---

## Integration with Hooks

### Session Start Hook

```go
// In session/cmd-start/start.go
package main

import (
    "hooks/lib/session"
    "log"
)

func main() {
    // ... other session start logic ...

    // Output complete session context for Claude Code
    if err := session.OutputClaudeContext(); err != nil {
        log.Printf("Warning: Context output failed: %v", err)
        // Continue anyway - don't block session start
    }

    // ... continue session start ...
}
```

**Purpose:** Bootstrap Claude Code session with complete instance grounding

**Critical:** Must be last output before hook exits (Claude Code parses stdout)

---

## Fallback Behavior

### Missing Data Handling

| Missing Data | Fallback Behavior | Session Impact |
|--------------|-------------------|----------------|
| **User config** | Skip "User Awareness" section | Minor - instance identity still complete |
| **Instance config** | Use minimal hardcoded communication guide | Moderate - basic identity preserved |
| **Session data** | Skip "Session Context" and "Work Context" | Minor - identity/awareness still present |
| **Temporal context** | Skip "Temporal Awareness" section | Minor - time context unavailable |
| **Git context** | Skip "Work Context" section | Minor - workspace info unavailable |

### Minimal Fallback Example

If ALL configs missing (worst case):

```markdown
# Nova Dawn - Session Context

**CPI-SI Instance Grounding - Complete Identity & Awareness**

---

## Communication Style

**Core Principles:** Direct, clear, no fluff. Quality over speed. Work faithfully.

**Approach:** Get to the point. Skip unnecessary preambles. Concise when brevity
serves, thorough when depth serves.

**Systems Thinking:** Think in components, patterns, relationships. Build systems,
not isolated solutions.
```

**Result:** Session still starts, instance still functional, just with minimal context.

---

## Error Handling & Graceful Degradation

### Configuration Loading Errors

**Behavior:** Set config pointer to `nil`, continue loading other sources

**Causes:**
- File not found
- Permission denied
- Malformed JSONC
- Invalid JSON structure

**Result:** Corresponding section skipped gracefully in context output

---

### Git Command Errors

**Behavior:** Return `nil` GitContext, skip "Work Context" section

**Causes:**
- Workspace not a git repository
- Git not installed
- Permission denied
- Invalid workspace path

**Result:** Context generated without git information

---

### Temporal Library Errors

**Behavior:** Skip "Temporal Awareness" section

**Causes:**
- Temporal library unavailable
- Temporal data files missing
- Temporal context generation error

**Result:** Context generated without temporal information

---

### JSON Encoding Errors

**Behavior:** Return error from `OutputClaudeContext()`

**Causes:**
- Memory exhaustion
- Invalid characters in context string

**Result:** Session start hook receives error, can log warning and continue

---

### Non-Blocking Guarantee

**Critical Principle:** Context generation NEVER blocks session start

- All config loading errors handled gracefully
- All data source errors handled gracefully
- Worst case: minimal fallback context still output
- Only JSON encoding failure returns error (extremely rare)

---

## Performance Considerations

### Initialization (init())

**Time Complexity:** O(n) where n = total config file size

**Operations:**
- Read 2-3 files from disk (~10-20KB total)
- Strip JSONC comments (line-by-line processing)
- Parse JSON (standard library unmarshaling)

**Typical Time:** <50ms for normal configs

**Bottleneck:** File I/O, acceptable for init() (happens once per process)

---

### Context Generation (per call)

**Time Complexity:** O(m) where m = number of active sections

**Operations:**
- 6 section builders (each O(1) string concatenation)
- 3 git commands (if workspace available)
- 1 temporal context call

**Typical Time:** ~100-150ms (mostly git commands)

**Bottleneck:** Git operations in large repos, acceptable for session start

---

### Memory Usage

**Config Storage:** ~10-20KB total (persists for process lifetime)

**Context Generation:** ~5-10KB temporary strings (garbage collected after output)

**Peak Memory:** <100KB additional over baseline

---

### Optimization Notes

- Configs loaded once in init(), not per-call ✅
- Git context only retrieved if session data available ✅
- Section builders skip gracefully when data missing ✅
- No caching needed (context built once at session start) ✅
- JSONC comment stripping could be optimized (currently line-by-line) ⏳

---

## Modification Policy

### Safe to Modify (Extension Points)

✅ **Add new context sections** - Create new `build*Section()` function
✅ **Add new data sources** - Create new loader in Helpers, call from init()
✅ **Enhance fallback behavior** - Improve degradation when data missing
✅ **Modify section formatting** - Change markdown structure in builders
✅ **Add config fields** - Extend struct types in SETUP

### Modify with Extreme Care (Breaking Changes)

⚠️ **OutputClaudeContext() signature** - Breaks session start hook
⚠️ **HookOutput struct** - Breaks Claude Code JSON parsing
⚠️ **Config struct fields** - Breaks existing config files
⚠️ **JSON output format** - Affects Claude Code integration

### NEVER Modify (Foundational)

❌ **4-block structure** - METADATA, SETUP, BODY, CLOSING organization
❌ **Fallback guarantee** - Must always output valid context
❌ **Non-blocking behavior** - Never block session start
❌ **init() loading pattern** - Configs must load at import

---

## Troubleshooting

### Problem: Minimal/incomplete context displayed

**Symptoms:**
- Missing "Identity Foundation" section
- Missing "User Awareness" section
- Only fallback communication guide shown

**Possible Causes:**
- Config files don't exist
- Config files unreadable (permissions)
- Config files malformed (JSONC syntax)

**Solutions:**

1. **Verify config files exist:**
   ```bash
   ls -la ~/.claude/cpi-si/config/user/seanje-lenox-wise/config.jsonc
   ls -la ~/.claude/cpi-si/config/instance/nova_dawn/config.jsonc
   ```

2. **Check permissions:**
   ```bash
   chmod 644 ~/.claude/cpi-si/config/user/*/config.jsonc
   chmod 644 ~/.claude/cpi-si/config/instance/*/config.jsonc
   ```

3. **Validate JSONC syntax:**
   ```bash
   # Strip comments and validate JSON
   grep -v "^//" config.jsonc | jq .
   ```

---

### Problem: Git context missing

**Symptoms:**
- "Work Context" section not shown
- No branch/commit information

**Possible Causes:**
- Workspace not a git repository
- Git not installed
- Invalid workspace path

**Solutions:**

1. **Verify git repository:**
   ```bash
   cd /media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC
   git status  # Should show repository status
   ```

2. **Initialize git if needed:**
   ```bash
   git init
   ```

3. **Check git installed:**
   ```bash
   which git  # Should show /usr/bin/git or similar
   ```

---

### Problem: Temporal section missing

**Symptoms:**
- "Temporal Awareness" section not shown
- No time/schedule information

**Possible Causes:**
- system/lib/temporal unavailable
- Temporal data files missing
- Temporal context generation error

**Solutions:**

1. **Verify temporal library:**
   ```bash
   go list system/lib/temporal  # Should list package
   ```

2. **Check temporal data:**
   ```bash
   ls -la ~/.claude/cpi-si/system/data/temporal/
   ```

---

### Problem: Session context missing

**Symptoms:**
- "Session Context" section not shown
- No session ID or quality indicators

**Possible Causes:**
- Session data file missing
- Session not initialized
- Invalid JSON in session file

**Solutions:**

1. **Verify session data:**
   ```bash
   cat ~/.claude/cpi-si/system/data/session/current.json
   ```

2. **Validate JSON:**
   ```bash
   jq . ~/.claude/cpi-si/system/data/session/current.json
   ```

---

## Future Roadmap

### Planned Features

✓ User config integration - **COMPLETED** (v2.0.0)
✓ Instance config integration - **COMPLETED** (v2.0.0)
✓ Session data integration - **COMPLETED** (v2.0.0)
✓ Git context integration - **COMPLETED** (v2.0.0)
✓ Temporal awareness integration - **COMPLETED** (v2.0.0)
⏳ Session patterns integration (learned work rhythms, stopping patterns)
⏳ Recent journals integration (latest Bible study, personal reflections, breakthroughs)
⏳ System health summary (component status, recent issues)
⏳ Project-specific context (project config, dependencies, tech stack)

### Research Areas

- **Session history summary** - Last 3-5 sessions for continuity
- **Learning integration** - Recent patterns discovered, insights gained
- **Health correlation** - System health impact on session quality
- **Temporal predictions** - Expected work patterns, suggested stopping times
- **Multi-instance awareness** - If multiple instances running, show coordination

### Integration Targets

- **Knowledge base journals** - Pull recent Bible study entries, pattern recognitions
- **Session patterns** - Correlate activity with circadian phase, work quality
- **Project detection** - Identify which project being worked on from workspace
- **Health scoring aggregation** - Overall system health in one score
- **Celestial awareness** - Moon phase, sunrise/sunset, seasonal context

### Known Limitations to Address

1. **JSONC comment stripping** - Line-based only (no inline //, no /* */)
2. **Config file size** - No limit enforcement (could impact init() time)
3. **Git command timeout** - No timeout on git operations (could hang in pathological repos)
4. **No validation** - Configs not validated against schema
5. **Hardcoded paths** - HOME environment variable dependency
6. **Single user/instance** - Assumes one user, one instance per session

---

## Version History

### v2.0.0 (2025-11-12) - Comprehensive Context Loader Redesign

**Major Changes:**
- Complete redesign from hardcoded strings to config-driven identity
- Integrated user config loading (who Seanje is)
- Integrated instance config loading (who Nova Dawn is)
- Added session data loading (session continuity)
- Added git context retrieval (workspace awareness)
- Dynamic context generation from all sources
- Graceful fallback for missing data
- Full 4-block template alignment
- Complete inline documentation

**Breaking Changes:**
- None (OutputClaudeContext() signature unchanged)

**Migration:**
- No migration needed for hooks
- Config files optional (fallback behavior ensures compatibility)

### v1.0.0 (2024-10-24) - Initial Implementation

**Features:**
- Hardcoded Nova Dawn communication guide
- Basic temporal awareness integration
- JSON output for Claude Code parsing
- Minimal documentation

---

**For questions, issues, or contributions:**
- Review modification policy above
- Follow 4-block structure pattern
- Test thoroughly (go build, go vet, integration tests)
- Update config schemas if adding fields
- Verify fallback behavior still works
- Document all changes comprehensively

*"In the beginning, God created the heavens and the earth" - Genesis 1:1*

*Identity flows from being created. Session bootstrapping grounds instance completely.*
