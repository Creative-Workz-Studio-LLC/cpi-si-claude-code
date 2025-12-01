# Temporary Data - Bootstrap Capture Space

**Purpose:** Manual data capture while building automation. Real data in temporary format awaiting proper home.

## Core Principle

**The work actually happened.** Sessions are real. Learning is valid. Patterns existed even without automated capture. This folder holds permanent data in temporary format until proper systems exist to store it correctly.

## What Belongs Here

**manual-captures/**

- Manual recordings while automation is being built
- Session notes before session-time tools exist
- Activity logs before proper activity tracking
- Pattern observations before pattern recognition system
- Any data we need to capture NOW but don't have proper system YET

**migration-notes/**

- Notes on what's being automated
- When automation is expected to be ready
- How to migrate manual data to permanent structure
- Validation criteria for migration

## Lifecycle

```bash
1. Identify what needs capturing
   ↓
2. Manually record in temp/manual-captures/
   ↓
3. Build proper automation informed by manual capture
   ↓
4. Validate automation works correctly
   ↓
5. Migrate manual data → proper permanent structure
   ↓
6. Archive/delete temp files (data preserved in permanent home)
```

**Key Insight:** We don't throw away bootstrap data - we MIGRATE it. The sessions happened, the learning is real, the patterns are valid. Once proper structure exists, manual data becomes part of permanent record.

## Naming Convention

**manual-captures/YYYY-MM-DD_type_description.{txt,json,md}**

Examples:

- `2025-11-06_sleep-tracking.txt`
- `2025-11-05_session-notes.md`
- `2025-11-07_manual-activity-log.json`

**migration-notes/YYYY-MM-DD_what-being-automated.md**

Examples:

- `2025-11-07_session-time-automation.md`
- `2025-11-08_activity-tracking-system.md`

## Why This Matters

**Pattern Validation**

- "Nova worked build-first in October" is a data point, even if manually recorded
- More data points = better pattern recognition

**Historical Continuity**

- Complete record from day one, not "started tracking [later date]"
- Memory system needs ALL sessions, not just post-automation

**Learning Accuracy**

- Manual capture informs automation design
- Bootstrap data validates automated capture works correctly

**Identity Confirmation**

- "I've been doing this for X months" (full history)
- Not "system only tracked 2 weeks" (incomplete picture)

## Current Contents

**manual-captures/**

- `2025-11-06_sleep-tracking.txt` - Sleep periods during session while building session-time tools

**migration-notes/**

- (None yet - will be created as we automate systems)

## Migration Examples

**Example 1: Sleep Tracking → Session Data**

```bash
Manual: 2025-11-06_sleep-tracking.txt
   ↓ session-time tools built
Migrate to: session/history/2025-11-06_session.json
   (manual sleep data becomes semi-downtime periods)
```

**Example 2: Manual Session Notes → Activity Log**

```bash
Manual: 2025-11-05_session-notes.md
   ↓ activity tracking built
Migrate to: session/current-log.json or activity.jsonl
   (notes become structured activity entries)
```

## Cleanup Policy

**DO NOT delete manual-captures/ files until:**

1. Proper automation exists
2. Automation has been validated
3. Manual data has been successfully migrated
4. Migration has been verified complete

**After successful migration:**

- Archive temp files (optional backup)
- Delete from temp/ (data preserved in permanent location)

## Development Notes

This folder supports build-first approach:

1. Capture data NOW (even manually)
2. Understand the pattern through captured data
3. Build automation based on real examples
4. Migrate bootstrap data into permanent system

Never lose learning. Never discard real history. Bootstrap manually, automate later, preserve everything.

---

*This folder is NOT a junk drawer. It's a staging area for permanent data awaiting proper structure.*
