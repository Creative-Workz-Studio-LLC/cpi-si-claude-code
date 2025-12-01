# Manual Captures

**Purpose:** Real data captured manually while automation is being built. Permanent data in temporary format.

## What Goes Here

Any data we need to record NOW but don't have proper automated system YET:

- Session notes before session-time tools
- Activity logs before activity tracking
- Pattern observations before pattern recognition
- Sleep tracking before downtime detection
- Stopping reasons before automated capture
- Work context before config inheritance
- Any learning that's happening in real-time

## Naming Convention

**Format:** `YYYY-MM-DD_type_description.{txt,json,md}`

**Examples:**

- `2025-11-06_sleep-tracking.txt`
- `2025-11-05_session-notes.md`
- `2025-11-07_manual-activity-log.json`
- `2025-11-08_pattern-observations.md`
- `2025-11-10_stopping-reasons.txt`

## File Format Guidelines

**Use whatever format captures the data effectively:**

- `.txt` - Quick notes, timestamps, simple logs
- `.json` - Structured data that will map to automated schemas
- `.md` - Narrative observations, longer-form notes

**Priority is capturing the data, not perfect format.** Automation will come later, informed by these captures.

## Current Manual Captures

- `2025-11-06_sleep-tracking.txt` - Sleep periods during session while building session-time tools

## Remember

**The work actually happened.** These captures are real data that will be migrated to permanent structure once automation exists. Don't skip recording because "we don't have the system yet" - that's exactly what this folder is for.

---

*Manual capture → Automation built → Data migrated → Learning preserved*
