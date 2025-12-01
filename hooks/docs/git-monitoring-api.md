# Git Repository Monitoring Library - API Documentation

**Component:** hooks/lib/session/git.go
**Version:** 2.0.0
**Created:** November 12, 2025
**Biblical Foundation:** Proverbs 27:23 (WEB) - "Know well the state of your flocks, and pay attention to your herds"

---

## Table of Contents

1. [Overview & Philosophy](#overview--philosophy)
2. [Design Rationale](#design-rationale)
3. [Public API](#public-api)
4. [Configuration System](#configuration-system)
5. [Integration Patterns](#integration-patterns)
6. [Modification Policy](#modification-policy)
7. [Troubleshooting](#troubleshooting)
8. [Future Roadmap](#future-roadmap)

---

## Overview & Philosophy

### What This Library Does

The Git Repository Monitoring library provides **configurable, non-intrusive status reporting** for workspace git repositories during Claude Code session start. It integrates with `system/lib/git` for repository information and `system/lib/display` for consistent formatted output.

### Biblical Foundation

> **"Know well the state of your flocks, and pay attention to your herds"** - Proverbs 27:23 (WEB)

This library embodies the biblical principle of **faithful stewardship through awareness**. Just as a shepherd knows the condition of their flock, a developer should know the state of their workspace. The library provides:

- **Awareness without anxiety** - Know what's uncommitted, ahead, behind
- **Non-intrusive monitoring** - Shows when needed, silent when clean
- **Configurable priorities** - Focus on what matters to your workflow
- **Faithful reporting** - Accurate, actionable information

### Core Design Principles

| Principle | Implementation | Why It Matters |
|-----------|----------------|----------------|
| **Non-blocking** | Never fails session start | Monitoring should enhance, never hinder |
| **Configurable** | JSONC config file with defaults | Users have different workflows and priorities |
| **Graceful degradation** | Works without config file | Always functional, even in minimal environment |
| **Privacy-respecting** | Only counts, not content | Respects workspace privacy |
| **Integration-focused** | Uses system libraries | Consistent formatting, no duplication |

---

## Design Rationale

### Why This Library Exists

**Problem:** Developers often start sessions unaware of workspace state:
- Uncommitted changes forgotten
- Local commits not pushed to remote
- Merge conflicts unresolved
- Stashes accumulating

**Solution:** **Gentle, configurable awareness** at session start showing actionable repository status.

### Architecture Decisions

#### 1. Library (Not Command)

**Decision:** Implemented as `hooks/lib/session` library, not standalone command.

**Rationale:**
- Session hooks orchestrate multiple awareness systems (git status, temporal context, etc.)
- Library provides reusable function other hooks can call
- Follows ladder pattern: hooks orchestrate, libraries provide functionality

**Alternative considered:** Standalone `git-status` command
**Why rejected:** Session start already runs multiple checks - extracting git logic into library allows composition

#### 2. Configuration-Driven

**Decision:** All messages, display options, and check priorities configurable via JSONC file.

**Rationale:**
- Different workflows have different priorities (some care about stashes, others don't)
- Message templates allow personalization/internationalization
- Configuration separated from code allows updates without recompilation
- JSONC format supports comments for self-documenting config

**Default location:** `~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc`

**Fallback:** Hardcoded defaults ensure functionality without config file.

#### 3. Integration with system/lib/git

**Decision:** Delegate repository queries to `system/lib/git`, don't execute git commands directly.

**Rationale:**
- Single source of truth for git operations
- Consistent error handling across all hooks
- Performance optimization (git library can cache results)
- Follows dependency ladder (libraries depend on lower-level libraries)

#### 4. Integration with system/lib/display

**Decision:** Use `display.Success()` formatters instead of raw `fmt.Println`.

**Rationale:**
- Consistent formatting across all hooks/libraries
- ANSI color support with graceful degradation
- Health scoring and logging built-in
- Reduces code duplication

---

## Public API

### CheckGitStatus(workspace string)

**Purpose:** Examine git repository state and report configured status checks.

**Signature:**
```go
func CheckGitStatus(workspace string)
```

**Parameters:**
- `workspace` (string): Directory path to check for git repository status

**Returns:**
- None (prints to stdout using display library)

**Behavior:**
1. Skip if monitoring disabled (`behavior.enabled = false`)
2. Skip if workspace is not a git repository
3. Gather repository information using `system/lib/git`
4. Check each enabled status check:
   - Uncommitted changes
   - Commits ahead of remote
   - Commits behind remote
   - Stashes present
   - Merge conflicts
5. Display results:
   - If issues found: Show header + bulleted list
   - If clean + `show_when_clean`: Show success message
   - If clean + not configured: Silent (no output)

**Example Usage:**

```go
package main

import "hooks/lib/session"

func main() {
    session.CheckGitStatus("/path/to/workspace")
}
```

**Example Output (Issues Detected):**

```
📝 Workspace Git Status:
   • 3 uncommitted change(s)
   • 2 commit(s) ahead of remote
```

**Example Output (Clean Repository, show_when_clean = true):**

```
📝 Workspace Git Status:
   ✓ Repository is clean
```

**Example Output (Clean Repository, show_when_clean = false):**

```
(no output - silent success)
```

**Health Impact:**
- Repository clean: +15 points (good workspace state)
- Issues detected: +10 points (successful monitoring)
- Check failed: -10 points (system error, e.g. git not installed)

---

## Configuration System

### Configuration File Structure

**Location:** `~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc`

**Format:** JSONC (JSON with comments)

```jsonc
{
  "metadata": {
    "name": "Git Repository Monitoring Configuration",
    "description": "Controls git status checks at session start",
    "version": "1.0.0",
    "author": "Seanje Lenox-Wise",
    "created": "2025-11-12",
    "last_updated": "2025-11-12"
  },
  "display": {
    "header_icon": "📝",
    "header_text": "Workspace Git Status",
    "show_when_clean": false,  // Set true to show message when repo clean
    "clean_message": "Repository is clean"
  },
  "checks": {
    "uncommitted_changes": true,  // Check for unstaged/staged changes
    "ahead_of_remote": true,      // Check for unpushed commits
    "behind_remote": true,        // Check if behind remote
    "stashes": true,              // Check for stashed changes
    "conflicts": true             // Check for merge conflicts
  },
  "messages": {
    "uncommitted_changes": "{count} uncommitted change(s)",
    "ahead_of_remote": "{count} commit(s) ahead of remote",
    "behind_remote": "{count} commit(s) behind remote",
    "stashes": "{count} stash(es)",
    "conflicts": "{count} merge conflict(s)"
  },
  "behavior": {
    "enabled": true,                // Master switch - set false to disable completely
    "check_on_session_start": true  // Future: other trigger points
  }
}
```

### Configuration Options Explained

#### Display Settings

| Option | Type | Default | Purpose |
|--------|------|---------|---------|
| `header_icon` | string | "📝" | Icon before header text |
| `header_text` | string | "Workspace Git Status" | Header displayed when issues found |
| `show_when_clean` | boolean | false | Show success message when repo clean |
| `clean_message` | string | "Repository is clean" | Message shown when clean (if enabled) |

#### Check Settings

Each check can be individually enabled/disabled:

| Check | What It Detects | When To Enable | When To Disable |
|-------|-----------------|----------------|-----------------|
| `uncommitted_changes` | Staged or unstaged files | Always (core awareness) | Never (fundamental) |
| `ahead_of_remote` | Local commits not pushed | Important for collaboration | Solo offline work |
| `behind_remote` | Remote commits not pulled | Team environments | Solo work, infrequent pulls |
| `stashes` | Work stashed away | If you use stashes | If you don't stash |
| `conflicts` | Merge conflicts present | Always (critical issue) | Never (must resolve) |

#### Message Templates

All messages support `{count}` placeholder for dynamic values:

```jsonc
"uncommitted_changes": "{count} uncommitted change(s)"
// Output: "3 uncommitted change(s)"

"ahead_of_remote": "You have {count} unpushed commits"
// Output: "You have 2 unpushed commits"
```

**Future expansion:** Additional placeholders planned (`{branch}`, `{files}`, `{remote}`)

#### Behavior Settings

| Option | Purpose | Values |
|--------|---------|--------|
| `enabled` | Master on/off switch | `true` = monitoring active, `false` = completely disabled |
| `check_on_session_start` | When to check | Currently always true (future: other triggers) |

### Customization Examples

#### Minimal Awareness (Only Critical Issues)

```jsonc
{
  "checks": {
    "uncommitted_changes": true,  // Must know
    "ahead_of_remote": false,     // Don't care about pushing
    "behind_remote": false,       // Pull manually when needed
    "stashes": false,             // Don't use stashes
    "conflicts": true             // Must know
  },
  "display": {
    "show_when_clean": false      // Silent when nothing wrong
  }
}
```

#### Verbose Awareness (Show Everything)

```jsonc
{
  "checks": {
    "uncommitted_changes": true,
    "ahead_of_remote": true,
    "behind_remote": true,
    "stashes": true,
    "conflicts": true
  },
  "display": {
    "show_when_clean": true,      // Always show status
    "clean_message": "✓ Workspace is clean and synced"
  }
}
```

#### Custom Messages (Personalized)

```jsonc
{
  "messages": {
    "uncommitted_changes": "⚠️  You forgot to commit {count} file(s)",
    "ahead_of_remote": "📤 Push your {count} commit(s) when you get a chance",
    "stashes": "💾 You have {count} stash(es) - might want to review those"
  }
}
```

---

## Integration Patterns

### Session Start Hook Integration

**Primary use case:** Called by `session/cmd-start/start.go` during session initialization.

```go
// In session start hook BODY
session.CheckGitStatus(workspace)
```

**Integration point in execution flow:**
1. Environment setup
2. Workspace analysis
3. **→ Git status check (this library)**
4. Temporal awareness
5. Identity grounding
6. Session initialization

### Custom Hook Integration

Any hook can call `CheckGitStatus()` to display current repository state:

```go
import "hooks/lib/session"

func customHook() {
    // ... other logic ...

    // Show current git status
    session.CheckGitStatus("/path/to/workspace")

    // ... continue ...
}
```

### Programmatic Configuration Override

For hooks that need dynamic behavior:

```go
// Temporarily disable monitoring
// (Note: Currently requires modifying config file - programmatic API planned for v3.0.0)
```

---

## Modification Policy

### Safe To Modify (Extension Points)

✅ **Add new git status checks:**
1. Add check field to `GitChecksConfig` type in SETUP
2. Add message template to `GitMessagesConfig` type in SETUP
3. Add check logic to `CheckGitStatus()` in BODY
4. Add default config values in `getDefaultGitConfig()`
5. Update configuration file schema documentation

✅ **Extend message template system:**
1. Add new placeholder to `formatGitMessage()` function
2. Update message template documentation
3. Test with various message templates

✅ **Add new display options:**
1. Add field to `GitDisplayConfig` type in SETUP
2. Update `getDefaultGitConfig()` with default value
3. Apply option in `CheckGitStatus()` display logic

### Modify With Extreme Care (Breaking Changes)

⚠️ **CheckGitStatus() signature** - Breaks all calling hooks if changed
⚠️ **Configuration structure** - Breaks existing config files if fields renamed/removed
⚠️ **Display output format** - May break hooks that parse output (none currently)

**Migration strategy if breaking changes needed:**
1. Deprecate old API, maintain for 2 versions
2. Add new API alongside old
3. Update all calling code to use new API
4. Remove old API in v4.0.0

### Never Modify (Foundational)

❌ **4-block structure** - METADATA, SETUP, BODY, CLOSING must remain
❌ **Non-blocking guarantee** - Monitoring must never block session start
❌ **Configuration fallback** - Must work without config file
❌ **Rails pattern** - Config loaded in init(), available throughout component

---

## Troubleshooting

### Common Issues

#### 1. No Output When Repository Has Changes

**Symptoms:**
- Repository has uncommitted files
- `git status` shows changes
- Session start shows no git status output

**Diagnosis:**
1. Check `behavior.enabled` in config (must be `true`)
2. Verify workspace path is a valid git repository
3. Check relevant checks enabled (e.g., `checks.uncommitted_changes = true`)
4. Verify config file exists and is readable

**Solution:**
- Edit config file: Set `behavior.enabled = true`
- Enable desired checks: Set `checks.uncommitted_changes = true`
- Test: `cat ~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc`

#### 2. Configuration Not Loading

**Symptoms:**
- Config file exists but default messages appear
- Changes to config file have no effect

**Diagnosis:**
1. Verify file location: `~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc`
2. Check file permissions (must be readable)
3. Validate JSONC syntax (comments must be `//` or `/* */`)
4. Check for trailing commas (valid in JSONC, but test parser)

**Solution:**
- Library falls back to defaults if config fails to load
- This is expected behavior, not an error
- Check syntax: Remove config file temporarily, see if defaults work
- If defaults work but config doesn't: Syntax error in config file

**Debug commands:**
```bash
# View config file
cat ~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc

# Test JSONC syntax (strip comments, validate JSON)
grep -v '^[[:space:]]*//' ~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc | jq .
```

#### 3. Message Placeholders Not Replaced

**Symptoms:**
- Messages show `{count}` literally instead of numbers
- Output: "You have {count} uncommitted change(s)"

**Diagnosis:**
- Only `{count}` placeholder supported currently
- Check message templates in configuration
- Verify `formatGitMessage()` is being called

**Solution:**
- Use `{count}` format in message templates
- Other placeholders planned for future versions
- Current limitation documented in roadmap

#### 4. Display Formatting Issues

**Symptoms:**
- No colors in output
- Emoji icons not displaying
- Formatting looks broken

**Diagnosis:**
- Terminal may not support ANSI colors
- Emoji support varies by terminal/font
- Display library handles graceful degradation

**Solution:**
- Display library automatically falls back to plain text
- No action needed - designed to work in all terminals
- If icons matter: Use terminal with emoji support

---

## Future Roadmap

### Planned Features

#### v2.1.0 - Enhanced Message Templates
- **Additional placeholders:** `{branch}`, `{files}`, `{remote}`
- **Conditional messages:** Different messages based on count thresholds
- **Example:** "1 file changed" vs "5 files changed"

#### v2.2.0 - Threshold-Based Warnings
- **Threshold configuration:** Only warn if count > threshold
- **Use case:** Don't show "1 uncommitted change" during active development
- **Example config:**
  ```jsonc
  "thresholds": {
    "uncommitted_changes": 5,  // Only warn if >5 files
    "ahead_of_remote": 3       // Only warn if >3 commits
  }
  ```

#### v2.3.0 - Branch Pattern Warnings
- **Branch-specific warnings:** Warn if working on protected branches
- **Use case:** Alert if committing directly to `main` or `master`
- **Example config:**
  ```jsonc
  "branch_warnings": {
    "protected_branches": ["main", "master", "production"],
    "warning_message": "⚠️  You're working on {branch} - consider using feature branch"
  }
  ```

#### v3.0.0 - Detailed File Listing
- **File listing option:** Show which files are uncommitted
- **Configurable detail level:** Count only, file paths, or full diff
- **Example output:**
  ```
  📝 Workspace Git Status:
     • 3 uncommitted changes:
       - src/components/Header.tsx (modified)
       - tests/Header.test.tsx (new file)
       - package.json (modified)
  ```

#### v3.1.0 - Untracked Files Check
- **New check:** Detect untracked files
- **Configurable:** Separate from uncommitted_changes
- **Use case:** Know if new files aren't being tracked

#### v3.2.0 - Detached HEAD Warning
- **Critical check:** Detect if HEAD is detached
- **Auto-enable:** Always enabled (critical state to know)
- **Warning level:** Higher priority than other checks

### Research Areas

#### CI/CD Status Integration
- Display build status from CI/CD systems
- Integration with GitHub Actions, GitLab CI
- Show: ✓ Build passing, ✗ Build failing, ⏳ Build running

#### Issue Tracker Integration
- Show related issues for current branch
- GitHub Issues, Jira integration
- Display: "Working on #123: Add user authentication"

#### Git Workflow Recommendations
- Suggest actions based on status
- Example: "3 commits ahead - run `git push`"
- Contextual, non-intrusive suggestions

#### Automatic Stash Suggestions
- Detect work that should be stashed
- Suggest stashing before branch switching
- Integration with `git stash` workflow

### Version History

**2.0.0** (2025-11-12) - Template Alignment & Configuration System
- Complete 8-rung METADATA block with Biblical Foundation
- Configuration system (git-monitoring.jsonc)
- Display library integration for consistent formatting
- Message templates with `{count}` placeholder support
- Show-when-clean option
- Individual check enable/disable controls
- Graceful fallback to defaults
- Refactored from hardcoded messages to configurable templates
- No breaking changes (maintains backward compatibility)

**1.0.0** (2024-10-24) - Initial Implementation
- Basic git status checking (uncommitted, ahead/behind, stashes, conflicts)
- Hardcoded display format and messages
- Direct fmt.Printf output

---

## Quick Reference

### Basic Usage

```go
import "hooks/lib/session"

func main() {
    session.CheckGitStatus("/path/to/workspace")
}
```

### Configuration File Location

```
~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc
```

### Minimal Config (Just Enable/Disable)

```jsonc
{
  "behavior": {
    "enabled": true
  }
}
```

### Most Common Customization

```jsonc
{
  "display": {
    "show_when_clean": true,
    "clean_message": "✓ Workspace is clean"
  }
}
```

### Disable Specific Checks

```jsonc
{
  "checks": {
    "uncommitted_changes": true,
    "ahead_of_remote": true,
    "behind_remote": false,    // Disable this check
    "stashes": false,          // Disable this check
    "conflicts": true
  }
}
```

---

*"Commit your work to the LORD, and your plans will be established" - Proverbs 16:3 (ESV)*

**This library helps you know your workspace state - faithful monitoring for faithful work.**
