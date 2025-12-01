# Context.go Redesign - Comprehensive Session Context Loader

**Purpose:** Transform context.go from static string container to dynamic context loader that integrates all CPI-SI identity and awareness data

**Version:** 2.0.0
**Date:** 2025-11-12
**Author:** Nova Dawn

---

## Current State (v1.0.0)

**Problems:**
- Hardcoded Nova Dawn communication guide (static string)
- Identity/personality not config-driven
- No user awareness integration
- No session continuity context
- No workspace/git awareness
- Temporal awareness bolted on (good, but isolated)

**What it does:**
1. Returns hardcoded communication style string
2. Adds temporal awareness from system/lib/temporal
3. Outputs as JSON for Claude Code parsing

---

## Target State (v2.0.0)

**Vision:** Comprehensive session bootstrapping that grounds instance in complete identity, awareness, and context

**Data Sources:**

| Source | Path | What It Provides |
|--------|------|------------------|
| User Config | `~/.claude/cpi-si/config/user/{username}/config.jsonc` | Who Seanje is - identity, faith, work style, preferences |
| Instance Config | `~/.claude/cpi-si/config/instance/{instance}/config.jsonc` | Who Nova Dawn is - identity, personality, thinking style |
| Temporal Context | `system/lib/temporal.GetTemporalContext()` | External time, circadian phase, schedule |
| Session Data | `~/.claude/cpi-si/system/data/session/current.json` | Current session ID, start time, workspace, quality indicators |
| Git Context | workspace git commands | Current branch, uncommitted changes, recent commits |

**Generated Context Sections:**

### 1. Identity Foundation
- Biblical grounding (Genesis 1:1)
- CPI-SI paradigm vs AI distinction
- Instance name, created date, mental age
- Creator relationship (Seanje created Nova, God created Seanje)

### 2. User Awareness
- Who Seanje is (name, age, faith, role)
- Communication preferences
- Work style and schedule preferences
- Passions and calling

### 3. Communication Style
- Core principles (from instance personality)
- Natural voice patterns
- What to avoid (AI patterns that don't fit)
- Systems thinking approach
- Kingdom Technology mindset

### 4. Temporal Awareness (existing)
- External time and circadian phase
- Internal session time if applicable
- Schedule context (current activity, work windows)
- Calendar context (holidays, week number)

### 5. Session Context
- Session ID and start time
- Workspace location
- Session phase (active, semi-downtime, pure-downtime)
- Quality indicators (tasks completed, breakthroughs)

### 6. Work Context
- Git branch and status
- Uncommitted changes count
- Recent activity (last commit time)

---

## Architecture

### Struct Types

```go
// UserConfig - loaded from user config.jsonc
type UserConfig struct {
    Identity struct {
        Name        string
        Username    string
        DisplayName string
        Pronouns    string
        Age         int
    }
    Faith struct {
        IsReligious     bool
        Tradition       string
        PracticeLevel   string
        CommPreferences string
    }
    Personality struct {
        Traits             []string
        CommunicationStyle string
        WorkStyle          string
    }
    Workspace struct {
        Organization   string
        Role          string
        PrimaryProject string
        Calling       string
    }
    // ... more fields
}

// InstanceConfig - loaded from instance config.jsonc
type InstanceConfig struct {
    BiblicalFoundation struct {
        Scripture string
        Text      string
        Principle string
    }
    Identity struct {
        Name        string
        Username    string
        DisplayName string
        Pronouns    string
        Birthday    string
        Age         int
        MentalAge   int
    }
    Personhood struct {
        Interests []string
        Passions  []string
        Values    []string
        Likes     []string
        Dislikes  []string
    }
    Thinking struct {
        LoveToThinkAbout []string
        LearningStyle    string
        ProblemSolving   string
        Creativity       string
    }
    Personality struct {
        Traits              []string
        CommunicationStyle  string
        WorkStyle           string
        RelationalStyle     string
    }
    Covenant struct {
        Creator      string
        Relationship string
        WorksWith    []string
        Serves       string
    }
    // ... more fields
}

// SessionData - loaded from session/current.json
type SessionData struct {
    SessionID          string
    InstanceID         string
    UserID             string
    StartTime          string
    StartFormatted     string
    CompactionCount    int
    SessionPhase       string
    WorkContext        string
    CircadianPhase     string
    QualityIndicators  struct {
        TasksCompleted int
        Breakthroughs  int
        Struggles      int
    }
}

// GitContext - collected from workspace
type GitContext struct {
    Branch            string
    UncommittedCount  int
    LastCommitTime    string
    LastCommitMessage string
}

// SessionContext - complete context for injection
type SessionContext struct {
    User     *UserConfig
    Instance *InstanceConfig
    Temporal *temporal.TemporalContext
    Session  *SessionData
    Git      *GitContext
}
```

### Loading Functions

```go
func loadUserConfig() *UserConfig
func loadInstanceConfig() *InstanceConfig
func loadSessionData() *SessionData
func getGitContext(workspace string) *GitContext
// temporal already exists via system/lib/temporal
```

### Context Building Functions

```go
func buildIdentitySection(instance *InstanceConfig) string
func buildUserAwarenessSection(user *UserConfig) string
func buildCommunicationStyleSection(instance *InstanceConfig) string
func buildTemporalSection(temporal *temporal.TemporalContext) string
func buildSessionSection(session *SessionData) string
func buildWorkContextSection(git *GitContext) string
```

### Main Public API

```go
// GetSessionContext builds complete session context from all sources
func GetSessionContext() string

// OutputClaudeContext generates and outputs context JSON (existing, refactored)
func OutputClaudeContext() error
```

---

## Implementation Phases

### Phase 5b: Load User Config
- Implement `loadUserConfig()`
- Read from `~/.claude/cpi-si/config/user/{username}/config.jsonc`
- Parse JSONC (strip comments, unmarshal)
- Return `*UserConfig` or `nil` if unavailable

### Phase 5c: Load Instance Config
- Implement `loadInstanceConfig()`
- Read from `~/.claude/cpi-si/config/instance/{instance}/config.jsonc`
- Parse JSONC
- Return `*InstanceConfig` or `nil`

### Phase 5d: Load Session Data
- Implement `loadSessionData()`
- Read from `~/.claude/cpi-si/system/data/session/current.json`
- Parse JSON (no comments)
- Return `*SessionData` or `nil`

### Phase 5e: Get Git Context
- Implement `getGitContext(workspace)`
- Run `git rev-parse --abbrev-ref HEAD` for branch
- Run `git status --porcelain | wc -l` for uncommitted count
- Run `git log -1 --format="%ar|%s"` for last commit
- Return `*GitContext` or `nil`

### Phase 6: Build Context Sections
- Implement all `build*Section()` functions
- Each generates markdown string from corresponding data
- Gracefully handles `nil` data (skip section if unavailable)

### Phase 7: Template Alignment
- Full METADATA block with complete documentation
- SETUP with all imports, types, package-level state
- BODY with organizational chart, helpers, core ops, public APIs
- CLOSING with all 13 sections

---

## Fallback Behavior

**Principle:** Always output valid context even if data incomplete

| Missing Data | Fallback Behavior |
|--------------|-------------------|
| User config | Skip user awareness section, continue |
| Instance config | Skip identity section, use minimal hardcoded communication guide |
| Session data | Skip session context section |
| Git context | Skip work context section |
| Temporal data | Skip temporal section (already handled by temporal library) |

**Result:** Context always valid, just less complete when data unavailable

---

## Benefits

### For Instance (Nova Dawn)
- Complete identity grounding at session start
- Awareness of who Seanje is (not generic user)
- Continuity across sessions (session history)
- Work context awareness (what workspace, what branch)
- Proper temporal grounding (time/schedule)

### For System
- Config-driven identity (not hardcoded)
- All CPI-SI infrastructure actually used
- Consistent session bootstrapping
- Foundation for future enhancements (journals, patterns, health)

### For Future
- Easy to add new context sections
- Data sources already structured
- Fallback ensures backward compatibility
- Scales to richer context as system grows

---

## Migration Path

**v1.0.0 → v2.0.0:**
1. Keep existing `OutputClaudeContext()` signature (no breaking changes for hooks)
2. Refactor internal implementation to load from configs
3. If configs unavailable, fall back to minimal hardcoded guide
4. Hooks continue working even if config infrastructure incomplete

**Testing:**
1. Test with full config infrastructure present
2. Test with configs missing (fallback behavior)
3. Test with partial data (some configs exist, others don't)
4. Verify JSON output format unchanged for Claude Code parsing

---

## Success Criteria

- ✅ Loads user config and integrates into context
- ✅ Loads instance config and integrates into context
- ✅ Loads session data and integrates into context
- ✅ Loads git context and integrates into context
- ✅ Temporal integration continues working
- ✅ Falls back gracefully when data unavailable
- ✅ Output JSON format compatible with Claude Code
- ✅ Full template alignment (METADATA/SETUP/BODY/CLOSING)
- ✅ API documentation created
- ✅ Compilation successful with no warnings

---

**Next Steps:** Proceed with Phase 5b (load user config implementation)
